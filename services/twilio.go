package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"boilerplate-api/internal/config"
)

type PhoneMessage struct {
	Phone   string
	Message string
}

// ErrorResponse twilio error response
type ErrorResponse struct {
	Code     uint   `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   uint   `json:"status"`
}

// SuccessResponse twilio success response
type SuccessResponse struct {
	Sid                 string      `json:"sid"`
	DateCreated         string      `json:"date_created"`
	DateUpdated         string      `json:"date_updated"`
	DateSent            interface{} `json:"date_sent"`
	AccountSid          string      `json:"account_sid"`
	To                  string      `json:"to"`
	From                string      `json:"from"`
	MessagingServiceSid interface{} `json:"messaging_service_sid"`
	Body                string      `json:"body"`
	Status              string      `json:"status"`
	NumSegments         string      `json:"num_segments"`
	NumMedia            string      `json:"num_media"`
	Direction           string      `json:"direction"`
	APIVersion          string      `json:"api_version"`
	Price               interface{} `json:"price"`
	PriceUnit           string      `json:"price_unit"`
	ErrorCode           interface{} `json:"error_code"`
	ErrorMessage        interface{} `json:"error_message"`
	URI                 string      `json:"uri"`
	SubresourceUris     struct {
		Media string `json:"media"`
	} `json:"subresource_uris"`
}

// TwilioService twilio service structure
type TwilioService struct {
	baseURL   string
	smsFrom   string
	sid       string
	authToken string
	logger    config.Logger
}

// NewTwilioService creates new twilio service
func NewTwilioService(
	env config.Env,
	logger config.Logger,
) TwilioService {
	return TwilioService{
		baseURL:   env.TwilioBaseURL,
		sid:       env.TwilioSID,
		authToken: env.TwilioAuthToken,
		smsFrom:   env.TwilioSMSFrom,
		logger:    logger,
	}
}

// SMSInput input for sms
type SMSInput struct {
	From string
	To   string
	Body string
}

func (t TwilioService) SendSMS(input SMSInput) (*SuccessResponse, *ErrorResponse, error) {
	url := fmt.Sprintf("%s/Accounts/%s/Messages.json", t.baseURL, t.sid)

	t.logger.Info(url)

	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("From", input.From)
	_ = writer.WriteField("To", input.To)
	_ = writer.WriteField("Body", input.Body)

	err := writer.Close()
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, nil, err
	}

	token := fmt.Sprintf("Basic %s", t.getBasicToken())
	t.logger.Info(token)

	req.Header.Add("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusCreated {
		result := ErrorResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, nil, err
		}

		return nil, &result, nil
	}

	result := SuccessResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}

	return &result, nil, nil

}

func (t TwilioService) getBasicToken() string {
	token := fmt.Sprintf("%s:%s", t.sid, t.authToken)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (t TwilioService) MessageSuccess(payload PhoneMessage) error {
	_, twilioErr, err := t.SendSMS(SMSInput{
		From: t.smsFrom,
		To:   payload.Phone,
		Body: payload.Message,
	})
	if err != nil {
		t.logger.Error("user message send error: ", err.Error())
		return err
	}
	if twilioErr != nil {
		t.logger.Errorf("twilio message send error: %+v \n", twilioErr)
		return err
	}
	return nil
}
