package rscni

import (
	"errors"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/bernarpa/goutils"
)

type RsCniScraper struct {
}

type rsGiornaliera struct {
	title, urlDettaglio, urlPdf string
	date                        time.Time
}

func NewScraper() *RsCniScraper {
	return new(RsCniScraper)
}

func (scra *RsCniScraper) Name() string {
	return "Rassegna Stampa CNI"
}

func (scra *RsCniScraper) Scrape(cfg *goutils.Cfg, logger *log.Logger) error {
	datadir, ok := cfg.Get("datadir")
	if !ok {
		return errors.New("Configuration value datadir is missing")
	}

	doc, err := goquery.NewDocument("http://www.centrostudicni.it/rassegna-stampa/rassegna-stampa-quotidiana")
	if err != nil {
		return err
	}

	giorni := make([]rsGiornaliera, 0, 10)
	doc.Find(".list-title").Each(func(_ int, s *goquery.Selection) {
		a := s.Find("a")
		title := strings.TrimSpace(a.Text())
		urlDettaglio, found := a.Attr("href")
		if found && urlDettaglio != "#" {
			urlDettaglio = "http://www.centrostudicni.it" + urlDettaglio
			giorni = append(giorni, rsGiornaliera{title: title, urlDettaglio: urlDettaglio})
		}
	})

	dateRegexp := regexp.MustCompile(`[0-9]{4}_[0-9]{2}_[0-9]{2}`)

	for i := 0; i < len(giorni); i++ {
		doc, err := goquery.NewDocument(giorni[i].urlDettaglio)
		if err != nil {
			return err
		}

		doc.Find("a").Each(func(_ int, a *goquery.Selection) {
			url, found := a.Attr("href")
			dateStr := dateRegexp.FindString(url)
			if found && strings.HasSuffix(strings.ToLower(url), ".pdf") && dateStr != "" {
				giorni[i].urlPdf = "http://www.centrostudicni.it" + url
				date, err := time.Parse("2006_01_02", dateStr)
				if err != nil {
					logger.Printf("Error while parsing date in URL '%s': %v\n", giorni[i].urlPdf, err)
					giorni[i].urlPdf = ""
				} else {
					giorni[i].date = date
				}
			}
		})
	}

	baseDir := path.Join(datadir, "issues")
	os.MkdirAll(baseDir, os.ModeDir)

	for _, giorno := range giorni {
		if giorno.urlPdf == "" {
			continue
		}

		dir := path.Join(baseDir, giorno.date.Format("2006-01"))
		pdf := path.Join(dir, "Rassegna Stampa CNI "+giorno.date.Format("2006-01-02")+".pdf")

		if goutils.PathExists(pdf) {
			continue
		}

		logger.Printf("Downloading %s...\n", giorno.title)
		os.MkdirAll(dir, os.ModeDir)
		goutils.DownloadHttpFile(giorno.urlPdf, pdf)
	}

	return nil
}
