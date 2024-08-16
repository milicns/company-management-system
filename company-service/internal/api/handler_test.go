package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/milicns/company-manager/company-service/internal/application"
	"github.com/milicns/company-manager/company-service/internal/model"
	"github.com/milicns/company-manager/company-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCompany(t *testing.T) {
	store := &MockStore{}
	service := application.NewService(store)
	producer := &MockProducer{events: make([]utils.Event, 0, 5)}
	handler := NewHandler(service, producer)
	company := model.Company{
		Name:           "LXC",
		Description:    "Description",
		EmployeeAmount: 15,
		Registered:     true,
		Type:           model.Corporations,
	}

	cmp, err := json.Marshal(company)
	if err != nil {
		t.Fatalf("Failed to marshal company to JSON: %v", err)
	}
	body := bytes.NewReader(cmp)

	server := httptest.NewServer(http.HandlerFunc(handler.Create))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", resp.StatusCode)
	}
	if len(producer.events) > 0 {
		t.Errorf("Expected to not produce")
	}
}

func TestGetCompany(t *testing.T) {
	store := &MockStore{}
	service := application.NewService(store)
	producer := &MockProducer{events: make([]utils.Event, 0, 5)}
	handler := NewHandler(service, producer)

	idStr := testObjectId.Hex()
	r := mux.NewRouter()
	r.HandleFunc("/{id}", handler.GetOne).Methods("GET")
	server := httptest.NewServer(r)

	defer server.Close()

	url := fmt.Sprintf("%s/%s", server.URL, idStr)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var company model.Company
	if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if company.Id != testObjectId {
		t.Errorf("Expected company ID %v, got %v", testObjectId, company.Id)
	}
}

func TestPatchCompany(t *testing.T) {
	store := &MockStore{}
	service := application.NewService(store)
	producer := &MockProducer{events: make([]utils.Event, 0, 5)}
	handler := NewHandler(service, producer)

	patchData := make(map[string]interface{})
	patchData["employeeamount"] = 0
	idStr := testObjectId.Hex()

	r := mux.NewRouter()
	r.HandleFunc("/{id}", handler.Patch).Methods("PATCH")
	patch, err := json.Marshal(patchData)
	if err != nil {
		t.Fatalf("Failed to marshal company to JSON: %v", err)
	}

	body := bytes.NewReader(patch)
	server := httptest.NewServer(r)
	defer server.Close()

	url := fmt.Sprintf("%s/%s", server.URL, idStr)

	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", resp.StatusCode)
	}
	if len(producer.events) > 0 {
		t.Errorf("Expected to not produce")
	}
}

func TestDeleteCompany(t *testing.T) {
	store := &MockStore{}
	service := application.NewService(store)
	producer := &MockProducer{events: make([]utils.Event, 0, 5)}
	handler := NewHandler(service, producer)

	idStr := testObjectId.Hex()

	r := mux.NewRouter()
	r.HandleFunc("/{id}", handler.Delete).Methods("DELETE")

	server := httptest.NewServer(r)
	defer server.Close()

	url := fmt.Sprintf("%s/%s", server.URL, idStr)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusExpectationFailed {
		t.Errorf("Expected status code 417, got %d", resp.StatusCode)
	}

	if len(producer.events) > 0 {
		t.Errorf("Expected to not produce")
	}
}

type MockStore struct {
}

func (s *MockStore) Create(company model.Company) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (s *MockStore) GetOne(id primitive.ObjectID) (*model.Company, error) {
	return &model.Company{
		Id:             testObjectId,
		Name:           "LXC United",
		Description:    "Description",
		EmployeeAmount: 100,
		Registered:     true,
		Type:           model.Corporations,
	}, nil
}

func (s *MockStore) Patch(patchData map[string]interface{}) error {
	return nil
}

func (s *MockStore) Delete(id primitive.ObjectID) (*model.Company, error) {
	return nil, errors.New("company was not found")
}

type MockProducer struct {
	events []utils.Event
}

func (p *MockProducer) Produce(event utils.Event) {
	p.events = append(p.events, event)
}

var testObjectId primitive.ObjectID = primitive.NewObjectID()
