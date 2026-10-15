package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	RoleName    string    `gorm:"-" json:"role"` // admin=0,supervisor=1,student=2,viewer=3
}

func (sm *staffMember) roleString() string {
	switch sm.Role {
	case admin:
		return "admin"
	case supervisor:
		return "supervisor"
	case student:
		return "student"
	case viewer:
		return "viewer"
	default:
		return "unknown"
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

	var staff []*staffMember
	if err := svc.DB.Order("last_name asc").Find(&staff).Error; err != nil {
		log.Printf("ERROR: unable to get staff: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, s := range staff {
		s.RoleName = s.roleString()
	}

	c.JSON(http.StatusOK, staff)
}
