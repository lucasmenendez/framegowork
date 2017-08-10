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
	if form, err := c.ParseMultiPartForm(); err == nil {
		fmt.Println(form)
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
