package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dashboardStats struct {
	DueInOneWeek            int64 `json:"dueInOneWeek"`
	Overdue                 int64 `json:"overdue"`
	ReadyForDelivery        int64 `json:"readyForDelivery"`
	ArchivesSpaceRequests   int64 `json:"asRequests"`
	ArchivesSpaceReviews    int64 `json:"asReviews"`
	ArchivesSpaceRejections int64 `json:"asRejections"`
}

func (svc *serviceContext) getDashboardStats(c *gin.Context) {
	log.Printf("INFO: looking up dashboard stats")

	var stats dashboardStats
	now := time.Now()
	inProcQ := svc.DB.Where("order_status != ? and order_status != ? and order_status != ?", "completed", "deferred", "canceled")

	// due in a week
	oneWeek := now.AddDate(0, 0, 7)
	baseQ := svc.DB.Table("orders").
		Joins("inner join units u on u.order_id=orders.id").
		Where("u.intended_use_id <> ?", 110).Where(inProcQ)
	err := baseQ.Where("date_due>=?", now.Format("2006-01-02")).
		Where("date_due<=?", oneWeek.Format("2006-01-02")).
		Distinct("orders.id").Count(&stats.DueInOneWeek).Error
	if err != nil {
		log.Printf("ERROR: unable to get orders due in a week: %s", err.Error())
	}

	// overdue
	oneYearAgo := now.AddDate(-1, 0, 0)
	baseQ = svc.DB.Table("orders").
		Joins("inner join units u on u.order_id=orders.id").
		Where("u.intended_use_id <> ?", 110).Where(inProcQ)
	err = baseQ.Where("date_request_submitted>?", oneYearAgo.Format("2006-01-02")).
		Where("date_due<?", now.Format("2006-01-02")).
		Distinct("orders.id").Count(&stats.Overdue).Error
	if err != nil {
		log.Printf("ERROR: unable to get overdue orders: %s", err.Error())
	}

	// ready for delivery
	err = svc.DB.Table("orders").Joins("inner join units u on u.order_id=orders.id").
		Where("u.intended_use_id <> ?", 110).
		Where("orders.email is not null and orders.email != ? and date_customer_notified is null", "").
		Where("order_status != ? and order_status != ?", "canceled", "completed").
		Distinct("orders.id").Count(&stats.ReadyForDelivery).Error
	if err != nil {
		log.Printf("ERROR: unable to get ready for delivery orders: %s", err.Error())
	}

	// archivesspace
	var asActive []archivesspaceReview
	err = svc.DB.Where("published_at is null").Find(&asActive).Error
	if err != nil {
		log.Printf("ERROR: unable to get active archivesspace review stats: %s", err.Error())
	} else {
		for _, asR := range asActive {
			switch asR.Status {
			case "rejected":
				stats.ArchivesSpaceRejections++
			case "requested":
				stats.ArchivesSpaceRequests++
			default:
				stats.ArchivesSpaceReviews++
			}
		}
	}

	c.JSON(http.StatusOK, stats)
}
