package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type event struct {
	ID          uint64    `json:"id"`
	JobStatusID uint64    `json:"-"`
	Level       uint      `json:"level"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"createdAt"`
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
	log.Printf("INFO: get job %d statuses starting from offset %d", pageSize, startIndex)
	var total int64
	svc.DB.Table("job_statuses").Count(&total)

	var jobs []jobStatus
	err := svc.DB.Offset(startIndex).Limit(pageSize).Order("started_at desc").Find(&jobs).Error
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
	}).Find(&js).Error
	if err != nil {
		log.Printf("ERROR: unable to get job details: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
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
