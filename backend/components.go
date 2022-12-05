package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	Ancestry          string        `json:"-"`
	ComponentTypeID   int64         `json:"-"`
	ComponentType     componentType `gorm:"foreignKey:ComponentTypeID" json:"componentType"`
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

	var topComponent *component
	if tgtCmp.Ancestry == "" {
		log.Printf("INFO: component %d is a top level component", cID)
		topComponent = tgtCmp
	} else {
		ancestryParts := strings.Split(tgtCmp.Ancestry, "/")
		topID, _ := strconv.ParseInt(ancestryParts[0], 10, 64)
		log.Printf("INFO: component %d is part of a tree rooted at %d", cID, topID)
		err = svc.DB.Preload("ComponentType").Find(&topComponent, topID).Error
		if err != nil {
			log.Printf("ERROR: unable to get component %d top level parent %d: %s", cID, topID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	log.Printf("INFO: get all children for component %d", topComponent.ID)
	var children []*component
	err = svc.DB.Preload("ComponentType").Where("ancestry like ?", fmt.Sprintf("%d%%", topComponent.ID)).Order("id asc").Find(&children).Error
	if err != nil {
		log.Printf("ERROR: unable to get children of %d: %s", topComponent.ID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: arrange %d children into a hierarchy", len(children))
	for _, c := range children {
		ancestryBits := strings.Split(c.Ancestry, "/")
		if len(ancestryBits) == 1 {
			// log.Printf("INFO: component %d is top level child of %d", c.ID, topComponent.ID)
			topComponent.Children = append(topComponent.Children, c)
		} else {
			// toss the head as that is the topComponent and not part of the children array
			ancestryBits = ancestryBits[1:]
			// log.Printf("INFO: component %d is nested child of %d: %v", c.ID, topComponent.ID, ancestryBits)
			parentComponent := topComponent
			for _, aID := range ancestryBits {
				tgtID, _ := strconv.ParseInt(aID, 10, 64)
				for _, pc := range parentComponent.Children {
					if pc.ID == tgtID {
						// log.Printf("INFO: found matching child %d in parent component %d", pc.ID, parentComponent.ID)
						parentComponent = pc
						break
					}
				}
			}
			// log.Printf("INFO: add %d to parent %d", c.ID, parentComponent.ID)
			parentComponent.Children = append(parentComponent.Children, c)
		}
	}

	log.Printf("INFO: find masterfiles related to componend %d", cID)
	var related []*masterFile
	err = svc.DB.Preload("Metadata").Where("component_id=?", cID).Find(&related).Error
	if err != nil {
		log.Printf("ERROR: unable to get master files related to component %d: %s", cID, err.Error())
	}

	for idx, mf := range related {
		mfPID := mf.PID
		mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		if mf.MetadataID != nil {
			mf.ViewerURL = fmt.Sprintf("%s/view/%s?unit=%d", svc.ExternalSystems.Curio, mf.Metadata.PID, mf.UnitID)
			if idx > 0 {
				mf.ViewerURL += fmt.Sprintf("&page=%d", (idx + 1))
			}
		} else {
			mf.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		}
	}

	resp := struct {
		Component   *component    `json:"component"`
		MasterFiles []*masterFile `json:"masterFiles,omitempty"`
	}{
		Component:   topComponent,
		MasterFiles: related,
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getComponentMasterFiles(c *gin.Context) {
	cID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	log.Printf("INFO: get master files related to component %d", cID)
	var related []*masterFile
	err := svc.DB.Preload("Metadata").Where("component_id=?", cID).Find(&related).Error
	if err != nil {
		log.Printf("ERROR: unable to get master files related to component %d: %s", cID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
	}
	for idx, mf := range related {
		mfPID := mf.PID
		mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		if mf.MetadataID != nil {
			mf.ViewerURL = fmt.Sprintf("%s/view/%s?unit=%d", svc.ExternalSystems.Curio, mf.Metadata.PID, mf.UnitID)
			if idx > 0 {
				mf.ViewerURL += fmt.Sprintf("&page=%d", (idx + 1))
			}
		} else {
			mf.ViewerURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, mfPID)
		}
	}
	c.JSON(http.StatusOK, related)
}
