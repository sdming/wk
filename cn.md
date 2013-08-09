GWK
===

简介
===

gwk(GO Web Server Kit)是GO语言的Web Server开发框架，简单易用，扩展性好，而且兼容Go App Engine。

安装
===

gwk只支持GO 1.1+版本，安装GO 1.1后，运行下面的命令即可。
  
  	go get github.com/sdming/wk

gwk依赖kiss、mcache和fsnotify三个package，如果没有自动安装成功的话，可以单独运行下面的命令安装:  

  	go get github.com/sdming/kiss   
  	go get github.com/sdming/mcache   
  	go get github.com/howeyc/fsnotify   


示例
===

gwk的文档比较简单，写的不是很详细，自带的demo写的比较全，主要的功能点都涉及到了。

另外Google App Engine上有一个展示gwk demo的网站，是用gwk框架搭建的，也是一个了解gwk的地方。网址是[http://gwk-demo.appspot.com]。 需要注意的是，这个demo网站的示例数据是放在内存的，多用户访问时会互相影响 。另外App Engine会自动管理服务实例以及会根据访问情况自动关闭或启动服务，示例数据也会受到影响。

启动服务
===

gwk不像revel那样是一个Web Server框架，需要自己写代码来启动gwk的服务。最简单的方式如下。

	server, err := wk.NewDefaultServer()

	if err != nil {
	    fmt.Println("NewDefaultServer error", err)
	    return
	}

	server.RouteTable.Get("/data/top/{count}").To(...)

	server.Start()

基本步骤就是:
* 创建HttpServer示例
* 注册路由
* 调用Start方法监听Http端口

接下来详细介绍gwk各个功能模块的用法，先从路由开始。

路由
===

gwk用RouteTable来存储注册的路由，RouteTable的定义如下:

	type RouteTable struct {
		Routes []*RouteRule
	}

当gwk接收到http请求时，按照顺序遍历RouteRule直到找到匹配的Route，如果没有找到则返回404。

RouteRule的定义如下:

	type RouteRule struct {

		// Methos is http method of request
		Method string

		// Pattern is path pattern
		Pattern string

		// Handler process request
		Handler Handler
	}

Method是Http method，比如:GET, POST, PUT, DELETE，*代表匹配所有的http method。  
Pattern是URL匹配的模式，具体的格式下面再讲。  
Handler是用来处理请求的代码，是一个接口，定义如下:  

	type Handler interface {
		Execute(ctx *HttpContext)
	}


gwk提供了若干方法来注册路由。一个最简单的方法是路由到一个func (*wk.HttpContext) (wk.HttpResult, error) 类型的函数，比如:

	// url: /data/top/10
	server.RouteTable.Get("/data/top/{count}").To(DataTopHandle)

上面的代码将GET /data/top/10这样的request path注册到一个func (*wk.HttpContext) (wk.HttpResult, error)类型的函数，例子中DataTopHandle的定义如下:

	func DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		if count, ok := ctx.RouteData.Int("count"); !ok {
			err = errors.New("parameter invalid:" + "count")
		} else {
			data := DataTop(count)
			result = wk.Json(data)
		}
		return
	}

例子中HttpContext是这次http请求的一个封装，HttpResult是这次请求返回的数据，详细的定义在下面会介绍。  

通过RouteData可以获得在request path中匹配到的参数，比如上面的{count}参数，RouteData的定义如下:  

	type RouteData map[string]string

RouteData提供了若干方法来简化Route参数的读取。

读取整数

	func (r RouteData) Int(name string) (int, bool)

读取整数，如果参数不存在或者不是有效的整数格式则返回一个缺省值

	func (r RouteData) IntOr(name string, v int) int

类似的方法还有
	
	Bool，BoolOr，Float，FloatOr，Str，StrOr


asp.net mvc 的开发人员应该会对"/data/top/{count}"这样的路由规则写法比较熟悉，gwk支持两种路由规则的写法。

* 正则表达式

		^/user/?(?P<action>[[:alnum:]]+)?/?(?P<arg>[[:alnum:]]+)?/?

	匹配类似/user/view/1这样的request path

* 正则表达式写起来比较麻烦，而且执行速度慢，你也可以用类似asp.net mvc的写法，比如

		"/data/top/{count}"   
		"/query/{year}-{month}-{day}"   
		"/basic#/{action}/{id}"   
		"/data/int/{p0}?"   
		"/data/range/{p0}-{p1}"   


	规则很简单，前缀匹配requets path，{}匹配的内容会被提取为参数。其中两个特殊字符需要介绍一下。

	* "#"字符代表精确匹配#之前的内容，#之后的为可选匹配，比如/basic和/basic/add，/basic/delete/1都匹配上面规则。

	* "?"字符代表匹配request path的结束，比如/data/int/1匹配上面的规则，/data/int/1/foo就不符合上面的规则。


gwk还提供了其他的方式来注册路由。

通过ToFunc方法注册路由到一个普通函数，比如

	// url: /data/int/1
	server.RouteTable.Get("/data/int/{p0}?").ToFunc(model.DataByInt)

	func DataByInt(i int) *Data {
		if i < 0 {
			i = 0
		}
		return newData(i)
	}


再比如

	// url: /data/range/1-9
	server.RouteTable.Get("/data/range/{p0}-{p1}").ToFunc(model.DataByIntRange)

	func DataByIntRange(start, end int) []*Data {
		data := make([]*Data, 2)
		data[0] = newData(start)
		data[1] = newData(end)
		return data
	}


因为GO的反射不能获得函数参数的名字，所以这里用p0,p1,p2...来代表函数的第0,1,2...个参数。

gwk会根据http请求中accept的内容来自动决定返回数据的格式，上面的例子中Data定义如下:
	
	type Data struct {
		Str   string
		Uint  uint64
		Int   int
		Float float32
		Byte  byte
	}

如果Request的Accept中包含字符串"xml"，则结果序列化为xml格式，如果包含"jsonp"，则结果序列化为jsonp，如果包含json，则序列化为json格式，详细信息可以参考下面的"格式化"一节

除了上面的p0,p1,p2指定参数的方式，可以用BindByNames按照名字来绑定函数的参数，比如:

	// url: /data/name/1
	server.RouteTable.Get("/data/name/{id}").ToFunc(model.DataByInt).
	    BindByNames("id")

