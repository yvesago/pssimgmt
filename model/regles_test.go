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

func initReglesValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Regles{
		Name:     "regle 1",
		Status: "ok",
		Descorig: "descorig regle 1",
	}
	dbmap.Create(&a1)
	a2 := Regles{
		Name:     "regle 2",
		Status: "ok",
		Descorig: "descorig regle 2",
	}
	dbmap.Create(&a2)

	d1 := Docs{
		Name: "doc 1",
	}
	dbmap.Create(&d1)
	d2 := Docs{
		Name: "doc 2",
	}
	dbmap.Create(&d2)

	i1 := Iso27002s{
		Name: "iso 1",
	}
	dbmap.Create(&i1)
	i2 := Iso27002s{
		Name: "iso 2",
	}
	dbmap.Create(&i2)

	ir1 := IsoRegleses{Regles: a2.ID, Iso: i1.ID}
	dbmap.Create(&ir1)
	ir2 := IsoRegleses{Regles: a2.ID, Iso: i2.ID}
	dbmap.Create(&ir2)

	dom1 := Domaine{
		Name: "domaine 1",
	}
	dbmap.Create(&dom1)
	dom2 := Domaine{
		Name: "domaine 2",
	}
	dbmap.Create(&dom2)

	dr1 := ReglesDomaineses{Regle: a1.ID, Domaine: dom1.ID, Applicable: 1, Modifdesc: "modif domaine 1 for regle 1 "}
	dbmap.Create(&dr1)

        //Themes
        t1 := Theme{
                Name: "theme 1",
        }
	dbmap.Create(&t1)
        t2 := Theme{
                Name: "theme 2",
        }
	dbmap.Create(&t2)

        rt1 := ReglesThemeses{Th: t2.ID, Regle: a1.ID}
        dbmap.Create(&rt1)

	return
}

