package model

/*import (
        "time"

        "strconv"
        "github.com/gin-gonic/gin"
        "github.com/jinzhu/gorm"

)*/

/*
DB Table Details
-------------------------------------


CREATE TABLE regles_themeses (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  regle integer NOT NULL  ,
  th integer NOT NULL
)

JSON Sample
-------------------------------------
{    "id": 2,    "regle": 82,    "th": 4}



*/

// ReglesThemeses struct is a row record of the regles_themeses table in the pssimgmt database
type ReglesThemeses struct {
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Regle int32 `gorm:"column:regle;type:integer;" json:"regle"`
	Th int32 `gorm:"column:th;type:integer;" json:"th"`
}

/*
// BeforeSave invoked before saving, return an error if field is not populated.
func (r *ReglesThemeses) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *ReglesThemeses) Prepare() {
}
*/
