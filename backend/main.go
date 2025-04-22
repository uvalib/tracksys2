package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const Version = "1.30.1"

func main() {
	// Load cfg
	log.Printf("===> TrackSys2 is starting up <===")
	cfg := getConfiguration()
	svc := initializeService(Version, cfg)

	// Set routes and start server
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Set routes and start server
	router.GET("/version", svc.getVersion)
	router.GET("/healthcheck", svc.healthCheck)
	router.GET("/authenticate", svc.authenticate)
	router.GET("/config", svc.getConfig)
	router.POST("/cleanup", svc.cleanupExpiredData)
	router.POST("/upload_search_image", svc.uploadSearchImage)
	router.GET("/pdf", svc.downloadPDF)

	router.GET("/script", svc.scriptRunner)

	api := router.Group("/api", svc.authMiddleware)
	{
		api.POST("/agency", svc.addAgency)

		api.GET("/aptrust", svc.getAPTrustSubmissions)
		api.GET("/archivesspace", svc.getArchivesSpaceReviews)
		api.GET("/hathitrust", svc.getHathiTrustSubmissions)
		api.PUT("/hathitrust", svc.updateHathiTrustSubmissions)

		api.POST("/collection-facet", svc.addCollectionFacet)
		api.GET("/collections", svc.getCollections)
		api.GET("/collections/candidates", svc.findCollectionCandidates)
		api.GET("/collections/:id", svc.getCollectionItems)
		api.GET("/collections/:id/csv", svc.exportCollectionCSV)
		api.DELETE("/collections/:id/items/:item", svc.removeCollectionItem)
		api.POST("/collections/:id/item", svc.addCollectionItem)
		api.GET("/collections/:id/aptrust", svc.getAPTrustCollectionStatus)

		api.GET("/components/:id", svc.getComponentTree)
		api.GET("/components/:id/masterfiles", svc.getComponentMasterFiles)

		api.GET("/customers", svc.getCustomers)
		api.POST("/customers", svc.addOrUpdateCustomer)

		api.GET("/dashboard", svc.getDashboardStats)

		api.GET("/equipment", svc.getEquipment)
		api.POST("/equipment", svc.createEquipment)
		api.POST("/equipment/:id/update", svc.updateEquipment)
		api.POST("/workstation", svc.createWorkstation)
		api.POST("/workstation/:id/update", svc.updateWorkstation)
		api.POST("/workstation/:id/setup", svc.updateWorkstationSetup)

		api.GET("/locations/:id/units", svc.getLocationUnits)

		api.GET("/masterfiles/:id", svc.getMasterFile)
		api.POST("/masterfiles/:id/update", svc.updateMasterFile)
		api.POST("/masterfiles/:id/tags", svc.addMasterFileTag)
		api.DELETE("/masterfiles/:id/tags", svc.removeMasterFileTag)

		api.GET("/metadata/sirsi", svc.lookupSirsiMetadata)
		api.GET("/metadata/archivesspace", svc.validateArchivesSpaceMetadata)
		api.GET("/metadata/:id", svc.getMetadata)
		api.POST("/metadata/:id", svc.updateMetadata)
		api.DELETE("/metadata/:id", svc.deleteMetadata)
		api.POST("/metadata/:id/hathitrust", svc.updateHathiTrustStatus)
		api.POST("/metadata/:id/xml", svc.uploadXMLMetadata)
		api.GET("/metadata/:id/xml", svc.getXMLMetadata)
		api.POST("/metadata", svc.createMetadata)
		api.GET("/metadata/:id/aptrust", svc.getAPTrustMetadataStatus)

		api.POST("/metadata/:id/archivesspace", svc.requestArchivesSpaceReview)
		api.POST("/metadata/:id/archivesspace/review", svc.beginArchivesSpaceReview)
		api.POST("/metadata/:id/archivesspace/resubmit", svc.resubmitArchivesSpaceReview)
		api.POST("/metadata/:id/archivesspace/publish", svc.publishArchivesSpace)
		api.POST("/metadata/:id/archivesspace/reject", svc.rejectArchivesSpaceSubmission)
		api.DELETE("/metadata/:id/archivesspace", svc.cancelArchivesSpaceSubmission)
		api.POST("/metadata/:id/archivesspace/notes", svc.updateArchivesSpaceSubmissionNotes)

		api.GET("/orders", svc.getOrders)
		api.POST("/orders", svc.createOrder)
		api.DELETE("/orders/:id", svc.deleteOrder)
		api.GET("/orders/:id", svc.getOrderDetails)
		api.DELETE("/orders/:id/items/:item", svc.deleteOrderItem)
		api.POST("/orders/:id/units", svc.addUnitToOrder)
		api.POST("/orders/:id/update", svc.updateOrder)
		api.POST("/orders/:id/fee/waive", svc.waiveFee)
		api.POST("/orders/:id/fee/accept", svc.acceptFee)
		api.POST("/orders/:id/fee/decline", svc.declineFee)
		api.POST("/orders/:id/approve", svc.approveOrder)
		api.POST("/orders/:id/defer", svc.deferOrder)
		api.POST("/orders/:id/resume", svc.resumeOrder)
		api.POST("/orders/:id/cancel", svc.cancelOrder)
		api.POST("/orders/:id/complete", svc.completeOrder)
		api.POST("/orders/:id/processor", svc.setOrderProcessor)
		api.POST("/invoices/:id/update", svc.updateInvoice)

		api.GET("/published/dpla", svc.getPublishedDPLA)
		api.GET("/published/virgo", svc.getPublishedVirgo)
		api.GET("/published/archivesspace", svc.getPublishedArchivesSpace)

		api.GET("/jobs", svc.getJobStatuses)
		api.DELETE("/jobs", svc.deleteJobStatuses)
		api.GET("/jobs/:id", svc.getJobDetails)

		api.GET("/units/:id", svc.getUnit)
		api.DELETE("/units/:id", svc.deleteUnit)
		api.GET("/units/:id/pdf", svc.requestPDF)
		api.GET("/units/:id/pdf/status", svc.getPDFStatus)
		api.GET("/units/:id/exists", svc.validateUnit)
		api.POST("/units/:id/exemplar/:mfid", svc.setExemplar)
		api.POST("/units/:id/project", svc.createProject)
		api.GET("/units/:id/masterfiles", svc.getUnitMasterfiles)
		api.GET("/units/:id/clone-sources", svc.getUnitCloneSources)
		api.POST("/units/:id/update", svc.updateUnit)
		api.GET("/units/:id/csv", svc.exportUnitCSV)

		api.GET("/search", svc.searchRequest)
		api.GET("/search/images", svc.imageSearchRequest)

		api.GET("/staff", svc.getStaff)
		api.POST("/staff", svc.addOrUpdateStaff)

		api.GET("/tags", svc.getTags)
		api.POST("/tags", svc.createTag)

	}

	// Note: in dev mode, this is never actually used. The front end is served
	// by yarn and it proxies all requests to the API to the routes above
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", cfg.port)
	versionMap := svc.lookupVersion()
	versionStr := fmt.Sprintf("%s-%s", versionMap["version"], versionMap["build"])
	log.Printf("INFO: start TrackSys2 v%s on port %s with CORS support enabled", versionStr, portStr)
	log.Fatal(router.Run(portStr))
}

func (svc *serviceContext) scriptRunner(c *gin.Context) {
	log.Printf("INFO: script runner called")
	c.String(http.StatusNotImplemented, "no script is available")
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
