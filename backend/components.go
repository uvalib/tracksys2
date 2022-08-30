package main

import "time"

type componentType struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Descrition string `json:"description"`
}

type component struct {
	ID              int64         `json:"id"`
	PID             string        `gorm:"column:pid" json:"pid"`
	Title           string        `json:"title"`
	Label           string        `json:"label"`
	ContentDesc     string        `json:"description"`
	Date            string        `json:"date"`
	Level           string        `json:"level"`
	Barcode         string        `json:"barcode"`
	EadIDAtt        string        `gorm:"column:ead_id_att" json:"eadID"`
	Ancestry        string        `json:"ancestry"`
	ComponentTypeID int64         `json:"-"`
	ComponentType   componentType `gorm:"foreignKey:ComponentTypeID" json:"componentType"`
	CreatedAt       time.Time     `json:"-"`
	UpdatedAt       time.Time     `json:"-"`
}
