[![GoDoc](https://godoc.org/github.com/lucasmenendez/shgf?status.svg)](https://godoc.org/github.com/lucasmenendez/shgf)
[![Build Status](https://travis-ci.org/lucasmenendez/shgf.svg?branch=master)](https://travis-ci.org/lucasmenendez/shgf)
[![Report](https://goreportcard.com/badge/github.com/lucasmenendez/shgf)](https://goreportcard.com/report/github.com/lucasmenendez/shgf)

# SHGF: Simple HTTP golang framework
Opinionated simple HTTP golang framework. Provides simple API to create a HTTP server and a group of functions to register new available routes with its handler by HTTP method.

## Install

```
go get github.com/lucasmenendez/shgf	
```

## Docs
Read all the documentation into [GoDoc](https://godoc.org/github.com/lucasmenendez/shgf) article.

## Example

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
