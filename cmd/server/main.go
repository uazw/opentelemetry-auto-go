package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()

    // Root route for backward compatibility
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("hello from otel-go-webapp (goframe)")
    })

    // Hello World endpoint
    s.BindHandler("/hello", func(r *ghttp.Request) {
        r.Response.Write("hello world")
    })

    s.SetPort(8080)
    s.Run()
}
