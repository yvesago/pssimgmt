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

func initThemeValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	//Themes
	a1 := Theme{
		Name:   "theme 1",
		Status: "ok",
	}
	dbmap.Create(&a1)
	a2 := Theme{
		Name:   "theme 2",
		Status: "ok",
	}
	dbmap.Create(&a2)
	a12 := Theme{
		Name:   "sub theme 12",
		Status: "ok",
		Parent: a1.ID,
		Ordre:  20,
	}
	dbmap.Create(&a12)
	a11 := Theme{
		Name:   "subtheme 11",
		Status: "ok",
		Parent: a1.ID,
		Ordre:  10,
	}
	dbmap.Create(&a11)

	// Regles
	r1 := Regles{
		Name:   "regle 1",
		Status: "ok",
		Axe1:   "1",
		Axe2:   "6",
	}
	dbmap.Create(&r1)
	r2 := Regles{
		Name:   "regle 2",
		Status: "ok",
		Axe1:   "6",
	}
	dbmap.Create(&r2)
	r3 := Regles{
		Status: "ok",
		Name:   "regle 3",
	}
	dbmap.Create(&r3)
	r4 := Regles{
		Status: "ok",
		Name:   "regle 4",
	}
	dbmap.Create(&r4)

	rt1 := ReglesThemeses{Th: a11.ID, Regle: r1.ID}
	dbmap.Create(&rt1)
	rt2 := ReglesThemeses{Th: a11.ID, Regle: r2.ID}
	dbmap.Create(&rt2)
	rt3 := ReglesThemeses{Th: a12.ID, Regle: r3.ID}
	dbmap.Create(&rt3)
	rt4 := ReglesThemeses{Th: a12.ID, Regle: r4.ID}
	dbmap.Create(&rt4)

	// Iso
	i1 := Iso27002s{
		Name: "iso 1",
	}
	dbmap.Create(&i1)
	i2 := Iso27002s{
		Name: "iso 2",
	}
	dbmap.Create(&i2)

	it1 := IsoThemes{Th: 2, Iso: 1}
	dbmap.Create(&it1)
	it2 := IsoThemes{Th: 2, Iso: 2}
	dbmap.Create(&it2)

	// Domaines
	dom1 := Domaine{
		Name: "domaine 1",
	}
	dbmap.Create(&dom1)
	dom2 := Domaine{
		Name: "domaine 2",
	}
	dbmap.Create(&dom2)

	dr1 := ReglesDomaineses{Regle: r1.ID, Domaine: dom1.ID, Conform: "3", Evolution: "0"}
	dbmap.Create(&dr1)
	dr2 := ReglesDomaineses{Regle: r2.ID, Domaine: dom1.ID, Conform: "2", Evolution: "1"}
	dbmap.Create(&dr2)

	return
}

