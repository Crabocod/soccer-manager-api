package middleware

import (
	"soccer_manager_service/pkg/i18n"

	"github.com/gin-gonic/gin"
)

const LocalizerKey = "localizer"

func I18nMiddleware(i18nManager *i18n.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLanguage := c.GetHeader("Accept-Language")
		localizer := i18nManager.GetLocalizer(acceptLanguage)

		c.Set(LocalizerKey, localizer)
		c.Next()
	}
}
