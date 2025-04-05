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
	Reality              EventName = "reality"

	Login     EventName = "login"
	DNSUpdate EventName = "dns_update"
)

// DataEvent
const (
	ProcessBegin  EventName = "process_begin"
	ProcessFinish EventName = "process_finish"

	LoginStart  EventName = "login_start"
	LoginFailed EventName = "login_failed"

	ClientIPDetected EventName = "client_ip_detected"
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
	ID           uint      `json:"id"`
	Timestamp    int64     `json:"timestamp"` // unix milli
	Name         EventName `json:"name"`
	EventType    EventType `json:"event_type"`
	EventContext string    `json:"event_context,omitempty"`
}

func (BaseEvent) implementWebhookEvent() {}

func (ev BaseEvent) GetID() uint {
	return ev.ID
}

var EventID atomic.Uint64

func NewBaseEvent(evName EventName, _type EventType, evContext string) BaseEvent {
	return BaseEvent{
		ID:           uint(EventID.Add(1)),
		Timestamp:    time.Now().UnixMilli(),
		Name:         evName,
		EventType:    _type,
		EventContext: evContext,
	}
}

type Property []PropertyElement

func NewProperty(el ...PropertyElement) Property {
	return el
}

type PropertyElement struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func (prop Property) Add(val ...PropertyElement) Property {
	return append(prop, val...)
}

type DataEvent struct {
	BaseEvent
	Property Property `json:"property,omitempty"`
}

func NewDataEvent(eventName EventName, evContext string, property Property) DataEvent {
	return DataEvent{
		BaseEvent: NewBaseEvent(eventName, TypeDataEvent, evContext),
		Property:  property,
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
	Status   ActionEventStatus `json:"status"`
	Property Property          `json:"property,omitempty"`

	// existence determined by status
	Value        string `json:"value,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewActionSuccessEvent(evName EventName, evContext string, property Property, value string) ActionEvent {
	return ActionEvent{
		BaseEvent: NewBaseEvent(evName, TypeActionEvent, evContext),
		Status:    ActionEventStatusSuccess,
		Property:  property,
		Value:     value,
	}
}

func NewActionFailureEvent(eventName EventName, evContext string, property Property, errMsg string) ActionEvent {
	return ActionEvent{
		BaseEvent:    NewBaseEvent(eventName, TypeActionEvent, evContext),
		Status:       ActionEventStatusFailure,
		Property:     property,
		ErrorMessage: errMsg,
	}
}
