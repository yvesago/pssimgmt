package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// override gorm.Model
type Model struct {
	ID        int32     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	CreatedBy int64     `gorm:"column:created_by;type:integer;" json:"created_by"`
	CreatedOn time.Time `gorm:"column:created_on;type:timestamp;" json:"created_on"`
	UpdatedBy int64     `gorm:"column:updated_by;type:integer;" json:"updated_by"`
	UpdatedOn time.Time `gorm:"column:updated_on;type:timestamp;" json:"updated_on"`
	//  DeletedAt *time.Time
}

type Config struct {
	Port        string
	CorsOrigin  string
	AuthURL     string
	CallbackURL string
	JwtPass     string
	DBname      string
	DBh         *gorm.DB
	Verbose     bool
	Debug       bool
}

// Init DBs return the new db handler DBh in config
func DBs(config Config) Config {
	config.DBh = InitDb(config.DBname, config.Verbose)
	return config
}

func InitDb(dbName string, verbose bool) *gorm.DB {

	//f := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbName)
	f := dbName
	//fmt.Println(f)
	db, err := gorm.Open(
		"sqlite3",
		f,
	)
	//defer db.Close()
	if err != nil {
		checkErr(nil, "Open failed")
	}

	db.LogMode(verbose)

	sqlDB := db.DB()
	//sqlDB.SetMaxOpenConns(10)
	//sqlDB.SetConnMaxLifetime(time.Hour)
	if verbose {
		fmt.Printf("%#v\n", sqlDB.Stats())
	}

	db.AutoMigrate(
		&Regles{}, &Theme{},
		&Docs{}, &DocRegleses{},
		&Domaine{}, &ReglesThemeses{},
		&Iso27002s{}, &IsoRegleses{},
		&Users{}, &IsoThemes{},
                &Versions{}, &Documents{},
		&ReglesDomaineses{})

	db.DB().Ping()
	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func ParseQuery(q map[string][]string) (string, string, string) {
	query := ""
	//fmt.Println(q)
	var andSearches []string
	for col, search := range q {
		if string([]rune(col)[0]) == "_" {
			continue
		}
		valid := regexp.MustCompile("^[A-Za-z0-9_.]+$")
		invalidSearch := regexp.MustCompile("[\\'\")(&`]")
		var searches []string
		for _, sencoded := range search {
			sdecoded, _ := url.QueryUnescape(sencoded)
			s := html.EscapeString(sdecoded)
			if col != "" && s != "" && valid.MatchString(col) && invalidSearch.MatchString(s) == false {
				switch {
				case col == "q":
					// XXX trick for code array autocomplete
					searches = append(searches, "code LIKE \"%"+s+"%\"")
				case col == "casid":
					// XXX trick for cas_id json
					searches = append(searches, "cas_id LIKE \"%"+s+"%\"")
				case col == "domaine_id":
					// XXX trick for domaine_id json
					searches = append(searches, "domaine LIKE \"%"+s+"%\"")
				case col == "user_1":
					searches = append(searches, "user1 = \""+s+"\"")
				case col == "id":
					searches = append(searches, "id = \""+s+"\"")
				case col == "axe":
					searches = append(searches, "axe1 = \""+s+"\" OR axe2 = \""+s+"\"")
				case col == "user":
					searches = append(searches, "user1 = \""+s+"\" OR user2 = \""+s+"\" OR user3 = \""+s+"\"")
				case col == "descriptions":
					searches = append(searches, "descorig LIKE \"%"+s+"%\" OR description LIKE \"%"+s+"%\"")
				default:
					searches = append(searches, col+" LIKE \"%"+s+"%\"")
				}
			}
		}
		if len(searches) > 0 {
			andSearches = append(andSearches, "("+strings.Join(searches, " OR ")+")")
		}

	}
	if len(andSearches) > 0 {
		query = query + " " + strings.Join(andSearches, " AND ")
	}

	order := ""
	if q["_sort"] != nil && q["_order"] != nil {
		sortField := q["_sort"][0]
		if sortField == "casid" { // XXX trick for cas_id json
			sortField = "cas_id"
		}
		if sortField == "user_1" {
			sortField = "user1"
		}
		if sortField == "user_2" {
			sortField = "user2"
		}
		if sortField == "user_2" {
			sortField = "user2"
		}
		// prevent SQLi
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if !valid.MatchString(sortField) {
			sortField = ""
		}
		/*if sortField == "created" || sortField == "updated" { // XXX trick for sqlite
			sortField = "datetime(" + sortField + "_on)"
		}*/
		sortOrder := q["_order"][0]
		if sortOrder != "ASC" {
			sortOrder = "DESC"
		}
		if sortField != "" {
			order = " ORDER BY " + sortField + " " + sortOrder
		}
	}

	limit := ""
	if q["_start"] != nil && q["_end"] != nil {
		start := q["_start"][0]
		end := q["_end"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		startInt, _ := strconv.Atoi(start)
		endInt, _ := strconv.Atoi(end)
		//startInt = startInt - 1 // indice start from 0

		if valid.MatchString(start) && valid.MatchString(end) && endInt > startInt {
			size := endInt - startInt
			//query = query + " LIMIT " + strconv.Itoa(startInt) + ", " + strconv.Itoa(size)

			limit = " LIMIT " + strconv.Itoa(size) + " OFFSET " + start
		}
	}

	return query, order, limit
}

// Auth
type AuthInfo struct {
	IP      string
	Role    string
	Login   string
	LoginID int64
	Doms    []int64
}

func (a AuthInfo) Log(message string) {
	log.Printf("[%s] %s (%s) : %s\n", a.IP, a.Login, a.Role, message)
}

func Auth(c *gin.Context) AuthInfo {
	var a AuthInfo
	a.IP = c.ClientIP()

	if login, ok := c.Get("Login"); ok {
		a.Login = login.(string)
	}

	if loginid, ok := c.Get("LoginID"); ok {
		a.LoginID = loginid.(int64)
	}

	if role, ok := c.Get("Role"); ok {
		a.Role = role.(string)
	}

	if doms, ok := c.Get("Doms"); ok {
		a.Doms = doms.([]int64)
	}

	return a
}

// Tools
func EqualArrayIds(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
