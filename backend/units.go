package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	UpdatedAt                   time.Time    `json:"-"`
	ProjectID                   int64        `gorm:"-" json:"projectID,omitempty"`
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
