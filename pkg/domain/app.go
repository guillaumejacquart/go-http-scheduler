package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

type checkType int

const (
	ResponseCheck = checkType(0)
	StatusCheck   = checkType(1)
)

// App is the app model
type App struct {
	gorm.Model
	Name           string     `json:"name" binding:"required"`
	URL            string     `json:"url"`
	Method         string     `json:"method"`
	Body           string     `json:"body"`
	Status         string     `json:"status"`
	CronExpression string     `json:"cronExpression"`
	LastUpDate     *time.Time `json:"lastUpDate"`
	Headers        []Header   `json:"headers"`
}
