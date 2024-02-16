package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "database/sql"
    "golang.org/x/crypto/bcrypt"
)

func HandleRoutes(router *gin.Engine, database *sql.DB) {
    router.GET("/", func(context *gin.Context) {
        context.File("../frontend/static/index.html")
    })
    
    router.GET("/register", func(context *gin.Context) {
        context.File("../frontend/static/register.html")
    })

    router.POST("/register", func(context *gin.Context) {
        usernameInput := context.PostForm("username-input")
        passwordInput := context.PostForm("password-input")

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordInput), bcrypt.DefaultCost)
        if err != nil {
            context.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        var count int
        err = database.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", usernameInput).Scan(&count)
        if err != nil {
            context.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        if count > 0 {
            context.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
            return
        }

        _, err = database.Exec("INSERT INTO users (username, password) VALUES (?, ?)", usernameInput, hashedPassword)
        if err != nil {
            context.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        context.Redirect(http.StatusSeeOther, "/login")
    })

    router.GET("/login", func(context *gin.Context) {
        context.File("../frontend/static/login.html")
    })
}
