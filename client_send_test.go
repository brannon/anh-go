package anh

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func withTestEndpoint(server *httptest.Server) ClientOption {
	return func(c *Client) error {
		c.endpoint = server.URL
		return nil
	}
}

type TestSuite_Client_SendDirectNotification struct {
	suite.Suite

	server *httptest.Server
	client *Client

	requestURL     *url.URL
	requestMethod  string
	requestHeaders http.Header
	requestBody    []byte

	responseStatus   int
	responseLocation string
}

func (suite *TestSuite_Client_SendDirectNotification) SetupTest() {
	var err error

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.requestURL = r.URL
		suite.requestMethod = r.Method
		suite.requestHeaders = r.Header
		suite.requestBody, _ = ioutil.ReadAll(r.Body)

		w.Header().Set("Location", suite.responseLocation)
		w.WriteHeader(suite.responseStatus)
	}))

	httpClient := suite.server.Client()

	suite.client, err = NewClient("hubName", withTestEndpoint(suite.server), WithHttpClient(httpClient))
	assert.NoError(suite.T(), err)

	// Setup defaults
	suite.responseStatus = http.StatusCreated
}

func (suite *TestSuite_Client_SendDirectNotification) TearDownTest() {
	suite.server.Close()
}

func (suite *TestSuite_Client_SendDirectNotification) TestRequest() {
	notification := &AppleNotification{}
	deviceToken := "device-token"

	_, err := suite.client.SendDirectNotification(context.Background(), notification, deviceToken)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "POST", suite.requestMethod)
	assert.Equal(suite.T(), "/hubName/messages", suite.requestURL.Path)
	assert.Equal(suite.T(), "true", suite.requestURL.Query().Get("direct"))
	assert.Equal(suite.T(), "application/json", suite.requestHeaders.Get("Content-Type"))
	assert.Equal(suite.T(), "device-token", suite.requestHeaders.Get("ServiceBusNotification-DeviceHandle"))
	assert.Equal(suite.T(), "apple", suite.requestHeaders.Get("ServiceBusNotification-Format"))
}

func (suite *TestSuite_Client_SendDirectNotification) TestRequest_Body() {
	notification := &AppleNotification{}
	notification.SetBodyFromString(`{"aps":{"alert":"Hello World!"}}`)
	deviceToken := "device-token"

	_, err := suite.client.SendDirectNotification(context.Background(), notification, deviceToken)

	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"aps":{"alert":"Hello World!"}}`, string(suite.requestBody))
}

func (suite *TestSuite_Client_SendDirectNotification) TestResponse_Success() {
	suite.responseLocation = "https://hubName/messages/notifications/notification-123"
	notification := &AppleNotification{}
	deviceToken := "device-token"

	notificationResult, err := suite.client.SendDirectNotification(context.Background(), notification, deviceToken)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), notificationResult)
	assert.Equal(suite.T(), "notification-123", notificationResult.GetNotificationId())
}

func (suite *TestSuite_Client_SendDirectNotification) TestResponse_HTTPError() {
	suite.responseStatus = http.StatusUnauthorized
	notification := &AppleNotification{}
	deviceToken := "device-token"

	_, err := suite.client.SendDirectNotification(context.Background(), notification, deviceToken)

	assert.ErrorIs(suite.T(), err, ErrInvalidCredentials)
}

func Test_Client_SendDirectNotification(t *testing.T) {
	suite.Run(t, new(TestSuite_Client_SendDirectNotification))
}
