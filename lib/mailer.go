package rscni

import (
	"errors"
	"io/ioutil"
	"log"
	"net/mail"
    "os/exec"
	"path"
	"strings"

	"github.com/bernarpa/goutils"
)

type RsCniMailer struct {
}

func NewMailer() *RsCniMailer {
	return new(RsCniMailer)
}

func (mailer *RsCniMailer) Name() string {
	return "Rassegna Stampa CNI"
}

func (mailer *RsCniMailer) Mail(cfg *goutils.Cfg, logger *log.Logger) error {
	datadir, ok := cfg.Get("datadir")
	if !ok {
		return errors.New("Configuration value datadir is missing")
	}

	// TODO: controllo dei parametri di cfg smtp.*
	smtpHost, okSmtpHost := cfg.Get("smtp.host")
	smtpPort, okSmtpPort := cfg.Get("smtp.port")
	smtpUser, okSmtpUser := cfg.Get("smtp.username")
	smtpPass, okSmtpPass := cfg.Get("smtp.password")
	smtpName, okSmtpName := cfg.Get("smtp.fromname")
	smtpEmail, okSmtpEmail := cfg.Get("smtp.fromemail")
	if !okSmtpHost || !okSmtpPort || !okSmtpUser || !okSmtpPass || !okSmtpName || !okSmtpEmail {
		return errors.New("Configuration value smtp.* missing (one or more among smtp.host, smtp.port, smtp.username, smtp.password, smtp.name, smtp.email)")
	}

	smtpAccount := goutils.SmtpAccount{
		Host:     smtpHost,
		Port:     smtpPort,
		User:     smtpUser,
		Password: smtpPass,
	}

	from := mail.Address{Name: smtpName, Address: smtpEmail}

	mlFile := path.Join(datadir, "ml.txt")
	if !goutils.PathExists(mlFile) {
		return errors.New("Cannot find ML file: " + mlFile)
	}

	lastTimeFile := path.Join(datadir, "last.txt")
	last := ""
	if goutils.PathExists(lastTimeFile) {
		if bytes, err := ioutil.ReadFile(lastTimeFile); err != nil {
			return err
		} else {
			last = strings.TrimSpace(string(bytes))
		}
	}

	var lastDir, lastFile string
	if last == "" {
		logger.Printf("Sending %s for the first time\n", mailer.Name())
	} else {
		logger.Printf("Last %s issue: %s", mailer.Name(), last)
		v := strings.SplitN(last, "/", 2)
		if len(v) < 2 {
			return errors.New("Bad last issue format: " + last)
		}

		lastDir = v[0]
		lastFile = v[1]
	}

	bytes, err := ioutil.ReadFile(mlFile)
	if err != nil {
		return err
	}

	recipients := make([]*mail.Address, 0, 10)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		r, err := mail.ParseAddress(line)
		if err != nil {
			return err
		}

		recipients = append(recipients, r)
	}

	issuesDir := path.Join(datadir, "issues")
	subdirs, err := ioutil.ReadDir(issuesDir)
	if err != nil {
		return err
	}

	for _, sd := range subdirs {
		if sd.Name() < lastDir {
			continue
		}

		logger.Printf("Exploring %s\n", sd.Name())
		subdir := path.Join(issuesDir, sd.Name())
		issues, err := ioutil.ReadDir(subdir)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			if issue.Name() <= lastFile {
				continue
			}

            issuePath := path.Join(subdir, issue.Name())
			attachment := goutils.Attachment{
				ContentType: "application/pdf",
				Path:        issuePath,
			}
			for _, r := range recipients {
				logger.Printf("Sending %s to %v\n", issue.Name(), r)

				err = goutils.SendMail(
					smtpAccount,
					from,
					[]mail.Address{*r},
					strings.Split(issue.Name(), ".")[0], // Subject
					"Vedi l'allegato.", // Body
					[]goutils.Attachment{attachment},
				)
				if err != nil {
					return err
				}
			}

            if onedrivePut, ok := cfg.Get("onedriveput"); ok {
                dest := "Documenti\\\\Ordine degli Ingegneri\\\\Rassegna Stampa CNI\\\\" + sd.Name()
                if out, err := exec.Command("python2", onedrivePut, issuePath, dest).Output(); err != nil {
                    logger.Println(string(out))
                    return err
                } else {
                    logger.Println(string(out))
                }
            }

            if dropboxUploader, ok := cfg.Get("dropbox_uploader"); ok {
                dest := "Documenti/Ordine degli Ingegneri/Rassegna Stampa CNI/" + sd.Name()
                if out, err := exec.Command(dropboxUploader, "upload", issuePath, dest).Output(); err != nil {
                    logger.Println(string(out))
                    return err
                } else {
                    logger.Println(string(out))
                }
            }

			current := sd.Name() + "/" + issue.Name()
			if current > last {
				last = current
			}
		}
	}

	if err := ioutil.WriteFile(lastTimeFile, []byte(last), 0444); err != nil {
		return err
	}

	return nil
}
