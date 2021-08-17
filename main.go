package fhttp

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
)

type FHttpContext struct {
	body  interface{}
	reply func(data interface{})
}

func replyRaw(ctx *fasthttp.RequestCtx, data interface{}) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return
	}

	ctx.SetBody(rawData)
}

type FHttpInstance struct {
	get   func(route string, options FHttpRouteOptions, callback func(context FHttpContext))
	post  func(route string, options FHttpRouteOptions, callback func(context FHttpContext))
	put   func(route string, options FHttpRouteOptions, callback func(context FHttpContext))
	patch func(route string, options FHttpRouteOptions, callback func(context FHttpContext))

	listen func(port uint16)
}

type FHttpRouteSchema struct {
	Body interface{}
}

type FHttpRouteOptions struct {
	Schema FHttpRouteSchema
}

func New() FHttpInstance {
	routes := make(map[string]func(context FHttpContext))
	var parser fastjson.Parser

	return FHttpInstance{
		get: func(route string, options FHttpRouteOptions, callback func(context FHttpContext)) {
			if !strings.HasPrefix(route, "/") {
				route = "/" + route
			}

			routes["GET"+route] = callback
		},
		listen: func(port uint16) {
			err := fasthttp.ListenAndServe(":"+strconv.Itoa(int(port)), func(ctx *fasthttp.RequestCtx) {
				var key = string(ctx.Method()) + string(ctx.Path())

				if routes[key] != nil {
					body, err := parser.ParseBytes(ctx.PostBody())
					if err != nil {
						// ctx.Error("Invalid body", fasthttp.)
						return
					}

					routes[key](FHttpContext{
						body: body,
						reply: func(data interface{}) {
							replyRaw(ctx, data)
						},
					})
				} else {
					ctx.Error("Unsupported path", fasthttp.StatusNotFound)
				}
			})

			if err != nil {
				panic(err)
			}
		},
	}
}
