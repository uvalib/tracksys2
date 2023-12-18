package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type archivesspaceReview struct {
	ID              int64        `json:"id"`
	MetadataID      int64        `gorm:"column:metadata_id" json:"metadataID"`
	Metadata        *metadata    `gorm:"foreignKey:MetadataID" json:"metadata,omitempty"`
	SubmitStaffID   int64        `gorm:"column:submit_staff_id" json:"-"`
	Submitter       staffMember  `gorm:"foreignKey:SubmitStaffID" json:"submitter"`
	SubmittedAt     time.Time    `json:"submittedAt"`
	ReviewStaffID   *int64       `gorm:"column:review_staff_id" json:"-"`
	Reviewer        *staffMember `gorm:"foreignKey:ReviewStaffID" json:"reviewer,omitempty"`
	ReviewStartedAt *time.Time   `json:"reviewStartedAt,omitempty"`
	Status          string       `json:"status"`
	Notes           string       `json:"notes"`
	PublishedAt     *time.Time   `json:"publishedAt,omitempty"`
}

type asReviewsResponse struct {
	Total         int64                 `json:"total"`
	ViewerBaseURL string                `json:"viewerBaseURL"`
	Reviews       []archivesspaceReview `json:"submissions"`
}

func (svc *serviceContext) beginArchivesSpaceReview(c *gin.Context) {
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s in archivesspace begin review request", c.Param("id"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	reviewID, _ := strconv.ParseInt(c.Query("reviewer"), 10, 64)
	if reviewID == 0 {
		log.Printf("ERROR: invalid reviewer id %s in archivesspace begin review request", c.Param("reviewer"))
		c.String(http.StatusBadRequest, "invalid reviewer")
		return
	}

	var reviewUser staffMember
	err := svc.DB.Find(&reviewUser, reviewID).Error
	if err != nil {
		log.Printf("ERROR: unable to find reviewer %d: %s", reviewID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var asR archivesspaceReview
	err = svc.DB.Joins("Metadata").Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to find review for metadata %d: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if asR.ReviewStartedAt == nil {
		now := time.Now()
		asR.ReviewStartedAt = &now
	}
	asR.ReviewStaffID = &reviewID
	asR.Reviewer = &reviewUser
	asR.Status = "review"
	err = svc.DB.Save(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to begin as review for metadata %d: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, asR)
}

func (svc *serviceContext) resubmitArchivesSpaceReview(c *gin.Context) {
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s in archivesspace review request", c.Param("id"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	log.Printf("INFO: resubmit metadata %d for archivesspace review", mdID)
	var asR archivesspaceReview
	err := svc.DB.Joins("Metadata").Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to find review for metadata %d: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	asR.Status = "requested"
	err = svc.DB.Save(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable resubmit metadata %d for archivespace review: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func (svc *serviceContext) requestArchivesSpaceReview(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user"), 10, 64)
	mdID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if mdID == 0 {
		log.Printf("ERROR: invalid metadata id %s in archivesspace review request", c.Param("id"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	if userID == 0 {
		log.Printf("ERROR: invalid user id %s in archivesspace review request", c.Param("user"))
		c.String(http.StatusBadRequest, "invalid id")
		return
	}
	log.Printf("INFO: user %d requests archivesspace review for metadata %d", userID, mdID)

	var submitter staffMember
	err := svc.DB.First(&submitter, userID).Error
	if err != nil {
		log.Printf("ERROR: unable to load staffmember %d: %s", userID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var mdRec metadata
	err = svc.DB.First(&mdRec, mdID).Error
	if err != nil {
		log.Printf("ERROR: unable to load metadata %d: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	asReview := archivesspaceReview{SubmitStaffID: userID, Submitter: submitter, MetadataID: mdID, SubmittedAt: time.Now(), Status: "requested"}
	err = svc.DB.Save(&asReview).Error
	if err != nil {
		log.Printf("ERROR: user %d unable to request archives spaces review for %d: %s", userID, mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, asReview)
}

func (svc *serviceContext) getArchivesSpaceReviews(c *gin.Context) {
	resp := asReviewsResponse{ViewerBaseURL: fmt.Sprintf("%s/view", svc.ExternalSystems.Curio)}
	countQ := "select count(id) as total from archivesspace_reviews"
	err := svc.DB.Raw(countQ).Scan(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get archivesspace reviews count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = svc.DB.Joins("Submitter").Joins("Reviewer").Joins("Metadata").Where("published_at is null").Order("submitted_at asc").Find(&resp.Reviews).Error
	if err != nil {
		log.Printf("ERROR: unable to get archivesspace reviews: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
