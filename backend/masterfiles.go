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
	MetadataID      int64          `gorm:"column:metadata_id" json:"-"`
	ContainerTypeID *int64         `json:"-"`
	ContainerType   *containerType `gorm:"foreignKey:ContainerTypeID" jon:"containerType"`
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
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
}

type tag struct {
	ID  uint64 `json:"id"`
	Tag string `json:"tag"`
}

type masterFile struct {
	ID                int64          `json:"id"`
	PID               string         `gorm:"column:pid" json:"pid"`
	MetadataID        *int64         `gorm:"column:metadata_id" json:"metadataID,omitempty"`
	Metadata          *metadata      `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	ImageTechMeta     *imageTechMeta `gorm:"foreignKey:MasterFileID" json:"techMetadata,omitempty"`
	UnitID            int64          `gorm:"column:unit_id" json:"unitID"`
	Unit              *unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	ComponentID       int64          `gorm:"column:component_id" json:"componentID"`
	Filename          string         `json:"filename"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Tags              []*tag         `gorm:"many2many:master_file_tags" json:"tags"`
	Filesize          int64          `json:"filesize"`
	MD5               string         `gorm:"column:md5" json:"md5"`
	OriginalMfID      int64          `gorm:"column:original_mf_id" json:"originalID"`
	DateArchived      *time.Time     `json:"dateArchived"`
	DeaccessionedAt   *time.Time     `json:"deaccessionedAt"`
	DeaccessionedByID *int64         `gorm:"column:deaccessioned_by_id" json:"-"`
	DeaccessionedBy   *staffMember   `gorm:"foreignKey:DeaccessionedByID" json:"deaccessionedBy"`
	DeaccessionNote   string         `json:"deaccessionNote"`
	TranscriptionText string         `json:"transcription"`
	DateDlIngest      *time.Time     `gorm:"column:date_dl_ingest" json:"dateDLIngest"`
	DateDlUpdate      *time.Time     `gorm:"column:date_dl_update" json:"dateDLUpdate"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	ThumbnailURL      string         `gorm:"-" json:"thumbnailURL,omitempty"`
	ViewerURL         string         `gorm:"-" json:"viewerURL,omitempty"`
	Exemplar          bool           `json:"exemplar"`
	Locations         []location     `gorm:"many2many:master_file_locations" json:"locations"`
}

func (svc *serviceContext) getMasterFile(c *gin.Context) {
	mfID := c.Param("id")
	log.Printf("INFO: get master file %s details", mfID)
	var mf masterFile
	err := svc.DB.Preload("ImageTechMeta").Preload("DeaccessionedBy").Preload("Tags").
		Preload("Metadata").Preload("Locations").Find(&mf, mfID).Error
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

	c.JSON(http.StatusOK, out)
}
