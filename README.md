# framework.go
Golang micro web framework.

## Install

```
go get github.com/lucasmenendez/framework.go	
```

## Demo

```go
package main

import (
	"fmt"
	f "github.com/lucasmenendez/framework.go"
)

func req(c f.Context) {
	var err error
	var form f.Form
	if form, err = c.ParseForm(); err == nil {
		if key, ok := form.Get("key"); ok {
			fmt.Println(key)
		}
		
		c.PlainWrite([]byte("Hello world!"), 200)
	} else {
		fmt.Println(err)
	}
}

func mid(c f.Context) {
	fmt.Println(c.Params)
	c.Continue()
}


func main() {
	server := f.New()
	server.SetPort(9999)
	server.DebugMode(true)

	server.POST("/request/:id", req, mid)
	server.Run()
}
```
