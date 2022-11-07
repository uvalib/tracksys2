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
const Version = "0.1.0"

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
	api := router.Group("/api")
	{
		api.GET("/config", svc.getConfig)

		api.GET("/customers", svc.getCustomers)
		api.POST("/customers", svc.addOrUpdateCustomer)

		api.GET("/masterfiles/:id", svc.getMasterFile)
		api.POST("/masterfiles/:id/update", svc.updateMasterFile)
		api.POST("/masterfiles/:id/tags", svc.addMasterFileTag)
		api.DELETE("/masterfiles/:id/tags", svc.removeMasterFileTag)

		api.GET("/metadata/sirsi", svc.lookupSirsiMetadata)
		api.GET("/metadata/archivesspace", svc.validateArchivesSpaceMetadata)
		api.GET("/metadata/:id", svc.getMetadata)
		api.POST("/metadata", svc.createMetadata)

		api.GET("/orders", svc.getOrders)
		api.POST("/orders", svc.createOrder)
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
	log.Printf("INFO: start TrackSys2 on port %s with CORS support enabled", portStr)
	log.Fatal(router.Run(portStr))
}
