package model

import (
	"time"
)

/*
DB Table Details
-------------------------------------

CREATE TABLE regles_domaineses (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL  ,
  regle integer NOT NULL  ,
  domaine integer NOT NULL  ,
  modifdesc text   ,
  supldesc text
, evolution varchar(255)  DEFAULT 0, conform varchar(255), applicable integer  DEFAULT 1, created_by integer, created_on timestamp, updated_by integer, updated_on timestamp)

-------------------------------------

*/

// ReglesDomaineseses struct is a row record of the regles_domaineses table in the pssimgmt database
type ReglesDomaineses struct {
	ID         int32     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:integer;" json:"id"`
	Regle      int32     `gorm:"column:regle;type:integer;" json:"regle"`
	Domaine    int32     `gorm:"column:domaine;type:integer;" json:"domaineID"`
	Modifdesc  string    `gorm:"column:modifdesc;type:text;" json:"modifdesc"`
	Supldesc   string    `gorm:"column:supldesc;type:text;" json:"supldesc"`
	Evolution  string    `gorm:"column:evolution;type:varchar;size:255;default:0;" json:"evolution"`
	Conform    string    `gorm:"column:conform;type:varchar;size:255;" json:"conform"`
	Applicable int64     `gorm:"column:applicable;type:integer;default:1;" json:"applicable"`
	CreatedBy  int64     `gorm:"column:created_by;type:integer;" json:"created_by"`
	CreatedOn  time.Time `gorm:"column:created_on;type:timestamp;" json:"created_on"`
	UpdatedBy  int64     `gorm:"column:updated_by;type:integer;" json:"updated_by"`
	UpdatedOn  time.Time `gorm:"column:updated_on;type:timestamp;" json:"updated_on"`
	Domaines   Domaine   `gorm:"-;" json:"domain"`
	Modif      string    `gorm:"-;" json:"modif"`
}