上面的代码告诉gwk，路由参数"id"是函数DataByInt的第一个参数。  

	// url: /data/namerange/1-9
	server.RouteTable.Get("/data/namerange/{start}-{end}").ToFunc(model.DataByIntRange).
	    BindByNames("start", "end")

上面的代码告诉gwk，路由参数"start"，"end"是函数DataByIntRange的第一个和第二个参数。  

也可以绑定到querypath或者form中的参数，比如:

	// url: /data/namerange/?start=1&end=9
	server.RouteTable.Get("/data/namerange/").ToFunc(model.DataByIntRange).
	    BindByNames("start", "end")


再比如:

	// url: get /data/set?str=string&uint=1024&int=32&float=3.14&byte=64
	server.RouteTable.Get("/data/set?").ToFunc(model.DataSet).
	    BindByNames("str", "uint", "int", "float", "byte")

	func DataSet(s string, u uint64, i int, f float32, b byte) *Data {
		return &Data{
			Str:   s,
			Uint:  u,
			Int:   i,
			Float: f,
			Byte:  b,
		}
	}


如果参数比较多，一个合适的方法是将参数定义为一个struct，然后调用BindToStruct来绑定参数。比如:

	// url: post /data/post?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	server.RouteTable.Post("/data/post?").ToFunc(model.DataPost).BindToStruct()

	func DataPost(data Data) string {
		return data.String()
	}

	// url: post /data/postptr?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	server.RouteTable.Post("/data/postptr?").ToFunc(model.DataPostPtr).BindToStruct()

	func DataPostPtr(data *Data) string {
		return data.String()
	}

如果觉得每个函数定义一个参数对象比较麻烦，也可以用匿名对象:

	server.RouteTable.Get("/data/anonymous?").ToFunc(DataAnonymous).BindToStruct()

	func DataAnonymous(data struct {
	    Str   string
	    Uint  uint64
	    Int   int
	    Float float32
	    Byte  byte
	},) string {
	 
	    return fmt.Sprintln(data)
	}
 

更方便的路由方式是注册一个到controller的路由，在controller一节会介绍。


格式化
===

gwk会根据http请求中accept的内容来自动决定返回数据的格式，如果Request的Accept中包含字符串"xml"，则结果序列化为xml格式，如果包含"jsonp"，则结果序列化为jsonp，如果包含json，则序列化为json格式。具体的例子可以看http://gwk-demo.appspot.com/doc/routedemo 

你也可以通过设置Formatter来指定输出的格式，Formatter的定义如下

	type FormatFunc func(*HttpContext, interface{}) (HttpResult, bool)

gwk默认支持两种序列化方式，xml和json。比如:

	// url: /data/int/1/xml
	server.RouteTable.Get("/data/int/{p0}/xml").ToFunc(model.DataByInt).ReturnXml()

	// url: /data/int/1/json
	server.RouteTable.Get("/data/int/{p0}/json").ToFunc(model.DataByInt).ReturnJson()

ReturnXml指定DataByInt的返回值格式化为xml，ReturnJson指定DataByInt的返回值格式化为json。

你也可以自定义序列化的方式，比如:

	// url: /data/int/1/kson
	server.RouteTable.Get("/data/int/{p0}/kson").ToFunc(model.DataByInt).Return(formatKson)

	func formatKson(ctx *wk.HttpContext, x interface{}) (wk.HttpResult, bool) {
		b, _ := kson.Marshal(x)
		return wk.Content(string(b), "text/plain"), true
	}

kson格式是gwk的配置文件采用的格式，后文会详细介绍，上面的代码返回的数据如下:

	{
		Str:string:1
		Uint:100
		Int:10
		Float:1.1
		Byte:1
	}

gwk还提供了注册全局FormatFunc的地方：

	type FormatList []FormatFunc
	var Formatters FormatList

你可以通过增删Formatters或者修改Formatters的顺序来调整默认的格式化方式。


Controller
===

当需要注册的路由方法比较多，而且之间有一定的逻辑关系时，可以定义一个类似的asp.net mvc的Controller对象，然后将路由指向这个对象。代码可以参考https://github.com/sdming/wk/blob/master/demo/basic/controller/basic.go

一个简单的注册Controller的例子如下:

	basic = NewBasicController()

	// url: /basic/xxx/xxx
	server.RouteTable.Path("/basic#/{action}/{id}").ToController(basic)

用正则表达式的话，例子如下

	srv.RouteTable.Regexp("*", "^/user/?(?P<action>[[:alnum:]]+)?/?(?P<arg>[[:alnum:]]+)?/?").ToController(NewUserController())

注册到controller的路由，一个特殊的路由参数是{action}，它指定了调用Controller的哪一个方法。比如/basic/delete/32对应的action是"delete"，调用controller的delete方法，一个例子如下:

	// url: /basic/delete/32
	func (c *BasicController) Delete(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		l := len(c.data)
		if i, ok := ctx.RouteData.Int("id"); ok {
			c.deleteByInt(i)
		}
		return wk.Data(l - len(c.data)), nil
	}
	     
gwk现在版本的controller只支持func (ctx *wk.HttpContext) (result wk.HttpResult, err error)类型的方法，已经基本够用了。

再看一个例子:

	// url: /basic/add/?int=32&str=string&uint=1024&float=1.1&byte=64
	func (c *BasicController) Add(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		data := &model.Data{
			Int:   ctx.FormIntOr("int", 0),
			Uint:  uint64(ctx.FormIntOr("uint", 0)),
			Str:   ctx.FormValue("str"),
			Float: float32(ctx.FormFloatOr("float", 0.0)),
			Byte:  byte(ctx.FormIntOr("byte", 0)),
		}
		c.data = append(c.data, data)
		return wk.Data(data.String()), nil
	}

FormIntOr，FormFloatOr等函数是为了方便读取Request Form数据，可以参考RouteData的对应函数。

wk.Data函数返回 *DataResult对象，*DataResult实现了wk.HttpResult接口， HttpResult接口的详细介绍见后面的章节。

	func Data(data interface{}) *DataResult {
		return &DataResult{
			Data: data,
		}
	}

再看一个返回json数据的例子

	// url: /basic/int/32
	func (c *BasicController) Int(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		if id, ok := ctx.RouteData.Int("id"); ok {
			return wk.Json(c.getByInt(id)), nil
		}
		return wk.Data(""), nil
	}

