package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type masterFileHit struct {
	ID           uint64 `json:"id"`
	PID          string `gorm:"column:pid" json:"pid"`
	UnitID       uint64 `json:"unitID"`
	Filename     string `json:"filename"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `gorm:"-" json:"thumbnailURL"`
	ImageURL     string `gorm:"-" json:"imageURL"`
}

type metadataHit struct {
	ID               uint64          `json:"id"`
	PID              string          `gorm:"column:pid" json:"pid"`
	Type             string          `json:"type"`
	Title            string          `json:"title"`
	CallNumber       string          `json:"callNumber"`
	Barcode          string          `json:"barcode"`
	CatalogKey       string          `json:"catalogKey"`
	CreatorName      string          `json:"creatorName"`
	DateDlIngest     *time.Time      `gorm:"column:date_dl_ingest" json:"-"`
	Virgo            bool            `gorm:"-" json:"virgo"`
	DPLA             bool            `json:"dpla"`
	VirgoURL         string          `gorm:"-" json:"virgoURL"`
	ExternalSystemID *int64          `json:"-"`
	ExternalSystem   *externalSystem `gorm:"foreignKey:ExternalSystemID" json:"externalSystem,omitempty"`
}

type orderHit struct {
	ID                  uint64   `json:"id"`
	OrderStatus         string   `json:"status"`
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
	Total int64          `json:"total"`
	Hits  []*metadataHit `json:"hits"`
}
type masterFileResp struct {
	Total int64            `json:"total"`
	Hits  []*masterFileHit `json:"hits"`
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

type filterData struct {
	Type   string   `json:"type"`
	Params []string `json:"params"`
}

func (svc *serviceContext) searchRequest(c *gin.Context) {
	// get the query and tag it for starts with and contains searches
	qStr := c.Query("q")
	matchStart := fmt.Sprintf("%s%%", qStr)
	matchAny := fmt.Sprintf("%%%s%%", qStr)

	// setup pagination
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 15
	}

	// when searcghing for items to be added to a collection, the request includes
	// a param collection=true. if this is the case, exclude items that are already part of a collection
	// these items will have their parent_metadata_id set greater than zero
	excludeCollectionItems, _ := strconv.ParseBool(c.Query("collection"))

	// limit search scope to an item type
	scope := c.Query("scope")
	if scope == "" {
		scope = "all"
	}
	if scope != "all" && scope != "masterfiles" && scope != "metadata" && scope != "orders" && scope != "components" {
		log.Printf("ERROR: invalid search scope %s specified", scope)
		c.String(http.StatusBadRequest, "invalid search scope")
		return
	}

	// search specific fields?
	field := c.Query("field")
	if field == "" {
		field = "all"
	}
	log.Printf("INFO: search %s.%s for [%s] starting from %d limit %d", scope, field, qStr, startIndex, pageSize)

	// extract filters into json structs
	filterStr := c.Query("filters")
	log.Printf("INFO: raw filters %s", filterStr)
	var filterQ *gorm.DB
	var filterTarget string
	if filterStr != "" {
		log.Printf("INFO: parse filters from query string")
		var filters filterData
		err := json.Unmarshal([]byte(filterStr), &filters)
		if err != nil {
			log.Printf("ERROR: unable to parse filters %s: %s", filterStr, err.Error())
			c.String(http.StatusBadRequest, "invalid filters")
			return
		}

		filterTarget = filters.Type
		for idx, f := range filters.Params {
			log.Printf("INFO: found filter %s", f)
			bits := strings.Split(f, "|")
			tgtField := bits[0]
			tgtVal, _ := url.QueryUnescape(bits[2])
			log.Printf("INFO: filter %s on %s", tgtField, tgtVal)
			if tgtField == "type" {
				typeBits := strings.Split(tgtVal, ":")
				if len(typeBits) == 1 {
					filterQ = svc.DB.Where("type=?", tgtVal)
				} else {
					filterQ = svc.DB.Where("type=? and external_system_id=?", "ExternalMetadata", typeBits[1])
				}
			} else if tgtField == "virgo" {
				if idx == 0 {
					if tgtVal == "true" {
						filterQ = svc.DB.Where("date_dl_ingest is not null")
					} else {
						filterQ = svc.DB.Where("date_dl_ingest is null")
					}
				} else {
					if tgtVal == "true" {
						filterQ = filterQ.Where("date_dl_ingest is not null")
					} else {
						filterQ = filterQ.Where("date_dl_ingest is null")
					}
				}
			} else if tgtField == "dpla" {
				if idx == 0 {
					if tgtVal == "true" {
						filterQ = svc.DB.Where("dpla=1")
					} else {
						filterQ = svc.DB.Where("dpla=0")
					}
				} else {
					if tgtVal == "true" {
						filterQ = filterQ.Where("dpla=1")
					} else {
						filterQ = filterQ.Where("dpla=0")
					}
				}
			} else {
				op := "="
				if bits[1] == "contains" {
					tgtVal = fmt.Sprintf("%%%s%%", tgtVal)
					op = "like"
				} else if bits[1] == "startsWith" {
					tgtVal = fmt.Sprintf("%s%%", tgtVal)
					op = "like"
				}
				if idx == 0 {
					filterQ = svc.DB.Where(fmt.Sprintf("%s %s ?", tgtField, op), tgtVal)
				} else {
					filterQ = filterQ.Where(fmt.Sprintf("%s %s ?", tgtField, op), tgtVal)
				}
			}
		}
	}

	// init empty results
	resp := searchResults{
		Components:  componentResp{Hits: make([]component, 0)},
		MasterFiles: masterFileResp{Hits: make([]*masterFileHit, 0)},
		Metadata:    metadataResp{Hits: make([]*metadataHit, 0)},
		Orders:      orderResp{Hits: make([]orderHit, 0)},
	}

	// query each type of object individually: components, master files, metadata and orders. Aggregate rresults in response struct
	if scope == "all" || scope == "components" {
		log.Printf("INFO: searching components for [%s]...", qStr)
		startTime := time.Now()
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

		if filterTarget == "components" && filterQ != nil {
			searchQ = searchQ.Where(filterQ)
		}

		searchQ.Count(&resp.Components.Total)
		subQ := "(select count(*) from master_files m where component_id=components.id) as mf_cnt"
		err := searchQ.Offset(startIndex).Limit(pageSize).Select("components.*", subQ).Find(&resp.Components.Hits).Error
		if err != nil {
			log.Printf("ERROR: component search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: component search found %d hits. Elapsed Time: %d (ms)", resp.Components.Total, elapsedMS)
	}

	if scope == "all" || scope == "masterfiles" {
		log.Printf("INFO: searching masterfiles for [%s]...", qStr)
		startTime := time.Now()
		searchQ := svc.DB.Table("master_files")
		if field == "all" {
			searchQ = searchQ.Where(
				svc.DB.Where("master_files.id=?", qStr).Or("pid=?", qStr).Or("unit_id=?", qStr).
					Or("filename like ?", matchAny).Or("title like ?", matchAny).
					Or("description like ?", matchAny),
			)
		} else if field == "pid" || field == "id" || field == "unit_id" {
			searchQ = searchQ.Where(fmt.Sprintf("%s=?", field), qStr)
		} else if field == "title" || field == "description" || field == "filename" {
			searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
		} else if field == "tag" {
			searchQ = searchQ.
				Joins("left outer join master_file_tags mt on mt.master_file_id = master_files.id").
				Joins("left outer join tags t on mt.tag_id = t.id").
				Where("t.tag like ?", matchAny)
		}

		if filterTarget == "masterfiles" && filterQ != nil {
			searchQ = searchQ.Where(filterQ)
		}

		searchQ.Count(&resp.MasterFiles.Total)
		err := searchQ.Offset(startIndex).Limit(pageSize).Find(&resp.MasterFiles.Hits).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
		for _, mf := range resp.MasterFiles.Hits {
			mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
			mf.ImageURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: masterfile search found %d hits. Elapsed Time: %d (ms)", resp.Components.Total, elapsedMS)
	}

	if scope == "all" || scope == "metadata" {
		log.Printf("INFO: searching metadata for [%s]...", qStr)
		startTime := time.Now()
		searchQ := svc.DB.Table("metadata")
		if field == "all" {
			searchQ = searchQ.Where(
				svc.DB.Where("metadata.id=?", qStr).Or("pid=?", qStr).Or("title like ?", matchAny).
					Or("barcode=?", qStr).Or("catalog_key=?", qStr).Or("call_number like ?", matchAny).
					Or("creator_name like ?", matchAny).Or("collection_id like ?", matchStart).Or("collection_facet like ?", matchStart),
			)
		} else {
			if field == "title" || field == "creator_name" || field == "call_number" {
				searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
			} else {
				searchQ = searchQ.Where(fmt.Sprintf("%s=?", field), qStr)
			}
		}

		if filterTarget == "metadata" && filterQ != nil {
			searchQ = searchQ.Where(filterQ)
		}
		if excludeCollectionItems {
			// only pick from records with no parent metadata id (items not in a collection)
			searchQ = searchQ.Where("parent_metadata_id=?", 0)
		}

		searchQ.Count(&resp.Metadata.Total)
		searchQ = searchQ.Preload("ExternalSystem")
		err := searchQ.Debug().Offset(startIndex).Limit(pageSize).Find(&resp.Metadata.Hits).Error
		if err != nil {
			log.Printf("ERROR: metadata search failed: %s", err.Error())
		}
		for _, md := range resp.Metadata.Hits {
			if md.DateDlIngest != nil {
				md.VirgoURL = fmt.Sprintf("%s/sources/uva_library/items/%s", svc.ExternalSystems.Virgo, md.CatalogKey)
				md.Virgo = true
			}
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: metadata search found %d hits. Elapsed Time: %d (ms)", resp.Components.Total, elapsedMS)
	}

	if (scope == "all" || scope == "orders") && field != "pid" {
		log.Printf("INFO: searching orders for [%s]...", qStr)
		startTime := time.Now()
		searchQ := svc.DB.Table("orders")
		if field == "all" {
			searchQ = searchQ.
				Joins("inner join customers on customer_id = customers.id").
				Joins("left outer join agencies on agency_id = agencies.id").
				Where(
					svc.DB.Where("orders.id=?", qStr).
						Or("customers.last_name like ?", matchStart).Or("agencies.name like ?", matchStart).
						Or("orders.staff_notes like ?", matchAny).Or("orders.special_instructions like ?", matchAny),
				)
		} else if field == "id" {
			searchQ = searchQ.Where("id=?", qStr)
		} else if field == "last_name" {
			searchQ = searchQ.Joins("inner join customers on customer_id = customers.id").
				Where("customers.last_name like ?", matchStart)
		} else if field == "agency" {
			searchQ = searchQ.Joins("left outer join agencies on agency_id = agencies.id").
				Where("agencies.name like ?", matchAny)
		} else {
			searchQ = searchQ.Where(fmt.Sprintf("%s like ?", field), matchAny)
		}

		if filterTarget == "orders" && filterQ != nil {
			searchQ = searchQ.Where(filterQ)
		}

		searchQ.Count(&resp.Orders.Total)
		err := searchQ.Preload("Customer").Preload("Agency").
			Offset(startIndex).Limit(pageSize).
			Find(&resp.Orders.Hits).Error
		if err != nil {
			log.Printf("ERROR: order search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: orders search found %d hits. Elapsed Time: %d (ms)", resp.Components.Total, elapsedMS)
	}

	c.JSON(http.StatusOK, resp)
}
