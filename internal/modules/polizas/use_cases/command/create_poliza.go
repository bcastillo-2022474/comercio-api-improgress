package command

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/shared"
)

type CreatePolizaInput struct {
	Correlativo      int64
	FechaDeclaracion string
	TipoCambioDolar  float64
	AduanaId         *int64
	TipoRegimenId    *int64
	Items            []poliza.Item
}

type CreatePolizaUseCase struct {
	uow shared.UnitOfWork[poliza.Repository]
	log *slog.Logger
}

func NewCreatePolizaUseCase(uow shared.UnitOfWork[poliza.Repository], log *slog.Logger) *CreatePolizaUseCase {
	return &CreatePolizaUseCase{uow: uow, log: log.With("use_case", "create_poliza")}
}

// Execute stores the declaration with its items and publishes poliza.created once the write commits.
func (uc *CreatePolizaUseCase) Execute(ctx context.Context, in CreatePolizaInput) (poliza.Poliza, error) {
	uc.log.InfoContext(ctx, "starting create poliza", "correlativo", in.Correlativo, "items", len(in.Items))

	p, err := poliza.NewPoliza(in.Correlativo, in.FechaDeclaracion, in.TipoCambioDolar, in.AduanaId, in.TipoRegimenId, in.Items)
	if err != nil {
		return poliza.Poliza{}, err
	}

	err = uc.uow.Execute(ctx, func(ctx context.Context, polizas poliza.Repository, pub shared.Publisher) error {
		if err := polizas.Insert(ctx, &p); err != nil {
			return err
		}
		pub.Publish(poliza.NewPolizaCreatedEvent(p))
		return nil
	})
	if err != nil {
		return poliza.Poliza{}, err
	}

	uc.log.InfoContext(ctx, "poliza created", "poliza_id", p.Id)
	return p, nil
}
