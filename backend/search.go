package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/gin-gonic/gin"
	manticore "github.com/manticoresoftware/manticoresearch-go"
	"golang.org/x/image/tiff"
)

type masterFileHit struct {
	ID           uint64 `json:"id"`
	PID          string `json:"pid"`
	UnitID       uint64 `json:"unit_id"`
	Filename     string `json:"filename"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url"`
	ImageURL     string `json:"image_url"`
	CallNumber   string `json:"call_number"`
	MetadataID   int64  `json:"metadata_id"`
	IsClone      bool   `json:"is_clone"`
}

type metadataHit struct {
	ID          uint64 `json:"id"`
	PID         string `json:"pid"`
	SystemName  string `json:"system_name"`
	Title       string `json:"title"`
	CallNumber  string `json:"call_number"`
	Barcode     string `json:"barcode"`
	CatalogKey  string `json:"catalog_key"`
	CreatorName string `json:"creator_name"`
	Virgo       bool   `json:"virgo"`
	DPLA        bool   `json:"dpla"`
	HathiTrust  bool   `json:"hathitrust"`
	VirgoURL    string `json:"virgo_url"`
}

type orderHit struct {
	ID                  uint64 `json:"id"`
	Status              string `json:"status"`
	Title               string `json:"title"`
	StaffNotes          string `json:"staff_notes"`
	SpecialInstructions string `json:"special_instructions"`
	CustomerName        string `json:"customer"`
	Agency              string `json:"agency"`
}

type unitHit struct {
	ID                          uint64 `json:"id"`
	Status                      string `json:"status"`
	StaffNotes                  string `json:"staff_notes"`
	SpecialInstructions         string `json:"special_instructions"`
	DateDLDeliverablesReady     string `json:"date_dl_deliverables_ready,omitempty"`
	DatePatronDeliverablesReady string `json:"date_patron_deliverables_ready,omitempty"`
}

type componentHit struct {
	ID              uint64 `json:"id"`
	PID             string `json:"pid"`
	Title           string `json:"title"`
	Label           string `json:"label"`
	Description     string `json:"description"`
	Date            string `json:"date"`
	FindingAid      string `json:"finding_aid"`
	MasterFileCount uint   `json:"mf_cnt"`
}

type componentResp struct {
	Total int64          `json:"total"`
	Hits  []componentHit `json:"hits"`
}
type metadataResp struct {
	Total int64         `json:"total"`
	Hits  []metadataHit `json:"hits"`
}
type masterFileResp struct {
	Total int64           `json:"total"`
	Hits  []masterFileHit `json:"hits"`
}
type orderResp struct {
	Total int64      `json:"total"`
	Hits  []orderHit `json:"hits"`
}
type unitResp struct {
	Total int64     `json:"total"`
	Hits  []unitHit `json:"hits"`
}

type searchResults struct {
	Components  componentResp  `json:"components"`
	MasterFiles masterFileResp `json:"masterFiles"`
	Metadata    metadataResp   `json:"metadata"`
	Orders      orderResp      `json:"orders"`
	Units       unitResp       `json:"units"`
}

type filterRequest struct {
	Type   string   `json:"type"`
	Params []string `json:"params"`
}

type filterParam struct {
	Field string
	Value any
	Exact bool
}

type searchFilter struct {
	Target string
	Params []filterParam
}

type searchContext struct {
	Query      string
	Filter     *searchFilter
	StartIndex int
	PageSize   int
}

type searchChannel struct {
	Type    string
	Results any
}

