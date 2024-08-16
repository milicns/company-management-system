package application

import (
	"fmt"

	"github.com/milicns/company-manager/company-service/internal/model"
)

type Validator struct {
	errors map[string]string
}

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

func (val *Validator) check(condition bool, field string, message string) {
	if !condition {
		val.errors[field] = message
	}
}

func (val *Validator) hasErrors() bool {
	return len(val.errors) > 0
}

func (val *Validator) getErrors() map[string]string {
	return val.errors
}

func validateCreate(company model.Company) *ValidationError {
	val := NewValidator()
	val.check(company.Name != "", "name", "can't be empty")
	val.check(company.EmployeeAmount > 0, "employeeamount", "number of employees must be greater than 0")
	val.check(len(company.Name) > 15, "name", "must be longer than 15 characters")
	val.check(len(company.Description) < 3000, "description", "description can't be longer than 3000 characters")
	if val.hasErrors() {
		return &ValidationError{Errors: val.getErrors()}
	}
	return nil
}

func validatePatchData(patchData map[string]interface{}) *ValidationError {
	val := NewValidator()

	if name, ok := patchData["name"].(string); ok {
		val.check(name != "", "name", "can't be empty")
		val.check(len(name) > 15, "name", "must be longer than 15 characters")
	}
	if description, ok := patchData["description"].(string); ok {
		val.check(len(description) < 3000, "description", "description can't be longer than 3000 characters")
	}

	if employeeAmountFloat, ok := patchData["employeeamount"].(float64); ok {
		employeeAmount := int16(employeeAmountFloat)
		val.check(employeeAmount > 0, "employeeamount", "number of employees must be greater than 0")
	}

	if val.hasErrors() {
		return &ValidationError{Errors: val.getErrors()}
	}
	return nil
}
