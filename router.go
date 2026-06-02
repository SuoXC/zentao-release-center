package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	center "github.com/yi-nology/zentao-release-center/biz/handler/release/center"
)

func customizedRegister(r *server.Hertz) {
	api := r.Group("/api")
	api.POST("/notify/preview", center.NotifyPreview)
	api.POST("/notify/send", center.NotifySend)
}
