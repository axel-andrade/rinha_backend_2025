package http_adapter

import (
	"github.com/valyala/fasthttp"
)

func ConfigureRoutes(handler *Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/payments":
			if ctx.IsPost() {
				handler.HandlePayments(ctx)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBodyString("Method Not Allowed")

		case "/payments-summary":
			if ctx.IsGet() {
				handler.HandleSummary(ctx)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBodyString("Method Not Allowed")

		case "/health":
			handler.HandleHealth(ctx)

		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Not Found")
		}
	}
}
