package model

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	//"net/url"
	"testing"
)

/*
  domaine 1 : regle1, regle 2
  |_domaine 2 : r1 => c1
    |_ domaine 3 :
    |_ domaine 4 : r1 => c3 

*/

func initReglesDomValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Regles{
		Name:     "regle 1",
		Descorig: "descorig regle 1",
	}
	dbmap.Create(&a1)
	a2 := Regles{
		Name:     "regle 2",
		Descorig: "descorig regle 2",
	}
	dbmap.Create(&a2)

	dom1 := Domaine{
		Name: "domaine 1",
	}
	dbmap.Create(&dom1)
	dom2 := Domaine{
		Parent: dom1.ID,
		Name: "sub domaine 2",
	}
	dbmap.Create(&dom2)
	dom3 := Domaine{
		Parent: dom2.ID,
		Name: "sub sub domaine 3",
	}
	dbmap.Create(&dom3)
	dom4 := Domaine{
		Parent: dom2.ID,
		Name: "sub sub domaine 4",
	}
	dbmap.Create(&dom4)

	//dr1 := ReglesDomaineses{Regle: a1.ID, Domaine: dom1.ID, Modifdesc: "modif domaine 1 for regle 1 "}
	dr1 := ReglesDomaineses{Regle: a1.ID, Domaine: dom2.ID, Applicable: 1, Conform: "1", Evolution: "1"}
	dbmap.Create(&dr1)

	//dr2 := ReglesDomaineses{Regle: a1.ID, Domaine: dom3.ID, Modif: "eval"}
	//dbmap.Create(&dr2)

	dr4 := ReglesDomaineses{Regle: a1.ID, Domaine: dom4.ID, Applicable: 1, Conform: "3", Modifdesc: "modif domaine 4 for regle 1 "}
	dbmap.Create(&dr4)

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

func TestInheritEval(t *testing.T) {
	defer deleteFile(config.DBname)
	initReglesDomValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{Role: "admin"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/regles"
	router.GET(urla+"/:id/:dom", GetRegleByDom)

	// Get one
	log.Println("= http GET Regle 1 in Dom 1")
	var a1 Regles
	req, _ := http.NewRequest("GET", urla+"/1/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(prettyPrint(a1))
	fmt.Println(a1.RegleDomaine.Applicable, prettyPrint(a1.RegleDomaine.Conform))
	assert.Equal(t, "", a1.RegleDomaine.Conform, "Conform not specified")

	// Get one
	log.Println("= http GET Regle 1 in Dom 2")
	var a2 Regles
	req, _ = http.NewRequest("GET", urla+"/1/2", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a2)
	//fmt.Println(prettyPrint(a2))
	fmt.Println(a2.RegleDomaine.Applicable, prettyPrint(a2.RegleDomaine.Conform))
	assert.Equal(t, "1", a2.RegleDomaine.Conform, "Conform")

	// Get one
	log.Println("= http GET Regle 1 in Dom 3")
	var a3 Regles
	req, _ = http.NewRequest("GET", urla+"/1/3", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(prettyPrint(a3.RegleDomaine.Domaines))
	fmt.Println(a3.RegleDomaine.Applicable,prettyPrint(a3.RegleDomaine.Conform))
	assert.Equal(t, "1", a3.RegleDomaine.Conform, "Conform inherit from 2")

	// Get one
	log.Println("= http GET Regle 1 in Dom 4")
	var a4 Regles
	req, _ = http.NewRequest("GET", urla+"/1/4", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a4)
	//fmt.Println(prettyPrint(a4))
	fmt.Println(a4.RegleDomaine.Applicable, prettyPrint(a4.RegleDomaine.Conform))
	assert.Equal(t, "3", a4.RegleDomaine.Conform, "Conform specified")

	/*dr1 = ReglesDomaineses{Applicable: 1, Conform: "3", Evolution: "1", Modif: "eval"}
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
	assert.Equal(t, "0", aRes1.RegleDomaine.Evolution, "Evaluation")*/

}

