// web/app/upload/upload.go

package upload

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler serves the upload page.
func Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		profile := session.Get("profile")

		if profile == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		// Render the upload.html page with the profile data.
		ctx.HTML(http.StatusOK, "upload.html", profile)
	}
}
