package model

import (
	//	"fmt"
	//"encoding/json"
	"errors"
	//"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE regles (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  ordre integer   ,
  code varchar(255) NOT NULL  ,
  descorig text   ,
  description text   ,
  notes text   ,
  status varchar(255)   ,
  created_by integer   ,
  created_on timestamp   ,
  updated_by integer   ,
  updated_on timestamp
, axe2 varchar(255), axe1 varchar(255))

*/

// Regles struct is a row record of the regles table in the pssimgmt database
type Regles struct {
	Model
	Name         string           `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Ordre        int64            `gorm:"column:ordre;type:integer;" json:"ordre"`
	Code         string           `gorm:"column:code;type:varchar;size:255;" json:"code"`
	Descorig     string           `gorm:"column:descorig;type:text;" json:"descorig"`
	Description  string           `gorm:"column:description;type:text;" json:"description"`
	Notes        string           `gorm:"column:notes;type:text;" json:"notes"`
	Status       string           `gorm:"column:status;type:varchar;size:255;" json:"status"`
	Axe2         string           `gorm:"column:axe2;type:varchar;size:255;" json:"axe2"`
	Axe1         string           `gorm:"column:axe1;type:varchar;size:255;" json:"axe1"`
	RegleDomaine ReglesDomaineses `gorm:"-" json:"regle_domaine"` // regles for domain : populate when reading theme by domain
	ReglesIso    []*Iso27002s     `gorm:"-" json:"regles_iso"`
	ReglesIsoIDs []int32          `gorm:"-" json:"iso_ids"`
	Docs         []*Docs          `gorm:"-" json:"docs"`
	DocsIDs      []int32          `gorm:"-" json:"docs_ids"`
	Theme        ThemeID          `gorm:"-" json:"theme"`
}

func ValidOptionalAxes(val string) bool {
	var validAxesStrId = []string{"0", "1", "2", "3", "4", "5"}
	if val == "" || Contains(validAxesStrId, val) {
		return true
	}
	return false
}

type ThemeID struct {
	ID   int32  `gorm:"-" json:"id"`
	Name string `gorm:"-" json:"name"`
}

func (r *Regles) AfterFind(tx *gorm.DB) (err error) {
	var is []*Iso27002s
	tx.Model("Iso27002s").Joins("join iso_regleses AS i ON i.regles=? AND iso27002s.ID = i.iso", r.ID).Find(&is)
	r.ReglesIso = is
	for _, i := range is {
		r.ReglesIsoIDs = append(r.ReglesIsoIDs, i.ID)
	}

	var ds []*Docs
	tx.Model("Docs").Joins("join doc_regleses AS d ON d.regles=? AND docs.ID = d.doc", r.ID).Find(&ds)
	r.Docs = ds
	for _, d := range ds {
		r.DocsIDs = append(r.DocsIDs, d.ID)
	}

	var t []*Theme
	e := tx.Model("Theme").Joins("join regles_themeses AS d ON d.regle=? and themes.ID = d.th", r.ID).Find(&t).Error
	if e == nil && len(t) != 0 {
		r.Theme.Name = t[0].Name
		r.Theme.ID = t[0].ID
	}
	return nil
}

func (r *Regles) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("regle = ?", r.ID).Delete(ReglesThemeses{})
	tx.Where("regle = ?", r.ID).Delete(ReglesDomaineses{})
	tx.Where("regles = ?", r.ID).Delete(IsoRegleses{})
	tx.Where("regles = ?", r.ID).Delete(DocRegleses{})
	return nil
}

func GetRegles(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetRegles")

	query := "SELECT * FROM regles"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("regles").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("regles").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var regles []Regles
	err := dbmap.Raw(query).Scan(&regles).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, regles)
	} else {
		c.JSON(404, gin.H{"error": "no log(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/regles
}

func GetRegle(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("GetRegle " + id)

	var regle Regles
	err := dbmap.First(&regle, id).Error

	if err == nil {
		c.JSON(200, regle)
	} else {
		c.JSON(404, gin.H{"error": "regle not found"})
	}

	// curl -i http://localhost:8080/api/v1/regles/1
}

func PostRegle(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("PostRegles")
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var regle Regles
	c.Bind(&regle)

	if regle.Name != "" { // XXX Check mandatory fields
		if ValidOptionalStatus(regle.Status) == false {
			c.JSON(400, gin.H{"error": "bad Status"})
			return
		}
		if ValidOptionalAxes(regle.Axe1) == false || ValidOptionalAxes(regle.Axe2) == false {
			c.JSON(400, gin.H{"error": "bad Axe"})
			return
		}

		regle.CreatedOn = time.Now()
		regle.CreatedBy = a.LoginID
		err := dbmap.Create(&regle).Error
		if err == nil {
			c.JSON(201, regle)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/regles
}

func UpdateRegle(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateRegle " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var regle Regles
	err := dbmap.First(&regle, id).Error
	if err == nil {
		var json Regles
		c.Bind(&json)

		//fmt.Printf("UPDATE\n%+v\n", json)
		//TODO : find fields via reflections
		//XXX custom fields mapping
		newregle := Regles{
			Name:        json.Name,
			Ordre:       json.Ordre,
			Descorig:    json.Descorig,
			Description: json.Description,
			Notes:       json.Notes,
			Status:      json.Status,
			Axe1:        json.Axe1,
			Axe2:        json.Axe2,
		}
		if newregle.Name != "" { // XXX Check mandatory fields
			if ValidOptionalStatus(newregle.Status) == false {
				c.JSON(400, gin.H{"error": "bad Status"})
				return
			}
			if ValidOptionalAxes(newregle.Axe1) == false || ValidOptionalAxes(newregle.Axe2) == false {
				c.JSON(400, gin.H{"error": "bad Axe"})
				return
			}

			if EqualArrayIds(regle.DocsIDs, json.DocsIDs) == false {
				dbmap.Where("regles = ?", regle.ID).Delete(DocRegleses{})
				for _, d := range json.DocsIDs {
					dr := DocRegleses{Regles: regle.ID, Doc: d}
					dbmap.Create(&dr)
				}
				newregle.DocsIDs = json.DocsIDs
			}
			if EqualArrayIds(regle.ReglesIsoIDs, json.ReglesIsoIDs) == false {
				dbmap.Where("regles = ?", regle.ID).Delete(IsoRegleses{})
				for _, d := range json.ReglesIsoIDs {
					dr := IsoRegleses{Regles: regle.ID, Iso: d}
					dbmap.Create(&dr)
				}
				newregle.ReglesIsoIDs = json.ReglesIsoIDs
			}

			if err := dbmap.Model(&regle).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":        newregle.Name,
					"Ordre":       newregle.Ordre,
					"Descorig":    newregle.Descorig,
					"Description": newregle.Description,
					"Notes":       newregle.Notes,
					"Status":      newregle.Status,
					"Axe1":        newregle.Axe1,
					"Axe2":        newregle.Axe2,
					"UpdatedOn":   time.Now(),
					"UpdatedBy":   a.LoginID,
				}).Error; err == nil {
				c.JSON(200, newregle)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "regle not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/regles/1
}

func DeleteRegle(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteRegle " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var regle Regles
	err := dbmap.First(&regle, id).Error

	if err == nil {
		e := dbmap.Delete(&regle)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "regle not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/regles/1
}

func GetRegleByDom(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")
	dom := c.Params.ByName("dom")

	a := Auth(c)
	a.Log("GetRegleByDom r: " + id + " d: " + dom)

	var regle Regles
	err := dbmap.First(&regle, id).Error
	if err == nil {
		domint, _ := strconv.Atoi(dom)
		if domint != 0 {
			tdom := dom
			init := 0
			for {
			    var dr ReglesDomaineses
				dbRresult := dbmap.Model("ReglesDomaineses").Where("regle = ? AND domaine = ?", regle.ID, tdom).First(&dr)
				//log.Printf(" => regle by dom : Dom %s, ID %d, Name «%s», Applicable: %d, Exist %v\n", tdom, regle.ID, regle.Name, dr.Applicable, dbRresult.Error)
				if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) == false && dr.Applicable == 1 {
					if init == 0 {
						regle.RegleDomaine = dr
					} else {
						regle.RegleDomaine.Applicable = 0
					}
					regle.RegleDomaine.Conform = dr.Conform
					regle.RegleDomaine.Evolution = dr.Evolution
					break
				}
				//find parent
				var pdom Domaine
				dbmap.Model("Domaine").Where("ID = ?", dom).First(&pdom)
				//log.Printf("  => regle by dom : found parent %d\n", pdom.Parent)

				if pdom.Parent == 0 || strconv.Itoa(int(pdom.Parent)) == tdom {
					break
				}
				tdom = strconv.Itoa(int(pdom.Parent))
				if init == 0 {
					init = 1
				}
			}

		}
		c.JSON(200, regle)

	} else {
		c.JSON(404, gin.H{"error": "regle not found"})
	}

}

/*func prettyPrint(i interface{}) string {
        s, _ := json.MarshalIndent(i, "", "\t")
        return string(s)
}*/

func UpdateRegleByDom(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")
	dom := c.Params.ByName("dom")

	a := Auth(c)
	a.Log("UpdateRegleByDom r: " + id + " d: " + dom)
	if a.Role != "admin" && a.Role != "cssi" {
		c.JSON(403, gin.H{"error": "admin or cssi role required"})
		return
	}

	domint, _ := strconv.Atoi(dom)
	idint, _ := strconv.Atoi(id)
	if domint == 0 || idint == 0 {
		c.JSON(404, gin.H{"error": "domaine and regle required"})
		return
	}

	var rdjson ReglesDomaineses
	c.Bind(&rdjson)
	//fmt.Println(prettyPrint(rdjson))

	var rd ReglesDomaineses
	dbmap.Model("ReglesDomaineses").Where("regle = ? AND domaine = ?", id, dom).Find(&rd)

	new_rd := rd

	if rdjson.Modif == "modif" && a.Role == "admin" {
		new_rd.Modifdesc = rdjson.Modifdesc
		new_rd.Supldesc = rdjson.Supldesc
	}

	if rdjson.Modif == "eval" && (a.Role == "admin" || a.Role == "cssi") {
		if Contains([]string{"0", "1", "2", "3"}, rdjson.Conform) == false {
			rdjson.Conform = "0"
		}
		if Contains([]string{"-1", "0", "1"}, rdjson.Evolution) == false {
			rdjson.Evolution = "0"
		}
		// validate (max 0 evolution fo conform = 3, min 0 evolution for conform 0)
		if rdjson.Applicable == 0 {
			rdjson.Conform = "0"
			rdjson.Evolution = "0"
		} else {
			if rdjson.Conform == "3" && rdjson.Evolution == "1" {
				rdjson.Evolution = "0"
			}
			if rdjson.Conform == "0" && rdjson.Evolution == "-1" {
				rdjson.Evolution = "0"
			}
		}
		new_rd.Applicable = rdjson.Applicable
		new_rd.Conform = rdjson.Conform
		new_rd.Evolution = rdjson.Evolution
	}

	if rd.ID == 0 {
		// Create
		new_rd.Regle = int32(idint)
		new_rd.Domaine = int32(domint)
		new_rd.CreatedOn = time.Now()
		new_rd.CreatedBy = a.LoginID
		//json.RegleDomaine = new_rd
		err := dbmap.Create(&new_rd).Error
		if err == nil {
			c.JSON(201, new_rd)
		} else {
			checkErr(err, "Insert failed")
		}
	} else {
		// Update
		if err := dbmap.Model(&rd).Updates(
			map[string]interface{}{ // gorm don't update null value in struct !!!
				"Regle":      int32(idint),
				"Domaine":    int32(domint),
				"Modifdesc":  new_rd.Modifdesc,
				"Supldesc":   new_rd.Supldesc,
				"Applicable": new_rd.Applicable,
				"Conform":    new_rd.Conform,
				"Evolution":  new_rd.Evolution,
				"UpdatedOn":  time.Now(),
				"UpdatedBy":  a.LoginID,
			}).Error; err == nil {
			//json.RegleDomaine = new_rd
			c.JSON(201, new_rd)
		} else {
			checkErr(err, "Updated failed")
		}
	}

}
