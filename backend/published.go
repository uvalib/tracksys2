package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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

	type resp struct {
		Metadata []metadata `json:"records"`
		Total    int64      `json:"total"`
	}
	out := resp{}

	pubQ := svc.DB.Table("metadata").Where(typeQuery)
	for _, filter := range filters {
		bits := strings.Split(filter, "|") // target | comparison | value
		tgtField := bits[0]
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
