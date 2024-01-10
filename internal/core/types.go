package core

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	PodId uint   `json:"pod_id"`
	Pod   Pod    `json:"pod"`
}
type Pod struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}
type Event struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}
type UserEvent struct {
	UserId  uint `json:"user_id"`
	EventId uint `json:"event_id"`
}

type FullUserEvent struct {
	User  User
	Event Event
}
type UserEventIdData struct {
	UserId  uint `json:"user_id"`
	EventId uint `json:"event_id"`
}

type UserEventNameData struct {
	UserName  string `json:"user_name"`
	EventName string `json:"event_name"`
}
