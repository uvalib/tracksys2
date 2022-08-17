package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type staffMember struct {
	ID          uint   `json:"id"`
	ComputingID string `json:"computingID"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	IsActive    bool   `json:"active"`
	Role        uint   `json:"roleID"`
	RoleName    string `gorm:"-" json:"role"`
}

func (sm *staffMember) roleString() string {
	roles := []string{"admin", "supervisor", "student", "viewer"}
	if sm.Role < 0 || sm.Role > uint(len(roles)-1) {
		return "viewer"
	}
	return roles[sm.Role]
}

func (svc *serviceContext) getStaff(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	log.Printf("INFO: get %d staff starting from offset %d", pageSize, startIndex)
	var total int64
	svc.DB.Table("staff_members").Count(&total)

	var staff []staffMember
	err := svc.DB.Offset(startIndex).Where("last_name <> ?", "").Limit(pageSize).Order("last_name asc").Find(&staff).Error
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
