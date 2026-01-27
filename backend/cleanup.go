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
	log.Printf("INFO: cleanup job logs older than 2 months")
	deleteThreshold := time.Now().AddDate(0, -2, 0)
	dateStr := deleteThreshold.Format("2006-01-02")

	log.Printf("INFO: scan for job statuses to delete")
	var delCount int64
	if err := svc.DB.Table("job_statuses").Where("status=? and ended_at < ?", "finished", dateStr).Count(&delCount).Error; err != nil {
		log.Printf("ERROR: unable to get count of old jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if delCount == 0 {
		log.Printf("INFO: there are no jobs to delete")
		c.String(http.StatusOK, "no messages to delete")
		return
	}

	log.Printf("INFO: delete %d old jobs ", delCount)
	if err := svc.DB.Exec("DELETE from job_statuses where status=? and ended_at < ?", "finished", dateStr).Error; err != nil {
		log.Printf("ERROR: unable to delete finished jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d jobs deleted", delCount))
}
