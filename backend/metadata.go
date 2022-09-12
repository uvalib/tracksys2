package main

import (
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
	CommercialUse  bool   `json:"commercialUse"`
	EducationalUse bool   `json:"educationalUse"`
	Modifications  bool   `json:"modifications"`
}

type externalSystem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	PublicURL string `json:"publicURL"`
	APIURL    string `gorm:"column:api_url" jjson:"apiURL"`
}

type preservationTier struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type metadata struct {
	ID                   int64              `json:"id"`
	PID                  string             `gorm:"column:pid" json:"pid"`
	Type                 string             `json:"type"`
	ParentMetadataID     uint64             `gorm:"column:parent_metadata_id" json:"parentID"`
	Title                string             `json:"title"`
	Barcode              string             `json:"barcode"`
	CallNumber           string             `json:"callNumber"`
	CatalogKey           string             `json:"catalogKey"`
	CreatorName          string             `json:"creatorName"`
	CreatorDeathDate     uint64             `json:"creatorDeathDate"`
	DescMetadata         string             `json:"descMetadata"`
	CollectionFacet      string             `json:"collectionFacet"`
	UseRightID           uint               `json:"-"`
	UseRight             useRight           `gorm:"foreignKey:UseRightID" json:"useRight"`
	UseRightRationale    string             `json:"useRightRationale"`
	OCRHintID            *uint              `json:"-"`
	OCRHint              *ocrHint           `gorm:"foreignKey:OCRHintID" json:"ocrHint"`
	OCRLanguageHint      string             `json:"ocrLanguageHint"`
	AvailabilityPolicyID uint               `json:"-"`
	AvailabilityPolicy   availabilityPolicy `gorm:"foreignKey:AvailabilityPolicyID" json:"availability"`
	ExternalSystemID     *uint              `json:"-"`
	ExternalSystem       *externalSystem    `gorm:"foreignKey:ExternalSystemID" json:"externalSystem"`
	SupplementalSystemID *uint              `json:"-"`
	SupplementalSystem   *externalSystem    `gorm:"foreignKey:SupplementalSystemID" json:"supplementalSystem"`
	PreservationTierID   uint               `json:"-"`
	PreservationTier     preservationTier   `gorm:"foreignKey:PreservationTierID" json:"preservationTier"`
	DPLA                 bool               `gorm:"column:dpla" json:"dpla"`
	IsManuscript         bool               `json:"isManuscript"`
	DateDLIngest         *time.Time         `gorm:"date_dl_ingest" json:"dateDLIngest"`
	DateDLUpdate         *time.Time         `gorm:"date_dl_update" json:"dateDLUpdate"`
	CreatedAt            *time.Time         `json:"-"`
	UpdatedAt            *time.Time         `json:"-"`
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
	err := svc.DB.Preload("UseRight").Preload("OCRHint").Preload("AvailabilityPolicy").
		Preload("ExternalSystem").Preload("SupplementalSystem").
		Preload("PreservationTier").Find(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to load metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type mdResp struct {
		Metadata metadata          `json:"metadata"`
		Extended *internalMetadata `json:"details"`
		VirgoURL string            `json:"virgoURL"`
		Error    string            `json:"error"`
	}
	out := mdResp{Metadata: md}
	if md.Type == "SirsiMetadata" || md.Type == "XmlMetadata" {
		parsedDetail, err := svc.getUVAMapData(md.PID)
		if err != nil {
			log.Printf("ERROR: unable to get extended metadata for sirsi/cml %s: %s", md.PID, err.Error())
			out.Error = err.Error()
		} else {
			out.Extended = parsedDetail
			if md.DateDLIngest != nil {
				out.VirgoURL = fmt.Sprintf("%s/sources/uva_library/items/%s", svc.ExternalSystems.Virgo, md.CatalogKey)
			}
		}
	}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getMetadataRelatedItems(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: get related units and orders for metadata %s", mdID)

	var units []unit
	err := svc.DB.Debug().Where("metadata_id=?", mdID).Preload("IntendedUse").
		Preload("Order").Preload("Order.Customer").Preload("Order.Agency").
		Find(&units).Error
	if err != nil {
		log.Printf("ERROR: unable to get units related to metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(units) == 0 {
		// No units but one master file is an indicator that descriptive XML
		// metadata was created specifically for the master file after initial ingest.
		// This is usually the case with image collections where each image has its own descriptive metadata.
		// In this case, there is no direct link from metadata to unit. Must find it by
		// going through the master file that this metadata describes
		var mfCnt int64
		err := svc.DB.Table("master_files").Where("metadata_id=?", mdID).Count(&mfCnt).Error
		if err != nil {
			log.Printf("ERROR: unable to get master file count for metadata %s: %s", mdID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if mfCnt == 1 {
			var mf masterFile
			err = svc.DB.Preload("Unit").Preload("Unit.Order").
				Preload("Unit.Order.Customer").Preload("Unit.Order.Agency").
				Preload("Unit.IntendedUse").Where("metadata_id=?", mdID).First(&mf).Error
			if err != nil {
				log.Printf("ERROR: unabel to get masterfile unit for metadata %s: %s", mdID, err.Error())
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			units = append(units, *mf.Unit)
		}
	}

	c.JSON(http.StatusOK, units)
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
