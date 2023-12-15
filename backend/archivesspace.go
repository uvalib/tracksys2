package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type archivesspaceReview struct {
	ID            int64
	MetadataID    int64        `gorm:"column:metadata_id" json:"metadataID"`
	Metadata      *metadata    `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	SubmitStaffID int64        `gorm:"column:submit_staff_id" json:"-"`
	Submitter     staffMember  `gorm:"foreignKey:SubmitStaffID" json:"submitter"`
	SubmittedAt   time.Time    `json:"submittedAt"`
	ReviewStaffID *int64       `gorm:"column:review_staff_id" json:"-"`
	Reviewer      *staffMember `gorm:"foreignKey:ReviewStaffID" json:"reviewer,omitempty"`
	Status        string       `json:"status"`
	Notes         string       `json:"notes"`
	PublishedAt   *time.Time   `json:"publishedAt,omitempty"`
}

type asReviewsResponse struct {
	Total   int64                 `json:"total"`
	Reviews []archivesspaceReview `json:"submissions"`
}

func (svc *serviceContext) getArchivesSpaceReviews(c *gin.Context) {
	// queryStr := c.Query("q")

	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	resp := asReviewsResponse{}
	countQ := "select count(id) as total from archivesspace_reviews"
	err := svc.DB.Raw(countQ).Scan(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get archivesspace reviews count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.DB.Debug().Joins("Submitter").Joins("Reviewer").Joins("Metadata").Order("submitted_at asc").Find(&resp.Reviews).Error
	if err != nil {
		log.Printf("ERROR: unable to get archivesspace reviews: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
