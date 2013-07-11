// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"github.com/sdming/wk/fsw"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)

// configViewEngine
func (srv *HttpServer) configViewEngine() error {
	if !srv.Config.ViewEnable || srv.Config.ViewDir == "" {
		return nil
	}

	ve, err := NewGoHtml(srv.Config.ViewDir)
	if err != nil {
		return err
	}
	if err = ve.Register(srv); err != nil {
		return err
	}
	srv.ViewEngine = ve
	Logger.Println("ViewEngine:", ve)
	return nil

}

// func Register(name string, engine ViewEngine) {
// 	ViewEngines[name] = engine
// }

// ViewEngine is wrap of executing template file
type ViewEngine interface {
	Register(server *HttpServer) error
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

	sync.Mutex
	watcher    *fsw.FsWatcher
	dependence map[string][]string
}

// NewGoHtml return a *GoHtml, it retur error if basePath doesn't exists
func NewGoHtml(basePath string) (*GoHtml, error) {
	if basePath == "" {
		return nil, errors.New("bash path cann't be empty")
	}

	if basePath != "" {
		basePath = cleanFilePath(basePath)
		if !strings.HasSuffix(basePath, "/") {
			basePath = basePath + "/"
		}
	}

	if !isDirExists(basePath) {
		return nil, errors.New("path doesn't exists " + basePath)
	}

	ve := &GoHtml{
		BasePath:       basePath,
		TemplatesCache: make(map[string]*template.Template),
		Funcs:          make(template.FuncMap),
		EnableCache:    false,
		dependence:     make(map[string][]string, 0),
	}

	ve.dependence["a"] = make([]string, 0)

	for name, fn := range TemplateFuncs {
		ve.Funcs[name] = fn
	}
	ve.Funcs["partial"] = ve.renderfile // render a template file

	return ve, nil
}

// Register initialize viewengine
func (ve *GoHtml) Register(server *HttpServer) error {
	if server.Config.ViewDir == "" {
		return nil
	}
	basePath := cleanFilePath(server.Config.ViewDir)
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}
	ve.BasePath = basePath

	ve.TemplatesCache = make(map[string]*template.Template)
	ve.Funcs = make(template.FuncMap)

	for name, fn := range TemplateFuncs {
		ve.Funcs[name] = fn
	}
	ve.Funcs["partial"] = ve.renderfile // render a template file

	if conf, ok := server.Config.PluginConfig.Child("gohtml"); ok {
		ve.EnableCache = conf.ChildBoolOrDefault("cache_enable", false)
		if ve.EnableCache {
			if fw, err := fsw.NewFsWatcher(ve.BasePath); err == nil {
				ve.watcher = fw
				ve.watcher.Listen(ve.notify)
			}
		}

		Logger.Println("gohtml", "cache_enable:", ve.EnableCache, "watcher", ve.watcher)
	}

	return nil
}

func (ve *GoHtml) notify(e fsw.Event) {
	ve.Lock()
	defer ve.Unlock()

	key := strings.ToLower(cleanFilePath(e.Name))
	delete(ve.TemplatesCache, key)

	if d, ok := ve.dependence[key]; ok && len(d) > 0 {
		for i := 0; i < len(d); i++ {
			delete(ve.TemplatesCache, d[i])
		}
	}
	delete(ve.dependence, key)

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

func (ve *GoHtml) getCache(file string) (t *template.Template, ok bool) {
	ve.Lock()
	defer ve.Unlock()

	key := strings.ToLower(cleanFilePath(path.Join(ve.BasePath, file)))
	t, ok = ve.TemplatesCache[key]
	return t, ok
}

func (ve *GoHtml) setCache(file string, t *template.Template) {
	ve.Lock()
	defer ve.Unlock()

	key := strings.ToLower(cleanFilePath(path.Join(ve.BasePath, file)))
	ve.TemplatesCache[key] = t
}

func (ve *GoHtml) setDepend(partial, template string) {
	ve.Lock()
	defer ve.Unlock()

	key := strings.ToLower(partial)
	d, ok := ve.dependence[key]

	if ok && d != nil {
		d = append(d, strings.ToLower(template))
		ve.dependence[key] = d
	} else {
		d = make([]string, 1, 31)
		d[0] = strings.ToLower(template)
		ve.dependence[key] = d
	}
}

// findTemplate
func (ve *GoHtml) findTemplate(file string) (t *template.Template, err error) {
	if ve.EnableCache {
		var ok bool

		if t, ok = ve.getCache(file); !ok {
			t, err = ve.parse(file)
			if err == nil {
				ve.setCache(file, t)
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
		return nil, errors.New("file can't be empty")
	}

	file = cleanFilePath(path.Join(ve.BasePath, file))
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
				fileToImport = path.Join(ve.BasePath, fileToImport)
				data.Write(ve.readImportFile(fileToImport))
				ve.setDepend(fileToImport, file)
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
	TemplateFuncs["import"] = importFile // import a temlate file, must be separate line
	TemplateFuncs["fv"] = formValue      // call *http.Request.FormValue
	TemplateFuncs["incl"] = include      // 
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

func formValue(req *http.Request, name string) string {
	if req == nil || name == "" {
		return ""
	}
	return req.FormValue(name)
}

func importFile(files ...string) template.HTML {
	return template.HTML("")
}

func include(values []string, v string) bool {
	if values == nil || len(values) == 0 {
		return false
	}
	for _, s := range values {
		if v == s {
			return true
		}
	}
	return false
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
