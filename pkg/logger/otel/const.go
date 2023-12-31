package otel

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Option func(l LogMessage) LogMessage

type LogMessage struct {
	Id          uuid.UUID `json:"uuid"`
	ServiceName string    `json:"serviceName"`
	TraceInfo   struct {
		TraceId string `json:"traceID"`
		SpanId  string `json:"spanID"`
	}
	Event struct {
		IsAuditEvent bool `json:"isAuditEvent"`
		EventInfo    struct {
			EventName        string `json:"eventName"`
			EventTriggerType string `json:"eventTriggerType"`
			Billable         bool   `json:"billable"`
		}
		EventResource interface{} `json:"eventResource"`
		EventResult   interface{} `json:"eventResult"`
		EventMessage  string      `json:"eventMessage"`
	}
	ErrorMessage string `json:"errorMessage"`
	Metadata     interface{}
}

func SetEventResource(resource interface{}) Option {
	return func(l LogMessage) LogMessage {
		l.Event.EventResource = resource
		return l
	}
}

func SetEventResult(result interface{}) Option {
	return func(l LogMessage) LogMessage {
		l.Event.EventResult = result
		return l
	}
}

func SetErrorMessage(e string) Option {
	return func(l LogMessage) LogMessage {
		l.ErrorMessage = e
		return l
	}
}

func SetMetadata(m string) Option {
	return func(l LogMessage) LogMessage {
		l.Metadata = m
		return l
	}
}

func NewLogMessage(
	span trace.Span,
	isAuditEvent bool,
	eventName string,
	eventTriggerType string,
	eventMessage string,
	billable bool,
	options ...Option,
) []byte {
	logMessage := LogMessage{}
	logMessage.Id, _ = uuid.NewV4()
	logMessage.ServiceName = "controller"
	logMessage.TraceInfo = struct {
		TraceId string "json:\"traceID\""
		SpanId  string "json:\"spanID\""
	}{
		TraceId: span.SpanContext().TraceID().String(),
		SpanId:  span.SpanContext().SpanID().String(),
	}
	logMessage.Event = struct {
		IsAuditEvent bool "json:\"isAuditEvent\""
		EventInfo    struct {
			EventName        string "json:\"eventName\""
			EventTriggerType string "json:\"eventTriggerType\""
			Billable         bool   "json:\"billable\""
		}
		EventResource interface{} "json:\"eventResource\""
		EventResult   interface{} "json:\"eventResult\""
		EventMessage  string      "json:\"eventMessage\""
	}{
		IsAuditEvent: isAuditEvent,
		EventInfo: struct {
			EventName        string "json:\"eventName\""
			EventTriggerType string "json:\"eventTriggerType\""
			Billable         bool   "json:\"billable\""
		}{
			EventName:        eventName,
			EventTriggerType: eventTriggerType,
			Billable:         billable,
		},
		EventMessage: eventMessage,
	}

	for _, o := range options {
		logMessage = o(logMessage)
	}

	bLogMessage, _ := json.Marshal(logMessage)

	return bLogMessage
}
