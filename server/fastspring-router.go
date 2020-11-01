package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type FastSpringRouter struct {
}

const (
	FastSpringEventOrderCompleted              = "order.completed"
	FastSpringEventSubscriptionActivated       = "subscription.activated"
	FastSpringEventSubscriptionCanceled        = "subscription.canceled"
	FastSpringEventSubscriptionChargeCompleted = "subscription.charge.completed"
	FastSpringEventSubscriptionChargeFailed    = "subscription.charge.failed"
	FastSpringEventSubscriptionDeactivated     = "subscription.deactivated"
	FastSpringEventSubscriptionPaymentOverdue  = "subscription.payment.overdue"
	FastSpringEventSubscriptionPaymentReminder = "subscription.payment.reminder"
	FastSpringEventSubscriptionTrialReminder   = "subscription.trial.reminder"
	FastSpringEventSubscriptionUncanceled      = "subscription.uncanceled"
	FastSpringEventSubscriptionUpdated         = "subscription.updated"
)

const (
	FastSpringProduct50Users = "50-users"
)

type FastSpringWebhookRequest struct {
	Events []FastSpringSubscriptionEvent `json:"events" validate:"required"`
}

type FastSpringSubscriptionEvent struct {
	EventID    string                          `json:"id" validate:"required"`
	Type       string                          `json:"type" validate:"required"`
	Live       bool                            `json:"live"`
	Processed  bool                            `json:"processed"`
	CreateTime int64                           `json:"created"`
	Data       FastSpringSubscriptionEventData `json:"data" validate:"required"`
}

type FastSpringSubscriptionEventData struct {
	SubscriptionID    string            `json:"id" validate:"required"`
	Active            bool              `json:"active"`
	State             string            `json:"state"`
	ChangeTime        int64             `json:"changed"`
	Currency          string            `json:"currency"`
	Quantity          int               `json:"quantity"`
	Price             float32           `json:"price"`
	SubscriptionStart int64             `json:"begin"`
	NextChargeDate    int64             `json:"nextChargeDate"`
	Account           FastSpringAccount `json:"account" validate:"required"`
	Product           FastSpringProduct `json:"product" validate:"required"`
}

type FastSpringAccount struct {
	AccountID string            `json:"id" validate:"required"`
	Contact   FastSpringContact `json:"contact" validate:"required"`
}

type FastSpringProduct struct {
	ProductID       string `json:"product" validate:"required"`
	ParentProductID string `json:"parent"`
}

func (router *FastSpringRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/webhook", router.webhook).Methods("POST")
}

func (router *FastSpringRouter) webhook(w http.ResponseWriter, r *http.Request) {
	m, err := router.getValidateRequest(r)
	if err != nil {
		log.Println(err)
		SendBadRequest(w)
		return
	}
	processedEvents := ""
	for _, event := range m.Events {
		if err := router.processEvent(event); err == nil {
			processedEvents += event.EventID + "\n"
		}
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(processedEvents))
}

func (router *FastSpringRouter) processEvent(event FastSpringSubscriptionEvent) error {
	//s, _ := json.MarshalIndent(event, "", "\t")
	//log.Println(string(s))
	persistEvent, err := router.prepareSubscriptionEvent(event)
	if err != nil {
		return err
	}
	existingEvent, err := GetSubscriptionRepository().GetProcessedByBrokerEventID(persistEvent.BrokerEventID)
	if existingEvent != nil {
		return fmt.Errorf("Event with ID %s already processed", persistEvent.BrokerEventID)
	}
	defer GetSubscriptionRepository().Create(persistEvent)
	switch eventType := strings.ToLower(event.Type); eventType {
	case FastSpringEventSubscriptionActivated:
		persistEvent.EventType = SubscriptionEventActivate
		return router.processActivateEvent(persistEvent)
	case FastSpringEventSubscriptionDeactivated:
		persistEvent.EventType = SubscriptionEventDeactivate
		return router.processDeactivateEvent(persistEvent)
	case FastSpringEventSubscriptionUpdated:
		persistEvent.EventType = SubscriptionEventUpdate
		return router.processUpdateEvent(persistEvent)
	default:
		return nil
	}
}

