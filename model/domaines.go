package model

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE domaines (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  description text
, user3 integer, user2 integer, user1 integer, parent integer)

JSON Sample
-------------------------------------
{    "id": 49,    "name": "tOcCrAsePSBsCiaadNKmGwaFc",    "description": "EHqLQXIovUmoxxiSOssFIQrlT",    "user_3": 35,    "user_2": 17,    "user_1": 43,    "parent": 22}



*/

// Domaines struct is a row record of the domaines table in the pssimgmt database
type Domaine struct {
	ID          int32      `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Name        string     `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Description string     `gorm:"column:description;type:text;" json:"description"`
	User3       int64      `gorm:"column:user3;type:integer;" json:"user_3"`
	User2       int64      `gorm:"column:user2;type:integer;" json:"user_2"`
	User1       int64      `gorm:"column:user1;type:integer;" json:"user_1"`
	Parent      int32      `gorm:"column:parent;type:integer;" json:"parent"`
	Children    []*Domaine `gorm:"-" json:"children"`
}

func (d *Domaine) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Model(d).Where("parent = ?", d.ID).Update("parent", 0)
	tx.Where("domaine = ?", d.ID).Delete(ReglesDomaineses{})
	return nil
}

func GetDomaines(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetDomaines")

	query := "SELECT * FROM domaines"

	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)

	count := 0
	if s != "" {
		dbmap.Table("domaines").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("domaines").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var domaines []Domaine
	err := dbmap.Raw(query).Scan(&domaines).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, domaines)
	} else {
		c.JSON(404, gin.H{"error": "no domaine(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/domaines
}

func GetDomainesTree(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetDomainesTree")

	query := "SELECT * FROM domaines"
	if a.Role != "admin" {
		query += fmt.Sprintf(" WHERE user1 = %d OR user2 = %d OR user3 = %d", a.LoginID, a.LoginID, a.LoginID)
	}

	var domaines []*Domaine
	err := dbmap.Raw(query).Scan(&domaines).Error

	if err == nil {

		dmap := make(map[int32]*Domaine)
		tree := &Domaine{ID: 0, Name: "Périmètres"}
		if a.Role == "admin" {
			for _, e := range domaines {
				dmap[e.ID] = e
				if e.Parent == 0 {
					tree.Children = append(tree.Children, e)
				}

			}
			for _, e := range dmap {
				if dmap[e.Parent] != nil {
					dmap[e.Parent].Children = append(dmap[e.Parent].Children, e)
				}
			}
		} else {
			for _, e := range domaines {
				e.Parent = 0
				tree.Children = append(tree.Children, e)
			}
		}

		bytes, _ := json.Marshal(tree)

		c.Header("X-Total-Count", "1")
		//c.JSON(200, bytes)
		c.Data(200, "application/json", bytes)
	} else {
		c.JSON(404, gin.H{"error": "no domaine(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/domaines

}

func GetDomaine(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetDomaine " + id)

	var domaine Domaine
	err := dbmap.First(&domaine, id).Error

	if err == nil {
		c.JSON(200, domaine)
	} else {
		c.JSON(404, gin.H{"error": "domaine not found"})
	}

	// curl -i http://localhost:8080/api/v1/domaines/1
}

func PostDomaine(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostDomaine")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var domaine Domaine
	c.Bind(&domaine)

	if domaine.Name != "" { // XXX Check mandatory fields
		err := dbmap.Create(&domaine).Error
		if err == nil {
			c.JSON(201, domaine)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/domaines
}

func UpdateDomaine(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateDomaine " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var domaine Domaine
	err := dbmap.First(&domaine, id).Error
	if err == nil {
		idint, _ := strconv.Atoi(id)

		var json Domaine
		c.BindJSON(&json)

		fmt.Printf(" => %+v\n", json)
		//TODO : find fields via reflections
		//XXX custom fields mapping
		newdomaine := Domaine{
			Name:        json.Name,
			Description: json.Description,
			User1:       json.User1,
			User2:       json.User2,
			User3:       json.User3,
			Parent:      json.Parent,
		}
		if json.Parent == int32(idint) {
			newdomaine.Parent = 0
		}
		if newdomaine.Name != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&domaine).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":        newdomaine.Name,
					"Description": newdomaine.Description,
					"User1":       newdomaine.User1,
					"User2":       newdomaine.User2,
					"User3":       newdomaine.User3,
					"Parent":      newdomaine.Parent,
				}).Error; err == nil {
				c.JSON(200, newdomaine)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "domaine not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/domaines/1
}

func DeleteDomaine(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteDomaine " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var domaine Domaine
	err := dbmap.First(&domaine, id).Error

	if err == nil {
		e := dbmap.Delete(&domaine)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "domaine not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/domaines/1
}
