package main

import "time"

type containerType struct {
	ID         int64
	Name       string
	hasFolders bool
}

type location struct {
	ID              int64
	MetadataID      int64 `gorm:"column:metadata_id"`
	ContainerTypeID int64
	ContainerType   containerType `gorm:"foreignKey:ContainerTypeID"`
	ContainerID     string        `gorm:"column:container_id"`
	FolderID        string        `gorm:"column:folder_id"`
	Notes           string
}

type imageTechMeta struct {
	ID           int64
	MasterFileID int64
	ImageFormat  string
	Width        uint
	Height       uint
	Resolution   uint
	ColorSpace   string
	Depth        uint
	Compression  string
	ColorProfile string
	Equipment    string
	Software     string
	Model        string
	ExifVersion  string
	CaptureDate  *time.Time
	ISO          uint `gorm:"column:iso"`
	ExposureBias string
	ExposureTime string
	Aperture     string
	FocalLength  float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type tag struct {
	ID  uint64
	Tag string
}

type masterFile struct {
	ID                int64
	PID               string        `gorm:"column:pid"`
	MetadataID        *int64        `gorm:"column:metadata_id"`
	ImageTechMeta     imageTechMeta `gorm:"foreignKey:MasterFileID"`
	UnitID            int64         `gorm:"column:unit_id"`
	Filename          string
	Title             string
	Description       string
	Locations         []location `gorm:"many2many:master_file_locations"`
	Tags              []*tag     `gorm:"many2many:master_file_tags"`
	Filesize          int64
	MD5               string `gorm:"column:md5"`
	OriginalMfID      *int64 `gorm:"column:original_mf_id"`
	DateArchived      *time.Time
	DeaccessionedAt   *time.Time
	DeaccessionedByID *int64 `gorm:"column:deaccessioned_by_id"`
	DeaccessionNote   string
	TranscriptionText string
	DateDlIngest      *time.Time `gorm:"column:date_dl_ingest"`
	DateDlUpdate      *time.Time `gorm:"column:date_dl_update"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (mf *masterFile) location() *location {
	if len(mf.Locations) == 0 {
		return nil
	}
	return &mf.Locations[0]
}
