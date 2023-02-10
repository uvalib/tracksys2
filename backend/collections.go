package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (svc *serviceContext) getCollectionRecords(c *gin.Context) {
	collectionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if collectionID == 0 {
		log.Printf("ERROR: bad collection id %s in get collection records request", c.Param("id"))
		c.String(http.StatusBadRequest, "invalid collection id")
		return
	}
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	qStr := c.Query("q")
	var queryClause *gorm.DB
	if qStr != "" {
		qLike := fmt.Sprintf("%%%s%%", qStr)
		queryClause = svc.DB.Where("title like ? ", qLike)
	}

	log.Printf("INFO: get collection records for collection %d, start %d limit %d, query [%s]", collectionID, startIndex, pageSize, qStr)

	var resp struct {
		Metadata []metadata `json:"records"`
		Total    int64      `json:"total"`
	}
	if queryClause != nil {
		err := svc.DB.Table("metadata").Where("parent_metadata_id=?", collectionID).Where(queryClause).Count(&resp.Total).Error
		if err != nil {
			log.Printf("ERROR: unable to get filtered collection records count for collection %d: %s", collectionID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		err = svc.DB.Where("parent_metadata_id=?", collectionID).Where(queryClause).Offset(startIndex).Limit(pageSize).Find(&resp.Metadata).Error
		if err != nil {
			log.Printf("ERROR: unable to get filtered collection records for collection %d: %s", collectionID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		err := svc.DB.Table("metadata").Where("parent_metadata_id=?", collectionID).Count(&resp.Total).Error
		if err != nil {
			log.Printf("ERROR: unable to get collection records count for collection %d: %s", collectionID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		err = svc.DB.Where("parent_metadata_id=?", collectionID).Offset(startIndex).Limit(pageSize).Find(&resp.Metadata).Error
		if err != nil {
			log.Printf("ERROR: unable to get collection records for collection %d: %s", collectionID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) addCollectionFacet(c *gin.Context) {
	var req struct {
		Facet string `json:"facet"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid sad collection facet request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: add new facet %s", req.Facet)
	newFacet := collectionFacet{Name: req.Facet, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = svc.DB.Create(&newFacet).Error
	if err != nil {
		log.Printf("ERROR: unable to create collection faccet %s: %s", req.Facet, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load collections facets after creating new facet %s", req.Facet)
	var out []collectionFacet
	err = svc.DB.Order("name asc").Find(&out).Error
	if err != nil {
		log.Printf("ERROR: unable to get updated facets: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, out)
}
