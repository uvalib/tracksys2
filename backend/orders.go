package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type invoice struct {
	ID                   int64      `json:"id"`
	OrderID              int64      `json:"-"`
	DateInvoice          time.Time  `json:"invoiceDate"`
	DateFeePaid          *time.Time `json:"dateFeePaid,omitempty"`
	DateSecondNoticeSent *time.Time `json:"dateNoticeSent,omitempty"`
	DateFeeDeclined      *time.Time `json:"dateFeeDeclined,omitempty"`
	FeeAmountPaid        int64      `json:"feeAmountPaid"`
	TransmittalNumber    string     `json:"transmittalNumber"`
	Notes                string     `json:"notes"`
	PermanentNonpayment  bool       `json:"permanentNonpayment"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"-"`
}

type orderItem struct {
	ID            int64        `json:"id"`
	OrderID       int64        `json:"-"`
	IntendedUseID *int64       `json:"-"`
	IntendedUse   *intendedUse `gorm:"foreignKey:IntendedUseID" json:"intendedUse"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Pages         string       `json:"pages"`
	CallNumber    string       `json:"callNumber"`
	Author        string       `json:"author"`
	Year          string       `json:"year"`
	Location      string       `json:"location"`
	SourceURL     string       `gorm:"column:source_url" json:"sourceURL"`
	Converted     bool         `json:"converted"`
}

type auditEvent struct {
	ID            int64       `json:"id"`
	StaffMemberID int64       `json:"-"`
	StaffMember   staffMember `gorm:"foreignKey:StaffMemberID" json:"staffMember"`
	AuditableID   int64       `json:"-"`
	AuditableType string      `json:"-"`
	Details       string      `json:"details"`
	CreatedAt     time.Time   `json:"createdAt"`
}

type order struct {
	ID                             int64      `json:"id"`
	OrderStatus                    string     `json:"status"`
	OrderTitle                     string     `json:"title"`
	DateDue                        time.Time  `json:"dateDue"`
	CustomerID                     *uint      `json:"-"`
	Customer                       *customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	AgencyID                       *uint      `json:"-"`
	Agency                         *agency    `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Fee                            *float64   `json:"fee,omitempty"`
	Invoice                        *invoice   `gorm:"-" json:"invoice,omitempty"`
	UnitCount                      int64      `gorm:"unitCount" json:"unitCount"`
	MasterFileCount                int64      `gorm:"masterFileCount" json:"masterFileCount"`
	Email                          string     `json:"email"`
	StaffNotes                     string     `json:"staffNotes"`
	SpecialInstructions            string     `json:"specialInstructions"`
	DateRequestSubmitted           time.Time  `json:"dateSubmitted"`
	DateOrderApproved              *time.Time `json:"dateApproved"`
	DateDeferred                   *time.Time `json:"dateDeferred"`
	DateCanceled                   *time.Time `json:"dateCanceled"`
	DateCustomerNotified           *time.Time `json:"dateCustomerNotified"`
	DatePatronDeliverablesComplete *time.Time `json:"datePatronDeliverablesComplete"`
	DateArchivingComplete          *time.Time `json:"dateArchivingComplete"`
	DateFinalizationBegun          *time.Time `json:"dateFinalizationBegun"`
	DateFeeEstimateSentToCustomer  *time.Time `json:"dateFeeEstimateSent"`
	DateCompleted                  *time.Time `json:"dateCompleted"`
	UpdatedAt                      time.Time  `json:"-"`
}

