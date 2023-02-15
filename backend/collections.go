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

func (svc *serviceContext) getCollections(c *gin.Context) {
	log.Printf("INFO: get all collections")
	type collectionHit struct {
		ID          uint64 `json:"id"`
		PID         string `gorm:"column:pid" json:"pid"`
		Type        string `json:"type"`
		Title       string `json:"title"`
		CallNumber  string `json:"callNumber"`
		Barcode     string `json:"barcode"`
		CatalogKey  string `json:"catalogKey"`
		CreatorName string `json:"creatorName"`
		RecordCount int64  `json:"recordCount"`
	}
	var resp struct {
		Total       int64           `json:"total"`
		Collections []collectionHit `json:"collections"`
	}

	log.Printf("INFO: get collections count")
	err := svc.DB.Table("metadata").Where("is_collection=?", true).Count(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get collections count %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: get collections details")
	err = svc.DB.Debug().Table("metadata").
		Joins("left join metadata mc on mc.parent_metadata_id = metadata.id").
		Select("metadata.*", "count(mc.id) as record_count").
		Where("metadata.is_collection=?", true).Group("metadata.id").Find(&resp.Collections).Error
	if err != nil {
		log.Printf("ERROR: unable to get collections %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getCollectionItems(c *gin.Context) {
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

func (svc *serviceContext) removeCollectionItem(c *gin.Context) {
	collectionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if collectionID == 0 {
		log.Printf("ERROR: bad collection id %s in remove collection item request", c.Param("id"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid collection id %s", c.Param("id")))
		return
	}
	itemID, _ := strconv.ParseInt(c.Param("item"), 10, 64)
	if itemID == 0 {
		log.Printf("ERROR: bad item id %s in remove collection item request", c.Param("item"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid collection item id %s", c.Param("item")))
		return
	}
	var tgtItem metadata
	err := svc.DB.Find(&tgtItem, itemID).Error
	if err != nil {
		log.Printf("ERROR: unable to load collection item %d: %s", itemID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if tgtItem.ParentMetadataID != collectionID {
		log.Printf("ERROR: item %d is not part of collection %d", itemID, collectionID)
		c.String(http.StatusBadRequest, fmt.Sprintf("item %d is not part of collection %d", itemID, collectionID))
		return
	}

	now := time.Now()
	tgtItem.ParentMetadataID = 0
	tgtItem.UpdatedAt = &now
	err = svc.DB.Model(&tgtItem).Select("ParentMetadataID", "UpdatedAt").Updates(tgtItem).Error
	if err != nil {
		log.Printf("ERROR: unable to remove %d from collection %d: %s", itemID, collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "removed")
}

func (svc *serviceContext) addCollectionItems(c *gin.Context) {
	collectionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if collectionID == 0 {
		log.Printf("ERROR: bad collection id %s in add collection items request", c.Param("id"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid collection id %s", c.Param("id")))
		return
	}
	var req struct {
		Items []int64 `json:"items"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid add collection %d items request: %s", collectionID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// NOTE: the collection lookup logic excludes items that are already part of a collection
	// so it is ok to add any ites that are included in the request
	log.Printf("INFO: add items to collection %d: %+v", collectionID, req.Items)

	// update parent_metadata_id of all metadata records in the request list to be collectionID
	err = svc.DB.Table("metadata").Where("id in ?", req.Items).Updates(metadata{ParentMetadataID: collectionID}).Error
	if err != nil {
		log.Printf("ERROR: unable to update collection %d items: %s", collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d items added to collection %d", len(req.Items), collectionID))
}

func (svc *serviceContext) addCollectionFacet(c *gin.Context) {
	var req struct {
		Facet string `json:"facet"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid add collection facet request: %s", err.Error())
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
