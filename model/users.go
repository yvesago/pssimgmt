package model

import (
	//	"encoding/json"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  cas_id text   ,
  name text   ,
  user_role varchar(255)  DEFAULT 'user' ,
  email text  DEFAULT '' ,
  email_confirmed boolean
)

-------------------------------------



*/

// Users struct is a row record of the users table in the pssimgmt database
type Users struct {
	ID       int32   `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	CasID    string  `gorm:"column:cas_id;type:text;" json:"casid"`
	Name     string  `gorm:"column:name;type:text;" json:"name"`
	UserRole string  `gorm:"column:user_role;type:varchar;size:255;default:'user';" json:"user_role"`
	Email    string  `gorm:"column:email;type:text;default:'';" json:"email"`
	Doms     []int32 `gorm:"-" json:"doms"`
}

func ValidUserRoles(val string) bool {
	if Contains([]string{"admin", "cssi", "reader", "guest"}, val) {
		return true
	}
	return false
}

func (u *Users) ByLogin(dbmap *gorm.DB, login string) error {
	err := dbmap.First(u, "cas_id = ?", login).Error
	if u.ID == 0 {
		newuser := Users{CasID: login, Name: login, UserRole: "guest"}
		e := dbmap.Create(&newuser).Error
		u.ID = newuser.ID
		return e
	}
	return err
}

func (u *Users) AfterFind(tx *gorm.DB) (err error) {
	var domids []int32
	tx.Model(&Domaine{}).Where("user1=? OR user2=? OR user3=?", u.ID, u.ID, u.ID).Pluck("DISTINCT(ID)", &domids)
	//fmt.Println(domids)
	u.Doms = domids
	return nil
}

/*func (u *Users) AfterCreate(tx *gorm.DB) (err error) {
        var doms []Domaine
        var domids []int32
        tx.Table("domaines").Select("id").Where("user1=? OR user2=? OR user3=?", u.ID, u.ID, u.ID).Scan(&doms)
	fmt.Println(doms)
	for _, d := range doms {
	   domids = append(domids, d.ID)
	}
        u.Doms = domids
        return nil
}*/

func GetUsers(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetUsers")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	query := "SELECT * FROM users"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("users").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("users").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var users []Users
	err := dbmap.Raw(query).Scan(&users).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/users
}

func GetUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetUser " + id)
	if a.Role != "admin" && strconv.Itoa(int(a.LoginID)) != id {
		c.JSON(403, gin.H{"error": "admin role or self user required"})
		return
	}

	var user Users
	err := dbmap.First(&user, id).Error

	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func PostUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostUser")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var user Users
	c.Bind(&user)

	if ValidUserRoles(user.UserRole) == false {
		user.UserRole = "guest"
	}

	if user.CasID != "" { // XXX Check mandatory fields
		err := dbmap.Create(&user).Error
		if err == nil {
			c.JSON(201, user)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateUser " + id)
	//if a.Role != "admin" && a.Role != "validator" {
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var user Users
	err := dbmap.First(&user, id).Error
	if err == nil {
		var json Users
		c.Bind(&json)

		if ValidUserRoles(json.UserRole) == false {
			json.UserRole = "guest"
		}

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newuser := Users{
			Name:     json.Name,
			CasID:    json.CasID,
			UserRole: json.UserRole,
			Email:    json.Email,
		}
		if newuser.CasID != "" { // XXX Check mandatory fields
			if err := dbmap.Model(&user).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":     newuser.Name,
					"CasID":    newuser.CasID,
					"UserRole": newuser.UserRole,
					"Email":    newuser.Email,
				}).Error; err == nil {
				c.JSON(200, newuser)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteUser " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var user Users
	err := dbmap.First(&user, id).Error

	if err == nil {
		e := dbmap.Delete(&user)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
}
