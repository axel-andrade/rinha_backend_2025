package http_router

import (
	"github.com/axel-andrade/rinha_backend_2025/internal/infra/bootstrap"
	"github.com/valyala/fasthttp"
)

func ConfigureRoutes(d *bootstrap.Dependencies) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/payments":
			if ctx.IsPost() {
				d.HTTPHandler.HandlePayments(ctx)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBodyString("Method Not Allowed")

		case "/payments-summary":
			if ctx.IsGet() {
				d.HTTPHandler.HandleSummary(ctx)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBodyString("Method Not Allowed")

		case "/health":
			d.HTTPHandler.HandleHealth(ctx)

		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Not Found")
		}
	}
}
