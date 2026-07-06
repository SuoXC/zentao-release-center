package main

import (
	"context"
	"embed"
	"io/fs"
	"path"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

//go:embed all:frontend/dist
var frontendDist embed.FS

// mountFrontend 注册前端静态资源（嵌入 frontend/dist）。
// 1) 真实存在的静态文件直接由嵌入 FS 返回；
// 2) 非 /api 路径且文件不存在时回退到 index.html（SPA 路由）。
// dev 期若 frontend/dist 尚未构建（占位占位 index.html 含 placeholder:true），跳过挂载。
func mountFrontend(h *server.Hertz) error {
	sub, err := fs.Sub(frontendDist, "frontend/dist")
	if err != nil {
		return nil
	}
	idx, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		return nil
	}
	if strings.Contains(string(idx), "placeholder:true") {
		return nil
	}

	h.GET("/favicon.svg", func(ctx context.Context, c *app.RequestContext) { serveAsset(c, sub, "favicon.svg") })
	h.GET("/icons.svg", func(ctx context.Context, c *app.RequestContext) { serveAsset(c, sub, "icons.svg") })

	h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		p := string(c.Request.URI().Path())
		if strings.HasPrefix(p, "/api") {
			c.JSON(404, map[string]string{"code": "404", "message": "route not found"})
			return
		}
		cleaned := strings.TrimPrefix(p, "/")
		if cleaned != "" && hasAsset(sub, cleaned) {
			serveAsset(c, sub, cleaned)
			return
		}
		c.Data(200, "text/html; charset=utf-8", idx)
	})
	return nil
}

func serveAsset(ctx *app.RequestContext, sub fs.FS, p string) {
	cleaned := path.Clean("/" + p)
	if cleaned == "/" || strings.Contains(cleaned, "..") {
		ctx.Status(404)
		return
	}
	data, err := fs.ReadFile(sub, cleaned[1:])
	if err != nil {
		ctx.Status(404)
		return
	}
	ctx.Data(200, guessContentType(cleaned), data)
}

func hasAsset(sub fs.FS, p string) bool {
	if strings.Contains(p, "..") {
		return false
	}
	if _, err := fs.Stat(sub, p); err == nil {
		return true
	}
	return false
}

func guessContentType(name string) string {
	switch {
	case strings.HasSuffix(name, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(name, ".png"):
		return "image/png"
	case strings.HasSuffix(name, ".jpg"), strings.HasSuffix(name, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(name, ".gif"):
		return "image/gif"
	case strings.HasSuffix(name, ".woff"):
		return "font/woff"
	case strings.HasSuffix(name, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(name, ".ttf"):
		return "font/ttf"
	case strings.HasSuffix(name, ".eot"):
		return "application/vnd.ms-fontobject"
	case strings.HasSuffix(name, ".json"):
		return "application/json; charset=utf-8"
	case strings.HasSuffix(name, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(name, ".html"):
		return "text/html; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}
