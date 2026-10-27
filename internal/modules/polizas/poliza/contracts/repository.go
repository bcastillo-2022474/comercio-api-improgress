// Package contracts holds the behavior every poliza.Repository must have. The in-memory fake and
// the SQL adapter both run it, so the two cannot drift apart.
package contracts

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"joao-api-comercio/internal/modules/polizas/poliza"
)

// RunRepository runs the contract. newRepo gets a fresh, empty repository for each subtest. The
// contract uses aduana 1 and the lookup rows 1 for SAC, country and unit, so the repository must
// have them.
func RunRepository(t *testing.T, newRepo func(t *testing.T) poliza.Repository) {
	t.Run("insert assigns an id", func(t *testing.T) {
		repo := newRepo(t)
		p := newPoliza(t, 1001, nil)

		mustInsert(t, repo, &p)
		if p.Id <= 0 {
			t.Fatalf("id = %d, want a positive id", p.Id)
		}
	})

	t.Run("two declarations get different ids", func(t *testing.T) {
		repo := newRepo(t)
		first, second := newPoliza(t, 1001, nil), newPoliza(t, 1002, nil)

		mustInsert(t, repo, &first)
		mustInsert(t, repo, &second)
		if first.Id == second.Id {
			t.Fatalf("both declarations got id %d", first.Id)
		}
	})

	t.Run("a correlativo already in use is rejected", func(t *testing.T) {
		repo := newRepo(t)
		first := newPoliza(t, 1001, nil)
		mustInsert(t, repo, &first)

		again := newPoliza(t, 1001, nil)
		err := repo.Insert(context.Background(), &again)
		var taken *poliza.CorrelativoTakenError
		if !errors.As(err, &taken) {
			t.Fatalf("err = %v, want *CorrelativoTakenError", err)
		}
	})

	t.Run("a saved declaration comes back with its header and items", func(t *testing.T) {
		repo := newRepo(t)
		aduana := int64(1)
		p := newPoliza(t, 1001, &aduana)
		p.Items = append(p.Items, poliza.Item{SacId: 1, PaisId: 1, TipoUnidadMedidaId: 1, CantidadFraccion: 2.5, TasaDai: 0.15, ValorDai: 12.5, ValorCifUds: 99.9, TasaCifCantidadFraccion: 39.96})
		mustInsert(t, repo, &p)

		got := mustGet(t, repo, p.Id)
		requireSame(t, got, p)
	})

	t.Run("an unknown id is absent, not an error", func(t *testing.T) {
		repo := newRepo(t)

		got, err := repo.GetById(context.Background(), 999999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsPresent {
			t.Fatalf("found %+v for an unknown id", got.Value)
		}
	})
}

func newPoliza(t *testing.T, correlativo int64, aduanaId *int64) poliza.Poliza {
	t.Helper()
	items := []poliza.Item{{SacId: 1, PaisId: 1, TipoUnidadMedidaId: 1, CantidadFraccion: 1, ValorCifUds: 10}}
	p, err := poliza.NewPoliza(correlativo, "2025-01-15", 7.8, aduanaId, nil, items)
	if err != nil {
		t.Fatalf("building poliza: %v", err)
	}
	return p
}

func mustInsert(t *testing.T, repo poliza.Repository, p *poliza.Poliza) {
	t.Helper()
	if err := repo.Insert(context.Background(), p); err != nil {
		t.Fatalf("inserting: %v", err)
	}
}

func mustGet(t *testing.T, repo poliza.Repository, id int64) poliza.Poliza {
	t.Helper()
	got, err := repo.GetById(context.Background(), id)
	if err != nil {
		t.Fatalf("getting poliza %d: %v", id, err)
	}
	if !got.IsPresent {
		t.Fatalf("poliza %d not found after insert", id)
	}
	return got.Value
}

// requireSame compares what the caller stored with what the repository returned, field by field.
func requireSame(t *testing.T, got, want poliza.Poliza) {
	t.Helper()
	if got.Id != want.Id || got.Correlativo != want.Correlativo ||
		got.FechaDeclaracion != want.FechaDeclaracion || got.TipoCambioDolar != want.TipoCambioDolar ||
		!samePointer(got.AduanaId, want.AduanaId) || !samePointer(got.TipoRegimenId, want.TipoRegimenId) ||
		!reflect.DeepEqual(got.Items, want.Items) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func samePointer(a, b *int64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
