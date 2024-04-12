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


CREATE TABLE versions (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  validationdate datetime   ,
  validationpar text   ,
  status varchar(255)  DEFAULT 'redaction' ,
  changelog text   ,
  created_by integer   ,
  created_on timestamp   ,
  updated_by integer   ,
  updated_on timestamp
)

JSON Sample
-------------------------------------
{    "id": 49,    "name": "EgeTQdUURFvGfukhxKVuFoqNH",    "validationdate": "2298-09-04T21:13:07.274677904+02:00",    "validationpar": "likstjdcJdUvlJIIeXVeaiDDm",    "status": "DEOOjBOEQoCBDDvghGlAJOyjB",    "changelog": "yRdkZWcIflWIsHrZvnOXXNpFc",    "created_by": 12,    "created_on": "2056-06-20T23:47:25.779788166+02:00",    "updated_by": 82,    "updated_on": "2026-10-10T16:04:04.586123973+02:00"}



*/

// Versions struct is a row record of the versions table in the pssimgmt database
type Versions struct {
	Model
	Name           string    `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Validationdate time.Time `gorm:"column:validationdate;type:timestamp;" json:"validationdate"`
	Validationpar  string    `gorm:"column:validationpar;type:text;" json:"validationpar"`
	Status         string    `gorm:"column:status;type:varchar;size:255;" json:"status"`
	Changelog      string    `gorm:"column:changelog;type:text;" json:"changelog"`
}

func GetVersions(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetVersions")

	query := "SELECT * FROM versions"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("versions").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("versions").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var versions []Versions
	err := dbmap.Raw(query).Scan(&versions).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, versions)
	} else {
		c.JSON(404, gin.H{"error": "no version(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/versions
}

func GetVersion(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetVersion " + id)

	var version Versions
	err := dbmap.First(&version, id).Error

	if err == nil {
		c.JSON(200, version)
	} else {
		c.JSON(404, gin.H{"error": "version not found"})
	}

	// curl -i http://localhost:8080/api/v1/versions/1
}

func PostVersion(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostVersion")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var version Versions
	c.Bind(&version)

	if version.Name != "" { // XXX Check mandatory fields
		version.CreatedOn = time.Now()
		version.CreatedBy = a.LoginID
		err := dbmap.Create(&version).Error
		if err == nil {
			c.JSON(201, version)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/versions
}

func UpdateVersion(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateVersion " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var version Versions
	err := dbmap.First(&version, id).Error
	if err == nil {
		var json Versions
		c.Bind(&json)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newversion := Versions{
			Name:           json.Name,
			Validationpar:  json.Validationpar,
			Validationdate: json.Validationdate,
			Changelog:      json.Changelog,
			Status:         json.Status,
		}
		if newversion.Name != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&version).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":           newversion.Name,
					"Validationpar":  newversion.Validationpar,
					"Validationdate": newversion.Validationdate,
					"Changelog":      newversion.Changelog,
					"Status":         newversion.Status,
					"UpdatedOn":      time.Now(),
					"UpdatedBy":      a.LoginID,
				}).Error; err == nil {
				c.JSON(200, newversion)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "version not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/versions/1
}

func DeleteVersion(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteVersion " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var version Versions
	err := dbmap.First(&version, id).Error

	if err == nil {
		e := dbmap.Delete(&version)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "version not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/versions/1
}
