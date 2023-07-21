package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "bwastartup/web/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database is good")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	
	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.100.2"})
	router.Use(cors.Default())

	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("mrsastartup", cookieStore))

	router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")
	
	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/image", "./web/assets/image")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)
	api.GET("/campaigns",campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id",campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService),campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService),campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService),campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users", userWebHandler.Create)

	router.GET("/users/edit/:id", userWebHandler.Edit)
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.Update)

	router.GET("/users/avatar/:id", authAdminMiddleware(), userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)

	router.GET("/campaigns", authAdminMiddleware(),campaignWebHandler.Index)
	router.GET("/campaigns/new", authAdminMiddleware(),campaignWebHandler.New)
	router.POST("/campaigns", authAdminMiddleware(),campaignWebHandler.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(),campaignWebHandler.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(),campaignWebHandler.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(),campaignWebHandler.Edit)
	router.POST("/campaigns/update/:id", authAdminMiddleware(),campaignWebHandler.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(),campaignWebHandler.Show)
	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login",sessionWebHandler.New)
	router.POST("/session",sessionWebHandler.Create)
	router.GET("/logout",sessionWebHandler.Destroy)

	log.Fatal(router.Run(":8080"))


	// gambaran struktur flow:
		// request
		// handler
		// Service
		// repository
		// database

	// gambaran flow register:
		// input
		// handler mapping input dari user -> struct input
		// Service mapping dari struct input ke struct User
		// repository save struct user ke db
		// database
	
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
	
		if !strings.Contains(authHeader, "Bearer"){
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
	
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
	
		token, err := authService.ValidateToken(tokenString)

		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
		c.Set("currentUser",user)
	}
}
// ambil nilai header authorization: bearer tokentokentoken
// dari header authorization, kita ambil nilai tokkennya saja
// kita validasi token
// kita ambil user_id
// ambil user dari db berdasarkan user_id lewat service
// kalau user ada set context isinya user

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// userID
		userIDSession := session.Get("userID")
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
  r := multitemplate.NewRenderer()

  layouts, err := filepath.Glob(templatesDir + "/layouts/*")
  if err != nil {
    panic(err.Error())
  }

  includes, err := filepath.Glob(templatesDir + "/**/*")
  if err != nil {
    panic(err.Error())
  }

  // Generate our templates map from our layouts/ and includes/ directories
  for _, include := range includes {
    layoutCopy := make([]string, len(layouts))
    copy(layoutCopy, layouts)
    files := append(layoutCopy, include)
    r.AddFromFiles(filepath.Base(include), files...)
  }
  return r
}