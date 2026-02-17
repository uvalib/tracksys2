package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (svc *serviceContext) getDeliveryStats(c *gin.Context) {
	log.Printf("INFO: get patron deliveries statistics")
	tgtYear := c.Query("year")
	if tgtYear == "" {
		log.Printf("ERROR: year is required")
		c.String(http.StatusBadRequest, "year is required")
		return
	}

	// NOTE: in the response, the totals are an array of 12 counts. each count corresponds to a month
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	var resp struct {
		Months []string `json:"months"`
		Total  []int    `json:"total"`
		Late   []int    `json:"late"`
		OnTime []int    `json:"onTime"`
	}
	resp.Months = months

	type order struct {
		ID            int64
		DateDue       time.Time
		DateCompleted time.Time
	}
	var completedOrders []order
	err := svc.DB.Debug().Joins("inner join units u on order_id = orders.id").
		Where("date_patron_deliverables_ready is not null").
		Where("order_status=?", "completed").
		Where("date_completed like ?", fmt.Sprintf("%s%%", tgtYear)).
		Distinct("orders.id", "date_due", "date_completed").
		Order("date_completed asc").
		Find(&completedOrders).Error
	if err != nil {
		log.Printf("ERROR: unable to get deliveries report: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// orders are sorted by completion date asc. track the current month and add the
	// count to the correct slot in the total/late/ontime arrays.
	// NOTE: the logic below does not account for skipped months as this should never happen
	currMonth := 0
	total := 0
	late := 0
	ontime := 0
	for _, o := range completedOrders {
		orderMonth := int(o.DateCompleted.Month())
		if currMonth != orderMonth {
			if currMonth > 0 {
				resp.Total = append(resp.Total, total)
				resp.Late = append(resp.Late, late)
				resp.OnTime = append(resp.OnTime, ontime)
			}
			total = 0
			late = 0
			ontime = 0
			currMonth = orderMonth
		}
		total++
		if o.DateCompleted.Before(o.DateDue) {
			ontime++
		} else {
			late++
		}
	}

	if total > 0 {
		resp.Total = append(resp.Total, total)
		resp.Late = append(resp.Late, late)
		resp.OnTime = append(resp.OnTime, ontime)
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getImageStats(c *gin.Context) {
	log.Printf("INFO: get image statistics")
	dateQStr := c.Query("date")
	if (strings.Contains(dateQStr, "TO") || strings.Contains(dateQStr, "AFTER") || strings.Contains(dateQStr, "BEFORE")) == false {
		log.Printf("ERROR: invalid date query [%s]", dateQStr)
		c.String(http.StatusBadRequest, fmt.Sprintf("%s is not valid", dateQStr))
		return
	}

	var imageResp struct {
		Total int64 `json:"total"`
		DL    int64 `json:"dl"`
		DPLA  int64 `json:"dpla"`
	}

	// all, dl, dpla
	for i := 0; i < 3; i++ {
		cntQuery := svc.DB.Table("master_files")
		addDateConstraint(cntQuery, "master_files.created_at", dateQStr)

		count := &imageResp.Total
		switch i {
		case 1:
			count = &imageResp.DL
			cntQuery = cntQuery.Where("master_files.date_dl_ingest is not null")
			log.Printf("INFO: get DL images")
		case 2:
			count = &imageResp.DPLA
			cntQuery = cntQuery.Joins("inner join metadata m on metadata_id=m.id and m.dpla=1")
			log.Printf("INFO: get DPLA images")
		default:
			log.Printf("INFO: get all images")
		}

		err := cntQuery.Count(count).Error
		if err != nil {
			log.Printf("ERROR: unable to image counts: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, imageResp)
}

func (svc *serviceContext) getMetadataStats(c *gin.Context) {
	log.Printf("INFO: get metadata statistics")
	dateQStr := c.Query("date")
	if (strings.Contains(dateQStr, "TO") || strings.Contains(dateQStr, "AFTER") || strings.Contains(dateQStr, "BEFORE")) == false {
		log.Printf("ERROR: invalid date query [%s]", dateQStr)
		c.String(http.StatusBadRequest, fmt.Sprintf("%s is not valid", dateQStr))
		return
	}

	type metadata struct {
		ID           int64
		Type         string
		DateDLIngest *time.Time `gorm:"column:date_dl_ingest"`
		DPLA         bool       `gorm:"column:dpla"`
	}
	type metadataDetail struct {
		Total int64 `json:"total"`
		SIRSI int64 `json:"sirsi"`
		XML   int64 `json:"xml"`
	}
	var metadataResp struct {
		All  metadataDetail `json:"all"`
		DL   metadataDetail `json:"DL"`
		DPLA metadataDetail `json:"DPLA"`
	}

	// get all the MD recs in the requested date range...
	var mdRecs []metadata
	mdQ := svc.DB.Select("id", "type", "date_dl_ingest", "dpla")
	addDateConstraint(mdQ, "created_at", dateQStr)
	err := mdQ.Find(&mdRecs).Error
	if err != nil {
		log.Printf("ERROR: unable to get metadata statistics: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// ...then split up the counts based on metadata type, DL and DPLA
	for _, md := range mdRecs {
		metadataResp.All.Total++
		if md.DPLA {
			metadataResp.DPLA.Total++
		}
		if md.DateDLIngest != nil {
			metadataResp.DL.Total++
		}

		switch md.Type {
		case "SirsiMetadata":
			metadataResp.All.SIRSI++
			if md.DPLA {
				metadataResp.DPLA.SIRSI++
			}
			if md.DateDLIngest != nil {
				metadataResp.DL.SIRSI++
			}
		case "XmlMetadata":
			metadataResp.All.XML++
			if md.DPLA {
				metadataResp.DPLA.XML++
			}
			if md.DateDLIngest != nil {
				metadataResp.DL.XML++
			}
		}
	}
	c.JSON(http.StatusOK, metadataResp)
}

func (svc *serviceContext) getPublishedStats(c *gin.Context) {
	log.Printf("INFO: get published statistics")

	type publishedVirgo struct {
		ID       uint64 `json:"id"`
		PID      string `gorm:"column:pid" json:"pid"`
		Title    string `json:"title"`
		ThumbURL string `gorm:"-" json:"thumbURL"`
	}
	type publishedAS struct {
		ID          uint64 `json:"id"`
		ExternalURL string `gorm:"column:external_uri" json:"externalURL"`
		Title       string `json:"title"`
	}
	type respData struct {
		Virgo         []*publishedVirgo `json:"virgo"`
		ArchivesSpace []*publishedAS    `json:"archivesSpace"`
	}
	var resp respData
	err := svc.DB.Table("metadata").Where("date_dl_ingest is not null").Where("type=?", "SirsiMetadata").
		Limit(20).Order("date_dl_ingest DESC").Find(&resp.Virgo).Error
	if err != nil {
		log.Printf("ERROR: unable to get recent virgo publication stats: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range resp.Virgo {
		v.ThumbURL = svc.getExemplarThumbURL(v.ID)
	}

	err = svc.DB.Table("metadata").Where("date_dl_ingest is not null").
		Where("type=?", "ExternalMetadata").Where("external_system_id=?", 1).
		Limit(20).Order("date_dl_ingest DESC").Find(&resp.ArchivesSpace).Error
	if err != nil {
		log.Printf("ERROR: unable to get recent archivesspace publication stats: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range resp.ArchivesSpace {
		v.ExternalURL = fmt.Sprintf("https://archives.lib.virginia.edu%s", v.ExternalURL)

	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) getExemplarThumbURL(mdID uint64) string {
	log.Printf("INFO: get exemplar for %d", mdID)
	var mf masterFile
	err := svc.DB.Preload("ImageTechMeta").Where("metadata_id=? and exemplar=1", mdID).First(&mf).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false {
			log.Printf("ERROR: unable to get examplar for %d: %s", mdID, err.Error())
			return ""
		}
		log.Printf("INFO: no exemplar set for metadata id %d; choosing first masterfile", mdID)
		err = svc.DB.Preload("ImageTechMeta").Where("metadata_id=?", mdID).Order("filename asc").First(&mf).Error
		if err != nil {
			log.Printf("ERROR: unable to get examplar for %d: %s", mdID, err.Error())
			return ""
		}
	}

	log.Printf("%+v", mf)

	// orientation is enum type: none: 0, flip_y_axis: 1, rotate90: 2, rotate180: 3, rotate270
	rotations := []string{"0", "!0", "90", "180", "270"}
	exemplarURL := fmt.Sprintf("%s/%s/full/!125,200/%s/default.jpg", svc.ExternalSystems.IIIF, mf.PID, rotations[mf.ImageTechMeta.Orientation])
	return exemplarURL
}

func (svc *serviceContext) getStorageStats(c *gin.Context) {
	log.Printf("INFO: get storage statistics")

	var storageResp struct {
		All float64 `json:"total"`
		DL  float64 `json:"dl"`
	}

	// all, dl
	for i := 0; i < 2; i++ {
		szQuery := svc.DB.Table("master_files")
		sizeGB := &storageResp.All
		if i == 1 {
			sizeGB = &storageResp.DL
			szQuery = szQuery.Where("master_files.date_dl_ingest is not null")
			log.Printf("INFO: get DL images size")
		} else {
			log.Printf("INFO: get all images size")
		}

		err := szQuery.Select("sum(filesize)/1073741824.0 as size_gb").Row().Scan(sizeGB)
		if err != nil {
			log.Printf("ERROR: unable to image size: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, storageResp)
}

func (svc *serviceContext) getArchiveStats(c *gin.Context) {
	log.Printf("INFO: get archive statistics")
	dateQStr := c.Query("date")
	if (strings.Contains(dateQStr, "TO") || strings.Contains(dateQStr, "AFTER") || strings.Contains(dateQStr, "BEFORE")) == false {
		log.Printf("ERROR: invalid date query [%s]", dateQStr)
		c.String(http.StatusBadRequest, fmt.Sprintf("%s is not valid", dateQStr))
		return
	}

	var archiveResp struct {
		Bound      int   `json:"bound"`
		Manuscript int64 `json:"manuscript"`
		Photo      int64 `json:"photo"`
	}

	log.Printf("INFO: get archive stats for bound items")
	var boundIDs []int64
	boundCntQ := svc.DB.Table("master_files").
		Joins("inner join metadata m on metadata_id = m.id").
		Select("m.id").
		Where("master_files.title=?", "Spine")
	addDateConstraint(boundCntQ, "date_archived", dateQStr)

	err := boundCntQ.Find(&boundIDs).Error
	if err != nil {
		log.Printf("ERROR: unable to get archived bound counts: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	archiveResp.Bound = len(boundIDs)

	log.Printf("INFO: get archive stats for photo/av items")
	photoCntQ := svc.DB.Table("master_files").
		Joins("inner join metadata m on master_files.metadata_id = m.id").
		Joins("inner join units u on u.id = master_files.unit_id").
		Joins("inner join orders o on o.id = u.order_id").
		Joins("inner join agencies a on a.id = o.agency_id")

	// Negatives orders and Fine arts agency
	photoCntQ.Where("(o.id = 8126 or o.id = 8125 or agency_id = 37)")

	// Not unbound sheets
	photoCntQ.Where("(m.call_number is null or (m.call_number is not null and m.call_number not like 'RG-%'))")
	addDateConstraint(photoCntQ, "master_files.date_archived", dateQStr)
	err = photoCntQ.Count(&archiveResp.Photo).Error
	if err != nil {
		log.Printf("ERROR: unable to get archived photo counts: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: get archive stats for manuscript items")
	unboundCntQ := svc.DB.Table("master_files").
		Joins("inner join metadata m on metadata_id = m.id").
		Where("(m.call_number like 'MSS%' or m.call_number like 'RG-%')"). // manuscript call numbers
		Where("master_files.title not like '%verso'").                     // skip back sides of pages
		Where("m.id != 3009").                                             // no visual history
		Where(                                                             // take numbered pages or pages that look like standard parts of MSS
			svc.DB.Where("master_files.title regexp '^[[:digit:]]+'").
				Or("master_files.title like 'front%'").
				Or("master_files.title like 'rear%'").
				Or("master_files.title like 'back%'").
				Or("master_files.title like 'title%'").
				Or("master_files.title like 'table%'").
				Or("master_files.title like 'blank%'").
				Or("master_files.title regexp '^(IX|IV|V?I{0,3})$'"),
		).
		Not("m.id in ?", boundIDs) // skip files that are part of a bound volume
	addDateConstraint(unboundCntQ, "master_files.date_archived", dateQStr)
	err = unboundCntQ.Count(&archiveResp.Manuscript).Error
	if err != nil {
		log.Printf("ERROR: unable to get manuscript photo counts: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, archiveResp)
}

func addDateConstraint(baseQ *gorm.DB, fieldName string, dateQStr string) {
	if strings.Contains(dateQStr, "TO") {
		bits := strings.Split(dateQStr, " ")
		baseQ.Where(fmt.Sprintf("%s >= ? and %s <= ?", fieldName, fieldName), bits[0], bits[2])
	} else if strings.Contains(dateQStr, "AFTER") {
		bits := strings.Split(dateQStr, " ")
		baseQ.Where(fmt.Sprintf("%s >= ?", fieldName), bits[1])
	} else if strings.Contains(dateQStr, "BEFORE") {
		bits := strings.Split(dateQStr, " ")
		baseQ.Where(fmt.Sprintf("%s <= ?", fieldName), bits[1])
	}
}
