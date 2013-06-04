gwk
=====

Web, Let's GO!

gwk is a smart &amp; lightweight web server engine 
gwk is webkit for web server  


Roadmap
---

* 0.1 fork gomvc, refactoring. april 2013  
* 0.2 configration framework. april 2013  
* 0.3 web api server. april 2013  
* 0.4 test on go 1.1. may 2013  
* 0.5 cookie, session. may 2013 
* 0.6 view engine. may 2013  
* 0.7 watch & load config. june 2013      
* 0.8 custome error & 404 page. june 2013    
* 0.9 file upload. june 2013       
* 1.0 go. july 2013   
* 1.x    

Requirements
---

go 1.1

Usage
---

go get github.com/sdming/kiss  
go get github.com/sdming/mcache
go get github.com/sdming/wk  

Document
---
Take time to translate documents to english.  

Getting Started
---

	// ./demo/basic/main.go for more detail

	server, err := wk.NewDefaultServer()

	if err != nil {
		fmt.Println("NewDefaultServer error", err)
		return
	}

	server.RouteTable.Get("/data/top/{count}").To(...)

	server.Start()


How to run demo  
---
	cd ./demo/basic  
	go run main.go  


Route
---

Go cann't get parameter name of function by reflect, so it's a littel tricky to create parameters when call function by reflect.  


	// url: /demo/xxx/xxx
	// route to controller
	server.RouteTable.Path("/demo/{action}/{id}").ToController(controller)

	// url: /data/top/10
	// func: DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error)
	// route to func (*wk.HttpContext) (wk.HttpResult, error)
	server.RouteTable.Get("/data/top/{count}").To(DataTopHandle)

	// url: /data/int/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/int/{p0}?").ToFunc(model.DataByInt)

	// url: /data/range/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/range/{p0}-{p1}").ToFunc(model.DataByIntRange)

	// url: /data/int/1/xml
	// func: DataByInt(i int) *Data
	// return xml
	server.RouteTable.Get("/data/int/{p0}/xml").ToFunc(model.DataByInt).ReturnXml()

	// url: /data/int/1/json
	// func: DataByInt(i int) *Data
	// return json
	server.RouteTable.Get("/data/int/{p0}/json").ToFunc(model.DataByInt).ReturnJson()

	// url: /data/int/1/kson
	// func: DataByInt(i int) *Data
	// return custome formatted data
	server.RouteTable.Get("/data/int/{p0}/kson").ToFunc(model.DataByInt).Return(formatKson)

	// url: /data/name/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/name/{id}").ToFunc(model.DataByInt).
		BindByNames("id")

	// url: /data/namerange/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/{start}-{end}").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: /data/namerange/?start=1&end=9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: post /data/post?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func: DataPost(data Data) string 
	// route http post to function, build struct parameter from form  
	server.RouteTable.Post("/data/post?").ToFunc(model.DataPost).BindToStruct()

	// url: post /data/postptr?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func DataPostPtr(data *Data) string
	// route http post to function, build struct parameter from form
	server.RouteTable.Post("/data/postptr?").ToFunc(model.DataPostPtr).BindToStruct()

	// url: delete /data/delete/1
	// func: DataDelete(i int) string 
	// route http delete to function
	server.RouteTable.Delete("/data/delete/{p0}").ToFunc(model.DataDelete)

	// url: get /data/set?str=string&uint=1024&int=32&float=3.14&byte=64
	// func: DataSet(s string, u uint64, i int, f float32, b byte) *Data 
	// test diffrent parameter type
	server.RouteTable.Get("/data/set?").ToFunc(model.DataSet).
		BindByNames("str", "uint", "int", "float", "byte")


Controller  
---

Route url like "/demo/{action}" to T, call it's method by named {action}.     

Current version only support method of type (*HttpContext) (result HttpResult, error)  

	// url: /demo/int/32
	func (c *DemoController) Int(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		if id, ok := ctx.RouteData.Int("id"); ok {
			return wk.Json(c.getByInt(id)), nil
		}
		return wk.Data(""), nil
	}


Formatter
---

If function doesn't return HttpResult, server will convert result to HttpResult according Accept of request. you can set Formatter for each function, or register global Formatter.  

	// return formatted (HttpResult, true) or return (nil, false) if doesn't format it
	type FormatFunc func(*HttpContext, interface{}) (HttpResult, bool)


Configration 
---

