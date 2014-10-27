package main

import (
	"log"
	"os"

	"github.com/bernarpa/goutils"
	lib "github.com/bernarpa/rscni/lib"
)

func main() {
	logger := log.New(os.Stderr, "RsCni: ", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println("Welcome to Rassegna Stampa CNI")

	logger.Println("Loading configuration...")
	cfg, err := goutils.NewCfg("Rassegna Stampa CNI", "RSCNICFG")
	if err != nil {
		logger.Printf("Cannot load configuration: %v\n", err)
	}

	logger.Printf("Config: %v\n", cfg)

	scraper := lib.NewScraper()
	if err := scraper.Scrape(cfg, logger); err != nil {
		logger.Printf("Error while scraping %s: %v\n", scraper.Name(), err)
	}

	mailer := lib.NewMailer()
	if err := mailer.Mail(cfg, logger); err != nil {
		logger.Printf("Error while mailing %s: %v\n", mailer.Name(), err)
	}
}
