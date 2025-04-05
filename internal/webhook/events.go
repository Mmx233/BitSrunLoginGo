package webhook

import (
	"sync/atomic"
	"time"
)

type EventName string

// ActionEvent
const (
	SettingsAcidDetected EventName = "settings_acid_detected"
	SettingsEncDetected  EventName = "settings_enc_detected"

	DNSUpdate EventName = "dns_update"
)

// DataEvent
const (
	ProcessBegin  EventName = "process_begin"
	ProcessFinish EventName = "process_finish"

	LoginStart   EventName = "login_start"
	LoginError   EventName = "login_error"
	LoginSuccess EventName = "login_success"
	LoginFailed  EventName = "login_failed"
)

type Event interface {
	implementWebhookEvent()
	GetID() uint
}

type EventType string

const (
	TypeActionEvent EventType = "action"
	TypeDataEvent   EventType = "data"
)

type BaseEvent struct {
	ID        uint      `json:"id"`
	Timestamp int64     `json:"timestamp"` // unix milli
	Name      EventName `json:"name"`
	EventType EventType `json:"event_type"`
}

func (BaseEvent) implementWebhookEvent() {}

func (ev BaseEvent) GetID() uint {
	return ev.ID
}

var EventID atomic.Uint64

func NewBaseEvent(name EventName, _type EventType) BaseEvent {
	return BaseEvent{
		ID:        uint(EventID.Add(1)),
		Timestamp: time.Now().UnixMilli(),
		Name:      name,
		EventType: _type,
	}
}

type ActionEventStatus string

const (
	ActionEventStatusSuccess ActionEventStatus = "success"
	ActionEventStatusFailure ActionEventStatus = "failure"
)

type ActionName string

type ActionEvent struct {
	BaseEvent
	Status     ActionEventStatus `json:"status"`
	ActionName string            `json:"action_name"`

	// existence determined by status
	Value        string `json:"value,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewActionSuccessEvent(eventName EventName, actionName, value string) ActionEvent {
	return ActionEvent{
		BaseEvent:  NewBaseEvent(eventName, TypeActionEvent),
		Status:     ActionEventStatusSuccess,
		ActionName: actionName,
		Value:      value,
	}
}

func NewActionFailureEvent(eventName EventName, actionName, errMsg string) ActionEvent {
	return ActionEvent{
		BaseEvent:    NewBaseEvent(eventName, TypeActionEvent),
		Status:       ActionEventStatusFailure,
		ActionName:   actionName,
		ErrorMessage: errMsg,
	}
}

type DataEvent struct {
	BaseEvent
	Property interface{} `json:"property"`
}

func NewDataEvent(eventName EventName, property any) DataEvent {
	return DataEvent{
		BaseEvent: NewBaseEvent(eventName, TypeDataEvent),
		Property:  property,
	}
}