func TestRegleByDom(t *testing.T) {
	defer deleteFile(config.DBname)
	initReglesValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/regles"
	router.GET(urla+"/:id/:dom", GetRegleByDom)
	router.PUT(urla+"/:id/:dom", UpdateRegleByDom)

	router.DELETE(urla+"/:id", DeleteRegle)
	var urlb = "/api/v1/iso"
	router.DELETE(urlb+"/:id", DeleteIso27002)

	// Get one
	log.Println("= http GET one Regle")
	var a1 Regles
	req, _ := http.NewRequest("GET", urla+"/1/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(prettyPrint(a1))

	// Get one
	log.Println("= http GET one Regle")
	var a2 Regles
	req, _ = http.NewRequest("GET", urla+"/2/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a2)
	//fmt.Println(prettyPrint(a2))

	// Update one
	b := new(bytes.Buffer)
	log.Println("= Update Regle with new RegleDomaine")
	dr2 := ReglesDomaineses{Modifdesc: "modif domaine 2 fo regle 2 ", Applicable: 1, Modif: "modif" }
	//a2.RegleDomaine = dr2
	//fmt.Println(prettyPrint(a2))
	json.NewEncoder(b).Encode(dr2)
	req, _ = http.NewRequest("PUT", urla+"/2/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http PUT success")
	var aRes Regles
	//json.Unmarshal(resp.Body.Bytes(), &aRes)
	//fmt.Println(prettyPrint(aRes))
	req, _ = http.NewRequest("GET", urla+"/2/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &aRes)
	//fmt.Println(prettyPrint(aRes))
	assert.Equal(t, dr2.Modifdesc, aRes.RegleDomaine.Modifdesc, "http PUT success")

	log.Println("= Update Regle with RegleDomaine")
	dr1 := ReglesDomaineses{Modifdesc: "modif domaine 1 for regle 1 ", Applicable: 1, Modif: "modif"}
	//a1.RegleDomaine = dr1
	//fmt.Println(prettyPrint(a1))
	json.NewEncoder(b).Encode(dr1)
	req, _ = http.NewRequest("PUT", urla+"/1/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http PUT success")
	var aRes1 Regles
	req, _ = http.NewRequest("GET", urla+"/1/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &aRes1)
	//fmt.Println(prettyPrint(aRes))
	assert.Equal(t, dr1.Modifdesc, aRes1.RegleDomaine.Modifdesc, "http PUT success")

	log.Println("= Update Eval Regle with RegleDomaine")
	dr12 := ReglesDomaineses{Applicable: 1, Modif: "eval"}
	json.NewEncoder(b).Encode(dr12)
	req, _ = http.NewRequest("PUT", urla+"/1/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http PUT success")
	req, _ = http.NewRequest("GET", urla+"/1/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &aRes1)
	assert.Equal(t, "0", aRes1.RegleDomaine.Conform, "Conform")
	assert.Equal(t, "0", aRes1.RegleDomaine.Evolution, "Evaluation")
	assert.Equal(t, dr1.Modifdesc, aRes1.RegleDomaine.Modifdesc, "unchanged Modifdesc")

        dr1 = ReglesDomaineses{Applicable: 1, Conform: "3", Evolution: "1", Modif: "eval"}
        json.NewEncoder(b).Encode(dr1)
        req, _ = http.NewRequest("PUT", urla+"/1/1", b)
        req.Header.Set("Content-Type", "application/json")
        resp = httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        assert.Equal(t, 201, resp.Code, "http PUT success")
        req, _ = http.NewRequest("GET", urla+"/1/1", nil)
        resp = httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        assert.Equal(t, 200, resp.Code, "http success")
        json.Unmarshal(resp.Body.Bytes(), &aRes1)
        assert.Equal(t, "3", aRes1.RegleDomaine.Conform, "Conform")
        assert.Equal(t, "0", aRes1.RegleDomaine.Evolution, "Evaluation")

}

func TestRelationsRegle(t *testing.T) {
	defer deleteFile(config.DBname)

	initReglesValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	b := new(bytes.Buffer)

	var urla = "/api/v1/regles"
	router.GET(urla+"/:id", GetRegle)
	router.DELETE(urla+"/:id", DeleteRegle)
	router.PUT(urla+"/:id", UpdateRegle)
	var urlb = "/api/v1/iso"
	router.DELETE(urlb+"/:id", DeleteIso27002)

	// Get one
	log.Println("= http GET one Regle")
	var a1 Regles
	req, _ := http.NewRequest("GET", urla+"/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	fmt.Println(prettyPrint(a1))
	assert.Equal(t,"theme 2",a1.Theme.Name,"regle 1 on Theme 2")

	var a2 Regles
	req, _ = http.NewRequest("GET", urla+"/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a2)
	//fmt.Println(prettyPrint(a2))
	assert.Equal(t, 2, len(a2.ReglesIso), "2 Iso for regles 2")

	// Update one
	log.Println("= Update Regle with Docs")
	a1.DocsIDs = []int32{1, 2}
	json.NewEncoder(b).Encode(a1)
	req, _ = http.NewRequest("PUT", urla+"/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var aRes Regles
	json.Unmarshal(resp.Body.Bytes(), &aRes)
	//fmt.Println(prettyPrint(aRes))
	assert.Equal(t, true, EqualArrayIds(a1.DocsIDs, aRes.DocsIDs), "DocsIDs Updated")

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
	assert.Equal(t, 1, len(a2.ReglesIso), "1 Iso for regles 2")
}

func TestAccessRegle(t *testing.T) {
	defer deleteFile(config.DBname)

	initReglesValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "reader"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/regles"
	router.POST(urla, PostRegle)
	router.GET(urla, GetRegles)
	router.GET(urla+"/:id", GetRegle)
	router.DELETE(urla+"/:id", DeleteRegle)
	router.PUT(urla+"/:id", UpdateRegle)

	b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET all Regles")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as []Regles
	json.Unmarshal(resp.Body.Bytes(), &as)
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one Regle")
	var a1 Regles
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	assert.Equal(t, as[0].Name, a1.Name, "a1 = a")

	// Add
	log.Println("= http POST Regle")
	var a = Regles{Name: "Name test"}
	json.NewEncoder(b).Encode(a)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http POST denied")

	// Delete one
	log.Println("= http DELETE one Regle")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Update one
	log.Println("= http PUT one Regle")
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
	r.GET(urla, GetRegles)
	// Get all
	log.Println("= http GET all Regles with an other role")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []Regles
	json.Unmarshal(resp.Body.Bytes(), &as2)
	assert.Equal(t, 2, len(as2), "2 results")
	//fmt.Println(as2)
}

func TestRegle(t *testing.T) {
	defer deleteFile(config.DBname)

	InitDb(config.DBname, config.Verbose)
	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))

	var urla = "/api/v1/regles"
	router.POST(urla, PostRegle)
	router.GET(urla, GetRegles)
	router.GET(urla+"/:id", GetRegle)
	router.DELETE(urla+"/:id", DeleteRegle)
	router.PUT(urla+"/:id", UpdateRegle)

	b := new(bytes.Buffer)
	// Add
	log.Println("= http POST Regle")
	var a = Regles{Name: "Name test", Status: "ok" }
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

	// Add second regle
	log.Println("= http POST more Regle")
	var a2 = Regles{Name: "Name test2"}
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")


	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = Regles{}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

        // Test bad status field
        log.Println("= Test bad status field")
        a2x = Regles{Name: "bad status field", Status: "xx"}
        json.NewEncoder(b).Encode(a2x)
        req, _ = http.NewRequest("POST", urla, b)
        req.Header.Set("Content-Type", "application/json")
        resp = httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        assert.Equal(t, 400, resp.Code, "http POST failed, bad status field")

        a2x = Regles{Name: "bad axe field", Axe1: "xx", Axe2: "0"}
        json.NewEncoder(b).Encode(a2x)
        req, _ = http.NewRequest("POST", urla, b)
        req.Header.Set("Content-Type", "application/json")
        resp = httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        assert.Equal(t, 400, resp.Code, "http POST failed, bad axe field")



	// Get all
	log.Println("= http GET all Regles")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Regles
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
	log.Println("= http GET one Regle")
	var a1 Regles
	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Regle")
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
	log.Println("= http PUT one Regle")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	var a3 Regles
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

	a2x = Regles{}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "Can't update missing mandatory field in /2")

	a2x = Regles{Name:"update name bas status", Status: "xx"}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "Can't update bad status field in /2")

	a2x = Regles{Name:"update name bas status", Status: "ok", Axe1: "xx"}
	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "Can't update bad axe field in /2")

}
