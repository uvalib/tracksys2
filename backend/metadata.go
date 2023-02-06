package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	ID                   int64               `json:"id"`
	PID                  string              `gorm:"column:pid" json:"pid"`
	Type                 string              `json:"type"`
	ParentMetadataID     uint64              `gorm:"column:parent_metadata_id" json:"parentID"`
	Title                string              `json:"title"`
	Barcode              *string             `json:"barcode"`
	CallNumber           *string             `json:"callNumber"`
	CatalogKey           *string             `json:"catalogKey"`
	CreatorName          *string             `json:"creatorName"`
	CreatorDeathDate     *uint64             `json:"creatorDeathDate"`
	DescMetadata         *string             `json:"descMetadata"`
	CollectionID         *string             `json:"collectionID"`    // internal usage to track a collection ID
	CollectionFacet      *string             `json:"collectionFacet"` // used at index to put item in collection in DL; EX: Ganon Project, McGregor
	UseRightID           *int64              `json:"-"`
	UseRight             *useRight           `gorm:"foreignKey:UseRightID" json:"useRight"`
	UseRightRationale    string              `json:"useRightRationale"`
	OCRHintID            *int64              `json:"-"`
	OCRHint              *ocrHint            `gorm:"foreignKey:OCRHintID" json:"ocrHint"`
	OCRLanguageHint      string              `json:"ocrLanguageHint"`
	AvailabilityPolicyID *int64              `json:"-"`
	AvailabilityPolicy   *availabilityPolicy `gorm:"foreignKey:AvailabilityPolicyID" json:"availabilityPolicy"`
	ExternalSystemID     *int64              `json:"-"`
	ExternalSystem       *externalSystem     `gorm:"foreignKey:ExternalSystemID" json:"externalSystem"`
	ExternalURI          *string             `gorm:"column:external_uri" json:"externalURI"`
	SupplementalSystemID *int64              `json:"-"`
	SupplementalSystem   *externalSystem     `gorm:"foreignKey:SupplementalSystemID" json:"supplementalSystem"`
	SupplementalURI      *string             `json:"supplementalURI"`
	PreservationTierID   *int64              `json:"-"`
	PreservationTier     *preservationTier   `gorm:"foreignKey:PreservationTierID" json:"preservationTier"`
	DPLA                 bool                `gorm:"column:dpla" json:"dpla"`
	IsManuscript         bool                `json:"isManuscript"`
	IsPersonalItem       bool                `json:"isPersonalItem"`
	DateDLIngest         *time.Time          `gorm:"date_dl_ingest" json:"dateDLIngest"`
	DateDLUpdate         *time.Time          `gorm:"date_dl_update" json:"dateDLUpdate"`
	CreatedAt            *time.Time          `json:"-"`
	UpdatedAt            *time.Time          `json:"-"`
}

func (m *metadata) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(m).Update("pid", fmt.Sprintf("tsb:%d", m.ID)).Error
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

type asMetadata struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	CreatedBy       string `json:"created_by"`
	CreateTime      string `json:"create_time"`
	Level           string `json:"level"`
	URL             string `json:"url"`
	Repo            string `json:"repo"`
	CollectionTitle string `json:"collection_title"`
	Language        string `json:"language"`
	Dates           string `json:"dates"`
	PublishedAt     string `json:"published_at,omitempty"`
}

type jstorMetadata struct {
	ID           string `json:"id"`
	SSID         string `json:"ssid"`
	Title        string `json:"title"`
	Description  string `json:"desc"`
	Creator      string `json:"creator"`
	Date         string `json:"date"`
	CollectionID string `json:"collectionID"`
	Collection   string `json:"collection"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type apolloMetadata struct {
	PID                  string `json:"pid"`
	Type                 string `json:"type"`
	Title                string `json:"title"`
	CollectionPID        string `json:"collectionPID"`
	CollectionTitle      string `json:"collectionTitle"`
	CollectionBarcode    string `json:"collectionBarcode"`
	CollectionCatalogKey string `json:"collectionCatalogKey"`
	ItemURL              string `json:"itemURL"`
	CollectionURL        string `json:"collectionURL"`
}

type apolloNode struct {
	PID  string `json:"pid"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
	Value string `json:"value"`
}

type apolloContainer struct {
	PID  string `json:"pid"`
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
	Children []apolloNode `json:"children"`
}

