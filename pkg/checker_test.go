package pkg

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/guillaumejacquart/go-http-scheduler/pkg/domain"
	"github.com/spf13/viper"
)

func TestRegisterChecks(t *testing.T) {
	apps, _ := getAllApps()
	for _, a := range apps {
		deleteApp(a.ID)
	}

	app := domain.App{
		Name:           "test1",
		URL:            "http://google.fr",
		Status:         "down",
		Method:         "GET",
		CronExpression: "* * * * *",
		JobStatus:      true,
	}

	insertApp(&app)

	go registerChecks()

	timer := time.NewTimer(time.Second * 3)

	<-timer.C

	assert.True(t, len(appsJob) >= 1)
}

func TestRegisterCheck(t *testing.T) {
	app := domain.App{
		Name:           "test2",
		URL:            "http://google.fr",
		Status:         "down",
		Method:         "GET",
		CronExpression: "* * * * *",
		JobStatus:      true,
	}

	insertApp(&app)

	registerCheck(app)

	timer := time.NewTimer(time.Second * 3)

	<-timer.C

	assert.True(t, len(appsJob) >= 1)
}

func TestCheckApp(t *testing.T) {
	app := domain.App{
		Name:           "test3",
		URL:            "http://google.fr",
		Status:         "down",
		Method:         "GET",
		CronExpression: "* * * * *",
		JobStatus:      true,
	}

	insertApp(&app)

	err := checkApp(app)

	assert.Equal(t, err, nil)

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestRunHTTPCheck(t *testing.T) {
	app := domain.App{}
	app.Name = "Test"
	app.URL = "http://google.fr"
	_, err := runHTTPCheck(app)

	if err != nil {
		t.Error("Google should be checked up")
	}
}

func TestUpdateCheckedApp(t *testing.T) {

	lastApp := domain.App{
		Name:           "test5",
		URL:            "http://google.fr",
		Method:         "GET",
		CronExpression: "* * * * *",
	}

	insertApp(&lastApp)

	app := domain.App{
		Name:           "test6",
		URL:            "http://google.fr",
		Method:         "GET",
		CronExpression: "* * * * *",
	}

	var err error

	updateCheckedApp(app, lastApp, &http.Response{StatusCode: 200}, err)

	newApp, newErr := getApp(lastApp.ID)

	if newErr != nil {
		t.Error("Gettin app failed")
	}

	if newApp.Status == "down" {
		t.Error("App status not updated")
	}
}

func TestAddHistory(t *testing.T) {
	viper.Set("history.enabled", true)
	lastApp := domain.App{
		Name:           "test7",
		URL:            "http://google.fr",
		Method:         "GET",
		CronExpression: "* * * * *",
	}

	insertApp(&lastApp)

	addHistory(lastApp, time.Now())

	history, newErr := getAppHistory(lastApp.ID)

	if newErr != nil {
		t.Error("Gettin app failed")
	}

	if len(history) == 0 {
		t.Error("App history not updated")
	}
}
