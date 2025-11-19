package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	manticore "github.com/manticoresoftware/manticoresearch-go"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type externalSystems struct {
	APTrust  string
	IIIFMan  string
	IIIF     string
	ILS      string
	Reports  string
	Projects string
	Curio    string
	Apollo   string
	Jobs     string
	PDF      string
	Virgo    string
	Solr     string
	XMLIndex string
}

// serviceContext contains common data used by all handlers
type serviceContext struct {
	Version         string
	HTTPClient      *http.Client
	DB              *gorm.DB
	Index           *manticore.SearchAPIService
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

type collectionFacet struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type containerType struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	HasFolders bool   `json:"hasFolders"`
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
			APTrust:  cfg.apTrustURL,
			IIIFMan:  cfg.iiifManifestURL,
			IIIF:     cfg.iiifURL,
			ILS:      cfg.ilsURL,
			Reports:  cfg.reportsURL,
			Projects: cfg.projectsURL,
			Curio:    cfg.curioURL,
			Apollo:   cfg.apolloURL,
			PDF:      cfg.pdfURL,
			Virgo:    cfg.virgoURL,
			Solr:     cfg.solrURL,
			Jobs:     cfg.jobsURL,
			XMLIndex: cfg.xmlIndexURL,
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

	log.Printf("INFO: configure db pool settings...")
	sqlDB, _ := gdb.DB()
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	ctx.DB = gdb
	log.Printf("INFO: DB Connection established")

	log.Printf("INFO: connect to search index...")
	mc := manticore.NewConfiguration()
	mc.Servers[0].URL = cfg.index
	apiClient := manticore.NewAPIClient(mc)
	ctx.Index = apiClient.SearchAPI
	log.Printf("INFO: search index connected")

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

func parseDateString(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, err
	}
	parsed, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, err
}

