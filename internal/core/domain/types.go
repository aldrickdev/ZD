package domain

type Event struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}
type User struct {
	// TODO: I think using email might be better because this would be the most
	// reliable way to match users across the entire application
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Pod  string `json:"pod"`
}
type UserEvent struct {
	UserID    uint   `json:"user_id"`
	EventName string `json:"event_name"`
	Points    uint   `json:"points"`
}
