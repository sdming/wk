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
	defaultAddress string = "0.0.0.0:8080"

	// defaultPubicDir is default directory of public content
	defaultPubicDir string = "public"

	// defaultConfigDir is default directory of config files
	defaultConfigDir string = "conf"

	// defaultConfigDir is default directory of config files
	defaultViewDir string = "views"

	// defaultReadTimeout
	defaultReadTimeout int = 30

	// appConfigFile is default filename of app coinfig
	appConfigFile string = "app.conf"

	// webConfigFile is default file name of server coinfig
	webConfigFile string = "web.conf"

	// pluginConfigFile is default file name of plugin coinfig
	pluginConfigFile string = "plugin.conf"

	// defaultSessionTimeout is default value of SessionTimeout
	defaultSessionTimeout int = 20 * 60

	// defaultSessionDriver is the name of default session driver
	defaultSessionDriver string = "session_default"
)

// WebConfig is configuration of go web server
type WebConfig struct {
	// ServerKey is the identify of server
	ServerKey string

	// Address is the address to listen
	Address string

	// RootDir is the route directory of web application
	RootDir string

	// Timeout is timeout of http handle in second
	Timeout int

	// PublicDir is directory of static files
	PublicDir string

	// ConfigDir is directory of config files
	ConfigDir string

	// ViewDir is base directory of views
	ViewDir string

	// AppConfig is app configuration data
	AppConfig *kson.Node

	// PluginConfig is configuration data of plugins
	PluginConfig *kson.Node

	// ReadTimeout is maximum duration before timing out read of the request, in second
	ReadTimeout int

	// WriteTimeout is maximum duration before timing out write of the response, in second
	WriteTimeout int

	// MaxHeaderBytes is maximum size of request headers
	MaxHeaderBytes int

	// SessionEnable is true if enable session
	SessionEnable bool

	// SessionTimeout, session timeout in second
	SessionTimeout int

	// SessionDriver is the name of driver
	SessionDriver string

	// ViewEnable, enable view engine or not
	ViewEnable bool

	// IndexesEnable like Indexes of apache 
	IndexesEnable bool

	// Debug is true if run in debug model
	Debug bool
}

// String
func (conf *WebConfig) String() string {
	if conf == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%#v", conf)
}

// init
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

// defaultConfig return *WebConfig with default value
func defaultConfig() *WebConfig {
	rootdir := defaultRootPath()
	publicdir := path.Join(rootdir, defaultPubicDir)
	confdir := path.Join(rootdir, defaultConfigDir)
	viewdir := path.Join(rootdir, defaultViewDir)

	conf := &WebConfig{
		Address:        defaultAddress,
		RootDir:        rootdir,
		PublicDir:      publicdir,
		ConfigDir:      confdir,
		ViewDir:        viewdir,
		ReadTimeout:    defaultReadTimeout,
		SessionEnable:  true,
		SessionTimeout: defaultSessionTimeout,
		SessionDriver:  defaultSessionDriver,
		ViewEnable:     true,
		IndexesEnable:  true,
	}

	return conf
}

// NewDefaultConfig return *WebConfig with default value
func NewDefaultConfig() *WebConfig {
	conf := defaultConfig()
	conf.init()
	return conf
}

// ReadDefaultConfigFile parse default config file and return *WebConfig
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

// ConfigFromFile parse file and return *WebConfig
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
