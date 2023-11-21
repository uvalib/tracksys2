package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (svc *serviceContext) getCollectionAPTrustStatus(c *gin.Context) {
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
		c.String(http.StatusBadRequest, fmt.Sprintf("collection %d has not been submitted tp aptrust", collMD.ID))
		return
	}

	var collectionMemberIDs []struct {
		ID    int64
		PID   string `gorm:"column:pid"`
		Title string
	}
	err = svc.DB.Raw("select id, pid, title from metadata where parent_metadata_id=?", collectionID).Scan(&collectionMemberIDs).Error
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

	for _, member := range collectionMemberIDs {
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
				break
			}
		}
		if !found {
			log.Printf("INFO: metadata %d missing from collection %s aptrust status response", member.ID, collectionID)
			// 		log.Printf("INFO: metadata %d missing from collection %s aptrust status response; look up separately", member.ID, collectionID)
			// 		resp, err := svc.requestAPTStatus(member.ID)
			// 		if err != nil {
			// 			log.Printf("ERROR: unable to get aptrust status for %d: %s", member.ID, err.Error())
			// 		} else {
			// 			fail := apTrustCollectionResult{APTrustResult: resp, MetadataID: member.ID, MetadataPID: member.PID, MetadataTitle: member.Title}
			// 			parsedStatusList = append(parsedStatusList, &fail)
			// 		}
		}
	}

	c.JSON(http.StatusOK, parsedStatusList)
}

func (svc *serviceContext) requestAPTStatus(mdID int64) (*APTrustResult, error) {
	raw, getErr := svc.getRequest(fmt.Sprintf("%s/metadata/%d/aptrust", svc.ExternalSystems.Jobs, mdID))
	if getErr != nil {
		return nil, fmt.Errorf(getErr.Message)
	}

	var parsedStatus APTrustResult
	err := json.Unmarshal(raw, &parsedStatus)
	if err != nil {
		return nil, err
	}

	return &parsedStatus, nil
}

func (svc *serviceContext) getAPTrustStatus(md *metadata) (*apTrustStatus, error) {
	log.Printf("INFO: check aptrust status for metadata %d", md.ID)

	var aptSubmission apTrustSubmission
	err := svc.DB.Where("metadata_id=?", md.ID).Limit(1).Find(&aptSubmission).Error
	if err != nil {
		return nil, fmt.Errorf("unable to get submission info: %s", err.Error())
	}

	// if ID is zeo, there is no submission record so the item has not yet been submitted
	if aptSubmission.ID == 0 {
		return nil, nil
	}

	if md.IsCollection {
		log.Printf("INFO: metadata %d is a collection; get overall submission status from the ts database, not aptrust", md.ID)
		out := apTrustStatus{RequestedAt: aptSubmission.RequestedAt, SubmittedAt: aptSubmission.SubmittedAt,
			FinishedAt: aptSubmission.ProcessedAt, Status: "Failed"}
		if aptSubmission.Success {
			out.Status = "Success"
		}
		return &out, nil
	}

	log.Printf("INFO: get status for metadata %d from aptrust", md.ID)
	raw, getErr := svc.getRequest(fmt.Sprintf("%s/metadata/%d/aptrust", svc.ExternalSystems.Jobs, md.ID))
	if getErr != nil {
		if getErr.StatusCode == 404 {
			out := apTrustStatus{Bag: aptSubmission.Bag,
				Status:      "Failed",
				Note:        "Bagging or submission failed; check job status logs for more details",
				RequestedAt: aptSubmission.RequestedAt, SubmittedAt: aptSubmission.SubmittedAt, FinishedAt: aptSubmission.ProcessedAt}
			return &out, nil
		}
		return nil, fmt.Errorf("%d:%s", getErr.StatusCode, getErr.Message)
	}

	var parsedStatus APTrustResult
	err = json.Unmarshal(raw, &parsedStatus)
	if err != nil {
		return nil, fmt.Errorf("malformed status: %s", err.Error())
	}

	// merge TS submit data and APT status info into one record and return it
	out := apTrustStatus{ID: parsedStatus.ID, Bag: aptSubmission.Bag, ETag: parsedStatus.ETag, ObjectIdentifier: parsedStatus.ObjectIdentifier,
		StorageOption: parsedStatus.StorageOption, Status: parsedStatus.Status, Note: parsedStatus.Note, RequestedAt: aptSubmission.RequestedAt,
		SubmittedAt: aptSubmission.SubmittedAt, FinishedAt: aptSubmission.ProcessedAt, GroupIdentifier: parsedStatus.GroupIdentifier}
	if parsedStatus.ProcessedAt != "0001-01-01T00:00:00Z" {
		finishedAt, err := time.Parse("2006-01-02T15:04:05Z", parsedStatus.ProcessedAt)
		if err != nil {
			log.Printf("ERROR: unable to parse aptrust finished time: %s", err.Error())
		} else {
			out.FinishedAt = &finishedAt
		}
	}
	return &out, nil
}