func (svc *serviceContext) searchRequest(c *gin.Context) {
	// setup the search context, which will be passed to each search function
	q := strings.TrimSpace(c.Query("q"))
	sc := searchContext{Query: q}

	tgtScope := c.Query("scope")
	if tgtScope != "all" && tgtScope != "masterfiles" && tgtScope != "metadata" && tgtScope != "orders" && tgtScope != "components" && tgtScope != "units" {
		log.Printf("ERROR: invalid search scope %s specified", tgtScope)
		c.String(http.StatusBadRequest, "invalid search scope")
		return
	}

	sc.StartIndex, _ = strconv.Atoi(c.Query("start"))
	sc.PageSize, _ = strconv.Atoi(c.Query("limit"))
	if sc.PageSize == 0 {
		sc.PageSize = 15
	}

	sc.Filter = svc.initFilter(c.Query("filters"))

	// cleanup format of PID values which may be pasted from other sources
	if strings.Index(q, "uva-lib-") == 0 {
		log.Printf("INFO: query %s appears to be a uva-lib pid formatted with a dash; update format and search on pid columns", q)
		sc.Query = strings.ReplaceAll(q, "uva-lib-", "uva-lib:")
	} else if strings.Index(q, "tsb-") == 0 {
		log.Printf("INFO: query %s appears to be a tsb pid formatted with a dash; update format and search on pid columns", q)
		sc.Query = strings.ReplaceAll(q, "tsb-", "tsb:")
	}

	log.Printf("INFO: search %s for [%s] starting from %d limit %d", tgtScope, sc.Query, sc.StartIndex, sc.PageSize)

	// query each type of object individually and await for all responses
	log.Printf("INFO: issue search requests...")
	pendingCount := 0
	channel := make(chan searchChannel)
	startTime := time.Now()
	if tgtScope == "all" || tgtScope == "masterfiles" {
		pendingCount++
		go svc.queryMasterFiles(&sc, channel)
	}
	if tgtScope == "all" || tgtScope == "components" {
		pendingCount++
		go svc.queryComponents(&sc, channel)
	}
	if tgtScope == "all" || tgtScope == "metadata" {
		pendingCount++
		go svc.queryMetadata(&sc, channel)
	}
	if tgtScope == "all" || tgtScope == "orders" {
		pendingCount++
		go svc.queryOrders(&sc, channel)
	}
	if tgtScope == "all" || tgtScope == "units" {
		pendingCount++
		go svc.queryUnits(&sc, channel)
	}

	log.Printf("INFO: await all search responses...")
	resp := searchResults{}
	for pendingCount > 0 {
		searchResp := <-channel
		pendingCount--
		if searchResp.Type == "masterFiles" {
			log.Printf("INFO received master files search response")
			mfResp, ok := searchResp.Results.(masterFileResp)
			if ok {
				resp.MasterFiles = mfResp
			}
		} else if searchResp.Type == "metadata" {
			log.Printf("INFO received metadata search response")
			mResp, ok := searchResp.Results.(metadataResp)
			if ok {
				resp.Metadata = mResp
			}
		} else if searchResp.Type == "components" {
			log.Printf("INFO received components search response")
			cResp, ok := searchResp.Results.(componentResp)
			if ok {
				resp.Components = cResp
			}
		} else if searchResp.Type == "orders" {
			log.Printf("INFO received orders search response")
			oResp, ok := searchResp.Results.(orderResp)
			if ok {
				resp.Orders = oResp
			}
		} else if searchResp.Type == "units" {
			log.Printf("INFO: received units search response")
			uResp, ok := searchResp.Results.(unitResp)
			if ok {
				resp.Units = uResp
			}
		}
	}
	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)
	log.Printf("INFO: all responses received. Elapsed Time: %d (ms)", elapsedMS)

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) initFilter(filterStr string) *searchFilter {
	log.Printf("INFO: raw filters [%s]", filterStr)
	out := searchFilter{Target: "none", Params: make([]filterParam, 0)}

	if filterStr != "" {
		//  Format: filters={"type":"TABLE_NAME","params":["FIELD_NAME|contains|sing","FIELD2|equals|true"]}'
		var reqFilter filterRequest
		err := json.Unmarshal([]byte(filterStr), &reqFilter)
		if err != nil {
			log.Printf("ERROR: invalid format for filter %s: %s", filterStr, err.Error())
			return &out
		}
		out.Target = reqFilter.Type
		for _, f := range reqFilter.Params {
			bits := strings.Split(f, "|")
			if bits[2] == "true" {
				out.Params = append(out.Params, filterParam{Field: bits[0], Value: 1, Exact: true})
			} else if bits[2] == "false" {
				out.Params = append(out.Params, filterParam{Field: bits[0], Value: 0, Exact: true})
			} else {
				out.Params = append(out.Params, filterParam{Field: bits[0], Value: bits[2], Exact: bits[1] == "equals"})
			}
		}
	}

	return &out
}

