// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// DefaultViewEngine
// TODO: global? move to httpserver?
var DefaultViewEngine ViewEngine

// configViewEngine
// TODO: remove DefaultViewEngine?
func (srv *HttpServer) configViewEngine() error {
	ve, err := NewGoHtml(srv.Config.ViewDir)
	if err != nil {
		return err
	}

	DefaultViewEngine = ve
	srv.ViewEngine = ve
	Logger.Printf("ViewEngine: %v; \n", ve)
	return nil
}

// func Register(name string, engine ViewEngine) {
// 	ViewEngines[name] = engine
// }

// ViewEngine is wrap of executing template file
type ViewEngine interface {
	Register(server *HttpServer)
	Execte(writer io.Writer, file string, data interface{}) error
}

var (
	importre *regexp.Regexp = regexp.MustCompile(`^[\s]*{{[\s]*import[\s]+"(?P<file>[\S]*)"[\s]*}}[\s]*$`)
)

// GoHtml is a ViewEngine that wrap "html/template"
type GoHtml struct {
	BasePath       string
	EnableCache    bool
	TemplatesCache map[string]*template.Template
	Funcs          template.FuncMap
}

// NewGoHtml return a *GoHtml, it retur error if basePath doesn't exists
func NewGoHtml(basePath string) (*GoHtml, error) {
	// if basePath == "" {
	// 	return nil, errors.New("bash path cann't be empty")
	// }

	if basePath != "" {
		basePath = cleanFilePath(basePath)
		if !strings.HasSuffix(basePath, "/") {
			basePath = basePath + "/"
		}
	}

	// if !isDirExists(basePath) {
	// 	return nil, errors.New("path doesn't exists " + basePath)
	// }

	ve := &GoHtml{
		BasePath:       basePath,
		TemplatesCache: make(map[string]*template.Template),
		Funcs:          make(template.FuncMap),
	}

	for name, fn := range TemplateFuncs {
		ve.Funcs[name] = fn
	}
	ve.Funcs["partial"] = ve.renderfile // render a template file

	return ve, nil
}

// Register initialize viewengine
func (ve *GoHtml) Register(server *HttpServer) {
	if server.Config.ViewDir == "" {
		return
	}
	basePath := cleanFilePath(server.Config.ViewDir)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	ve.BasePath = basePath
}

// String
func (ve *GoHtml) String() string {
	if ve == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GoHtml: BasePath=%s ", ve.BasePath)
}

// Execte execute template
func (ve *GoHtml) Execte(wr io.Writer, file string, data interface{}) error {
	t, err := ve.findTemplate(file)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}

// findTemplate
func (ve *GoHtml) findTemplate(file string) (t *template.Template, err error) {
	if ve.EnableCache {
		file = strings.ToLower(file)
		t = ve.TemplatesCache[file]
		if t == nil {
			t, err = ve.parse(file)
			if err == nil {
				ve.TemplatesCache[file] = t
			}
		}
	} else {
		t, err = ve.parse(file)
	}
	return
}

func (ve *GoHtml) readImportFile(file string) []byte {
	if file == "" {
		return nil
	}
	file = path.Join(ve.BasePath, file)
	if !isFileExists(file) {
		return nil
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return b
}

// parse
func (ve *GoHtml) parse(file string) (*template.Template, error) {
	if file == "" {
		return nil, errors.New("file cann't be empty")
	}

	file = path.Join(ve.BasePath, file)
	file = cleanFilePath(file)
	if !isFileExists(file) {
		return nil, errors.New("file doesn't exists " + file)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(b)
	var data bytes.Buffer

	for {
		var line []byte
		end := false

		if line, err = buff.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				return nil, err
			} else {
				end = true
			}
		}

		if importre.Match(line) {
			fileToImport := importre.FindStringSubmatch(string(line))[1]
			if fileToImport != "" {
				data.Write(ve.readImportFile(fileToImport))
			}
			continue
		}

		if _, err = data.Write(line); err != nil {
			return nil, err
		}

		if end {
			break
		}
	}

	t := template.New(file)
	//t.Funcs(TemplateFuncs)
	t.Funcs(ve.Funcs)

	return t.Parse(string(data.Bytes()))
}

var (
	TemplateFuncs template.FuncMap = make(map[string]interface{})
	//ViewEngines   map[string]ViewEngine = make(map[string]ViewEngine)
)

func init() {
	initFuncs()
}

func initFuncs() {
	TemplateFuncs["eq"] = equal          // equal
	TemplateFuncs["eqs"] = equalAsString // convert to string and compare
	TemplateFuncs["gt"] = greater        // greater
	TemplateFuncs["le"] = less           // less

	TemplateFuncs["set"] = setmap        // set map[string]interface{}
	TemplateFuncs["raw"] = raw           // unescaped html
	TemplateFuncs["selected"] = selected // output "selected" or ""
	TemplateFuncs["checked"] = checked   // output "checked" or ""
	TemplateFuncs["nl2br"] = nl2br       // replace \n to <br/>
	TemplateFuncs["jsvar"] = jsvar       // convert data to javascript variable, like var name = {...}
	TemplateFuncs["import"] = importfile // import a temlate file, must be separate line
	//TemplateFuncs["render"] = renderfile //
}

func (ve *GoHtml) renderfile(file string, data interface{}) template.HTML {

	t, err := ve.findTemplate(file)
	if err != nil {
		return template.HTML(err.Error())
	}
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, data)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(buffer.String())
}

// func renderfile(file string, data interface{}) template.HTML {
// 	fmt.Println("render partial", file, data)
// 	return template.HTML("TODO:render")
// }

func importfile(files ...string) template.HTML {
	return template.HTML("")
}

func raw(text string) template.HTML {
	return template.HTML(text)
}

func equalAsString(a, b interface{}) bool {
	if a == b {
		return true
	}

	return fmt.Sprint(a) == fmt.Sprint(b)
}

func equal(a, b interface{}) bool {
	return gotype.Equal(a, b)
}

func greater(a, b interface{}) bool {
	return gotype.Greater(a, b)
}

func less(a, b interface{}) bool {
	return gotype.Less(a, b)
}

func selected(selected bool) template.HTMLAttr {
	if selected {
		return "selected"
	}
	return ""
}

func checked(checked bool) template.HTMLAttr {
	if checked {
		return "checked"
	}
	return ""
}

func nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br/>", -1))
}

func setmap(data map[string]interface{}, name string, value interface{}) template.HTML {
	data[name] = value
	return template.HTML("")
}

func jsvar(name string, data interface{}) template.JS {
	buffer := new(bytes.Buffer)
	buffer.WriteString(" var ")
	buffer.WriteString(name)
	buffer.WriteString(" = ")
	encoder := json.NewEncoder(buffer)
	encoder.Encode(data)
	buffer.WriteString(";")
	return template.JS(buffer.String())
}