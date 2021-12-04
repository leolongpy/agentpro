package http

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func agentStatusRoutes(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "ok")
}
