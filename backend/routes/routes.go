package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "database/sql"
    "github.com/gin-contrib/sessions"
    "golang.org/x/crypto/bcrypt"
)

func HandleRoutes(router *gin.Engine, database *sql.DB) {
    router.GET("/", func(context *gin.Context) {
        session := sessions.Default(context)
        username := session.Get("username")

        context.HTML(http.StatusOK, "index.html", gin.H {
            "isLoggedIn": username != nil,
        })
    })

    router.GET("/register", func(context *gin.Context) {
        session := sessions.Default(context)
        username := session.Get("username")

        context.HTML(http.StatusOK, "register.html", gin.H {
            "isLoggedIn": username != nil,
        })
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
        session := sessions.Default(context)
        username := session.Get("username")

        context.HTML(http.StatusOK, "login.html", gin.H {
            "isLoggedIn": username != nil,
        })
    })

    router.POST("/login", func(context *gin.Context) {
        usernameInput := context.PostForm("username-input")
        passwordInput := context.PostForm("password-input")

        var storedPassword string
        err := database.QueryRow("SELECT password FROM users WHERE username = ?", usernameInput).Scan(&storedPassword)
        if err != nil {
            if err == sql.ErrNoRows {
                context.JSON(http.StatusUnauthorized, gin.H{"error": "Username does not exist"})
                return
            }
            context.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(passwordInput))
        if err != nil {
            context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
            return
        }

        session := sessions.Default(context)
        session.Set("username", usernameInput)
        session.Save()

        context.Redirect(http.StatusSeeOther, "/profile/" + usernameInput)
    })

    router.POST("/logout", func(context *gin.Context) {
        session := sessions.Default(context)
        session.Clear()
        session.Save()

        context.SetCookie("user_session", "", -1, "/", "", false, true)

        context.Redirect(http.StatusSeeOther, "/")
    })

    router.GET("/profile/:username", func(context *gin.Context) {
        requestedUsername := context.Param("username")

        var count int
        err := database.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", requestedUsername).Scan(&count)
        if err != nil {
            context.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        if count == 1 {
            session := sessions.Default(context)
            username := session.Get("username")

            context.HTML(http.StatusOK, "profile.html", gin.H {
                "username": requestedUsername,
                "isLoggedIn": username != nil,
            })
        } else {
            context.String(http.StatusNotFound, "User not found")
        }
    })
}
