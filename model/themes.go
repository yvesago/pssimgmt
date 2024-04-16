package model

import (
	"encoding/json"
	"errors"
	"fmt"
	//"log"
	"math"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

/*
DB Table Details
-------------------------------------


CREATE TABLE themes (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  name varchar(255) NOT NULL  ,
  ordre integer   ,
  parent integer   ,
  descorig text   ,
  description text   ,
  notes text   ,
  status varchar(255)   ,
  created_by integer   ,
  created_on timestamp   ,
  updated_by integer   ,
  updated_on timestamp
)

*/

// Themes struct is a row record of the themes table in the pssimgmt database
type Theme struct {
	Model
	Name        string `gorm:"column:name;type:varchar;size:255;" json:"name"`
	Ordre       int32  `gorm:"column:ordre;type:integer;" json:"ordre"`
	Parent      int32  `gorm:"column:parent;type:integer;" json:"parent"`
	Descorig    string `gorm:"column:descorig;type:text;" json:"descorig"`
	Description string `gorm:"column:description;type:text;" json:"description"`
	Notes       string `gorm:"column:notes;type:text;" json:"notes"`
	Status      string `gorm:"column:status;type:varchar;size:255;" json:"status"`
	//Conforme     [6]int       `gorm:"-" json:"conforme"`
	//Evolution    [6]int       `gorm:"-" json:"evolution"`
	//Axes         [6]string    `gorm:"-" json:"axes"`
	Evaluation   Evaluation   `gorm:"-" json:"evaluation"`
	Regles       []*Regles    `gorm:"-" json:"regles"` // gorm:"ForeignKey:Regle;AssociationForeignKey:regle"`
	ReglesIDs    []int32      `gorm:"-" json:"regles_ids"`
	IsoThemes    []*Iso27002s `gorm:"-" json:"iso_themes"`
	IsoThemesIDs []int32      `gorm:"-" json:"iso_ids"`
	Children     []*Theme     `gorm:"-" json:"children"`
	Dom          string       `gorm:"-"` // interal use for thmeme by domaine
}

type Evaluation struct {
	Axes      [6]string `gorm:"-" json:"axes"`
	Conforme  [6]int    `gorm:"-" json:"conforme"`
	Evolution [6]int    `gorm:"-" json:"evolution"`
}

var ValidAxes = [6]string{
	"Gouvernance",
	"Maîtrise des risques",
	"Maîtrise des systèmes",
	"Protection des systèmes",
	"Gestion des incidents",
	"Évaluation",
}

var ValidStatus = []string{
	"ok",
	"nok",
	"etu",
}

func ValidOptionalStatus(val string) bool {
	if val == "" || Contains(ValidStatus, val) {
		return true
	}
	return false
}

func (t *Theme) AfterFind(tx *gorm.DB) (err error) {
	//fmt.Println(t)
	var is []*Iso27002s
	tx.Model("Iso27002s").Joins("join iso_themes AS i ON i.th=? AND iso27002s.ID = i.iso", t.ID).Find(&is)
	t.IsoThemes = is
	for _, i := range is {
		t.IsoThemesIDs = append(t.IsoThemesIDs, i.ID)
	}

	var rsint []int32
	tx.Raw("SELECT r.id from regles AS r JOIN regles_themeses AS rt ON r.ID = rt.regle AND rt.th=?", t.ID).Pluck("DISTINCT(ID)", &rsint)
	/*for _, r := range rsint {
		t.ReglesIDs = append(t.ReglesIDs, r)
	}*/
	t.ReglesIDs = rsint
	/*var rs []*Regles
	tx.Model("Regles").Joins("JOIN regles_themeses AS rt ON regles.ID = rt.regle AND rt.th=?", t.ID).Find(&rs)
	for i, r := range rs {
		rs[i].Theme.Name = t.Name
		rs[i].Theme.ID = t.ID
		t.ReglesIDs = append(t.ReglesIDs, r.ID)
	}
	t.Regles = rs*/

	if t.Dom == "" {
		return nil
	}

	var rs []*Regles
	tx.Model("Regles").Joins("JOIN regles_themeses AS rt ON regles.ID = rt.regle AND rt.th=?", t.ID).Find(&rs)
	//dbmap.Raw("SELECT * from regles AS r JOIN regles_themeses AS rt ON r.ID = rt.regle AND rt.th=?", id).Find(&rs)

	for i, r := range rs {
		if rs[i].Status != "ok" {
			continue
		}
		rs[i].Theme.Name = t.Name
		rs[i].Theme.ID = t.ID
		t.ReglesIDs = append(t.ReglesIDs, r.ID)
		tdom := t.Dom
		init := 0
		for {
			var dr ReglesDomaineses
			dbRresult := tx.Model("ReglesDomaineses").Where("regle = ? AND domaine = ?", r.ID, tdom).Find(&dr)
			//log.Printf(" themes by dom  => ID %d, Name «%s», Applicable: %d, Exist %v, Parent %d\n", t.ID, t.Name, dr.Applicable, dbRresult.Error, t.Parent)
			if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) == false && dr.Applicable == 1 {
				if init == 0 {
					r.RegleDomaine = dr
				} else {
					r.RegleDomaine.Applicable = 0
				}
				r.RegleDomaine.Conform = dr.Conform
				r.RegleDomaine.Evolution = dr.Evolution
				break
			}
			//find parent
			var pdom Domaine
			tx.Model("Domaine").Where("ID = ?", tdom).First(&pdom)
			//log.Printf(" dom %v, Parent %d\n", tdom, pdom.Parent)

			if pdom.Parent == 0 || strconv.Itoa(int(pdom.Parent)) == tdom {
				break
			}
			tdom = strconv.Itoa(int(pdom.Parent))
			if init == 0 {
				init = 1
			}
		}

	}
	t.Regles = rs

	conforme := 0
	evolution := 0
	axes := make(map[string]Evaluation)
	var rsok_count [7]int
	for _, r := range rs {
		if r.Status != "ok" {
			continue
		}
		rd := r.RegleDomaine
		/*if rd.Applicable != 1 { //XXX disable inherit
			continue
		}*/
		//fmt.Printf("Conform: %s\nEvolution: %s\n Axe 1: %s\n Axe 2: %s\n", rd.Conform, rd.Evolution, r.Axe1, r.Axe2)
		cint, _ := strconv.Atoi(rd.Conform)
		eint, _ := strconv.Atoi(rd.Evolution)

		if r.Axe1 != "" {
			a1int, _ := strconv.Atoi(r.Axe1)
			a1int--
			if a1int >= 0 && a1int < 6 {
				rsok_count[a1int]++
				t := axes[r.Axe1]
				t.Conforme[a1int] += cint
				t.Evolution[a1int] += eint
				axes[r.Axe1] = t
			}
		}

		if r.Axe2 != "" {
			a2int, _ := strconv.Atoi(r.Axe2)
			a2int--
			if a2int >= 0 && a2int < 6 {
				rsok_count[a2int]++
				t2 := axes[r.Axe2]
				t2.Conforme[a2int] += cint
				t2.Evolution[a2int] += eint
				axes[r.Axe2] = t2
			}
		}

		conforme += cint
		evolution += eint
	}
	//fmt.Printf(" => Conform: %d, Evolution: %d\n", conforme, evolution)
	//fmt.Printf(" => %+v\n", axes)
	//fmt.Println(prettyPrint(axes))

	var te Evaluation
	for ai, a := range axes {
		aint, _ := strconv.Atoi(ai)
		aint--
		//fmt.Printf("++ %d : %v => %s\n", aint, a, Axes[aint])
		if aint >= 0 && aint < 6 {
			te.Axes[aint] = ValidAxes[aint]
			te.Conforme[aint] = 10 * a.Conforme[aint] / (3 * rsok_count[aint])
			te.Evolution[aint] = 10 * (a.Conforme[aint] + a.Evolution[aint]) / (3 * rsok_count[aint])
		}
	}
	//fmt.Println(prettyPrint(t))
	//fmt.Printf(" => %+v\n", te)
	t.Evaluation = te

	return nil
}

