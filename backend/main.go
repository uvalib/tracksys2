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
const Version = "1.33.0"

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
	router.POST("/upload_search_image", svc.uploadSearchImage)
	router.GET("/pdf", svc.downloadPDF)

	router.POST("/script", svc.scriptRunner)

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

		// statistics reporting support
		api.GET("/stats/archive", svc.getArchiveStats)
		api.GET("/stats/images", svc.getImageStats)
		api.GET("/stats/metadata", svc.getMetadataStats)
		api.GET("/stats/published", svc.getPublishedStats)
		api.GET("/stats/storage", svc.getStorageStats)
		api.GET("/stats/deliveries", svc.getDeliveryStats)

		// master file audit report
		api.GET("/reports/audit", svc.getAuditReport)
	}

	cleanup := router.Group("/cleanup")
	{
		cleanup.POST("/expired-jobs", svc.cleanupExpiredJobLogs)
		cleanup.POST("/canceled-units", svc.cleanupCanceledUnits)
		cleanup.POST("/canceled-orders", svc.cleanupCanceledOrders)
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
