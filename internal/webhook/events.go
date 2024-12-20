package webhook

type EventName string

// Action Events
const (
	SettingsAcidDetected EventName = "settings_acid_detected"
	SettingsEncDetected  EventName = "settings_enc_detected"

	DNSUpdate EventName = "dns_update"
)

// Data Events
const (
	ProcessBegin  EventName = "process_begin"
	ProcessFinish EventName = "process_finish"

	LoginStart   EventName = "login_start"
	LoginError   EventName = "login_error"
	LoginSuccess EventName = "login_success"
	LoginFailed  EventName = "login_failed"
)

type Event interface {
	isWebhookEvent()
}

type BaseEvent struct {
	ID        uint      `json:"id"`
	Timestamp int64     `json:"timestamp"` // unix
	Name      EventName `json:"name"`
}

func (BaseEvent) isWebhookEvent() {}

type ActionEventStatus string

const (
	ActionEventStatusSuccess ActionEventStatus = "success"
	ActionEventStatusFailure ActionEventStatus = "failure"
)

type ActionEvent struct {
	BaseEvent
	Status     ActionEventStatus `json:"status"`
	ActionName string            `json:"action_name"`

	// existence determined by status
	Value        string `json:"value,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type DataEvent struct {
	BaseEvent
	Property interface{} `json:"property"`
}