func (svc *serviceContext) queryMasterFiles(sc *searchContext, channel chan searchChannel) {
	resp := masterFileResp{Hits: make([]masterFileHit, 0)}

	newQ := newQuery("masterfiles", sc, int32(sc.StartIndex), int32(sc.PageSize))
	mResp, _, err := svc.Index.Search(context.Background()).SearchRequest(*newQ).Execute()
	if err != nil {
		log.Printf("ERROR: masterfiles search failed: %s", err.Error())
		channel <- searchChannel{Type: "masterFiles", Results: resp}
		return
	}
	resp.Total = int64(mResp.Hits.GetTotal())

	for _, h := range mResp.Hits.GetHits() {
		b, _ := json.Marshal(h.GetSource())

		var hitObj masterFileHit
		uErr := json.Unmarshal(b, &hitObj)
		if uErr != nil {
			log.Printf("ERROR: unable to unmarshal response; %s", uErr)
		} else {
			hitObj.ID = uint64(h.GetId())
			hitObj.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, hitObj.PID)
			hitObj.ImageURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, hitObj.PID)
			resp.Hits = append(resp.Hits, hitObj)
		}
	}

	channel <- searchChannel{Type: "masterFiles", Results: resp}
}

func (svc *serviceContext) queryUnits(sc *searchContext, channel chan searchChannel) {
	resp := unitResp{Hits: make([]unitHit, 0)}

	newQ := newQuery("units", sc, int32(sc.StartIndex), int32(sc.PageSize))
	mResp, _, err := svc.Index.Search(context.Background()).SearchRequest(*newQ).Execute()
	if err != nil {
		log.Printf("ERROR: units search failed: %s", err.Error())
		channel <- searchChannel{Type: "units", Results: resp}
		return
	}
	resp.Total = int64(mResp.Hits.GetTotal())

	for _, h := range mResp.Hits.GetHits() {
		b, _ := json.Marshal(h.GetSource())

		var hitObj unitHit
		uErr := json.Unmarshal(b, &hitObj)
		if uErr != nil {
			log.Printf("ERROR: unable to unmarshal unit response; %s", uErr)
		} else {
			hitObj.ID = uint64(h.GetId())
			resp.Hits = append(resp.Hits, hitObj)
		}
	}

	channel <- searchChannel{Type: "units", Results: resp}
}

func newQuery(table string, sc *searchContext, offset, limit int32) *manticore.SearchRequest {
	searchRequest := manticore.NewSearchRequest(table)
	searchRequest.SetLimit(limit)
	searchRequest.SetOffset(offset)

	searchQuery := manticore.NewSearchQuery()

	// Docs on search filter setup: https://manual.manticoresearch.com/Searching/Filters
	// search will be comprised of a list of queryFilters. The first filter is a
	// query_string with the text entered in the main query field on the UI. This
	// is a full-text search
	mustFiters := make([]manticore.QueryFilter, 0)
	qf := manticore.NewQueryFilter()
	qf.SetQueryString(sc.Query)
	mustFiters = append(mustFiters, *qf)

	if sc.Filter.Target == table {
		for _, fp := range sc.Filter.Params {
			if fp.Exact {
				// exact matches are for filtering on attributes rather than full text fields
				// they use a different setup; a query filter with an Equals match rather than a query string
				qf := manticore.NewQueryFilter()
				filter := map[string]any{fp.Field: fp.Value}
				qf.SetEquals(filter)
				mustFiters = append(mustFiters, *qf)
			} else {
				// non-exact matches are added as column-specific full text index searches (queryString)
				qf := manticore.NewQueryFilter()
				qf.SetQueryString(fmt.Sprintf("@%s %s", fp.Field, fp.Value))
				mustFiters = append(mustFiters, *qf)
			}
		}
	}

	// add all the QueryFilters collected above to a boolean must query
	boolFilter := manticore.NewBoolFilter()
	boolFilter.SetMust(mustFiters)
	searchQuery.SetBool(*boolFilter)

	searchRequest.Query = searchQuery
	b, _ := searchRequest.MarshalJSON()
	log.Printf("INFO: %s query details %s", table, b)
	return searchRequest
}

