package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type hathitrustStatus struct {
	ID                  uint       `json:"id"`
	MetadataID          int64      `json:"metadataID"`
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
		pCreateDate, err := time.Parse("2006-01-02", req.PackageCreatedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid package create date", req.PackageCreatedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid package create date", req.PackageCreatedAt))
			return
		}
		md.HathiTrustStatus.PackageCreatedAt = &pCreateDate
		updates = append(updates, "PackageCreatedAt")
	}

	if req.PackageSubmittedAt != "" {
		pSubmitDate, err := time.Parse("2006-01-02", req.PackageSubmittedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid package submit date", req.PackageSubmittedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid package submit date", req.PackageSubmittedAt))
			return
		}
		md.HathiTrustStatus.PackageSubmittedAt = &pSubmitDate
		updates = append(updates, "PackageSubmittedAt")
	}

	if req.FinishedAt != "" {
		finishDate, err := time.Parse("2006-01-02", req.FinishedAt)
		if err != nil {
			log.Printf("ERROR: %s is not a valid finish date", req.FinishedAt)
			c.String(http.StatusBadRequest, fmt.Sprintf("%s is not a valid finish date", req.FinishedAt))
			return
		}
		md.HathiTrustStatus.FinishedAt = &finishDate
		updates = append(updates, "FinishedAt")
	}

	if req.MetadataSubmittedAt != "" {
		mSubmitDate, err := time.Parse("2006-01-02", req.MetadataSubmittedAt)
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
