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
	Apollo   string
	Jobs     string
	PDF      string
	Virgo    string
	TSAPI    string
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

type agency struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type category struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type containerType struct {
	ID         int64
	Name       string
	hasFolders bool
}

type workflow struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"-"`
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
			Apollo:   cfg.apolloURL,
			PDF:      cfg.pdfURL,
			Virgo:    cfg.virgoURL,
			TSAPI:    cfg.tracksysAPIURL,
			Jobs:     cfg.jobsURL,
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
	type searchField struct {
		Field string `json:"value"`
		Label string `json:"label"`
	}

	type cfgData struct {
		Version                string `json:"version"`
		ReportsURL             string `json:"reportsURL"`
		ProjectsURL            string `json:"projectsURL"`
		IIIFURL                string `json:"iiifURL"`
		IIIFManifestURL        string `json:"iiifManifestURL"`
		CurioURL               string `json:"curioURL"`
		PdfURL                 string `json:"pdfURL"`
		JobsURL                string `json:"jobsURL"`
		ControlledVocabularies struct {
			AcademicStatuses []academicStatus `json:"academicStatuses"`
			Agencies         []agency         `json:"agencies"`
			Categories       []category       `json:"categories"`
			ContainerTypes   []containerType  `json:"containerTypes"`
			IntendedUses     []intendedUse    `json:"intendedUses"`
			Workflows        []workflow       `json:"workflows"`
		} `json:"controlledVocabularies"`
		SearchFields map[string][]searchField `json:"searchFields"`
	}
	resp := cfgData{Version: Version,
		CurioURL:        svc.ExternalSystems.Curio,
		IIIFURL:         svc.ExternalSystems.IIIF,
		IIIFManifestURL: svc.ExternalSystems.IIIFMan,
		ReportsURL:      svc.ExternalSystems.Reports,
		PdfURL:          svc.ExternalSystems.PDF,
		ProjectsURL:     svc.ExternalSystems.Projects,
		JobsURL:         svc.ExternalSystems.Jobs,
	}

	log.Printf("INFO: load academic statuses")
	err := svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.AcademicStatuses).Error
	if err != nil {
		log.Printf("ERROR: unable to get academic statuses: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load agencies")
	err = svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.Agencies).Error
	if err != nil {
		log.Printf("ERROR: unable to get agencies: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load categories")
	err = svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.Categories).Error
	if err != nil {
		log.Printf("ERROR: unable to get categories: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load container types")
	err = svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.ContainerTypes).Error
	if err != nil {
		log.Printf("ERROR: unable to get container types: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load intended uses")
	err = svc.DB.Order("description asc").Where("is_approved=?", 1).Find(&resp.ControlledVocabularies.IntendedUses).Error
	if err != nil {
		log.Printf("ERROR: unable to get intended uses: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load workflows")
	err = svc.DB.Order("name asc").Where("active=?", 1).Find(&resp.ControlledVocabularies.Workflows).Error
	if err != nil {
		log.Printf("ERROR: unable to get workflows: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: define global search fields")
	resp.SearchFields = map[string][]searchField{
		"all": {
			searchField{Field: "all", Label: "All fields"},
			searchField{Field: "id", Label: "ID"},
			searchField{Field: "pid", Label: "PID"},
			searchField{Field: "title", Label: "Title"},
		},
		"components": {
			searchField{Field: "all", Label: "All fields"},
			searchField{Field: "pid", Label: "PID"},
			searchField{Field: "title", Label: "Title"},
			searchField{Field: "label", Label: "Label"},
			searchField{Field: "content_desc", Label: "Description"},
			searchField{Field: "date", Label: "Date"},
			searchField{Field: "barcode", Label: "Barcode"},
			searchField{Field: "ead_id_att", Label: "Finding Aid"},
		},
		"masterfiles": {
			searchField{Field: "all", Label: "All fields"},
			searchField{Field: "pid", Label: "PID"},
			searchField{Field: "unit_id", Label: "Unit ID"},
			searchField{Field: "filename", Label: "Filename"},
			searchField{Field: "title", Label: "Title"},
			searchField{Field: "description", Label: "Description"},
			searchField{Field: "tag", Label: "Tag"},
		},
		"metadata": {
			searchField{Field: "all", Label: "All fields"},
			searchField{Field: "pid", Label: "PID"},
			searchField{Field: "title", Label: "Title"},
			searchField{Field: "barcode", Label: "Barcode"},
			searchField{Field: "call_number", Label: "Call number"},
			searchField{Field: "catalog_key", Label: "Catalog key"},
			searchField{Field: "creator_name", Label: "Creator name"},
		},
		"orders": {
			searchField{Field: "all", Label: "All fields"},
			searchField{Field: "id", Label: "Order ID"},
			searchField{Field: "order_title", Label: "Title"},
			searchField{Field: "last_name", Label: "Customer last name"},
			searchField{Field: "agency", Label: "Agency"},
			searchField{Field: "staff_notes", Label: "Staff notes"},
			searchField{Field: "special_instructions", Label: "Special instructions"},
		},
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
