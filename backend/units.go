package main

import (
	"encoding/csv"
	"encoding/json"
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
	Order                       *order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	MetadataID                  *int64       `json:"metadataID,omitempty"`
	Metadata                    *metadata    `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	UnitStatus                  string       `json:"status"`
	IntendedUseID               *int64       `json:"-"`
	IntendedUse                 *intendedUse `gorm:"foreignKey:IntendedUseID" json:"intendedUse"`
	PatronSourceURL             string       `json:"patronSourceURL"`
	IncludeInDL                 bool         `gorm:"column:include_in_dl" json:"includeInDL"`
	RemoveWatermark             bool         `json:"removeWatermark"`
	Reorder                     bool         `json:"reorder"`
	CompleteScan                bool         `json:"completeScan"`
	ThrowAway                   bool         `json:"throwAway"`
	OCRMasterFiles              bool         `gorm:"column:ocr_master_files" json:"ocrMasterFiles"`
	NumMasterFiles              uint         `json:"masterFilesCount"`
	Attachments                 []attachment `gorm:"foreignKey:UnitID" json:"attachments"`
	SpecialInstructions         string       `json:"specialInstructions"`
	StaffNotes                  string       `json:"staffNotes"`
	DateArchived                *time.Time   `json:"dateArchived"`
	DatePatronDeliverablesReady *time.Time   `json:"datePatronDeliverablesReady"`
	DateDLDeliverablesReady     *time.Time   `gorm:"column:date_dl_deliverables_ready" json:"dateDLDeliverablesReady"`
	CreatedAt                   time.Time    `json:"-"`
	UpdatedAt                   time.Time    `json:"-"`
	ProjectID                   int64        `gorm:"-" json:"projectID,omitempty"`
	LastError                   *lastError   `gorm:"-" json:"lastError,omitempty"`
	RelatedUnitIDs              []int64      `gorm:"-" json:"relatedUnits,omitempty"`
	HasText                     bool         `gorm:"-" json:"hasText"`
}

type projectLookupResponse struct {
	Exists      bool   `json:"exists"`
	ProjectID   int64  `json:"projectID"`
	Workflow    string `json:"workflow"`
	CurrentStep string `json:"currentStep"`
	Finished    bool   `json:"finished"`
}

func (svc *serviceContext) validateUnit(c *gin.Context) {
	unitID := c.Param("id")
	log.Printf("INFO: validate unit %s exists", unitID)
	var cnt int64
	err := svc.DB.Table("units").Where("id=?", unitID).Count(&cnt).Error
	if err != nil {
		log.Printf("INFO: error validating unit %s: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if cnt == 0 {
		log.Printf("INFO: unit %s not found", unitID)
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.String(http.StatusOK, "exists")
}

func (svc *serviceContext) deleteUnit(c *gin.Context) {
	unitID := c.Param("id")
	log.Printf("INFO: delete unit %s details", unitID)
	var mfCount int64
	err := svc.DB.Table("master_files").Where("unit_id=?", unitID).Count(&mfCount).Error
	if err != nil {
		log.Printf("ERROR: unable to determine if unit %s has master files: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.DB.Delete(&unit{}, unitID).Error
	if err != nil {
		log.Printf("ERROR: unable to delete unit %s: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// check id a project exists, and cancel it if so
	var lookupResp projectLookupResponse
	respBytes, reqErr := svc.getRequest(fmt.Sprintf("%s/projects/lookup?unit=%s", svc.ExternalSystems.Projects, unitID))
	if reqErr != nil {
		log.Printf("ERROR: unable to determine if unit %s has a project: %s", unitID, reqErr.Message)
	} else {
		if err := json.Unmarshal(respBytes, &lookupResp); err != nil {
			log.Printf("ERROR: unable to parse response for project lookup: %s", err.Error())
		} else if lookupResp.Exists {
			log.Printf("INFO: unit %s is associated with project %d; cancel it", unitID, lookupResp.ProjectID)
			if rErr := svc.projectsPost(fmt.Sprintf("projects/%d/cancel", lookupResp.ProjectID), getJWT(c)); rErr != nil {
				log.Printf("ERROR: unable to cancel project %d: %s", lookupResp.ProjectID, rErr.Message)
			} else {
				log.Printf("INFO: project %d associated with deleted unit %s has been canceled", lookupResp.ProjectID, unitID)
			}
		}
	}

	c.String(http.StatusOK, "deleted")
}

func (svc *serviceContext) getUnit(c *gin.Context) {
	unitID := c.Param("id")
	log.Printf("INFO: get unit %s details", unitID)
	var unitDetail unit
	err := svc.DB.Preload("IntendedUse").Preload("Attachments").Preload("Order").
		Preload("Metadata").Preload("Metadata.OCRHint").First(&unitDetail, unitID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("INFO: unit %s not found", unitID)
			c.String(http.StatusNotFound, fmt.Sprintf("%s not found", unitID))
		} else {
			log.Printf("ERROR: unable to get detauls for unit %s: %s", unitID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	log.Printf("INFO: check for project associated with unit %s", unitID)
	var lookupResp projectLookupResponse
	respBytes, reqErr := svc.getRequest(fmt.Sprintf("%s/projects/lookup?unit=%s", svc.ExternalSystems.Projects, unitID))
	if reqErr != nil {
		log.Printf("ERROR: lookup project for unit %s failed: %s", unitID, reqErr.Message)
	} else {
		if err := json.Unmarshal(respBytes, &lookupResp); err != nil {
			log.Printf("ERROR: unable to parse response for project lookup: %s", err.Error())
		} else if lookupResp.Exists {
			log.Printf("INFO: unit %s is associated with project %d", unitID, lookupResp.ProjectID)
			unitDetail.ProjectID = lookupResp.ProjectID
		}
	}

	log.Printf("INFO: check if unit %d has any ocr/transcription text", unitDetail.ID)
	var txtCnt int64
	err = svc.DB.Table("master_files").Where("unit_id=? and transcription_text is not null and transcription_text != ?", unitDetail.ID, "").Count(&txtCnt).Error
	if err != nil {
		log.Printf("ERROR: unable to determine if unit %d has text associated with its masterfiles: %s", unitDetail.ID, err.Error())
	} else {
		unitDetail.HasText = txtCnt > 0
	}

	log.Printf("INFO: check for recent errors for unit %d", unitDetail.ID)
	var lastStatus jobStatus
	err = svc.DB.Where("originator_type=?", "Unit").Where("originator_id=?", unitDetail.ID).Order("created_at desc").Limit(1).Find(&lastStatus).Error
	if err != nil {
		log.Printf("ERROR: failed to find job statuses for unit %d", unitDetail.ID)
	} else {
		if lastStatus.Status == "failure" {
			le := lastError{ID: lastStatus.ID, Error: lastStatus.Error}
			unitDetail.LastError = &le
		}
	}

	log.Printf("INFO: get a list if other unit ids that belong to the same order as unit %d", unitDetail.ID)
	err = svc.DB.Table("units").Where("order_id=?", unitDetail.OrderID).Select("id").Find(&unitDetail.RelatedUnitIDs).Error
	if err != nil {
		log.Printf("ERROR: unable to find related units for unit %d: %s", unitDetail.ID, err.Error())
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

	// Only 1 unit per metadata record can be flagged for inclusion in the DL (Virgo) enforce this now
	if req.IncludeInDL {
		err = svc.validateIncludeInDL(*unitDetail.MetadataID, unitDetail.ID)
		if err != nil {
			log.Printf("ERROR: unit %d failed include in dl validation: %s", unitDetail.ID, err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	var lookupResp projectLookupResponse
	respBytes, reqErr := svc.getRequest(fmt.Sprintf("%s/projects/lookup?unit=%s", svc.ExternalSystems.Projects, unitID))
	if reqErr != nil {
		log.Printf("ERROR: lookup project for unit %s with status %s failed: %s", unitID, req.Status, reqErr.Message)
	} else {
		if err := json.Unmarshal(respBytes, &lookupResp); err != nil {
			log.Printf("ERROR: unable to parse response for project lookup: %s", err.Error())
		}
	}

	// do not allow unit to be set to done if it is associated with a project that is not on the finalize step
	if lookupResp.Exists && req.Status == "done" && lookupResp.Finished == false && lookupResp.CurrentStep != "Finialize" {
		log.Printf("INFO: cannot set unit to done; it is tied to project %d on step %s", lookupResp.ProjectID, lookupResp.CurrentStep)
		c.String(http.StatusPreconditionFailed, "Cannot set status to done; associated project is in progress")
		return
	}

	updateMasterFileMetadata := false
	if unitDetail.MetadataID == nil && req.MetadataID != 0 {
		log.Printf("INFO: unit %d update changes metadata from none to %d; master files must be updated", unitDetail.ID, req.MetadataID)
		updateMasterFileMetadata = true
	} else if unitDetail.MetadataID != nil && unitDetail.MetadataID != &req.MetadataID {
		log.Printf("INFO: unit %d update changes metadata from %d to %d; master files must be updated", unitDetail.ID, *unitDetail.MetadataID, req.MetadataID)
		updateMasterFileMetadata = true
	}

	log.Printf("INFO: update unit %d with %+v", unitDetail.ID, req)
	unitDetail.UnitStatus = req.Status
	unitDetail.PatronSourceURL = req.PatronSourceURL
	unitDetail.SpecialInstructions = req.SpecialInstructions
	unitDetail.StaffNotes = req.StaffNotes
	unitDetail.CompleteScan = req.CompleteScan
	unitDetail.ThrowAway = req.ThrowAway
	unitDetail.OrderID = req.OrderID
	unitDetail.MetadataID = &req.MetadataID
	unitDetail.IntendedUseID = &req.IntendedUseID
	unitDetail.OCRMasterFiles = req.OCRMasterFiles
	unitDetail.RemoveWatermark = req.RemoveWaterMark
	unitDetail.IncludeInDL = req.IncludeInDL
	err = svc.DB.Model(&unitDetail).
		Select(
			"UnitStatus", "PatronSourceURL", "SpecialInstructions", "StaffNotes", "CompleteScan", "ThrowAway",
			"OrderID", "MetadataID", "IntendedUseID", "OcrMasterFiles", "RemoveWatermark", "IncludeInDL").
		Updates(unitDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to update unit %d: %s", unitDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if req.Status == "canceled" || req.Status == "done" {
		if lookupResp.Exists {
			log.Printf("INFO: unit %s is associated with project %d", unitID, lookupResp.ProjectID)
			updateURL := fmt.Sprintf("projects/%d/done", lookupResp.ProjectID)
			if req.Status == "canceled" {
				updateURL = fmt.Sprintf("projects/%d/cancel", lookupResp.ProjectID)
			}
			log.Printf("INFO: update status of project %d to reflect unit status %s", lookupResp.ProjectID, req.Status)
			if rErr := svc.projectsPost(updateURL, getJWT(c)); rErr != nil {
				log.Printf("ERROR: unable to update project %d status: %s", lookupResp.ProjectID, rErr.Message)
			}
		}
	}

	if updateMasterFileMetadata {
		log.Printf("INFO: update masterfiles metadata to %d", req.MetadataID)
		mfErr := svc.DB.Debug().Exec("update master_files set metadata_id=? where unit_id=?", req.MetadataID, unitDetail.ID).Error
		if mfErr != nil {
			log.Printf("ERROR: unable to update unit %d masterfiles with new metadata: %s", unitDetail.ID, mfErr.Error())
			c.String(http.StatusInternalServerError, mfErr.Error())
			return
		}
	}

	svc.DB.Preload("IntendedUse").Preload("Attachments").Preload("Order").Preload("Metadata").Preload("Metadata.OCRHint").Find(&unitDetail, unitID)
	svc.DB.Table("units").Where("order_id=?", unitDetail.OrderID).Select("id").Find(&unitDetail.RelatedUnitIDs)
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

	var currMFMetadataID int64
	metadataPageMap := make(map[int64]int)
	for _, mf := range masterFiles {
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
			if currMFMetadataID == 0 {
				// this is the first setting just init page for this meatdata record to 1
				currMFMetadataID = mf.Metadata.ID
				metadataPageMap[mf.Metadata.ID] = 1
			} else {
				if currMFMetadataID != mf.Metadata.ID {
					_, ok := metadataPageMap[mf.Metadata.ID]
					if ok == false {
						metadataPageMap[mf.Metadata.ID] = 1
					}
				}
			}
			mf.ViewerURL = fmt.Sprintf("%s/view/%s?unit=%s", svc.ExternalSystems.Curio, mf.Metadata.PID, unitID)
			if metadataPageMap[mf.Metadata.ID] > 1 {
				mf.ViewerURL += fmt.Sprintf("&page=%d", metadataPageMap[mf.Metadata.ID])
			}
			metadataPageMap[mf.Metadata.ID] = metadataPageMap[mf.Metadata.ID] + 1
		} else {
			mf.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		}
	}
	c.JSON(http.StatusOK, masterFiles)
}

func (svc *serviceContext) setExemplar(c *gin.Context) {
	unitID := c.Param("id")
	mfID := c.Param("mfid")
	log.Printf("INFO: set master file %s as exemplar for unit %s", mfID, unitID)
	var exemplar masterFile
	now := time.Now()
	err := svc.DB.Find(&exemplar, mfID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("INFO: unable to set master file %s as exemplar; file not found", mfID)
			c.String(http.StatusNotFound, fmt.Sprintf("master file %s not found", mfID))
		} else {
			log.Printf("ERROR: unable to set master file %s as exemplar: %s", mfID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}
	err = svc.DB.Exec("update master_files set exemplar=? where unit_id=?", false, unitID).Error
	if err != nil {
		log.Printf("ERROR: unable to clear unit %s exemplar: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	exemplar.Exemplar = true
	exemplar.UpdatedAt = now
	err = svc.DB.Model(&exemplar).Select("Exemplar", "UpdatedAt").Updates(exemplar).Error
	if err != nil {
		log.Printf("ERROR: unable to set master file %s as exemplar: %s", mfID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "exemplar set")
}

func (svc *serviceContext) getUnitCloneSources(c *gin.Context) {
	tgtUnitID := c.Param("id")
	log.Printf("INFO: get clone sources for unit %s", tgtUnitID)
	var tgtUnit unit
	err := svc.DB.Find(&tgtUnit, tgtUnitID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: unit %s not found", tgtUnitID)
			c.String(http.StatusNotFound, fmt.Sprintf("unit %s not found", tgtUnitID))
		} else {
			log.Printf("ERROR: unable to retrieve unit %s: %s", tgtUnitID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	var units []unit
	log.Printf("INFO: unit %d clone source has metadata [%d]", tgtUnit.ID, *tgtUnit.MetadataID)
	err = svc.DB.Where("reorder=? and units.id<>? and throw_away=? and units.metadata_id=?", false, tgtUnit.ID, false, tgtUnit.MetadataID).Find(&units).Error
	if err != nil {
		log.Printf("ERROR: unable to get source units for clone of unit %s: %s", tgtUnitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, units)
}

func (svc *serviceContext) exportUnitCSV(c *gin.Context) {
	uID := c.Param("id")
	log.Printf("INFO: export csv for unit %s", uID)
	var mfs []masterFile
	err := svc.DB.Where("unit_id=?", uID).Preload("Locations").Preload("Locations.ContainerType").Find(&mfs).Error
	if err != nil {
		log.Printf("ERROR: unable to get master files for unit %s csv export: %s", uID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	cw := csv.NewWriter(c.Writer)
	csvHead := []string{"pid", "title", "description", "location", "iiif"}
	cw.Write(csvHead)
	for _, mf := range mfs {
		line := make([]string, 0)
		line = append(line, mf.PID)
		line = append(line, mf.Title)
		line = append(line, mf.Description)
		if len(mf.Locations) > 0 {
			loc := mf.Locations[0]
			locStr := fmt.Sprintf("%s %s, Folder %s", loc.ContainerType.Name, loc.ContainerID, loc.FolderID)
			line = append(line, locStr)
		} else {
			line = append(line, "")
		}
		iifURL := fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mf.PID)
		line = append(line, iifURL)
		cw.Write(line)
	}

	cw.Flush()
}

func (svc *serviceContext) validateIncludeInDL(metadataID int64, tgtUnitID int64) error {
	log.Printf("INFO: validate  include in dl setting for metadata %d", metadataID)

	var dlCnt int64
	err := svc.DB.Table("units").Where("metadata_id=? and include_in_dl=? and id != ?", metadataID, 1, tgtUnitID).Count(&dlCnt).Error
	if err != nil {
		return err
	}

	if dlCnt > 0 {
		return fmt.Errorf("metadata %d already has another unit flagged for inclusion in virgo", metadataID)
	}

	return nil
}

type unitLocations struct {
	ID            int64      `json:"id"`
	OrderID       int64      `gorm:"column:order_id" json:"orderID"`
	OrderTitle    string     `gorm:"column:order_title" json:"orderTitle"`
	UnitStatus    string     `json:"unitStatus"`
	DateArchived  *time.Time `json:"dateArchived"`
	IntendedUseID int64      `json:"intendedUseID"`
	CompleteScan  bool       `json:"completeScan"`
	StaffNotes    string     `json:"staffNotes"`
}

func (svc *serviceContext) getLocationUnits(c *gin.Context) {
	tgtLocID := c.Param("id")
	log.Printf("INFO: get units associated with location %s", tgtLocID)
	sql := "select u.id,u.order_id,o.order_title, u.unit_status, u.date_archived, u.intended_use_id, u.complete_scan,u.staff_notes "
	sql += " from master_file_locations l "
	sql += " inner join master_files m on m.id = master_file_id "
	sql += " inner join units u on u.id = unit_id "
	sql += " inner join orders o on o.id = order_id "
	sql += " where location_id=? group by u.id"
	var out []unitLocations
	err := svc.DB.Raw(sql, tgtLocID).Scan(&out).Error
	if err != nil {
		log.Printf("ERROR: unable to get units related to locatipon %s: %s", tgtLocID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, out)
}
