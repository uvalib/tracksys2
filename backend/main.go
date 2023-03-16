package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const Version = "1.3.9"

// func (svc *serviceContext) scriptHack(c *gin.Context) {
// 	f, err := os.Open("./data/as_records.csv")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	defer f.Close()
// 	csvReader := csv.NewReader(f)
// 	asRecs, err := csvReader.ReadAll()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	out, err := os.Create("./data/as_update.sql")
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	priorURI := ""
// 	priorID := ""
// 	dupPrinted := false
// 	for _, rec := range asRecs {
// 		if rec[2] == priorURI {
// 			if dupPrinted == false {
// 				log.Printf("DUPLICATE URI: %s : %s", priorID, priorURI)
// 				dupPrinted = true
// 			}
// 			log.Printf("DUPLICATE URI: %s : %s", rec[0], rec[2])
// 			// continue
// 		} else {
// 			priorURI = rec[2]
// 			priorID = rec[0]
// 			dupPrinted = false
// 		}
// 		tgtURL := fmt.Sprintf("https://dpg-jobs.lib.virginia.edu/archivesspace/lookup?uri=%s", rec[2])
// 		log.Printf("%s : %s", rec[0], tgtURL)

// 		rawDetail, getErr := svc.getRequest(tgtURL)
// 		if getErr != nil {
// 			log.Fatal(fmt.Sprintf("get AS data failed %d: %s", getErr.StatusCode, getErr.Message))
// 		}
// 		var asObj asMetadata
// 		parseErr := json.Unmarshal(rawDetail, &asObj)
// 		if parseErr != nil {
// 			log.Fatal(parseErr.Error())
// 		}

// 		newTitle := asObj.Title
// 		if asObj.Dates != "" {
// 			newTitle = fmt.Sprintf("%s, %s", asObj.Title, asObj.Dates)
// 		}
// 		sql := fmt.Sprintf("update metadata set title = \"%s\" where id = %s;\n", newTitle, rec[0])
// 		out.WriteString(sql)
// 	}
// }

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

	// router.GET("/script", svc.scriptHack)

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