再看一个直接读取post的数据，解析成json的例子:

	// url: post /basic/post
	func (c *BasicController) Post(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		var body []byte
		if body, err = ctx.ReadBody(); err != nil {
			return nil, err
		}
	 
		data := &model.Data{}
		err = json.Unmarshal(body, data)
		if err != nil {
			return nil, err
		}
		c.data = append(c.data, data)
		return wk.Data(true), nil
	}


如果没有找到路由参数{action}，gwk会将http method作为action，也就是把"get"，"post"，"delete"作为action，如果找不到对应的方法，会将{action}{method}组合起来作为action，这在一些场合还是比较实用的，如果还找不到对应的方法，则将"default"作为action，如果还是找不到对应的方法，则返回404.

事件
---

如果Controller的每个Action执行前后需要执行一些相同的代码怎么办？这就需要用到ActionSubscriber接口:

	type ActionSubscriber interface {
		OnActionExecuting(action *ActionContext)
		OnActionExecuted(action *ActionContext)
		OnException(action *ActionContext)
	}

OnActionExecuting在具体的Action执行之前执行，OnActionExecuted在具体的Action执行之后执行，OnException在具体的Action执行出错后执行。

通过ActionSubscriber可以做权限验证，数据验证，记录日志，同一错误处理等等。


HttpResult
===

凡是实现了HttpResult接口的对象，都可以作为gwk返回Web客户端的内容。HttpResult接口定义非常简单，只有一个方法:

	type HttpResult interface {
		Execute(ctx *HttpContext) error
	}

func Execute(ctx *HttpContext) error 方法定义了应该怎么样将数据返回客户端，*HttpContext 是当前http请求的上下文对象，后文会详细介绍。

gwk内置了支持几种常用的HttpResult。

ContentResult
---
	
	type ContentResult struct {
		ContentType string
		Data interface{}
	}

	func Content(data interface{}, contentType string) *ContentResult {
		return &ContentResult{
			Data:        data,
			ContentType: contentType,
		}
	}


ContentResult对应了raw html数据，直接将Data原样写入到http response中，如果你定义了ContentType参数，会在写Data之前先写http header:Content-Type。

如果Data实现了WriterTo、Reader接口，或者Data是[]byte 或者string，直接将Data写入Response，如果不是的话，gwk调用fmt.Fprintln将Data写入Response。  


JsonResult
---
	
	func Json(a interface{}) *JsonResult 


JsonResult顾名思义，先将数据序列化为json格式，再写入Response，默认会将http header的Content-Type设置为"application/json"，你也可以先给Content-Type设置一个值来阻止gwk设置Content-Type。


XmlResult
---

	func Xml(a interface{}) *XmlResult 


XmlResult将数据序列化为xml格式再写入Response，默认会将Content-Type设置为"text/xml"。

FileResult
---

	func File(path string) *FileResult

FileResult对应静态文件，实际上就是调用http.ServeFile来输出静态文件。 FileResult的path支持两种方式：绝对路径和相对路径，例子如下:

	func FileAbsolute(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		return wk.File(path.Join(ctx.Server.Config.RootDir, "public/humans.txt")), nil
	}

	func FileRelative(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		return wk.File("~/public/humans.txt"), nil
	}

如果path以~/开头则为相对路径，否则即为绝对路径。


FileStreamResult
---

	func FileStream(contentType, downloadName string, reader io.Reader, modtime time.Time) *FileStreamResult {
		return &FileStreamResult{
			ContentType:  contentType,
			DownloadName: downloadName,
			Data:         reader,
			ModifyTime:   modtime,
		}
	}


FileStreamResult对应一个Stream文件，如果设置了DownloadName参数，则将其作为浏览器保存文件的默认文件名，实际就是设置http header:"Content-Disposition"。

FileStreamResult内部是调用ServeContent。一个简单的例子如下： 
	

	// url: get /file/time.txt
	server.RouteTable.Get("/file/time.txt").To(FileHelloTime)

	func FileHelloTime(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		s := "hello, time is " + time.Now().String()
		reader := strings.NewReader(s)
		return wk.FileStream("", "hellotime.txt", reader, time.Now()), nil
	}


BundleResult
---
	
	func FileJsBundling(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		files := []string{"xxx/js/main.js", "xxx/js/plugins.js"}	
		return &wk.BundleResult{Files: files}, nil
	}


BundleResult是将若干相同类型的文件打包成一个文件返回，藉此提升响应速度。BundleResult原先是一个用来演示如何自定义HttpResult的demo，现在集成到gwk中。现在的版本还只支持绝对路径，后续的版本可能会支持相对路径。


RedirectResult
---

	func Redirect(urlStr string, permanent bool) *RedirectResult 

RedirectResult用来做http重定向，根据permanent参数决定返回 http.StatusMovedPermanently还是http.StatusFound。

NotFoundResult
---

NotFoundResult默认返回http.StatusNotFound，如果你开启了自定义404页面功能，则按如下逻辑返回：

1. 如果Request的http header "accept"包含"text/html"，先找public路径下的404.html，如果存在则返回404.html的内容。

2. 如果启用View引擎，并且views目录下存在404.html，则解析模板404.html返回。

3. 如果Request的http header "accept"包含"text/plain"，并且public路径下存在的404.txt，则返回404.txt的内容。

4. 如果上面的情况都不成立，则返回http.StatusNotFound。

开启了自定义404页面功能的方法是设置config的NotFoundPageEnable为true。设置方式见"配置"章节。


ErrorResult 
---

	func Error(message string) *ErrorResult {
		return &ErrorResult{
			Message: message,
		}
	}

ErrorResult顾名思义返回错误信息，默认返回http.StatusInternalServerError，如果你开启了自定义错误页面功能，则按如下逻辑返回：

1. 如果Request的http header "accept"包含"text/html"，先找public路径下的error.html，如果存在则返回error.html的内容。

2. 如果启用View引擎，并且views目录下存在error.html，则解析模板error.html返回。

3. 如果Request的http header "accept"包含"text/plain"，并且public路径下存在的error.txt，则返回error.txt的内容。

4. 如果上面的情况都不成立，则返回http.http.StatusInternalServerError。

开启了自定义错误页面功能的方法是这是config的ErrorPageEnable为true。


ViewResult
---

	func View(file string) *ViewResult {
		return &ViewResult{
			File: file,
		}
	}

