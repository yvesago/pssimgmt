package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"

	. "model"
)

// SetConfig: gin Middlware to push some config values
func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Port", config.Port)
		//c.Set("CorsOrigin", config.CorsOrigin)
		c.Set("Verbose", config.Verbose)
		c.Set("DBmap", config.DBh)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-total-count")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func servermain(config Config) {
	if config.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	if config.Debug == true {
		r.Use(gin.Logger())
	}

	r.Use(SetConfig(config))
	//r.Use(Databases(config))
	//        r.Use(CORS())
	/*	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Cache-Control", "public, max-age=10")
		}
	}())*/

	r.Use(cors.Middleware(cors.Config{
		Origins:         config.CorsOrigin,
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Access-Control-Allow-Methods, Access-Control-Allow-Headers, Access-Control-Allow-Origin, Origin, Authorization, Content-Type, X-MyToken, Bearer",
		ExposedHeaders:  "x-total-count",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	authMiddleware := CreateMiddlware(config)
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	authURL, err := url.Parse(config.AuthURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	casOptions := CasOptions{
		URL:         authURL,
		SendService: true,
	}

	casAuth := CasMiddlewareFunc(&casOptions)
	auth := r.Group("/cas")
	auth.Use(casAuth)
	{
		auth.GET("/login", authMiddleware.LoginHandler)
		auth.GET("/logout", func(c *gin.Context) { c.Next() })
	}

	//r.GET("/refresh-token", authMiddleware.RefreshHandler)
	admin := r.Group("api/v1")
	//admin.Use()
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		admin.GET("/todos", GetTodos)
		admin.OPTIONS("/todos", Options)     // POST

		admin.GET("/themestree/:dom", GetThemesTree)
		admin.GET("/themes", GetThemes)
		admin.GET("/themes/:id", GetTheme)
		admin.GET("/themes/:id/:dom", GetThemeByDom)
		admin.POST("/themes", PostTheme)
		admin.PUT("/themes/:id", UpdateTheme)
		admin.DELETE("/themes/:id", DeleteTheme)
		admin.OPTIONS("/themes", Options)     // POST
		admin.OPTIONS("/themes/:id", Options) // PUT, DELETE

		admin.GET("/regles", GetRegles)
		admin.GET("/regles/:id", GetRegle)
		admin.GET("/regles/:id/:dom", GetRegleByDom)
		admin.PUT("/regles/:id/:dom", UpdateRegleByDom)
		admin.POST("/regles", PostRegle)
		admin.PUT("/regles/:id", UpdateRegle)
		admin.DELETE("/regles/:id", DeleteRegle)
		admin.OPTIONS("/regles", Options)          // POST
		admin.OPTIONS("/regles/:id", Options)      // PUT, DELETE
		admin.OPTIONS("/regles/:id/:dom", Options) // PUT, DELETE

		admin.GET("/domainestree", GetDomainesTree)
		admin.GET("/domaines", GetDomaines)
		admin.GET("/domaines/:id", GetDomaine)
		admin.POST("/domaines", PostDomaine)
		admin.PUT("/domaines/:id", UpdateDomaine)
		admin.DELETE("/domaines/:id", DeleteDomaine)
		admin.OPTIONS("/domaines", Options)     // POST
		admin.OPTIONS("/domaines/:id", Options) // PUT, DELETE

		admin.GET("/users", GetUsers)
		admin.GET("/users/:id", GetUser)
		admin.POST("/users", PostUser)
		admin.PUT("/users/:id", UpdateUser)
		admin.DELETE("/users/:id", DeleteUser)
		admin.OPTIONS("/users", Options)     // POST
		admin.OPTIONS("/users/:id", Options) // PUT, DELETE

		admin.GET("/docs", GetDocs)
		admin.GET("/docs/:id", GetDoc)
		admin.POST("/docs", PostDoc)
		admin.PUT("/docs/:id", UpdateDoc)
		admin.DELETE("/docs/:id", DeleteDoc)
		admin.OPTIONS("/docs", Options)     // POST
		admin.OPTIONS("/docs/:id", Options) // PUT, DELETE

		admin.GET("/isos", GetIso27002s)
		admin.GET("/isos/:id", GetIso27002)
		admin.POST("/isos", PostIso27002)
		admin.PUT("/isos/:id", UpdateIso27002)
		admin.DELETE("/isos/:id", DeleteIso27002)
		admin.OPTIONS("/isos", Options)     // POST
		admin.OPTIONS("/isos/:id", Options) // PUT, DELETE

		admin.GET("/versions", GetVersions)
		admin.GET("/versions/:id", GetVersion)
		admin.POST("/versions", PostVersion)
		admin.PUT("/versions/:id", UpdateVersion)
		admin.DELETE("/versions/:id", DeleteVersion)
		admin.OPTIONS("/versions", Options)     // POST
		admin.OPTIONS("/versions/:id", Options) // PUT, DELETE

		admin.GET("/documents", GetDocuments)
		admin.GET("/documents/:id", GetDocument)
		admin.POST("/documents", PostDocument)
		admin.PUT("/documents/:id", UpdateDocument)
		admin.DELETE("/documents/:id", DeleteDocument)
		admin.OPTIONS("/documents", Options)     // POST
		admin.OPTIONS("/documents/:id", Options) // PUT, DELETE
	}

	/*if config.TLScert != "" && config.TLSkey != "" {
		if config.Debug == true {
			fmt.Println("Listening and serving HTTPS on ", config.Port)
		}
		err := http.ListenAndServeTLS(config.Port, config.TLScert, config.TLSkey, r)
		if err != nil {
			fmt.Println("ListenAndServe: ", err)
			os.Exit(0)
		}
	} else {
		r.Run(config.Port)
	}*/
	r.Run(config.Port)
}

// Options - common response for rest options
func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime) // Add date to logs
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s\n\n  Default behaviour: start daemon\n\n", os.Args[0])
		//flag.SortFlags = false
		//flag.MarkHidden("divIP")
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Usage = Usage

	// Parameters
	confPtr := flag.StringP("conf", "c", "", "*Mandatory* Json config file")
	debugPtr := flag.BoolP("debug", "d", false, "Debug mode")
	verbosePtr := flag.BoolP("verbose", "v", false, "Verbose mode, need Debug mode")
	flag.Parse()
	conf := *confPtr
	Debug := *debugPtr
	Verbose := *verbosePtr

	if Debug == false { // Verbose need Debug
		Verbose = false
	}

	// Load config from file
	file, err := os.Open(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "=========\nError: %s\n=========\n", err)
		Usage()
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "=========\nError: %s\n=========\n", err)
		Usage()
	}
	config.Debug = Debug
	config.Verbose = Verbose

	config = DBs(config)

	var wg sync.WaitGroup
	wg.Add(1)
	go servermain(config)
	wg.Wait()

}
