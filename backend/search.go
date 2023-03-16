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
	ID           uint64   `json:"id"`
	PID          string   `gorm:"column:pid" json:"pid"`
	UnitID       uint64   `json:"unitID"`
	Filename     string   `json:"filename"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ThumbnailURL string   `gorm:"-" json:"thumbnailURL"`
	ImageURL     string   `gorm:"-" json:"imageURL"`
	MetadataID   int64    `json:"-"`
	Metadata     metadata `gorm:"foreignKey:MetadataID" json:"metadata"`
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

type searchFilter struct {
	Target string
	Query  *gorm.DB
}

type searchContext struct {
	Query                  string
	QueryAny               string
	QueryStart             string
	QueryType              string // general, id or pid
	IntQuery               int64
	Scope                  string
	Field                  string
	Filter                 *searchFilter
	StartIndex             int
	PageSize               int
	ExcludeCollectionItems bool
}

func (svc *serviceContext) searchRequest(c *gin.Context) {
	// setup the search context, starting with the query param
	q := strings.TrimSpace(c.Query("q"))
	sc := searchContext{QueryType: "general", Query: q, QueryAny: fmt.Sprintf("%%%s%%", q), QueryStart: fmt.Sprintf("%s%%", q)}
	sc.IntQuery, _ = strconv.ParseInt(sc.Query, 10, 64)
	if sc.IntQuery > 0 {
		log.Printf("INFO: query %s appears to be an id; use specialized id handling during search", sc.Query)
		sc.QueryType = "id"
	} else {
		if strings.Index(q, "tsm:") == 0 || strings.Index(q, "tsb:") == 0 || strings.Index(q, "uva-lib:") == 0 {
			log.Printf("INFO: query %s appears to be a pid; just search on pid columns", q)
			sc.QueryType = "pid"
		}
	}

	// setup pagination
	sc.StartIndex, _ = strconv.Atoi(c.Query("start"))
	sc.PageSize, _ = strconv.Atoi(c.Query("limit"))
	if sc.PageSize == 0 {
		sc.PageSize = 15
	}

	// when searcghing for items to be added to a collection, the request includes
	// a param collection=true. if this is the case, exclude items that are already part of a collection
	// these items will have their parent_metadata_id set greater than zero
	sc.ExcludeCollectionItems, _ = strconv.ParseBool(c.Query("collection"))

	// limit search scope to an item type
	sc.Scope = c.Query("scope")
	if sc.Scope == "" {
		sc.Scope = "all"
	}
	if sc.Scope != "all" && sc.Scope != "masterfiles" && sc.Scope != "metadata" && sc.Scope != "orders" && sc.Scope != "components" {
		log.Printf("ERROR: invalid search scope %s specified", sc.Scope)
		c.String(http.StatusBadRequest, "invalid search scope")
		return
	}

	// search specific fields?
	sc.Field = c.Query("field")
	if sc.Field == "" {
		sc.Field = "all"
	}
	log.Printf("INFO: search %s.%s for [%s] starting from %d limit %d", sc.Scope, sc.Field, sc.Query, sc.StartIndex, sc.PageSize)

	// extract filterdata into a db query that can be appended later
	filter, err := svc.initFilter(c.Query("filters"))
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// init empty results
	resp := searchResults{
		Components:  componentResp{Hits: make([]component, 0)},
		MasterFiles: masterFileResp{Hits: make([]*masterFileHit, 0)},
		Metadata:    metadataResp{Hits: make([]*metadataHit, 0)},
		Orders:      orderResp{Hits: make([]orderHit, 0)},
	}

	// query each type of object individually: components, master files, metadata and orders. Aggregate rresults in response struct
	if sc.Scope == "all" || sc.Scope == "components" {
		log.Printf("INFO: searching components for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Debug().Table("components")
		if sc.QueryType != "pid" {
			if filter.Target == "components" {
				searchQ = searchQ.Where(filter.Query)
			}

			var fieldQ *gorm.DB
			if sc.Field == "all" {
				fieldQ = svc.DB.Or("title like ?", sc.QueryAny).Or("label like ?", sc.QueryAny).
					Or("content_desc like ?", sc.QueryAny).Or("date like ?", sc.QueryAny).Or("ead_id_att=?", sc.Query)
				if sc.QueryType == "id" {
					fieldQ = fieldQ.Or("components.id=?", sc.IntQuery)
				}
			} else if sc.Field == "ead_id_att" {
				fieldQ = svc.DB.Where("ead_id_att=?", sc.Query)
			} else {
				fieldQ = svc.DB.Where(fmt.Sprintf("%s like ?", sc.Field), sc.QueryAny)
			}
			searchQ.Where(fieldQ)
		} else {
			searchQ.Where("pid=?", sc.Query)
		}

		searchQ.Count(&resp.Components.Total)
		subQ := "(select count(*) from master_files m where component_id=components.id) as mf_cnt"
		err := searchQ.Offset(sc.StartIndex).Limit(sc.PageSize).Select("components.*", subQ).Find(&resp.Components.Hits).Error
		if err != nil {
			log.Printf("ERROR: component search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: component search found %d hits. Elapsed Time: %d (ms)", resp.Components.Total, elapsedMS)
	}

	if sc.Scope == "all" || sc.Scope == "masterfiles" {
		log.Printf("INFO: searching masterfiles for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Debug().Table("master_files").Joins("inner join metadata md on md.id=metadata_id")
		if sc.QueryType != "pid" {
			if filter.Target == "masterfiles" {
				searchQ = searchQ.Where(filter.Query)
			}

			var fieldQ *gorm.DB
			if sc.Field == "all" {
				if sc.QueryType == "id" {
					fieldQ = svc.DB.Or("master_files.id=?", sc.IntQuery).Or("unit_id=?", sc.IntQuery)
				} else {
					fieldQ = svc.DB.Or("filename like ?", sc.QueryAny).Or("master_files.title like ?", sc.QueryAny).
						Or("description like ?", sc.QueryAny).Or("md.call_number = ?", sc.Query)
				}
			} else if sc.Field == "unit_id" {
				fieldQ = svc.DB.Where("unit_id=?", sc.IntQuery)
			} else if sc.Field == "call_number" {
				fieldQ = svc.DB.Where("call_number=?", sc.Query)
			} else if sc.Field == "title" || sc.Field == "description" || sc.Field == "filename" {
				fieldQ = svc.DB.Where(fmt.Sprintf("master_files.%s like ?", sc.Field), sc.QueryAny)
			} else if sc.Field == "tag" {
				searchQ = searchQ.
					Joins("left outer join master_file_tags mt on mt.master_file_id = master_files.id").
					Joins("left outer join tags t on mt.tag_id = t.id")
				fieldQ = svc.DB.Where("t.tag like ?", sc.QueryAny)
			}
			searchQ.Where(fieldQ)
		} else {
			searchQ.Where("pid=?", sc.Query)
		}

		searchQ.Count(&resp.MasterFiles.Total)
		err := searchQ.Preload("Metadata").Offset(sc.StartIndex).Limit(sc.PageSize).Find(&resp.MasterFiles.Hits).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
		for _, mf := range resp.MasterFiles.Hits {
			mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
			mf.ImageURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: masterfile search found %d hits. Elapsed Time: %d (ms)", resp.MasterFiles.Total, elapsedMS)
	}

	if sc.Scope == "all" || sc.Scope == "metadata" {
		log.Printf("INFO: searching metadata for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Debug().Table("metadata")
		if sc.QueryType != "pid" {
			if filter.Target == "metadata" {
				searchQ = searchQ.Where(filter.Query)
			}

			var fieldQ *gorm.DB
			if sc.Field == "all" {
				fieldQ = svc.DB.Or("title like ?", sc.QueryAny).
					Or("barcode=?", sc.Query).Or("catalog_key=?", sc.Query).Or("call_number like ?", sc.QueryAny).
					Or("creator_name like ?", sc.QueryAny).Or("collection_id like ?", sc.QueryStart).Or("collection_facet like ?", sc.QueryStart)
				if sc.QueryType == "id" {
					fieldQ = fieldQ.Or("metadata.id=?", sc.IntQuery)
				}
			} else {
				if sc.Field == "title" || sc.Field == "creator_name" || sc.Field == "call_number" {
					fieldQ = svc.DB.Where(fmt.Sprintf("%s like ?", sc.Field), sc.QueryAny)
				} else {
					fieldQ = svc.DB.Where(fmt.Sprintf("%s=?", sc.Field), sc.Query)
				}
			}
			searchQ.Where(fieldQ)
		} else {
			searchQ.Where("pid=?", sc.Query)
		}

		if sc.ExcludeCollectionItems {
			// only pick from records with no parent metadata id (items not in a collection)
			searchQ = searchQ.Where("parent_metadata_id=?", 0)
		}

		searchQ.Count(&resp.Metadata.Total)
		searchQ = searchQ.Preload("ExternalSystem")
		err := searchQ.Debug().Offset(sc.StartIndex).Limit(sc.PageSize).Find(&resp.Metadata.Hits).Error
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
		log.Printf("INFO: metadata search found %d hits. Elapsed Time: %d (ms)", resp.Metadata.Total, elapsedMS)
	}

	if (sc.Scope == "all" || sc.Scope == "orders") && sc.QueryType != "pid" {
		log.Printf("INFO: searching orders for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Debug().Table("orders").
			Joins("inner join customers on customer_id = customers.id").
			Joins("left outer join agencies on agency_id = agencies.id")
		if filter.Target == "orders" {
			searchQ = searchQ.Where(filter.Query)
		}

		var fieldQ *gorm.DB
		if sc.Field == "all" {
			if sc.QueryType == "id" {
				fieldQ = svc.DB.Where("orders.id=?", sc.IntQuery)
			} else {
				fieldQ = svc.DB.Or("customers.last_name like ?", sc.QueryStart).Or("agencies.name like ?", sc.QueryStart).
					Or("orders.staff_notes like ?", sc.QueryAny).Or("orders.special_instructions like ?", sc.QueryAny)
			}
		} else if sc.Field == "id" {
			fieldQ = svc.DB.Where("orders.id=?", sc.IntQuery)
		} else if sc.Field == "last_name" {
			fieldQ = svc.DB.Where("customers.last_name like ?", sc.QueryStart)
		} else if sc.Field == "agency" {
			fieldQ = svc.DB.Where("agencies.name like ?", sc.QueryAny)
		} else {
			fieldQ = svc.DB.Where(fmt.Sprintf("%s like ?", sc.Field), sc.QueryAny)
		}
		searchQ.Where(fieldQ)

		searchQ.Count(&resp.Orders.Total)
		err := searchQ.Preload("Customer").Preload("Agency").
			Offset(sc.StartIndex).Limit(sc.PageSize).
			Find(&resp.Orders.Hits).Error
		if err != nil {
			log.Printf("ERROR: order search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: orders search found %d hits. Elapsed Time: %d (ms)", resp.Orders.Total, elapsedMS)
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) initFilter(filterStr string) (*searchFilter, error) {
	log.Printf("INFO: raw filters [%s]", filterStr)
	out := searchFilter{Target: "none"}

	if filterStr != "" {
		log.Printf("INFO: parse filters from query string")
		var filters filterData
		err := json.Unmarshal([]byte(filterStr), &filters)
		if err != nil {
			return nil, fmt.Errorf("unable to parse filter %s: %s", filterStr, err.Error())
		}

		out.Target = filters.Type
		for idx, f := range filters.Params {
			log.Printf("INFO: found filter %s", f)
			bits := strings.Split(f, "|")
			tgtField := bits[0]
			tgtVal, _ := url.QueryUnescape(bits[2])
			log.Printf("INFO: filter %s on %s", tgtField, tgtVal)

			if tgtField == "type" {
				typeBits := strings.Split(tgtVal, ":")
				if len(typeBits) == 1 {
					out.Query = svc.DB.Where("type=?", tgtVal)
				} else {
					out.Query = svc.DB.Where("type=? and external_system_id=?", "ExternalMetadata", typeBits[1])
				}
			} else if tgtField == "virgo" {
				if idx == 0 {
					if tgtVal == "true" {
						out.Query = svc.DB.Where("date_dl_ingest is not null")
					} else {
						out.Query = svc.DB.Where("date_dl_ingest is null")
					}
				} else {
					if tgtVal == "true" {
						out.Query = out.Query.Where("date_dl_ingest is not null")
					} else {
						out.Query = out.Query.Where("date_dl_ingest is null")
					}
				}
			} else if tgtField == "dpla" {
				if idx == 0 {
					if tgtVal == "true" {
						out.Query = svc.DB.Where("dpla=1")
					} else {
						out.Query = svc.DB.Where("dpla=0")
					}
				} else {
					if tgtVal == "true" {
						out.Query = out.Query.Where("dpla=1")
					} else {
						out.Query = out.Query.Where("dpla=0")
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
					out.Query = svc.DB.Where(fmt.Sprintf("%s %s ?", tgtField, op), tgtVal)
				} else {
					out.Query = out.Query.Where(fmt.Sprintf("%s %s ?", tgtField, op), tgtVal)
				}
			}
		}
	}
	return &out, nil
}