func TestThemeByDom(t *testing.T) {
	defer deleteFile(config.DBname)

	initThemeValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/themes"
	//router.GET(urla+"tree", GetThemesTree)
	//router.GET(urla+"/:id", GetTheme)
	router.GET(urla+"/:id/:dom", GetThemeByDom)

	// Get theme by dom
	log.Println("= http GET Theme 4 for Dom 1")
	req, _ := http.NewRequest("GET", urla+"/4/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var at Theme
	json.Unmarshal(resp.Body.Bytes(), &at)
	fmt.Println(at.Evaluation)
	//fmt.Println(prettyPrint(at))
	//fmt.Printf("%+v\n", at.Evaluation)
	tres := `	[4] subtheme 11 1
		(1) regle 1
		(2) regle 2
`
	assert.Equal(t, 2, len(at.Regles), "2 rules")
	assert.Equal(t, tres, at.String(), "2 rules")
	assert.Equal(t, "Gouvernance", at.Evaluation.Axes[0], "Gouvernance first axe")
	assert.Equal(t, "Évaluation", at.Evaluation.Axes[5], "Évaluation second axe")
	assert.Equal(t, 10, at.Evaluation.Conforme[0], "Gouvernance 5/10 conforme")
	assert.Equal(t, 10, at.Evaluation.Evolution[0], "Gouvernance 5/10 evolution")
	assert.Equal(t, 8, at.Evaluation.Conforme[5], "Évaluation 8/10 conforme")
	assert.Equal(t, 10, at.Evaluation.Evolution[5], "Évaluation 10/10 evolution")

	log.Println("= http GET Theme 4 for Dom 2")
	req, _ = http.NewRequest("GET", urla+"/4/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &at)
	fmt.Println(at.Evaluation)
}

func TestRelationsTheme(t *testing.T) {
	defer deleteFile(config.DBname)

	initThemeValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	b := new(bytes.Buffer)

	var urla = "/api/v1/themes"
	router.GET(urla+"/:id", GetTheme)
	router.DELETE(urla+"/:id", DeleteTheme)
	router.PUT(urla+"/:id", UpdateTheme)
	var urlb = "/api/v1/iso"
	router.DELETE(urlb+"/:id", DeleteIso27002)

	// Get one
	log.Println("= http GET one Theme")
	var a1 Theme
	req, _ := http.NewRequest("GET", urla+"/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	var a2 Theme
	req, _ = http.NewRequest("GET", urla+"/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a2)
	//fmt.Println(prettyPrint(a2))
	assert.Equal(t, 2, len(a2.IsoThemes), "2 Iso for theme 2")

	// Update one
	log.Println("= Update Theme with Regles")
	a1.ReglesIDs = []int32{1, 2}
	json.NewEncoder(b).Encode(a1)
	req, _ = http.NewRequest("PUT", urla+"/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var aRes Theme
	json.Unmarshal(resp.Body.Bytes(), &aRes)
	//fmt.Println(prettyPrint(aRes))
	assert.Equal(t, true, EqualArrayIds(a1.ReglesIDs, aRes.ReglesIDs), "ReglesIDs Updated")

	// Hook delete Iso
	// Delete one Iso
	log.Println("= http DELETE one Iso")
	req, _ = http.NewRequest("DELETE", urlb+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE ok")

	req, _ = http.NewRequest("GET", urla+"/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a2)
	//fmt.Println(prettyPrint(a2))
	assert.Equal(t, 1, len(a2.IsoThemes), "1 Iso for theme 2")
}

func TestThemeTree(t *testing.T) {
	defer deleteFile(config.DBname)

	initThemeValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/themes"
	router.GET(urla+"tree", GetThemesTree)
	router.GET(urla+"/:id", GetTheme)
	router.GET(urla+"/:id/:dom", GetThemeByDom)

	// Get Theme Tree
	log.Println("= http GET Theme Tree")
	req, _ := http.NewRequest("GET", urla+"tree", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as Theme
	json.Unmarshal(resp.Body.Bytes(), &as)
	fmt.Println(&as)
	//fmt.Printf("%+v\n", as)
	//fmt.Println(prettyPrint(&as))
	tres := `[0] PSSI 0
[1] theme 1 0
	[4] subtheme 11 1
	[3] sub theme 12 1
[2] theme 2 0
`
	assert.Equal(t, "PSSI", as.Name, "1 result")
	assert.Equal(t, tres, as.String(), "1 result")

	// Get theme by dom
	log.Println("= http GET Theme By Dom")
	req, _ = http.NewRequest("GET", urla+"/4/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var at Theme
	json.Unmarshal(resp.Body.Bytes(), &at)
	fmt.Println(&at)
	//fmt.Printf("%+v\n", at)
	tres = `	[4] subtheme 11 1
		(1) regle 1
		(2) regle 2
`
	assert.Equal(t, 2, len(at.Regles), "2 rules")
	assert.Equal(t, tres, at.String(), "2 rules")
	assert.Equal(t, "subtheme 11", at.Regles[0].Theme.Name, "theme name for regle")
	//fmt.Println(prettyPrint(at.Regles[0]))

}

func TestAccessTheme(t *testing.T) {
	defer deleteFile(config.DBname)

	initThemeValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/themes"
	router.POST(urla, PostTheme)
	router.GET(urla, GetThemes)
	router.GET(urla+"/:id", GetTheme)
	router.DELETE(urla+"/:id", DeleteTheme)
	router.PUT(urla+"/:id", UpdateTheme)

	b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET all Theme")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as []Theme
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 4, len(as), "4 results")
	//fmt.Println(prettyPrint(&as))

	// Get one
	log.Println("= http GET one Theme")
	var a1 Theme
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	assert.Equal(t, as[0].Name, a1.Name, "a1 = a")
	//fmt.Println(prettyPrint(&a1))

	// Add
	log.Println("= http POST Theme")
	var a = Theme{Name: "Name test"}
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http POST denied")

	// Delete one
	log.Println("= http DELETE one Theme")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Update one
	log.Println("= http PUT one Theme")
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
	r.GET(urla, GetThemes)
	// Get all
	log.Println("= http GET all Theme with an other role")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []Theme
	json.Unmarshal(resp.Body.Bytes(), &as2)
	assert.Equal(t, 4, len(as2), "4 results")
	//fmt.Println(as2)
}

func TestTheme(t *testing.T) {
	defer deleteFile(config.DBname)

	InitDb(config.DBname, config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{LoginID: 1, Login: "adminlogin", Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/themes"
	router.POST(urla, PostTheme)
	router.GET(urla, GetThemes)
	router.GET(urla+"/:id", GetTheme)
	router.DELETE(urla+"/:id", DeleteTheme)
	router.PUT(urla+"/:id", UpdateTheme)

	b := new(bytes.Buffer)
	// Add
	log.Println("= http POST Theme")
	var a = Theme{Name: "Name test", Status: "ok"}
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

	// Add second theme
	log.Println("= http POST more Theme")
	var a2 = Theme{Name: "Name test2", Status: "ok"}
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = Theme{Description: "missing name"}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

	// Test bad status field
	log.Println("= Test bad status field")
	a2x = Theme{Name: "bad status field", Status: "xx"}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, bad status field")

	// Get all
	log.Println("= http GET all Theme")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Theme
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
	assert.Equal(t, ` (name LIKE "%t%")`, query, "Parse query")
	assert.Equal(t, " ORDER BY created_on ASC", s, "Parse query")

	// Get one
	log.Println("= http GET one Theme")
	var a1 Theme
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Theme")
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
	log.Println("= http PUT one Theme")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var a3 Theme
	req, _ = http.NewRequest("GET", urla+"/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	assert.Equal(t, a2.Name, a3.Name, "a2 Name updated")
	//fmt.Printf("%+v\n", a3)

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
