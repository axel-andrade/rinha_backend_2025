package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func main() {
	handler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("ok")
		case "/health":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString(`{"status":"ok"}`)
			ctx.SetContentType("application/json")
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("not found")
		}
	}

	fmt.Println("🚀 Servidor rodando em http://localhost:9999")
	err := fasthttp.ListenAndServe(":9999", handler)
	if err != nil {
		panic(err)
	}
}
