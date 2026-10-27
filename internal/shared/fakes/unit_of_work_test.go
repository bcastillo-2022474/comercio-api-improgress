package fakes_test

import (
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	userfakes "joao-api-comercio/internal/modules/identity/user/fakes"
	"joao-api-comercio/internal/shared"
	"joao-api-comercio/internal/shared/contracts"
	sharedfakes "joao-api-comercio/internal/shared/fakes"
	outboxfakes "joao-api-comercio/internal/shared/outbox/fakes"
)

func TestFakeUnitOfWorkContract(t *testing.T) {
	contracts.RunUnitOfWork(t, func(*testing.T) contracts.Harness {
		users := userfakes.NewRepository()
		outbox := outboxfakes.NewRepository()
		return contracts.Harness{
			UoW:    sharedfakes.NewUnitOfWork[user.Repository](users, outbox, time.Now),
			Users:  users,
			Outbox: outbox,
		}
	})
}

var _ shared.UnitOfWork[user.Repository] = (*sharedfakes.UnitOfWork[user.Repository])(nil)
