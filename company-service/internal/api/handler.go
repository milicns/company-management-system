package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/milicns/company-manager/company-service/internal/application"
	"github.com/milicns/company-manager/company-service/internal/model"
	"github.com/milicns/company-manager/company-service/internal/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service  *application.CompanyService
	producer application.Producer
}

func NewHandler(service *application.CompanyService, producer application.Producer) *Handler {
	return &Handler{
		service:  service,
		producer: producer,
	}
}

func (handler *Handler) Create(writer http.ResponseWriter, req *http.Request) {
	var company model.Company
	err := json.NewDecoder(req.Body).Decode(&company)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	objId, err := handler.service.Create(company)
	if err != nil {
		if valErr, ok := err.(application.ValidationError); ok {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			if encodeErr := json.NewEncoder(writer).Encode(valErr); encodeErr != nil {
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error encoding error to JSON: %v", encodeErr)
				return
			}
			return
		}
		http.Error(writer, err.Error(), http.StatusExpectationFailed)
		return
	}

	id := objId.Hex()
	response := struct {
		Id string `json:"id"`
	}{Id: id}

	go func() {
		event := utils.NewEvent(company.Name, "create")
		handler.producer.Produce(event)
	}()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if encodeErr := json.NewEncoder(writer).Encode(response); encodeErr != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding id to JSON: %v", encodeErr)
		return
	}
}

func (handler *Handler) GetOne(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	company, err := handler.service.GetOne(objId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusExpectationFailed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if encodeErr := json.NewEncoder(writer).Encode(*company); encodeErr != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding company to JSON: %v", encodeErr)
		return
	}
}

func (handler *Handler) Patch(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var patchData map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&patchData); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	patchData["id"] = objId
	err = handler.service.Patch(patchData)
	if err != nil {
		if valErr, ok := err.(application.ValidationError); ok {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)

			if encodeErr := json.NewEncoder(writer).Encode(valErr); encodeErr != nil {
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error encoding error to JSON: %v", encodeErr)
				return
			}
			return
		}
		http.Error(writer, err.Error(), http.StatusExpectationFailed)
		return
	}

	go func() {
		event := utils.NewEvent(id, "update")
		handler.producer.Produce(event)
	}()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *Handler) Delete(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deleted, err := handler.service.Delete(objId)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusExpectationFailed)
		return
	}

	go func() {
		event := utils.NewEvent(deleted.Id.Hex(), "delete")
		handler.producer.Produce(event)
	}()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
