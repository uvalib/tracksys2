package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type subField struct {
	XMLName xml.Name `xml:"subfield"`
	Code    string   `xml:"code,attr"`
	Value   string   `xml:",chardata"`
}

type dataField struct {
	XMLName   xml.Name   `xml:"datafield"`
	Tag       string     `xml:"tag,attr"`
	Subfields []subField `xml:"subfield"`
	Value     string     `xml:",chardata"`
}

type controlField struct {
	XMLName xml.Name `xml:"controlfield"`
	Tag     string   `xml:"tag,attr"`
	Value   string   `xml:",chardata"`
}

type marcRecord struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	ControlFields []controlField `xml:"controlfield"`
	DataFields    []dataField    `xml:"datafield"`
}

type marcMetadata struct {
	XMLName xml.Name `xml:"collection"`
	Record  marcRecord
}

type sirsiResponse struct {
	CatalogKey        string `json:"catalogKey"`
	Barcode           string `json:"barcode"`
	CallNumber        string `json:"callNumber"`
	Title             string `json:"title"`
	CreatorName       string `json:"creatorName"`
	CreatorType       string `json:"creatorType"`
	Year              string `json:"year"`
	PublicationPlace  string `json:"publicationPlace"`
	Location          string `json:"location"`
	CollectionID      string `json:"collectionID"`
	UseRightName      string `json:"useRightName"`
	UseRightURI       string `json:"useRighURI"`
	UseRightStatement string `json:"useRightStatement"`
}

func (svc *serviceContext) lookupSirsiMetadata(c *gin.Context) {
	barcode := strings.TrimSpace(strings.ToUpper(c.Query("barcode")))
	catKey := strings.TrimSpace(strings.ToLower(c.Query("ckey")))
	if barcode == "" && catKey == "" {
		log.Printf("ERROR: sirsi lookup requires barcode or catkey")
		c.String(http.StatusBadRequest, "barcode or ckey required")
		return
	}
	resp, err := svc.doSirsiLookup(catKey, barcode)
	if err != nil {
		return
	}

	out := struct {
		*sirsiResponse
		ExistingPID string `json:"existingPID"`
		ExistingID  int64  `json:"existingID"`
		Exists      bool   `json:"exists"`
	}{
		sirsiResponse: resp,
	}

	var existMD metadata
	err = svc.DB.Where("barcode = ?", resp.Barcode).Limit(1).Find(&existMD).Error
	if err != nil {
		log.Printf("ERROR: failed check for existing metadata with barcode %s: %s", resp.Barcode, err.Error())
	} else {
		if existMD.ID > 0 {
			log.Printf("INFO: metadata with barcode [%s] catkey [%s] already exists with id [%d]", barcode, catKey, existMD.ID)
			out.Exists = true
			out.ExistingID = existMD.ID
			out.ExistingPID = existMD.PID
		}
	}
	c.JSON(http.StatusOK, out)
}

type solrDocument struct {
	FullRecord string `json:"fullrecord"`
}

type solrResponseHeader struct {
	Status int `json:"status,omitempty"`
}

type solrResponseDocuments struct {
	NumFound int            `json:"numFound,omitempty"`
	Start    int            `json:"start,omitempty"`
	Docs     []solrDocument `json:"docs,omitempty"`
}

type solrResponse struct {
	Header   solrResponseHeader    `json:"responseHeader,omitempty"`
	Response solrResponseDocuments `json:"response,omitempty"`
}

