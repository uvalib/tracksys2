package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (svc *serviceContext) cleanupCanceledUnits(c *gin.Context) {
	log.Printf("INFO: cleanup canceled units tha are older than 3 months")

	maxDel, _ := strconv.ParseUint(c.Query("max"), 10, 64)
	if maxDel == 0 {
		maxDel = 500
	}
	log.Printf("INFO: limit to %d deletions", maxDel)

	deleteThreshold := time.Now().AddDate(0, -3, 0)
	dateStr := deleteThreshold.Format("2006-01-02")

	var unitResp []struct {
		ID        int64
		UpdatedAt time.Time
		Reorder   bool
		FileCount int64
	}

	qStr := "select u.id, u.updated_at, reorder, count(f.id) as file_count from units u "
	qStr += " left join master_files f on f.unit_id = u.id  "
	qStr += " where unit_status=? and u.updated_at < ? group by u.id"
	if err := svc.DB.Raw(qStr, "canceled", dateStr).Scan(&unitResp).Error; err != nil {
		log.Printf("ERROR: find canceled units failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(unitResp) == 0 {
		log.Printf("INFO: there are no canceled units to cleanup")
		c.String(http.StatusOK, "done")
		return
	}

	log.Printf("INFO: found %d canceled units older than %s", len(unitResp), dateStr)
	go func() {
		delCnt := 0
		projCnt := 0
		failed := make([]int64, 0)
		hasFiles := make([]int64, 0)

		for _, u := range unitResp {
			log.Printf("INFO: delete unit %d canceled on %s", u.ID, u.UpdatedAt.Format("2006-01-02"))
			if u.FileCount > 0 {
				if u.Reorder {
					log.Printf("INFO: unit %d is a reorder with %d masterfiles; delete them first", u.ID, u.FileCount)
					var mfIDs []int64
					if err := svc.DB.Raw("select id from master_files where unit_id=?", u.ID).Scan(&mfIDs).Error; err != nil {
						log.Printf("ERROR: unable to get unit %d masterfiles: %s", u.ID, err.Error())
						continue
					}
					if err := svc.DB.Exec("delete from image_tech_meta where master_file_id in ?", mfIDs).Error; err != nil {
						log.Printf("ERROR: unable remove unit %d image tech metadata: %s", u.ID, err.Error())
						continue
					}
					delQ := "delete from master_files where unit_id=?"
					if err := svc.DB.Exec(delQ, u.ID).Error; err != nil {
						log.Printf("ERROR: unable to delete %d masterfiles for reorder unit %d: %s", u.FileCount, u.ID, err.Error())
						continue
					}
				} else {
					log.Printf("INFO: unit %d has %d masterfiles, not deleting", u.ID, u.FileCount)
					hasFiles = append(hasFiles, u.ID)
					continue
				}
			}

			projResp := svc.getUnitProject(u.ID)
			if projResp.Exists {
				log.Printf("INFO: canceled unit %d is associated with project %d; cancel it", u.ID, projResp.ProjectID)
				if rErr := svc.projectsPost(fmt.Sprintf("projects/%d/cancel", projResp.ProjectID), getJWT(c)); rErr != nil {
					log.Printf("ERROR: unable to cancel project %d: %s", projResp.ProjectID, rErr.Message)
				} else {
					projCnt++
				}
			}

			if err := svc.DB.Delete(&unit{}, u.ID).Error; err != nil {
				log.Printf("ERROR: unable to delete unit %d: %s", u.ID, err.Error())
				failed = append(failed, u.ID)
			} else {
				delCnt++
			}
			if delCnt >= int(maxDel) {
				log.Printf("INFO: max deletions (%d) reached; stopping", maxDel)
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		log.Printf("INFO: unit cleanup done. %d units deleted. %d projects deleted", delCnt, projCnt)
		if len(failed) > 0 {
			log.Printf("INFO: %d failed: %v", len(failed), failed)
		}
		if len(hasFiles) > 0 {
			log.Printf("INFO: %d units with masterfiles: %v", len(hasFiles), hasFiles)
		}
	}()

	c.String(http.StatusOK, "process started to retire %d of %d canceled units older than %s", maxDel, len(unitResp), dateStr)
}

func (svc *serviceContext) cleanupCanceledOrders(c *gin.Context) {
	log.Printf("INFO: cleanup canceled orders")
	maxDel, _ := strconv.ParseUint(c.Query("max"), 10, 64)
	if maxDel == 0 {
		maxDel = 500
	}
	log.Printf("INFO: limit to %d deletions", maxDel)

	deleteThreshold := time.Now().AddDate(0, -3, 0)
	dateStr := deleteThreshold.Format("2006-01-02")

	var orderResp []struct {
		ID           int64
		InvoiceID    int64
		DateCanceled *time.Time
		UpdatedAt    time.Time
		Units        int64
	}

	qStr := "select o.id,o.date_canceled, o.updated_at, (select id from invoices where order_id=o.id) as invoice_id, count(u.id) as units from orders o"
	qStr += " left join units u on u.order_id = o.id"
	qStr += " where order_status=? group by o.id"
	if err := svc.DB.Debug().Raw(qStr, "canceled").Scan(&orderResp).Error; err != nil {
		log.Printf("ERROR: find canceled units failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(orderResp) == 0 {
		log.Printf("INFO: there are no canceled orders to cleanup")
		c.String(http.StatusOK, "done")
		return
	}

	log.Printf("INFO: found %d canceled orders", len(orderResp))
	go func() {
		delCnt := 0
		failed := make([]int64, 0)
		hasUnits := make([]int64, 0)
		for _, o := range orderResp {
			log.Printf("INFO: check if order %d can be deleted", o.ID)
			if o.Units > 0 {
				log.Printf("INFO: order %d has %d units and cannot be deleted", o.ID, o.Units)
				hasUnits = append(hasUnits, o.ID)
				continue
			}

			cancelDate := o.UpdatedAt.Format("2006-01-02")
			if o.DateCanceled != nil {
				cancelDate = o.DateCanceled.Format("2006-01-02")
			}
			if cancelDate < dateStr {
				log.Printf("INFO: order %d canceled at %s can be deleted", o.ID, cancelDate)
				if o.InvoiceID > 0 {
					log.Printf("INFO: canceled order %d has invoice %d; delete it", o.ID, o.InvoiceID)
					if err := svc.DB.Exec("delete from invoices where id = ?", o.InvoiceID).Error; err != nil {
						log.Printf("ERROR: unable remove order %d invoice %d: %s", o.ID, o.InvoiceID, err.Error())
						continue
					}
				}
				if err := svc.DB.Delete(&order{}, o.ID).Error; err != nil {
					log.Printf("ERROR: unable to delete order %d: %s", o.ID, err.Error())
					failed = append(failed, o.ID)
				} else {
					delCnt++
				}
			}

			if delCnt >= int(maxDel) {
				log.Printf("INFO: max deletions (%d) reached; stopping", maxDel)
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		log.Printf("INFO: oder cleanup done. %d orders deleted", delCnt)
		if len(failed) > 0 {
			log.Printf("INFO: %d failed: %v", len(failed), failed)
		}
		if len(hasUnits) > 0 {
			log.Printf("INFO: %d orders with units: %v", len(hasUnits), hasUnits)
		}
	}()

	c.String(http.StatusOK, "process started to retire %d of %d canceled orders older than %s", maxDel, len(orderResp), dateStr)
}

func (svc *serviceContext) cleanupExpiredJobLogs(c *gin.Context) {
	log.Printf("INFO: cleanup job logs older than 2 months")
	deleteThreshold := time.Now().AddDate(0, -2, 0)
	dateStr := deleteThreshold.Format("2006-01-02")

	log.Printf("INFO: scan for job statuses to delete")
	var delCount int64
	if err := svc.DB.Table("job_statuses").Where("status=? and ended_at < ?", "finished", dateStr).Count(&delCount).Error; err != nil {
		log.Printf("ERROR: unable to get count of old jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if delCount == 0 {
		log.Printf("INFO: there are no jobs to delete")
		c.String(http.StatusOK, "no messages to delete")
		return
	}

	log.Printf("INFO: delete %d old jobs ", delCount)
	if err := svc.DB.Exec("DELETE from job_statuses where status=? and ended_at < ?", "finished", dateStr).Error; err != nil {
		log.Printf("ERROR: unable to delete finished jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d jobs deleted", delCount))
}
