package wk_test

import (
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk"
	"testing"
)

type Config struct {
	Log_Level string
	Listen    uint
	Roles     []Role
	Db_Log    Db
	Env       map[string]string
}

type Role struct {
	Name  string
	Allow []string
	Deny  []string
}

type Db struct {
	Driver   string
	Host     string
	User     string
	Password string
}

func Equal(t *testing.T, name string, expect, actual interface{}) {
	if expect != actual {
		t.Errorf("%s Equal fail, expect %v, actual %v ", name, expect, actual)
	}
	return
}

func TestConfig(t *testing.T) {
	conf, err := wk.ReadDefaultConfigFile()

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%#v", conf)

	Equal(t, "ServerKey", "demoServer", conf.ServerKey)
	Equal(t, "Address", "127.0.0.0:80", conf.Address)
	Equal(t, "Timeout", 1020, conf.Timeout)
	Equal(t, "ReadTimeout", 1021, conf.ReadTimeout)
	Equal(t, "WriteTimeout", 1022, conf.WriteTimeout)
	Equal(t, "MaxHeaderBytes", 1023, conf.MaxHeaderBytes)

	Equal(t, "SessionEnable", true, conf.SessionEnable)
	Equal(t, "SessionTimeout", 3600, conf.SessionTimeout)
	Equal(t, "SessionDriver", "session_default", conf.SessionDriver)

	if v, err := conf.AppConfig.MustChild("key_string").String(); err != nil {
		t.Errorf("app config %s error %v", "key_string", err)
	} else {
		Equal(t, "key_string", "demo", v)
	}

	if v, err := conf.AppConfig.MustChild("key_int").Int(); err != nil {
		t.Errorf("app config %s error %v", "key_int", err)
	} else {
		Equal(t, "key_int", int64(101), v)
	}

	if v, err := conf.AppConfig.MustChild("key_bool").Bool(); err != nil {
		t.Errorf("app config %s error %v", "key_bool", err)
	} else {
		Equal(t, "key_bool", true, v)
	}

	if v, err := conf.AppConfig.MustChild("key_float").Float(); err != nil {
		t.Errorf("app config %s error %v", "key_float", err)
	} else {
		Equal(t, "key_float", 3.14, v)
	}

	if v, err := conf.AppConfig.MustChild("key_map").Map(); err != nil {
		t.Errorf("app config %s error %v", "key_map", err)
	} else {
		Equal(t, "key_map.key1", "key1 value", v["key1"])
		Equal(t, "key_map.key2", "key2 value", v["key2"])
	}

	if v, err := conf.AppConfig.MustChild("key_array").Slice(); err != nil {
		t.Errorf("app config %s error %v", "key_array", err)
	} else {
		Equal(t, "key_array.0", "item 1", v[0])
		Equal(t, "key_array.1", "item 2", v[1])
	}

	db := &Db{}
	if err := conf.AppConfig.MustChild("key_struct").Value(db); err != nil {
		t.Errorf("app config %s error %v", "key_struct", err)
	} else {
		Equal(t, "key_struct.Driver", "mysql", db.Driver)
		Equal(t, "key_struct.Host", "127.0.0.1", db.Host)
		Equal(t, "key_struct.User", "user", db.User)
		Equal(t, "key_struct.Password", "password", db.Password)
	}

	c := &Config{}
	if err := conf.AppConfig.MustChild("key_config").Value(c); err != nil {
		t.Errorf("app config %s error %v", "key_config", err)
	} else {
		t.Logf("%#v", c)
	}
}

type Option struct {
	Int     int
	String  string
	Float64 float64
	Bool    bool
	Map     map[string]string
}

func TestPluginConfig(t *testing.T) {
	conf, err := wk.ReadDefaultConfigFile()

	if err != nil {
		t.Error(err)
		return
	}

	node := conf.PluginConfig.MustChild("session_debug")
	dump := node.Dump()
	t.Log("session_debug", dump)

	option := &Option{}
	kson.Unmarshal([]byte(dump), option)

	if err := conf.PluginConfig.MustChild("session_debug").Value(option); err != nil {
		t.Errorf("plugin config %s error %v", "session_debug", err)
	} else {
		Equal(t, "session_debug.Driver", 1024, option.Int)
		Equal(t, "session_debug.String", "string demo", option.String)
		Equal(t, "session_debug.Float32", 3.14, option.Float64)
		Equal(t, "session_debug.Bool", true, option.Bool)
		Equal(t, "session_debug.Map.key1", "key1 value", option.Map["key1"])
		Equal(t, "session_debug.Map.key2", "key2 value", option.Map["key2"])
	}

}
