package main

import (
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

type marcMetadata struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	ControlFields []controlField `xml:"controlfield"`
	DataFields    []dataField    `xml:"datafield"`
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
	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) doSirsiLookup(catKey, barcode string) (*sirsiResponse, error) {
	// prefer catkey over barcode
	qp := fmt.Sprintf("barcode=%s", barcode)
	if catKey != "" {
		re := regexp.MustCompile(`^u`)
		cKey := re.ReplaceAll([]byte(catKey), []byte(""))
		qp = fmt.Sprintf("ckey=%s", cKey)
	}

	log.Printf("INFO: lookup sirsi marc metadata with [%s]", qp)
	url := fmt.Sprintf("https://ils.lib.virginia.edu/uhtbin/getMarc?%s&type=xml", qp)
	rawResp, err := svc.getRequest(url)
	if err != nil {
		return nil, fmt.Errorf("getMarc failed %d: %s", err.StatusCode, err.Message)
	}

	var parsed marcMetadata
	parseErr := xml.Unmarshal(rawResp, &parsed)
	if parseErr != nil {
		return nil, parseErr
	}

	if len(parsed.ControlFields) == 0 && len(parsed.DataFields) == 0 {
		return nil, fmt.Errorf("no matches found in sirsi")
	}

	log.Printf("INFO: extract fields from raw marc response")
	resp := sirsiResponse{CatalogKey: catKey}

	// catkey is in 001 of control fields. find it first
	for _, cf := range parsed.ControlFields {
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
	subtitleRegex := regexp.MustCompile(`\s*\/$`)              // strip trailing /
	pubRegex := regexp.MustCompile(`(?:^\[|\]$|\.$|\]\.$|:$)`) // strip [] and trailing . or :
	for _, df := range parsed.DataFields {
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
					resp.Title = strings.TrimSpace(sf.Value)
				}
				if sf.Code == "b" {
					sub := strings.TrimSpace(sf.Value)
					sub = strings.TrimSpace(subtitleRegex.ReplaceAllString(sub, ""))
					resp.Title += fmt.Sprintf(" %s", sub)
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
		if err != nil {
			log.Printf("ERROR: unable to load CNE data: %s", dbErr.Error())
		} else {
			resp.UseRightName = cne.Name
			resp.UseRightURI = cne.URI
			resp.UseRightStatement = cne.Statement
		}
	}

	return &resp, nil
}