func (svc *serviceContext) queryMetadata(sc *searchContext, channel chan searchChannel) {
	resp := metadataResp{Hits: make([]metadataHit, 0)}
	newQ := newQuery("metadata", sc, int32(sc.StartIndex), int32(sc.PageSize))
	mResp, _, err := svc.Index.Search(context.Background()).SearchRequest(*newQ).Execute()
	if err != nil {
		log.Printf("ERROR: metadata search failed: %s", err.Error())
		channel <- searchChannel{Type: "metadata", Results: resp}
		return
	}
	resp.Total = int64(mResp.Hits.GetTotal())

	for _, h := range mResp.Hits.GetHits() {
		b, _ := json.Marshal(h.GetSource())
		var hitObj metadataHit
		uErr := json.Unmarshal(b, &hitObj)
		if uErr != nil {
			log.Printf("ERROR: unable to unmarshall metadata hit; %s", uErr)
		} else {
			hitObj.ID = uint64(h.GetId())
			if hitObj.Virgo {
				if hitObj.SystemName == "SirsiMetadata" {
					hitObj.VirgoURL = fmt.Sprintf("%s/sources/uva_library/items/%s", svc.ExternalSystems.Virgo, hitObj.CatalogKey)
				} else if hitObj.SystemName == "XmlMetadata" {
					hitObj.VirgoURL = fmt.Sprintf("%s/sources/images/items/%s", svc.ExternalSystems.Virgo, hitObj.PID)
				}
			}
			resp.Hits = append(resp.Hits, hitObj)
		}
	}
	channel <- searchChannel{Type: "metadata", Results: resp}
}

func (svc *serviceContext) queryOrders(sc *searchContext, channel chan searchChannel) {
	resp := orderResp{Hits: make([]orderHit, 0)}
	newQ := newQuery("orders", sc, int32(sc.StartIndex), int32(sc.PageSize))
	mResp, _, err := svc.Index.Search(context.Background()).SearchRequest(*newQ).Execute()
	if err != nil {
		log.Printf("ERROR: orders search failed: %s", err.Error())
		channel <- searchChannel{Type: "orders", Results: resp}
		return
	}
	resp.Total = int64(mResp.Hits.GetTotal())

	for _, h := range mResp.Hits.GetHits() {
		b, _ := json.Marshal(h.GetSource())
		var hitObj orderHit
		uErr := json.Unmarshal(b, &hitObj)
		if uErr != nil {
			log.Printf("ERROR: unable to unmarshall order hit; %s", uErr)
		} else {
			hitObj.ID = uint64(h.GetId())
			resp.Hits = append(resp.Hits, hitObj)
		}
	}

	channel <- searchChannel{Type: "orders", Results: resp}
}

func (svc *serviceContext) queryComponents(sc *searchContext, channel chan searchChannel) {
	resp := componentResp{Hits: make([]componentHit, 0)}
	newQ := newQuery("components", sc, int32(sc.StartIndex), int32(sc.PageSize))
	mResp, _, err := svc.Index.Search(context.Background()).SearchRequest(*newQ).Execute()
	if err != nil {
		log.Printf("ERROR: components search failed: %s", err.Error())
		channel <- searchChannel{Type: "components", Results: resp}
		return
	}
	resp.Total = int64(mResp.Hits.GetTotal())

	for _, h := range mResp.Hits.GetHits() {
		b, _ := json.Marshal(h.GetSource())
		var hitObj componentHit
		uErr := json.Unmarshal(b, &hitObj)
		if uErr != nil {
			log.Printf("ERROR: unable to unmarshall component hit; %s", uErr)
		} else {
			hitObj.ID = uint64(h.GetId())
			resp.Hits = append(resp.Hits, hitObj)
		}
	}

	channel <- searchChannel{Type: "components", Results: resp}
}

