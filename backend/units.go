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

type step struct {
	ID          int64  `json:"id"`
	WorkflowID  int64  `json:"workflowID"`
	StepType    uint   `json:"stepType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NextStepID  uint   `json:"nextStepID"`
	FailStepID  uint   `json:"failStepID"`
	OwnerType   uint   `json:"ownerType"`
}

type project struct {
	ID              int64      `json:"id"`
	WorkflowID      int64      `json:"-"`
	ContainerTypeID *int64     `json:"-"`
	UnitID          int64      `json:"-"`
	CurrentStepID   int64      `json:"-"`
	AddedAt         *time.Time `json:"addedAt,omitempty"`
	CategoryID      int64      `json:"-"`
	ItemCondition   uint       `json:"itemCondition"`
	ConditionNote   string     `json:"conditionNote,omitempty"`
}

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

	var project struct {
		ID int64
	}
	err = svc.DB.Table("projects").Select("id").Where("unit_id=?", unitID).Limit(1).Find(&project).Error
	if err != nil {
		log.Printf("ERROR: unable to get project id for unit %s: %s", unitID, err.Error())
	} else {
		unitDetail.ProjectID = project.ID
	}

	log.Printf("INFO: check if unit %d has any ocr/transcription text", unitDetail.ID)
	var txtCnt int64
	err = svc.DB.Table("master_files").Where("unit_id=? and transcription_text is not null and transcription_text != ?", unitDetail.ID, "").Count(&txtCnt).Error
	if err != nil {
		log.Printf("ERROR: unabble to determine if unit %d has text associated with its masterfiles: %s", unitDetail.ID, err.Error())
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

func (svc *serviceContext) createProject(c *gin.Context) {
	unitID := c.Param("id")
	var unitDetail unit
	err := svc.DB.Find(&unitDetail, unitID).Error
	if err != nil {
		log.Printf("ERROR: unable to get unit %s details before update: %s", unitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var projCnt int64
	err = svc.DB.Table("projects").Where("unit_id=?", unitDetail.ID).Count(&projCnt).Error
	if err != nil {
		log.Printf("ERROR: unable to determine if a project already exists for unit %d: %s", unitDetail.ID, err.Error())
	} else if projCnt > 0 {
		log.Printf("INFO: unable to create project for unit %d as it already has a project", unitDetail.ID)
		c.String(http.StatusConflict, "a project already exists for this unit")
		return
	}

	var req struct {
		WorkflowID      int64  `json:"workflowID"`
		ContainerTypeID int64  `json:"containerTypeID"`
		CategoryID      int64  `json:"categoryID"`
		Condition       uint   `json:"condition"`
		Notes           string `json:"notes"`
	}
	err = c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create project request for unit %s: %s", unitID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: lookup first step of new project for unit %d, workflow %d", unitDetail.ID, req.WorkflowID)
	var firstStep step
	err = svc.DB.Where("workflow_id=? and step_type=0", req.WorkflowID).First(&firstStep).Error
	if err != nil {
		log.Printf("ERROR: unable to get first step for workflow %d: %s", req.WorkflowID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: create project for unit %d", unitDetail.ID)
	now := time.Now()
	newProj := project{
		WorkflowID:    req.WorkflowID,
		UnitID:        unitDetail.ID,
		CurrentStepID: firstStep.ID,
		AddedAt:       &now,
		CategoryID:    req.CategoryID,
		ItemCondition: req.Condition,
		ConditionNote: req.Notes,
	}
	if req.ContainerTypeID != 0 {
		newProj.ContainerTypeID = &req.ContainerTypeID
	}
	err = svc.DB.Create(&newProj).Error
	if err != nil {
		log.Printf("ERROR: unable to create project for unit %d: %s", unitDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: new project %d created for unit %d", newProj.ID, unitDetail.ID)
	c.String(http.StatusOK, fmt.Sprintf("%d", newProj.ID))
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
