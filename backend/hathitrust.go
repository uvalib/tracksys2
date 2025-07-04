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

type hathitrustStatus struct {
	ID                  uint       `json:"id"`
	MetadataID          int64      `json:"metadataID"`
	Metadata            *metadata  `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	RequestedAt         time.Time  `json:"requestedAt"`
	PackageCreatedAt    *time.Time `json:"packageCreatedAt"`
	PackageSubmittedAt  *time.Time `json:"packageSubmittedAt"`
	PackageStatus       string     `json:"packageStatus"`
	MetadataSubmittedAt *time.Time `json:"metadataSubmittedAt"`
	MetadataStatus      string     `json:"metadataStatus"`
	FinishedAt          *time.Time `json:"finishedAt"`
	Notes               string     `json:"notes"`
}

type hatiTrustUpdateRequest struct {
	RequestedAt         string `json:"requestedAt"`
	PackageCreatedAt    string `json:"packageCreatedAt"`
	PackageSubmittedAt  string `json:"packageSubmittedAt"`
	PackageStatus       string `json:"packageStatus"`
	MetadataSubmittedAt string `json:"metadataSubmittedAt"`
	MetadataStatus      string `json:"metadataStatus"`
	FinishedAt          string `json:"finishedAt"`
	Notes               string `json:"notes"`
}

type hathiTrustBatchUpdateRequest struct {
	OderID    int64   `json:"orderID"`
	StatusIDs []int64 `json:"statusIDs"`
	Field     string  `json:"field"`
	Value     string  `json:"value"`
}

type hathiTrustSubmissionsResonse struct {
	Total       int64              `json:"total"`
	Submissions []hathitrustStatus `json:"submissions"`
}

func (svc *serviceContext) getHathiTrustSubmissions(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	queryStr := c.Query("q")

	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "pid"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sortField := fmt.Sprintf("hathitrust_statuses.%s", sortBy)
	switch sortBy {
	case "pid":
		sortField = "Metadata.pid"
	case "title":
		sortField = "Metadata.title"
	case "barcode":
		sortField = "Metadata.barcode"
	}

	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	log.Printf("INFO: get %d hathitrust submissions starting from offset %d order %s; query=[%s]", pageSize, startIndex, orderStr, queryStr)

	var filterQ *gorm.DB
	filterStr := c.Query("filters")
	if filterStr != "" {
		log.Printf("INFO: raw filters [%s]", filterStr)
		var filters []string
		err := json.Unmarshal([]byte(filterStr), &filters)
		if err != nil {
			log.Printf("ERROR: unable to parse filter payload %s: %s", filterStr, err.Error())
			c.String(http.StatusBadRequest, fmt.Sprintf("invalid filters param: %s", filterStr))
			return
		}
		for _, filter := range filters {
			bits := strings.Split(filter, "|") // target | comparison | value
			tgtField := bits[0]
			tgtVal, _ := url.QueryUnescape(bits[2])
			log.Printf("INFO: filter %s=%s", tgtField, tgtVal)
			switch tgtField {
			case "metadataStatus":
				tgtField = "metadata_status"
			case "packageStatus":
				tgtField = "package_status"
			default:
				continue
			}
			if filterQ == nil {
				filterQ = svc.DB.Where(fmt.Sprintf("%s=?", tgtField), tgtVal)
			} else {
				filterQ = filterQ.Where(fmt.Sprintf("%s=?", tgtField), tgtVal)
			}
		}
	}

	resp := hathiTrustSubmissionsResonse{}
	err := svc.DB.Table("hathitrust_statuses").Count(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get hathi submissions count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	searchQ := svc.DB.Joins("Metadata").Order(orderStr).Offset(startIndex).Limit(pageSize).
		Where("title like ? or barcode like ?", fmt.Sprintf("%%%s%%", queryStr), fmt.Sprintf("%s%%", queryStr))
	if filterQ != nil {
		searchQ = searchQ.Where(filterQ)
	}
	err = searchQ.Debug().Find(&resp.Submissions).Error
	if err != nil {
		log.Printf("ERROR: unable to get hathi submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) updateHathiTrustSubmissions(c *gin.Context) {
	log.Printf("INFO: received batch hathitrust update request")

	var req hathiTrustBatchUpdateRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid batch update hathitrust request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(req.StatusIDs) == 0 && req.OderID == 0 {
		log.Printf("INFO: batch update is missing required id list or order id")
		c.String(http.StatusBadRequest, "either order id or and id list is required")
		return
	}

	statusIDs := req.StatusIDs
	if req.OderID > 0 {
		log.Printf("INFO: batch update hathitrust %s=%s for order %d", req.Field, req.Value, req.OderID)
		idQ := "select h.id from orders o inner join units u on u.order_id = o.id inner join metadata m on m.id = u.metadata_id "
		idQ += "inner join hathitrust_statuses h on h.metadata_id = m.id where o.id=? and m.hathitrust = 1"
		err = svc.DB.Raw(idQ, req.OderID).Scan(&statusIDs).Error
		if err != nil {
			log.Printf("ERROR: unable to get hathitrust status ids for update: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = svc.DB.Exec(fmt.Sprintf("update hathitrust_statuses set %s=? where id in ?", req.Field), req.Value, statusIDs).Error
	if err != nil {
		log.Printf("ERROR: unable to update hathitrust status hathitrust %s=%s: %s", req.Field, req.Value, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "updated")
}

func (svc *serviceContext) updateHathiTrustStatus(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received hathitrust update request for metadata %s", mdID)

	var req hatiTrustUpdateRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid update hathitrust request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var md metadata
	err = svc.DB.Preload("HathiTrustStatus").First(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to get metadata %s: %s ", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if md.HathiTrustStatus == nil {
		log.Printf("ERROR: metadata %d does not have a hathtrust status record", md.ID)
		c.String(http.StatusBadRequest, "HathiTrust status not found")
		return
	}

	md.HathiTrustStatus.MetadataStatus = req.MetadataStatus
	md.HathiTrustStatus.PackageStatus = req.PackageStatus
	md.HathiTrustStatus.Notes = req.Notes
	updates := []string{"MetadataStatus", "PackageStatus", "Notes"}

	if req.PackageCreatedAt != "" {
		pCreateDate, err := parseDateString(req.PackageCreatedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid package create date", req.PackageCreatedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid package create date", req.PackageCreatedAt))
			return
		}
		md.HathiTrustStatus.PackageCreatedAt = &pCreateDate
		updates = append(updates, "PackageCreatedAt")
	}

	if req.PackageSubmittedAt != "" {
		pSubmitDate, err := parseDateString(req.PackageSubmittedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid package submit date", req.PackageSubmittedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid package submit date", req.PackageSubmittedAt))
			return
		}
		md.HathiTrustStatus.PackageSubmittedAt = &pSubmitDate
		updates = append(updates, "PackageSubmittedAt")
	}

	if req.FinishedAt != "" {
		finishDate, err := parseDateString(req.FinishedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid finish date", req.FinishedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid finish date", req.FinishedAt))
			return
		}
		md.HathiTrustStatus.FinishedAt = &finishDate
		updates = append(updates, "FinishedAt")
	}

	if req.MetadataSubmittedAt != "" {
		mSubmitDate, err := parseDateString(req.MetadataSubmittedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid metadata submit date", req.MetadataSubmittedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid metadata submit date", req.MetadataSubmittedAt))
			return
		}
		md.HathiTrustStatus.MetadataSubmittedAt = &mSubmitDate
		updates = append(updates, "MetadataSubmittedAt")
	}

	err = svc.DB.Model(md.HathiTrustStatus).Select(updates).Updates(md.HathiTrustStatus).Error
	if err != nil {
		log.Printf("ERROR: hathiutrust status update for metadata %d failed: %s", md.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, md.HathiTrustStatus)
}

func (svc *serviceContext) flagMetadataForHathiTrust(mdID int64) error {
	log.Printf("INFO: flag metadata %d for inclusion in hathitrust", mdID)
	var existCnt int64
	err := svc.DB.Table("hathitrust_statuses").Where("metadata_id=?", mdID).Count(&existCnt).Error
	if err != nil {
		return fmt.Errorf("unable to determine if metadata %d has hathitrust status: %s", mdID, err.Error())
	}
	if existCnt > 0 {
		return fmt.Errorf("hathitrust status already exists for metadata %d", mdID)
	}

	err = svc.DB.Model(&metadata{ID: mdID}).Update("hathitrust", 1).Error
	if err != nil {
		return fmt.Errorf("unable to flag metadata %d for hathitrust: %s", mdID, err.Error())
	}

	htStatus := hathitrustStatus{MetadataID: mdID, RequestedAt: time.Now()}
	err = svc.DB.Create(&htStatus).Error
	if err != nil {
		return fmt.Errorf("unable to create hathitrust status for metadata %d: %s", mdID, err.Error())
	}
	return nil
}

func (svc *serviceContext) removeMetadataFromHathiTrust(mdID int64) error {
	log.Printf("INFO: remove hathitrust flag and status for matadata %d", mdID)

	var htStatus hathitrustStatus
	err := svc.DB.Where("metadata_id=?", mdID).Limit(1).Find(&htStatus).Error
	if err != nil {
		return fmt.Errorf("unable to find existing hathitrust status for metadata %d: %s", mdID, err.Error())
	}

	if htStatus.PackageSubmittedAt != nil || htStatus.MetadataSubmittedAt != nil {
		return fmt.Errorf("metadata %d has started the hathitrust submission process and cannot be unpublished", mdID)
	}

	if htStatus.ID > 0 {
		err = svc.DB.Delete(&htStatus).Error
		if err != nil {
			return fmt.Errorf("unable to delete hathitrust status for metadata %d: %s", mdID, err.Error())
		}
	}

	err = svc.DB.Model(&metadata{ID: mdID}).Update("hathitrust", 0).Error
	if err != nil {
		return fmt.Errorf("unable to remove hathitrust flag from metadata %d: %s", mdID, err.Error())
	}

	return nil
}
