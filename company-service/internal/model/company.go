package model

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id             primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	EmployeeAmount int16              `json:"employeeamount"`
	Registered     bool               `json:"registered"`
	Type           CompanyType        `json:"type"`
}

type CompanyType int8

const (
	Corporations CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship
)

var companyTypeToString = map[CompanyType]string{
	Corporations:       "Corporations",
	NonProfit:          "NonProfit",
	Cooperative:        "Cooperative",
	SoleProprietorship: "SoleProprietorship",
}

var stringToCompanyType = map[string]CompanyType{
	"Corporations":       Corporations,
	"NonProfit":          NonProfit,
	"Cooperative":        Cooperative,
	"SoleProprietorship": SoleProprietorship,
}

func (c CompanyType) MarshalJSON() ([]byte, error) {
	str, ok := companyTypeToString[c]
	if !ok {
		return nil, errors.New("invalid type value")
	}
	return json.Marshal(str)
}

func (c *CompanyType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	val, ok := stringToCompanyType[str]
	if !ok {
		return errors.New("invalid type value, only following are available: Corporations, NonProfit, Cooperative, Sole Proprietorship")
	}
	*c = val
	return nil
}
