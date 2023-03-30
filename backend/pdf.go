package main

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

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
	includeText, _ := strconv.ParseBool(c.Query("text"))
	pageStr := c.Query("pages")
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

	// Write the PDF to the temp director
	destDir := fmt.Sprintf("/tmp/%s", token)
	os.MkdirAll(destDir, 0777)
	pdfFileName := filepath.Join(destDir, fmt.Sprintf("%s.pdf", token))
	log.Printf("INFO: write PDF to %s", pdfFileName)
	writeErr := os.WriteFile(pdfFileName, resp, 0644)
	if writeErr != nil {
		log.Printf("ERROR: unable to write PDF: %s", writeErr.Error())
		c.String(http.StatusInternalServerError, writeErr.Error())
		return
	}

	if includeText == false {
		c.Header("Content-Type", "application/pdf")
		c.File(pdfFileName)
		os.RemoveAll(destDir)
		return
	}

	log.Printf("INFO: include text for masterfiles [%s]", pageStr)
	textFileName := filepath.Join(destDir, fmt.Sprintf("%s.txt", token))
	var textList []string

	if pageStr == "all" {
		dbErr := svc.DB.Debug().Table("master_files").Where("unit_id = ?", unitID).Select("transcription_text").Find(&textList).Error
		if dbErr != nil {
			log.Printf("ERROR: unable to get all master file text for unit %d: %s", unitID, dbErr.Error())
		}
	} else {
		dbErr := svc.DB.Debug().Table("master_files").Where("id in ?", strings.Split(pageStr, ",")).Select("transcription_text").Find(&textList).Error
		if dbErr != nil {
			log.Printf("ERROR: unable to get master file text for unit %d, pages %s: %s", unitID, pageStr, dbErr.Error())
		}
	}

	log.Printf("INFO: write transcription text to %s", textFileName)
	txtFile, osErr := os.Create(textFileName)
	if osErr != nil {
		log.Printf("ERROR: unable to create file %s for transcription text: %s", textFileName, osErr.Error())
	} else {
		for _, txt := range textList {
			txtFile.WriteString(txt)
			txtFile.WriteString("\n\n")
		}
		txtFile.Close()
	}

	zipFileName := filepath.Join("/tmp/", fmt.Sprintf("%s.zip", token))
	log.Printf("INFO: create zip of PDF/Text result %s", zipFileName)
	zipFile, osErr := os.Create(zipFileName)
	if osErr != nil {
		log.Printf("ERROR: unable to create %s: %s", zipFileName, osErr.Error())
		c.Header("Content-Type", "application/pdf")
		c.File(pdfFileName)
	} else {
		zipWriter := zip.NewWriter(zipFile)
		addFileToZip(zipWriter, destDir, fmt.Sprintf("%s.pdf", token))
		if textFileName != "" {
			addFileToZip(zipWriter, destDir, fmt.Sprintf("%s.txt", token))
		}
		zipWriter.Close()
		zipFile.Close()
		c.Header("Content-Type", "application/zip")
		c.File(zipFileName)
	}

	log.Printf("INFO: cleaning up temp files")
	os.Remove(zipFileName)
	os.RemoveAll(destDir)
}

func addFileToZip(zw *zip.Writer, filePath string, fileName string) {
	fileToZip, err := os.Open(path.Join(filePath, fileName))
	if err != nil {
		log.Printf("ERROR: unable to open %s for inclusion in zip: %s", filePath, err.Error())
		return
	}
	defer fileToZip.Close()
	zipFileWriter, err := zw.Create(fileName)
	if _, err := io.Copy(zipFileWriter, fileToZip); err != nil {
		log.Printf("ERROR: unable to copy contenst of  %s to zip: %s", filePath, err.Error())
		return
	}
}
