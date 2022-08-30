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

type componentResp struct {
	Total int64       `json:"total"`
	Hits  []component `json:"hits"`
}
type metadataResp struct {
	Total int64         `json:"total"`
	Hits  []metadataHit `json:"hits"`
}
type masterFileResp struct {
	Total int64           `json:"total"`
	Hits  []masterFileHit `json:"hits"`
}
type orderResp struct {
	Total int64      `json:"total"`
	Hits  []orderHit `json:"hits"`
}

type searchResults struct {
	Components  componentResp  `json:"components"`
	MasterFiles masterFileResp `json:"masterFiles"`
	Metadata    metadataResp   `json:"metadata"`
	Orders      orderResp      `json:"orders"`
}

func (svc *serviceContext) searchRequest(c *gin.Context) {
	qStr := c.Query("q")
	matchStart := fmt.Sprintf("%s%%", qStr)
	matchAny := fmt.Sprintf("%%%s%%", qStr)
	scope := c.Query("scope")
	if scope == "" {
		scope = "all"
	}
	if scope != "all" && scope != "masterfiles" && scope != "metadata" && scope != "orders" && scope != "components" {
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
	resp := searchResults{
		Components:  componentResp{Hits: make([]component, 0)},
		MasterFiles: masterFileResp{Hits: make([]masterFileHit, 0)},
		Metadata:    metadataResp{Hits: make([]metadataHit, 0)},
		Orders:      orderResp{Hits: make([]orderHit, 0)}}

	if scope == "all" || scope == "components" {
		searchQ := svc.DB.Table("components")
		if field == "all" {
			searchQ = searchQ.Where(
				svc.DB.Where("components.id=?", qStr).Or("pid=?", qStr).
					Or("title like ?", matchAny).Or("label like ?", matchAny).Or("content_desc like ?", matchAny).
					Or("date like ?", matchAny).Or("ead_id_att=?", qStr),
			)
		} else if field == "id" || field == "pid" || field == "ead_id_att" {
			searchQ = searchQ.Where(fmt.Sprintf("%s=?", field), qStr)
		} else {
			searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
		}
		searchQ.Count(&resp.Components.Total)
		err := searchQ.Limit(hitLimit).Find(&resp.Components.Hits).Error
		if err != nil {
			log.Printf("ERROR: component search failed: %s", err.Error())
		}
	}

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
		}
		searchQ.Count(&resp.MasterFiles.Total)
		err := searchQ.Limit(hitLimit).Find(&resp.MasterFiles.Hits).Error
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

		searchQ.Count(&resp.Metadata.Total)
		err := searchQ.Limit(hitLimit).Find(&resp.Metadata.Hits).Error
		if err != nil {
			log.Printf("ERROR: metadata search failed: %s", err.Error())
		}
	}

	if scope == "all" || scope == "orders" {
		searchQ := svc.DB.Table("orders")
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
		searchQ.Count(&resp.Orders.Total)
		err := searchQ.Preload("Customer").Preload("Agency").Limit(hitLimit).Find(&resp.Orders.Hits).Error
		if err != nil {
			log.Printf("ERROR: order search failed: %s", err.Error())
		}
	}

	c.JSON(http.StatusOK, resp)
}
