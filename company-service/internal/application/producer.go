package application

import "github.com/milicns/company-manager/company-service/internal/utils"

type Producer interface {
	Produce(utils.Event)
}
