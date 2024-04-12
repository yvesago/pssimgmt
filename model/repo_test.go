package model

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sort"
)

// Set test config
var config = Config{
	DBname:  "_test.sqlite",
	Verbose: false, // Set with cmd line in release mode
}

func deleteFile(file string) {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func SetConfig(config Config, a AuthInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Login", a.Login)     // to test acces rights
		c.Set("LoginID", a.LoginID) // to test acces rights
		c.Set("Role", a.Role)
		c.Set("Verbose", config.Verbose)
		c.Set("DBmap", config.DBh)
		c.Next()
	}
}

type Element struct {
	Id       int32  `json:"id"` // <--- Field tags
	Name     string `json:"name"`
	Order    int32  `json:"order"`
	parent   int32
	P        *Element
	Children []*Element `json:"chidren"`
}

type ByOrder []*Element

func (e ByOrder) Len() int           { return len(e) }
func (e ByOrder) Less(i, j int) bool { return e[i].Order < e[j].Order }
func (e ByOrder) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (el *Element) String() string {
	s := ""
	if el.parent != 0 {
		s += "\t"
		if el.P != nil && el.P.parent != 0 {
			s += "\t"
		}
	}
	//s += fmt.Sprintf("(%d) %s\n",el.Order, el.Name)
	s += fmt.Sprintf("%s\n", el.Name)
	sort.Sort(ByOrder(el.Children))
	for _, child := range el.Children {
		s += child.String()
	}
	//s += "</" + el.tag + ">"
	return s
}

func (t *Theme) String() string {
	s := ""
	if t.Parent != 0 {
		s += "\t"
	}
	s += fmt.Sprintf("[%d] %s %d\n", t.ID, t.Name, t.Parent)
	for _, child := range t.Children {
		s += child.String()
	}
	if t.Regles != nil {
		for _, r := range t.Regles {
			s += r.String()
		}
	}
	return s
}

func (r *Regles) String() string {
	s := fmt.Sprintf("\t\t(%d) %s\n", r.ID, r.Name)
	if r.Descorig != "" {
		s += fmt.Sprintf("\t\t%s\n", r.Descorig)
	}
	return s
}
