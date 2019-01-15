package authentic_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"gopkg.in/h2non/gock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/articulate/authentic-go"
)

const (
	wellKnown = "/.well-known/openid-configuration"
	testUrl   = "https://authentic.articulate.com"

	badIss = "eyJraWQiOiJEYVgxMWdBcldRZWJOSE83RU1QTUw1VnRUNEV3cmZrd2M1U2xHaVd2VXdBIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMH0V" +
		"kanlqc3NidDJTMVFWcjBoNyIsInZlciI6MSwiaXNzIjoiaHR0cHM6Ly9iYWQtaXNzLmNvbSIsImF1ZCI6IjBvYWRqeWs1MjNobFpmeWIxMGg3IiwiaWF0" +
		"IjoxNTE2NjM3MDkxLCJleHAiOjE1MTY2NDA2OTEsImp0aSI6IklELmM4amh6b2t5MGZGTlByOExfU0NycnBnVFRVeUFvY3RIdjY5T0tTbWY1R0EiLCJhb" +
		"XIiOlsicHdkIl0sImlkcCI6IjAwb2NnNHRidTZGSzJEaDVHMGg3Iiwibm9uY2UiOiIyIiwiYXV0aF90aW1lIjoxNTE2NjM3MDkxLCJ0ZW5hbnRJZCI6Im" +
		"Q0MmUzM2ZkLWYwNWUtNGE0ZS05MDUwLTViN2IyZTgwMDgzNCJ9.Senilj3Z8Z99b-UVnnxwWKjYIn4jNrE-BmZAuR7Qb3nkxS7N-r7WnAQ-4vuqtD5Fyy" +
		"s1zOFUxoO6jyMvhWbhNlPmYaBQk7InKZU6ABayrijfv7OJSQKzs0Q7EQbgtW4T27Gqp6G4Rp9l7O472lgwapTV_L2IUqYNP7aC3FAFcqmpP_KFyeKj-zc" +
		"wil6aszPgxzMA3Rp33BqQfuhIJKSYqWQT6pkDXkjM3pLxaHRfrRahQ2F0M190iCvBJMc4b82TVoQQu5uJbb1mD97wwlSvMFYCHN_51g9IY5BabZcOv4h0" +
		"T3-XqFxPNbS8PZVfBikumkhqD5b4zjA-3ddgPw2GkA"

	token = "eyJraWQiOiJEYVgxMWdBcldRZWJOSE83RU1QTUw1VnRUNEV3cmZrd2M1U2xHaVd2VXdBIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMHVkanl" +
		"qc3NidDJTMVFWcjBoNyIsInZlciI6MSwiaXNzIjoiaHR0cHM6Ly9hdXRoZW50aWMuYXJ0aWN1bGF0ZS5jb20vIiwiYXVkIjoiMG9hZGp5azUyM2hsWmZ5" +
		"YjEwaDciLCJpYXQiOjE1MTY2MzcwOTEsImV4cCI6MTUxNjY0MDY5MSwianRpIjoiSUQuYzhqaHpva3kwZkZOUHI4TF9TQ3JycGdUVFV5QW9jdEh2NjlPS" +
		"1NtZjVHQSIsImFtciI6WyJwd2QiXSwiaWRwIjoiMDBvY2c0dGJ1NkZLMkRoNUcwaDciLCJub25jZSI6IjIiLCJhdXRoX3RpbWUiOjE1MTY2MzcwOTEsIn" +
		"RlbmFudElkIjoiZDQyZTMzZmQtZjA1ZS00YTRlLTkwNTAtNWI3YjJlODAwODM0In0.NEVqz-jJIyaEgho3uQYOvWC52s_50AV--FHwBWm9BftucQ5G4bS" +
		"HL7szeaPc3HT0VrhFUntRLlJHzw7pZvRJG2WExj6HJi-Ug3LDwQOj47Gf_ywlEydBAQz7u98JK2ZJcCP16-lIOM1J-fUz-SpFqI4RcO5MLiiEPnMqsXS-" +
		"EkPd8Y27G64PnHnNjaY3sLrOc9peeD5Xh82TSjeMFFAPpiYNtTCixnfZeQCCtxOCPhiDYAwDSxaLbrOcDAYdO0ytKQ9dBfFoY0AzJNqgJUOPVeeC_AgEJ" +
		"eLIaSKVJAFqZAB8t5VagvVGIqcu7TaMCOmOZx_5A8Xc9JVmRoKDAMlizQ"

	badToken = "cyJraWQiOiJEYVgxMWdBcldRZWJOSE83RU1QTUw1VnRUNEV3cmZrd2M1U2xHaVd2VXdBIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIwMHV" +
		"kanlqc3NidDJTMVFWcjBoNyIsInZlciI6MSwiaXNzIjoiaHR0cHM6Ly9hdXRoZW50aWMuYXJ0aWN1bGF0ZS5jb20vIiwiYXVkIjoiMG9hZGp5azUyM2hs" +
		"WmZ5YjEwaDciLCJpYXQiOjE1MTY2MzcwOTEsImV4cCI6MTUxNjY0MDY5MSwianRpIjoiSUQuYzhqaHpva3kwZkZOUHI4TF9TQ3JycGdUVFV5QW9jdEh2N" +
		"jlPS1NtZjVHQSIsImFtciI6WyJwd2QiXSwiaWRwIjoiMDBvY2c0dGJ1NkZLMkRoNUcwaDciLCJub25jZSI6IjIiLCJhdXRoX3RpbWUiOjE1MTY2MzcwOT" +
		"EsInRlbmFudElkIjoiZDQyZTMzZmQtZjA1ZS00YTRlLTkwNTAtNWI3YjJlODAwODM0In0.NEVqz-jJIyaEgho3uQYOvWC52s_50AV--FHwBWm9BftucQ5" +
		"G4bSHL7szeaPc3HT0VrhFUntRLlJHzw7pZvRJG2WExj6HJi-Ug3LDwQOj47Gf_ywlEydBAQz7u98JK2ZJcCP16-lIOM1J-fUz-SpFqI4RcO5MLiiEPnMq" +
		"sXS-EkPd8Y27G64PnHnNjaY3sLrOc9peeD5Xh82TSjeMFFAPpiYNtTCixnfZeQCCtxOCPhiDYAwDSxaLbrOcDAYdO0ytKQ9dBfFoY0AzJNqgJUOPVeeC_" +
		"AgEJeLIaSKVJAFqZAB8t5VagvVGIqcu7TaMCOmOZx_5A8Xc9JVmRoKDAMlizQ"
)

