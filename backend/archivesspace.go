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
	ReviewStaffID   int64        `gorm:"column:review_staff_id" json:"-"`
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

type asRequest struct {
	ReviewInfo archivesspaceReview
	Staff      staffMember
}

func (svc *serviceContext) beginArchivesSpaceReview(c *gin.Context) {
	asR, err := svc.validateASRequest(c.Param("id"), c.Query("user"))
	if err != nil {
		log.Printf("ERROR: as review request failed: %s", err.Error())
	}
	log.Printf("INFO: user %d:%s requets archivesspace review for metadata %d", asR.Staff.ID, asR.Staff.ComputingID, asR.ReviewInfo.MetadataID)

	if asR.ReviewInfo.ReviewStartedAt == nil {
		now := time.Now()
		asR.ReviewInfo.ReviewStartedAt = &now
	}
	asR.ReviewInfo.ReviewStaffID = int64(asR.Staff.ID)
	asR.ReviewInfo.Status = "review"
	err = svc.DB.Save(&asR.ReviewInfo).Error
	if err != nil {
		log.Printf("ERROR: unable to begin as review for metadata %d: %s", asR.ReviewInfo.MetadataID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	asR.ReviewInfo.Reviewer = &asR.Staff
	c.JSON(http.StatusOK, asR.ReviewInfo)
}

func (svc *serviceContext) publishArchivesSpace(c *gin.Context) {
	reqInfo, err := svc.validateASRequest(c.Param("id"), c.Query("user"))
	if err != nil {
		log.Printf("ERROR: as publish request failed %s", err.Error())
	}
	log.Printf("INFO: user %s requests archivesspace publish for metadata %d to ", reqInfo.Staff.ComputingID, reqInfo.ReviewInfo.MetadataID)

	pubPayload := struct {
		UserID     int64 `json:"userID"`
		MetadataID int64 `json:"metadataID"`
	}{
		UserID:     int64(reqInfo.Staff.ID),
		MetadataID: reqInfo.ReviewInfo.MetadataID,
	}
	url := fmt.Sprintf("%s/archivesspace/publish", svc.ExternalSystems.Jobs)
	_, asErr := svc.postJSON(url, pubPayload)
	if asErr != nil {
		log.Printf("ERROR: unable to publish metadata %d: %d %s", reqInfo.ReviewInfo.MetadataID, asErr.StatusCode, asErr.Message)
		c.String(asErr.StatusCode, asErr.Message)
		return
	}

	now := time.Now()
	reqInfo.ReviewInfo.Status = "published"
	reqInfo.ReviewInfo.PublishedAt = &now
	err = svc.DB.Save(&reqInfo.ReviewInfo).Error
	if err != nil {
		log.Printf("ERROR: unable to update metadata published status for metadata %d: %s", reqInfo.ReviewInfo.MetadataID, err.Error())
	}
	c.String(http.StatusOK, "published")
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

func (svc *serviceContext) rejectArchivesSpaceSubmission(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received archivesspace reject request for metadata %s", mdID)

	var req struct {
		UserID int64  `json:"userID"`
		Notes  string `json:"notes"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid archivesspace reject request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("INFO: user %d rejects archivesspace submission %s with notes [%s]", req.UserID, mdID, req.Notes)
	var asR archivesspaceReview
	err = svc.DB.Joins("Metadata").Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to load submission info for metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	asR.Status = "rejected"
	asR.Notes = req.Notes
	err = svc.DB.Save(&asR).Error
	if err != nil {
		log.Printf("ERROR: reject as submission %s failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, asR)
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

func (svc *serviceContext) cancelArchivesSpaceSubmission(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received archivesspace submission cancel request for metadata %s", mdID)

	var asR archivesspaceReview
	err := svc.DB.Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to load submission info for metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = svc.DB.Delete(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to cancel archivesspace submission for metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "deleted")
}

func (svc *serviceContext) updateArchivesSpaceSubmissionNotes(c *gin.Context) {
	mdID := c.Param("id")
	log.Printf("INFO: received archivesspace update notes for metadata %s", mdID)

	var req struct {
		Notes string `json:"notes"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid archivesspace notes request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var asR archivesspaceReview
	err = svc.DB.Joins("Metadata").Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		log.Printf("ERROR: unable to load submission info for metadata %s: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	asR.Notes = req.Notes
	err = svc.DB.Save(&asR).Error
	if err != nil {
		log.Printf("ERROR: update notes for as submission %s failed: %s", mdID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, asR)
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

func (svc *serviceContext) validateASRequest(mdParam, userParam string) (*asRequest, error) {
	mdID, _ := strconv.ParseInt(mdParam, 10, 64)
	if mdID == 0 {
		return nil, fmt.Errorf("invalid metadata id %s", mdParam)
	}
	userID, _ := strconv.ParseInt(userParam, 10, 64)
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id %s", userParam)
	}

	var reqUser staffMember
	err := svc.DB.Find(&reqUser, userID).Error
	if err != nil {
		return nil, fmt.Errorf("unable to get user id %d: %s", userID, err.Error())
	}

	var asR archivesspaceReview
	err = svc.DB.Joins("Metadata").Where("metadata_id=?", mdID).First(&asR).Error
	if err != nil {
		return nil, fmt.Errorf("user %s is unable to get as reciew recird for metadata %d: %s", reqUser.ComputingID, mdID, err.Error())
	}

	out := asRequest{Staff: reqUser, ReviewInfo: asR}
	return &out, nil
}
