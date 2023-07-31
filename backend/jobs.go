package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type eventLevel uint

const (
	info  eventLevel = 0
	warn  eventLevel = 1
	err   eventLevel = 2
	fatal eventLevel = 3
)

type event struct {
	ID          uint64     `json:"id"`
	JobStatusID uint64     `json:"-"`
	Level       eventLevel `json:"level"`
	Text        string     `json:"text"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type jobStatus struct {
	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	OriginatorType string     `json:"originatorType"`
	OriginatorID   string     `json:"originatorID"`
	Failures       uint64     `json:"failures"`
	Error          string     `json:"error"`
	Events         []event    `gorm:"foreignKey:JobStatusID" json:"events"`
	StartedAt      time.Time  `json:"startedAt"`
	EndedAt        *time.Time `json:"finishedAt"`
}

func (svc *serviceContext) getJobStatuses(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	queryStr := c.Query("q")
	var queryClause *gorm.DB
	if queryStr != "" {
		qLike := fmt.Sprintf("%s%%", queryStr)
		queryClause = svc.DB.Where("name like ? or status=? or originator_id like ? or originator_type like ?", qLike, queryStr, qLike, qLike)
	}

	log.Printf("INFO: get job %d statuses starting from offset %d", pageSize, startIndex)
	var total int64
	var jobs []jobStatus
	var err error
	if queryClause != nil {
		svc.DB.Where(queryClause).Table("job_statuses").Count(&total)
		err = svc.DB.Where(queryClause).Offset(startIndex).Limit(pageSize).Order("started_at desc").Find(&jobs).Error
	} else {
		svc.DB.Table("job_statuses").Count(&total)
		err = svc.DB.Offset(startIndex).Limit(pageSize).Order("started_at desc").Find(&jobs).Error
	}

	if err != nil {
		log.Printf("ERROR: unable to get job statuses: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type resp struct {
		Jobs  []jobStatus `json:"jobs"`
		Total int64       `json:"total"`
	}
	out := resp{Jobs: jobs, Total: total}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getJobDetails(c *gin.Context) {
	jobID, _ := strconv.Atoi(c.Param("id"))
	log.Printf("INFO: get job %d details", jobID)
	js := jobStatus{ID: uint64(jobID)}
	err := svc.DB.Preload("Events", func(db *gorm.DB) *gorm.DB {
		return db.Order("events.created_at ASC")
	}).First(&js).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("INFO: job %d not found: %s", jobID, err.Error())
			c.String(http.StatusNotFound, err.Error())
		} else {
			log.Printf("ERROR: unable to get job details: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	type resp struct {
		Status string  `json:"status"`
		Error  string  `json:"error"`
		Events []event `json:"events"`
	}
	out := resp{Status: js.Status, Error: js.Error, Events: js.Events}
	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) deleteJobStatuses(c *gin.Context) {
	type delJobRequest struct {
		Jobs []uint64 `json:"jobs"`
	}
	var delReq delJobRequest
	err := c.BindJSON(&delReq)
	if err != nil {
		log.Printf("ERROR: invalid delete jobs request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.DB.Delete(jobStatus{}, delReq.Jobs).Error
	if err != nil {
		log.Printf("ERROR: unable to delete jobs %v: %s", delReq.Jobs, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, delReq)
}
