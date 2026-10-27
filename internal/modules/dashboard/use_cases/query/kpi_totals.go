package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dashboard/readmodel"
)

type KpiTotalsUseCase struct {
	totals readmodel.TotalsQuery
	log    *slog.Logger
}

func NewKpiTotalsUseCase(totals readmodel.TotalsQuery, log *slog.Logger) *KpiTotalsUseCase {
	return &KpiTotalsUseCase{totals: totals, log: log.With("use_case", "kpi_totals")}
}

func (uc *KpiTotalsUseCase) Execute(ctx context.Context) (readmodel.KpiTotals, error) {
	uc.log.InfoContext(ctx, "starting kpi totals")

	totals, err := uc.totals.KpiTotals(ctx)
	if err != nil {
		return readmodel.KpiTotals{}, err
	}

	uc.log.InfoContext(ctx, "kpi totals finished", "polizas", totals.PolizasCount)
	return totals, nil
}
