package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// apTrustSubmissionis a TrackSys DB record generated when a metadata record is submitted to APTrust
type apTrustSubmission struct {
	ID          int64      `json:"-"`
	MetadataID  int64      `gorm:"column:metadata_id" json:"-"`
	Bag         string     `json:"bag"`
	RequestedAt time.Time  `json:"requestedAt"`
	SubmittedAt *time.Time `json:"submittedAt"`
	ProcessedAt *time.Time `json:"processedAt"`
	Success     bool       `json:"success"`
}

// APTrustResult is the JSON workitem status object retruned by the APTrust API
type APTrustResult struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ETag             string `json:"etag"`
	ObjectIdentifier string `json:"object_identifier"`
	GroupIdentifier  string `json:"bag_group_identifier"`
	AltIdentifier    string `json:"alt_identifier"`
	StorageOption    string `json:"storage_option"`
	Note             string `json:"note"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	ProcessedAt      string `json:"date_processed"`
}

type apTrustCollectionResult struct {
	*APTrustResult
	MetadataID    int64  `json:"metadata_id"`
	MetadataPID   string `json:"metadata_pid"`
	MetadataTitle string `json:"metadata_title"`
}

type apTrustSubmissionsResonse struct {
	Total       int64 `json:"total"`
	Submissions []struct {
		ID          int64      `json:"id"`
		MetadataID  int64      `gorm:"column:metadata_id" json:"metadataID"`
		PID         string     `gorm:"column:pid" json:"pid"`
		Title       string     `json:"title"`
		RequestedAt time.Time  `json:"requestedAt"`
		ProcessedAt *time.Time `json:"processedAt"`
		Success     bool       `json:"success"`
	} `json:"submissions"`
}

// apTrust status is a combination of the JSON APTrust work item record from the API and the
// TrackSys DB submission record. It provides the complete picture of a submission
type apTrustStatus struct {
	ID               int64      `json:"id"`
	Bag              string     `json:"bag"`
	ETag             string     `json:"etag"`
	ObjectIdentifier string     `json:"objectIdentifier"`
	GroupIdentifier  string     `json:"groupIdentifier"`
	StorageOption    string     `json:"storage"`
	Note             string     `json:"note"`
	Status           string     `json:"status"`
	RequestedAt      time.Time  `json:"requestedAt"`
	SubmittedAt      *time.Time `json:"submittedAt"`
	FinishedAt       *time.Time `json:"finishedAt"`
}

func (svc *serviceContext) getAPTrustSubmissions(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if pageSize == 0 {
		pageSize = 30
	}
	queryStr := c.Query("q")

	sortBy := c.Query("by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := c.Query("order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sortField := fmt.Sprintf("apt.%s", sortBy)
	if sortBy == "requestedAt" {
		sortField = "apt.requested_at"
	} else if sortBy == "processedAt" {
		sortField = "apt.processed_at"
	} else if sortBy == "pid" {
		sortField = "m.pid"
	} else if sortBy == "title" {
		sortField = "m.title"
	}

	orderStr := fmt.Sprintf("%s %s", sortField, sortOrder)
	log.Printf("INFO: get %d aptrust submissions starting from offset %d order %s; query=[%s]", pageSize, startIndex, orderStr, queryStr)

	resp := apTrustSubmissionsResonse{}
	joinQ := "inner join metadata m on m.id = metadata_id"
	countQ := "select count(apt.id) as total from ap_trust_submissions apt " + joinQ
	countQ += " where is_collection=0"
	err := svc.DB.Raw(countQ).Scan(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get aptrust submissions count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	searchQ := "select apt.id as id, m.id as metadata_id, m.pid as pid, m.title as title, requested_at, processed_at, success from ap_trust_submissions apt "
	searchQ += joinQ
	if queryStr != "" {
		searchQ += fmt.Sprintf(" where is_collection=0 and m.title like '%%%s%%'", queryStr)
	} else {
		searchQ += " where is_collection=0"
	}
	searchQ += fmt.Sprintf(" order by %s", orderStr)
	searchQ += fmt.Sprintf(" limit %d,%d", startIndex, pageSize)
	err = svc.DB.Debug().Raw(searchQ).Scan(&resp.Submissions).Error
	if err != nil {
		log.Printf("ERROR: unable to get aptrust submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getAPTrustMetadataStatus(c *gin.Context) {
	metadataID := c.Param("id")
	log.Printf("INFO: get metadata %s aptrust status", metadataID)
	var mdRec metadata
	err := svc.DB.Joins("APTrustSubmission").Find(&mdRec, metadataID).Error
	if err != nil {
		log.Printf("ERROR: unable to load collection %s: %s", metadataID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if mdRec.APTrustSubmission == nil {
		log.Printf("INFO: metadata %d has not been submitted to aptrust", mdRec.ID)
		c.String(http.StatusBadRequest, fmt.Sprintf("metadata %d has not been submitted to aptrust", mdRec.ID))
		return
	}

	if mdRec.IsCollection {
		log.Printf("INFO: metadata %d is a collection; get overall submission status from the ts database, not aptrust", mdRec.ID)
		out := apTrustStatus{RequestedAt: mdRec.APTrustSubmission.RequestedAt, SubmittedAt: mdRec.APTrustSubmission.SubmittedAt,
			FinishedAt: mdRec.APTrustSubmission.ProcessedAt, Status: "Failed"}
		if mdRec.APTrustSubmission.Success {
			out.Status = "Success"
		}
		c.JSON(http.StatusOK, out)
		return
	}

	raw, getErr := svc.getRequest(fmt.Sprintf("%s/metadata/%d/aptrust", svc.ExternalSystems.Jobs, mdRec.ID))
	if getErr != nil {
		if getErr.StatusCode == 404 {
			log.Printf("INFO: no aptrust status found for metadata %d", mdRec.ID)
			out := apTrustStatus{Bag: mdRec.APTrustSubmission.Bag,
				Status:      "Failed",
				Note:        "Bagging or submission failed; check job status logs for more details",
				RequestedAt: mdRec.APTrustSubmission.RequestedAt, SubmittedAt: mdRec.APTrustSubmission.SubmittedAt}
			c.JSON(http.StatusOK, out)
			return
		}
		log.Printf("ERROR: unable to get aptrust status: %s", getErr.Message)
		c.String(http.StatusInternalServerError, fmt.Sprintf("%d:%s", getErr.StatusCode, getErr.Message))
		return
	}

	var parsedStatus APTrustResult
	err = json.Unmarshal(raw, &parsedStatus)
	if err != nil {
		log.Printf("ERROR: unable to parse aptrust status: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// merge TS submit data and APT status info into one record and return it
	out := apTrustStatus{ID: parsedStatus.ID, Bag: mdRec.APTrustSubmission.Bag, ETag: parsedStatus.ETag,
		ObjectIdentifier: parsedStatus.ObjectIdentifier, StorageOption: parsedStatus.StorageOption,
		Status: parsedStatus.Status, Note: parsedStatus.Note, RequestedAt: mdRec.APTrustSubmission.RequestedAt,
		SubmittedAt: mdRec.APTrustSubmission.SubmittedAt, FinishedAt: mdRec.APTrustSubmission.ProcessedAt,
		GroupIdentifier: parsedStatus.GroupIdentifier}
	if parsedStatus.ProcessedAt != "0001-01-01T00:00:00Z" {
		finishedAt, err := time.Parse("2006-01-02T15:04:05Z", parsedStatus.ProcessedAt)
		if err != nil {
			log.Printf("ERROR: unable to parse aptrust finished time: %s", err.Error())
		} else {
			out.FinishedAt = &finishedAt
		}
	}

	// If the status is finished but TS submit record has not been updated, update it now
	if mdRec.APTrustSubmission.ProcessedAt == nil && (parsedStatus.Status == "Success" || parsedStatus.Status == "Failed" || parsedStatus.Status == "Canceled") {
		mdRec.APTrustSubmission.Success = (parsedStatus.Status == "Success")
		mdRec.APTrustSubmission.ProcessedAt = out.FinishedAt
		err = svc.DB.Save(&mdRec.APTrustSubmission).Error
		if err != nil {
			log.Printf("ERROR: update aptrust status for %d failed: %s", mdRec.ID, err.Error())
		}
	}
	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getAPTrustCollectionStatus(c *gin.Context) {
	collectionID := c.Param("id")
	log.Printf("INFO: get collection %s aptrust status", collectionID)
	var collMD metadata
	err := svc.DB.Joins("APTrustSubmission").Find(&collMD, collectionID).Error
	if err != nil {
		log.Printf("ERROR: unable to load collection %s: %s", collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if collMD.APTrustSubmission == nil {
		log.Printf("INFO: collection %d has not been submitted to aptrust", collMD.ID)
		c.String(http.StatusBadRequest, fmt.Sprintf("collection %d has not been submitted to aptrust", collMD.ID))
		return
	}

	var collectionMemberInfo []struct {
		ID           int64
		PID          string `gorm:"column:pid"`
		Title        string
		SubmissionID int64 `gorm:"column:submission_id"`
		ProcessedAt  *time.Time
		Success      bool
	}
	err = svc.DB.Raw(`
	select metadata.id, pid, title, ap_trust_submissions.id as submission_id, processed_at, success from metadata
	inner join ap_trust_submissions on metadata_id=metadata.id where parent_metadata_id=?`, collectionID).
		Scan(&collectionMemberInfo).Error
	if err != nil {
		log.Printf("ERROR: unable to get collection %s member id list: %s", collectionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	raw, getErr := svc.getRequest(fmt.Sprintf("%s/metadata/%s/aptrust", svc.ExternalSystems.Jobs, collectionID))
	if getErr != nil {
		log.Printf("ERROR: unable to get collection %s aptrust status: %s", collectionID, getErr.Message)
		c.String(http.StatusInternalServerError, getErr.Message)
		return
	}

	var parsedStatusList []*apTrustCollectionResult
	err = json.Unmarshal(raw, &parsedStatusList)
	if err != nil {
		log.Printf("ERROR: unable to parse aptrust response: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// walk the list of aptrust submissions from the DB abd compare with the APTrust status response
	for _, member := range collectionMemberInfo {
		found := false
		strID := fmt.Sprintf("%d", member.ID)
		for _, aptStatus := range parsedStatusList {
			idBits := strings.Split(aptStatus.ObjectIdentifier, "-")
			objID := idBits[len(idBits)-1]
			if strID == objID {
				found = true
				aptStatus.MetadataID = member.ID
				aptStatus.MetadataPID = member.PID
				aptStatus.MetadataTitle = member.Title
				if member.ProcessedAt == nil && (aptStatus.Status == "Success" || aptStatus.Status == "Failed" || aptStatus.Status == "Canceled") {
					log.Printf("INFO: metadata %d has completed processing; update submission record", member.ID)
					finishedAt, err := time.Parse("2006-01-02T15:04:05Z", aptStatus.ProcessedAt)
					if err != nil {
						log.Printf("ERROR: unable to parse finished date %s for metadata %d", aptStatus.ProcessedAt, member.ID)
					} else {
						update := apTrustSubmission{ID: member.SubmissionID, ProcessedAt: &finishedAt, Success: aptStatus.Status == "Success"}
						err = svc.DB.Model(&update).Select("ProcessedAt", "Success").Updates(update).Error
						if err != nil {
							log.Printf("ERROR: unable to update status for metadata %d: %s", member.ID, err.Error())
						}
					}
				}
				break
			}
		}
		if !found {
			log.Printf("INFO: metadata %d missing from collection %s aptrust status response", member.ID, collectionID)
			// add a failed status
			failedReslt := APTrustResult{Status: "Failed", Note: "Bagging failed; see job logs for details"}
			fail := apTrustCollectionResult{APTrustResult: &failedReslt, MetadataID: member.ID, MetadataPID: member.PID, MetadataTitle: member.Title}
			parsedStatusList = append(parsedStatusList, &fail)
		}
	}

	c.JSON(http.StatusOK, parsedStatusList)
}
