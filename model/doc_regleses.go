package model

/*import (
)*/


/*
DB Table Details
-------------------------------------


CREATE TABLE doc_regleses (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  doc integer NOT NULL  ,
  regles integer NOT NULL
)

JSON Sample
-------------------------------------
{    "id": 86,    "doc": 33,    "regles": 81}



*/

// DocRegleses struct is a row record of the doc_regleses table in the pssimgmt database
type DocRegleses struct {
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Doc int32 `gorm:"column:doc;type:integer;" json:"doc"`
	Regles int32 `gorm:"column:regles;type:integer;" json:"regles"`
}

