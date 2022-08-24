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
	log.Printf("INFO: search %s for [%s]", scope, qStr)
	hitLimit := 25

	resp := searchResults{MasterFiles: make([]masterFileHit, 0), Metadata: make([]metadataHit, 0), Orders: make([]orderHit, 0)}
	if scope == "all" || scope == "masterfiles" {
		searchQ := svc.DB.Debug().Table("master_files").
			Joins("left outer join master_file_tags mt on mt.master_file_id = master_files.id").
			Joins("left outer join tags t on mt.tag_id = t.id").
			Where(
				svc.DB.Where("pid=?", qStr).Or("t.tag=?", qStr).
					Or("filename like ?", matchStart).Or("title like ?", matchAny).
					Or("description like ?", matchAny),
			)
		err := searchQ.Limit(hitLimit).Find(&resp.MasterFiles).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
	}

	if scope == "all" || scope == "metadata" {
		err := svc.DB.Debug().Table("metadata").
			Where(
				svc.DB.Where("pid=?", qStr).Or("title like ?", matchAny).
					Or("barcode=?", qStr).Or("catalog_key=?", qStr).Or("call_number like ?", matchStart).
					Or("creator_name like ?", matchAny),
			).
			Limit(hitLimit).Find(&resp.Metadata).Error
		if err != nil {
			log.Printf("ERROR: metadata search failed: %s", err.Error())
		}
	}

	if scope == "all" || scope == "orders" {
		err := svc.DB.Debug().Table("orders").Preload("Customer").Preload("Agency").
			Joins("inner join customers c on customer_id = c.id").
			Joins("left outer join agencies a on agency_id = a.id").
			Where(
				svc.DB.Where("c.last_name like ?", matchStart).Or("a.name=?", qStr).
					Or("staff_notes like ?", matchAny).Or("special_instructions like ?", matchAny),
			).
			Limit(hitLimit).Find(&resp.Orders).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
	}

	c.JSON(http.StatusOK, resp)
}
