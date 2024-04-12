package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func initDocsValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Docs{
		Name: "doc 1",
	}
	dbmap.Create(&a1)
	a2 := Docs{
		Name: "doc 2",
	}
	dbmap.Create(&a2)
	return
}

func TestAccessDoc(t *testing.T) {
	defer deleteFile(config.DBname)

	initDocsValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/docs"
	router.POST(urla, PostDoc)
	router.GET(urla, GetDocs)
	router.GET(urla+"/:id", GetDoc)
	router.DELETE(urla+"/:id", DeleteDoc)
	router.PUT(urla+"/:id", UpdateDoc)

	b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET all Docs")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as []Docs
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one Doc")
	var a1 Docs
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	assert.Equal(t, as[0].Name, a1.Name, "a1 = a")

	// Add
	log.Println("= http POST Doc")
	var a = Docs{Name: "Name test"}
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http POST denied")

	// Delete one
	log.Println("= http DELETE one Doc")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Update one
	log.Println("= http PUT one Doc")
	a.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http PUT denied")

	// Change user
	userauth2 := AuthInfo{Role: "guest"}
	r := gin.New()
	r.Use(SetConfig(config, userauth2))
	r.GET(urla, GetDocs)
	// Get all
	log.Println("= http GET all Docs with an other role")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []Docs
	json.Unmarshal(resp.Body.Bytes(), &as2)
	assert.Equal(t, 2, len(as2), "2 results")
	//fmt.Println(as2)
}

func TestDoc(t *testing.T) {
	defer deleteFile(config.DBname)

	InitDb(config.DBname, config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/docs"
	router.POST(urla, PostDoc)
	router.GET(urla, GetDocs)
	router.GET(urla+"/:id", GetDoc)
	router.DELETE(urla+"/:id", DeleteDoc)
	router.PUT(urla+"/:id", UpdateDoc)

	b := new(bytes.Buffer)
	// Add
	log.Println("= http POST Doc")
	var a = Docs{Name: "Name test"}
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", urla, b)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second doc
	log.Println("= http POST more Doc")
	var a2 = Docs{Name: "Name test2"}
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = Docs{}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

	// Get all
	log.Println("= http GET all Docs")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Docs
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	log.Println("= Test parsing query")
	s := "http://127.0.0.1:8080/api?name=t&_order=ASC&_sort=created_on"
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query, s, _ := ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, " (name LIKE \"%t%\")", query, "Parse query")
	assert.Equal(t, " ORDER BY created_on ASC", s, "Parse query")

	// Get one
	log.Println("= http GET one Doc")
	var a1 Docs
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Doc")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE success")

	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all for count success")
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 1, len(as), "1 result")

	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")

	// Update one
	log.Println("= http PUT one Doc")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var a3 Docs
	req, _ = http.NewRequest("GET", urla+"/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	assert.Equal(t, a2.Name, a3.Name, "a2 Name updated")

	req, _ = http.NewRequest("PUT", urla+"/1", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "Can't update missing /1")

	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "Can't update missing mandatory field in /2")

}
