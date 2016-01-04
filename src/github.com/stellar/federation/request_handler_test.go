package federation

import (
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
  "github.com/stretchr/testify/mock"
)

type MockDatabase struct {
  mock.Mock
}
func (m *MockDatabase) Get(dest interface{}, query string, args ...interface {}) error {
  a := m.Called(dest, query, args[0])
  return a.Error(0)
}

func TestRequestHandler(t *testing.T) {
  mockDatabase := new(MockDatabase)

  app := App{
    config: Config{
      Domain: "acme.com",
      FederationQuery: "FederationQuery",
      ReverseFederationQuery: "ReverseFederationQuery",
    },
    database: mockDatabase,
  }

  requestHandler := RequestHandler{config: &app.config, database: app.database}
  testServer := httptest.NewServer(http.HandlerFunc(requestHandler.Main))
  defer testServer.Close()

  Convey("Given federation request", t, func() {
    Convey("When record exists", func() {
      username := "test"
      accountId := "GD3YBOYIUVLU2VGK4GW5J3L4O5JCS626KG53OIFSXX7UOBS6NPCJIR2T"

      responseRecord := Record{}

      mockDatabase.On("Get", &responseRecord, app.config.FederationQuery, username).Return(nil).Run(func(args mock.Arguments) {
        record := args.Get(0).(*Record)
        record.Username = username
        record.AccountId = accountId
      })

      Convey("it should return correct values", func() {
        response := GetResponse(testServer, "?type=name&q="+username+"*"+app.config.Domain);
        json.Unmarshal(response, &responseRecord)

        So(responseRecord.Username,  ShouldEqual, username+"*"+app.config.Domain)
        So(responseRecord.AccountId, ShouldEqual, accountId)

        mockDatabase.AssertExpectations(t)
      })
    })

    Convey("When record does not exist", func() {
      username := "not-exist"
      responseRecord := Record{}

      mockDatabase.On("Get", &responseRecord, app.config.FederationQuery, username).Return(errors.New("sql: no rows in result set"))

      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=name&q="+username+"*"+app.config.Domain);
        CheckErrorResponse(response, "not_found", "Account not found")
        mockDatabase.AssertExpectations(t)
      })
    })

    Convey("When domain is invalid", func() {
      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=name&q=test*other.com");
        CheckErrorResponse(response, "not_found", "Incorrect domain")
        mockDatabase.AssertNotCalled(t, "Get")
      })
    })

    Convey("When domain is empty", func() {
      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=name&q=test");
        CheckErrorResponse(response, "not_found", "Incorrect domain")
        mockDatabase.AssertNotCalled(t, "Get")
      })
    })

    Convey("When no `q` param", func() {
      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=id");
        CheckErrorResponse(response, "invalid_request", "Invalid request")
        mockDatabase.AssertNotCalled(t, "Get")
      })
    })
  })

  Convey("Given reverse federation request", t, func() {
    Convey("When record exists", func() {
      username := "test"
      accountId := "GD3YBOYIUVLU2VGK4GW5J3L4O5JCS626KG53OIFSXX7UOBS6NPCJIR2T"

      responseRecord := Record{}

      mockDatabase.On("Get", &responseRecord, app.config.ReverseFederationQuery, accountId).Return(nil).Run(func(args mock.Arguments) {
        record := args.Get(0).(*Record)
        record.Username = username
        record.AccountId = accountId
      })

      Convey("it should return correct values", func() {
        response := GetResponse(testServer, "?type=id&q="+accountId);
        json.Unmarshal(response, &responseRecord)

        So(responseRecord.Username,  ShouldEqual, username+"*"+app.config.Domain)
        So(responseRecord.AccountId, ShouldEqual, accountId)

        mockDatabase.AssertExpectations(t)
      })
    })

    Convey("When record does not exist", func() {
      accountId := "GCKWDG2RWKPJNLLPLNU5PYCYN3TLKWI2SWAMSGFGSTVHCJX5P2EVMFGS"
      responseRecord := Record{}

      mockDatabase.On("Get", &responseRecord, app.config.ReverseFederationQuery, accountId).Return(errors.New("sql: no rows in result set"))

      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=id&q="+accountId);
        CheckErrorResponse(response, "not_found", "Account not found")
        mockDatabase.AssertExpectations(t)
      })
    })

    Convey("When no `q` param", func() {
      Convey("it should return error response", func() {
        response := GetResponse(testServer, "?type=id");
        CheckErrorResponse(response, "invalid_request", "Invalid request")
        mockDatabase.AssertNotCalled(t, "Get")
      })
    })
  })

  Convey("Given request with invalid type", t, func() {
    Convey("it should return error response", func() {
      response := GetResponse(testServer, "?type=invalid");
      CheckErrorResponse(response, "invalid_request", "Invalid request")
      mockDatabase.AssertNotCalled(t, "Get")
    })
  })
}

func GetResponse(testServer *httptest.Server, query string) []byte {
  res, err := http.Get(testServer.URL+query)
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

  So(errorResponse.Code,  ShouldEqual, code)
  So(errorResponse.Message, ShouldEqual, message)
}
