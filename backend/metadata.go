package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type internalMetadata struct {
	Title            string `json:"title"`
	CallNumber       string `json:"callNumber"`
	CreatorName      string `json:"creatorName"`
	CreatorType      string `json:"creatorType"`
	Year             string `json:"year"`
	PublicationPlace string `json:"publicationPlace"`
	Location         string `json:"location"`
	PreviewURL       string `json:"previewURL"`
	ObjectURL        string `json:"objectURL"`
}

type uvaMAP struct {
	Doc struct {
		Field []struct {
			Text   string `xml:",chardata"`
			Name   string `xml:"name,attr"`
			Type   string `xml:"type,attr"`
			Access string `xml:"access,attr"`
		} `xml:"field"`
	} `xml:"doc"`
}

func (svc *serviceContext) getMetadata(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: get metadata %s details", mdID)

	var md metadata
	err := svc.DB.Preload("UseRight").Preload("OCRHint").
		Preload("AvailabilityPolicy").Preload("ExternalSystem").Find(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to load metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type mdResp struct {
		Metadata metadata          `json:"metadata"`
		Extended *internalMetadata `json:"details"`
		Error    string            `json:"error"`
	}
	out := mdResp{Metadata: md}
	if md.Type == "SirsiMetadata" || md.Type == "XmlMetadata" {
		ext, err := svc.getUVAMapData(md.PID)
		if err != nil {
			log.Printf("ERROR: unable to get extended metadata for sirsi/cml %s: %s", md.PID, err.Error())
			out.Error = err.Error()
		} else {
			out.Extended = ext
		}
	}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getUVAMapData(pid string) (*internalMetadata, error) {
	url := fmt.Sprintf("%s/api/metadata/%s?type=uvamap", svc.ExternalSystems.TSAPI, pid)
	resp, err := svc.getRequest(url)
	if err != nil {
		return nil, fmt.Errorf("%d %s", err.StatusCode, err.Message)
	}
	var uvamap uvaMAP
	parseErr := xml.Unmarshal(resp, &uvamap)
	if err != nil {
		return nil, fmt.Errorf("unable to parse uvamp response: %s", parseErr.Error())
	}
	var detail internalMetadata
	for _, f := range uvamap.Doc.Field {
		switch f.Name {
		case "displayTitle":
			detail.Title = f.Text
		case "callNumber":
			detail.CallNumber = f.Text
		case "creator":
			detail.CreatorName = f.Text
			detail.CreatorType = f.Type
		case "keyDate":
			detail.Year = f.Text
		case "physLocation":
			detail.Location = f.Text
		case "pubProdDistPlace":
			if detail.PublicationPlace == "" {
				detail.PublicationPlace = f.Text
			} else {
				joined := fmt.Sprintf("%s, %s", f.Text, detail.PublicationPlace)
				detail.PublicationPlace = joined
			}
		case "uri":
			if f.Access == "raw object" {
				detail.ObjectURL = f.Text
			} else if f.Access == "preview" {
				url := strings.Replace(f.Text, "!125,200", "!240,385", 1)
				detail.PreviewURL = url
			}
		}
	}

	return &detail, nil
}
