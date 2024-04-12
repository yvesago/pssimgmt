package model

/*import (
	"github.com/jinzhu/gorm"
)*/

/*
DB Table Details
-------------------------------------


CREATE TABLE iso_regleses (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  iso integer NOT NULL  ,
  regles integer NOT NULL
)

JSON Sample
-------------------------------------
{    "id": 24,    "iso": 22,    "regles": 46}



*/

// IsoRegleses struct is a row record of the iso_regleses table in the pssimgmt database
type IsoRegleses struct {
	ID     int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Iso    int32 `gorm:"column:iso;type:integer;" json:"iso"`
	Regles int32 `gorm:"column:regles;type:integer;" json:"regles"`
}
