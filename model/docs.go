package model

import (
	//	"encoding/json"
	"time"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE docs (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  ordre integer   ,
  url text   ,
  description text   ,
  created_by integer   ,
  notes text   ,
  created_on timestamp   ,
  status varchar(255)   ,
  updated_by integer   ,
  updated_on timestamp
)

JSON Sample
-------------------------------------
{    "id": 47,    "name": "MNjtsBJieLjArysBgASIbmgAh",    "ordre": 17,    "url": "qBkgVOxrGtlOiPFSaSjttxmLs",    "description": "HWhUEeKHgsEsUpQLoUZsjeTLR",    "created_by": 85,    "notes": "TPxxpyfQEYDluZKnbgWCkbsNV",    "created_on": "2309-03-20T07:06:15.724800286+01:00",    "status": "WYgSAwAEIHdGmTEHvkyRUfUfZ",    "updated_by": 4,    "updated_on": "2310-01-11T11:53:38.82728523+01:00"}



*/

// Docs struct is a row record of the docs table in the pssimgmt database
type Docs struct {
	Model
	Name        string `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Ordre       int64  `gorm:"column:ordre;type:integer;" json:"ordre"`
	URL         string `gorm:"column:url;type:text;" json:"url"`
	Description string `gorm:"column:description;type:text;" json:"description"`
	Notes       string `gorm:"column:notes;type:text;" json:"notes"`
	Status      string `gorm:"column:status;type:varchar;size:255;" json:"status"`
}

func (d *Docs) BeforeDelete(tx *gorm.DB) (err error) {
  tx.Where("doc = ?", d.ID).Delete(DocRegleses{})
  return nil
}


func GetDocs(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetDocs")

	query := "SELECT * FROM docs"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("docs").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("docs").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var docs []Docs
	err := dbmap.Raw(query).Scan(&docs).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, docs)
	} else {
		c.JSON(404, gin.H{"error": "no doc(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/docs
}

func GetDoc(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetDoc " + id)

	var doc Docs
	err := dbmap.First(&doc, id).Error

	if err == nil {
		c.JSON(200, doc)
	} else {
		c.JSON(404, gin.H{"error": "doc not found"})
	}

	// curl -i http://localhost:8080/api/v1/docs/1
}

func PostDoc(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostDoc")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var doc Docs
	c.Bind(&doc)

	if doc.Name != "" { // XXX Check mandatory fields
		doc.CreatedOn = time.Now()
		doc.CreatedBy = a.LoginID
		err := dbmap.Create(&doc).Error
		if err == nil {
			c.JSON(201, doc)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/docs
}

func UpdateDoc(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateDoc " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var doc Docs
	err := dbmap.First(&doc, id).Error
	if err == nil {
		var json Docs
		c.Bind(&json)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newdoc := Docs{
			Name:        json.Name,
			Ordre:       json.Ordre,
			URL:         json.URL,
			Description: json.Description,
			Notes:       json.Notes,
			Status:      json.Status,
		}
		if newdoc.Name != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&doc).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":        newdoc.Name,
					"Ordre":       newdoc.Ordre,
					"URL":         newdoc.URL,
					"Description": newdoc.Description,
					"Notes":       newdoc.Notes,
					"Status":      newdoc.Status,
					"UpdatedOn":   time.Now(),
					"UpdatedBy":   a.LoginID,
				}).Error; err == nil {
				c.JSON(200, newdoc)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "doc not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/docs/1
}

func DeleteDoc(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteDoc " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var doc Docs
	err := dbmap.First(&doc, id).Error

	if err == nil {
		e := dbmap.Delete(&doc)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "doc not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/docs/1
}
