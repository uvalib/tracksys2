package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type masterFileHit struct {
	ID          uint64 `json:"id"`
	PID         string `gorm:"column:pid" json:"pid"`
	Filename    string `json:"filename"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type metadataHit struct {
	ID          uint64 `json:"id"`
	PID         string `gorm:"column:pid" json:"pid"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	CallNumber  string `json:"callNumber"`
	Barcode     string `json:"barcode"`
	CatalogKey  string `json:"catalogKey"`
	CreatorName string `json:"creatorName"`
}

type orderHit struct {
	ID                  uint64   `json:"id"`
	OrderTitle          string   `json:"title"`
	StaffNotes          string   `json:"notes"`
	SpecialInstructions string   `json:"specialInstructions"`
	CustomerID          uint     `json:"-"`
	Customer            customer `gorm:"foreignKey:CustomerID" json:"customer"`
	AgencyID            *uint    `json:"-"`
	Agency              *agency  `gorm:"foreignKey:AgencyID" json:"agency"`
}

func (orderHit) TableName() string {
	return "orders"
}

type searchResults struct {
	MasterFiles []masterFileHit `json:"masterFiles"`
	Metadata    []metadataHit   `json:"metadata"`
	Orders      []orderHit      `json:"orders"`
}

func (svc *serviceContext) searchRequest(c *gin.Context) {
	qStr := c.Query("q")
	matchStart := fmt.Sprintf("%s%%", qStr)
	matchAny := fmt.Sprintf("%%%s%%", qStr)
	scope := c.Query("scope")
	if scope == "" {
		scope = "all"
	}
	if scope != "all" && scope != "masterfiles" && scope != "metadata" && scope != "orders" {
		log.Printf("ERROR: invalid search scope %s specified", scope)
		c.String(http.StatusBadRequest, "invalid search scope")
		return
	}
	field := c.Query("field")
	if field == "" {
		field = "all"
	}
	log.Printf("INFO: search %s.%s for [%s]", scope, field, qStr)
	hitLimit := 30
	resp := searchResults{MasterFiles: make([]masterFileHit, 0), Metadata: make([]metadataHit, 0), Orders: make([]orderHit, 0)}

	if scope == "all" || scope == "masterfiles" {
		searchQ := svc.DB.Table("master_files")
		if field == "all" {
			searchQ = searchQ.Where(
				svc.DB.Where("master_files.id=?", qStr).Or("pid=?", qStr).
					Or("filename like ?", matchStart).Or("title like ?", matchAny).
					Or("description like ?", matchAny),
			)
		} else if field == "pid" {
			searchQ = searchQ.Where("pid=?", qStr)
		} else if field == "title" || field == "description" || field == "filename" {
			searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
		} else if field == "tag" {
			searchQ = searchQ.
				Joins("left outer join master_file_tags mt on mt.master_file_id = master_files.id").
				Joins("left outer join tags t on mt.tag_id = t.id").
				Where("t.tag=?", qStr)
		}
		err := searchQ.Limit(hitLimit).Find(&resp.MasterFiles).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
	}

	if scope == "all" || scope == "metadata" {
		searchQ := svc.DB.Table("metadata")
		if field == "all" {
			searchQ = searchQ.Where(
				svc.DB.Where("id=?", qStr).Or("pid=?", qStr).Or("title like ?", matchAny).
					Or("barcode=?", qStr).Or("catalog_key=?", qStr).Or("call_number like ?", matchStart).
					Or("creator_name like ?", matchAny),
			).Limit(hitLimit)
		} else {
			if field == "title" || field == "creator_name" {
				searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny).Limit(hitLimit)
			} else if field == "call_number" {
				searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchStart).Limit(hitLimit)
			} else {
				searchQ = searchQ.Where(fmt.Sprintf("%s=?", field), qStr).Limit(hitLimit)
			}
		}
		err := searchQ.Limit(hitLimit).Find(&resp.Metadata).Error
		if err != nil {
			log.Printf("ERROR: metadata search failed: %s", err.Error())
		}
	}

	if scope == "all" || scope == "orders" {
		searchQ := svc.DB.Table("orders").Preload("Customer").Preload("Agency")
		if field == "all" {
			searchQ = searchQ.
				Joins("inner join customers on customer_id = customers.id").
				Joins("left outer join agencies on agency_id = agencies.id").
				Joins("left outer join units on units.order_id = orders.id").
				Where(
					svc.DB.Where("orders.id=?", qStr).Or("units.id=?", qStr).
						Or("customers.last_name like ?", matchStart).Or("agencies.name=?", qStr).
						Or("orders.staff_notes like ?", matchAny).Or("orders.special_instructions like ?", matchAny),
				)
		} else if field == "id" {
			searchQ = searchQ.Where("id=?", qStr)
		} else if field == "unit_id" {
			searchQ = searchQ.Joins("left outer join units on units.order_id = orders.id").
				Where("units.id=?", qStr)
		} else if field == "last_name" {
			searchQ = searchQ.Joins("inner join customers on customer_id = customers.id").
				Where("customers.last_name like ?", matchStart)
		} else if field == "agency" {
			searchQ = searchQ.Joins("left outer join agencies on agency_id = agencies.id").
				Where("agencies.name like ?", matchAny)
		} else {
			searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
		}
		err := searchQ.Limit(hitLimit).Find(&resp.Orders).Error
		if err != nil {
			log.Printf("ERROR: order search failed: %s", err.Error())
		}
	}

	c.JSON(http.StatusOK, resp)
}
