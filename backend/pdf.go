package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (svc *serviceContext) requestPDF(c *gin.Context) {
	unitID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pagesStr := c.Query("pages")
	log.Printf("INFO: request pdf for unit %d pages [%s]", unitID, pagesStr)

	var tgtUnit unit
	dbErr := svc.DB.Preload("Metadata").First(&tgtUnit, unitID).Error
	if dbErr != nil {
		log.Printf("ERROR: unable to get target unit %d: %s", unitID, dbErr.Error())
		c.String(http.StatusBadRequest, dbErr.Error())
		return
	}

	url := fmt.Sprintf("%s/%s?unit=%d&embed=1", svc.ExternalSystems.PDF, tgtUnit.Metadata.PID, unitID)
	token := fmt.Sprintf("%x", md5.Sum([]byte(c.Param("id"))))
	if pagesStr != "" {
		token = fmt.Sprintf("%x", md5.Sum([]byte(pagesStr)))
		url += fmt.Sprintf("&pages=%s", pagesStr)
	}
	url += fmt.Sprintf("&token=%s", token)

	_, err := svc.getRequest(url)
	if err != nil {
		log.Printf("ERROR: pdf request for %s pages [%s] failed: %d:%s", tgtUnit.Metadata.PID, pagesStr, err.StatusCode, err.Message)
		c.String(err.StatusCode, err.Message)
		return
	}

	log.Printf("INFO: pdf for %s requested successfully; now check status", tgtUnit.Metadata.PID)
	resp, statusErr := svc.checkPDFStatus(tgtUnit, token)
	if statusErr != nil {
		log.Printf("ERROR: %s", statusErr.Error())
		c.String(http.StatusInternalServerError, statusErr.Error())
		return
	}

	// responses: READY, FAILED, PROCESSING, percentage% (includes the percent symbol)
	out := struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}{
		Status: resp,
		Token:  token,
	}
	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) getPDFStatus(c *gin.Context) {
	unitID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	token := c.Query("token")
	var tgtUnit unit
	dbErr := svc.DB.Preload("Metadata").First(&tgtUnit, unitID).Error
	if dbErr != nil {
		log.Printf("ERROR: unable to get target unit %d: %s", unitID, dbErr.Error())
		c.String(http.StatusBadRequest, dbErr.Error())
		return
	}

	resp, statusErr := svc.checkPDFStatus(tgtUnit, token)
	if statusErr != nil {
		log.Printf("ERROR: %s", statusErr.Error())
		c.String(http.StatusInternalServerError, statusErr.Error())
		return
	}
	c.String(http.StatusOK, resp)
}

func (svc *serviceContext) checkPDFStatus(tgtUnit unit, token string) (string, error) {
	url := fmt.Sprintf("%s/%s/status", svc.ExternalSystems.PDF, tgtUnit.Metadata.PID)
	if token != "" {
		url += fmt.Sprintf("?token=%s", token)
	} else {
		url += fmt.Sprintf("?unit=%d", tgtUnit.ID)
	}
	resp, err := svc.getRequest(url)
	if err != nil {
		return "", fmt.Errorf("%d:%s", err.StatusCode, err.Message)
	}
	return string(resp), nil
}

func (svc *serviceContext) downloadPDF(c *gin.Context) {
	unitID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	token := c.Query("token")
	log.Printf("INFO: download pdf for unit %d token [%s]", unitID, token)

	var tgtUnit unit
	dbErr := svc.DB.Preload("Metadata").First(&tgtUnit, unitID).Error
	if dbErr != nil {
		log.Printf("ERROR: unable to get target unit %d for pdf download: %s", unitID, dbErr.Error())
		c.String(http.StatusBadRequest, dbErr.Error())
		return
	}

	url := fmt.Sprintf("%s/%s/download", svc.ExternalSystems.PDF, tgtUnit.Metadata.PID)
	if token != "" {
		url += fmt.Sprintf("?token=%s", token)
	} else {
		url += fmt.Sprintf("?unit=%d", unitID)
	}
	resp, err := svc.getRequest(url)
	if err != nil {
		log.Printf("ERROR: unable to download pdf for unit %d: %d:%s", unitID, err.StatusCode, err.Message)
		c.String(err.StatusCode, err.Message)
		return
	}
	// c.Data(http.StatusOK, "application/pdf", resp)
	fileName := fmt.Sprintf("/tmp/%s.pdf", token)
	writeErr := os.WriteFile(fileName, resp, 0644)
	if writeErr != nil {
		log.Printf("ERROR: unbale to write PDF: %s", writeErr.Error())
	}
	c.Header("Content-Type", "application/pdf")
	c.File(fileName)
	os.Remove(fileName)
}