ViewResult解析html模板并且输出到Response，因为这一块内容比较多，在"View引擎"一节单独介绍。

JsonpResult
---

返回Jsonp格式的数据，目前还没有实现。

NotModifiedResult
---

返回http.StatusNotModified

ChanResult
---

可以用来模拟BigPipe，定义如下

	type ChanResult struct {
		Wait        sync.WaitGroup
		Chan        chan string
		ContentType string
		Start       []byte
		End         []byte
		Timeout     time.Duration
	}

ChanResult会先输出Start，然后读取Chan中的字符串输出到客户端，最后输出End。


自定义HttpResult
---

自定义HttpResult十分简单，只要实现Execute(ctx *HttpContext) error方法就可以了，Go的interface机制让使用第三方的HttpResult或者开发一个HttpResult给别人使用变得很简单。

gwk的demo中包含一个自定义HttpResult的例子[QrCodeResult](https://github.com/sdming/wk/blob/master/demo/basic/model/qr.go])，可以将文本转化为二维码显示，这个例子不兼容App Engine，只能在线下运行demo程序看效果。



模板引擎
===

作为Web Engine框架，模板引擎是必不可少的，gwk的模板引擎基于Go自带的Html Template，在此基础上添加了一些新的功能。

* 内存中缓存编译的模板
* 内置了一系列Template Func
* 支持模板layout
* 支持partial view

先看几个具体的模板定义的例子，对gwk的模板有个直观的印象。

layout文件:_layout.html   

	<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<title>{{.title}}</title>
	<script type="text/javascript">
	</script>
	<style>
	</style>
	{{template "head" .}}
	</head>

	<body>
	<div id="header">
	{{partial "nav.html" .user }}
	</div>

	{{/* a comment */}}

	{{template "body" .}}

	<div id="footer">
	build by gwk
	</div>

	<script type="text/javascript">
	</script>

	{{template "script" .}}

	</body>
	</html>


模板文件:basic.html
	
	{{set . "title" "title demo" }}
	{{import "_layout.html" }}

	{{define "head" }}
	<script type="text/javascript">
	</script>
	<style>
	div{padding: 10px;}
	</style>
	{{end}}

	{{define "body" }}

	<h1>hello gwk!</h1>

	{{raw  "<!-- <script><style><html>  -->"}}

	<div>
	<lable for="selected">selected</lable>
	<select id="selected">
		<option value="" ></option>
		<option value="selected" {{selected true}}>selected</option>
	</select>
	<lable for="notselected">not selected</lable>
	<select id="notselected">
		<option value="" ></option>
		<option value="notselected" {{selected false}}>not selected</option>
	</select>
	</div>

	<div>
	<input id="checked" type="checkbox" {{checked true}}>checked</input>
	<input id="notchecked" type="checkbox" {{checked false}}>not checked</input>
	</div>

	<ul>
	<li id="eq">eq 123 123 = {{eq 123 123}}</li>
	<li id="eq">eqs "123" 123 = {{eqs "123" 123}}</li>
	<li id="gt">gt 3.14 3 = {{gt 3.14 3}}</li>
	<li id="le">le 1.1 2 = {{le 1.1 2}}</li>
	</ul>

	<div>{{nl2br "a\nb\nc" }}</div>

	<div id="settest-before">settest-before = {{.settest}}</div>
	{{set . "settest" "true"}}
	<div id="settest-after">settest-after = {{.settest}}</div>

	{{partial "user.html" .user}}

	{{end}}

	{{define "script" }}

	<script>
	{{jsvar "user" .user}}
	</script>

	{{end}}



partial view文件:nav.html 

	<div id="nav">Hi {{.Name}}</div>


另外一个partial view文件:user.html

	<ul id="div-{{.Name}}">
	<li>name:{{.Name}} </li>
	<li>age:{{.Age}}</li>
	<li><a href="{{.Web}}">web</a></li>
	<li><a href="mailto:{{.Email}}">email</a></li>
	</ul>


最后的输出应该类似下面的html

	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>title demo</title>
		<script type="text/javascript">
		</script>
		<style>
		</style>
		<script type="text/javascript">
		</script>
		<style>
			div{padding: 10px;}
		</style>
	</head>

	<body>

	<div id="header">
		<div id="nav">Hi Gopher</div>
	</div>

	<h1>hello gwk!</h1>

	<!-- <script><style><html>  -->

	<div>
		<lable for="selected">selected</lable>
		<select id="selected">
			<option value="" ></option>
			<option value="selected" selected>selected</option>
		</select>
		<lable for="notselected">not selected</lable>
		<select id="notselected">
			<option value="" ></option>
			<option value="notselected" >not selected</option>
		</select>
	</div>

	<div>
		<input id="checked" type="checkbox" checked>checked</input>
		<input id="notchecked" type="checkbox" >not checked</input>
	</div>

	<ul>
		<li id="eq">eq 123 123 = true</li>
		<li id="eq">eqs "123" 123 = true</li>
		<li id="gt">gt 3.14 3 = true</li>
		<li id="le">le 1.1 2 = true</li>
	</ul>

	<div>a<br/>b<br/>c</div>

	<div id="settest-before">settest-before = </div>
	<div id="settest-after">settest-after = true</div>

	<ul id="div-Gopher">
		<li>name:Gopher </li>
		<li>age:3</li>
		<li><a href="http://golang.org">web</a></li>
		<li><a href="mailto:gopher@golang.org">email</a></li>
	</ul>

	<div id="footer">
		build by gwk
	</div>

	<script type="text/javascript">
	</script>

	<script>
	 var user = {"Name":"Gopher","Age":3,"Web":"http://golang.org","Email":"gopher@golang.org"};
	</script>

	</body>
	</html>


更多模板的例子可以参考https://github.com/sdming/wk/tree/master/demo/basic/views/user


Template Func
---

gwk默认添加了若干Template Func

* eq: 		判断是否相等  
* eqs:		转化成字符串，再判断是否相等  
* gt:		大于  
* le:		小于  
* set:		设置map[string]interface{}元素的值  
* raw:		输出非转义的字符串
* selected:	输出字符串"selected"或者""  
* checked：	输出字符串"checked"或者""  
* nl2br:	将字符串中的"\n"替换为`"<br/>"`
* jsvar:	将Go的变量转化为javascript中的变量定义 
* import:	导入模板文件
* fv:		调用*http.Request.FormValue
* incl：	判断一个[]string中是否包含字符串v
* partial:	调用一个partial view


模板layout
---

你gwk中你可以定义若干个模板layout，然后在每个具体的模板文件中调用函数"import"引用某个layout文件，layout文件的路径为相对于模板根目录的相对路径。

	{{import "_layout.html" }}

需要注意的是，import要在模板输出具体内容之前调用才有效。


调用另一个模板
---

在gwk的模板文件中，可以通过函数partial调用另一个模板文件，这对于web服务端模块化开发来说很有用。在上面例子中定义了一个模板文件user.html来显示user对象的信息，在其他模板文件中就可以直接使用user.html了。

	{{partial "user.html" .user}}



模板缓存
---

默认配置下，gwk在第一次访问某个模板文件时会缓存编译后的模板*template.Template，后续访问这个模板时直接从缓存中读取*template.Template对象，如果模板的物理文件被修改，gwk会从缓存中删除对应的*template.Template对象。gwk使用fsnotify来监控物理文件，详细信息可以访问fsnotify的项目主页。

需要注意的是fsnofity在App Engine上不起作用，其实App Engine的更新机制也决定了不需要物理文件变更监控这样的功能。

你可以在plugin.conf关闭模板缓存功能，配置代码类似:

	#GoHtml config
	gohtml: {
		cache_enable:	true
	}
	# -->end GoHtml


接下来介绍gwk的内部实现机制。


GWK内部机制
===

gwk的内部逻辑十分简单，它的核心对象是HttpServer，HttpContent, HttpProcessor。下面分别介绍。

HttpServer
---

前面的例子里已经演示了怎么用函数NewDefaultServer创建一个缺省配置的HttpServer实例，NewDefaultServer会调用函数ReadDefaultConfigFile来尝试读取默认的配置文件，如果./conf/web.conf存在，则解析这个配置文件，如果解析出错或者文件不存在则调用函数NewDefaultConfig获得缺省配置。

你也可以用NewHttpServer传入WebConfig参数创建HttpServer实例。

	func NewHttpServer(config *WebConfig) (srv *HttpServer, err error) 

WebConfig的定义在后面介绍。

创建好HttpServer后，调用它的Start方法来监听http请求，启动Web服务。如果你的代码运行在Google App Engine这样不需要监听端口的平台，可以调用Setup方法来初始化HttpServer。Start方法内部实际上是先调用Setup方法，再调用http.Server的ListenAndServe方法。

HttpServer内部会创建一个http.Server实例，可以通过InnerServer方法来获得这个http.Server。

HttpServer有一个Variables字段，如果你有什么整个HttpServer共享的全局变量，可以放在Variables中。

	//server variables
	Variables map[string]interface{}


HttpServer实现了http.Handler接口，ServeHTTP函数的内部流程是:

1. 创建HttpContext
2. 循环执行每一个HttpProcessor的Execute方法

我们先介绍HttpContext。

HttpContext
---

HttpContext是对当前http请求的封装，定义如下:

	type HttpContext struct {
		Server *HttpServer
		Request *http.Request
		Resonse http.ResponseWriter
		Method string
		RequestPath string
		PhysicalPath string
		RouteData RouteData
		ViewData map[string]interface{}
		Result HttpResult
		Error error
		Flash map[string]interface{}
		Session Session
		SessionIsNew bool
	}

* Server:		当前的HttpServer
* Request:		http.Request
* Resonse: 		http.ResponseWriter
* Method:		http请求的method，比如GET，PUT，DELETE，POST...
* RequestPath:	http请求url的path部分
* PhysicalPath: 请求对应的物理文件，只有请求的是静态文件时，该字段才有值
* RouteData: 	RequestPath解析后的路由参数
* ViewData:		存放传递给View模板的数据
* Result:		此次请求的HttpResult
* Error:		此次请求中的Error
* Flash:		可以存放临时变量，生命周期为此次请求
* Session:		Session
* SessionIsNew:	Session是否此次在请求创建


HttpContext还定义了若干方法简化一些常见的操作:

RouteValue，读取RouteData里的数据

	func (ctx *HttpContext) RouteValue(name string) (string, bool) 

FormValue，调用http.Request的FormValue，FV也是相同的逻辑

	func (ctx *HttpContext) FV(name string) string 
	func (ctx *HttpContext) FormValue(name string) string

另外还有FormInt，FormIntOr，FormBool，FormBoolOr，FormFloat，FormFloatOr，前面已经做过介绍。

ReqHeader，读取Http request header的数据

	func (ctx *HttpContext) ReqHeader(name string) string 

SetHeader，设置Http resonse header的数据

	func (ctx *HttpContext) SetHeader(key string, value string) 
	
AddHeader，向http response header添加数据

	func (ctx *HttpContext) AddHeader(key string, value string) 

ContentType，设置http response header的"Content-Type"	
	
	func (ctx *HttpContext) ContentType(ctype string) 
	
Status，设置返回http status code

	func (ctx *HttpContext) Status(code int) 

Accept，读取http request header的"Accept"	

	func (ctx *HttpContext) Accept() string 

Write，调用http response的Write方法

	func (ctx *HttpContext) Write(b []byte) (int, error) 
	
Expires，设置http response header的"Expires"

	func (ctx *HttpContext) Expires(t string) 
	
SetCookie，设置cookie

	func (ctx *HttpContext) SetCookie(cookie *http.Cookie) 
	
Cookie，读取Cookie

	func (ctx *HttpContext) Cookie(name string) (*http.Cookie, error) 
	
SessionId，返回SessionId，只有启用了Session才有效

	func (ctx *HttpContext) SessionId() string 

GetFlash，读取Flash中的变量

	func (ctx *HttpContext) GetFlash(key string) (v interface{}, ok bool)

SetFlash，设置Flash中的变量

	func (ctx *HttpContext) SetFlash(key string, v interface{})

ReadBody，读取整个http request的内容

	func (ctx *HttpContext) ReadBody() ([]byte, error)

Flush，Flush当前Response中的数据到客户端

	func (ctx *HttpContext) Flush() {
	
前面介绍的ChanResult就是调用Flush把内容输出到客户端，代码基本逻辑如下:

	ctx.Write(c.Start)
	ctx.Flush()

	if c.Timeout < time.Millisecond {
		c.Timeout = defaultChanResultTimeout
	}

	waitchan := make(chan bool)
	donechan := make(chan bool)

	go func() {
		for s := range c.Chan {
			ctx.Write([]byte(s))
			ctx.Flush()
		}
		donechan <- true
	}()

	go func() {
		c.Wait.Wait()
		close(c.Chan)
		waitchan <- true
	}()

	select {
	case <-waitchan:
	case <-time.After(c.Timeout):
	}

	<-donechan
	ctx.Write(c.End)


HttpProcessor
---

HttpProcessor的定义如下

	type HttpProcessor interface {
		Execute(ctx *HttpContext)
		Register(server *HttpServer)
	}

Execute负责处理http请求，Register会在HttpServer初始化时调用一次，如果你的HttpProcessor需要执行一些初始化代码，可以放在Register方法中。

调用RegisterProcessor可以注册一个HttpProcessor

	func RegisterProcessor(name string, p HttpProcessor)

注册的HttpProcessor存在ProcessTable类型的全局变量中

	type ProcessTable []*Process

	type Process struct {
		Name string
		Path string
		Method string 
		Handler HttpProcessor
	}


如果一个Processor需要特定的条件才执行，可以设置它的Path和Method字段，Method是要匹配的http method，既GET、PUT、POST、DELETE...，"*"或者""匹配所有的http method，Path是要匹配的Request Path，目前版本是前缀匹配，以后可能改成支持通配符。

HttpServer启动时，默认注册三个HttpProcessor：StaticProcessor、RouteProcessor、RenderProcessor。

StaticProcessor
---

StaticProcessor负责处理静态文件，如果请求的路径能匹配到物理文件，则将HttpContext的的Result设置为FileResult。gwk只会将public子目录下的文件看做静态文件。

StaticProcessor支持缓存静态文件以及自定义http response header。缓存静态文件在缓存一节详细介绍，自定义输出的http header是指为每个静态文件的Response设置你定义的http header，比如统一为静态文件设置Cache-Control。下面是配置的例子：

	#static processor config
	static_processor: {

		cache_enable:	true
		cache_expire:	3600

		header:	{
			Cache-Control: 	max-age=43200
			X-Title: 		gwk-demo
		}
	}
	# -->end static processor


RouteProcessor
---

RouteProcessor负责按照你定义的路由规则调用具体的处理代码，逻辑很简单，只有几十行代码。

RenderProcessor
---

RenderProcessor负责执行HttpResult的Execute，也只有几十行代码。HttpResult没有赋值的话则返回404错误。

自定义HttpProcessor
---

你可以增删HttpProcessor或者调整顺序来改变默认的处理逻辑，比如你的程序是一个web api服务，不需要处理静态文件，则可以去掉RouteProcessor。ProcessTable定义了Append、InsertBefore、InsertAfter、Remove方法来简化对HttpProcessor的调整。

http gzip压缩
---

CompressProcessor可以对http输出做gzip压缩，需要注册到RenderProcessor之前才有效，其本质是用compressResponseWriter来代替默认的ResponseWriter。

	type CompressProcessor struct {
		Enable   bool
		Level    int
		MimeType string
	}

	type compressResponseWriter struct {
		rw            http.ResponseWriter
		writer        compresser
		contentType   string
		format        string
		headerWritten bool
	}

CompressProcessor设计时考虑能够按照MimeType或者RequestPath来过滤需要压缩的内容，但一直没实现，因为访问量小流量小的网站开不开启gzip压缩意义不大，访问量大的网站一般会用一些http反向代理或者http缓存的服务，自己没必要处理gzip压缩。

通过自定义HttpProcessor，你可以为全网站做统一的权限验证，访问限制，日志处理，错误处理等等。

除了HttpProcessor，你还可以通过gwk的事件机制来实现这些逻辑。


事件
===

gwk支持事件系统，但并没有硬编码有哪些事件，而是采用了比较松散的定义方式。

订阅事件有两种方式: 调用On函数或者OnFunc函数

	func On(moudle, name string, handler Subscriber) 

	func OnFunc(moudle, name string, handler func(*EventContext))


参数moudle是指订阅哪一个模块触发的事件，参数name是指订阅事件的名字，参数handler是处理事件的对象实例，是Subscriber类型的对象，Subscriber接口定义如下:

	type Subscriber interface {
		On(e *EventContext)
	}


	type SubscriberFunc func(*EventContext)

	func (f SubscriberFunc) On(e *EventContext) {
		f(e)
	}

EventContext定义如下:

	type EventContext struct {
		Moudle  string
		Name    string
		Source  interface{}
		Data    interface{}
		Context *HttpContext
	}

* Moudle：	触发事件的模块名
* Name：	事件名
* Source:	触发事件的变量
* Data：    事件附带的参数，每个事件可能不同，由Source负责赋值
* Context： HttpContext

如果想要触发一个自定义事件，要调用HttpServer的Fire方法：

	func (srv *HttpServer) Fire(moudle, name string, source, data interface{}, context *HttpContext) 

参数说明参照EventContext的定义。

使用事件系统可以做权限验证，日志、同一错误处理等等，十分方便。

demo/basic项目中的event.go演示了如何使用事件：

	wk.OnFunc("*", "*", eventTraceFunc)

这段代码调用OnFunc订阅了所有的事件，在eventTraceFunc中记录所有事件的触发时间并存在HttpContext的Flash字段中，在Server端结束所有处理前把这些数据返回客户端，这样客户端就能得到每个代码段的执行时间。返回的数据格式如下：

	_webserver	 start_request	 	0 ns 
	    _static	 start_execute	 	13761 ns 
	    _static	 end_execute	 	24829 ns 
	    _route	 start_execute	 	27988 ns 
	        _route	 start_action	50774 ns 
	        _route	 end_action	 	62984 ns 
	    _route	 end_execute	 	64255 ns 
	    _render	 start_execute	 	66379 ns 
	        _render	 start_result	68203 ns 
	        _render	 end_result	 	27631463 ns 
	    _render	 end_execute	 	27634149 ns 
	_webserver	 end_request	 	27636472 ns 

上面的数据列出了默认情况下gwk会触发的所有事件。

上面的例子给出了profile代码执行事件的一种思路。


配置
===

前面的例子都是基于gwk的默认配置，接下来将如何自定义配置以及如何使用gwk的配置框架。

gwk默认读取文件.conf/web.conf作为配置，如果文件不存在则采用预定义的默认配置。WebConfig的定义如下：


	type WebConfig struct {
		// 你可以给每一个Server设一个单独的名字，默认为""
		ServerKey string

		// 要监听的地址，默认为"0.0.0.0:8080"
		Address string

		// 根目录，默认为当前的工作目录
		RootDir string

		// 执行超时时间设置
		Timeout int

		// 静态文件的根目录，默认为RootDir下的public目录
		PublicDir string

		// 配置文件所在的目录，默认为RootDir下的conf目录
		ConfigDir string

		// View模板文件所在的目录，默认为RootDir下的views目录
		ViewDir string

		// 解析ConfigDir目录下的app.conf
		AppConfig *kson.Node

		// 解析ConfigDir目录下的plugin.conf
		PluginConfig *kson.Node

		// 读取Request的超时时间(秒)
		ReadTimeout int

		// 写Response的超时时间(秒)
		WriteTimeout int

		// Request headers的最大值
		MaxHeaderBytes int

		// 是否启用session
		SessionEnable bool

		// session的过期时间(秒)
		SessionTimeout int

		// SessionDriver is the name of driver
		SessionDriver string

		// 是否启用View引擎
		ViewEnable bool

		// 是否允许目录浏览，类似apache的Indexes 
		IndexesEnable bool

		// 是否允许自定义404页面
		NotFoundPageEnable bool

		// 是否允许自定义错误页面
		ErrorPageEnable bool

		// 是否开启Debug模式
		Debug bool

	}

如果ConfigDir目录下存在app.conf和plugin.conf文件，gwk解析这两个文件并将解析好的内容存在AppConfig字段和PluginConfig字段，建议app.conf存放程序的配置数据，plugin.conf存放gwk各模块的配置数据。

如果app.conf文件存在，gwk会使用fsnotify监控这个文件，如果文件改动就重新解析并刷新AppConfig字段。

kson
===

gwk的配置文件采用自创的kson格式，类似json或者yaml，项目地址在https://github.com/sdming/kiss/tree/master/kson，详细的例子请看项目的readme.md

kson特点是

* 首先方便人类阅读
* 字符串不需要用""，除非存在特殊字符
* 不需要用","分割字段，默认回车就是分隔符
* 类似yaml但是不依赖缩进
* 支持普通类型、map、slice、struct的序列化和反序列化
* 支持注释，#开始的行会被看做注释，不会被解析

先看一个配置数据的例子

	#app config file demo
	 
	#string
	key_string: demo
	 
	#string
	key_int:    101
	 
	#bool
	key_bool:   true
	 
	#float
	key_float:  3.14
	 
	#map
	key_map:    {
	    key1:   key1 value
	    key2:   key2 value
	}
	 
	#array
	key_array:  [
	    item 1      
	    item 2
	]
	 
	#struct
	key_struct: {
	    Driver:     mysql           
	    Host:       127.0.0.1
	    User:       user
	    Password:   password            
	}
	 
	#composite
	key_config: {   
	    Log_Level:  debug
	    Listen:     8000
	 
	    Roles: [
	        {
	            Name:   user
	            Allow:  [
	                /user       
	                /order
	            ]
	        } 
	        {
	            Name:   *               
	            Deny:   [
	                /user
	                /order
	            ]
	        } 
	    ]
	 
	    Db_Log: {
	        Driver:     mysql           
	        Host:       127.0.0.1
	        User:       user
	        Password:   password
	        Database:   log
	    }
	 
	    Env:    {
	        auth:       http://auth.io
	        browser:    ie, chrome, firefox, safari
	    }
	}


对应的Go代码的定义

	type Driver struct {
	    Driver   string
	    Host     string
	    User     string
	    Password string
	    A        string
	    B        string
	}
	 
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


kson格式的数据解析后存在kson.Node类型的实例中，具体的定义请参考kson项目的说明，这里只介绍kson.Node几个常用方法。

Dump

将node里的数据dump为kson格式的文本

	func (c *ConfigController) Dump(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    return wk.Data(c.node.MustChild("key_config").Dump()), nil
	}
        
Child

根据name返回node的子节点

 
	func (c *ConfigController) Child(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    _, ok := c.node.Child("key_string")
	    return wk.Data(ok), nil
	}

Query

查询node的子节点，现版本只支持按照节点名查询，以后可能支持按照属性查询比如 name[@field=xxx]

 
	func (c *ConfigController) Query(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    n, ok := c.node.Query("key_config Db_Log Host")
	    if ok {
	        return wk.Data(n.Literal), nil
	    }
	    return wk.Data(ok), nil
	}
        

ChildStringOrDefault

将子节点的内容解析为字符串返回，如果子节点不存在则返回默认值，类似的方法还有ChildIntOrDefault, ChildUintOrDefault, ChildFloatOrDefault, ChildBoolOrDefault, ChildStringOrDefault等

	func (c *ConfigController) ChildStringOrDefault(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    s := c.node.ChildStringOrDefault("key_string_not", "default value")
	    return wk.Data(s), nil
	}
        

ChildInt

将子节点的内容解析为Int64返回，如果子节点不存在则panic，类似的方法还有ChildInt, ChildUint, ChildFloat, ChildBool, ChildString等


	func (c *ConfigController) ChildInt(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    i := c.node.ChildInt("key_int")
	    return wk.Data(i), nil
	}
        

Bool

将节点的值解析为bool返回，类似的方法还有Int, Uint, Float, Bool, String等

	func (c *ConfigController) Bool(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    b, err := c.node.MustChild("key_bool").Bool()
	    if err != nil {
	        return nil, err
	    }
	    return wk.Data(b), nil
	}
        

Slice

将子节点的内容解析为[]string
	        
	func (c *ConfigController) Slice(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    data, err := c.node.MustChild("key_array").Slice()
	    if err != nil {
	        return nil, err
	    }
	    return wk.Data(data), nil
	}


Map

将子节点的内容解析为map[string]string

            
	func (c *ConfigController) Map(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    data, err := c.node.MustChild("key_map").Map()
	    if err != nil {
	        return nil, err
	    }
	    return wk.Data(data), nil
	}
        

Value

将子节点的内容解析到一个interface{}，传入的参数必须是可以通过reflect赋值的。

           
	func (c *ConfigController) Value(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    v := Driver{
	        Driver:   "driver",
	        Host:     "host",
	        User:     "user",
	        Password: "password",
	        A:        "aaa",
	        B:        "bbb",
	    }
	 
	    err := c.node.MustChild("key_struct").Value(&v)
	    if err != nil {
	        return nil, err
	    }
	    return wk.Data(v), nil
	}
        

接下来是一个解析复杂格式的例子

            
	func (c *ConfigController) Composite(ctx *wk.HttpContext) (wk.HttpResult, error) {
	    conf := &Config{}
	    err := c.node.MustChild("key_config").Value(conf)
	    if err != nil {
	        return nil, err
	    }
	    return wk.Data(conf), nil
	}
        

kson支持常见数据格式(不承诺支持所有的数据格式)，而且解析速度比json要快。


Session
===

Go的net/http本身不带session的机制，需要开发人员自行实现，gwk实现了内存中的session存储机制，如果需要将session存在其他地方比如redis或者memcache需要实现gwk的session.Driver接口。

session.Driver
---

session.Driver的接口如下


	type Driver interface {
		// 初始化
		Init(options string) error

		// Driver的名字
		Name() string

		// 添加key，如果重复返回false,error
		Add(sessionId, key string, value interface{}) (bool, error)

		// 读取key的值，如果不存在返回nil,false,nil，如果报错返回nil,false,error
		Get(sessionId, key string) (interface{}, bool, error)

		// 添加key，如果存在则更新
		Set(sessionId, key string, value interface{}) error

		// 移除key
		Remove(sessionId, key string) error

		// 根据sessionid创建新的session
		New(sessionId string, timeout time.Duration) error

		// 移除整个session
		Abandon(sessionId string) error

		// 判断sessionid是否存在
		Exists(sessionId string) (bool, error)

		// 返回session中所有key
		Keys(sessionId string) ([]string, error)
	}

gwk的Driver接口相比其他的框架要复杂一点，主要是为了Driver的开发人员可以实现更精确的控制。

自定义的session.Driver可以通过函数Register注册。

	func Register(name string, driver Driver)

gwk内置了In-memory的session.Driver的实现， 注册的名字为"session_default"，是基于开源项目MCache的，MCache的详细信息请参考其项目主页 [https://github.com/sdming/mcache]。

你可以通过修改web.conf文件或者直接修改WebConfig实例来启用或者关闭session机制，配置项如下:

* SessionEnable: 	是否启用session
* SessionTimeout: 	session的超时时间
* SessionDriver:	session driver的名字

如果你的session.Driver需要保存配置信息，请放在plugin.conf文件，gwk初始化时会将Config.PluginConfig.Child(session_driver_name)的数据作为options参数调用Driver的Init方法。

使用Session
---

你可以通过HttpContext的字段Session来获得session实例，字段SessionIsNew来获取session是否当前的请求创建的，方法SessionId()获取session的id。

./demo中的session.go文件演示了如何操作session。

获取session id

	// url: /session/id
	func (c *Session) Id(ctx *wk.HttpContext) (wk.HttpResult, error) {
		id := ctx.SessionId()
		return wk.Data(id), nil
	}

添加session


	// url: /session/add?k=test&v=101
	func (c *Session) Add(ctx *wk.HttpContext) (wk.HttpResult, error) {
		ok, err := ctx.Session.Add(ctx.FV("k"), ctx.FV("v"))
		return wk.Data(ok), err
	}

读取session中key的值

	// url: /session/get?k=test
	func (c *Session) Get(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
		v, _, err := ctx.Session.Get(ctx.FV("k"))
		return wk.Data(v), err
	}


设置session 

	// url: /session/set?k=test&v=101
	func (c *Session) Set(ctx *wk.HttpContext) (wk.HttpResult, error) {
		err := ctx.Session.Set(ctx.FV("k"), ctx.FV("v"))
		return wk.Data(err == nil), err
	}

移除session中的key

	// url: /session/remove?k=test
	func (c *Session) Remove(ctx *wk.HttpContext) (wk.HttpResult, error) {
		err := ctx.Session.Remove(ctx.FV("k"))
		return wk.Data(err == nil), err
	}

放弃整个session

	// url: /session/abandon
	func (c *Session) Abandon(ctx *wk.HttpContext) (wk.HttpResult, error) {
		err := ctx.Session.Abandon()
		return wk.Data(true), err
	}

返回session中所有的key

	// url: /session/keys
	func (c *Session) Keys(ctx *wk.HttpContext) (wk.HttpResult, error) {
		keys, err := ctx.Session.Keys()
		return wk.Data(fmt.Sprintln(keys)), err
	}

另外在session.go中还包含一个如何注册自定义session.Driver的例子。


缓存
===

view模板的缓存以及配置数据的缓存前文已经讲过，除此之外gwk可以将静态文件的内容缓存到内存中。这种缓存策略并不一定很有用，如果网站规模小流量不大，缓存静态文件的收益有限，而网站达到一定规模，为了提升性能，静态文件常部署在单独的静态文件服务器或者借助CDN，另外Go的内核用sendfile来处理静态文件，如果将其内容缓存到内存就没有办法用到sendfile的优势了。

开启gwk静态文件缓存的配制方法如下：
	
	#plugin.conf

	#static processor config
	static_processor: {

		#开启静态文件缓存，默认是false
		cache_enable:	true

		#缓存1小时(3600秒)，默认是86400秒
		cache_expire:	3600
	}
	# -->end static processor


gzip压缩可以参考前文的CompressProcessor部分


gwk并不内置供开发人员调用的Cache功能，如果需要in-memory的第三方缓存库，可以参考上文提到的MCache，项目在 [https://github.com/sdming/mcache]

日志
===

gwk本身不实现复杂的日志功能，只是公开了一个log.Logger类型的字段Logger，所有的日志信息会被记录到这个Logger中，另外你还可以通过设置LogLevel来调整记录日志的级别，默认为LogError，支持的日志级别为：

	const (
		LogError = iota
		LogDebug
	)


验证
===

虽然很多Web框架提供了验证功能，但gwk还没有这方面的计划。

ORM
===

gwk关注Web开发，短时间内不会包含ORM的功能，需要访问数据库的开发人员可以关注开源项目(kdb)[https://github.com/sdming/kdb]，项目刚启动，功能大概完成了30%。


Performance benchmark
==

Incoming


关键词
===

Golang, Go, Web Framework, Web Server Kit, gwk, MVC

