package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type collectionCandidate struct {
	ID               uint64 `json:"id"`
	PID              string `gorm:"column:pid" json:"pid"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	CallNumber       string `json:"callNumber"`
	Barcode          string `json:"barcode"`
	CatalogKey       string `json:"catalogKey"`
	ExternalSystemID int64  `json:"externalSystemID"`
}

func (svc *serviceContext) getCollections(c *gin.Context) {
	log.Printf("INFO: get all collections")
	type collectionHit struct {
		ID           uint64 `json:"id"`
		PID          string `gorm:"column:pid" json:"pid"`
		Type         string `json:"type"`
		Title        string `json:"title"`
		CallNumber   string `json:"callNumber"`
		Barcode      string `json:"barcode"`
		CatalogKey   string `json:"catalogKey"`
		CreatorName  string `json:"creatorName"`
		CollectionID string `gorm:"column:collection_id" json:"collectionID"`
		RecordCount  int64  `json:"recordCount"`
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

func (svc *serviceContext) getCollectionUnits(collectionID int64) ([]*unit, error) {
	out := make([]*unit, 0)

	var inCollectionIDs []uint64
	if err := svc.DB.Raw("select id from metadata where parent_metadata_id=?", collectionID).Scan(&inCollectionIDs).Error; err != nil {
		return out, err
	}
	// NOTE: Manually calculate the master files count and return it as num_master_files instead of using the inaccurate cache
	mfCnt := "(select count(*) from master_files m inner join units u on u.id=m.unit_id where u.id=units.id) as num_master_files"
	err := svc.DB.
		Where("metadata_id in ? and (intended_use_id=? or intended_use_id=?) and unit_status != ?", inCollectionIDs, 110, 101, "canceled").
		Select("units.*", mfCnt).Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("unable to get related units for %d: %s", collectionID, err.Error())
	}

	if len(out) == 0 {
		log.Printf("INFO: no units directly found for collection metadata %d; searching master files...", collectionID)
		q := "select  u.*,count(m.id) as num_master_files from master_files m "
		q += " inner join units u on u.id = m.unit_id"
		q += " where (intended_use_id=? or intended_use_id=?) and unit_status != ?"
		q += " and m.metadata_id in ? group by u.id"
		if err := svc.DB.Raw(q, 110, 101, "canceled", inCollectionIDs).Scan(&out).Error; err != nil {
			return nil, fmt.Errorf("unable to get related units for %d: %s", collectionID, err.Error())
		}
	}

	return out, nil
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
	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}
	sortField := fmt.Sprintf("metadata.%s", sortBy)
	if sortBy == "callNumber" {
		sortField = "call_number"
	}
	if sortBy == "aptStatus" {
		sortField = "APTrustSubmission.success"
	}
	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)

	qStr := c.Query("q")
	var queryClause *gorm.DB
	if qStr != "" {
		qLike := fmt.Sprintf("%%%s%%", qStr)
		queryClause = svc.DB.Where("title like ? ", qLike)
	}

	log.Printf("INFO: get collection records for collection %d, start %d limit %d, order %s, query [%s]", collectionID, startIndex, pageSize, orderStr, qStr)

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
		err = svc.DB.Joins("APTrustSubmission").Where("parent_metadata_id=?", collectionID).Where(queryClause).
			Offset(startIndex).Limit(pageSize).Order(orderStr).Find(&resp.Metadata).Error
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
		err = svc.DB.Joins("APTrustSubmission").Where("parent_metadata_id=?", collectionID).
			Offset(startIndex).Limit(pageSize).Order(orderStr).Find(&resp.Metadata).Error
		if err != nil {
			log.Printf("ERROR: unable to get collection records for collection %d: %s", collectionID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) findCollectionCandidates(c *gin.Context) {
	qStr := c.Query("q")
	qAny := fmt.Sprintf("%%%s%%", qStr)
	qStart := fmt.Sprintf("%s%%", qStr)
	log.Printf("INFO: find collection candidates matching [%s]", qStr)

	var out struct {
		Total int64                 `json:"total"`
		Hits  []collectionCandidate `json:"hits"`
	}

	searchQ := svc.DB.Debug().Table("metadata").Where("parent_metadata_id=?", 0)
	fieldQ := svc.DB.Or("title like ?", qAny).Or("barcode=?", qStr).Or("catalog_key=?", qStr).Or("call_number like ?", qStart).Or("id=?", qStr)
	searchQ.Where(fieldQ)

	searchQ.Count(&out.Total)
	err := searchQ.Limit(100).Find(&out.Hits).Error
	if err != nil {
		log.Printf("ERROR: collection candidate search failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
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

func (svc *serviceContext) addCollectionItem(c *gin.Context) {
	collectionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if collectionID == 0 {
		log.Printf("INFO: bad collection id %s in add collection item request", c.Param("id"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid collection id %s", c.Param("id")))
		return
	}
	mdID, _ := strconv.ParseInt(c.Query("rec"), 10, 64)
	if mdID == 0 {
		log.Printf("INFO: bad record id %s in add collection item request", c.Query("rec"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid rec param %s", c.Query("rec")))
		return
	}

	log.Printf("INFO: add metadata %d to collection %d", mdID, collectionID)
	var md metadata
	err := svc.DB.Find(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to get metadata %d to add to collection %d: %s", mdID, collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if md.ParentMetadataID != 0 {
		log.Printf("INFO: invalid request to add metadata %d to collection %d; it is already in collection %d", mdID, collectionID, md.ParentMetadataID)
		c.String(http.StatusBadRequest, fmt.Sprintf("this record is already part of collection %d", md.ParentMetadataID))
		return
	}

	md.ParentMetadataID = collectionID
	err = svc.DB.Model(&md).Select("ParentMetadataID").Updates(md).Error
	if err != nil {
		log.Printf("ERROR: unable to update parent_metadata_id: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("metadata has been to collection %d", collectionID))
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

func (svc *serviceContext) exportCollectionCSV(c *gin.Context) {
	collectionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if collectionID == 0 {
		log.Printf("ERROR: bad collection id %s in add collection items request", c.Param("id"))
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid collection id %s", c.Param("id")))
		return
	}
	log.Printf("INFO: export a list of pids for items in collection %d", collectionID)
	var resp []metadata
	err := svc.DB.Select("pid").Where("parent_metadata_id=?", collectionID).Find(&resp).Error
	if err != nil {
		log.Printf("ERROR: unable to get collection %d items: %s", collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	cw := csv.NewWriter(c.Writer)
	csvHead := []string{"pid"}
	cw.Write(csvHead)
	for _, md := range resp {
		line := []string{fmt.Sprintf("%s", md.PID)}
		cw.Write(line)
	}
	cw.Flush()
}