func (svc *serviceContext) cleanupExpiredData(c *gin.Context) {
	log.Printf("INFO: cleanup job logs and deleted messages older than 2 months")
	lastMonth := time.Now().AddDate(0, -2, 0)
	out := struct {
		DeletedJobs     int64 `json:"deletedJobs"`
		DeletedMessages int64 `json:"deletedMessages"`
	}{
		DeletedJobs:     0,
		DeletedMessages: 0,
	}

	log.Printf("INFO: scan for job statuses to delete")
	var oldStatuses []jobStatus
	err := svc.DB.Where("status=? and ended_at < ?", "finished", lastMonth).Find(&oldStatuses).Error
	if err != nil {
		log.Printf("ERROR: unable to get count of old jobs: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	out.DeletedJobs = int64(len(oldStatuses))
	if out.DeletedJobs > 0 {
		log.Printf("INFO: delete %d expired jobs", out.DeletedJobs)
		jsIDs := make([]uint64, 0)
		for _, js := range oldStatuses {
			jsIDs = append(jsIDs, js.ID)
		}
		err = svc.DB.Where("id in ?", jsIDs).Delete(&jobStatus{}).Error
		if err != nil {
			log.Printf("ERROR: unable to delete old jobs: %s", err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	log.Printf("INFO: scan for messages to delete")
	err = svc.DB.Table("messages").Where("deleted=? and deleted_at < ?", 1, lastMonth).Count(&out.DeletedMessages).Error
	if err != nil {
		log.Printf("ERROR: unable to get count of old messages: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: delete %d deleted messages", out.DeletedMessages)
	err = svc.DB.Exec("DELETE from messages where deleted=? and deleted_at < ?", 1, lastMonth).Error
	if err != nil {
		log.Printf("ERROR: unable to delete messages: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}

func (svc *serviceContext) addAgency(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid create agency request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: create new agency [%s] - [%s]", req.Name, req.Desc)
	newAgency := agency{Name: req.Name, Description: req.Desc}
	err = svc.DB.Create(&newAgency).Error
	if err != nil {
		log.Printf("ERROR: unable to create agency %s: %s", req.Name, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var agencies []agency
	err = svc.DB.Order("name asc").Find(&agencies).Error
	if err != nil {
		log.Printf("ERROR: unable to get updated agencies: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, agencies)
}

func (svc *serviceContext) getConfig(c *gin.Context) {
	log.Printf("INFO: get service configuration")
	type searchField struct {
		Field string `json:"value"`
		Label string `json:"label"`
	}

	type cfgData struct {
		Version                string `json:"version"`
		APTtustURL             string `json:"apTrustURL"`
		ReportsURL             string `json:"reportsURL"`
		ProjectsURL            string `json:"projectsURL"`
		IIIFURL                string `json:"iiifURL"`
		IIIFManifestURL        string `json:"iiifManifestURL"`
		CurioURL               string `json:"curioURL"`
		JobsURL                string `json:"jobsURL"`
		ControlledVocabularies struct {
			AcademicStatuses     []academicStatus     `json:"academicStatuses"`
			Agencies             []agency             `json:"agencies"`
			AvailabilityPolicies []availabilityPolicy `json:"availabilityPolicies"`
			Categories           []category           `json:"categories"`
			CollectionFacets     []collectionFacet    `json:"collectionFacets"`
			ContainerTypes       []containerType      `json:"containerTypes"`
			ExternalSyatems      []externalSystem     `json:"externalSystems"`
			IntendedUses         []intendedUse        `json:"intendedUses"`
			OCRHints             []ocrHint            `json:"ocrHints"`
			OCRLanguageHints     []ocrLanguageHint    `json:"ocrLanguageHints"`
			PreservationTiers    []preservationTier   `json:"preservationTiers"`
			UseRights            []useRight           `json:"useRights"`
			Workflows            []workflow           `json:"workflows"`
		} `json:"controlledVocabularies"`
	}

	vMap := svc.lookupVersion()
	resp := cfgData{Version: fmt.Sprintf("%s-%s", vMap["version"], vMap["build"]),
		APTtustURL:      svc.ExternalSystems.APTrust,
		CurioURL:        svc.ExternalSystems.Curio,
		IIIFURL:         svc.ExternalSystems.IIIF,
		IIIFManifestURL: svc.ExternalSystems.IIIFMan,
		ReportsURL:      svc.ExternalSystems.Reports,
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

	log.Printf("INFO: load availability policies")
	err = svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.AvailabilityPolicies).Error
	if err != nil {
		log.Printf("ERROR: unable to get availability policies: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load collection facets")
	err = svc.DB.Order("name asc").Find(&resp.ControlledVocabularies.CollectionFacets).Error
	if err != nil {
		log.Printf("ERROR: unable to get collection facets: %s", err.Error())
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

	log.Printf("INFO: load external systems")
	err = svc.DB.Order("id asc").Find(&resp.ControlledVocabularies.ExternalSyatems).Error
	if err != nil {
		log.Printf("ERROR: unable to get external systems: %s", err.Error())
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

	log.Printf("INFO: load ocr hints")
	err = svc.DB.Order("id asc").Find(&resp.ControlledVocabularies.OCRHints).Error
	if err != nil {
		log.Printf("ERROR: unable to get ocr hints: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load ocr language hints")
	resp.ControlledVocabularies.OCRLanguageHints = make([]ocrLanguageHint, 0)
	f, err := os.Open("./data/languages.csv")
	if err != nil {
		log.Printf("ERROR: unable to load ocr language hints: %s", err.Error())
	} else {
		defer f.Close()
		csvReader := csv.NewReader(f)
		langRecs, err := csvReader.ReadAll()
		if err != nil {
			log.Printf("ERROR: unable to parse languages file: %s", err.Error())
		} else {
			for _, rec := range langRecs {
				resp.ControlledVocabularies.OCRLanguageHints = append(resp.ControlledVocabularies.OCRLanguageHints, ocrLanguageHint{Code: rec[0], Language: rec[1]})
			}
		}
	}

	log.Printf("INFO: load preservation tiers")
	err = svc.DB.Order("id asc").Find(&resp.ControlledVocabularies.PreservationTiers).Error
	if err != nil {
		log.Printf("ERROR: unable to get preservtion tiers: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: load use rights")
	err = svc.DB.Order("id asc").Find(&resp.ControlledVocabularies.UseRights).Error
	if err != nil {
		log.Printf("ERROR: unable to get use rights: %s", err.Error())
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

	c.JSON(http.StatusOK, resp)
}

// GetVersion reports the version of the serivce
func (svc *serviceContext) getVersion(c *gin.Context) {
	vMap := svc.lookupVersion()
	c.JSON(http.StatusOK, vMap)
}

func (svc *serviceContext) lookupVersion() map[string]string {
	build := "unknown"
	// working directory is the bin directory, and build tag is in the root
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = svc.Version
	vMap["build"] = build
	return vMap
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

func (svc *serviceContext) postJSON(url string, jsonPayload any) ([]byte, *RequestError) {
	log.Printf("INFO: POST json request: %s", url)
	startTime := time.Now()

	payload, _ := json.Marshal(jsonPayload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("Content-type", "application/json")

	rawResp, rawErr := svc.HTTPClient.Do(req)
	resp, err := handleAPIResponse(url, rawResp, rawErr)
	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)

	if err != nil {
		log.Printf("ERROR: Failed response from POST %s - %d:%s. Elapsed Time: %d (ms)",
			url, err.StatusCode, err.Message, elapsedMS)
	} else {
		log.Printf("INFO: Successful POST response from %s. Elapsed Time: %d (ms)", url, elapsedMS)
	}
	return resp, err
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
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36")
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		status := resp.StatusCode
		errMsg := string(bodyBytes)
		return nil, &RequestError{StatusCode: status, Message: errMsg}
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	return bodyBytes, nil
}

func (svc *serviceContext) awaitJobCompletion(jobID int64) {
	statusURL := fmt.Sprintf("%s/jobs/%d", svc.ExternalSystems.Jobs, jobID)
	done := false
	for done == false {
		time.Sleep(15 * time.Second)
		rawResp, err := svc.getRequest(statusURL)
		if err != nil {
			log.Printf("ERROR: unable to get status for job %d: %s", jobID, err.Message)
		} else {
			var js jobStatus
			pErr := json.Unmarshal(rawResp, &js)
			if pErr != nil {
				log.Printf("ERROR: unable to parse job %d status response %s: %s", jobID, rawResp, pErr.Error())
				done = true
			} else {
				switch js.Status {
				case "finished":
					log.Printf("INFO: job %d has completed", jobID)
					done = true
				case "failure":
					log.Printf("INFO: job %d failed: %s", jobID, js.Error)
					done = true
				}
			}
		}
	}
}

func saveUploadedFile(formFile *multipart.FileHeader, destFile string) error {
	frmFile, err := formFile.Open()
	if err != nil {
		return fmt.Errorf("unable to open uploaded image %s: %s", formFile.Filename, err.Error())
	}
	defer frmFile.Close()
	out, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("unable to create temp file %s for uploaded image %s: %s", destFile, formFile.Filename, err.Error())
	}
	defer out.Close()
	_, err = io.Copy(out, frmFile)
	if err != nil {
		return err
	}
	return nil
}
