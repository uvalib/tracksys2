package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type auditDetail struct {
	Label string `json:"label"`
	Total int64  `json:"total"`
}

type auditStats struct {
	TotalAudited int64         `json:"totalAudited"`
	Results      []auditDetail `json:"results"`
}

func (svc *serviceContext) getAuditReport(c *gin.Context) {
	tgtYear := c.Query("year")
	if tgtYear == "" {
		tgtYear = "all"
	}
	log.Printf("INFO: get audit results for year [%s]", tgtYear)
	resp := auditStats{TotalAudited: 0, Results: make([]auditDetail, 0)}

	log.Printf("INFO: lookup total audit count")
	auditQ := svc.DB.Table("master_file_audits")
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err := auditQ.Count(&resp.TotalAudited).Error
	if err != nil {
		log.Printf("ERROR: unable to get total audit count for year[%s]: %s", tgtYear, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: lookup successful audit count")
	successData := auditDetail{Label: "No Errors", Total: 0}
	auditQ = svc.DB.Table("master_file_audits").Where("archive_exists=? and checksum_exists=? and checksum_match=? and iiif_exists=?", 1, 1, 1, 1)
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err = auditQ.Count(&successData.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get successful audit count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp.Results = append(resp.Results, successData)

	log.Printf("INFO: lookup missing archive count")
	noArchive := auditDetail{Label: "Missing Archive", Total: 0}
	auditQ = svc.DB.Table("master_file_audits").Where("archive_exists=?", 0)
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err = auditQ.Count(&noArchive.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get missing archive count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp.Results = append(resp.Results, noArchive)

	log.Printf("INFO: lookup missing checksum count")
	noChecksum := auditDetail{Label: "Missing Checksum", Total: 0}
	auditQ = svc.DB.Table("master_file_audits").Where("checksum_exists=?", 0)
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err = auditQ.Count(&noChecksum.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get missing checksum count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp.Results = append(resp.Results, noChecksum)

	log.Printf("INFO: lookup missing iiif count")
	noIIIF := auditDetail{Label: "Missing IIIF", Total: 0}
	auditQ = svc.DB.Table("master_file_audits").Where("iiif_exists=?", 0)
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err = auditQ.Count(&noIIIF.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get missing iiif count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp.Results = append(resp.Results, noIIIF)

	log.Printf("INFO: lookup failed checksum count")
	badChecksum := auditDetail{Label: "Checksum Mismatch", Total: 0}
	auditQ = svc.DB.Table("master_file_audits").Where("archive_exists=? and checksum_exists=? and checksum_match=?", 1, 1, 0)
	if tgtYear != "all" {
		auditQ = auditQ.Joins("inner join master_files mf on mf.id=master_file_id").Where("year(mf.created_at) = ?", tgtYear)
	}
	err = auditQ.Count(&badChecksum.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get missing iiif count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp.Results = append(resp.Results, badChecksum)

	c.JSON(http.StatusOK, resp)
}
