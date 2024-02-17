package main

import (
    "fmt"
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/Mjturn/social-media-website/backend/routes"
)

func main() {
    if err := godotenv.Load("../.env"); err != nil {
        log.Fatal(err)
    }

    databaseUsername := os.Getenv("DATABASE_USERNAME")
    databasePassword := os.Getenv("DATABASE_PASSWORD")
    databaseHost := os.Getenv("DATABASE_HOST")
    databasePort := os.Getenv("DATABASE_PORT")
    databaseName := os.Getenv("DATABASE_NAME")
    
    databaseConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseUsername, databasePassword, databaseHost, databasePort, databaseName)
    
    database, err := sql.Open("mysql", databaseConnectionString)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    if err := database.Ping(); err != nil {
        log.Fatal(err)
    }

    router := gin.Default()
    router.Static("/static", "../frontend/static")
    router.LoadHTMLGlob("../frontend/*.html")
    store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET_KEY")))
    router.Use(sessions.Sessions("user_session", store))
    routes.HandleRoutes(router, database)
    
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
