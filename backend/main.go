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
const Version = "1.8.1"

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

	router.GET("/script", svc.scriptRunner)

	api := router.Group("/api", svc.authMiddleware)
	{
		api.POST("/collection-facet", svc.addCollectionFacet)
		api.GET("/collections", svc.getCollections)
		api.GET("/collections/:id", svc.getCollectionItems)
		api.GET("/collections/:id/export", svc.exportCollectionItems)
		api.DELETE("/collections/:id/items/:item", svc.removeCollectionItem)
		api.POST("/collections/:id/items", svc.addCollectionItems)

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

		api.GET("/masterfiles/:id", svc.getMasterFile)
		api.POST("/masterfiles/:id/update", svc.updateMasterFile)
		api.POST("/masterfiles/:id/tags", svc.addMasterFileTag)
		api.DELETE("/masterfiles/:id/tags", svc.removeMasterFileTag)

		api.GET("/metadata/sirsi", svc.lookupSirsiMetadata)
		api.GET("/metadata/archivesspace", svc.validateArchivesSpaceMetadata)
		api.GET("/metadata/:id", svc.getMetadata)
		api.POST("/metadata/:id", svc.updateMetadata)
		api.DELETE("/metadata/:id", svc.deleteMetadata)
		api.POST("/metadata/:id/xml", svc.uploadXMLMetadata)
		api.GET("/metadata/:id/xml", svc.getXMLMetadata)
		api.POST("/metadata", svc.createMetadata)

		api.GET("/orders", svc.getOrders)
		api.POST("/orders", svc.createOrder)
		api.DELETE("/orders/:id", svc.deleteOrder)
		api.GET("/orders/:id", svc.getOrderDetails)
		api.DELETE("/orders/:id/items/:item", svc.deleteOrderItem)
		api.POST("/orders/:id/units", svc.addUnitToOrder)
		api.POST("/orders/:id/update", svc.updateOrder)
		api.POST("/orders/:id/fee/accept", svc.acceptFee)
		api.POST("/orders/:id/fee/decline", svc.declineFee)
		api.POST("/orders/:id/approve", svc.approveOrder)
		api.POST("/orders/:id/defer", svc.deferOrder)
		api.POST("/orders/:id/resume", svc.resumeOrder)
		api.POST("/orders/:id/cancel", svc.cancelOrder)
		api.POST("/orders/:id/complete", svc.completeOrder)
		api.POST("/orders/:id/processor", svc.setOrderProcessor)
		api.POST("/invoices/:id/update", svc.updateInvoice)

		api.GET("/jobs", svc.getJobStatuses)
		api.DELETE("/jobs", svc.deleteJobStatuses)
		api.GET("/jobs/:id", svc.getJobDetails)

		api.GET("/units/:id", svc.getUnit)
		api.GET("/units/:id/pdf", svc.requestPDF)
		api.GET("/units/:id/pdf/status", svc.getPDFStatus)
		api.GET("/units/:id/pdf/download", svc.downloadPDF)
		api.GET("/units/:id/exists", svc.validateUnit)
		api.POST("/units/:id/exemplar/:mfid", svc.setExemplar)
		api.POST("/units/:id/project", svc.createProject)
		api.GET("/units/:id/masterfiles", svc.getUnitMasterfiles)
		api.GET("/units/:id/clone-sources", svc.getUnitCloneSources)
		api.POST("/units/:id/update", svc.updateUnit)

		api.GET("/search", svc.searchRequest)

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
