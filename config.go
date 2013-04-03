// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/kson"
	"os"
	"path"
)

const (
	// defaultAddress is default address to listen
	defaultAddress string = "127.0.0.1:8080"

	// defaultPubicDir is default path of public content
	defaultPubicDir string = "public"

	// defaultConfigDir is default path of config files
	defaultConfigDir string = "conf"

	// appConfigFile is default app coinfig file
	appConfigFile string = "app.conf"

	// webConfigFile is default web coinfig file
	webConfigFile string = "web.conf"

	// pluginConfigFile is default plugin coinfig file
	pluginConfigFile string = "plugin.conf"
)

// web server config
type WebConfig struct {
	// ServerKey to identify a server
	ServerKey string

	// Address is the address to listen
	Address string

	// RootDir is the route of web application
	RootDir string

	// HandleTimeout in Second
	Timeout int

	// PublicDir is path of static files
	PublicDir string

	// ConfigDir is path of config files
	ConfigDir string

	// AppConfig is config data for web app
	AppConfig *kson.Node

	// PluginConfig is config data for plugins
	PluginConfig *kson.Node

	// maximum duration before timing out read of the request, in second
	ReadTimeout int

	// maximum duration before timing out write of the response, in second
	WriteTimeout int

	// maximum size of request headers
	MaxHeaderBytes int
}

// String 
func (conf *WebConfig) String() string {
	if conf == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%#v", conf)
}

//  
func (conf *WebConfig) init() {
	if conf == nil {
		return
	}

	if isDirExists(conf.ConfigDir) {
		appFile := path.Join(conf.ConfigDir, appConfigFile)
		if isFileExists(appFile) {
			conf.AppConfig, _ = kson.ParseFile(appFile)
		}

		pluginFile := path.Join(conf.ConfigDir, pluginConfigFile)
		if isFileExists(pluginFile) {
			conf.PluginConfig, _ = kson.ParseFile(pluginFile)
		}
	}

}

func defaultConfig() *WebConfig {
	rootdir := defaultRootPath()
	publicdir := path.Join(rootdir, defaultPubicDir)
	confdir := path.Join(rootdir, defaultConfigDir)

	conf := &WebConfig{
		Address:   defaultAddress,
		RootDir:   rootdir,
		PublicDir: publicdir,
		ConfigDir: confdir,
	}

	return conf
}

// DefaultConfig return web config with default value
func NewDefaultConfig() *WebConfig {
	conf := defaultConfig()
	conf.init()
	return conf
}

// defaultRootPath return default root of web application
func ReadDefaultConfigFile() (conf *WebConfig, err error) {
	root := defaultRootPath()
	if root == "" {
		err = errors.New("root path is empty")
		return
	}

	confdir := path.Join(root, defaultConfigDir)
	if !isDirExists(confdir) {
		err = errors.New("conf directory is not exists:" + confdir)
		return
	}

	conffile := path.Join(confdir, webConfigFile)
	if !isFileExists(conffile) {
		err = errors.New("web conf file is not exists:" + conffile)
		return
	}

	return ConfigFromFile(conffile)
}

// defaultRootPath return default root of web application
func defaultRootPath() string {
	pwd, err := os.Getwd()
	if err == nil {
		return pwd
	}
	return path.Dir(os.Args[0])
}

// configFromFile can read config value from a file
func ConfigFromFile(file string) (conf *WebConfig, err error) {
	if !isFileExists(file) {
		err = errors.New("file is not exist:" + file)
		return
	}

	var node *kson.Node
	if node, err = kson.ParseFile(file); err != nil {
		return
	}
	conf = defaultConfig()
	if err = node.Value(conf); err != nil {
		return
	}
	conf.init()
	return

}
