package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type publishedReport struct {
	Metadata []metadata `json:"records"`
	Total    int64      `json:"total"`
}

func (svc *serviceContext) getPublishedVirgo(c *gin.Context) {
	log.Printf("INFO: get items published to virgo")
	svc.getPublished(c, "date_dl_ingest is not null and external_system_id is null")
}

func (svc *serviceContext) getPublishedDPLA(c *gin.Context) {
	log.Printf("INFO: get items published to dpla")
	svc.getPublished(c, "date_dl_ingest is not null and dpla=1")
}

func (svc *serviceContext) getPublishedArchivesSpace(c *gin.Context) {
	log.Printf("INFO: get items published to archivesspace")
	svc.getPublished(c, "date_dl_ingest is not null and external_system_id=1")
}

func (svc *serviceContext) getPublishedReport(c *gin.Context) {
	reportType := c.Param("type")
	typeQuery := ""
	csvHead := make([]string, 0)
	switch reportType {
	case "virgo":
		log.Printf("INFO: generate report of all items published to virgo")
		typeQuery = "date_dl_ingest is not null and external_system_id is null"
		csvHead = append(csvHead, "PID", "Type", "Title", "Creator Name", "Barcode", "Call Number", "Catalog Key")
	case "dpla":
		log.Printf("INFO: generate report of all items published to the dpla")
		typeQuery = "date_dl_ingest is not null and dpla=1"
		csvHead = append(csvHead, "PID", "Type", "Title", "Creator Name")
	case "archivesspace":
		log.Printf("INFO: generate report of all items published to archivesspace")
		typeQuery = "date_dl_ingest is not null and external_system_id=1"
		csvHead = append(csvHead, "PID", "Type", "Title", "Creator Name")
	default:
		c.String(http.StatusBadRequest, fmt.Sprintf("unknown report type %s requested", reportType))
		return
	}

	var out []metadata
	if err := svc.DB.Where(typeQuery).Find(&out).Order("date_dl_ingest desc").Error; err != nil {
		log.Printf("ERROR: unable to geterate published report: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: %d published %s records found", len(out), reportType)

	c.Header("Content-Type", "text/csv")
	cw := csv.NewWriter(c.Writer)
	cw.Write(csvHead)
	for _, md := range out {
		line := []string{md.PID, md.Type, md.Title}
		creator := ""
		if md.CreatorName != nil {
			creator = *md.CreatorName
		}
		line = append(line, creator)
		if reportType == "virgo" {
			barcode := ""
			if md.Barcode != nil {
				barcode = *md.Barcode
			}
			callNum := ""
			if md.CallNumber != nil {
				callNum = *md.CallNumber
			}
			catKey := ""
			if md.CatalogKey != nil {
				catKey = *md.CatalogKey
			}
			line = append(line, barcode, callNum, catKey)
		}
		cw.Write(line)
	}
	cw.Flush()
}

// support pagination and filtering only
// fields to return: PID, Type, Title, Creator Name, Barcode, Call Number, Catalog Key
func (svc *serviceContext) getPublished(c *gin.Context, typeQuery string) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 25
	}

	// filter fomat: ["status|equals|active","customer|startsWith|foster","agency|equals|142"]
	filterStr := c.Query("filters")
	log.Printf("INFO: raw filters [%s]", filterStr)
	var filters []string
	err := json.Unmarshal([]byte(filterStr), &filters)
	if err != nil {
		log.Printf("ERROR: unable to parse filter payload %s: %s", filterStr, err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid filters param: %s", filterStr))
		return
	}

	out := publishedReport{}
	pubQ := svc.DB.Table("metadata").Where(typeQuery)
	for _, filter := range filters {
		bits := strings.Split(filter, "|") // target | comparison | value
		tgtField := bits[0]
		switch tgtField {
		case "creatorName":
			tgtField = "creator_name"
		case "callNumber":
			tgtField = "call_number"
		case "catalogKey":
			tgtField = "catalog_key"
		}
		comparison := bits[1]
		tgtVal, _ := url.QueryUnescape(bits[2])
		log.Printf("INFO: filter %s %s %s", tgtField, comparison, tgtVal)
		switch comparison {
		case "equals":
			pubQ = pubQ.Where(fmt.Sprintf("%s=?", tgtField), tgtVal)
		case "startsWith":
			pubQ = pubQ.Where(fmt.Sprintf("%s like ?", tgtField), fmt.Sprintf("%s%%", tgtVal))
		default:
			pubQ = pubQ.Where(fmt.Sprintf("%s like ?", tgtField), fmt.Sprintf("%%%s%%", tgtVal))
		}
	}

	err = pubQ.Count(&out.Total).Error
	if err != nil {
		log.Printf("ERROR: get published count failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = pubQ.Offset(startIndex).Limit(pageSize).Find(&out.Metadata).Order("date_dl_ingest desc").Error
	if err != nil {
		log.Printf("ERROR: published query failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}
