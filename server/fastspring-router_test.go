package main

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFastSpringCRUD(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	fastSpringAccountID := uuid.New().String()
	GetSettingsRepository().Set(org.ID, SettingFastSpringAccountID.Name, fastSpringAccountID)

	// Activate
	fastSpringSubscriptionID := uuid.New().String()
	eventID := uuid.New().String()
	now := time.Now().UnixNano()/1000000 - (60000 * 60)
	payload := `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.activated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + fastSpringSubscriptionID + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 2,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req := newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, eventID+"\n", res.Body.String())
	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, 150, maxUsers)
	isActive, _ := GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, true, isActive)
	sID, _ := GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, fastSpringSubscriptionID, sID)

	// Update
	eventID = uuid.New().String()
	now = time.Now().UnixNano()/1000000 - (60000 * 30)
	payload = `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.updated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + fastSpringSubscriptionID + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 5,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req = newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, eventID+"\n", res.Body.String())
	maxUsers, _ = GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, 300, maxUsers)
	isActive, _ = GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, true, isActive)
	sID, _ = GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, fastSpringSubscriptionID, sID)

	// Deactivate
	eventID = uuid.New().String()
	now = time.Now().UnixNano()/1000000 - (60000 * 10)
	payload = `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.deactivated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + fastSpringSubscriptionID + `",
					"active": false,
					"state": "deactivated",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 5,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req = newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, eventID+"\n", res.Body.String())
	maxUsers, _ = GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, SettingDefaultSubscriptionMaxUsers, maxUsers)
	isActive, _ = GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, false, isActive)
	sID, _ = GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, "", sID)

}

func TestFastSpringInvalidProduct(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	fastSpringAccountID := uuid.New().String()
	GetSettingsRepository().Set(org.ID, SettingFastSpringAccountID.Name, fastSpringAccountID)

	// Activate
	eventID := uuid.New().String()
	now := time.Now().UnixNano()/1000000 - (60000 * 60)
	payload := `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.activated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + uuid.New().String() + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 2,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "fantasy-name",
						"parent": ""
					}
				}
			}
		]
	}`
	req := newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, "", res.Body.String())

	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, SettingDefaultSubscriptionMaxUsers, maxUsers)
	isActive, _ := GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, false, isActive)
	sID, _ := GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, "", sID)
}

func TestFastSpringInvalidAccountID(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	fastSpringAccountID := uuid.New().String()
	GetSettingsRepository().Set(org.ID, SettingFastSpringAccountID.Name, fastSpringAccountID)

	// Activate
	eventID := uuid.New().String()
	now := time.Now().UnixNano()/1000000 - (60000 * 60)
	payload := `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.activated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + uuid.New().String() + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 2,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + uuid.New().String() + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req := newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, "", res.Body.String())

	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, SettingDefaultSubscriptionMaxUsers, maxUsers)
	isActive, _ := GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, false, isActive)
	sID, _ := GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, "", sID)
}

func TestFastSpringEventAlreadyProcessed(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	fastSpringAccountID := uuid.New().String()
	GetSettingsRepository().Set(org.ID, SettingFastSpringAccountID.Name, fastSpringAccountID)

	// Activate
	fastSpringSubscriptionID := uuid.New().String()
	eventID := uuid.New().String()
	now := time.Now().UnixNano()/1000000 - (60000 * 60)
	payload := `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.activated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + fastSpringSubscriptionID + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 2,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req := newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, eventID+"\n", res.Body.String())
	maxUsers, _ := GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, 150, maxUsers)
	isActive, _ := GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, true, isActive)
	sID, _ := GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, fastSpringSubscriptionID, sID)

	// Send again
	req = newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	res = executeTestRequest(req)
	checkTestResponseCode(t, http.StatusAccepted, res.Code)
	checkTestString(t, "", res.Body.String())
	maxUsers, _ = GetSettingsRepository().GetInt(org.ID, SettingSubscriptionMaxUsers.Name)
	checkTestInt(t, 150, maxUsers)
	isActive, _ = GetSettingsRepository().GetBool(org.ID, SettingActiveSubscription.Name)
	checkTestBool(t, true, isActive)
	sID, _ = GetSettingsRepository().Get(org.ID, SettingFastSpringSubscriptionID.Name)
	checkTestString(t, fastSpringSubscriptionID, sID)
}

func TestFastSpringInvalidSignature(t *testing.T) {
	clearTestDB()
	org := createTestOrg("test.com")
	fastSpringAccountID := uuid.New().String()
	GetSettingsRepository().Set(org.ID, SettingFastSpringAccountID.Name, fastSpringAccountID)
	GetConfig().FastSpringValidate = true
	defer func() {
		GetConfig().FastSpringValidate = false
	}()

	// Activate
	fastSpringSubscriptionID := uuid.New().String()
	eventID := uuid.New().String()
	now := time.Now().UnixNano()/1000000 - (60000 * 60)
	payload := `{
		"events": [
			{
				"id": "` + eventID + `",
				"type": "subscription.activated",
				"live": true,
				"processed": false,
				"created": ` + strconv.FormatInt(now, 10) + `,
				"data": {
					"id": "` + fastSpringSubscriptionID + `",
					"active": true,
					"state": "active",
					"changed": ` + strconv.FormatInt(now, 10) + `,
					"currency": "EUR",
					"quantity": 2,
					"price": 100.0,
					"begin": ` + strconv.FormatInt(now, 10) + `,
					"nextChargeDate": ` + strconv.FormatInt(now+(60000*60*24*30), 10) + `,
					"account": {
						"id": "` + fastSpringAccountID + `",
						"contact": {
							"first": "Foo",
							"last": "Bar",
							"email": "foo@bar.com",
							"company": "Foo Bar Ltd.",
							"phone": ""
						}
					},
					"product": {
						"product": "` + FastSpringProduct50Users + `",
						"parent": ""
					}
				}
			}
		]
	}`
	req := newHTTPRequest("POST", "/fastspring/webhook", "", bytes.NewBufferString(payload))
	req.Header.Set("X-FS-Signature", uuid.New().String())
	res := executeTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}
