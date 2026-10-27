package command

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/polizas/poliza"
	polizafakes "joao-api-comercio/internal/modules/polizas/poliza/fakes"
	sharedfakes "joao-api-comercio/internal/shared/fakes"
	outboxfakes "joao-api-comercio/internal/shared/outbox/fakes"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func newCreate(repo *polizafakes.Repository) (*CreatePolizaUseCase, *outboxfakes.Repository) {
	outbox := outboxfakes.NewRepository()
	uow := sharedfakes.NewUnitOfWork[poliza.Repository](repo, outbox, time.Now)
	return NewCreatePolizaUseCase(uow, discard), outbox
}

func validInput(correlativo int64) CreatePolizaInput {
	return CreatePolizaInput{
		Correlativo:      correlativo,
		FechaDeclaracion: "2025-01-15",
		TipoCambioDolar:  7.8,
		Items:            []poliza.Item{{SacId: 1, PaisId: 1, TipoUnidadMedidaId: 1, CantidadFraccion: 1, ValorCifUds: 10}},
	}
}

func TestCreatePoliza_StoresAndPublishesEvent(t *testing.T) {
	create, outbox := newCreate(polizafakes.NewRepository())

	p, err := create.Execute(context.Background(), validInput(1001))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Id <= 0 {
		t.Errorf("id = %d, want a positive id", p.Id)
	}

	records := outbox.Records()
	if len(records) != 1 || records[0].EventName != "poliza.created" {
		t.Fatalf("outbox = %v, want one poliza.created record", records)
	}
}

func TestCreatePoliza_RejectionsWriteNothing(t *testing.T) {
	repo := polizafakes.NewRepository()
	create, outbox := newCreate(repo)
	if _, err := create.Execute(context.Background(), validInput(1001)); err != nil {
		t.Fatalf("seeding: %v", err)
	}

	noItems := validInput(1002)
	noItems.Items = nil
	tests := []struct {
		name   string
		input  CreatePolizaInput
		target any
	}{
		{name: "no items", input: noItems, target: new(*poliza.NoItemsError)},
		{name: "correlativo taken", input: validInput(1001), target: new(*poliza.CorrelativoTakenError)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := len(outbox.Records())
			_, err := create.Execute(context.Background(), tt.input)
			if !errors.As(err, tt.target) {
				t.Fatalf("err = %v, want %T", err, tt.target)
			}
			if len(outbox.Records()) != before {
				t.Error("a rejected declaration wrote an outbox record")
			}
		})
	}
}
