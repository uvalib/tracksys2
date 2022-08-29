package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type invoice struct {
	ID          int64
	OrderID     int64
	DateInvoice time.Time
	DateFeePaid *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type order struct {
	ID                             int64           `json:"id"`
	OrderStatus                    string          `json:"status"`
	OrderTitle                     string          `json:"title"`
	DateDue                        time.Time       `json:"dateDue"`
	CustomerID                     uint            `json:"-"`
	Customer                       customer        `gorm:"foreignKey:CustomerID" json:"customer"`
	AgencyID                       uint            `json:"-"`
	Agency                         agency          `gorm:"foreignKey:AgencyID" json:"agency"`
	Fee                            sql.NullFloat64 `json:"fee"`
	Invoices                       []invoice       `gorm:"foreignKey:OrderID"  json:"invoices"`
	UnitCount                      int64           `gorm:"unitCount" json:"unitCount"`
	Email                          string          `json:"email"`
	StaffNotes                     string          `json:"staffNotes"`
	DateRequestSubmitted           time.Time       `json:"dateSubmitted"`
	DateOrderApproved              *time.Time      `json:"dateOrderApproved"`
	DateCustomerNotified           *time.Time      `json:"dateCustomerNotified"`
	DatePatronDeliverablesComplete *time.Time      `json:"datePatronDeliverablesComplete"`
	DateArchivingComplete          *time.Time      `json:"dateArchivingComplete"`
	DateFinalizationBegun          *time.Time      `json:"dateFinalizationBegun"`
	DateFeeEstimateSentToCustomer  *time.Time      `json:"dateFeeEstimateSent"`
	UpdatedAt                      time.Time       `json:"-"`
}

func (svc *serviceContext) getOrders(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	filter := c.Query("filter")
	if filter == "" {
		filter = "active"
	}
	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}
	sortField := sortBy
	if sortBy == "dateDue" {
		sortField = "date_due"
	} else if sortBy == "dateSubmitted" {
		sortField = "date_request_submitted"
	} else if sortBy == "title" {
		sortField = "order_title"
	} else if sortBy == "unitCount" {
		sortField = "unit_count"
	}
	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	log.Printf("INFO: get %d %s orders starting from offset %d order %s", pageSize, filter, startIndex, orderStr)

	filterQ := svc.DB.Table("orders")

	dateNow := time.Now().Format("2006-01-02")
	if filter == "active" {
		filterQ = filterQ.Where("order_status!=? and order_status!=?", "canceled", "completed")
	} else if filter == "await" {
		filterQ = filterQ.Where("order_status=? or order_status=?", "requested", "await_fee")
	} else if filter == "deferred" {
		filterQ = filterQ.Where("order_status=?", "deferred")
	} else if filter == "canceled" {
		filterQ = filterQ.Where("order_status=?", "canceled")
	} else if filter == "complete" {
		filterQ = filterQ.Where("order_status=?", "completed")
	} else if filter == "due_today" {
		filterQ = filterQ.Where("date_due=?", dateNow).
			Where("order_status!=?", "completed").Where("order_status!=?", "deferred").Where("order_status!=?", "canceled")
	} else if filter == "due_week" {
		dateWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		filterQ = filterQ.Where("date_due>=?", dateNow).Where("date_due<=?", dateWeek).
			Where("order_status!=?", "completed").Where("order_status!=?", "deferred").Where("order_status!=?", "canceled")
	} else if filter == "overdue" {
		filterQ = filterQ.Where("date_due<?", dateNow).
			Where("order_status!=?", "completed").Where("order_status!=?", "deferred").Where("order_status!=?", "canceled")
	}

	var total int64
	filterQ.Count(&total)

	var o []*order
	unitCnt := "(select count(*) from units where order_id=orders.id) as unit_count"
	err := filterQ.Debug().Preload("Customer").Preload("Agency").
		Select("*", unitCnt).
		Offset(startIndex).Limit(pageSize).Order(orderStr).Find(&o).Error
	if err != nil {
		log.Printf("ERROR: unable to get orders: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type resp struct {
		Jobs  []*order `json:"orders"`
		Total int64    `json:"total"`
	}
	out := resp{Jobs: o, Total: total}

	c.JSON(http.StatusOK, out)
}
