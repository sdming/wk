// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"fmt"
	"github.com/sdming/mcache"
	"github.com/sdming/wk/fsw"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// StaticProcessor handle request of static file
type StaticProcessor struct {
	server      *HttpServer
	fileCache   *fileCache
	cacheEnable bool
	cacheExpire int
	header      map[string]string
}

// newStaticProcessor return default  *StaticProcessor
func newStaticProcessor() *StaticProcessor {
	return &StaticProcessor{}
}

// Register
func (p *StaticProcessor) Register(server *HttpServer) {
	p.server = server

	if conf, ok := server.Config.PluginConfig.Child("static_processor"); ok {
		p.cacheEnable = conf.ChildBoolOrDefault("cache_enable", false)
		if p.cacheEnable {
			p.cacheExpire = int(conf.ChildIntOrDefault("cache_expire", 86400))
			if fc, err := newFileCache(server.Config.PublicDir); err != nil {
				Logger.Println("static processor new file cache fail", err)
			} else {
				p.fileCache = fc
			}
		}

		if headNode, ok := conf.Child("header"); ok {
			p.header, _ = headNode.Map()
		}

		Logger.Println("satic processor", "cache_enable:", p.cacheEnable, "cache_expire:", p.cacheExpire, "header:", p.header)
	}
}

// setHeader set customer header ti response
func (p *StaticProcessor) setHeader(ctx *HttpContext) {
	if p.header == nil {
		return
	}

	for k, v := range p.header {
		if ctx.ResHeader(k) == "" {
			ctx.SetHeader(k, v)
		}
	}
}

// Execute set FileResult if request file does exist
func (p *StaticProcessor) Execute(ctx *HttpContext) {
	if ctx.Result != nil || ctx.Error != nil {
		return
	}

	physicalPath := path.Join(p.server.Config.PublicDir, ctx.RequestPath)
	if p.cacheEnable && p.fileCache != nil {
		if f, ok := p.fileCache.get(cleanFilePath(physicalPath)); ok {
			ctx.PhysicalPath = physicalPath
			ctx.Result = f
			p.setHeader(ctx)
			return
		}
	}

	info, err := os.Stat(physicalPath)
	if err != nil {
		return
	}

	if (info.IsDir() && p.server.Config.IndexesEnable && ctx.RequestPath != "/") || !info.IsDir() {
		ctx.PhysicalPath = physicalPath
		ctx.Result = File(physicalPath)
		p.setHeader(ctx)

		if !info.IsDir() {
			p.trySetFileCache(cleanFilePath(physicalPath))
		}
		return
	}
}

func (p *StaticProcessor) trySetFileCache(file string) {
	if p.cacheEnable && p.fileCache != nil {
		if f, err := newCachedFile(file); err == nil {
			p.fileCache.set(file, f)
		} else {
			// log?
		}
	}
}

// fileCache can cache static files
type fileCache struct {
	cache   *mcache.MCache
	root    string
	expire  int
	watcher *fsw.FsWatcher
}

func newFileCache(root string) (*fileCache, error) {
	if info, err := os.Stat(root); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, fmt.Errorf("%s isn't a directory", root)
	}

	fc := &fileCache{
		root:   root,
		cache:  mcache.NewMCache(),
		expire: 60,
	}
	if fw, err := fsw.NewFsWatcher(root); err == nil {
		fc.watcher = fw
		fc.watcher.Listen(fc.notify)
	}
	return fc, nil
}

func (fc *fileCache) notify(e fsw.Event) {
	name := cleanFilePath(e.Name)
	if e.Mode&fsw.Delete == fsw.Delete {
		fc.cache.Delete(name)
	}
	if e.Mode&fsw.Modify == fsw.Modify {
		//refresh??
		fc.cache.Delete(name)
	}
	if e.Mode&fsw.Rename == fsw.Rename {
		fc.cache.Delete(name)
	}
	if e.Mode&fsw.Create == fsw.Create {
		//add?
	}
}

// String
func (fc *fileCache) String() string {
	if fc == nil {
		return "<nil>"
	}
	return fmt.Sprintf("root: %s; expire:%d", fc.root, fc.expire)
}

func (fc *fileCache) get(name string) (*cachedFile, bool) {
	a, ok := fc.cache.Get(name)
	if !ok {
		return nil, false
	}

	f, ok := a.(*cachedFile)
	return f, ok
}

func (fc *fileCache) set(name string, f *cachedFile) {
	d := time.Duration(time.Second * time.Duration(fc.expire))
	fc.cache.SetAbs(name, f, d)
}

// cachedFile is cached content of static file
type cachedFile struct {
	content []byte
	name    string
	modtime time.Time
}

// String
func (f *cachedFile) String() string {
	if f == nil {
		return "<nil>"
	}
	return f.name
}

// 
func (f *cachedFile) Execute(ctx *HttpContext) error {
	http.ServeContent(ctx.Resonse, ctx.Request, f.name, f.modtime, bytes.NewReader(f.content))
	return nil
}

// ContentType return mime type of file
func (f *cachedFile) Type() string {
	return mime.TypeByExtension(filepath.Ext(f.name))

	// var buf [1024]byte
	// n, _ := io.ReadFull(content, buf[:])
	// b := buf[:n]
	// ctype = DetectContentType(b)
	// _, err := content.Seek(0, os.SEEK_SET)
}

func newCachedFile(name string) (*cachedFile, error) {
	var f *os.File
	var err error

	f, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var info os.FileInfo
	info, err = f.Stat()

	if err != nil {
		return nil, err
	} else if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", name)
	}

	var content []byte
	content, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &cachedFile{
		content: content,
		name:    name,
		modtime: info.ModTime(),
	}, nil
}
