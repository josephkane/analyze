package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/storage"
)

func TestPluginHandler_ReturnResultsSuccessfully(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	fixturePlugins1 := newPluginFixture("123456798")
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockStorage(t, map[string]string{
		models.PluginPrefix + "123456798": fixturePlugins1.string(),
	}), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}

	p := toPlugin(t, rr.Body)
	if reflect.DeepEqual(*p, fixturePlugins1) {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), fixturePlugins1.string())
	}
}

func TestPluginHandler_ReturnInternalError(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockBrokenStorage(t), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}

func TestPluginHandler_ReturnNotFound(t *testing.T) {
	analyzeApi := api.GetTestAPI(t)
	//TODO: create interface for logger, and use dummy logger for tests
	analyzeApi.GetPluginHandler = handlers.NewPluginHandler(storage.GetMockStorage(t, nil), logrus.New())
	server := api.NewServer(analyzeApi)
	server.ConfigureAPI()

	h := server.GetHandler()

	req, err := http.NewRequest("GET", "/api/v1/plugins/123456798", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v, body: %v", status, http.StatusOK, rr.Body.String())
	}
}
