package main

import (
	"log"
	"net/http"
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
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Status    uint        `json:"status"` // [:active, :inactive, :retired]
	Equipment []equipment `gorm:"many2many:workstation_equipment" json:"equipment"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
}

func (svc *serviceContext) getEquipment(c *gin.Context) {
	log.Printf("INfO: get all equppment")
	var resp struct {
		Workstations []workstation `json:"workstations"`
		Equipment    []equipment   `json:"equipment"`
	}
	err := svc.DB.Preload("Equipment").Order("name asc").Find(&resp.Workstations).Error
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
