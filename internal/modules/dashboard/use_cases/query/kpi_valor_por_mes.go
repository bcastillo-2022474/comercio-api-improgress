package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dashboard/kpi"
	"joao-api-comercio/internal/modules/dashboard/readmodel"
)

// ValorPorMes is one year of import value, with all twelve months present. Months with no
// declarations are zero, so a chart always gets twelve points.
type ValorPorMes struct {
	Year  int64                  `json:"year"`
	Meses []readmodel.MonthValue `json:"meses"`
}

type KpiValorPorMesUseCase struct {
	valorPorMes readmodel.ValorPorMesQuery
	log         *slog.Logger
}

func NewKpiValorPorMesUseCase(valorPorMes readmodel.ValorPorMesQuery, log *slog.Logger) *KpiValorPorMesUseCase {
	return &KpiValorPorMesUseCase{valorPorMes: valorPorMes, log: log.With("use_case", "kpi_valor_por_mes")}
}

// Execute reads one year. Year 0 means the latest year in the data.
func (uc *KpiValorPorMesUseCase) Execute(ctx context.Context, year int64) (ValorPorMes, error) {
	uc.log.InfoContext(ctx, "starting kpi valor por mes", "year", year)

	if year == 0 {
		latest, err := uc.valorPorMes.LatestYear(ctx)
		if err != nil {
			return ValorPorMes{}, err
		}
		year = latest
	}
	if year < 1900 || year > 2100 {
		return ValorPorMes{}, &kpi.InvalidYearError{Year: year}
	}

	rows, err := uc.valorPorMes.ValorPorMes(ctx, year)
	if err != nil {
		return ValorPorMes{}, err
	}

	meses := make([]readmodel.MonthValue, 12)
	for i := range meses {
		meses[i] = readmodel.MonthValue{Mes: int64(i + 1)}
	}
	for _, row := range rows {
		if row.Mes >= 1 && row.Mes <= 12 {
			meses[row.Mes-1].ValorCifUsd = row.ValorCifUsd
		}
	}

	uc.log.InfoContext(ctx, "kpi valor por mes finished", "year", year)
	return ValorPorMes{Year: year, Meses: meses}, nil
}