func (svc *serviceContext) getOrderDetails(c *gin.Context) {
	oID := c.Param("id")
	log.Printf("INFO: get order %s details", oID)
	var oDetail order
	err := svc.DB.Preload("Agency").Preload("Customer").Find(&oDetail, oID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: lookup invoice for order %d", oDetail.ID)
	var invDetail invoice
	err = svc.DB.Where("order_id=?", oID).Order("created_at desc").First(&invDetail).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: unable to get invoice for order %s: %s", oID, err.Error())
		} else {
			log.Printf("INFO: no invoice for order %s", oID)
		}
	} else {
		oDetail.Invoice = &invDetail
	}

	type oResp struct {
		Order  order        `json:"order"`
		Units  []unit       `json:"units"`
		Items  []orderItem  `json:"items"`
		Events []auditEvent `json:"events"`
	}
	out := oResp{Order: oDetail}
	err = svc.DB.Where("order_id=?", oID).Preload("IntendedUse").Find(&out.Units).Error
	if err != nil {
		log.Printf("ERROR: unable to get units for order %s: %s", oID, err.Error())
	}
	err = svc.DB.Where("order_id=?", oID).Preload("IntendedUse").Find(&out.Items).Error
	if err != nil {
		log.Printf("ERROR: unable to get items for order %s: %s", oID, err.Error())
	}
	err = svc.DB.Where("auditable_type=?", "Order").Where("auditable_id=?", oDetail.ID).Preload("StaffMember").Find(&out.Events).Error
	if err != nil {
		log.Printf("ERROR: unable to get audit events for order %s: %s", oID, err.Error())
	}

	c.JSON(http.StatusOK, out)
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
	sortField := fmt.Sprintf("orders.%s", sortBy)
	if sortBy == "dateDue" {
		sortField = "date_due"
	} else if sortBy == "dateSubmitted" {
		sortField = "date_request_submitted"
	} else if sortBy == "title" {
		sortField = "order_title"
	} else if sortBy == "unitCount" {
		sortField = "unit_count"
	} else if sortBy == "masterFileCount" {
		sortField = "master_file_count"
	}
	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	log.Printf("INFO: get %d %s orders starting from offset %d order %s", pageSize, filter, startIndex, orderStr)

	// set up filtering....
	filterQ := svc.DB.Table("orders").Joins("inner join customers c on c.id=orders.customer_id").Joins("left outer join agencies a on a.id = orders.agency_id")
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

	// set up query...
	queryStr := c.Query("q")
	var qObj *gorm.DB
	if queryStr != "" {
		queryAny := fmt.Sprintf("%%%s%%", queryStr)
		queryStart := fmt.Sprintf("%s%%", queryStr)
		qObj = svc.DB.Where("order_title like ?", queryAny).Or("staff_notes like ?", queryAny).
			Or("special_instructions like ?", queryAny).Or("c.last_name like ?", queryStart).Or("a.name like ?", queryStart)
		filterQ = filterQ.Where(qObj)
	}

	var total int64
	filterQ.Count(&total)

	var o []*order
	unitCnt := "(select count(*) from units where order_id=orders.id) as unit_count"
	mfCnt := "(select count(*) from master_files m inner join units u on u.id=m.unit_id where u.order_id=orders.id) as master_file_count"
	err := filterQ.Preload("Agency").Preload("Customer").
		Select("orders.*", unitCnt, mfCnt).
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

func (svc *serviceContext) updateOrder(c *gin.Context) {
	orderID := c.Param("id")
	log.Printf("INFO: update order %s", orderID)
	var oDetail order
	err := svc.DB.Preload("Agency").Preload("Customer").Find(&oDetail, orderID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", orderID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var updateRequest struct {
		Status              string  `json:"status"`
		DateDue             string  `json:"dateDue"`
		Title               string  `json:"title"`
		SpecialInstructions string  `json:"specialInstructions"`
		StaffNotes          string  `json:"staffNotes"`
		Fee                 *string `json:"fee"`
		AgencyID            uint    `json:"agencyID"`
		CustomerID          uint    `json:"customerID"`
	}
	err = c.BindJSON(&updateRequest)
	if err != nil {
		log.Printf("ERROR: invalid update order %s request: %s", orderID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	oDetail.OrderStatus = updateRequest.Status
	oDetail.DateDue, _ = time.Parse("2006-01-02", updateRequest.DateDue)
	oDetail.OrderTitle = updateRequest.Title
	oDetail.SpecialInstructions = updateRequest.SpecialInstructions
	oDetail.StaffNotes = updateRequest.StaffNotes
	oDetail.Fee = nil
	if updateRequest.Fee != nil {
		floatFee, _ := strconv.ParseFloat(*updateRequest.Fee, 64)
		if floatFee > 0 {
			oDetail.Fee = &floatFee
		}
	}
	oDetail.AgencyID = nil
	if updateRequest.AgencyID != 0 {
		svc.DB.Find(&oDetail.Agency, updateRequest.AgencyID)
		oDetail.AgencyID = &updateRequest.AgencyID
	}
	oDetail.CustomerID = nil
	if updateRequest.CustomerID != 0 {
		svc.DB.Find(&oDetail.Customer, updateRequest.CustomerID)
		oDetail.CustomerID = &updateRequest.CustomerID
	}

	err = svc.DB.Model(&oDetail).Select("OrderStatus", "DateDue", "OrderTitle", "SpecialInstructions", "StaffNotes", "Fee", "AgencyID", "CustomerID").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to update order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: order %d updated", oDetail.ID)
	svc.DB.Preload("Agency").Preload("Customer").Find(&oDetail, orderID)
	c.JSON(http.StatusOK, oDetail)
}