func (svc *serviceContext) imageSearchRequest(c *gin.Context) {
	pHashQ, parmErr := strconv.ParseUint(strings.TrimSpace(c.Query("phash")), 10, 64)
	if parmErr != nil {
		log.Printf("INFO: invalid phash param: %s", parmErr.Error())
		c.String(http.StatusBadRequest, parmErr.Error())
		return
	}
	distance, parmErr := strconv.ParseInt(strings.TrimSpace(c.Query("distance")), 10, 64)
	if parmErr != nil {
		log.Printf("INFO: invalid distance param: %s", parmErr.Error())
		c.String(http.StatusBadRequest, parmErr.Error())
		return
	}

	type similarImageHit struct {
		ID            int64  `gorm:"column:id" json:"id"`
		PID           string `gorm:"column:pid" json:"pid"`
		Filename      string `gorm:"column:filename" json:"filename"`
		Title         string `gorm:"column:title" json:"title"`
		Description   string `gorm:"column:description" json:"description"`
		Distance      int64  `gorm:"column:distance" json:"distance"`
		UnitID        int64  `gorm:"column:unit_id" json:"unitID"`
		MetadatID     int64  `gorm:"column:md_id" json:"metadataID"`
		MetadataPID   string `gorm:"column:md_pid" json:"metadataPID"`
		MetadataTitle string `gorm:"column:md_title" json:"metadataTitle"`
		ThumbnailURL  string `gorm:"-" json:"thumbnailURL"`
		ImageURL      string `gorm:"-" json:"imageURL"`
	}
	type similarResult struct {
		Hits  []*similarImageHit `json:"hits"`
		Total int64              `json:"total"`
	}

	log.Printf("INFO: searching for images matching pHash [%d] with distance [%d]", pHashQ, distance)
	resp := similarResult{Hits: make([]*similarImageHit, 0)}
	startTime := time.Now()
	distQ := "WHERE BIT_COUNT(phash ^ ?) <= ?"
	err := svc.DB.Raw(fmt.Sprintf("SELECT count(id) FROM master_files %s", distQ), pHashQ, distance).Scan(&resp.Total).Error
	if err != nil {
		log.Printf("ERROR: unable to get image search hit count: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	imgFields := "m.id as id, m.pid as pid, m.filename as filename, m.title as title, m.description as description, BIT_COUNT(phash ^ ?) as distance"
	mdFields := "m.unit_id as unit_id, m2.id as md_id, m2.title as md_title, m2.pid as md_pid"
	orderClause := "order by distance asc limit 0,50"
	err = svc.DB.Raw(fmt.Sprintf("SELECT %s, %s FROM master_files m left join metadata m2 on m2.id=metadata_id %s %s",
		imgFields, mdFields, distQ, orderClause), pHashQ, pHashQ, distance).Scan(&resp.Hits).Error
	if err != nil {
		log.Printf("ERROR: unable to get image search hits: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for _, mf := range resp.Hits {
		imgPID := mf.PID
		mf.ThumbnailURL = fmt.Sprintf("%s/%s/full/!125,200/0/default.jpg", svc.ExternalSystems.IIIF, imgPID)
		mf.ImageURL = fmt.Sprintf("%s/%s/full/full/0/default.jpg", svc.ExternalSystems.IIIF, imgPID)
	}

	elapsedNanoSec := time.Since(startTime)
	elapsedMS := int64(elapsedNanoSec / time.Millisecond)
	log.Printf("INFO: masterfile search found %d hits. Elapsed Time: %d (ms)", resp.Total, elapsedMS)
	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) uploadSearchImage(c *gin.Context) {
	log.Printf("INFO: received image search upload")
	formFile, err := c.FormFile("imageSearch")
	if err != nil {
		log.Printf("ERROR: unable to get upload image: %s", err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("unable to get file: %s", err.Error()))
		return
	}

	log.Printf("INFO: receive image %s", formFile.Filename)
	destFile := path.Join("/tmp", formFile.Filename)
	err = c.SaveUploadedFile(formFile, destFile)
	if err != nil {
		log.Printf("ERROR: unable to save %s: %s", formFile.Filename, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	imgFile, err := os.Open(destFile)
	if err != nil {
		log.Printf("ERROR: unable to open %s for phash generation: %s", destFile, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	defer imgFile.Close()
	fileType := strings.ToUpper(path.Ext(destFile))
	fileType = strings.Replace(fileType, ".", "", 1)
	var imgData image.Image

	if fileType == "TIF" {
		imgData, err = tiff.Decode(imgFile)
	} else if fileType == "JPG" {
		imgData, err = jpeg.Decode(imgFile)
	} else if fileType == "PNG" {
		imgData, err = png.Decode(imgFile)
	} else if fileType == "GIF" {
		imgData, err = gif.Decode(imgFile)
	}
	if err != nil {
		log.Printf("ERROR: unable to decode %s for phash generation: %s", destFile, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	imgHash, err := goimagehash.DifferenceHash(imgData)
	if err != nil {
		log.Printf("ERROR: unable to calculate pHash for %s: %s", destFile, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	pHash := imgHash.GetHash()
	os.Remove(destFile)

	c.String(http.StatusOK, fmt.Sprintf("%d", pHash))
}
