package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type event struct {
	ID          uint      `json:"-"`
	JobStatusID uint      `json:"-"`
	Level       uint      `json:"level"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"createdAt"`
}

type jobStatus struct {
	ID             uint       `json:"id"`
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
	pageSize := 30
	pageQ := c.Query("page")
	if pageQ == "" {
		pageQ = "1"
	}
	page, err := strconv.Atoi(pageQ)
	if err != nil {
		log.Printf("ERROR: invalid page %s specified, default to 1", pageQ)
		page = 1
	}
	offset := (page - 1) * pageSize
	log.Printf("INFO: get job %d statuses starting from offset %d", pageSize, offset)
	var total int64
	svc.DB.Table("job_statuses").Count(&total)

	var jobs []jobStatus
	err = svc.DB.Offset(offset).Limit(pageSize).Order("started_at desc").Find(&jobs).Error
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
	jobID := c.Param("id")
	log.Printf("INFO: get job %s details", jobID)
	var out []event
	err := svc.DB.Order("created_at asc").Find(&out).Error
	if err != nil {
		log.Printf("ERROR: unable to get job details: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, out)
}
