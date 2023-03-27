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
	NumMasterFiles   uint            `json:"masterFilesCount,omitempty"`
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
	Query      string
	QueryAny   string
	QueryStart string
	QueryType  string // general, id or pid
	IntQuery   int64
	Scope      string
	Field      string
	Filter     *searchFilter
	StartIndex int
	PageSize   int
	Collection bool
}

type searchChannel struct {
	Type    string
	Results interface{}
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
	sc.Collection, _ = strconv.ParseBool(c.Query("collection"))

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

	// extract filter data into a db query that can be appended later
	var err error
	sc.Filter, err = svc.initFilter(c.Query("filters"))
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// query each type of object individually: components, master files, metadata and orders
	log.Printf("INFO: issue search requests...")
	pendingCount := 4
	channel := make(chan searchChannel)
	startTime := time.Now()
	go svc.queryMasterFiles(&sc, channel)
	go svc.queryComponents(&sc, channel)
	go svc.queryMetadata(&sc, channel)
	go svc.queryOrders(&sc, channel)

	log.Printf("INFO: await all searchs responses...")
	resp := searchResults{}
	for pendingCount > 0 {
		searchResp := <-channel
		pendingCount--
		if searchResp.Type == "masterFiles" {
			log.Printf("INFO received master files search response")
			mfResp, ok := searchResp.Results.(masterFileResp)
			if ok {
				resp.MasterFiles = mfResp
			}
		} else if searchResp.Type == "metadata" {
			log.Printf("INFO received metadata search response")
			mResp, ok := searchResp.Results.(metadataResp)
			if ok {
				resp.Metadata = mResp
			}
		} else if searchResp.Type == "components" {
			log.Printf("INFO received components search response")
			cResp, ok := searchResp.Results.(componentResp)
			if ok {
				resp.Components = cResp
			}
		} else if searchResp.Type == "orders" {
			log.Printf("INFO received orders search response")
			oResp, ok := searchResp.Results.(orderResp)
			if ok {
				resp.Orders = oResp
			}
		}
	}
	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)
	log.Printf("INFO: all responses received. Elapsed Time: %d (ms)", elapsedMS)

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

