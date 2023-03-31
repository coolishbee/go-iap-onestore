package onestore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SandboxURL      string = "https://sbpp.onestore.co.kr"
	ProductionURL   string = "https://apis.onestore.co.kr"
	ContentTypeForm string = "application/x-www-form-urlencoded"
	ContentTypeJson string = "application/json"
)

type Client struct {
	OAuthURL     string
	VerifyURL    string
	URL          string
	ClientID     string
	ClientSecret string
	httpCli      *http.Client
}

type OAuthResp struct {
	Client_id    string `json:"client_id"`
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type IAPResponse struct {
	// ConsumptionState: Consumption status of purchased products:
	// 0. non-consumption 1. consumption
	ConsumptionState int `json:"consumptionState"`

	// DeveloperPayload: Payment unique identifier provided by the developer.
	DeveloperPayload string `json:"developerPayload"`

	// PurchaseState: The purchase state of the order. Possible values are:
	// 0. Purchased 1. Canceled
	PurchaseState int `json:"purchaseState"`

	// PurchaseTime: The time the product was purchased, in
	// milliseconds since the epoch (Jan 1, 1970).
	PurchaseTime int64 `json:"purchaseTime"`

	// PurchaseId: The purchase token generated to identify this purchase.
	PurchaseId string `json:"purchaseId"`

	// AcknowledgeState: The acknowledgement state of the inapp product.
	// Possible values are: 0. Yet to be acknowledged 1. Acknowledged
	AcknowledgeState int `json:"acknowledgeState"`

	// Quantity: The quantity associated with the purchase of the inapp
	// product.
	Quantity int `json:"quantity"`
}

type IAPResponseError struct {
	IAPError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func getURL(purchaseToken string) string {
	if strings.Contains(purchaseToken, "SANDBOX") {
		return SandboxURL
	}
	return ProductionURL
}

func New(client_id, client_secret, purchaseToken string) *Client {
	client := &Client{
		URL:          getURL(purchaseToken),
		ClientID:     client_id,
		ClientSecret: client_secret,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return client
}

func (c *Client) Verify(
	ctx context.Context,
	packageName string,
	productID string,
	token string,
) (IAPResponse, error) {
	result := IAPResponse{}
	oauthResp := OAuthResp{}
	//accessToken url
	oauthUrl := fmt.Sprintf("%s/v7/oauth/token", c.URL)
	formData := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"client_credentials"},
	}

	req, err := http.NewRequest("POST", oauthUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", ContentTypeForm)
	req = req.WithContext(ctx)
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respError := IAPResponseError{}
		err = json.NewDecoder(resp.Body).Decode(&respError)
		if err != nil {
			fmt.Println(err.Error())
			return result, err
		}
		return result, fmt.Errorf("code : %s, msg : %s", respError.IAPError.Code, respError.IAPError.Message)
	}

	//respBody, _ := io.ReadAll(resp.Body)
	//log.Println(string(text))
	err = json.NewDecoder(resp.Body).Decode(&oauthResp)
	if err != nil {
		return result, err
	}

	//getPurchaseDetails url
	verifyUrl := fmt.Sprintf("%s/v7/apps/%s/purchases/inapp/products/%s/%s", c.URL, packageName, productID, token)
	req, err = http.NewRequest("GET", verifyUrl, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthResp.Access_token))
	req.Header.Set("Content-Type", ContentTypeJson)
	req = req.WithContext(ctx)

	resp, err = c.httpCli.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respError := IAPResponseError{}
		err = json.NewDecoder(resp.Body).Decode(&respError)
		if err != nil {
			fmt.Println(err.Error())
			return result, err
		}
		return result, fmt.Errorf("code : %s, msg : %s", respError.IAPError.Code, respError.IAPError.Message)
	}
	// respBody, _ := io.ReadAll(resp.Body)
	// log.Println(string(respBody))
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
