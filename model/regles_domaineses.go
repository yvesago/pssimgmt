package model

import (
	"time"
        "github.com/gin-gonic/gin"
        "github.com/jinzhu/gorm"
        "strconv"
)

/*
DB Table Details
-------------------------------------

CREATE TABLE regles_domaineses (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  regle integer NOT NULL  ,
  domaine integer NOT NULL  ,
  modifdesc text   ,
  supldesc text
, evolution varchar(255)  DEFAULT 0, conform varchar(255), applicable integer  DEFAULT 1, created_by integer, created_on timestamp, updated_by integer, updated_on timestamp)

-------------------------------------

*/

// ReglesDomaineseses struct is a row record of the regles_domaineses table in the pssimgmt database
type ReglesDomaineses struct {
	ID         int32     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Regle      int32     `gorm:"column:regle;type:integer;" json:"regle"`
	Domaine    int32     `gorm:"column:domaine;type:integer;" json:"domaine_id"`
	Modifdesc  string    `gorm:"column:modifdesc;type:text;" json:"modifdesc"`
	Supldesc   string    `gorm:"column:supldesc;type:text;" json:"supldesc"`
	Evolution  string    `gorm:"column:evolution;type:varchar;size:255;default:0;" json:"evolution"`
	Conform    string    `gorm:"column:conform;type:varchar;size:255;" json:"conform"`
	Applicable int64     `gorm:"column:applicable;type:integer;" json:"applicable"`
	CreatedBy  int64     `gorm:"column:created_by;type:integer;" json:"created_by"`
	CreatedOn  time.Time `gorm:"column:created_on;type:timestamp;" json:"created_on"`
	UpdatedBy  int64     `gorm:"column:updated_by;type:integer;" json:"updated_by"`
	UpdatedOn  time.Time `gorm:"column:updated_on;type:timestamp;" json:"updated_on"`
	Domaines   Domaine   `gorm:"-;" json:"domain"`
	Modif      string    `gorm:"-;" json:"modif"`
        User1      int32     `gorm:"-" json:"user_1"`
}


func GetTodos(c *gin.Context) {
        dbmap := c.MustGet("DBmap").(*gorm.DB)

        a := Auth(c)
        a.Log("GetTodos")
        /*if a.Role != "admin" {
                c.JSON(403, gin.H{"error": "admin role required"})
                return
        }*/

        query := "SELECT * FROM regles_domaineses LEFT JOIN domaines ON regles_domaineses.domaine = domaines.ID"
        // Parse query string
        //  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
        q := c.Request.URL.Query()
        s, o, l := ParseQuery(q)
        count := 0
        if s != "" {
                //dbmap.Table("regles_domaineses").Where(s).Count(&count)
                dbmap.Table("regles_domaineses").Joins("join domaines ON regles_domaineses.domaine = domaines.ID").Where(s).Count(&count)
                query = query + " WHERE " + s
        } else {
                dbmap.Table("regles_domaineses").Count(&count)
        }
        if o != "" {
                query = query + o
        }
        if l != "" {
                query = query + l
        }

        var regles_domaineses []ReglesDomaineses
        err := dbmap.Raw(query).Scan(&regles_domaineses).Error

        if err == nil {
                c.Header("X-Total-Count", strconv.Itoa(count))
                c.JSON(200, regles_domaineses)
        } else {
                c.JSON(404, gin.H{"error": "no regles_domainese(s) into the table"})
        }

        // curl -i http://localhost:8080/api/v1/regles_domaineses
}

func DeleteTodo(c *gin.Context) {
        dbmap := c.MustGet("DBmap").(*gorm.DB)
        id := c.Params.ByName("id")

        a := Auth(c)
        a.Log("DeleteReglesDomainese " + id)
        if a.Role != "admin" {
                c.JSON(403, gin.H{"error": "admin role required"})
                return
        }

        var regles_domainese ReglesDomaineses
        err := dbmap.First(&regles_domainese, id).Error

        if err == nil {
                e := dbmap.Delete(&regles_domainese)

                if e != nil {
                        c.JSON(200, gin.H{"id #" + id: "deleted"})
                } else {
                        checkErr(nil, "Delete failed")
                }
        } else {
                c.JSON(404, gin.H{"error": "regles_domainese not found"})
        }

        // curl -i -X DELETE http://localhost:8080/api/v1/regles_domaineses/1
}
