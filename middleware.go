package authentic

import (
	"github.com/gin-gonic/gin"
)

type (
	// ErrorResponse represents the default not authorized error
	ErrorResponse struct {
		Message string   `json:"message"`
		Causes  []string `json:"causes"`
	}

	middlewareCreator struct {
		Validator   Validator
		FailureHook gin.HandlerFunc
		SuccessHook gin.HandlerFunc
	}
)

// OnFailure setup call when validation fails, as opposed to using the default 401 JSON response
func (m *middlewareCreator) OnFailure(hook gin.HandlerFunc) MiddlewareCreator {
	m.FailureHook = hook
	return m
}

// OnSuccess setup call when validation works
func (m *middlewareCreator) OnSuccess(hook gin.HandlerFunc) MiddlewareCreator {
	m.SuccessHook = hook
	return m
}

func (m *middlewareCreator) notAuthorizedError() *ErrorResponse {
	return &ErrorResponse{
		Message: "Unauthorized",
		Causes:  []string{"Invalid session"},
	}
}