type apolloResp struct {
	Collection apolloContainer `json:"collection"`
	Item       apolloContainer `json:"item"`
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

type metadataDetailResponse struct {
	Metadata     *metadata         `json:"metadata"`
	Units        []*unit           `json:"units"`
	MasterFiles  []*masterFile     `json:"masterFiles,omitempty"`
	Extended     *internalMetadata `json:"details"`
	ArchiveSpace *asMetadata       `json:"asDetails"`
	JSTOR        *jstorMetadata    `json:"jstorDetails"`
	Apollo       *apolloMetadata   `json:"apolloDetails"`
	VirgoURL     string            `json:"virgoURL"`
	Error        string            `json:"error"`
}

type metadataRequest struct {
	Type                 string `json:"type"`
	ExternalSystemID     int64  `json:"externSystemID"`
	ExternalURI          string `json:"externalURI"`
	Title                string `json:"title"`
	CallNumber           string `json:"callNumber"`
	Author               string `json:"author"`
	CatalogKey           string `json:"catalogKey"`
	Barcode              string `json:"barcode"`
	PersonalItem         bool   `json:"personalItem"`
	Manuscript           bool   `json:"manuscript"`
	OCRHint              int64  `json:"ocrHint"`
	OCRLanguageHint      string `json:"ocrLanguageHint"`
	PreservationTierID   int64  `json:"preservationTier"`
	AvailabilityPolicyID int64  `json:"availabilityPolicy"`
	UseRightID           int64  `json:"useRight"`
	UseRightRationale    string `json:"useRightRationale"`
	DPLA                 bool   `json:"inDPLA"`
	CollectionID         string `json:"collectionID"`
	CollectionFacet      string `json:"CollectionFacet"`
}

var modsTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<mods xmlns="http://www.loc.gov/mods/v3" version="3.6">
   <titleInfo>
      <title>[TITLE]</title>
   </titleInfo>
`
var modsAuthor = `   <name>
      <role>
         <roleTerm type="text">creator</roleTerm>
      </role>
      <namePart>[AUTHOR]</namePart>
   </name>
`

func (svc *serviceContext) getMetadata(c *gin.Context) {
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s", c.Param("iid"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	log.Printf("INFO: get metadata %d details", mdID)
	resp, err := svc.loadMetadataDetails(mdID)
	if err != nil {
		log.Printf("ERROR: get metadata %d failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Metadata.ID == 0 {
		log.Printf("INFO: metadata %d not found", mdID)
		c.String(http.StatusNotFound, "metadata record not found")
		return
	}

	log.Printf("INFO: get related units and orders for metadata %d", resp.Metadata.ID)
	// NOTE: Manually calculate the master files count and return it as num_master_files instead of using the inaccurate cache
	mfCnt := "(select count(*) from master_files m inner join units u on u.id=m.unit_id where u.metadata_id=units.metadata_id) as num_master_files"
	err = svc.DB.Where("metadata_id=?", resp.Metadata.ID).Preload("IntendedUse").
		Preload("Order").Preload("Order.Customer").Preload("Order.Agency").
		Select("units.*", mfCnt).Find(&resp.Units).Error
	if err != nil {
		log.Printf("ERROR: unable to get related units for %d: %s", resp.Metadata.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(resp.Units) == 0 {
		// No units but one master file is an indicator that descriptive XML
		// metadata was created specifically for the master file after initial ingest.
		// This is usually the case with image collections where each image has its own descriptive metadata.
		// In this case, there is no direct link from metadata to unit. Must find it by
		// going through the master file that this metadata describes
		// NOTE: this also applies to AS metadata. Often, there is one large unit represeting a box of meterial.
		// master files a grouped by page number and assigned to an AS metadata record. This AS metadata record will have
		// no units associated with it, but 1 or more master files. There is special handing in IIIF manifest to process these.
		log.Printf("INFO: no units directly found for metadata %d; searching master files...", resp.Metadata.ID)
		err = svc.DB.Preload("Unit").Preload("Unit.Order").
			Preload("Unit.Order.Customer").Preload("Unit.Order.Agency").
			Preload("Unit.IntendedUse").Where("metadata_id=?", resp.Metadata.ID).
			Find(&resp.MasterFiles).Error
		if err != nil {
			log.Printf("ERROR: unable to get masterfile unit for metadata %d: %s", resp.Metadata.ID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("INFO: %d master files found that are associated with metadata %d", len(resp.MasterFiles), resp.Metadata.ID)
		for _, f := range resp.MasterFiles {
			f.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, f.PID)
			f.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, f.PID)
			unique := true
			for _, u := range resp.Units {
				if u.ID == f.UnitID {
					unique = false
					break
				}
			}
			if unique {
				resp.Units = append(resp.Units, f.Unit)
			}
		}
	}

	for _, u := range resp.Units {
		u.MetadataID = &resp.Metadata.ID
		u.Metadata = resp.Metadata
	}
	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) createMetadata(c *gin.Context) {
	var req metadataRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create metadata request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: create %s metadata request", req.Type)

	// create new record and set common attributes
	createTime := time.Now()
	newMD := metadata{Type: req.Type, Title: req.Title, IsPersonalItem: req.PersonalItem, IsManuscript: req.Manuscript, CreatedAt: &createTime}
	if req.Author != "" {
		newMD.CreatorName = &req.Author
	}
	if req.CollectionID != "" {
		newMD.CollectionID = &req.CollectionID
	}
	if req.CollectionFacet != "" {
		newMD.CollectionFacet = &req.CollectionFacet
	}
	if req.OCRHint > 0 {
		newMD.OCRHintID = &req.OCRHint
		if req.OCRHint == 1 && req.OCRLanguageHint != "" {
			newMD.OCRLanguageHint = req.OCRLanguageHint
		}
	}
	if req.PreservationTierID > 0 {
		newMD.PreservationTierID = &req.PreservationTierID
	}

	// For non-external, set digital library attributes
	if req.Type != "ExternalMetadata" {
		newMD.AvailabilityPolicyID = &req.AvailabilityPolicyID
		newMD.UseRightID = &req.UseRightID
		newMD.UseRightRationale = req.UseRightRationale
		newMD.DPLA = req.DPLA
	}

	if req.Type == "XmlMetadata" {
		xmlMD := strings.Replace(modsTemplate, "[TITLE]", req.Title, 1)
		if req.Author != "" {
			xmlMD += strings.Replace(modsAuthor, "[AUTHOR]", req.Author, 1)
		}
		xmlMD += "</mods>"
		newMD.DescMetadata = &xmlMD
	} else if req.Type == "SirsiMetadata" {
		newMD.Barcode = &req.Barcode
		newMD.CallNumber = &req.CallNumber
		newMD.CatalogKey = &req.CatalogKey
	} else if req.Type == "ExternalMetadata" {
		newMD.ExternalURI = &req.ExternalURI
		newMD.ExternalSystemID = &req.ExternalSystemID
	}

	err = svc.DB.Create(&newMD).Error
	if err != nil {
		log.Printf("ERROR: unable to create metadata record: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load full details of newly created metadata record %d", newMD.ID)
	resp, err := svc.loadMetadataDetails(newMD.ID)
	if err != nil {
		log.Printf("ERROR: get metadata %d failed: %s", newMD.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) deleteMetadata(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received request to delete metadata %s", mdID)

	hasRelated, err := svc.metadataHasRelatedItems(mdID, "units")
	if err != nil {
		log.Printf("ERROR: unable to determine if metadata %s has any units: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if hasRelated {
		log.Printf("INFO: unable to delete metadata %s because it has related units", mdID)
		c.String(http.StatusPreconditionFailed, "order has units and cannont be deleted")
		return
	}
	hasRelated, err = svc.metadataHasRelatedItems(mdID, "master_files")
	if err != nil {
		log.Printf("ERROR: unable to determine if metadata %s has any master files: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if hasRelated {
		log.Printf("INFO: unable to delete metadata %s because it has related master files", mdID)
		c.String(http.StatusPreconditionFailed, "order has units and cannont be deleted")
		return
	}

	log.Printf("INFO: metadata %s has no related items, proceeding with delete", mdID)
	err = svc.DB.Delete(&metadata{}, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to delete metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: metadata %s has been deleted", mdID)
	c.String(http.StatusOK, "deleted")
}

func (svc *serviceContext) metadataHasRelatedItems(mdID string, tableName string) (bool, error) {
	var cnt int64
	err := svc.DB.Table(tableName).Where("metadata_id=?", mdID).Count(&cnt).Error
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func (svc *serviceContext) updateMetadata(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received update request for metadata %s", mdID)
	var md metadata
	err := svc.DB.Preload("ExternalSystem").Find(&md, mdID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("INFO: metadata %s not found", mdID)
			c.String(http.StatusNotFound, fmt.Sprintf("metadata %s not found", mdID))
		} else {
			log.Printf("ERROR: unable to get metadata %s: %s ", mdID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("INFO: update request for metadata %s is a valid metadata record", mdID)

	var req metadataRequest
	err = c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create metadata request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: update metadata %d request with %+v", md.ID, req)
	fields := []string{"DPLA", "IsManuscript", "IsPersonalItem"}
	md.DPLA = req.DPLA
	md.IsManuscript = req.Manuscript
	md.IsPersonalItem = req.PersonalItem
	if req.OCRHint > 0 {
		md.OCRHintID = &req.OCRHint
		fields = append(fields, "OCRHintID")
	}
	if req.OCRLanguageHint != "" {
		md.OCRLanguageHint = req.OCRLanguageHint
		fields = append(fields, "OCRLanguageHint")
	}
	if req.PreservationTierID > 0 {
		md.PreservationTierID = &req.PreservationTierID
		fields = append(fields, "PreservationTierID")
	}
	if req.AvailabilityPolicyID > 0 {
		md.AvailabilityPolicyID = &req.AvailabilityPolicyID
		fields = append(fields, "AvailabilityPolicyID")
	}
	if req.UseRightID > 0 {
		md.UseRightID = &req.UseRightID
		fields = append(fields, "UseRightID")
	}
	if req.UseRightRationale != "" {
		md.UseRightRationale = req.UseRightRationale
		fields = append(fields, "UseRightRationale")
	}
	if req.CollectionID != "" {
		md.CollectionID = &req.CollectionID
		fields = append(fields, "CollectionID")
	}
	if req.CollectionFacet != "" {
		md.CollectionFacet = &req.CollectionFacet
		fields = append(fields, "CollectionFacet")
	}

	if md.Type == "SirsiMetadata" {
		md.Barcode = &req.Barcode
		md.CallNumber = &req.CallNumber
		md.CatalogKey = &req.CatalogKey
		md.Title = req.Title
		md.CreatorName = &req.Author
		fields = append(fields, "Barcode", "CallNumber", "CatalogKey", "Title", "CreatorName")
	}
	if md.Type == "ExternalMetadata" && md.ExternalSystem.Name == "ArchivesSpace" {
		md.ExternalURI = &req.ExternalURI
		fields = append(fields, "ExternalURI")
	}

	err = svc.DB.Model(&md).Select(fields).Updates(md).Error
	if err != nil {
		log.Printf("ERROR: unable to update metadata %d: %s", md.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: metadata %d updated, reloadaing details", md.ID)
	resp, err := svc.loadMetadataDetails(md.ID)
	if err != nil {
		log.Printf("ERROR: unable load updated metadata %d: %s", md.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, *resp)
}

func (svc *serviceContext) loadMetadataDetails(mdID int64) (*metadataDetailResponse, error) {
	var md metadata
	err := svc.DB.Preload("UseRight").Preload("OCRHint").Preload("AvailabilityPolicy").
		Preload("ExternalSystem").Preload("SupplementalSystem").
		Preload("PreservationTier").Limit(1).Find(&md, mdID).Error
	if err != nil {
		return nil, err
	}

	out := metadataDetailResponse{Metadata: &md, Units: make([]*unit, 0)}
	if md.ID == 0 {
		return &out, nil
	}

	if md.Type == "SirsiMetadata" || md.Type == "XmlMetadata" {
		log.Printf("INFO: get extended sirsi/xml metadata for %d", md.ID)
		parsedDetail, err := svc.getUVAMapData(md.PID)
		if err != nil {
			log.Printf("ERROR: unable to get extended metadata for sirsi/xml %s: %s", md.PID, err.Error())
			out.Error = err.Error()
		} else {
			out.Extended = parsedDetail
			if md.Type == "SirsiMetadata" && md.CatalogKey != nil {
				out.VirgoURL = fmt.Sprintf("%s/sources/uva_library/items/%s", svc.ExternalSystems.Virgo, *md.CatalogKey)
			}
			if md.Type == "XmlMetadata" && md.DateDLIngest != nil {
				out.VirgoURL = fmt.Sprintf("%s/sources/images/items/%s", svc.ExternalSystems.Virgo, md.PID)
			}
			log.Printf("INFO: look for metdata %d exemplar", mdID)
			var exemplar masterFile
			err = svc.DB.Where("metadata_id=? and exemplar=?", mdID, 1).Limit(1).Find(&exemplar).Error
			if err != nil {
				log.Printf("ERROR: unable to find exemplar for metadata %d: %s", mdID, err.Error())
			} else {
				if exemplar.ID > 0 {
					log.Printf("INFO: metadata %d has exemplar [%s]", mdID, exemplar.PID)
					out.Extended.PreviewURL = fmt.Sprintf("%s/%s/full/!240,385/0/default.jpg", svc.ExternalSystems.IIIF, exemplar.PID)
				} else {
					log.Printf("INFO: metadata %d does not have an exemplar", mdID)
				}
			}
		}
	} else {
		if md.ExternalSystem.Name == "ArchivesSpace" {
			log.Printf("INFO: get external ArchivesSpace metadata for %s", md.PID)
			raw, getErr := svc.getRequest(fmt.Sprintf("%s/archivesspace/lookup?pid=%s&uri=%s", svc.ExternalSystems.Jobs, md.PID, *md.ExternalURI))
			if getErr != nil {
				log.Printf("ERROR: unable to get archivesSpace metadata for %s: %s", md.PID, getErr.Message)
			} else {
				log.Printf("INFO: raw as response: %s", raw)
				var asData asMetadata
				err := json.Unmarshal(raw, &asData)
				if err != nil {
					log.Printf("ERROR: unable to parse AS response for %s: %s", md.PID, err.Error())
				} else {
					out.ArchiveSpace = &asData
				}
			}
		} else if md.ExternalSystem.Name == "JSTOR Forum" {
			log.Printf("INFO: get external JSTOR Forum metadata for %s", md.PID)
			var mfInfo struct {
				Filename string
			}
			err := svc.DB.Table("master_files").Where("metadata_id=?", md.ID).Select("filename").First(&mfInfo).Error
			if err != nil {
				log.Printf("ERROR: unable to get master file associated with jstor metadata %s: %s", md.PID, err.Error())
			} else {
				tgtFilename := strings.TrimSuffix(mfInfo.Filename, filepath.Ext(mfInfo.Filename))
				raw, getErr := svc.getRequest(fmt.Sprintf("%s/jstor/lookup?filename=%s", svc.ExternalSystems.Jobs, tgtFilename))
				if getErr != nil {
					log.Printf("ERROR: unable to get jstor metadata for %s: %s", md.PID, getErr.Message)
				} else {
					var jsData jstorMetadata
					err := json.Unmarshal(raw, &jsData)
					if err != nil {
						log.Printf("ERROR: unable to parse jstor response for %s: %s", md.PID, err.Error())
					} else {
						out.JSTOR = &jsData
					}
				}
			}
		} else if md.ExternalSystem.Name == "Apollo" {
			log.Printf("INFO: get external apollo metadata for %s", md.PID)
			raw, getErr := svc.getRequest(fmt.Sprintf("%s%s", svc.ExternalSystems.Apollo, *md.ExternalURI))
			if getErr != nil {
				log.Printf("ERROR: unable to get apollo metadata for %s: %s", md.PID, getErr.Message)
			} else {
				var apResp apolloResp
				err := json.Unmarshal(raw, &apResp)
				if err != nil {
					log.Printf("ERROR: unable to parse apollo response for %s: %s", md.PID, err.Error())
				} else {
					apollo := apolloMetadata{CollectionPID: apResp.Collection.PID, PID: apResp.Item.PID, Type: apResp.Item.Type.Name}
					apollo.ItemURL = fmt.Sprintf("%s/collections/%s?item=%s", svc.ExternalSystems.Apollo, apResp.Collection.PID, apResp.Item.PID)
					apollo.CollectionURL = fmt.Sprintf("%s/collections/%s", svc.ExternalSystems.Apollo, apResp.Collection.PID)
					for _, c := range apResp.Collection.Children {
						if c.Type.Name == "title" {
							apollo.CollectionTitle = c.Value
						} else if c.Type.Name == "barcode" {
							apollo.CollectionBarcode = c.Value
						} else if c.Type.Name == "catalogKey" {
							apollo.CollectionCatalogKey = c.Value
						}
					}
					for _, c := range apResp.Item.Children {
						if c.Type.Name == "title" {
							apollo.Title = c.Value
						}
					}
					out.Apollo = &apollo
				}
			}
		}
	}

	return &out, nil
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
			}
		}
	}

	return &detail, nil
}

func (svc *serviceContext) validateArchivesSpaceMetadata(c *gin.Context) {
	uri := c.Query("uri")
	if uri == "" {
		log.Printf("INFO: invalid archivesspace lookup request; missing uri param")
		c.String(http.StatusBadRequest, "uri param is required")
		return
	}
	log.Printf("INFO: validate archivesspace uri %s", uri)
	raw, getErr := svc.getRequest(fmt.Sprintf("%s/archivesspace/validate?url=%s", svc.ExternalSystems.Jobs, uri))
	if getErr != nil {
		log.Printf("ERROR: unable to validate archivesSpace uri %s: %d - %s", uri, getErr.StatusCode, getErr.Message)
		c.String(getErr.StatusCode, getErr.Message)
		return
	}

	var resp struct {
		URI    string     `json:"uri"`
		Detail asMetadata `json:"detail"`
	}
	resp.URI = string(raw)

	log.Printf("INFO: lookup details for %s", resp.URI)
	rawDetail, getErr := svc.getRequest(fmt.Sprintf("%s/archivesspace/lookup?uri=%s", svc.ExternalSystems.Jobs, resp.URI))
	if getErr != nil {
		log.Printf("ERROR: unable to validate archivesSpace uri %s: %d - %s", uri, getErr.StatusCode, getErr.Message)
		c.String(getErr.StatusCode, getErr.Message)
		return
	}
	parseErr := json.Unmarshal(rawDetail, &resp.Detail)
	if parseErr != nil {
		log.Printf("ERROR: unable to parse archivespace details for %s: %s", resp.URI, parseErr.Error())
		c.String(http.StatusInternalServerError, parseErr.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getXMLMetadata(c *gin.Context) {
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s for get xml", c.Param("iid"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	log.Printf("INFO: get xml metadata %d", mdID)
	var md metadata
	err := svc.DB.Find(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: get xml metadata %d failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, *md.DescMetadata)
}

func (svc *serviceContext) uploadXMLMetadata(c *gin.Context) {
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s for xml upload", c.Param("iid"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	log.Printf("INFO: get metadata %d details", mdID)
	var md metadata
	err := svc.DB.Find(&md, mdID).Error
	if err != nil {
		log.Printf("ERROR: get metadata %d for xml upload failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if md.Type != "XmlMetadata" {
		log.Printf("ERROR: xml upload to metadata %d is invalid; wrong metadata type %s", mdID, md.Type)
		c.String(http.StatusBadRequest, fmt.Sprintf("%d is not an  XML metadata record", md.ID))
		return
	}

	formFile, err := c.FormFile("xml")
	if err != nil {
		log.Printf("ERROR: Unable to get uploaded xml file: %s", err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("unable to get file: %s", err.Error()))
		return
	}

	savedXMLFile := path.Join("/tmp", formFile.Filename)
	err = c.SaveUploadedFile(formFile, savedXMLFile)
	if err != nil {
		log.Printf("ERROR: Unable to read uploaded xml file %s for metadata %d: %s", formFile.Filename, mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	xmlBytes, err := os.ReadFile(savedXMLFile)
	if err != nil {
		log.Printf("ERROR: unable to read uploaded xml file %s: %s", formFile.Filename, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	descMetadata := string(xmlBytes)
	md.DescMetadata = &descMetadata
	err = svc.DB.Model(&md).Select("DescMetadata").Updates(md).Error
	if err != nil {
		log.Printf("ERROR: update xml metadata %d failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if md.DateDLIngest != nil || md.DateDLUpdate != nil {
		log.Printf("INFO: call xml reindexing hook for previously published metadata %s", md.PID)
		_, putErr := svc.putRequest(fmt.Sprintf("%s/%d", svc.ExternalSystems.XMLIndex, md.ID))
		if putErr != nil {
			log.Printf("ERROR: request to reindex %s failed: %d:%s", md.PID, putErr.StatusCode, putErr.Message)
		} else {
			log.Printf("INFO: %s was successfully queued for reindex; update dates", md.PID)
			now := time.Now()
			md.DateDLUpdate = &now
			err = svc.DB.Model(&md).Select("DateDLUpdate").Updates(md).Error
			if err != nil {
				log.Printf("ERROR: update xml publish date for %s failed: %s", md.PID, err.Error())
			}
		}
	}

	c.String(http.StatusOK, *md.DescMetadata)
}
