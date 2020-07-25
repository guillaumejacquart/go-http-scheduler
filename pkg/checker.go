package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/guillaumejacquart/go-http-scheduler/pkg/domain"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var c = cron.New()
var appsJob = make(map[uint]cron.EntryID)

func registerChecks() {
	for _, e := range c.Entries() {
		c.Remove(e.ID)
	}

	apps, err := getAllApps()

	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		registerCheck(a)
	}
	c.Run()
}

func registerCheck(a domain.App) {
	log.Println(a.ID)
	jobID, exists := appsJob[a.ID]
	log.Printf("Exist %v", exists)
	if exists {
		log.Println("Job for app already exists, cleaning it")
		c.Remove(jobID)
		delete(appsJob, a.ID)
	}

	eID, _ := c.AddFunc(a.CronExpression, func() {
		checkApp(a)
	})

	appsJob[a.ID] = eID
}

func checkApp(a domain.App) error {
	lastApp, _ := getApp(a.ID)

	resp, err := runHTTPCheck(lastApp)
	updateCheckedApp(a, lastApp, resp, err)

	return err
}

func runHTTPCheck(a domain.App) (*http.Response, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var body io.Reader
	if a.Body != "" {
		body = strings.NewReader(a.Body)
	}
	req, err := http.NewRequest(a.Method, a.URL, body)

	if len(a.Headers) > 0 {
		for _, h := range a.Headers {
			req.Header.Add(h.Name, h.Value)
		}
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func updateCheckedApp(a domain.App, lastApp domain.App, resp *http.Response, err error) {
	var status string
	nowDate := time.Now()

	if err != nil || resp.StatusCode >= 400 {
		status = "down"
	} else {
		status = "up"
	}

	fmt.Println("App", lastApp.URL, "is", status)

	updateAppStatus(lastApp.ID, status)

	if lastApp.Status != a.Status {
		addHistory(lastApp, nowDate)
	}
}

func addHistory(app domain.App, date time.Time) {
	if viper.GetBool("history.enabled") {
		history := domain.History{
			AppID:  app.ID,
			Date:   date,
			Status: app.Status,
		}
		insertHistory(history)
	}
}
