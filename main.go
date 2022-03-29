package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rapstartup/auth"
	"rapstartup/handler"
	"rapstartup/user"

	webHandler "rapstartup/web/handler"
)

func main() {
	dsn := "rapdeploy:rahasia1234@tcp(127.0.0.1:3306)/rapstartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	userWebHandler := webHandler.NewUserHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("rapstartup", cookieStore))

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)

	router.Run()
}

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIDSession := session.Get("userID")

		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}
