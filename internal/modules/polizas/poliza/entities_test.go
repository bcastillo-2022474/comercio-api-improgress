package poliza_test

import (
	"errors"
	"testing"

	"joao-api-comercio/internal/modules/polizas/poliza"
)

func validItem() poliza.Item {
	return poliza.Item{SacId: 1, PaisId: 1, TipoUnidadMedidaId: 1, CantidadFraccion: 1, ValorCifUds: 10}
}

func TestNewPoliza_Rejections(t *testing.T) {
	tests := []struct {
		name        string
		correlativo int64
		fecha       string
		tipoCambio  float64
		items       []poliza.Item
		target      any
	}{
		{name: "zero correlativo", correlativo: 0, fecha: "2025-01-15", tipoCambio: 7.8, items: []poliza.Item{validItem()}, target: new(*poliza.InvalidCorrelativoError)},
		{name: "bad date", correlativo: 1, fecha: "15/01/2025", tipoCambio: 7.8, items: []poliza.Item{validItem()}, target: new(*poliza.InvalidFechaError)},
		{name: "zero tipo de cambio", correlativo: 1, fecha: "2025-01-15", tipoCambio: 0, items: []poliza.Item{validItem()}, target: new(*poliza.InvalidTipoCambioError)},
		{name: "no items", correlativo: 1, fecha: "2025-01-15", tipoCambio: 7.8, items: nil, target: new(*poliza.NoItemsError)},
		{name: "item without sac", correlativo: 1, fecha: "2025-01-15", tipoCambio: 7.8, items: []poliza.Item{{PaisId: 1, TipoUnidadMedidaId: 1}}, target: new(*poliza.InvalidItemError)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := poliza.NewPoliza(tt.correlativo, tt.fecha, tt.tipoCambio, nil, nil, tt.items)
			if !errors.As(err, tt.target) {
				t.Fatalf("err = %v, want %T", err, tt.target)
			}
		})
	}
}

func TestNewPoliza_AcceptsAValidDeclaration(t *testing.T) {
	p, err := poliza.NewPoliza(1001, "2025-01-15", 7.8, nil, nil, []poliza.Item{validItem()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Id != 0 || p.Correlativo != 1001 || len(p.Items) != 1 {
		t.Errorf("got %+v", p)
	}
}
