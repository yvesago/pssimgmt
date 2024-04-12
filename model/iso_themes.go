package model

/*
DB Table Details
-------------------------------------


CREATE TABLE iso_themes (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  iso integer NOT NULL  ,
  th integer NOT NULL
)

JSON Sample
-------------------------------------
{    "id": 79,    "iso": 14,    "th": 75}



*/

// IsoThemes struct is a row record of the iso_themes table in the pssimgmt database
type IsoThemes struct {
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Iso int32 `gorm:"column:iso;type:integer;" json:"iso"`
	Th int32 `gorm:"column:th;type:integer;" json:"th"`
}

