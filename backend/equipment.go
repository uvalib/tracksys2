package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type equipment struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serialNumber"`
	Status       uint      `json:"status"` // [:active, :inactive, :retired]
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type workstation struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Status       uint        `json:"status"` // [:active, :inactive, :retired]
	Equipment    []equipment `gorm:"many2many:workstation_equipment" json:"equipment"`
	ProjectCount int64       `gorm:"column:proj_cnt" json:"projectCount"`
	CreatedAt    time.Time   `json:"-"`
	UpdatedAt    time.Time   `json:"-"`
}

func (svc *serviceContext) getEquipment(c *gin.Context) {
	log.Printf("INfO: get all equppment")
	var resp struct {
		Workstations []workstation `json:"workstations"`
		Equipment    []equipment   `json:"equipment"`
	}

	// only count projects that are less than one year old that have been started, but not finished
	lastYear := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	projQ := fmt.Sprintf("(select count(*) from projects p where workstations.id = p.workstation_id and finished_at is null and started_at > '%s') as proj_cnt", lastYear)
	err := svc.DB.Debug().Preload("Equipment").Select(projQ, "workstations.*").
		Where("status != ?", 2).
		Order("name asc").Find(&resp.Workstations).Error
	if err != nil {
		log.Printf("ERROR: unable to get workstations: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.DB.Order("type asc, name asc").Find(&resp.Equipment).Error
	if err != nil {
		log.Printf("ERROR: unable to get equipment: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) updateWorkstation(c *gin.Context) {
	wsID := c.Param("id")
	tgtStatus, _ := strconv.ParseInt(c.Query("status"), 10, 8)
	var tgtWS workstation
	err := svc.DB.Find(&tgtWS, wsID).Error
	if err != nil {
		log.Printf("ERROR: unable to load workstation %s: %s", wsID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: update workstation %d status to %d", tgtWS.ID, tgtStatus)
	tgtWS.Status = uint(tgtStatus)
	err = svc.DB.Model(&tgtWS).Update("status", tgtStatus).Error
	if err != nil {
		log.Printf("ERROR: unable to update workstation %d status: %s", tgtWS.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "udpated")
}
