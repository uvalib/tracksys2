package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type apTrustResult struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ETag             string `json:"etag"`
	ObjectIdentifier string `json:"object_identifier"`
	AltIdentifier    string `json:"alt_identifier"`
	StorageOption    string `json:"storage"`
	Note             string `json:"note"`
	Status           string `json:"status"`
	QueuedAt         string `json:"queued_at"`
	ProcessedAt      string `json:"date_processed"`
}

func (svc *serviceContext) updateAPTrustStatus(md *metadata) {
	log.Printf("INFO: check aptrust status for metadata %d", md.ID)

	// possible APTrust status values: Cancelled, Failed, Pending, Started, Success, Suspended
	// in addition to the APTrust defined statuses, TrackSys adds: Baggit (bagging in process), Submit (submitting bag to S3 for preservation)
	// status will be created/updated if:
	//    1. there is no status in the DB
	//    2. the DB status is not Success
	if md.APTrustStatus != nil && md.APTrustStatus.Status == "Success" {
		log.Printf("INFO: metadata %d has aptrust status success; nothing more to do", md.ID)
		return
	}

	log.Printf("INFO: get status for metadata %d from aptrust", md.ID)
	raw, getErr := svc.getRequest(fmt.Sprintf("%s/metadata/%d/aptrust", svc.ExternalSystems.Jobs, md.ID))
	if getErr != nil {
		if getErr.StatusCode == 404 {
			log.Printf("INFO: metadata %d has no aptrust status", md.ID)
		} else {
			log.Printf("ERROR: aptrust status request for metadata %d failed: %d:%s", md.ID, getErr.StatusCode, getErr.Message)
		}
		return
	}

	var parsedStatus apTrustResult
	err := json.Unmarshal(raw, &parsedStatus)
	if err != nil {
		log.Printf("ERROR: unable to parse aptrust status response: %s", err.Error())
		return
	}

	if md.APTrustStatus == nil {
		log.Printf("INFO: create new aptrust status record for metadata %d", md.ID)
		aptStatus := apTrustStatus{MetadataID: md.ID}
		md.APTrustStatus = &aptStatus
	} else {
		log.Printf("INFO: update aptrust status for metadata %d", md.ID)
	}

	md.APTrustStatus.Etag = parsedStatus.ETag
	md.APTrustStatus.ObjectID = parsedStatus.ObjectIdentifier
	md.APTrustStatus.Status = parsedStatus.Status
	md.APTrustStatus.Note = parsedStatus.Note

	// only update initial submit date if it is not set and there is valid data in the response
	if md.APTrustStatus.SubmittedAt.IsZero() && parsedStatus.QueuedAt != "0001-01-01T00:00:00Z" {
		parsedSubmit, err := time.Parse("2006-01-02T15:04:05Z", parsedStatus.QueuedAt)
		if err != nil {
			log.Printf("ERROR: unable to parse aptrust queued time %s, default to now: %s", parsedStatus.QueuedAt, err.Error())
			parsedSubmit = time.Now()
		}
		md.APTrustStatus.SubmittedAt = parsedSubmit
	}

	// only update processed date if there is valid data in teh response
	if parsedStatus.ProcessedAt != "0001-01-01T00:00:00Z" {
		parsedDone, err := time.Parse("2006-01-02T15:04:05Z", parsedStatus.ProcessedAt)
		if err != nil {
			log.Printf("ERROR: unable to parse aptrust processed time %s, default to now: %s", parsedStatus.ProcessedAt, err.Error())
			parsedDone = time.Now()
		}
		md.APTrustStatus.FinishedAt = &parsedDone
	}

	// NOTE: Save will create a record if there is no primary key set
	err = svc.DB.Save(md.APTrustStatus).Error
	if err != nil {
		log.Printf("ERROR: unable to save aptrust status record for metadata %d: %s", md.ID, err.Error())
	} else {
		log.Printf("INFO: save aptrust status record for metadata %d", md.ID)
	}
}
