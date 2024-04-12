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
	"sort"
	"testing"
)

func TestAllDomains(t *testing.T) {
	defer deleteFile(config.DBname)
	initDomaineValues(config.DBname, false) //config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var url = "/admin/api/v1/domains"
	//router.POST(url, PostDomain)
	router.GET(url, GetDomaines)
	router.GET(url+"tree", GetDomainesTree)
	//router.GET(url+"/:id", GetDomain)
	//router.DELETE(url+"/:id", DeleteDomain)
	//router.PUT(url+"/:id", UpdateDomain)

	log.Println("= http GET all Domains")
	//req, err := http.NewRequest("GET", url+"?name=domain&_sort=created_on&_order=DESC&_start=0&_end=10", nil)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	var as []Domaine
	json.Unmarshal(resp.Body.Bytes(), &as)

	fmt.Println(prettyPrint(as))
	elmap := make(map[int32]*Element)
	tree := &Element{Id: 0, Name: "head"}
	for _, e := range as {
		el := &Element{Id: e.ID, Name: e.Name, parent: e.Parent}
		elmap[e.ID] = el
		if e.Parent == 0 {
			tree.Children = append(tree.Children, el)
		}
	}
	for _, e := range elmap {

		if e.parent != 0 {
			e.P = elmap[e.parent]
			elmap[e.parent].Children = append(elmap[e.parent].Children, e)
		}
	}

	//var menu []Entry
	//json.Unmarshal(resp.Body.Bytes(), &menu)
	//fmt.Printf("\n\n%+v\n",menu)
	//fmt.Println(prettyPrint(menu))
	sort.Sort(ByOrder(tree.Children))
	fmt.Println(tree)
	//fmt.Println(prettyPrint(tree))
	log.Println("= http GET Domaines Tree")
	//req, err := http.NewRequest("GET", url+"?name=domain&_sort=created_on&_order=DESC&_start=0&_end=10", nil)
	req, err = http.NewRequest("GET", url+"tree", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	var dtree Domaine
	json.Unmarshal(resp.Body.Bytes(), &dtree)

	fmt.Println(prettyPrint(dtree))

	assert.Equal(t, 5, len(as), "5 results")

}

func initDomaineValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Domaine{
		Name: "domaine 1",
	}
	dbmap.Create(&a1)
	a2 := Domaine{
		Name: "domaine 2",
	}
	dbmap.Create(&a2)

	as1 := Domaine{
		Name: "sub domaine 1",
		Parent: a1.ID,
	}
	dbmap.Create(&as1)

	as11 := Domaine{
		Name: "sub sub domaine 1",
		Parent: as1.ID,
	}
	dbmap.Create(&as11)
	as12 := Domaine{
		Name: "sub sub domaine 2",
		Parent: as1.ID,
	}
	dbmap.Create(&as12)
	return
}

func TestAccessDomaine(t *testing.T) {
	defer deleteFile(config.DBname)

	initDomaineValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/domaines"
	router.POST(urla, PostDomaine)
	router.GET(urla, GetDomaines)
	router.GET(urla+"/:id", GetDomaine)
	router.DELETE(urla+"/:id", DeleteDomaine)
	router.PUT(urla+"/:id", UpdateDomaine)

	b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET all Domaine")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as []Domaine
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 5, len(as), "5 results")

	// Get one
	log.Println("= http GET one Domaine")
	var a1 Domaine
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	assert.Equal(t, as[0].Name, a1.Name, "a1 = a")

	// Add
	log.Println("= http POST Domaine")
	var a = Domaine{Name: "Name test"}
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http POST denied")

	// Delete one
	log.Println("= http DELETE one Domaine")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Update one
	log.Println("= http PUT one Domaine")
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
	r.GET(urla, GetDomaines)
	// Get all
	log.Println("= http GET all Domaine with an other role")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []Domaine
	json.Unmarshal(resp.Body.Bytes(), &as2)
	assert.Equal(t, 5, len(as2), "5 results")
	//fmt.Println(as2)
}

func TestDomaine(t *testing.T) {
	defer deleteFile(config.DBname)

	InitDb(config.DBname, config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/domaines"
	router.POST(urla, PostDomaine)
	router.GET(urla, GetDomaines)
	router.GET(urla+"/:id", GetDomaine)
	router.DELETE(urla+"/:id", DeleteDomaine)
	router.PUT(urla+"/:id", UpdateDomaine)

	b := new(bytes.Buffer)
	// Add
	log.Println("= http POST Domaine")
	var a = Domaine{Name: "Name test"}
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

	// Add second domaine
	log.Println("= http POST more Domaine")
	var a2 = Domaine{Name: "Name test2"}
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = Domaine{}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

	// Get all
	log.Println("= http GET all Domaine")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Domaine
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
	log.Println("= http GET one Domaine")
	var a1 Domaine
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Domaine")
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
	log.Println("= http PUT one Domaine")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var a3 Domaine
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
