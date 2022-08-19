package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type externalSystems struct {
	IIIFMan  string
	IIIF     string
	Reports  string
	Projects string
	Curio    string
	Sirsi    string
	PDF      string
}

// serviceContext contains common data used by all handlers
type serviceContext struct {
	Version         string
	HTTPClient      *http.Client
	DB              *gorm.DB
	JWTKey          string
	ExternalSystems externalSystems
	DevAuthUser     string
}

// RequestError contains http status code and message for a failed HTTP request
type RequestError struct {
	StatusCode int
	Message    string
}

// InitializeService sets up the service context for all API handlers
func initializeService(version string, cfg *configData) *serviceContext {
	ctx := serviceContext{Version: version,
		ExternalSystems: externalSystems{
			IIIFMan:  cfg.iiifManifestURL,
			IIIF:     cfg.iiifURL,
			Reports:  cfg.reportsURL,
			Projects: cfg.projectsURL,
			Curio:    cfg.curioURL,
			Sirsi:    cfg.sirsiURL,
			PDF:      cfg.pdfURL,
		},
		JWTKey:      cfg.jwtKey,
		DevAuthUser: cfg.devAuthUser}

	log.Printf("INFO: connecting to DB...")
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.db.User, cfg.db.Pass, cfg.db.Host, cfg.db.Name)
	gdb, err := gorm.Open(mysql.Open(connectStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	ctx.DB = gdb
	log.Printf("INFO: DB Connection established")

	log.Printf("INFO: create HTTP client...")
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 600 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	ctx.HTTPClient = &http.Client{
		Transport: defaultTransport,
		Timeout:   5 * time.Second,
	}
	log.Printf("INFO: HTTP Client created")
	return &ctx
}

func (svc *serviceContext) getConfig(c *gin.Context) {
	log.Printf("INFO: get service configuration")
	type cfgData struct {
		Version                string `json:"version"`
		ReportsURL             string `json:"reportsURL"`
		ProjectsURL            string `json:"projectsURL"`
		IiifURL                string `json:"iiifURL"`
		CurioURL               string `json:"curioURL"`
		PdfURL                 string `json:"pdfURL"`
		ControlledVocabularies struct {
			AcademicStatuses []academicStatus `json:"academicStatuses"`
		} `json:"controlledVocabularies"`
	}
	resp := cfgData{Version: Version,
		CurioURL:    svc.ExternalSystems.Curio,
		IiifURL:     svc.ExternalSystems.IIIF,
		ReportsURL:  svc.ExternalSystems.Reports,
		PdfURL:      svc.ExternalSystems.PDF,
		ProjectsURL: svc.ExternalSystems.Projects,
	}

	log.Printf("INFO: load academic statuses")
	dbResp := svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.AcademicStatuses)
	if dbResp.Error != nil {
		log.Printf("ERROR: unable to get academic statuses: %s", dbResp.Error.Error())
		c.String(http.StatusInternalServerError, dbResp.Error.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetVersion reports the version of the serivce
func (svc *serviceContext) getVersion(c *gin.Context) {
	build := "unknown"
	// working directory is the bin directory, and build tag is in the root
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = svc.Version
	vMap["build"] = build
	c.JSON(http.StatusOK, vMap)
}

// HealthCheck reports the health of the serivce
func (svc *serviceContext) healthCheck(c *gin.Context) {
	type hcResp struct {
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
	}
	hcMap := make(map[string]hcResp)
	hcMap["tracksys2"] = hcResp{Healthy: true}

	hcMap["database"] = hcResp{Healthy: true}
	sqlDB, err := svc.DB.DB()
	if err != nil {
		hcMap["database"] = hcResp{Healthy: false, Message: err.Error()}
	} else {
		err := sqlDB.Ping()
		if err != nil {
			hcMap["database"] = hcResp{Healthy: false, Message: err.Error()}
		}
	}

	c.JSON(http.StatusOK, hcMap)
}

func (svc *serviceContext) getRequest(url string) ([]byte, *RequestError) {
	return svc.sendRequest("GET", url, nil)
}
func (svc *serviceContext) putRequest(url string) ([]byte, *RequestError) {
	return svc.sendRequest("PUT", url, nil)
}
func (svc *serviceContext) postFormRequest(url string, payload *url.Values) ([]byte, *RequestError) {
	return svc.sendRequest("POST", url, payload)
}

func (svc *serviceContext) sendRequest(verb string, url string, payload *url.Values) ([]byte, *RequestError) {
	log.Printf("INFO: %s request: %s", verb, url)
	startTime := time.Now()

	var req *http.Request
	if verb == "POST" && payload != nil {
		req, _ = http.NewRequest("POST", url, strings.NewReader(payload.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, _ = http.NewRequest(verb, url, nil)
	}

	rawResp, rawErr := svc.HTTPClient.Do(req)
	resp, err := handleAPIResponse(url, rawResp, rawErr)
	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)

	if err != nil {
		log.Printf("ERROR: Failed response from %s %s - %d:%s. Elapsed Time: %d (ms)",
			verb, url, err.StatusCode, err.Message, elapsedMS)
	} else {
		log.Printf("INFO: Successful response from %s %s. Elapsed Time: %d (ms)", verb, url, elapsedMS)
	}
	return resp, err
}

func handleAPIResponse(logURL string, resp *http.Response, err error) ([]byte, *RequestError) {
	if err != nil {
		status := http.StatusBadRequest
		errMsg := err.Error()
		if strings.Contains(err.Error(), "Timeout") {
			status = http.StatusRequestTimeout
			errMsg = fmt.Sprintf("%s timed out", logURL)
		} else if strings.Contains(err.Error(), "connection refused") {
			status = http.StatusServiceUnavailable
			errMsg = fmt.Sprintf("%s refused connection", logURL)
		}
		return nil, &RequestError{StatusCode: status, Message: errMsg}
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		status := resp.StatusCode
		errMsg := string(bodyBytes)
		return nil, &RequestError{StatusCode: status, Message: errMsg}
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return bodyBytes, nil
}
