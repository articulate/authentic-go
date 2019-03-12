package authentic

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// CreateGinMiddleware generates a configured Gin Middleware for validating JWTs
func (m *middlewareCreator) CreateGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		rx := regexp.MustCompile("^[B|b]earer\\s*")
		jwt := rx.ReplaceAllString(authHeader, "")
		result := m.Validator.ValidateToken(jwt)

		if !result.Valid || result.Expired {
			if m.FailureHook != nil {
				m.FailureHook(c)
				return
			}

			if result.Expired {
				c.AbortWithStatusJSON(http.StatusForbidden, m.notAuthorizedError())
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, m.notAuthorizedError())
			return
		}

		if m.SuccessHook != nil {
			m.SuccessHook(c)
		}
	}
}
