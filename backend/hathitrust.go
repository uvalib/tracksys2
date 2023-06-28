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
