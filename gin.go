package authentic

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

// CreateGinMiddleware generates a configured Gin Middleware for validating JWTs
func (m *middlewareCreator) CreateGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		rx := regexp.MustCompile("^[B|b]earer\\s*")
		jwt := rx.ReplaceAllString(authHeader, "")
		if !m.Validator.IsValid(jwt) {
			if m.FailureHook != nil {
				m.FailureHook(c)
				return
			}

			c.AbortWithStatusJSON(401, m.notAuthorizedError())
			return
		}

		if m.SuccessHook != nil {
			m.SuccessHook(c)
		}
	}
}
