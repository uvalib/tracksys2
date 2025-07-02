package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type invoice struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"-"`
	DateInvoice       time.Time  `json:"invoiceDate"`
	DateFeePaid       *time.Time `json:"dateFeePaid,omitempty"`
	DateFeeDeclined   *time.Time `json:"dateFeeDeclined,omitempty"`
	FeeAmountPaid     *float64   `json:"feeAmountPaid"`
	TransmittalNumber string     `json:"transmittalNumber"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"-"`
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

type order struct {
	ID                             int64        `json:"id"`
	OrderStatus                    string       `json:"status"`
	OrderTitle                     string       `json:"title"`
	DateDue                        time.Time    `json:"dateDue"`
	CustomerID                     *uint        `json:"-"`
	Customer                       *customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProcessorID                    *int64       `json:"-"`
	Processor                      *staffMember `gorm:"foreignKey:ProcessorID" json:"processor,omitempty"`
	AgencyID                       *uint        `json:"-"`
	Agency                         *agency      `gorm:"foreignKey:AgencyID" json:"agency,omitempty"`
	Fee                            *float64     `json:"fee,omitempty"`
	FeeWaived                      bool         `json:"feeWaived"`
	Invoice                        *invoice     `gorm:"-" json:"invoice,omitempty"`
	UnitCount                      int64        `json:"unitCount"`       // NOTE: this is different than the cached count field units_count
	MasterFileCount                int64        `json:"masterFileCount"` // NOTE: this is different than the cached count master_files_count
	Email                          string       `json:"email"`
	StaffNotes                     string       `json:"staffNotes"`
	SpecialInstructions            string       `json:"specialInstructions"`
	DateRequestSubmitted           time.Time    `json:"dateSubmitted"`
	DateOrderApproved              *time.Time   `json:"dateApproved"`
	DateDeferred                   *time.Time   `json:"dateDeferred"`
	DateCanceled                   *time.Time   `json:"dateCanceled"`
	DateFeeWaived                  *time.Time   `json:"dateFeeWaived,omitempty"`
	DateCustomerNotified           *time.Time   `json:"dateCustomerNotified"`
	DatePatronDeliverablesComplete *time.Time   `json:"datePatronDeliverablesComplete"`
	DateArchivingComplete          *time.Time   `json:"dateArchivingComplete"`
	DateFinalizationBegun          *time.Time   `json:"dateFinalizationBegun"`
	DateFeeEstimateSentToCustomer  *time.Time   `json:"dateFeeEstimateSent"`
	DateCompleted                  *time.Time   `json:"dateCompleted"`
	UpdatedAt                      time.Time    `json:"-"`
}

type orderRequest struct {
	Status              string  `json:"status"`
	DateDue             string  `json:"dateDue"`
	Title               string  `json:"title"`
	SpecialInstructions string  `json:"specialInstructions"`
	StaffNotes          string  `json:"staffNotes"`
	Fee                 float64 `json:"fee"`
	AgencyID            uint    `json:"agencyID"`
	CustomerID          uint    `json:"customerID"`
}

func (svc *serviceContext) deleteOrder(c *gin.Context) {
	oID := c.Param("id")
	log.Printf("INFO: delete order %s", oID)
	var unitCnt int64
	err := svc.DB.Table("units").Where("order_id=?", oID).Count(&unitCnt).Error
	if err != nil {
		log.Printf("ERROR: unable to determine if unit %s has any units: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if unitCnt > 0 {
		log.Printf("INFO: unable to delete order %s because it already has %d units", oID, unitCnt)
		c.String(http.StatusPreconditionFailed, "order has units and cannont be deleted")
		return
	}

	err = svc.DB.Delete(&order{}, oID).Error
	if err != nil {
		log.Printf("ERROR: unable to delete order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.DB.Where("order_id=?", oID).Delete(orderItem{}).Error
	if err != nil {
		log.Printf("ERROR: unable to delete order %s items: %s", oID, err.Error())
	}

	c.String(http.StatusOK, "deleted")
}

func (svc *serviceContext) createOrder(c *gin.Context) {
	var req orderRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create order request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	dueDate, _ := parseDateString(req.DateDue)
	newOrder := order{OrderStatus: "requested", DateDue: dueDate, OrderTitle: req.Title,
		SpecialInstructions: req.SpecialInstructions, StaffNotes: req.StaffNotes,
		CustomerID: &req.CustomerID, DateRequestSubmitted: time.Now()}

	omitFields := []string{"UnitCount", "MasterFileCount", "Fee"}
	if req.AgencyID > 0 {
		newOrder.AgencyID = &req.AgencyID
	} else {
		omitFields = append(omitFields, "AgencyID")
	}

	err = svc.DB.Omit(omitFields...).Create(&newOrder).Error
	if err != nil {
		log.Printf("ERROR: unable to create new order %+v: %s", req, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: order %d updated", newOrder.ID)
	svc.DB.Preload("Agency").Preload("Customer").
		Preload("Customer.AcademicStatus").Preload("Customer.Addresses").
		Find(&newOrder, newOrder.ID)
	c.JSON(http.StatusOK, newOrder)
}

func (svc *serviceContext) loadOrder(orderID string) (*order, error) {
	log.Printf("INFO: load order %s details", orderID)
	var oDetail order
	err := svc.DB.Preload("Agency").Preload("Customer").Preload("Processor").
		Preload("Customer.AcademicStatus").Preload("Customer.AcademicStatus").Preload("Customer.Addresses").
		Limit(1).Find(&oDetail, orderID).Error
	if err != nil {
		return nil, err
	}

	if oDetail.ID == 0 {
		return &oDetail, nil
	}

	// Due date is stored as a date rather than datetime, it has no timezone info in the field. When it is
	// pulled from the DB into a time.Time variable, the zone defaults to UTC.  When is converted on the client, the time
	// difference in UTC vs EST/EDT bumps the date back one day. To fix, fomat the date as a string without time info (scrap the UTC)
	// then parse that date onto a datetime in the local timezone.
	dueStr := oDetail.DateDue.Format("2006-01-02")
	oDetail.DateDue, _ = parseDateString(dueStr)

	log.Printf("INFO: lookup invoice for order %d", oDetail.ID)
	var invDetail invoice
	err = svc.DB.Where("order_id=?", orderID).Order("created_at desc").First(&invDetail).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: unable to get invoice for order %s: %s", orderID, err.Error())
		} else {
			log.Printf("INFO: no invoice for order %s", orderID)
		}
	} else {
		oDetail.Invoice = &invDetail
	}

	return &oDetail, nil
}

func (svc *serviceContext) getOrderDetails(c *gin.Context) {
	oID := c.Param("id")
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if oDetail.ID == 0 {
		log.Printf("INFO: order %s not found", oID)
		c.String(http.StatusNotFound, "not found")
		return
	}

	type oResp struct {
		Order *order      `json:"order"`
		Units []unit      `json:"units"`
		Items []orderItem `json:"items"`
	}
	out := oResp{Order: oDetail}
	// NOTE: Manually calculate the master files count and return it as num_master_files instead of using the inaccurate cache
	mfCnt := "(select count(*) from master_files m inner join units u on u.id=m.unit_id where u.id=units.id) as num_master_files"
	err = svc.DB.Where("order_id=?", oID).Preload("IntendedUse").
		Preload("Metadata").Preload("Metadata.HathiTrustStatus").
		Select("units.*", mfCnt).Find(&out.Units).Error
	if err != nil {
		log.Printf("ERROR: unable to get units for order %s: %s", oID, err.Error())
	}
	err = svc.DB.Where("order_id=?", oID).Preload("IntendedUse").Find(&out.Items).Error
	if err != nil {
		log.Printf("ERROR: unable to get items for order %s: %s", oID, err.Error())
	}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getOrders(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
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
	switch sortBy {
	case "dateDue":
		sortField = "date_due"
	case "dateSubmitted":
		sortField = "date_request_submitted"
	case "title":
		sortField = "order_title"
	case "unitCount":
		sortField = "unit_count"
	case "masterFileCount":
		sortField = "master_file_count"
	}
	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	log.Printf("INFO: get %d orders starting from offset %d order %s", pageSize, startIndex, orderStr)

	// set up filtering....
	filterStr := c.Query("filters")
	log.Printf("INFO: raw filters [%s]", filterStr)
	var filters []string
	err := json.Unmarshal([]byte(filterStr), &filters)
	if err != nil {
		log.Printf("ERROR: unable to parse filter payload %s: %s", filterStr, err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid filters param: %s", filterStr))
		return
	}

	filterQ := svc.DB.Table("orders").Joins("inner join customers c on c.id=orders.customer_id").
		Joins("left outer join staff_members p on p.id = processor_id").
		Joins("left outer join agencies a on a.id = orders.agency_id")
	dateNow := time.Now().Format("2006-01-02")
	ownerID := c.Query("owner")
	if ownerID != "" {
		filterQ.Where("processor_id=?", ownerID)
	}

	for _, filter := range filters {
		bits := strings.Split(filter, "|") // target | comparison | value
		tgtField := bits[0]
		comparison := bits[1]
		tgtVal, _ := url.QueryUnescape(bits[2])
		log.Printf("INFO: filter %s %s %s", tgtField, comparison, tgtVal)
		if tgtField == "status" {
			switch tgtVal {
			case "active":
				filterQ = filterQ.Where("order_status!=? and order_status!=?", "canceled", "completed")
			case "await":
				filterQ = filterQ.Where("order_status=? or order_status=?", "requested", "await_fee")
			case "deferred":
				filterQ = filterQ.Where("order_status=?", "deferred")
			case "canceled":
				filterQ = filterQ.Where("order_status=?", "canceled")
			case "complete":
				filterQ = filterQ.Where("order_status=?", "completed")
			case "due_week":
				dateWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
				filterQ = filterQ.Joins("inner join units u on u.order_id=orders.id").
					Where("u.intended_use_id <> ?", 110).
					Where("date_due>=?", dateNow).Where("date_due<=?", dateWeek).
					Where("order_status!=?", "completed").Where("order_status!=?", "deferred").Where("order_status!=?", "canceled").Distinct("orders.id")
			case "overdue":
				oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
				filterQ = filterQ.Joins("inner join units u on u.order_id=orders.id").
					Where("u.intended_use_id <> ?", 110).
					Where("date_request_submitted>?", oneYearAgo).Where("date_due<?", dateNow).
					Where("order_status!=?", "completed").Where("order_status!=?", "deferred").Where("order_status!=?", "canceled").Distinct("orders.id")
			case "ready":
				filterQ = filterQ.Joins("inner join units u on u.order_id=orders.id").
					Where("u.intended_use_id <> ?", 110).
					Where("orders.email is not null and orders.email != ? and date_customer_notified is null", "").
					Where("order_status != ? and order_status != ?", "canceled", "completed").Distinct("orders.id")
			}
		} else if tgtField == "customer" {
			filterQ = filterQ.Where("c.last_name like ?", fmt.Sprintf("%s%%", tgtVal))
		} else if tgtField == "agency" {
			filterQ = filterQ.Where("orders.agency_id = ?", tgtVal)
		} else if tgtField == "processor" {
			filterQ = filterQ.Where("p.last_name like ?", fmt.Sprintf("%s%%", tgtVal))
		}
	}

	// set up query...
	queryStr := c.Query("q")
	var qObj *gorm.DB
	if queryStr != "" {
		queryAny := fmt.Sprintf("%%%s%%", queryStr)
		queryStart := fmt.Sprintf("%s%%", queryStr)
		qObj = svc.DB.Where("order_title like ?", queryAny).Or("orders.staff_notes like ?", queryAny).
			Or("orders.special_instructions like ?", queryAny).Or("orders.id like ?", queryStart)
		filterQ = filterQ.Where(qObj)
	}

	var total int64
	filterQ.Count(&total)

	var o []*order
	// NOTE: the DB originally had fields named units_count and master_files_count which cached counts. These were often wrong
	// so they are ignored below. Instead, calculate the actual counts and store them in slightly differet names (not plural) to
	// avoid conflicts with the cached fields. The cached fields are ignored
	unitCnt := "(select count(*) from units where order_id=orders.id) as unit_count"
	mfCnt := "(select count(*) from master_files m inner join units u on u.id=m.unit_id where u.order_id=orders.id) as master_file_count"
	err = filterQ.Preload("Agency").Preload("Customer").Preload("Customer.AcademicStatus").Preload("Processor").
		Select("orders.*", unitCnt, mfCnt).
		Offset(startIndex).Limit(pageSize).Order(orderStr).Find(&o).Error
	if err != nil {
		log.Printf("ERROR: unable to get orders: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type resp struct {
		Orders []*order `json:"orders"`
		Total  int64    `json:"total"`
	}
	out := resp{Orders: o, Total: total}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) waiveFee(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for fee waive")
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s waives fee for order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.FeeWaived = true
	oDetail.DateFeeWaived = &now
	oDetail.Fee = nil
	err = svc.DB.Model(oDetail).Select("FeeWaived", "DateFeeWaived", "Fee").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to waive fee for order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) acceptFee(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for fee accept")
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s accepts fee for order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.OrderStatus = "approved"
	oDetail.DateOrderApproved = &now
	err = svc.DB.Model(oDetail).Select("OrderStatus", "DateOrderApproved").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to accept fee for order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: order has been accepted; remove all order items")
	err = svc.DB.Where("order_id=?", oDetail.ID).Delete(orderItem{}).Error
	if err != nil {
		log.Printf("ERROR: unable to delete order %s items: %s", oID, err.Error())
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) setOrderProcessor(c *gin.Context) {
	oID := c.Param("id")
	staffID, _ := strconv.ParseInt(c.Query("staff"), 10, 64)
	if staffID == 0 {
		log.Printf("ERROR: staff param required for to set processor for order %s", oID)
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s to set processor %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	oDetail.ProcessorID = &staffID
	err = svc.DB.Model(oDetail).Select("ProcessorID").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to set processor for order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	oDetail, _ = svc.loadOrder(oID)
	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) completeOrder(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for complete order %s", oID)
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s to mark as complete: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// check if this has patron deliverables (units with intended use NOT equal to 110)
	var orderUnits []unit
	patronOrder := false
	allUnitsArchived := true
	var latestArchiveDate *time.Time
	err = svc.DB.Where("order_id=? and unit_status <> ?", oDetail.ID, "canceled").Order("date_archived desc").Find(&orderUnits).Error
	if err != nil {
		log.Printf("ERROR: unable to determine if order %d has patron deliverables: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for idx, u := range orderUnits {
		if *u.IntendedUseID != 110 {
			patronOrder = true
		}
		if idx == 0 {
			latestArchiveDate = u.DateArchived
		}
		if u.DateArchived == nil {
			allUnitsArchived = false
		}
	}

	if patronOrder == true {
		log.Printf("INFO: staff %s completes patron order %d", staffID, oDetail.ID)
		if oDetail.DateCustomerNotified == nil {
			if oDetail.DatePatronDeliverablesComplete == nil {
				log.Printf("INFO: order %d cannot be completed because deliverables have not been generated", oDetail.ID)
				c.String(http.StatusPreconditionFailed, "deliverables have not been generated")
			} else {
				log.Printf("INFO: order %d cannot be completed because customer has not been notified", oDetail.ID)
				c.String(http.StatusPreconditionFailed, "customer has not been notified")
			}
			return
		}
	} else {
		log.Printf("INFO: staff %s completes digital collection building order %d", staffID, oDetail.ID)
		if oDetail.DateArchivingComplete == nil {
			if oDetail.DateFinalizationBegun == nil {
				log.Printf("INFO: order %d cannot be completed because it has not been finalized", oDetail.ID)
				c.String(http.StatusPreconditionFailed, "order has not been finalized")
				return
			}
			if allUnitsArchived == false {
				log.Printf("INFO: order %d cannot be completed because not al units have been archived", oDetail.ID)
				c.String(http.StatusPreconditionFailed, "not all units have been archived")
				return
			}
		}
	}

	now := time.Now()
	oDetail.OrderStatus = "completed"
	oDetail.DateCompleted = &now
	if oDetail.DateArchivingComplete == nil {
		oDetail.DateArchivingComplete = latestArchiveDate
	}
	err = svc.DB.Model(&oDetail).Select("OrderStatus", "DateCompleted", "DateArchivingComplete").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to mark order  %d as complete: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) approveOrder(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for approve order %s", oID)
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s for approval: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s approves order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.OrderStatus = "approved"
	oDetail.DateOrderApproved = &now
	err = svc.DB.Model(oDetail).Select("OrderStatus", "DateOrderApproved").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to approve order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: remove all order items from approved order %d", oDetail.ID)
	err = svc.DB.Where("order_id=?", oDetail.ID).Delete(orderItem{}).Error
	if err != nil {
		log.Printf("ERROR: unable to delete order %s items: %s", oID, err.Error())
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) cancelOrder(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for cancel order %s", oID)
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s for cancelation: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s cancels order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.OrderStatus = "canceled"
	oDetail.DateCanceled = &now
	err = svc.DB.Model(oDetail).Select("OrderStatus", "DateCanceled").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to cancel order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: order %d was canceled; cancel all associated units", oDetail.ID)
	err = svc.DB.Model(unit{}).Where("order_id = ?", oDetail.ID).Updates(unit{UnitStatus: "canceled"}).Error
	if err != nil {
		log.Printf("ERROR: unable to cancel units related to canceled order %d: %s", oDetail.ID, err.Error())
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) deferOrder(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for fee defer")
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s for defer: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s defers order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.OrderStatus = "deferred"
	oDetail.DateDeferred = &now
	err = svc.DB.Model(oDetail).Select("OrderStatus", "DateDeferred").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to defer order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) resumeOrder(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for resume order")
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s for defer: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s resumes order %d", staffID, oDetail.ID)
	oDetail.OrderStatus = "requested"
	if oDetail.DateOrderApproved != nil {
		oDetail.OrderStatus = "approved"
	}
	err = svc.DB.Model(oDetail).Select("OrderStatus").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to resume order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) declineFee(c *gin.Context) {
	oID := c.Param("id")
	staffID := c.Query("staff")
	if staffID == "" {
		log.Printf("ERROR: staff param required for fee decline")
		c.String(http.StatusBadRequest, "staff param is required")
		return
	}
	oDetail, err := svc.loadOrder(oID)
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", oID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: staff %s declines fee for order %d", staffID, oDetail.ID)
	now := time.Now()
	oDetail.OrderStatus = "canceled"
	oDetail.DateCanceled = &now
	err = svc.DB.Model(oDetail).Select("OrderStatus", "DateCanceled").Updates(oDetail).Error
	if err != nil {
		log.Printf("ERROR: unable to decline fee for order %d: %s", oDetail.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	oDetail.Invoice.DateFeeDeclined = &now
	err = svc.DB.Model(oDetail.Invoice).Select("DateFeeDeclined").Updates(oDetail.Invoice).Error
	if err != nil {
		log.Printf("ERROR: unable to updated invoice declined time for order %d: %s", oDetail.ID, err.Error())
	}

	c.JSON(http.StatusOK, oDetail)
}

func (svc *serviceContext) updateInvoice(c *gin.Context) {
	invoiceID := c.Param("id")
	log.Printf("INFO: update invoice %s", invoiceID)
	var inv invoice
	err := svc.DB.Find(&inv, invoiceID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve invoice %s: %s", invoiceID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var updateRequest struct {
		DateFeePaid       string `json:"dateFeePaid"`
		DateFeeDeclined   string `json:"dateFeeDeclined"`
		FeeAmountPaid     string `json:"feeAmountPaid"`
		TransmittalNumber string `json:"transmittalNumber"`
		Notes             string `json:"notes"`
	}
	err = c.BindJSON(&updateRequest)
	if err != nil {
		log.Printf("ERROR: invalid update invoice %s request: %s", invoiceID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if updateRequest.DateFeePaid != "" {
		paid, _ := parseDateString(updateRequest.DateFeePaid)
		inv.DateFeePaid = &paid
	} else {
		inv.DateFeePaid = nil
	}
	if updateRequest.DateFeeDeclined != "" {
		paid, _ := parseDateString(updateRequest.DateFeeDeclined)
		inv.DateFeeDeclined = &paid
	} else {
		inv.DateFeeDeclined = nil
	}
	if updateRequest.FeeAmountPaid != "" {
		fee, _ := strconv.ParseFloat(updateRequest.FeeAmountPaid, 64)
		inv.FeeAmountPaid = &fee
	} else {
		inv.FeeAmountPaid = nil
	}
	inv.TransmittalNumber = updateRequest.TransmittalNumber
	inv.Notes = updateRequest.Notes

	err = svc.DB.Model(&inv).Select("DateFeePaid", "DateFeeDeclined", "FeeAmountPaid", "PermanentNonPayment", "TransmittalNumber", "Notes").Updates(inv).Error
	if err != nil {
		log.Printf("ERROR: unable to update invoice %d: %s", inv.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: order %d updated", inv.ID)
	c.JSON(http.StatusOK, inv)
}
func (svc *serviceContext) deleteOrderItem(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("item")
	log.Printf("INFO: discard item %s from order %s", itemID, orderID)
	err := svc.DB.Delete(&orderItem{}, itemID).Error
	if err != nil {
		log.Printf("ERROR: unable to delete item %s: %s", itemID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "deleted")
}

func (svc *serviceContext) addUnitToOrder(c *gin.Context) {
	orderID := c.Param("id")
	log.Printf("INFO: add unit order %s", orderID)
	var addReq struct {
		ItemID              int64  `json:"itemID"`
		MetadataID          int64  `json:"metadataID"`
		IntendedUseID       int64  `json:"intendedUseID"`
		SourceURL           string `json:"sourceURL"`
		SpecialInstructions string `json:"specialInstructions"`
		StaffNotes          string `json:"staffNotes"`
		CompleteScan        bool   `json:"completeScan"`
		ThrowAway           bool   `json:"throwAway"`
		IncludeInDL         bool   `json:"includeInDL"`
	}
	err := c.BindJSON(&addReq)
	if err != nil {
		log.Printf("ERROR: invalid add unit %s request: %s", orderID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: validate order %s before adding unit", orderID)
	var tgtOrder order
	err = svc.DB.Find(&tgtOrder, orderID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve order %s: %s", orderID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: validate metadata %d for new unit in order %s", addReq.MetadataID, orderID)
	var md metadata
	err = svc.DB.Find(&md, addReq.MetadataID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve metadata %d for new unit in order %s: %s", addReq.MetadataID, orderID, err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("metadata record %d not found: %s", addReq.MetadataID, err.Error()))
		return
	}

	log.Printf("INFO: validate intended use %d for new unit in order %s", addReq.MetadataID, orderID)
	var iu intendedUse
	err = svc.DB.Find(&iu, addReq.IntendedUseID).Error
	if err != nil {
		log.Printf("ERROR: unable to retrieve intended use %d for new unit in order %s: %s", addReq.IntendedUseID, orderID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if addReq.IncludeInDL {
		err = svc.validateIncludeInDL(md.ID, 0)
		if err != nil {
			log.Printf("ERROR: cannnot add new unit flagged for inclusion in virgo: %s", err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	if addReq.ItemID != 0 {
		log.Printf("INFO: update order item %d for new unit in order %s", addReq.MetadataID, orderID)
		var tgtItem orderItem
		err = svc.DB.Find(&tgtItem, addReq.ItemID).Error
		if err != nil {
			log.Printf("ERROR: unable to retrieve orderItem %d for new unit in order %s: %s", addReq.ItemID, orderID, err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if tgtItem.Converted == false {
			tgtItem.Converted = true
			svc.DB.Model(&tgtItem).Select("Converted").Updates(tgtItem)
		} else {
			log.Printf("INFO: item %d has already been marked as converted", tgtItem.ID)
		}
	}

	log.Printf("INFO: create new unit in order %s", orderID)
	newUnit := unit{UnitStatus: "approved", MetadataID: &md.ID, PatronSourceURL: addReq.SourceURL,
		SpecialInstructions: addReq.SpecialInstructions, StaffNotes: addReq.StaffNotes, CompleteScan: addReq.CompleteScan,
		ThrowAway: addReq.ThrowAway, IncludeInDL: addReq.IncludeInDL, CreatedAt: time.Now(), OrderID: tgtOrder.ID,
		IntendedUseID: &iu.ID, IntendedUse: &iu, Metadata: &md,
	}
	err = svc.DB.Omit("num_master_files").Create(&newUnit).Error
	if err != nil {
		log.Printf("ERROR: unable to create new unit for order %s", orderID)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: Update unit count for order %s", orderID)
	var unitCnt int64
	err = svc.DB.Table("units").Where("order_id=?", tgtOrder.ID).Count(&unitCnt).Error
	if err != nil {
		log.Printf("ERROR: unable to get unit count for order %d: %s", tgtOrder.ID, err.Error())
	} else {
		tgtOrder.UnitCount = unitCnt
		err = svc.DB.Model(&tgtOrder).Select("UnitsCount").Updates(tgtOrder).Error
		if err != nil {
			log.Printf("ERROR: unable to update unit count for order %d: %s", tgtOrder.ID, err.Error())
		}
	}

	c.JSON(http.StatusOK, newUnit)
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

	var updateRequest orderRequest
	err = c.BindJSON(&updateRequest)
	if err != nil {
		log.Printf("ERROR: invalid update order %s request: %s", orderID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: submitted due date [%s]", updateRequest.DateDue)

	oDetail.OrderStatus = updateRequest.Status
	oDetail.DateDue, err = parseDateString(updateRequest.DateDue)
	if err != nil {
		log.Printf("ERROR: unable to parse due date %s: %s", updateRequest.DateDue, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: parsed due date [%v]", oDetail.DateDue)
	oDetail.OrderTitle = updateRequest.Title
	oDetail.SpecialInstructions = updateRequest.SpecialInstructions
	oDetail.StaffNotes = updateRequest.StaffNotes
	oDetail.Fee = nil
	if updateRequest.Fee > 0 {
		oDetail.Fee = &updateRequest.Fee
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
	if oDetail.OrderStatus == "canceled" {
		log.Printf("INFO: order %d was canceled; cancel all associated units", oDetail.ID)
		err = svc.DB.Model(unit{}).Where("order_id = ?", oDetail.ID).Updates(unit{UnitStatus: "canceled"}).Error
		if err != nil {
			log.Printf("ERROR: unable to cancel units related to canceled order %d: %s", oDetail.ID, err.Error())
		}
	} else {
		log.Printf("INFO: order %d updated", oDetail.ID)
	}
	updated, _ := svc.loadOrder(orderID)
	c.JSON(http.StatusOK, updated)
}
