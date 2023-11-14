package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
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

// apTrustResult is the JSON workitem status object retruned by the APTrust API
type apTrustResult struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ETag             string `json:"etag"`
	ObjectIdentifier string `json:"object_identifier"`
	AltIdentifier    string `json:"alt_identifier"`
	StorageOption    string `json:"storage_option"`
	Note             string `json:"note"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	ProcessedAt      string `json:"date_processed"`
}

// apTrust status is a combination of the JSON APTrust work item record from the API and the
// TrackSys DB submission record. It provides the complete picture of a submission
type apTrustStatus struct {
	ID               int64      `json:"id"`
	Bag              string     `json:"bag"`
	ETag             string     `json:"etag"`
	ObjectIdentifier string     `json:"objectIdentifier"`
	StorageOption    string     `json:"storage"`
	Note             string     `json:"note"`
	Status           string     `json:"status"`
	RequestedAt      time.Time  `json:"requestedAt"`
	SubmittedAt      *time.Time `json:"submittedAt"`
	FinishedAt       *time.Time `json:"finishedAt"`
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

	var parsedStatus apTrustResult
	err = json.Unmarshal(raw, &parsedStatus)
	if err != nil {
		return nil, fmt.Errorf("malformed status: %s", err.Error())
	}

	// merge TS submit data and APT status info into one record and return it
	out := apTrustStatus{ID: parsedStatus.ID, Bag: aptSubmission.Bag, ETag: parsedStatus.ETag, ObjectIdentifier: parsedStatus.ObjectIdentifier,
		StorageOption: parsedStatus.StorageOption, Status: parsedStatus.Status, Note: parsedStatus.Note, RequestedAt: aptSubmission.RequestedAt,
		SubmittedAt: aptSubmission.SubmittedAt, FinishedAt: aptSubmission.ProcessedAt}
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
