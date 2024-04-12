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

func initUsersValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Users{
		CasID: "user1",
	}
	dbmap.Create(&a1)
	a2 := Users{
		CasID: "user2",
	}
	dbmap.Create(&a2)

	da1 := Domaine{
		Name:  "domaine 1",
		User1: int64(a1.ID),
		User2: int64(a1.ID),
		User3: int64(a2.ID),
	}
	dbmap.Create(&da1)
	da2 := Domaine{
		Name:  "domaine 2",
		User2: int64(a1.ID),
	}
	dbmap.Create(&da2)
	return
}

func TestAccessUser(t *testing.T) {
	defer deleteFile(config.DBname)

	initUsersValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/users"
	router.POST(urla, PostUser)
	router.GET(urla, GetUsers)
	router.GET(urla+"/:id", GetUser)
	router.DELETE(urla+"/:id", DeleteUser)
	router.PUT(urla+"/:id", UpdateUser)

	b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET all Users")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http success")
	/*var as []Users
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 2, len(as), "2 results")*/

	// Get one
	log.Println("= http GET one User")
	var a1 Users
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//assert.Equal(t, as[0].Name, a1.Name, "a1 = a")

	// Add
	log.Println("= http POST User")
	var a = Users{CasID: "user1"}
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http POST denied")

	// Delete one
	log.Println("= http DELETE one User")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Update one
	log.Println("= http PUT one User")
	a.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http PUT denied")

	// Create new user by login
	log.Println("= Create new user by login")
	var u = Users{}
	e := u.ByLogin(config.DBh, "newlogin")
	assert.Equal(t, nil, e, "New user by login whithout error")
	assert.Equal(t, int32(3), u.ID, "New user by login whith ID 3")

	// Change user
	userauth2 := AuthInfo{Role: "admin"}
	r := gin.New()
	r.Use(SetConfig(config, userauth2))
	r.GET(urla, GetUsers)
	r.GET(urla+"/:id", GetUser)
	// Get all
	log.Println("= http GET all Users with an other role")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []Users
	json.Unmarshal(resp.Body.Bytes(), &as2)
	assert.Equal(t, 3, len(as2), "3 results")

	// Get one
	log.Println("= http GET one User")
	//var a1 Users
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(prettyPrint(a1))
	assert.Equal(t, 2, len(a1.Doms), "user 1 manage 2 domains")
	assert.Equal(t, int32(2), a1.Doms[1], "user 1 manage domaine 2 ")
}

func TestUser(t *testing.T) {
	defer deleteFile(config.DBname)

	InitDb(config.DBname, config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/users"
	router.POST(urla, PostUser)
	router.GET(urla, GetUsers)
	router.GET(urla+"/:id", GetUser)
	router.DELETE(urla+"/:id", DeleteUser)
	router.PUT(urla+"/:id", UpdateUser)

	b := new(bytes.Buffer)
	// Add
	log.Println("= http POST User")
	var a = Users{CasID: "user1"}
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

	// Add second user
	log.Println("= http POST more User")
	var a2 = Users{CasID: "user2"}
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = Users{}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

	// Get all
	log.Println("= http GET all Users")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Users
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	log.Println("= Test parsing query")
	s := "http://127.0.0.1:8080/api?name=t&id=1&id=2&_order=ASC&_sort=created_on"
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query, s, _ := ParseQuery(q)
	//fmt.Println(query)
	//assert.Equal(t, ` (id = "1" OR id = "2") AND (name LIKE "%t%")`, query, "Parse query")
	assert.Equal(t, ` (name LIKE "%t%") AND (id = "1" OR id = "2")`, query, "Parse query")
	assert.Equal(t, " ORDER BY created_on ASC", s, "Parse query")

	s = "http://127.0.0.1:8080/api?casid=t&axe=1&descriptions=t&_sort=casid&_order=DESC&_end=120&_start=90"
	u, _ = url.Parse(s)
	q, _ = url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query, s, l := ParseQuery(q)
	//fmt.Println(query)
	//assert.Equal(t, ` (cas_id LIKE "%t%") AND (axe1 = "1" OR axe2 = "1") AND (descorig LIKE "%t%" OR description LIKE "%t%")`, query, "Parse query")
	//assert.Equal(t, ` (descorig LIKE "%t%" OR description LIKE "%t%") AND (cas_id LIKE "%t%") AND (axe1 = "1" OR axe2 = "1")`, query, "Parse query")
	assert.Equal(t, ` (cas_id LIKE "%t%") AND (axe1 = "1" OR axe2 = "1") AND (descorig LIKE "%t%" OR description LIKE "%t%")`, query, "Parse query")
	assert.Equal(t, " ORDER BY cas_id DESC", s, "Parse query")
	assert.Equal(t, ` LIMIT 30 OFFSET 90`, l, "Parse query")

	// Get one
	log.Println("= http GET one User")
	var a1 Users
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one User")
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
	log.Println("= http PUT one User")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var a3 Users
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
