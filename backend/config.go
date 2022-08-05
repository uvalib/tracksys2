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

type smtpConfig struct {
	Host     string
	Port     int
	User     string
	Pass     string
	Sender   string
	fakeSMTP bool
}

type configData struct {
	port        int
	db          dbConfig
	smtp        smtpConfig
	reportsURL  string
	projectsURL string
	devAuthUser string
	jwtKey      string
}

func getConfiguration() *configData {
	var config configData
	flag.IntVar(&config.port, "port", 8080, "Port to offer service on (default 8085)")
	flag.StringVar(&config.jwtKey, "jwtkey", "", "JWT signature key")
	flag.StringVar(&config.reportsURL, "reports", "https://dpg-reporting.lib.virginia.edu", "DPG reports URL")
	flag.StringVar(&config.projectsURL, "projects", "https://dpg-imaging.lib.virginia.edu/", "DPG projects URL")

	// DB connection params
	flag.StringVar(&config.db.Host, "dbhost", "", "Database host")
	flag.IntVar(&config.db.Port, "dbport", 3306, "Database port")
	flag.StringVar(&config.db.Name, "dbname", "", "Database name")
	flag.StringVar(&config.db.User, "dbuser", "", "Database user")
	flag.StringVar(&config.db.Pass, "dbpass", "", "Database password")

	// SMTP settings
	flag.StringVar(&config.smtp.Host, "smtphost", "", "SMTP Host")
	flag.IntVar(&config.smtp.Port, "smtpport", 0, "SMTP Port")
	flag.StringVar(&config.smtp.User, "smtpuser", "", "SMTP User")
	flag.StringVar(&config.smtp.Pass, "smtppass", "", "SMTP Password")
	flag.StringVar(&config.smtp.Sender, "smtpsender", "digitalservices@virginia.edu", "SMTP sender email")
	flag.BoolVar(&config.smtp.fakeSMTP, "stubsmtp", false, "Log email insted of sending (dev mode)")

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
	log.Printf("[CONFIG] reports       = [%s]", config.reportsURL)
	log.Printf("[CONFIG] projects      = [%s]", config.projectsURL)
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
