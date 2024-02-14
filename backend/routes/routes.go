package routes

import (
    "github.com/gin-gonic/gin"
)

func HandleRoutes(router *gin.Engine) {
    router.GET("/", func(context *gin.Context) {
        context.File("../frontend/static/index.html")
    })
    
    router.GET("/register", func(context *gin.Context) {
        context.File("../frontend/static/register.html")
    })
    
    router.GET("/login", func(context *gin.Context) {
        context.File("../frontend/static/login.html")
    })
}
