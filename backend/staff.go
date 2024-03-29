package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type staffRole uint

const (
	admin      staffRole = 0
	supervisor staffRole = 1
	student    staffRole = 2
	viewer     staffRole = 3
)

type staffMember struct {
	ID          uint      `json:"id"`
	ComputingID string    `json:"computingID"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"active"`
	Role        staffRole `json:"roleID"`
	RoleName    string    `gorm:"-" json:"role"`
}

func (sm *staffMember) roleString() string {
	switch sm.Role {
	case admin:
		return "admin"
	case supervisor:
		return "supervisor"
	case student:
		return "student"
	default:
		return "viewer"
	}
}

func (svc *serviceContext) addOrUpdateStaff(c *gin.Context) {
	var staffReq staffMember
	err := c.BindJSON(&staffReq)
	if err != nil {
		log.Printf("ERROR: invalid staff add/update request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: add or update staff request: %+v", staffReq)
	if staffReq.ID == 0 {
		err = svc.DB.Create(&staffReq).Error
		if err != nil {
			log.Printf("ERROR: unable to create staff member: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		err = svc.DB.Save(&staffReq).Error
		if err != nil {
			log.Printf("ERROR: unable to update staff member %d: %s", staffReq.ID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, staffReq)
}

func (svc *serviceContext) getStaff(c *gin.Context) {
	adminsOnly, _ := strconv.ParseBool(c.Query("admins"))
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "lastName"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// just get active admin stafff
	if adminsOnly {
		log.Printf("INFO: get all active admin/supervisor staff")
		var out []staffMember
		// 0=admin, 1=supervisor
		err := svc.DB.Where("is_active=? and (role=? or role=?)", 1, 0, 1).Order("last_name asc").Find(&out).Error
		if err != nil {
			log.Printf("ERROR: unable to get active admin users: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		log.Printf("INFO: %d active admin staff found", len(out))
		c.JSON(http.StatusOK, out)
		return
	}

	// normal paged, filterd and sorted retrieval
	fieldMap := map[string]string{"lastName": "last_name", "email": "email", "computingID": "computing_id"}
	sortField := fieldMap[sortBy]
	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	queryStr := c.Query("q")
	var qObj *gorm.DB
	if queryStr != "" {
		qObj = svc.DB.Where("last_name like ?", fmt.Sprintf("%s%%", queryStr)).
			Or("email like ?", fmt.Sprintf("%s%%", queryStr)).
			Or("computing_id like ?", fmt.Sprintf("%s%%", queryStr))
	}

	log.Printf("INFO: get %d staff starting from offset %d", pageSize, startIndex)
	var total int64
	countQ := svc.DB.Table("staff_members").Where("last_name <> ?", "")
	if queryStr != "" {
		countQ.Where(qObj)
	}
	countQ.Count(&total)

	var staff []staffMember
	mainQ := svc.DB.Offset(startIndex).Where("last_name <> ?", "")
	if queryStr != "" {
		mainQ.Where(qObj)
	}
	err := mainQ.Order(orderStr).Limit(pageSize).Find(&staff).Error
	if err != nil {
		log.Printf("ERROR: unable to get staff: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type resp struct {
		Staff []staffMember `json:"staff"`
		Total int64         `json:"total"`
	}
	out := resp{Total: total}
	for _, s := range staff {
		s.RoleName = s.roleString()
		out.Staff = append(out.Staff, s)
	}

	c.JSON(http.StatusOK, out)
}
