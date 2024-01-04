package domain

type Event struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Pod  string `json:"pod"`
}
type UserEvent struct {
	UserID  uint `json:"user_id"`
	EventID uint `json:"event_id"`
}

type UserEventData struct {
	UserID  uint `json:"user_id"`
	EventID uint `json:"event_id"`
}

type MarqueeData struct {
	UserName  string `json:"user_name"`
	EventName string `json:"event_name"`
}