func (svc *serviceContext) doSirsiLookup(catKey, barcode string) (*sirsiResponse, error) {
	// prefer catkey over barcode
	url := fmt.Sprintf("%s/select?fl=fullrecord&q=barcode_a:%s", svc.ExternalSystems.Solr, barcode)
	if catKey != "" {
		url = fmt.Sprintf("%s/select?fl=fullrecord&q=id:%s", svc.ExternalSystems.Solr, catKey)
	}

	respStr, err := svc.getRequest(url)
	if err != nil {
		return nil, fmt.Errorf("getMarc from solr failed %d: %s", err.StatusCode, err.Message)
	}

	var solr solrResponse
	jErr := json.Unmarshal(respStr, &solr)
	if jErr != nil {
		return nil, jErr
	}
	rawMarc := []byte(solr.Response.Docs[0].FullRecord)

	var parsed marcMetadata
	parseErr := xml.Unmarshal(rawMarc, &parsed)
	if parseErr != nil {
		return nil, parseErr
	}

	if len(parsed.Record.ControlFields) == 0 && len(parsed.Record.DataFields) == 0 {
		return nil, fmt.Errorf("no matches found in sirsi")
	}

	log.Printf("INFO: extract fields from raw marc response")
	resp := sirsiResponse{CatalogKey: catKey}

	// catkey is in 001 of control fields. find it first
	for _, cf := range parsed.Record.ControlFields {
		if cf.Tag == "001" {
			resp.CatalogKey = cf.Value
			break
		}
	}

	type bcHashData struct {
		Barcode    string
		CallNumber string
		Location   string
	}

	// the remaining data ins in the datafields
	titleRegex := regexp.MustCompile(`\s*\/$`)                 // strip trailing /
	pubRegex := regexp.MustCompile(`(?:^\[|\]$|\.$|\]\.$|:$)`) // strip [] and trailing . or :
	for _, df := range parsed.Record.DataFields {
		if df.Tag == "100" {
			for _, sf := range df.Subfields {
				if sf.Code == "a" {
					resp.CreatorName = strings.TrimSpace(sf.Value)
					break
				}
			}
			resp.CreatorType = "personal"
		}
		if (df.Tag == "110" || df.Tag == "111") && resp.CreatorType == "" {
			for _, sf := range df.Subfields {
				if sf.Code == "a" {
					resp.CreatorName = strings.TrimSpace(sf.Value)
				}
				if sf.Code == "b" {
					sub := strings.TrimSpace(sf.Value)
					resp.CreatorName += fmt.Sprintf(" %s", sub)
				}
			}
			resp.CreatorType = "corporate"
		}
		if df.Tag == "245" {
			// Title; main is in 245a, subtitle in 245b
			for _, sf := range df.Subfields {
				if sf.Code == "a" {
					trimmedTitle := strings.TrimSpace(sf.Value)
					resp.Title = titleRegex.ReplaceAllString(trimmedTitle, "")
				}
			}
		}
		if df.Tag == "260" {
			// publication info, a=place, c = year
			for _, sf := range df.Subfields {
				if sf.Code == "a" {
					resp.PublicationPlace = strings.TrimSpace(sf.Value)
					resp.PublicationPlace = strings.TrimSpace(pubRegex.ReplaceAllString(resp.PublicationPlace, ""))
				}
				if sf.Code == "c" && resp.Year == "" {
					resp.Year = strings.TrimSpace(sf.Value)
					resp.Year = strings.TrimSpace(pubRegex.ReplaceAllString(resp.Year, ""))
				}
			}
		}
		if df.Tag == "852" {
			// 852c is collectionID
			for _, sf := range df.Subfields {
				if sf.Code == "c" {
					resp.CollectionID = strings.TrimSpace(sf.Value)
				}
			}
		}

		if df.Tag == "856" {
			// use rights are held in 856 r (uri), t (name/statement), u (item uri)
			for _, sf := range df.Subfields {
				if sf.Code == "t" {
					if resp.UseRightName == "" {
						resp.UseRightName = strings.TrimSpace(sf.Value)
					} else {
						resp.UseRightStatement = strings.TrimSpace(sf.Value)
					}
				}
				if sf.Code == "r" {
					resp.UseRightURI = strings.TrimSpace(sf.Value)
				}
			}
		}

		if df.Tag == "999" {
			// 999 repeats, 1 per barcode. Match the queried barcode or just pick first
			// subfields: i=barcode,l=location, a=call number
			bcData := bcHashData{}
			for _, sf := range df.Subfields {
				if sf.Code == "a" {
					bcData.CallNumber = strings.TrimSpace(sf.Value)
				}
				if sf.Code == "i" {
					bcData.Barcode = strings.ToUpper(strings.TrimSpace(sf.Value))
				}
				if sf.Code == "l" {
					bcData.Location = strings.TrimSpace(sf.Value)
				}
			}
			if bcData.Barcode == barcode || barcode == "" && resp.Barcode == "" {
				resp.Barcode = bcData.Barcode
				resp.CallNumber = bcData.CallNumber
				resp.Location = bcData.Location
			}
		}
	}

	if resp.CollectionID == "" {
		resp.CollectionID = resp.CallNumber
	}

	if resp.UseRightName == "" {
		log.Printf("INFO: no use right data found in sirsi response; default to CNE")
		var cne useRight
		dbErr := svc.DB.First(&cne, 1).Error
		if dbErr != nil {
			log.Printf("ERROR: unable to load CNE data: %s", dbErr.Error())
		} else {
			resp.UseRightName = cne.Name
			resp.UseRightURI = cne.URI
			resp.UseRightStatement = cne.Statement
		}
	}

	return &resp, nil
}