func (svc *serviceContext) queryMasterFiles(sc *searchContext, channel chan searchChannel) {
	resp := masterFileResp{Hits: make([]*masterFileHit, 0)}
	if sc.Scope == "all" || sc.Scope == "masterfiles" {
		log.Printf("INFO: searching masterfiles for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Table("master_files").Joins("inner join metadata md on md.id=metadata_id")
		if sc.QueryType != "pid" {
			if sc.Filter.Target == "masterfiles" {
				searchQ = searchQ.Where(sc.Filter.Query)
			}

			var fieldQ *gorm.DB
			if sc.Field == "all" {
				if sc.QueryType == "id" {
					fieldQ = svc.DB.Or("master_files.id=?", sc.IntQuery).Or("unit_id=?", sc.IntQuery)
				} else {
					fieldQ = svc.DB.Or("md.call_number like ?", sc.QueryStart).
						Or("master_files.title like ?", sc.QueryAny).
						Or("description like ?", sc.QueryAny)
				}
			} else if sc.Field == "unit_id" {
				fieldQ = svc.DB.Where("unit_id=?", sc.IntQuery)
			} else if sc.Field == "call_number" {
				fieldQ = svc.DB.Where("call_number like ?", sc.QueryStart)
			} else if sc.Field == "title" || sc.Field == "description" {
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

		searchQ.Count(&resp.Total)
		err := searchQ.Preload("Metadata").Offset(sc.StartIndex).Limit(sc.PageSize).Find(&resp.Hits).Error
		if err != nil {
			log.Printf("ERROR: masterfile search failed: %s", err.Error())
		}
		for _, mf := range resp.Hits {
			mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
			mf.ImageURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: masterfile search found %d hits. Elapsed Time: %d (ms)", resp.Total, elapsedMS)
	}
	channel <- searchChannel{Type: "masterFiles", Results: resp}
}

func (svc *serviceContext) queryMetadata(sc *searchContext, channel chan searchChannel) {
	resp := metadataResp{Hits: make([]*metadataHit, 0)}
	if sc.Scope == "all" || sc.Scope == "metadata" {
		log.Printf("INFO: searching metadata for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Table("metadata")
		if sc.QueryType != "pid" {
			if sc.Filter.Target == "metadata" {
				searchQ = searchQ.Where(sc.Filter.Query)
			}

			if sc.Collection == true {
				// this is a special case search for metadata records that are candatates for inclusion in
				// a collection. Don't include external metadata and only search a few fields
				log.Printf("INFO: this is a search for internal metadata records that are candidates for being part of a collection")
				mfCnt := "(select count(*) from master_files m inner join units mu on mu.id = m.unit_id where mu.metadata_id=metadata.id and mu.intended_use_id=110) as num_master_files"
				searchQ = searchQ.Debug().Joins("inner join units u on u.metadata_id = metadata.id").Where("type != ? and u.intended_use_id=?", "ExternalMetadata", 110)
				fieldQ := svc.DB.Or("title like ?", sc.QueryAny).Or("barcode like ?", sc.QueryStart).Or("catalog_key like ?", sc.QueryStart).
					Or("call_number like ?", sc.QueryStart).Or("u.id=?", sc.IntQuery)
				searchQ.Where("parent_metadata_id=?", 0).Where(fieldQ).Group("metadata.id").Select("metadata.*", mfCnt)

			} else {
				var fieldQ *gorm.DB
				if sc.Field == "all" {
					fieldQ = svc.DB.Or("title like ?", sc.QueryAny).
						Or("barcode=?", sc.Query).Or("catalog_key=?", sc.Query).Or("call_number like ?", sc.QueryStart).
						Or("creator_name like ?", sc.QueryStart).Or("collection_id like ?", sc.QueryStart).Or("collection_facet like ?", sc.QueryStart)
					if sc.QueryType == "id" {
						fieldQ = fieldQ.Or("metadata.id=?", sc.IntQuery)
					}
				} else {
					if sc.Field == "title" || sc.Field == "creator_name" {
						fieldQ = svc.DB.Where(fmt.Sprintf("%s like ?", sc.Field), sc.QueryAny)
					} else if sc.Field == "call_number" {
						fieldQ = svc.DB.Where("call_number like ?", sc.QueryStart)
					} else {
						fieldQ = svc.DB.Where(fmt.Sprintf("%s=?", sc.Field), sc.Query)
					}
				}
				searchQ.Where(fieldQ)
			}
		} else {
			searchQ.Where("pid=?", sc.Query)
		}

		var err error
		searchQ.Count(&resp.Total)
		if sc.Collection == false {
			searchQ = searchQ.Preload("ExternalSystem")
			err = searchQ.Offset(sc.StartIndex).Limit(sc.PageSize).Find(&resp.Hits).Error
		} else {
			// this is a search foe metdata records to include in a collection. get them all as this query
			// is targetd and reaults will be small
			err = searchQ.Limit(500).Find(&resp.Hits).Error
		}
		if err != nil {
			log.Printf("ERROR: metadata search failed: %s", err.Error())
		}
		for _, md := range resp.Hits {
			if md.DateDlIngest != nil {
				md.VirgoURL = fmt.Sprintf("%s/sources/uva_library/items/%s", svc.ExternalSystems.Virgo, md.CatalogKey)
				md.Virgo = true
			}
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: metadata search found %d hits. Elapsed Time: %d (ms)", resp.Total, elapsedMS)
	}
	channel <- searchChannel{Type: "metadata", Results: resp}
}

func (svc *serviceContext) queryOrders(sc *searchContext, channel chan searchChannel) {
	resp := orderResp{Hits: make([]orderHit, 0)}
	if (sc.Scope == "all" || sc.Scope == "orders") && sc.QueryType != "pid" {
		log.Printf("INFO: searching orders for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Table("orders").
			Joins("inner join customers on customer_id = customers.id").
			Joins("left outer join agencies on agency_id = agencies.id")
		if sc.Filter.Target == "orders" {
			searchQ = searchQ.Where(sc.Filter.Query)
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

		searchQ.Count(&resp.Total)
		err := searchQ.Preload("Customer").Preload("Agency").
			Offset(sc.StartIndex).Limit(sc.PageSize).
			Find(&resp.Hits).Error
		if err != nil {
			log.Printf("ERROR: order search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: orders search found %d hits. Elapsed Time: %d (ms)", resp.Total, elapsedMS)
	}
	channel <- searchChannel{Type: "orders", Results: resp}
}

func (svc *serviceContext) queryComponents(sc *searchContext, channel chan searchChannel) {
	resp := componentResp{Hits: make([]component, 0)}
	if sc.Scope == "all" || sc.Scope == "components" {
		log.Printf("INFO: searching components for [%s]...", sc.Query)
		startTime := time.Now()
		searchQ := svc.DB.Table("components")
		if sc.QueryType != "pid" {
			if sc.Filter.Target == "components" {
				searchQ = searchQ.Where(sc.Filter.Query)
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

		searchQ.Count(&resp.Total)
		subQ := "(select count(*) from master_files m where component_id=components.id) as mf_cnt"
		err := searchQ.Offset(sc.StartIndex).Limit(sc.PageSize).Select("components.*", subQ).Find(&resp.Hits).Error
		if err != nil {
			log.Printf("ERROR: component search failed: %s", err.Error())
		}
		elapsedNanoSec := time.Since(startTime)
		elapsedMS := int64(elapsedNanoSec / time.Millisecond)
		log.Printf("INFO: component search found %d hits. Elapsed Time: %d (ms)", resp.Total, elapsedMS)
	}
	channel <- searchChannel{Type: "components", Results: resp}
}
