package application

import (
	"testing"

	"github.com/milicns/company-manager/company-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMock struct {
}

func TestValidator(t *testing.T) {
	var tests = []model.Company{
		{
			Id:             primitive.NewObjectID(),
			Name:           "Tech Innovators Inc.",
			Description:    "A leading tech company specializing in innovative solutions.",
			EmployeeAmount: 0,
			Registered:     true,
			Type:           model.NonProfit,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "Manufacture Masters Ltd.",
			Description:    "Description.",
			EmployeeAmount: 2000,
			Registered:     true,
			Type:           model.Cooperative,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "Lc Limited",
			Description:    "Description.",
			EmployeeAmount: 12,
			Registered:     false,
			Type:           model.SoleProprietorship,
		},
	}
	for _, tt := range tests {
		t.Run("validator test", func(t *testing.T) {
			valErr := validateCreate(tt)

			if valErr != nil {
				errors := valErr.Errors
				if len(errors) > 1 {
					t.Errorf("Expected 1 or 0 errors but got %d errors", len(errors))
				}
			}
		})
	}
}
