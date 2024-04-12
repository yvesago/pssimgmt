package model

import (
        "github.com/gin-gonic/gin"
        "github.com/jinzhu/gorm"
        "strconv"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE iso27002s (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  code varchar(255) NOT NULL  ,
  descorig text
)

JSON Sample
-------------------------------------
{    "id": 65,    "name": "nwaGWHqVaNGJXCFswFVMLpcJe",    "code": "cFEYeTNMCJNMurpeMDsVhmGOm",    "descorig": "OJqlpiHktjxTpoaopynBNxVEt"}



*/

// Iso27002s struct is a row record of the iso27002s table in the pssimgmt database
type Iso27002s struct {
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Name string `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Code string `gorm:"column:code;type:varchar;size:255;" json:"code"`
	Descorig string `gorm:"column:descorig;type:text;" json:"descorig"`
}


func (i *Iso27002s) BeforeDelete(tx *gorm.DB) (err error) {
  tx.Where("iso = ?", i.ID).Delete(IsoThemes{})
  tx.Where("iso = ?", i.ID).Delete(IsoRegleses{})
  return nil
}



func GetIso27002s(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

        a := Auth(c)
        a.Log("GetIso27002s")

	query := "SELECT * FROM iso27002s"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("iso27002s").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("iso27002s").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var iso27002s []Iso27002s
	err := dbmap.Raw(query).Scan(&iso27002s).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, iso27002s)
	} else {
		c.JSON(404, gin.H{"error": "no iso27002(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/iso27002s
}

func GetIso27002(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

        a := Auth(c)
        a.Log("GetIso27002 " + id)

	var iso27002 Iso27002s
	err := dbmap.First(&iso27002, id).Error

	if err == nil {
		c.JSON(200, iso27002)
	} else {
		c.JSON(404, gin.H{"error": "iso27002 not found"})
	}

	// curl -i http://localhost:8080/api/v1/iso27002s/1
}

func PostIso27002(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

        a := Auth(c)
        a.Log("PostIso27002")
        if a.Role != "admin" {
                c.JSON(403, gin.H{"error": "admin role required"})
                return
        }


	var iso27002 Iso27002s
	c.Bind(&iso27002)

	if iso27002.Name != "" { // XXX Check mandatory fields
		err := dbmap.Create(&iso27002).Error
		if err == nil {
			c.JSON(201, iso27002)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/iso27002s
}

func UpdateIso27002(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

        a := Auth(c)
        a.Log("UpdateIso27002 " + id)
        if a.Role != "admin" {
                c.JSON(403, gin.H{"error": "admin role required"})
                return
        }


	var iso27002 Iso27002s
	err := dbmap.First(&iso27002, id).Error
	if err == nil {
		var json Iso27002s
		c.Bind(&json)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newiso27002 := Iso27002s{
			Name: json.Name,
		}
		if newiso27002.Name != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&iso27002).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name": newiso27002.Name,
				}).Error; err == nil {
				c.JSON(200, newiso27002)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "iso27002 not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/iso27002s/1
}

func DeleteIso27002(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

        a := Auth(c)
        a.Log("DeleteIso27002 " + id)
        if a.Role != "admin" {
                c.JSON(403, gin.H{"error": "admin role required"})
                return
        }

	var iso27002 Iso27002s
	err := dbmap.First(&iso27002, id).Error

	if err == nil {
		e := dbmap.Delete(&iso27002)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "iso27002 not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/iso27002s/1
}
