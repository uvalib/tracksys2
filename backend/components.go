package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type componentType struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Descrition string `json:"description"`
}

type component struct {
	ID                int64         `json:"id"`
	ParentComponentID int64         `json:"-"`
	PID               string        `gorm:"column:pid" json:"pid"`
	Title             string        `json:"title"`
	Label             string        `json:"label"`
	ContentDesc       string        `json:"description"`
	Date              string        `json:"date"`
	Level             string        `json:"level"`
	Barcode           string        `json:"barcode"`
	EadIDAtt          string        `gorm:"column:ead_id_att" json:"eadID"`
	Ancestry          string        `json:"ancestry"`
	ComponentTypeID   int64         `json:"-"`
	ComponentType     componentType `gorm:"foreignKey:ComponentTypeID" json:"componentType"`
	DateDLIngest      *time.Time    `gorm:"date_dl_ingest" json:"dateDLIngest"`
	DateDLUpdate      *time.Time    `gorm:"date_dl_update" json:"dateDLUpdate"`
	Children          []*component  `gorm:"-" json:"children,omitempty"`
}

func (svc *serviceContext) getComponentTree(c *gin.Context) {
	cID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	log.Printf("INFO: get component tree for component %d", cID)
	var tgtCmp *component
	err := svc.DB.Preload("ComponentType").Find(&tgtCmp, cID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("INFO: component %d not found", cID)
			c.String(http.StatusNotFound, fmt.Sprintf("%d not found", cID))
		} else {
			log.Printf("ERROR:unable to load component %d: %s", cID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ancestryParts := strings.Split(tgtCmp.Ancestry, "/")
	if len(ancestryParts) == 0 {
		log.Printf("INFO: component %d is a top level component", cID)
		c.JSON(http.StatusOK, tgtCmp)
		return
	}

	topID, _ := strconv.ParseInt(ancestryParts[0], 10, 64)
	log.Printf("INFO: component %d is part of a tree rooted at %d", cID, topID)
	var topComponent *component
	err = svc.DB.Preload("ComponentType").Find(&topComponent, topID).Error
	if err != nil {
		log.Printf("ERROR: unable to get component %d top level parent %d: %s", cID, topID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.getComponentChildren(topComponent)
	if err != nil {
		log.Printf("ERROR: unable to get component %d children: %s", topID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, topComponent)
}

func (svc *serviceContext) getComponentChildren(parentC *component) error {
	log.Printf("INFO: get children for %d", parentC.ID)
	var children []*component
	err := svc.DB.Preload("ComponentType").Where("parent_component_id=?", parentC.ID).Find(&children).Error
	if err != nil {
		return err
	}
	parentC.Children = children

	for _, cmp := range children {
		err := svc.getComponentChildren(cmp)
		if err != nil {
			return err
		}
	}
	return nil
}