func (router *FastSpringRouter) processActivateEvent(event *SubscriptionEvent) error {
	GetSettingsRepository().Set(event.OrganizationID, SettingFastSpringSubscriptionID.Name, event.BrokerSubscriptionID)
	GetSettingsRepository().Set(event.OrganizationID, SettingActiveSubscription.Name, "1")
	GetSettingsRepository().Set(event.OrganizationID, SettingSubscriptionMaxUsers.Name, strconv.Itoa(event.MaxUsers))
	event.Processed = true
	return nil
}

func (router *FastSpringRouter) processDeactivateEvent(event *SubscriptionEvent) error {
	GetSettingsRepository().Set(event.OrganizationID, SettingFastSpringSubscriptionID.Name, "")
	GetSettingsRepository().Set(event.OrganizationID, SettingActiveSubscription.Name, "0")
	GetSettingsRepository().Set(event.OrganizationID, SettingSubscriptionMaxUsers.Name, strconv.Itoa(SettingDefaultSubscriptionMaxUsers))
	event.Processed = true
	return nil
}

func (router *FastSpringRouter) processUpdateEvent(event *SubscriptionEvent) error {
	GetSettingsRepository().Set(event.OrganizationID, SettingSubscriptionMaxUsers.Name, strconv.Itoa(event.MaxUsers))
	event.Processed = true
	return nil
}

func (router *FastSpringRouter) prepareSubscriptionEvent(event FastSpringSubscriptionEvent) (*SubscriptionEvent, error) {
	org := router.getOrgByAccountID(event.Data.Account.AccountID)
	if org == nil {
		return nil, fmt.Errorf("Organization not found for account with FastSpring Account ID %s", event.Data.Account.AccountID)
	}
	if event.Data.Product.ProductID != FastSpringProduct50Users {
		return nil, fmt.Errorf("Invalid product ID %s", event.Data.Product.ProductID)
	}
	e := &SubscriptionEvent{
		OrganizationID:       org.ID,
		EventTime:            router.timeFromMillis(event.CreateTime),
		ActivationTime:       router.timeFromMillis(event.Data.ChangeTime),
		BrokerEventID:        event.EventID,
		BrokerSubscriptionID: event.Data.SubscriptionID,
		BrokerCustomerID:     event.Data.Account.AccountID,
		MaxUsers:             router.getSubscriptionMaxUsers(event.Data.Quantity),
		Price:                event.Data.Price,
		Processed:            false,
	}
	return e, nil
}

func (router *FastSpringRouter) getSubscriptionMaxUsers(orderQuantity int) int {
	return (orderQuantity * 50) + SettingDefaultSubscriptionMaxUsers
}

func (router *FastSpringRouter) getOrgByAccountID(fastSpringAccountID string) *Organization {
	orgIDs, err := GetSettingsRepository().GetOrganizationIDsByValue(SettingFastSpringAccountID.Name, fastSpringAccountID)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(orgIDs) == 0 {
		return nil
	}
	org, err := GetOrganizationRepository().GetOne(orgIDs[0])
	if err != nil {
		log.Println(err)
		return nil
	}
	return org
}

func (router *FastSpringRouter) validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func (router *FastSpringRouter) getValidateRequest(r *http.Request) (*FastSpringWebhookRequest, error) {
	if r.Body == nil {
		return nil, errors.New("Body is nil")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if GetConfig().FastSpringValidate {
		signature, err := base64.StdEncoding.DecodeString(r.Header.Get("X-FS-Signature"))
		if err != nil {
			return nil, err
		}
		if !router.validMAC(body, signature, []byte(GetConfig().FastSpringHash)) {
			return nil, err
		}
	}
	//log.Println(string(body))
	var m FastSpringWebhookRequest
	if err = json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	err = GetValidator().Struct(m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (router *FastSpringRouter) timeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}
