package router

import (
	"encoding/gob"
	"github.com/KodeKunstenaars/Share2Teach/internal/authenticator"
	"github.com/KodeKunstenaars/Share2Teach/internal/middleware"
	"github.com/KodeKunstenaars/Share2Teach/web/app/callback"
	"github.com/KodeKunstenaars/Share2Teach/web/app/login"
	"github.com/KodeKunstenaars/Share2Teach/web/app/upload"
	"github.com/KodeKunstenaars/Share2Teach/web/app/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"net/http"

	"github.com/KodeKunstenaars/Share2Teach/web/app/logout"
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", middleware.IsAuthenticated, user.Handler)
	router.GET("/logout", logout.Handler)

	router.GET("/upload", middleware.IsAuthenticated, upload.Handler())

	return router
}
