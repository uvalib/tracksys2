package main

import (
	"flag"
	"log"
)

type dbConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

type configData struct {
	port            int
	db              dbConfig
	virgoURL        string
	reportsURL      string
	projectsURL     string
	iiifManifestURL string
	iiifURL         string
	ilsURL          string
	jobsURL         string
	apolloURL       string
	jstorURL        string
	curioURL        string
	pdfURL          string
	solrURL         string
	xmlIndexURL     string
	apTrustURL      string
	devAuthUser     string
	jwtKey          string
}

func getConfiguration() *configData {
	var config configData
	flag.IntVar(&config.port, "port", 8080, "Port to offer service on (default 8085)")
	flag.StringVar(&config.jwtKey, "jwtkey", "", "JWT signature key")
	flag.StringVar(&config.apTrustURL, "aptrust", "https://repo.aptrust.org", "APTrust URL")
	flag.StringVar(&config.reportsURL, "reports", "https://dpg-reporting.lib.virginia.edu", "DPG reports URL")
	flag.StringVar(&config.projectsURL, "projects", "https://dpg-imaging.lib.virginia.edu", "DPG projects URL")
	flag.StringVar(&config.virgoURL, "virgo", "https://search.lib.virginia.edu", "Virgo URL")
	flag.StringVar(&config.iiifManifestURL, "iiifman", "https://iiifman.lib.virginia.edu", "IIIF manifest URL")
	flag.StringVar(&config.iiifURL, "iiif", "https://iiif.lib.virginia.edu/iiif", "IIIF URL")
	flag.StringVar(&config.ilsURL, "ils", "https://ils-connector.lib.virginia.edu", "ILS Connector API URL")
	flag.StringVar(&config.curioURL, "curio", "https://curio.lib.virginia.edu", "Curio URL")
	flag.StringVar(&config.pdfURL, "pdf", "https://pdfservice.lib.virginia.edu/pdf", "PDF service URL")
	flag.StringVar(&config.solrURL, "solr", "http://virgo4-solr-production-replica-private.internal.lib.virginia.edu:8080/solr/test_core", "Solr URL")
	flag.StringVar(&config.jobsURL, "jobs", "http://dockerprod1.lib.virginia.edu:8710", "URL for job processing")
	flag.StringVar(&config.apolloURL, "apollo", "https://apollo.lib.virginia.edu", "URL for Apollo")
	flag.StringVar(&config.xmlIndexURL, "xmlhook", "https://virgo4-image-tracksys-reprocess-ws.internal.lib.virginia.edu/api/reindex", "XML index webhook")

	// DB connection params
	flag.StringVar(&config.db.Host, "dbhost", "", "Database host")
	flag.IntVar(&config.db.Port, "dbport", 3306, "Database port")
	flag.StringVar(&config.db.Name, "dbname", "", "Database name")
	flag.StringVar(&config.db.User, "dbuser", "", "Database user")
	flag.StringVar(&config.db.Pass, "dbpass", "", "Database password")

	// dev user
	flag.StringVar(&config.devAuthUser, "devuser", "", "Authorized computing id for dev")

	flag.Parse()

	if config.db.Host == "" {
		log.Fatal("Parameter dbhost is required")
	}
	if config.db.Name == "" {
		log.Fatal("Parameter dbname is required")
	}
	if config.db.User == "" {
		log.Fatal("Parameter dbuser is required")
	}
	if config.db.Pass == "" {
		log.Fatal("Parameter dbpass is required")
	}
	if config.jwtKey == "" {
		log.Fatal("Parameter jwtkey is required")
	}

	log.Printf("[CONFIG] port          = [%d]", config.port)
	log.Printf("[CONFIG] aptrust       = [%s]", config.apTrustURL)
	log.Printf("[CONFIG] solr          = [%s]", config.solrURL)
	log.Printf("[CONFIG] reports       = [%s]", config.reportsURL)
	log.Printf("[CONFIG] projects      = [%s]", config.projectsURL)
	log.Printf("[CONFIG] virgo         = [%s]", config.virgoURL)
	log.Printf("[CONFIG] jobs          = [%s]", config.jobsURL)
	log.Printf("[CONFIG] apollo        = [%s]", config.apolloURL)
	log.Printf("[CONFIG] iiifman       = [%s]", config.iiifManifestURL)
	log.Printf("[CONFIG] iiif          = [%s]", config.iiifURL)
	log.Printf("[CONFIG] ils           = [%s]", config.ilsURL)
	log.Printf("[CONFIG] curio         = [%s]", config.curioURL)
	log.Printf("[CONFIG] pdf           = [%s]", config.pdfURL)
	log.Printf("[CONFIG] xmlhook       = [%s]", config.xmlIndexURL)
	log.Printf("[CONFIG] dbuser        = [%s]", config.db.User)
	log.Printf("[CONFIG] dbhost        = [%s]", config.db.Host)
	log.Printf("[CONFIG] dbport        = [%d]", config.db.Port)
	log.Printf("[CONFIG] dbname        = [%s]", config.db.Name)
	log.Printf("[CONFIG] dbuser        = [%s]", config.db.User)
	if config.devAuthUser != "" {
		log.Printf("[CONFIG] devuser       = [%s]", config.devAuthUser)
	}

	return &config
}
