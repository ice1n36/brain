package models

type MobileAppNetTrafficRequest struct {
	AppID      string `json:"app_id,omitempty"`
	AppVersion string `json:"app_version,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	Host       string `json:"host,omitempty"`
	Port       string `json:"port,omitempty"`
}
