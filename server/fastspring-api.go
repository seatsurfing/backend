package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FastSpringAPI struct {
}

const FastSpringAPIBackendURL = "https://api.fastspring.com"
const FastSpringAPISuccess = "success"

type FastSpringCreateAccountRequest struct {
	Contact  FastSpringContact `json:"contact"`
	Language string            `json:"language"`
	Country  string            `json:"country"`
	Lookup   FastSpringLookup  `json:"lookup"`
}

type FastSpringContact struct {
	Firstname string `json:"first"`
	Lastname  string `json:"last"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`
}

type FastSpringLookup struct {
	CustomID string `json:"custom"`
}

type FastSpringCreateAccountResponse struct {
	AccountID string            `json:"account"`
	Action    string            `json:"action"`
	Result    string            `json:"result"`
	URL       string            `json:"url"`
	Errors    map[string]string `json:"error"`
}

type FastSpringAuthenticateResponse struct {
	Accounts []FastSpringCreateAccountResponse `json:"accounts"`
}

type FastSpringSessionRequest struct {
	AccountID string           `json:"account"`
	Items     []FastSpringItem `json:"items"`
}

type FastSpringItem struct {
	ProductID string `json:"product"`
	Quantity  int    `json:"quantity"`
}

type FastSpringSessionResponse struct {
	SessionID string           `json:"id"`
	Currency  string           `json:"currency"`
	Expiry    int64            `json:"expires"`
	AccountID string           `json:"account"`
	Subtotal  float32          `json:"subtotal"`
	Items     []FastSpringItem `json:"items"`
}

func (api *FastSpringAPI) CreateAccount(request *FastSpringCreateAccountRequest) (string, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	res, err := api.getAPIResponse("POST", "/accounts", bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
	}
	if res == nil {
		return "", err
	}
	defer res.Body.Close()
	var resBody *FastSpringCreateAccountResponse
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(responseData, &resBody); err != nil {
		return "", err
	}
	if strings.ToLower(resBody.Result) != FastSpringAPISuccess {
		if val, ok := resBody.Errors["custom"]; ok {
			if strings.Contains(val, "already exists") {
				lastPos := strings.LastIndex(val, "/")
				return val[lastPos+1:], nil
			}
		}
		return "", errors.New("Got result string " + resBody.Result)
	}
	return resBody.AccountID, nil
}

func (api *FastSpringAPI) Authenticate(accountID string) (string, error) {
	if accountID == "" {
		return "", errors.New("empty account ID")
	}
	res, err := api.getAPIResponse("GET", "/accounts/"+accountID+"/authenticate", nil)
	if err != nil {
		log.Println(err)
	}
	if res == nil {
		return "", err
	}
	defer res.Body.Close()
	var resBody *FastSpringAuthenticateResponse
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(responseData, &resBody); err != nil {
		return "", err
	}
	if len(resBody.Accounts) != 1 {
		return "", errors.New("got invalid account count")
	}
	account := resBody.Accounts[0]
	if strings.ToLower(account.Result) != FastSpringAPISuccess {
		return "", errors.New("got result string " + account.Result)
	}
	return account.URL + "#/subscriptions", nil
}

func (api *FastSpringAPI) StartCheckoutSession(accountID, productID string, quantity int, live bool) (string, error) {
	request := &FastSpringSessionRequest{
		AccountID: accountID,
		Items: []FastSpringItem{
			{
				ProductID: productID,
				Quantity:  quantity,
			},
		},
	}
	sessionID, err := api.startSession(request)
	if err != nil {
		return "", err
	}
	if live {
		return "https://seatsurfing.onfastspring.com/session/" + sessionID, nil
	} else {
		return "https://seatsurfing.test.onfastspring.com/session/" + sessionID, nil
	}
}

func (api *FastSpringAPI) startSession(request *FastSpringSessionRequest) (string, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	res, err := api.getAPIResponse("POST", "/sessions", bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
	}
	if res == nil {
		return "", err
	}
	defer res.Body.Close()
	var resBody *FastSpringSessionResponse
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(responseData, &resBody); err != nil {
		return "", err
	}
	return resBody.SessionID, nil
}

func (api *FastSpringAPI) getAPIResponse(method, endpoint string, request io.Reader) (*http.Response, error) {
	url := FastSpringAPIBackendURL + endpoint
	log.Println("Performing FastSpring API Request to " + url)
	req, err := http.NewRequest(method, url, request)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Seatsurfing App Backend")
	req.SetBasicAuth(GetConfig().FastSpringUsername, GetConfig().FastSpringPassword)
	if request != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return res, errors.New("Got http status code " + strconv.Itoa(res.StatusCode))
	}
	return res, nil
}
