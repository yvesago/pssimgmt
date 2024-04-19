package model

import (
	"encoding/json"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initTodoValues(connString string, verbose bool) {
	dbmap := InitDb(connString, verbose)
	a1 := Users{
		CasID: "user1",
	}
	dbmap.Create(&a1)
	a2 := Users{
		CasID: "user2",
	}
	dbmap.Create(&a2)

	dom1 := Domaine{
		Name:  "domaine 1",
		User1: int64(a1.ID),
		User2: int64(a1.ID),
		User3: int64(a2.ID),
	}
	dbmap.Create(&dom1)
	dom2 := Domaine{
		Name:  "domaine 2",
		User2: int64(a1.ID),
	}
	dbmap.Create(&dom2)

    r1 := Regles{
        Name:     "regle 1",
        Status: "ok",
        Descorig: "descorig regle 1",
    }
    dbmap.Create(&r1)
    r2 := Regles{
        Name:     "regle 2",
        Status: "ok",
        Descorig: "descorig regle 2",
    }
    dbmap.Create(&r2)

    dr1 := ReglesDomaineses{Regle: r1.ID, Domaine: dom1.ID, Applicable: 1, Modifdesc: "modif domaine 1 for regle 1 "}
    dbmap.Create(&dr1)

    dr2 := ReglesDomaineses{Regle: r2.ID, Domaine: dom2.ID, Applicable: 1, Modifdesc: "modif domaine 2 for regle 2 "}
    dbmap.Create(&dr2)

	return
}

func TestTodo(t *testing.T) {
	defer deleteFile(config.DBname)

	initTodoValues(config.DBname, false) //config.Verbose)

	config = DBs(config)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	userauth := AuthInfo{LoginID: 1, Role: "cssi"}
	router.Use(SetConfig(config, userauth))
	//router.Use(Database(config.DBname))

	var urla = "/api/v1/todos"
	router.GET(urla, GetTodos)
	router.DELETE(urla+"/:id", DeleteTodo)

	//b := new(bytes.Buffer)

	// Get all
	log.Println("= http GET Todos for user 1")
	req, _ := http.NewRequest("GET", urla, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as []ReglesDomaineses
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Printf("%+v\n", as)
	assert.Equal(t, 2, len(as), "2 results")

	// Change user
	userauth = AuthInfo{LoginID: 2, Role: "cssi"}
	router2 := gin.New()
	router2.Use(SetConfig(config, userauth))
	router2.GET(urla, GetTodos)
	log.Println("= http GET Todos for user 2")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	router2.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as2 []ReglesDomaineses
	json.Unmarshal(resp.Body.Bytes(), &as2)
	//fmt.Printf("%+v\n", as2)
	assert.Equal(t, 1, len(as2), "1 results")

	// Delete one
	log.Println("= http DELETE one Todo")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 403, resp.Code, "http DELETE denied")

	// Change user
	userauth = AuthInfo{LoginID: 2, Role: "admin"}
	r := gin.New()
	r.Use(SetConfig(config, userauth))
	r.GET(urla, GetTodos)
	r.DELETE(urla+"/:id", DeleteTodo)
	// Get all
	log.Println("= http GET all Todo for admin")
	req, _ = http.NewRequest("GET", urla, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	var as3 []ReglesDomaineses
	json.Unmarshal(resp.Body.Bytes(), &as3)
	assert.Equal(t, 2, len(as3), "2 results")

	// Delete one
	log.Println("= http DELETE one Todo")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	//fmt.Printf("%+v\n", resp)
	assert.Equal(t, 200, resp.Code, "http DELETE ")

}

