package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler_main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}                                                                      //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	r.Use(cors.New(config))

	//TODO:
	r.POST("/transfer", Transfer)
	go (*IpLimitMap).IpMonitor()
	r.Run(":" + "9708")
}
