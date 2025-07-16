package models

type Notification struct {
	From     string `json:"from"`
	Message  string `json:"message"`
	DialogID int32  `json:"dialog_id"`
}
