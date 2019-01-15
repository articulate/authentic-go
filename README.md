# Authentic

Golang clone of [@authentic/authentic](https://github.com/articulate/authentic) and [authentic-rb](https://github.com/articulate/authentic-rb). A simple library to validate JWTs against issuer's JWK.

## Installation

``` bash
go get -u github.com/articulate/authentic-go
```

## Usage

*Note* `Validator` caches keys in memory and a new instance results in a cleared cache.

| Provider | Sample `iss_whitelist` |
| -------- | ------------------- |
| [Auth0](https://auth0.com/) | `[ 'https://${tenant}.auth0.com/' ]` |
| [Okta](https://www.okta.com/) | `[ 'https://${tenant}.oktapreview.com/oauth2/${authServerId}' ]` |

```golang
import (
  "fmt"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/articulate/authentic-go"
)

func main() {
  // Create validator with ISS from environment var AUTHENTIC_ISS_WHITELIST and use default cache max age
  validator := authentic.NewValidator()

  // Or create validator specifying valid ISS and max JWK cache age
  validator = authentic.NewValidator().
      WithWhitelist("https://org.auth0.com/", "'https://org.okta.com/").
      WithCacheMaxAge(time.Hour * 4)

  // Validate token
  if !validator.IsValid("some jwt from somewhere") {
    fmt.Println("Your token is invalid! Get on out of here!")
  }

  // Let's create some Gin middleware with configured validator. When validation fails, a 401 response generated and
  // route is short circuited.
  router := gin.New()
  api := r.Group("/super/secret")

  middlewareOptions := &authentic.middlewareOptions{}
  api.Use(validator.CreateGinMiddleware(middlewareOptions))
}
```

## Options

| Name              | Default         | Required |
| ----------------- | --------------- | -------- |
| ISSWhitelist      | N/A             | y        |
| CacheMaxAge       | 10 * time.Hours | n        |
