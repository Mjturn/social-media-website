package main

import (
    "github.com/gin-gonic/gin"
    "github.com/Mjturn/social-media-website/backend/routes"
)

func main() {
    router := gin.Default()
    router.Static("/static", "../frontend/static")
    routes.HandleRoutes(router)
    router.Run(":8080")
}
