package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"rapstartup/auth"
	"rapstartup/campaign"
	"rapstartup/handler"
	"rapstartup/helper"
	"rapstartup/user"
	webHandler "rapstartup/web/handler"
	"strings"
)

func main() {

	dsn := os.Getenv("DSN")
	if dsn == "" {
		//dsn := "rapdeploy:rahasia1234@tcp(127.0.0.1:3306)/rapstartup?charset=utf8mb4&parseTime=True&loc=Local"
		//dsn := "root:@tcp(127.0.0.1:3306)/rapstartup?charset=utf8mb4&parseTime=True&loc=Local"
		dsn = "root:@tcp(rapstartup-mariadb:3306)/rapstartup?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandlder := webHandler.NewCampaignHandler(campaignService, userService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}

	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("rapstartup", cookieStore))

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users", userWebHandler.Create)

	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandlder.Index)
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHandlder.New)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandlder.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandlder.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandlder.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHandlder.Edit)
	router.POST("/campaigns/update/:id", authAdminMiddleware(), campaignWebHandlder.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHandlder.Show)
	//router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)
	router.GET("/logout", sessionWebHandler.Destroy)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
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
