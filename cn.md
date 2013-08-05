GWK
===

简介
===

gwk(GO Web Server Kit)是GO语言的Web Server开发框架，简单易用，扩展性好，而且兼容Go App Engine。

安装
===

gwk只支持GO 1.1+版本，安装GO 1.1后，运行下面的命令即可。
  
  	go get github.com/sdming/wk

gwk依赖kiss和mcache两个package，如果没有自动安装成功的话，可以单独运行下面的命令安装:  

  	go get github.com/sdming/kiss  
  	go get github.com/sdming/mcache  


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

各种流行的MVC类开发框架比较多，controller应该不用做过多的介绍，接下来介绍HttpResult接口。








