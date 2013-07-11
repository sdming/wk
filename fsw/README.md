fsw
=====

fsw is File System Watcher, it monitoring folders and notify listeners if any thing changed. 


Requirements
---

go 1.1

Usage
---

go get github.com/howeyc/fsnotify
go get github.com/sdming/fsw  

Document
---
TODO:  


Getting Started
---

	package main

	import (
		"fmt"
		"github.com/sdming/wk/fsw"
		"time"
	)

	func main() {
		fw, err := fsw.NewFsWatcher(`d:\ddd\a`, `d:\ddd\b`)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("watch", fw)

		fw.Listen(func(e fsw.Event) {
			fmt.Println("event", e)
		})

		<-time.After(time.Minute)
		fw.Close()
		fmt.Println("close")
		<-time.After(time.Minute)
	}

	

Contributing
---
* github.com/sdming

License
---
Apache License 2.0  


About
----

