package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dashboardStats struct {
	DueToday         int64 `json:"dueToday"`
	DueInOneWeek     int64 `json:"dueInOneWeek"`
	Overdue          int64 `json:"overdue"`
	ReadyForDelivery int64 `json:"readyForDelivery"`
}

func (svc *serviceContext) getDashboardStats(c *gin.Context) {
	log.Printf("INFO: looking up dashboard stastics")

	var stats dashboardStats
	now := time.Now()
	inProcQ := svc.DB.Where("order_status != ? and order_status != ? and order_status != ?", "completed", "deferred", "canceled")
	baseQ := svc.DB.Table("orders").
		Joins("inner join units u on u.order_id=orders.id").
		Where("u.intended_use_id <> ?", 110).Where(inProcQ)

	// due today
	err := baseQ.Debug().Where("date_due=?", now.Format("2006-01-02")).Distinct("orders.id").Count(&stats.DueToday).Error
	if err != nil {
		log.Printf("ERROR: unable to get orders due today: %s", err.Error())
	}

	/// due in a week
	oneWeek := now.AddDate(0, 0, 7)
	err = baseQ.Debug().Where("date_due>?", now.Format("2006-01-02")).
		Where("date_due<=?", oneWeek.Format("2006-01-02")).
		Distinct("orders.id").Count(&stats.DueInOneWeek).Error
	if err != nil {
		log.Printf("ERROR: unable to get orders due in a week: %s", err.Error())
	}

	// overdue
	oneYearAgo := now.AddDate(-1, 0, 0)
	err = baseQ.Debug().Where("date_request_submitted>?", oneYearAgo.Format("2006-01-02")).
		Where("date_due<?", now.Format("2006-01-02")).
		Distinct("orders.id").Count(&stats.Overdue).Error
	if err != nil {
		log.Printf("ERROR: unable to get overdue orders: %s", err.Error())
	}

	// ready for delivery
	err = svc.DB.Debug().Table("orders").Joins("inner join units u on u.order_id=orders.id").
		Where("u.intended_use_id <> ?", 110).
		Where("orders.email is not null and date_customer_notified is null and order_status=?", "approved").
		// Where("order_status != ? and order_status != ?", "canceled", "completed").
		Distinct("orders.id").Count(&stats.ReadyForDelivery).Error
	if err != nil {
		log.Printf("ERROR: unable to get ready for delivery orders: %s", err.Error())
	}

	c.JSON(http.StatusOK, stats)
}
