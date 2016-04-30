package federation

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/federation/config"
	"github.com/stellar/federation/db"
	"github.com/stretchr/testify/mock"
)

type MockDriver struct {
	mock.Mock
}

func (m *MockDriver) Init(url config.ConfigDatabase) error {
	a := m.Called(url)
	return a.Error(0)
}

func (m *MockDriver) GetByStellarAddress(name, query string) (*db.FederationRecord, error) {
	a := m.Called(name, query)
	record := a.Get(0)
	if record == nil {
		return nil, a.Error(1)
	} else {
		return a.Get(0).(*db.FederationRecord), a.Error(1)
	}
}

func (m *MockDriver) GetByAccountId(accountId, query string) (*db.ReverseFederationRecord, error) {
	a := m.Called(accountId, query)
	record := a.Get(0)
	if record == nil {
		return nil, a.Error(1)
	} else {
		return a.Get(0).(*db.ReverseFederationRecord), a.Error(1)
	}
}

func TestRequestHandler(t *testing.T) {
	mockDriver := new(MockDriver)

	app := App{
		config: config.Config{
			Domain: "acme.com",
			Queries: config.ConfigQueries{
				Federation:        "FederationQuery",
				ReverseFederation: "ReverseFederationQuery",
			},
		},
		driver: mockDriver,
	}

	requestHandler := RequestHandler{config: &app.config, driver: app.driver}
	testServer := httptest.NewServer(http.HandlerFunc(requestHandler.Main))
	defer testServer.Close()

	Convey("Given federation request", t, func() {
		Convey("When record exists", func() {
			username := "test"
			accountId := "GD3YBOYIUVLU2VGK4GW5J3L4O5JCS626KG53OIFSXX7UOBS6NPCJIR2T"

			record := db.FederationRecord{AccountId: accountId}
			mockDriver.On("GetByStellarAddress", username, app.config.Queries.Federation).Return(&record, nil)

			Convey("it should return correct values", func() {
				response := GetResponse(testServer, "?type=name&q="+username+"*"+app.config.Domain)
				var responseObject Response
				json.Unmarshal(response, &responseObject)

				So(responseObject.StellarAddress, ShouldEqual, username+"*"+app.config.Domain)
				So(responseObject.AccountId, ShouldEqual, accountId)

				mockDriver.AssertExpectations(t)
			})

		})

		Convey("When record does not exist", func() {
			username := "not-exist"

			mockDriver.On("GetByStellarAddress", username, app.config.Queries.Federation).Return(nil, nil)

			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=name&q="+username+"*"+app.config.Domain)

				CheckErrorResponse(response, "not_found", "Account not found")
				mockDriver.AssertExpectations(t)
			})
		})

		Convey("When domain is invalid", func() {
			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=name&q=test*other.com")
				CheckErrorResponse(response, "not_found", "Incorrect Domain")
				mockDriver.AssertNotCalled(t, "Get")
			})
		})

		Convey("When domain is empty", func() {
			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=name&q=test")
				CheckErrorResponse(response, "not_found", "Incorrect Domain")
				mockDriver.AssertNotCalled(t, "Get")
			})
		})

		Convey("When no `q` param", func() {
			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=id")
				CheckErrorResponse(response, "invalid_request", "Invalid request")
				mockDriver.AssertNotCalled(t, "Get")
			})
		})

	})

	Convey("Given reverse federation request", t, func() {

		Convey("When record exists", func() {
			username := "test"
			accountId := "GD3YBOYIUVLU2VGK4GW5J3L4O5JCS626KG53OIFSXX7UOBS6NPCJIR2T"

			record := db.ReverseFederationRecord{Name: username}
			mockDriver.On("GetByAccountId", accountId, app.config.Queries.ReverseFederation).Return(&record, nil)

			Convey("it should return correct values", func() {
				response := GetResponse(testServer, "?type=id&q="+accountId)
				var responseObject Response
				json.Unmarshal(response, &responseObject)

				So(responseObject.StellarAddress, ShouldEqual, username+"*"+app.config.Domain)
				So(responseObject.AccountId, ShouldEqual, accountId)

				mockDriver.AssertExpectations(t)
			})
		})

		Convey("When record does not exist", func() {
			accountId := "GCKWDG2RWKPJNLLPLNU5PYCYN3TLKWI2SWAMSGFGSTVHCJX5P2EVMFGS"

			mockDriver.On("GetByAccountId", accountId, app.config.Queries.ReverseFederation).Return(nil, nil)

			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=id&q="+accountId)
				CheckErrorResponse(response, "not_found", "Account not found")
				mockDriver.AssertExpectations(t)
			})
		})

		Convey("When no `q` param", func() {
			Convey("it should return error response", func() {
				response := GetResponse(testServer, "?type=id")
				CheckErrorResponse(response, "invalid_request", "Invalid request")
				mockDriver.AssertNotCalled(t, "Get")
			})
		})

	})

	Convey("Given request with invalid type", t, func() {
		Convey("it should return error response", func() {
			response := GetResponse(testServer, "?type=invalid")
			CheckErrorResponse(response, "invalid_request", "Invalid request")
			mockDriver.AssertNotCalled(t, "Get")
		})
	})

}

func GetResponse(testServer *httptest.Server, query string) []byte {
	res, err := http.Get(testServer.URL + query)
	if err != nil {
		panic(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	return response
}

func CheckErrorResponse(response []byte, code string, message string) {
	errorResponse := Error{}
	json.Unmarshal(response, &errorResponse)

	So(errorResponse.Code, ShouldEqual, code)
	So(errorResponse.Message, ShouldEqual, message)
}
