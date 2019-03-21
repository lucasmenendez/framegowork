# SHGF: Simple HTTP golang framework
Only another golang micro web framework.

## Install

```
go get github.com/lucasmenendez/shgf	
```

## Demo

```go
package main

import (
	"fmt"

	"github.com/lucasmenendez/shgf"
)

func req(ctx *shgf.Context) (res *shgf.Response) {
	var err error
	if res, err = shgf.NewResponse(200); err != nil {
		fmt.Println(err)
	} else if err = ctx.ParseParams(); err != nil {
		fmt.Println(err)
	} else if ctx.Params["bar"] == "bar" {
		if err = res.JSON(ctx.Params); err != nil {
			fmt.Println(err)
		}

		return
	}
	res, _ = shgf.NewResponse(403)

	return
}

func mid(ctx *shgf.Context) (res *shgf.Response) {
	return ctx.Next()
}

func main() {
	s, err := shgf.New("0.0.0.0", 9999, true)
	if err != nil {
		fmt.Println(err)
	}

	s.GET("/foo/<string:bar>", req, mid)
	s.Listen()
}
```
