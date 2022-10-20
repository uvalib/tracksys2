package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (svc *serviceContext) getUnitCloneSources(c *gin.Context) {
	tgtUnitID := c.Param("id")
	log.Printf("INFO: get clone sources for unit %s", tgtUnitID)
	var tgtUnit unit
	err := svc.DB.Find(&tgtUnit, tgtUnitID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("ERROR: unit %s not found", tgtUnitID)
			c.String(http.StatusNotFound, fmt.Sprintf("unit %s not found", tgtUnitID))
		} else {
			log.Printf("ERROR: unable to retrieve unit %s: %s", tgtUnitID, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	var units []unit
	err = svc.DB.Where("master_files_count>0 and reorder=false and id<>? and throw_away=false", tgtUnit.ID).
		Where("metadata_id = ?", tgtUnit.MetadataID).Find(&units).Error
	if err != nil {
		log.Printf("ERROR: unable to get source units for clone of unit %s: %s", tgtUnitID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, units)
}
