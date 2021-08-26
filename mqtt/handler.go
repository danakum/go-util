package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/danakum/go-util/error-handler"
	"github.com/danakum/go-util/log"
	PahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"
)

var currentSubscription subscription

type EventType string

type Header struct {
	Type      string `json:"type"`
	Version   int    `json:"version"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Expiry    int64  `json:"expiry,omitempty"`
	MessageID int64  `json:"message_id,omitempty"`
}

type Qos int

const (
	AT_MOST_ONCE Qos = iota
	AT_LEAST_ONCE
	EXACTLY_ONCE
)

type Event interface {
	Type() EventType
	Header() Header
	Body() interface{}
	Expired() bool
}

type MqttEventHandler struct {
	clusterId string
	mappings      map[EventType]func() Event
	subscriptions map[*PahoMqtt.Client]subscription
	client PahoMqtt.Client
	mu            *sync.RWMutex
}

type subscription struct {
	typ     EventType
	topic   string
	qos     Qos
	handler func(event Event) error
	client  *PahoMqtt.Client
}

func NewMqttHandler(clientId string, filePath string, clusterId string) *MqttEventHandler {
	//Init mqtt connection

	h := &MqttEventHandler{
		clusterId:clusterId,
		mappings:      make(map[EventType]func() Event, 0),
		subscriptions: make(map[*PahoMqtt.Client]subscription, 0),
		mu:            &sync.RWMutex{},
	}

	h.client = Init(clientId, filePath, func(c PahoMqtt.Client) {
		if currentSubscription.handler != nil {
			err := h.Subscribe(currentSubscription.typ, currentSubscription.topic, currentSubscription.qos, currentSubscription.handler)
			if err != nil {
				log.Error(`Re-subscribe error`)
				return
			}
			log.Info(`Re-subscribed`)
		}
	})

	return h
}

func (h *MqttEventHandler) RegisterEvent(typ EventType, handler func() Event) error {

	_, ok := h.mappings[typ]
	if ok {
		return errors.New(`Mqtt: event type ` + string(typ) + ` already registered`)
	}

	h.mappings[typ] = handler

	return nil

}

func (h *MqttEventHandler) Subscribe(typ EventType, topic string, qos Qos, handler func(event Event) error) (err error) {

	evGetter, ok := h.mappings[typ]
	if !ok {
		return errors.New(`mqtt: Event handler dose not exist for type ` + string(typ))
	}

	ev := evGetter()

	if token := h.client.Subscribe(topic, byte(qos), func(c PahoMqtt.Client, message PahoMqtt.Message) {
		err = json.Unmarshal(message.Payload(), ev)
		if err != nil {
			log.Error(err)
			return
		}

		createdAt := time.Unix(0, ev.Header().CreatedAt)
		MeasureEndToEndLatency(createdAt, h.clusterId, topic)
		CountConsumed(h.clusterId, topic)

		timeTaken := time.Now().Sub(time.Unix(0, ev.Header().CreatedAt))
		log.Trace(fmt.Sprintf(`%s Received after %v miliseconds`, typ, timeTaken))

		//Handle event on a separate go routine
		go func() {
			if err = handler(ev); err != nil && !error_handler.IsDomain(err) {
				log.Error(`Mqtt Event handler failed for : `, `event`, ev.Type(), `err: `, err)
			}
		}()

	}); token.Wait() && token.Error() != nil {
		log.Error(`Cannot subscribed to event : `, token.Error())
	} else {
		h.mu.Lock()
		currentSubscription = subscription{
			topic:   topic,
			typ:     typ,
			qos:     qos,
			handler: handler,
		}
		h.mu.Unlock()
		log.Info(`Mqtt subscription enabled for topic ` + topic)
	}

	return
}

func (h *MqttEventHandler) Publish(event Event, topic string, qos Qos, retained bool) (err error) {

	payload, err := json.Marshal(event)
	if err != nil {
		log.Error(`Cannot marshal event : `, err)
		return
	}

	begin := time.Now()
	if token := h.client.Publish(topic, byte(qos), retained, payload); token.Wait() && token.Error() != nil {
		log.Error(`Cannot publish to event : `, token.Error())
		CountProducerErrors(token.Error(), h.clusterId, topic)
		return token.Error()
	}

	MeasureProducerLatency(begin, h.clusterId, topic)
	CountProduced(h.clusterId, topic)
	return
}
