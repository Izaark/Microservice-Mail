package models

import (
	"time"
)

type ObjUserInfo struct {
	Email       string      `json:"email" `
	ObjTemplate ObjTemplate `json:"template" `
}
type ObjTemplate struct {
	Body string    `json:"body" `
	Date time.Time `json:"time"`
}
