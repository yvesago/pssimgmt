package model

import (
	//	"encoding/json"
	"time"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)


// Documents struct is a row record of the documents table in the pssimgmt database
type Documents struct {
	Model
	Name        string `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Titre       string `gorm:"column:titre;type:varchar;size:255;" json:"titre"`
	Description string `gorm:"column:description;type:text;" json:"description"`
	Notes       string `gorm:"column:notes;type:text;" json:"notes"`
}

func GetDocuments(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetDocuments")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	query := "SELECT * FROM documents"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("documents").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("documents").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var documents []Documents
	err := dbmap.Raw(query).Scan(&documents).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, documents)
	} else {
		c.JSON(404, gin.H{"error": "no document(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/documents
}

func GetDocument(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetDocument " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var document Documents
	err := dbmap.First(&document, id).Error

	if err == nil {
		c.JSON(200, document)
	} else {
		c.JSON(404, gin.H{"error": "document not found"})
	}

	// curl -i http://localhost:8080/api/v1/documents/1
}

func PostDocument(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostDocument")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var document Documents
	c.Bind(&document)

	if document.Name != "" { // XXX Check mandatory fields
		document.CreatedOn = time.Now()
		document.CreatedBy = a.LoginID
		err := dbmap.Create(&document).Error
		if err == nil {
			c.JSON(201, document)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/documents
}

func UpdateDocument(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateDocument " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var document Documents
	err := dbmap.First(&document, id).Error
	if err == nil {
		var json Documents
		c.Bind(&json)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newdocument := Documents{
			Name:        json.Name,
			Titre:       json.Titre,
			Description: json.Description,
			Notes:       json.Notes,
		}
		if newdocument.Name != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&document).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":        newdocument.Name,
					"Titre":       newdocument.Titre,
					"Description": newdocument.Description,
					"Notes":       newdocument.Notes,
					"UpdatedOn":   time.Now(),
					"UpdatedBy":   a.LoginID,
				}).Error; err == nil {
				c.JSON(200, newdocument)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "document not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/documents/1
}

func DeleteDocument(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteDocument " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var document Documents
	err := dbmap.First(&document, id).Error

	if err == nil {
		e := dbmap.Delete(&document)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "document not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/documents/1
}
