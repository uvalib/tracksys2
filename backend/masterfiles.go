package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type location struct {
	ID              int64          `json:"id"`
	MetadataID      int64          `gorm:"column:metadata_id" json:"metadataID"`
	ContainerTypeID *int64         `json:"-"`
	ContainerType   *containerType `gorm:"foreignKey:ContainerTypeID" json:"containerType"`
	ContainerID     string         `gorm:"column:container_id" json:"containerID"`
	FolderID        string         `gorm:"column:folder_id" json:"folderID"`
	Notes           string         `json:"notes"`
}

type imageTechMeta struct {
	ID           int64      `json:"id"`
	MasterFileID int64      `json:"-"`
	ImageFormat  string     `json:"imageFormat"`
	Width        uint       `json:"width"`
	Height       uint       `json:"height"`
	Resolution   uint       `json:"resolution"`
	ColorSpace   string     `json:"colorSpace"`
	Depth        uint       `json:"depth"`
	Compression  string     `json:"compression"`
	ColorProfile string     `json:"colorProfile"`
	Equipment    string     `json:"equipment"`
	Software     string     `json:"software"`
	Model        string     `json:"model"`
	ExifVersion  string     `json:"exifVersion"`
	CaptureDate  *time.Time `json:"captureDate"`
	ISO          uint       `gorm:"column:iso" json:"iso"`
	ExposureBias string     `json:"exposureBias"`
	ExposureTime string     `json:"exposureTime"`
	Aperture     string     `json:"aperture"`
	FocalLength  float64    `json:"focalLength"`
	Orientation  uint       `json:"orientation"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
}

type masterFileAudit struct {
	ID            int64       `json:"id"`
	MasterFileID  int64       `json:"masterFileID"`
	MasterFile    *masterFile `gorm:"foreignKey:MasterFileID" json:"masterFile,omitempty"`
	ArchiveExists bool        `json:"archiveExists"`
	ChecksumMatch bool        `json:"checksumMatch"`
	AuditChecksum string      `json:"auditChecksum"`
	IIIFExists    bool        `gorm:"column:iiif_exists" json:"iiifExists"`
	AuditedAt     time.Time   `json:"auditedAt"`
}

type masterFile struct {
	ID                int64            `json:"id"`
	PID               string           `gorm:"column:pid" json:"pid"`
	MetadataID        *int64           `gorm:"column:metadata_id" json:"metadataID,omitempty"`
	Metadata          *metadata        `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	ImageTechMeta     *imageTechMeta   `gorm:"foreignKey:MasterFileID" json:"techMetadata,omitempty"`
	UnitID            int64            `gorm:"column:unit_id" json:"unitID"`
	Unit              *unit            `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	ComponentID       int64            `gorm:"column:component_id" json:"componentID"`
	Filename          string           `json:"filename"`
	Title             string           `json:"title"`
	Description       string           `json:"description"`
	Sensitive         bool             `json:"sensitive"`
	Tags              []tag            `gorm:"many2many:master_file_tags" json:"tags"`
	Filesize          int64            `json:"filesize"`
	MD5               string           `gorm:"column:md5" json:"md5"`
	PHash             *uint64          `gorm:"column:phash" json:"-"`
	OriginalMfID      int64            `gorm:"column:original_mf_id" json:"originalID"`
	DateArchived      *time.Time       `json:"dateArchived"`
	DeaccessionedAt   *time.Time       `json:"deaccessionedAt"`
	DeaccessionedByID *int64           `gorm:"column:deaccessioned_by_id" json:"-"`
	DeaccessionedBy   *staffMember     `gorm:"foreignKey:DeaccessionedByID" json:"deaccessionedBy"`
	DeaccessionNote   string           `json:"deaccessionNote"`
	TranscriptionText string           `json:"transcription"`
	DateDlIngest      *time.Time       `gorm:"column:date_dl_ingest" json:"dateDLIngest"`
	DateDlUpdate      *time.Time       `gorm:"column:date_dl_update" json:"dateDLUpdate"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"-"`
	ThumbnailURL      string           `gorm:"-" json:"thumbnailURL,omitempty"`
	ViewerURL         string           `gorm:"-" json:"viewerURL,omitempty"`
	Exemplar          bool             `json:"exemplar"`
	Locations         []location       `gorm:"many2many:master_file_locations" json:"locations"`
	Audit             *masterFileAudit `gorm:"foreignKey:MasterFileID" json:"audit,omitempty"`
}

func (svc *serviceContext) getMasterFile(c *gin.Context) {
	mfID := c.Param("id")
	log.Printf("INFO: get master file %s details", mfID)
	var mf masterFile
	err := svc.DB.Preload("ImageTechMeta").Preload("DeaccessionedBy").Preload("Tags").
		Preload("Metadata").Preload("Metadata.OCRHint").Preload("Audit").
		Preload("Locations").Preload("Locations.ContainerType").Find(&mf, mfID).Error
	if err != nil {
		log.Printf("ERROR: unable to get masterfile %s: %s", mfID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type mfResp struct {
		MasterFile masterFile `json:"masterFile"`
		OrderID    uint64     `json:"orderID"`
		ThumbURL   string     `json:"thumbURL"`
		ViewerURL  string     `json:"viewerURL"`
		PrevID     int64      `json:"prevID,omitempty"`
		NextID     int64      `json:"nextID,omitempty"`
	}
	out := mfResp{MasterFile: mf}

	mfPID := mf.PID
	if mf.OriginalMfID > 0 {
		var originalMF masterFile
		err := svc.DB.Find(&originalMF, mf.OriginalMfID).Error
		if err != nil {
			log.Printf("ERROR: unbale to get original masterfile %d for clone %d: %s", mf.OriginalMfID, mf.ID, err.Error())
		}
		mfPID = originalMF.PID
	}

	out.ThumbURL = fmt.Sprintf("%s/%s/full/!240,385/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
	if mf.Metadata == nil {
		out.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
	} else {
		mfRegex := regexp.MustCompile(fmt.Sprintf(`^%09d_\w{4,}\.`, mf.UnitID))
		pageNum := 0
		if mfRegex.MatchString(mf.Filename) {
			pagePart := strings.Split(mf.Filename, "_")[1]
			pageNum, _ = strconv.Atoi(strings.Split(pagePart, ".")[0])
		}
		if mf.DateDlIngest != nil {
			out.ViewerURL = fmt.Sprintf("%s/view/%s", svc.ExternalSystems.Curio, mf.Metadata.PID)
			if pageNum > 0 {
				out.ViewerURL += fmt.Sprintf("?page=%d", pageNum)
			}
		} else {
			out.ViewerURL = fmt.Sprintf("%s/view/%s?unit=%d", svc.ExternalSystems.Curio, mf.Metadata.PID, mf.UnitID)
			if pageNum > 0 {
				out.ViewerURL += fmt.Sprintf("&page=%d", pageNum)
			}
		}
	}

	err = svc.DB.Table("units").Select("order_id").Where("id=?", mf.UnitID).Find(&out.OrderID).Error
	if err != nil {
		log.Printf("ERROR: unable to get masterfile %s order info: %s", mfID, err.Error())
	}

	log.Printf("INFO: get a sorted list of other masterfile ids that belong to unit %d", out.MasterFile.UnitID)
	mfIDs := make([]int64, 0)
	err = svc.DB.Table("master_files").Where("unit_id=?", out.MasterFile.UnitID).Select("id").Order("filename asc").Find(&mfIDs).Error
	if err != nil {
		log.Printf("ERROR: unable to find masterfiles for unit %d: %s", out.MasterFile.UnitID, err.Error())
	} else {
		for idx, mfID := range mfIDs {
			if mfID == out.MasterFile.ID {
				if idx > 0 {
					out.PrevID = mfIDs[idx-1]
				}
				if idx < len(mfIDs)-1 {
					out.NextID = mfIDs[idx+1]
				}
				break
			}
		}
	}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) updateMasterFile(c *gin.Context) {
	mfID := c.Param("id")
	var mf masterFile
	err := svc.DB.Preload("Locations").Preload("Locations.ContainerType").Preload("ImageTechMeta").Find(&mf, mfID).Error
	if err != nil {
		log.Printf("ERROR: unable to get masterfile %s: %s", mfID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var req struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		Orientation     uint   `json:"orientation"`
		MetadataID      int64  `json:"metadataID"`
		UpdateLocation  bool   `json:"updateLocation"`
		ContainerTypeID int64  `json:"containerTypeID"`
		ContainerID     string `json:"containerID"`
		FolderID        string `json:"folderID"`
		Notes           string `json:"notes"`
	}
	err = c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid update unit %s request: %s", mfID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: update masterfile %d with %+v", mf.ID, req)
	mf.Title = req.Title
	mf.Description = req.Description
	mf.ImageTechMeta.Orientation = req.Orientation
	mf.MetadataID = &req.MetadataID

	err = svc.DB.Model(&mf).Select("Title", "Description", "MetadataID").Updates(mf).Error
	if err != nil {
		log.Printf("ERROR: unable to update master file %d: %s", mf.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	svc.DB.Model(&mf.ImageTechMeta).Select("Orientation").Updates(mf.ImageTechMeta)
	if req.UpdateLocation {
		if len(mf.Locations) == 0 {
			newLoc := location{ContainerTypeID: &req.ContainerTypeID, ContainerID: req.ContainerID,
				FolderID: req.FolderID, Notes: req.Notes, MetadataID: *mf.MetadataID}
			svc.DB.Create(&newLoc)
			svc.DB.Exec("INSERT into master_file_locations (master_file_id, location_id) values (?,?)", mf.ID, newLoc.ID)
		} else {
			loc := mf.Locations[0]
			loc.ContainerTypeID = &req.ContainerTypeID
			loc.ContainerID = req.ContainerID
			loc.FolderID = req.FolderID
			loc.Notes = req.Notes
			svc.DB.Model(&loc).Select("ContainerTypeID", "ContainerID", "FolderID", "Notes").Updates(loc)
		}
	}

	svc.DB.Preload("ImageTechMeta").Preload("DeaccessionedBy").Preload("Tags").Preload("Metadata").Preload("Locations").Preload("Locations.ContainerType").Find(&mf, mfID)
	c.JSON(http.StatusOK, mf)
}

func (svc *serviceContext) addMasterFileTag(c *gin.Context) {
	mfID := c.Param("id")
	tID := c.Query("tag")
	log.Printf("INFO: add tag id %s to masterfile %s", tID, mfID)

	var mf masterFile
	err := svc.DB.Find(&mf, mfID).Error
	if err != nil {
		log.Printf("ERROR: unable to load masterfile %s to add tag: %s", mfID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var tgtTag tag
	err = svc.DB.Find(&tgtTag, tID).Error
	if err != nil {
		log.Printf("ERROR: unable to load tag %s for master file add: %s", tID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	svc.DB.Exec("INSERT into master_file_tags (master_file_id, tag_id) values (?,?)", mf.ID, tgtTag.ID)
	c.String(http.StatusOK, "added")
}

func (svc *serviceContext) removeMasterFileTag(c *gin.Context) {
	mfID := c.Param("id")
	tID := c.Query("tag")
	log.Printf("INFO: remove tag id %s from masterfile %s", tID, mfID)
	err := svc.DB.Exec("DELETE from master_file_tags where master_file_id=? and tag_id=?", mfID, tID).Error
	if err != nil {
		log.Printf("ERROR: unable to remove tag %s from master file %s: %s", tID, mfID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "removed")
}