var (
	oidc      interface{}
	keys      interface{}
	validator authentic.Validator
)

func jsonFixture(file string) interface{} {
	var body interface{}
	fixturePath, _ := filepath.Abs(path.Join("fixtures", file))
	data, _ := ioutil.ReadFile(fixturePath)
	json.Unmarshal(data, &body)

	return body
}

var _ = Describe("authentic", func() {
	Describe("Validator", func() {
		BeforeEach(func() {
			oidc = jsonFixture("oidc.json")
			keys = jsonFixture("keys.json")
			gock.New(testUrl).
				Times(1).
				Get(wellKnown).
				Reply(200).
				JSON(oidc)
			gock.New(testUrl).
				Times(1).
				Get("/v1/keys").
				Reply(200).
				JSON(keys)
			validator = authentic.NewValidator().
				WithWhitelist("https://org.auth0.com/", "https://org.okta.com/")
		})

		It("validates JWT against JWK", func() {
			Expect(validator.IsValid(token)).To(BeTrue())
		})

		It("fails to validate bad token", func() {
			Expect(validator.IsValid(badToken)).To(BeFalse())
		})

		It("fails to validate valid token with wrong iss", func() {
			Expect(validator.IsValid(badIss)).To(BeFalse())
		})

		It("caches key and only makes one request", func() {
			Expect(validator.IsValid(token)).To(BeTrue())
			Expect(validator.IsValid(token)).To(BeTrue())
		})

		It("retrieves key after cache is stale", func() {
			gock.New(testUrl).
				Times(2).
				Get(wellKnown).
				Reply(200).
				JSON(oidc)
			gock.New(testUrl).
				Times(2).
				Get("/v1/keys").
				Reply(200).
				JSON(keys)
			validator = authentic.NewValidator().WithCacheMaxAge(time.Microsecond)
			Expect(validator.IsValid(token)).To(BeTrue())
			time.Sleep(time.Microsecond)
			Expect(validator.IsValid(token)).To(BeTrue())
		})

		It("serves stale when request fails, but tries again subsequentially", func() {
			validator = authentic.NewValidator().WithCacheMaxAge(time.Microsecond)
			Expect(validator.IsValid(token)).To(BeTrue())
			gock.New(testUrl).
				Times(1).
				Get(wellKnown).
				Reply(500)
			time.Sleep(time.Microsecond)
			Expect(validator.IsValid(token)).To(BeTrue())
			gock.New(testUrl).
				Get(wellKnown).
				Reply(200).
				JSON(oidc)
			Expect(validator.IsValid(token)).To(BeTrue())
		})

		It("returns false when no OIDC config is returned", func() {
			gock.Flush()
			gock.New(testUrl).
				Times(1).
				Get(wellKnown).
				Reply(500)
			Expect(validator.IsValid(token)).To(BeFalse())
		})

		It("returns false when no keys is returned", func() {
			gock.Flush()
			gock.New(testUrl).
				Times(1).
				Get(wellKnown).
				Reply(200).
				JSON(oidc)
			gock.New(testUrl).
				Times(1).
				Get("/v1/keys").
				Reply(500)
			Expect(validator.IsValid(token)).To(BeFalse())
		})

		Context("middleware", func() {
			var (
				body        *authentic.ErrorResponse
				rec         *httptest.ResponseRecorder
				mockContext *gin.Context
			)

			BeforeEach(func() {
				body = &authentic.ErrorResponse{}
				rec = httptest.NewRecorder()
				mockContext, _ = gin.CreateTestContext(rec)
				mockContext.Request = &http.Request{
					Header: http.Header{
						"Authorization": []string{"Bearer " + token},
					},
				}
			})

			It("returns 401 response in Gin middleware", func() {
				mockContext.Request.Header["Authorization"] = []string{"Bearer " + badToken}
				authentic.NewMiddlewareCreator().CreateGinMiddleware()(mockContext)
				json.NewDecoder(rec.Body).Decode(&body)
				Expect(rec.Code).To(Equal(401))
				Expect(body.Message).To(Equal("Unauthorized"))
			})

			It("does not respond with a valid token", func() {
				authentic.NewMiddlewareCreator().CreateGinMiddleware()(mockContext)
				Expect(rec.Code).To(Equal(200))
			})
		})
	})
})
