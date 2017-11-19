package event

import (
	"errors"
	"log"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/phil-inc/plib/core/util"
)

//PhilEvent defines Phil event
type PhilEvent struct {
	UserID    string                 `json:"userId,omitempty"`
	RxID      string                 `json:"rxId,omitempty"`
	RxIDs     []string               `json:"rxIds,omitempty"`
	EventType string                 `json:"eventType,omitempty"`
	EventData map[string]interface{} `json:"eventData,omitempty"`
	Date      *time.Time             `json:"date,omitempty"`
}

//NewUserEvent creates new user event based on userID and type
func NewUserEvent(userID string, eventType string) *PhilEvent {
	pe := new(PhilEvent)
	pe.UserID = userID
	pe.EventType = eventType
	pe.Date = util.NowUTC()

	return pe
}

//NewRxEvent returns new rx event
func NewRxEvent(rxID string, eventType string) *PhilEvent {
	pe := new(PhilEvent)
	pe.RxID = rxID
	pe.EventType = eventType
	pe.Date = util.NowUTC()

	return pe
}

//AddStringData adds key value data to event data
func (p *PhilEvent) AddStringData(key, value string) {
	if p.EventData == nil {
		p.EventData = make(map[string]interface{}, 0)
	}
	p.EventData[key] = value
}

//AddData adds key value data to event data
func (p *PhilEvent) AddData(key string, value interface{}) {
	if p.EventData == nil {
		p.EventData = make(map[string]interface{}, 0)
	}
	p.EventData[key] = value
}

//StringData fetch string data from events data map
func (p PhilEvent) StringData(key string) string {
	if val, ok := p.EventData[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		default:
			return ""
		}
	}
	return ""
}

//BoolData fetch string data from events data map
func (p PhilEvent) BoolData(key string) bool {
	if val, ok := p.EventData[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		default:
			return false
		}
	}
	return false
}

//MapData fetch string data from events data map
func (p PhilEvent) MapData(key string) map[string]interface{} {
	if val, ok := p.EventData[key]; ok {
		return val.(map[string]interface{})
	}
	return make(map[string]interface{})
}

//Handler interface that represents the event handler
type Handler interface {
	HandleEvent(philEvent *PhilEvent) error
	IsEventSupported(eventType string) bool
}

//ErrorHanlder - handles event processing error
type ErrorHandler interface {
	HandleError(err error, handlerName string, panicErr bool)
}

var handlers = make([]Handler, 0)
var errHandlers = make([]ErrorHandler, 0)

//RegisterHandler register event handler
func RegisterHandler(handler Handler) bool {
	exists := false
	if len(handlers) > 0 {
		for _, h := range handlers {
			if util.TypeName(h) == util.TypeName(handler) {
				exists = false
			}
		}
	}
	if !exists {
		handlers = append(handlers, handler)
	}
	return !exists
}

//RegisterHandlers register all the event handlers
func RegisterHandlers(hs []Handler) {
	for _, h := range handlers {
		handlers = append(handlers, h)
	}
}

//RegisterErrorHandler register error handler
func RegisterErrorHandler(handler ErrorHandler) {
	if handler != nil {
		errHandlers = append(errHandlers, handler)
	}
}

//MainEventHandler - top level event handler that delegates to other event handlers
type MainEventHandler struct {
	PhilEvent *PhilEvent
}

//HandleEvent - handles main event
func (e MainEventHandler) HandleEvent(philEvent *PhilEvent) error {
	errList := make([]error, 0)
	for _, h := range handlers {
		newHandler := reflect.Zero(reflect.TypeOf(h)).Interface().(Handler)
		if newHandler.IsEventSupported(philEvent.EventType) {
			log.Printf("[INFO] event handled by %s", util.TypeName(newHandler))
			err := newHandler.HandleEvent(philEvent)
			if err != nil {
				eh := errHandlers[0]
				eh.HandleError(err, util.TypeName(newHandler), false)
				errList = append(errList, err)
			}
		}
	}
	if len(errList) > 0 {
		return util.WrapErrors("Error handling event", errList)
	}
	return nil
}

//IsEventSupported returns true for all event types since it's the main handler
func (e MainEventHandler) IsEventSupported(eventType string) bool {
	return true
}

//Run implements jobrunner.Job interface
func (e MainEventHandler) Run() {
	err := e.HandleEvent(e.PhilEvent)
	if err != nil {
		log.Printf("[ERROR][EVENT.%s]:: %s", e.PhilEvent.EventType, err)
	}
}

//Publish used for handling events generated locally
func Publish(philEvent *PhilEvent) {
	log.Printf("[INFO][EVENT] Phil event payload: %s", util.ToJSON(philEvent))
	HandleEvent(philEvent)
}

//HandleEvent handles phil event
func HandleEvent(philEvent *PhilEvent) {
	handler := new(MainEventHandler)
	handler.PhilEvent = philEvent

	//start event handler in a diferent goroutine
	jobrunner.Now(handler)
}

//HandlePanic handles panic for event handlers
func HandlePanic(eventHandlerName string) {
	if r := recover(); r != nil {
		log.Printf("PANIC: %s", debug.Stack())
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("Unknown panic")
		}
		if err != nil {
			eh := errHandlers[0]
			eh.HandleError(err, eventHandlerName, true)
		}
	}
}
