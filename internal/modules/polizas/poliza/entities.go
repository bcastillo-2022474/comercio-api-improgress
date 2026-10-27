// Package poliza is the import declaration aggregate: a declaration header and the goods it brings in.
package poliza

import (
	"time"
)

const dateLayout = "2006-01-02"

// Item is one good in a declaration. Its SAC, country and unit are references to lookup tables.
type Item struct {
	SacId                   int64
	PaisId                  int64
	TipoUnidadMedidaId      int64
	CantidadFraccion        float64
	TasaDai                 float64
	ValorDai                float64
	ValorCifUds             float64
	TasaCifCantidadFraccion float64
}

func (i Item) valid() bool {
	return i.SacId > 0 && i.PaisId > 0 && i.TipoUnidadMedidaId > 0 &&
		i.CantidadFraccion >= 0 && i.TasaDai >= 0 && i.ValorDai >= 0 &&
		i.ValorCifUds >= 0 && i.TasaCifCantidadFraccion >= 0
}

type Poliza struct {
	Id               int64
	Correlativo      int64
	FechaDeclaracion string // YYYY-MM-DD
	TipoCambioDolar  float64
	// Customs office and regime are optional: the source data leaves them empty on some declarations.
	AduanaId      *int64
	TipoRegimenId *int64
	Items         []Item
}

// ValidFecha checks that fecha is a YYYY-MM-DD date.
func ValidFecha(fecha string) error {
	if _, err := time.Parse(dateLayout, fecha); err != nil {
		return &InvalidFechaError{Fecha: fecha}
	}
	return nil
}

// NewPoliza checks a new declaration. It has no Id until the repository stores it.
func NewPoliza(correlativo int64, fecha string, tipoCambio float64, aduanaId, tipoRegimenId *int64, items []Item) (Poliza, error) {
	if correlativo <= 0 {
		return Poliza{}, &InvalidCorrelativoError{Correlativo: correlativo}
	}
	if err := ValidFecha(fecha); err != nil {
		return Poliza{}, err
	}
	if tipoCambio <= 0 {
		return Poliza{}, &InvalidTipoCambioError{TipoCambio: tipoCambio}
	}
	if len(items) == 0 {
		return Poliza{}, &NoItemsError{}
	}
	for i, item := range items {
		if !item.valid() {
			return Poliza{}, &InvalidItemError{Index: i}
		}
	}

	return Poliza{
		Correlativo:      correlativo,
		FechaDeclaracion: fecha,
		TipoCambioDolar:  tipoCambio,
		AduanaId:         aduanaId,
		TipoRegimenId:    tipoRegimenId,
		Items:            items,
	}, nil
}
