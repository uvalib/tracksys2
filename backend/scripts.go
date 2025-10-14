package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type scriptRequest struct {
	Script string `json:"script"`
	JWT    string `json:"jwt"`
	Params any    `json:"params"`
}

//	curl -X POST http://localhost:8085 -H "Content-Type: application/json" \
//	     --data '{"script": "deliver", "jwt": "token", "params": 7735}'
func (svc *serviceContext) scriptRunner(c *gin.Context) {
	log.Printf("INFO: script runner called")

	var req scriptRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid script request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: validating jwt...")
	jwtClaims := jwtClaims{}
	_, jwtErr := jwt.ParseWithClaims(req.JWT, &jwtClaims, func(token *jwt.Token) (any, error) {
		return []byte(svc.JWTKey), nil
	})

	if jwtErr != nil {
		log.Printf("ERROR: authentication failed; token validation failed: %+v", jwtErr)
		c.String(http.StatusForbidden, "forbidden")
		return
	}

	if jwtClaims.Role != "admin" {
		log.Printf("ERROR: authentication failed; running scripts requires admin permissions: %+v", jwtErr)
		c.String(http.StatusForbidden, "forbidden")
		return
	}

	log.Printf("INFO: requested execution of script [%s]", req.Script)
	var scriptErr error
	switch req.Script {
	case "deliver":
		scriptErr = svc.deliverOrderMasterfiles(req.Params)
	default:
		scriptErr = fmt.Errorf("script %s is unknown", req.Script)
	}

	if scriptErr != nil {
		log.Printf("ERROR: script request failed: %s", scriptErr.Error())
		c.String(http.StatusBadRequest, scriptErr.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%s success", req.Script))
}

func (svc *serviceContext) deliverOrderMasterfiles(params any) error {
	log.Printf("INFO: script [deliver] started with params %v", params)
	log.Printf("TYPE OF PARAMS %s", reflect.TypeOf(params))
	fOrder, ok := params.(float64)
	if !ok {
		return fmt.Errorf("invalid deliver script params %v", params)
	}
	orderID := int64(fOrder)
	log.Printf("Deliver master files from order %d", orderID)
	var units []unit
	err := svc.DB.Where("order_id=?", orderID).Preload("Metadata").Find(&units).Error
	if err != nil {
		return err
	}

	csvF, _ := os.Create(fmt.Sprintf("order_%d.csv", orderID))
	csvW := csv.NewWriter(csvF)
	defer csvF.Close()
	csvHead := []string{"order_id", "unit_id", "title", "image_count", "download_url"}
	csvW.Write(csvHead)

	type deliverReq struct {
		Filename  string `json:"filename"`
		ComputeID string `json:"computeID"`
		Deliver   bool   `json:"deliver"`
	}

	cnt := 0
	for _, u := range units {
		log.Printf("INFO: script [deliver] process unit #%d id: %d", (cnt + 1), u.ID)
		cnt++
		var mfCount int
		svc.DB.Raw("select count(*) from master_files where unit_id=?", u.ID).Scan(&mfCount)
		url := fmt.Sprintf("https://digiservdelivery.lib.virginia.edu/order_%d/%d.zip", u.OrderID, u.ID)
		row := []string{fmt.Sprintf("%d", u.OrderID), fmt.Sprintf("%d", u.ID), u.Metadata.Title, fmt.Sprintf("%d", mfCount), url}
		csvW.Write(row)

		postURL := fmt.Sprintf("https://dpg-jobs.lib.virginia.edu/units/%d/copy", u.ID)
		payload := deliverReq{Filename: "all", ComputeID: "lf6f", Deliver: true}
		resp, pErr := svc.postJSON(postURL, payload)
		if pErr != nil {
			log.Printf("ERROR: delivery post for unit %d failed %s", u.ID, pErr.Message)
			continue
		}

		if string(resp) == "" {
			log.Printf("unit %d required no processing", u.ID)
			continue
		}
		jobID, parseErr := strconv.ParseInt(string(resp), 10, 64)
		if parseErr != nil {
			log.Printf("ERROR: unable to parse unit %d delive job id from response %s: %s", u.ID, resp, parseErr.Error())
			continue
		}

		log.Printf("INFO: script [deliver] unit %d delivery has been requested and has job %d; wait for it to complete", u.ID, jobID)
		svc.awaitJobCompletion(jobID)
	}

	csvW.Flush()
	log.Printf("INFO: script [deliver] completed with %d units processed", cnt)
	return nil
}

// SCRIPT TO VERIFY BARCODES ARE SUBMITTED TO HATHITRUST AND ARE DISCOVERABLE
// func (svc *serviceContext) scriptRunner(c *gin.Context) {
// 	log.Printf("INFO: script runner called")
// 	bytes, err := os.ReadFile("./verified_missing.txt")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	count := 0
// 	missing := make([]string, 0)
// 	templateMARC := "https://ils.lib.virginia.edu/uhtbin/getMarc?barcode=BC&hathitrust=yes&type=xml"
// 	templateURL := "https://babel.hathitrust.org/cgi/ls?q1=sdr-uva.KEY&field1=ocr&a=srchls&ft=ft&lmt=all"
// 	for _, bc := range strings.Split(string(bytes), "\n") {
// 		cleanBC := strings.TrimSpace(bc)
// 		if cleanBC == "" {
// 			continue
// 		}
// 		count++

// 		// get the MARC data from Sirsi
// 		url := strings.ReplaceAll(templateMARC, "BC", cleanBC)
// 		resp, reqErr := svc.getRequest(url)
// 		if reqErr != nil {
// 			log.Fatal(fmt.Sprintf("%d: %s", reqErr.StatusCode, reqErr.Message))
// 		}

// 		// Parse cat key from the first control field
// 		// find this: <controlfield tag="001">u500878</controlfield>
// 		respStr := string(resp)
// 		cfIdx := strings.Index(respStr, "<controlfield")
// 		if cfIdx == -1 {
// 			log.Fatal("<controlfield not found")
// 		}
// 		trimmed := respStr[cfIdx+24:]
// 		endIdx := strings.Index(trimmed, "<")
// 		catKey := trimmed[:endIdx]

// 		// search HT using specially formatted cat key
// 		htURL := strings.ReplaceAll(templateURL, "KEY", catKey)
// 		resp, reqErr = svc.getRequest(htURL)
// 		if reqErr != nil {
// 			log.Fatal(fmt.Sprintf("%d: %s", reqErr.StatusCode, reqErr.Message))
// 		}

// 		if strings.Contains(string(resp), "No results") {
// 			log.Printf("INFO: ..... NOT FOUND")
// 			missing = append(missing, cleanBC)
// 		} else {
// 			log.Printf("INFO: ..... OK")
// 		}
// 	}

// 	log.Printf("INFO: MISSING %s", strings.Join(missing, ","))
// 	c.String(http.StatusOK, fmt.Sprintf("%d processed, %d missing", count, len(missing)))
// }

// SCRIPT TO ADD TO A COLLECTION BASED ON ARCHIVESSPACE URI
// func (svc *serviceContext) scriptRunner(c *gin.Context) {
// 	log.Printf("INFO: script runner called")

// 	// //randolph_nicholas
// 	// asCollectionID := 1395
// 	// tsCollectionID := 108160

// 	// randolph
// 	asCollectionID := 1426
// 	tsCollectionID := 108159

// 	var memberURIs []string
// 	memberSQL := "select external_uri from metadata where parent_metadata_id=?"
// 	err := svc.DB.Raw(memberSQL, tsCollectionID).Scan(&memberURIs).Error
// 	if err != nil {
// 		log.Fatal(err.Error())
// 		return
// 	}
// 	log.Printf("INFO: collection member uri list: %v", memberURIs)

// 	bytes, reqErr := svc.getRequest(fmt.Sprintf("%s/archivesspace/collections/%d/records", svc.ExternalSystems.Jobs, asCollectionID))
// 	if reqErr != nil {
// 		log.Fatal(fmt.Sprintf("%d:%s", reqErr.StatusCode, reqErr.Message))
// 	}

// 	uris := make([]string, 0)
// 	processedCount := 0
// 	for _, extURI := range strings.Split(string(bytes), "\n") {
// 		processedCount++
// 		alreadyMember := false
// 		for _, memberURI := range memberURIs {
// 			if memberURI == extURI {
// 				alreadyMember = true
// 				break
// 			}
// 		}

// 		if alreadyMember == false {
// 			uris = append(uris, extURI)
// 			log.Printf("INFO: %s added to the pending member list", extURI)
// 			if len(uris) == 50 {
// 				log.Printf("INFO: process chunk of %d items to add to collection", len(uris))
// 				sql := "update metadata set parent_metadata_id=? where parent_metadata_id=? and external_system_id=? and external_uri in ?"
// 				err := svc.DB.Exec(sql, tsCollectionID, 0, 1, uris).Error
// 				if err != nil {
// 					log.Fatal(err.Error())
// 					break
// 				}
// 				uris = make([]string, 0)
// 			}
// 		} else {
// 			log.Printf("INFO: %s is already part iof the colelction; skipping", extURI)
// 		}
// 	}

// 	if len(uris) > 0 {
// 		log.Printf("INFO: process final chunk of %d items to add to collection", len(uris))
// 		sql := "update metadata set parent_metadata_id=? where parent_metadata_id=? and external_system_id=? and external_uri in ?"
// 		err := svc.DB.Debug().Exec(sql, tsCollectionID, 0, 1, uris).Error
// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}
// 	}

// 	log.Printf("INFO: done. %d uris processed", processedCount)

// 	c.String(http.StatusOK, fmt.Sprintf("%d uris processed", processedCount))
// }

// SAMPLE SCRIPT TO BULK ADD ITEMS TO A COLLECTION
// func (svc *serviceContext) scriptRunner(c *gin.Context) {
// 	log.Printf("INFO: script runner called")
// 	bytes, err := os.ReadFile("./vickery_units.txt")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	var collectionMD metadata
// 	err = svc.DB.First(&collectionMD, 108065).Error
// 	if err != nil {
// 		log.Fatal("cannot find collection record: " + err.Error())
// 	}

// 	cnt := 0
// 	for _, unitStr := range strings.Split(string(bytes), "\n") {
// 		unitID, _ := strconv.ParseInt(unitStr, 10, 64)
// 		if unitID == 0 {
// 			continue
// 		}
// 		log.Printf("INFO: add %d to collection", unitID)
// 		var tgtUnit unit
// 		err = svc.DB.Preload("Metadata").First(&tgtUnit, unitID).Error
// 		if err != nil {
// 			log.Printf("ERROR: unable to get unit %d: %s", unitID, err.Error())
// 			break
// 		}

// 		if tgtUnit.Metadata.ParentMetadataID != 0 {
// 			log.Printf("WARNING: unit %d metadata %s already has parent metadata %d; skipping", unitID, tgtUnit.Metadata.PID, tgtUnit.Metadata.ParentMetadataID)
// 			continue
// 		}

// 		tgtUnit.Metadata.ParentMetadataID = collectionMD.ID
// 		err = svc.DB.Model(tgtUnit.Metadata).Update("parent_metadata_id", collectionMD.ID).Error
// 		if err != nil {
// 			log.Printf("ERROR: unable to update metadata %d parent: %s", tgtUnit.Metadata.ID, err.Error())
// 			break
// 		}

// 		cnt++
// 	}

// 	c.String(http.StatusOK, fmt.Sprintf("%d records added to collection 108065", cnt))
// }

// SCRIPT TO FLAG ORDER METADATA FOR HATHITRUST PUBLISH
// log.Printf("INFO: script runner called")
// orderStr := c.Query("order")
// tgtID, _ := strconv.ParseInt(orderStr, 10, 64)

// err := svc.flagOrderForHathTrust(tgtID)
// if err != nil {
// 	log.Printf("ERROR: publish order to hathitrust failed: %s", err.Error())
// 	c.String(http.StatusInternalServerError, err.Error())
// 	return
// }

// c.String(http.StatusOK, fmt.Sprintf("flagged order %d", tgtID))

// SAMPLE republish AS records
// log.Printf("INFO: republish AS records from 4/17 and 4/18")
// var asRecs []metadata
// err := svc.DB.Where("type=? and external_system_id=? and date_dl_ingest like ?", "ExternalMetadata", 1, "2023-04-18%").Find(&asRecs).Error
// if err != nil {
// 	log.Printf("ERROR: unable to find recently published as records: %s", err.Error())
// 	c.String(http.StatusInternalServerError, err.Error())
// 	return
// }

// published := 0
// errors := 0
// for _, md := range asRecs {
// 	log.Printf("INFO: publish %s to AS", md.PID)
// 	payload := struct {
// 		UserID     string `json:"userID"`
// 		MetadataID string `json:"metadataID"`
// 	}{
// 		UserID:     "120", /// lf6f
// 		MetadataID: fmt.Sprintf("%d", md.ID),
// 	}
// 	url := fmt.Sprintf("%s/archivesspace/publish?immediate=1", svc.ExternalSystems.Jobs)
// 	asResp, asErr := svc.postJSON(url, payload)
// 	if asErr != nil {
// 		log.Printf("ERROR: AS publish %d failed %d: %s", md.ID, asErr.StatusCode, asErr.Message)
// 		errors++
// 	} else {
// 		log.Printf("INFO: %s", asResp)
// 		published++
// 	}
// }

// c.String(http.StatusOK, fmt.Sprintf("DONE. %d total records, %d published, %d failed", len(asRecs), published, errors))

// SAMPLE for updating AS date_dl_ingested
//
// log.Printf("INFO: update date_dl_ingest for published AS metadata records")
// var asMD []metadata
// err := svc.DB.Where("external_system_id=? and date_dl_ingest is null", 1).Find(&asMD).Error
// if err != nil {
// 	log.Printf("ERROR: unable to get AS metadata records: %s", err.Error())
// 	c.String(http.StatusInternalServerError, err.Error())
// 	return
// }

// cnt := 0
// errors := 0
// for _, md := range asMD {
// 	log.Printf("INFO: update %s: %s", md.PID, *md.ExternalURI)
// 	raw, getErr := svc.getRequest(fmt.Sprintf("%s/archivesspace/lookup?pid=%s&uri=%s", svc.ExternalSystems.Jobs, md.PID, *md.ExternalURI))
// 	if getErr != nil {
// 		log.Printf("ERROR: unable to get archivesSpace metadata for %s: %s", md.PID, getErr.Message)
// 	} else {
// 		var parsed asMetadata
// 		err := json.Unmarshal(raw, &parsed)
// 		if err != nil {
// 			log.Printf("ERROR: unable to parse AS response for %s: %s", md.PID, err.Error())
// 			errors++
// 		} else {
// 			if parsed.PublishedAt != "" {
// 				parsedDate, err := time.Parse("2006-01-02T15:04:05Z", parsed.PublishedAt)
// 				if err != nil {
// 					log.Printf("ERROR: unable to parse date %s: %s", parsed.PublishedAt, err.Error())
// 					errors++
// 				} else {
// 					log.Printf("INFO: set AS published date: %v", parsedDate)
// 					md.DateDLIngest = &parsedDate
// 					err := svc.DB.Model(&md).Select("DateDLIngest").Updates(md).Error
// 					if err != nil {
// 						log.Printf("ERROR: unable to update %s: %s", md.PID, err.Error())
// 						errors++
// 					} else {
// 						cnt++
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// log.Printf("INFO: %d records updated", len(asMD))
// c.String(http.StatusOK, fmt.Sprintf("%d errors, %d records updated", errors, cnt))
