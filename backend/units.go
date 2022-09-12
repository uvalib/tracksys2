package main

import "time"

type intendedUse struct {
	ID                    int64  `json:"id"`
	Description           string `json:"name"`
	DeliverableFormat     string `json:"deliverableFormat"`
	DeliverableResolution string `json:"deliverableResolution"`
}

type unit struct {
	ID                          int64        `json:"id"`
	OrderID                     int64        `json:"-"`
	Order                       order        `json:"order"`
	MetadataID                  *int64       `json:"metadataID"`
	UnitStatus                  string       `json:"status"`
	IntendedUseID               *int64       `json:"-"`
	IntendedUse                 *intendedUse `gorm:"foreignKey:IntendedUseID" json:"intendedUse"`
	IncludeInDL                 bool         `gorm:"column:include_in_dl" json:"includeInDL"`
	RemoveWatermark             bool         `json:"removeWatermark"`
	Reorder                     bool         `json:"reorder"`
	CompleteScan                bool         `json:"completeScan"`
	ThrowAway                   bool         `json:"throwAway"`
	OcrMasterFiles              bool         `gorm:"column:ocr_master_files" json:"ocrMasterFiles"`
	MasterFiles                 []masterFile `gorm:"foreignKey:UnitID" json:"masterFiles"`
	StaffNotes                  string       `json:"staffNotes"`
	MasterFilesCount            uint         `json:"masterFilesCount"`
	UnitExtentActual            uint         `json:"unitExtentActual"`
	DateArchived                *time.Time   `json:"dateArchived"`
	DatePatronDeliverablesReady *time.Time   `json:"datePatronDeliverablesReady"`
	DateDLDeliverablesReady     *time.Time   `gorm:"column:date_dl_deliverables_ready" json:"dateDLDeliverablesReady"`
	UpdatedAt                   time.Time    `json:"-"`
}