wk can manage configration for you. below is config example, more example at test/config_test.go

	#app config file demo

	#string
	key_string: demo

	#string
	key_int: 	101

	#bool
	key_bool: 	true

	#float
	key_float:	3.14

	#map
	key_map:	{
		key1:	key1 value
		key2:	key2 value
	}

	#array
	key_array:	[
		item 1		
		item 2
	]

	#struct
	key_struct:	{
		Driver:		mysql			
		Host: 		127.0.0.1
		User:		user
		Password:	password			
	}

	#composite
	key_config:	{	
		Log_Level:	debug
		Listen:		8000

		Roles: [
			{
				Name:	user
				Allow:	[
					/user		
					/order
				]
			} 
			{
				Name:	*				
				Deny: 	[
					/user
					/order
				]
			} 
		]

		Db_Log:	{
			Driver:		mysql			
			Host: 		127.0.0.1
			User:		user
			Password:	password
			Database:	log
		}

		Env:	{
			auth:		http://auth.io
			browser:	ie, chrome, firefox, safari
		}
	}


Session  
---

Example of session: ./demo/basic/controller/session.go  

By default, session store in memory(mcache),  if you want to store somewhere else(redis, memcahe...), need to:

1.	implement interface: session.Driver
2. 	call session.Register to register it
3. 	set SessionDriver of web.conf

Basic Conception
---

HttpContext is wrap of http.Request & http.Response  

HttpProcessor handle request and build HttpResult(maybe)
	
	type HttpProcessor interface {
		Execute(ctx *HttpContext)

		// Register is called once when server init  
		Register(server *HttpServer)
	}


HttpResult know how to write http.Response

	type HttpResult interface {
		Execute(ctx *HttpContext)
	}


The lifecyle is 

1. receive request, create HttpContext   
3. run each HttpProcessor  
4. execute HttpResult  

HttpProcessor
---

You can add or remove HttpProcessor before server starting.  

	type ProcessTable []*Process

	func (pt *ProcessTable) Append(p *Process)

	func (pt *ProcessTable) Remove(name string)

	func (pt *ProcessTable) Insert(p *Process, index int) 



http result
---

* ContentResult: 	html raw 
* JsonResult: 		application/json
* XmlResult: 		application/xml
* JsonpResult: 		application/jsonp
* ViewResult 		view
* FileResult 		static file
* FileStreamResult 	stream file
* RedirectResult 	redirect
* NotFoundResult  	404
* ErrorResult 		error
* BundleResult		bundle of files  

Event
---

Call Fire() to fire an event

	Fire(moudle, name string, source, data interface{}, context *HttpContext) 

Call On() to listen events 

	On(moudle, name string, handler Subscriber) 


View engine
---

Template Funcs  
* eq: 	equal   
* eqs: 	compare as string  
* get:	greater  
* le:	less  
* set:	set map[string]interface{}
* raw:	unescaped html  
* incl:	include or not
* selected:	return "selected" or ""
* checked: 	return "checked" or ""
* nl2br:	replace "\n" with "<br/>" 
* jsvar:	create javascript variable, like var name = {...}
* import:	import a template file
* partial:	call a template

you can find examples in folder "./test/views/" or "./demo/basic/views/"

Example
---

TODO:  
* 	basic example  
	run default http server   
	file: ./demo/basic/main.go  

* 	rest api example  
	run rest http api server  
	file: ./demo/rest/main.go  

* 	httpresult example  
	how to write a http result to return qrcode image   
	file: ./demo/basic/model/qr.go  

*	event example  
	how to listen events   
	file: ./demo/basic/model/event.go    

* 	custom processor  
	how to register a Processor to compress http response  
	file: ./compress.go   

* 	file stream example  
	how to return a file stream    
	file: ./demo/basic/model/file.go  

	how to bundling several js files into one file
	file: ./demo/basic/model/file.go  	
	
* 	BigPipe example  
	how to simulate BigPipe & comet   
	file: ./demo/basic/model/bigpipe.go (need fix bug?)     
 
* 	session example  
	how to add, get,remove... session
	file: ./demo/basic/controller/session.go     

* 	view example  
	how to use viewengine  
	file: ./demo/basic/controller/user.go   


ORM
---
Maybe, maybe not. don't have a plan yet. focus on web server.    

Validation
---
No

Css & js bundling
---
Do we really need it ? 
./demo/basic/model/file.go is a very sample of how to bundle files.   

Cache, gzip
---
nginx, haproxy, Varnish can provide awesome service

Contributing
---
* github.com/sdming

License
---
Apache License 2.0  


About
----

