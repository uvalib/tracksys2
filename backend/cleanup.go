package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (svc *serviceContext) cleanupCanceledUnits(c *gin.Context) {
	log.Printf("INFO: cleanup canceled units")
	c.String(http.StatusNotImplemented, "not implemented")
}

func (svc *serviceContext) cleanupCanceledOrders(c *gin.Context) {
	log.Printf("INFO: cleanup canceled orders")
	c.String(http.StatusNotImplemented, "not implemented")
}

func (svc *serviceContext) cleanupExpiredJobLogs(c *gin.Context) {
	log.Printf("INFO: cleanup job logs and deleted messages older than 2 months")
	lastMonth := time.Now().AddDate(0, -2, 0)

	log.Printf("INFO: scan for job statuses to delete")
	var oldStatuses []jobStatus
	err := svc.DB.Where("status=? and ended_at < ?", "finished", lastMonth).Find(&oldStatuses).Error
	if err != nil {
		log.Printf("ERROR: unable to get count of old jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if len(oldStatuses) > 0 {
		log.Printf("INFO: delete %d expired jobs", len(oldStatuses))
		jsIDs := make([]uint64, 0)
		for _, js := range oldStatuses {
			jsIDs = append(jsIDs, js.ID)
		}
		err = svc.DB.Where("id in ?", jsIDs).Delete(&jobStatus{}).Error
		if err != nil {
			log.Printf("ERROR: unable to delete old jobs: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("%d jobs deleted", len(oldStatuses)))
}
