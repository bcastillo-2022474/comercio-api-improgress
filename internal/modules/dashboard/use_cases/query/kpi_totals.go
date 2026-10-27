package query

import (
	"context"
	"log/slog"
	"time"

	"joao-api-comercio/internal/modules/dashboard/readmodel"
	"joao-api-comercio/internal/shared/cache"
)

type KpiTotalsUseCase struct {
	totals readmodel.TotalsQuery
	log    *slog.Logger
	cache  *cache.TTL[struct{}, readmodel.KpiTotals]
}

func NewKpiTotalsUseCase(totals readmodel.TotalsQuery, log *slog.Logger, ttl time.Duration) *KpiTotalsUseCase {
	return &KpiTotalsUseCase{totals: totals, log: log.With("use_case", "kpi_totals"), cache: cache.NewTTL[struct{}, readmodel.KpiTotals](ttl, nil)}
}

func (uc *KpiTotalsUseCase) Execute(ctx context.Context) (readmodel.KpiTotals, error) {
	uc.log.InfoContext(ctx, "starting kpi totals")

	if cached, ok := uc.cache.Get(struct{}{}); ok {
		uc.log.InfoContext(ctx, "kpi served from cache")
		return cached, nil
	}

	totals, err := uc.totals.KpiTotals(ctx)
	if err != nil {
		return readmodel.KpiTotals{}, err
	}

	uc.cache.Set(struct{}{}, totals)
	uc.log.InfoContext(ctx, "kpi totals finished", "polizas", totals.PolizasCount)
	return totals, nil
}
