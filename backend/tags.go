package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tag struct {
	ID  int64  `json:"id"`
	Tag string `json:"tag"`
}

func (svc *serviceContext) getTags(c *gin.Context) {
	tagQ := c.Query("q")

	log.Printf("INFO: get tags using query [%s]", tagQ)
	var tags []tag
	q := svc.DB.Order("tag asc")
	if tagQ != "" {
		q = q.Where("tag like ?", fmt.Sprintf("%s%%", tagQ))
	}
	err := q.Find(&tags).Error
	if err != nil {
		log.Printf("ERROR: get tags failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (svc *serviceContext) createTag(c *gin.Context) {
	var req struct {
		Tag string `json:"tag"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create tag request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	newTag := tag{Tag: req.Tag}
	err = svc.DB.Create(&newTag).Error
	if err != nil {
		log.Printf("Unable to create tag %s: %s", newTag.Tag, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var tags []tag
	svc.DB.Order("tag asc").Find(&tags)
	c.JSON(http.StatusOK, tags)
}
