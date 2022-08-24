package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type apTrustStatus struct {
	ID          int64     `json:"-"`
	MetadataPID string    `gorm:"column:metadata_pid" json:"metadataPID"`
	Etag        string    `json:"etag"`
	ObjectID    string    `json:"objectID"`
	FinishedAt  time.Time `json:"submittedAt"`
}

type availabilityPolicy struct {
	ID   int64  `json:"id"`
	PID  string `gorm:"column:pid" json:"pid"`
	Name string `json:"name"`
}

type ocrHint struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	OCRCandidate bool   `json:"ocrCandidate"`
}

type ocrLanguageHint struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type useRight struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	URI            string `json:"uri"`
	Statement      string `json:"statement"`
	CommercialUse  bool   `json:"commercial_use"`
	EducationalUse bool   `json:"educational_use"`
	Modifications  bool   `json:"modifications"`
}

type externalSystem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	PublicURL string `json:"public_url"`
	APIURL    string `gorm:"column:api_url" jjson:"api_url"`
}

type metadata struct {
	ID                   int64              `json:"id"`
	PID                  string             `gorm:"column:pid" json:"pid"`
	Type                 string             `json:"type"`
	Title                string             `json:"title"`
	Barcode              string             `json:"barcode,omitempty"`
	CallNumber           string             `json:"call_number,omitempty"`
	CatalogKey           string             `json:"catalog_key,omitempty"`
	CreatorName          string             `json:"creator_name"`
	DescMetadata         string             `json:"desc_metadata,omitempty"`
	CollectionFacet      string             `json:"collection_facet,omitempty"`
	UseRightID           uint               `json:"-"`
	UseRight             useRight           `gorm:"foreignKey:UseRightID" json:"use_right"`
	OCRHintID            uint               `json:"-"`
	OCRHint              ocrHint            `gorm:"foreignKey:OCRHintID" json:"ocr_hint"`
	OCRLanguageHint      string             `json:"ocrLanguageHint"`
	AvailabilityPolicyID uint               `json:"-"`
	AvailabilityPolicy   availabilityPolicy `gorm:"foreignKey:AvailabilityPolicyID" json:"availability_policy"`
	ExternalSystemID     uint               `json:"-"`
	ExternalSystem       externalSystem     `gorm:"foreignKey:ExternalSystemID" json:"external_system"`
	DateDLIngest         sql.NullTime       `gorm:"date_dl_ingest" json:"date_dl_ingest"`
	UpdatedAt            sql.NullTime       `json:"updated_at"`
}

type metadataVersion struct {
	ID            int64  `json:"id"`
	MetadataID    int64  `json:"metadata_id"`
	StaffMemberID int64  `json:"staff_member_id"`
	DescMetadata  string `json:"desc_metadata"`
	VersionTag    string `json:"version_tag"`
	Comment       string `json:"comment"`
}

func (svc *serviceContext) getMetadata(c *gin.Context) {

}
