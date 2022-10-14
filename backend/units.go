package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type intendedUse struct {
	ID                    int64  `json:"id"`
	Description           string `json:"name"`
	DeliverableFormat     string `json:"deliverableFormat"`
	DeliverableResolution string `json:"deliverableResolution"`
}

type attachment struct {
	ID          int64  `json:"id"`
	UnitID      int64  `json:"unitID"`
	Filename    string `json:"filename"`
	MD5         string `gorm:"column:md5" json:"md5"`
	Description string `json:"description"`
}

type lastError struct {
	ID    uint64 `json:"jobID"`
	Error string `json:"error"`
}

type unit struct {
	ID                          int64        `json:"id"`
	OrderID                     int64        `json:"orderID"`
	Order                       *order       `json:"order,omitempty"`
	MetadataID                  *int64       `json:"metadataID,omitempty"`
	Metadata                    *metadata    `json:"metadata,omitempty"`
	UnitStatus                  string       `json:"status"`
	IntendedUseID               *int64       `json:"-"`
	IntendedUse                 *intendedUse `gorm:"foreignKey:IntendedUseID" json:"intendedUse"`
	PatronSourceURL             string       `json:"patronSourceURL"`
	IncludeInDL                 bool         `gorm:"column:include_in_dl" json:"includeInDL"`
	RemoveWatermark             bool         `json:"removeWatermark"`
	Reorder                     bool         `json:"reorder"`
	CompleteScan                bool         `json:"completeScan"`
	ThrowAway                   bool         `json:"throwAway"`
	OcrMasterFiles              bool         `gorm:"column:ocr_master_files" json:"ocrMasterFiles"`
	MasterFiles                 []masterFile `gorm:"foreignKey:UnitID" json:"masterFiles"`
	Attachments                 []attachment `gorm:"foreignKey:UnitID" json:"attachments"`
	SpecialInstructions         string       `json:"specialInstructions"`
	StaffNotes                  string       `json:"staffNotes"`
	MasterFilesCount            uint         `json:"masterFilesCount"`
	DateArchived                *time.Time   `json:"dateArchived"`
	DatePatronDeliverablesReady *time.Time   `json:"datePatronDeliverablesReady"`
	DateDLDeliverablesReady     *time.Time   `gorm:"column:date_dl_deliverables_ready" json:"dateDLDeliverablesReady"`
	CreatedAt                   time.Time    `json:"-"`
	UpdatedAt                   time.Time    `json:"-"`
	ProjectID                   int64        `gorm:"-" json:"projectID,omitempty"`
	LastError                   *lastError   `gorm:"-" json:"lastError,omitempty"`
}

func (svc *serviceContext) getUnit(c *gin.Context) {
	unitID := c.Param("id")
	log.Printf("INFO: get unit %s details", unitID)
	var unitDetail unit
	err := svc.DB.Preload("IntendedUse").Preload("Attachments").Preload("Metadata").Find(&unitDetail, unitID).Error
	if err != nil {
		log.Printf("ERROR: unable to get detauls for unit %s: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var project struct {
		ID int64
	}
	err = svc.DB.Table("projects").Select("id").Where("unit_id=?", unitID).First(&project).Error
	if err != nil {
		log.Printf("ERROR: unable to get project id for unit %s: %s", unitID, err.Error())
	} else {
		unitDetail.ProjectID = project.ID
	}

	log.Printf("INFO: check for recent errors for unit %d", unitDetail.ID)
	var lastStatus jobStatus
	err = svc.DB.Where("originator_type=?", "Unit").Where("originator_id=?", unitDetail.ID).Order("created_at desc").First(&lastStatus).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false {
			log.Printf("ERROR: failed to find job statuses for unit %d", unitDetail.ID)
		}
	} else {
		if lastStatus.Status == "failure" {
			le := lastError{ID: lastStatus.ID, Error: lastStatus.Error}
			unitDetail.LastError = &le
		}
	}

	c.JSON(http.StatusOK, unitDetail)
}

func (svc *serviceContext) updateUnit(c *gin.Context) {
	unitID := c.Param("id")
	var unitDetail unit
	err := svc.DB.Find(&unitDetail, unitID).Error
	if err != nil {
		log.Printf("ERROR: unable to get unit %s details before update: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var req struct {
		Status              string `json:"status"`
		PatronSourceURL     string `json:"patronSourceURL"`
		SpecialInstructions string `json:"specialInstructions"`
		StaffNotes          string `json:"staffNotes"`
		CompleteScan        bool   `json:"completeScan"`
		ThrowAway           bool   `json:"throwAway"`
		OrderID             int64  `json:"orderID"`
		MetadataID          int64  `json:"metadataID"`
		IntendedUseID       int64  `json:"intendedUseID"`
		OCRMasterFiles      bool   `json:"ocrMasterFiles"`
		RemoveWaterMark     bool   `json:"removeWatermark"`
		IncludeInDL         bool   `json:"includeInDL"`
	}
	err = c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid update unit %s request: %s", unitID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: update unit %d", unitDetail.ID)
	unitDetail.UnitStatus = req.Status
	unitDetail.PatronSourceURL = req.Status
	unitDetail.SpecialInstructions = req.SpecialInstructions
	unitDetail.StaffNotes = req.StaffNotes
	unitDetail.CompleteScan = req.CompleteScan
	unitDetail.ThrowAway = req.ThrowAway
	unitDetail.OrderID = req.OrderID
	unitDetail.MetadataID = &req.MetadataID
	unitDetail.IntendedUseID = &req.IntendedUseID
	unitDetail.OcrMasterFiles = req.OCRMasterFiles
	unitDetail.RemoveWatermark = req.RemoveWaterMark
	unitDetail.IncludeInDL = req.IncludeInDL
	err = svc.DB.Model(&unitDetail).
		Select(
			"Status", "PatronSourceURL", "SpecialInstructions", "StaffNotes", "CompleteScan", "ThrowAway",
			"OrderID", "MetadataID", "IntendedUseID", "OcrMasterFiles", "RemoveWatermark", "IncludeInDL").
		Updates(unitDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to update unit %d: %s", unitDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	svc.DB.Preload("IntendedUse").Preload("Metadata").Preload("Attachments").Find(&unitDetail, unitID)
	c.JSON(http.StatusOK, unitDetail)
}

func (svc *serviceContext) getUnitMasterfiles(c *gin.Context) {
	unitID := c.Param("id")
	log.Printf("INFO: get unit %s masterfiles", unitID)

	var masterFiles []*masterFile
	err := svc.DB.Where("unit_id=?", unitID).Preload("Metadata").Order("filename asc").Find(&masterFiles).Error
	if err != nil {
		log.Printf("ERROR: unable to get materfiles for unit %s: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for idx, mf := range masterFiles {
		mfPID := mf.PID
		if mf.OriginalMfID > 0 {
			var originalMF masterFile
			err := svc.DB.Find(&originalMF, mf.OriginalMfID).Error
			if err != nil {
				log.Printf("ERROR: unbale to get original masterfile %d for clone %d: %s", mf.OriginalMfID, mf.ID, err.Error())
			}
			mfPID = originalMF.PID
		}

		mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		if mf.MetadataID != nil {
			mf.ViewerURL = fmt.Sprintf("%s/view/%s?unit=%s", svc.ExternalSystems.Curio, mf.Metadata.PID, unitID)
			if idx > 0 {
				mf.ViewerURL += fmt.Sprintf("&page=%d", (idx + 1))
			}
		} else {
			mf.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		}
	}
	c.JSON(http.StatusOK, masterFiles)
}