func (t *Theme) BeforeDelete(tx *gorm.DB) (er error) {
	tx.Model(t).Where("parent = ?", t.ID).Update("parent", 0)
	tx.Where("th = ?", t.ID).Delete(ReglesThemeses{})
	tx.Where("th = ?", t.ID).Delete(IsoThemes{})
	return nil
}

type ByOrdre []*Theme

func (e ByOrdre) Len() int           { return len(e) }
func (e ByOrdre) Less(i, j int) bool { return e[i].Ordre < e[j].Ordre }
func (e ByOrdre) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func GetThemes(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	a := Auth(c)
	a.Log("GetThemes")

	query := "SELECT * FROM themes"
	// Parse query string
	//  receive : map[_sort:[id] q:[wx] _order:[DESC] _start:[0] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	count := 0
	if s != "" {
		dbmap.Table("themes").Where(s).Count(&count)
		query = query + " WHERE " + s
	} else {
		dbmap.Table("themes").Count(&count)
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	var themes []Theme
	err := dbmap.Raw(query).Scan(&themes).Error

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(count))
		c.JSON(200, themes)
	} else {
		c.JSON(404, gin.H{"error": "no theme(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/themes
}

func GetThemesTree(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	query := "SELECT * FROM themes"
	dom := c.Params.ByName("dom")
	domint, _ := strconv.Atoi(dom)

	a := Auth(c)
	a.Log("GetThemesTree")
	//login, _ := c.Get("Login")
	//role, _ := c.Get("Role")
	//log.Printf("  ==> %s - %s <==\n", login, role)

	var themes []*Theme
	err := dbmap.Raw(query).Scan(&themes).Error

	if err == nil {

		tmap := make(map[int32]*Theme)
		tree := &Theme{Name: "PSSI", Evaluation: Evaluation{Axes: ValidAxes}}
		var evalcount [6]int
		var conform [6]int
		var evolution [6]int
		//tree := &domaines[0]
		for _, e := range themes {
			tmap[e.ID] = e
			if e.Parent == 0 {
				tree.Children = append(tree.Children, e)
			}
			if domint != 0 {
				var t Theme
				t.Dom = dom
				dbmap.First(&t, e.ID)
				//fmt.Printf("%+v\n", t.Evaluation)
				for i, e := range t.Evaluation.Axes {
					if e != "" {
						evalcount[i]++
					}
				}
				for i, _ := range t.Evaluation.Axes {
					conform[i] += t.Evaluation.Conforme[i]
					evolution[i] += t.Evaluation.Evolution[i] //(t.Evaluation.Evolution[i] - t.Evaluation.Conforme[i])
				}
			}

		}

		if domint != 0 {
			//fmt.Printf("%+v %+v %+v\n", evalcount, conform, evolution)
			for i, _ := range conform {
				if evalcount[i] != 0 {
					conform[i] = int(math.RoundToEven(float64(conform[i]) / float64(evalcount[i])))
					evolution[i] = int(math.RoundToEven(float64(evolution[i]) / float64(evalcount[i])))
				}
			}
			tree.Evaluation.Conforme = conform
			tree.Evaluation.Evolution = evolution
			//fmt.Printf("%d %+v %+v\n", evalcount, conform, evolution)
		}

		for _, e := range tmap {
			if tmap[e.Parent] != nil {
				tmap[e.Parent].Children = append(tmap[e.Parent].Children, e)
				sort.Sort(ByOrdre(tmap[e.Parent].Children))
			}
		}

		sort.Sort(ByOrdre(tree.Children))
		bytes, _ := json.Marshal(tree)

		//c.Header("X-Total-Count", "1")
		//c.JSON(200, bytes)
		c.Data(200, "application/json", bytes)
	} else {
		c.JSON(404, gin.H{"error": "no theme(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/domaines

}

func GetThemeByDom(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")
	dom := c.Params.ByName("dom")

	a := Auth(c)
	a.Log("GetThemeByDom " + id + " dom: " + dom)

	var theme Theme
	theme.Dom = dom

	err := dbmap.First(&theme, id).Error

	if err == nil {
		c.JSON(200, theme)
	} else {
		c.JSON(404, gin.H{"error": "theme not found"})
	}

	// curl -i http://localhost:8080/api/v1/themes/1
}

func GetTheme(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	var theme Theme
	err := dbmap.First(&theme, id).Error

	if err == nil {
		c.JSON(200, theme)
	} else {
		c.JSON(404, gin.H{"error": "theme not found"})
	}

	// curl -i http://localhost:8080/api/v1/themes/1
}

func PostTheme(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)

	var theme Theme
	c.Bind(&theme)

	a := Auth(c)
	a.Log(fmt.Sprintf("PostTheme by id: %d", a.LoginID))
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	if theme.Name != "" { // XXX Check mandatory fields
		if ValidOptionalStatus(theme.Status) == false {
			c.JSON(400, gin.H{"error": "bad Status"})
			return
		}
		theme.CreatedOn = time.Now()
		theme.CreatedBy = a.LoginID
		err := dbmap.Create(&theme).Error
		if err == nil {
			c.JSON(201, theme)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}
	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/themes
}

func UpdateTheme(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("UpdateTheme " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var theme Theme
	err := dbmap.First(&theme, id).Error
	if err == nil {
		var json Theme
		c.Bind(&json)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		newtheme := Theme{
			Name:        json.Name,
			Parent:      json.Parent,
			Descorig:    json.Descorig,
			Description: json.Description,
			Notes:       json.Notes,
			Status:      json.Status,
		}
		if newtheme.Name != "" { // XXX Check mandatory fields
			if ValidOptionalStatus(newtheme.Status) == false {
				c.JSON(400, gin.H{"error": "bad Status"})
				return
			}
			if EqualArrayIds(theme.ReglesIDs, json.ReglesIDs) == false {
				dbmap.Where("th = ?", theme.ID).Delete(ReglesThemeses{})
				for _, d := range json.ReglesIDs {
					dr := ReglesThemeses{Th: theme.ID, Regle: d}
					dbmap.Create(&dr)
				}
				newtheme.ReglesIDs = json.ReglesIDs
			}
			if EqualArrayIds(theme.IsoThemesIDs, json.IsoThemesIDs) == false {
				dbmap.Where("th = ?", theme.ID).Delete(IsoThemes{})
				for _, d := range json.IsoThemesIDs {
					dr := IsoThemes{Th: theme.ID, Iso: d}
					dbmap.Create(&dr)
				}
				newtheme.IsoThemesIDs = json.IsoThemesIDs
			}

			if err := dbmap.Model(&theme).Updates(
				map[string]interface{}{ // gorm don't update null value in struct !!!
					"Name":        newtheme.Name,
					"Parent":      json.Parent,
					"Descorig":    json.Descorig,
					"Description": json.Description,
					"Notes":       json.Notes,
					"Status":      json.Status,
					"UpdatedBy":   a.LoginID,
					"UpdatedOn":   time.Now(),
				}).Error; err == nil {
				c.JSON(200, newtheme)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "theme not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/themes/1
}

func DeleteTheme(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorm.DB)
	id := c.Params.ByName("id")

	a := Auth(c)
	a.Log("DeleteTheme " + id)
	if a.Role != "admin" {
		c.JSON(403, gin.H{"error": "admin role required"})
		return
	}

	var theme Theme
	err := dbmap.First(&theme, id).Error

	if err == nil {
		e := dbmap.Delete(&theme)

		if e != nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(nil, "Delete failed")
		}
	} else {
		c.JSON(404, gin.H{"error": "theme not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/themes/1
}
