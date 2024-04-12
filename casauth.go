package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	goCas "github.com/go-cas/cas"
	"github.com/jinzhu/gorm"

	. "model"
)

type casMiddleware struct {
	casClient *goCas.Client
	handler   http.Handler
}

func (casMiddleware casMiddleware) authed(c *gin.Context) {
	c.Set("CASUsername", goCas.Username(c.Request))
	//c.Set("CASAttributes", goCas.Attributes(c.Request))
	c.Next()
	return
}

func (casMiddleware casMiddleware) middlewareFunc(c *gin.Context) {
	casMiddleware.handler.ServeHTTP(c.Writer, c.Request)
	c.Header("Cache-Control", "no-cache, private, max-age=0")
	//c.Header("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	c.Header("Pragma", "no-cache")
	c.Header("X-Accel-Expires", "0")
	if goCas.IsAuthenticated(c.Request) {
		casMiddleware.authed(c)
		return
	}
	c.Abort()
}

type CasOptions = goCas.Options

func CasMiddlewareFunc(options *CasOptions) gin.HandlerFunc {
	casClient := goCas.NewClient((*goCas.Options)(options))
	log.Println("goCas")
	rawHandler := func(res http.ResponseWriter, req *http.Request) {
		if goCas.IsAuthenticated(req) {
			log.Println("goCas.IsAuthenticated")
			log.Println(goCas.Username(req))
			return
		}
		/*if goCas.IsNewLogin(req) {
		                        log.Println("goCas.IsNewLogin")
		                        log.Println(goCas.Username(req))
					return
				}*/
		casClient.RedirectToLogin(res, req)
	}
	return casMiddleware{
		casClient: casClient,
		handler:   casClient.HandleFunc(rawHandler),
	}.middlewareFunc
}

/*
  Auth jwt
*/

func CreateMiddlware(config Config) *jwt.GinJWTMiddleware {
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.JwtPass),
		Timeout:     3 * time.Hour,
		MaxRefresh:  3 * time.Hour,
		IdentityKey: "IDuser",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			//log.Println("  data: PayloadFunc")
			//log.Printf("%+v\n", data)
			if v, ok := data.(Users); ok {
				//log.Println(v.CasID)
				//log.Println(v.Name)
				//log.Println(v.UserRole)
				return jwt.MapClaims{
					"IDuser":   strconv.FormatInt(int64(v.ID), 10),
					"Name":     v.CasID,
					"FullName": v.Name,
					"Role":     v.UserRole,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			username, err := c.Get("CASUsername")
			//log.Println("   jwt ==Auth== Authenticator")
			//log.Println("   ", username)
			var user Users
			dbmap := c.MustGet("DBmap").(*gorm.DB)
			user.ByLogin(dbmap, username.(string))
			//log.Printf("%+v\n",user)
			if err != false {
				log.Println("CASUsername :", user.CasID)
				return user, nil
			}
			return user, jwt.ErrFailedAuthentication
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			//log.Println("   jwt ==Auth== IdentityHandler")
			claims := jwt.ExtractClaims(c)
			i, _ := strconv.ParseInt(claims["IDuser"].(string), 10, 64)
			c.Set("LoginID", i)
			c.Set("Login", claims["Name"])
			c.Set("Role", claims["Role"])
			/*return &User{
				UserName: claims[identityKey].(string),
			}*/
			//log.Printf("%+v\n", claims)
			return nil
		},
		Authorizator: func(userId interface{}, c *gin.Context) bool {
			//log.Println("   jwt ==Auth== Authorizator")
			jwtClaims := jwt.ExtractClaims(c)
			/*log.Println("userId: ", userId)
			log.Println("jwtClaims Name: ", jwtClaims["Name"])*/
			if jwtClaims == nil {
				return false
			}
			if Contains([]string{"admin", "cssi", "reader"}, jwtClaims["Role"].(string)) {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			log.Println(message)
			//location := url.URL{Path: config.CallbackURL+"?"+message}
			location, _ := url.Parse(config.CallbackURL + "?" + message)
			fmt.Println(location.String())
			//c.SetCookie("token",message,10, "/", "localhost", false, false ) //true, true) //secure bool,httpOnly bool,
			c.Header("Cache-Control", "no-cache, private, max-age=0")
			//c.Header("Expires", time.Unix(0, 0).Format(http.TimeFormat))
			c.Header("Pragma", "no-cache")
			c.Header("X-Accel-Expires", "0")
			c.Redirect(301, location.String())
			//c.Redirect(301, location.RequestURI())
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
	return authMiddleware
}
