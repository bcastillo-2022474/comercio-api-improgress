-- +goose Up
-- Lookup tables and import tables. The provided comercio_exterior.db already has them, filled by
-- setup.py. IF NOT EXISTS makes this a no-op there, and creates them on a fresh database.
CREATE TABLE IF NOT EXISTS aduanas (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tipo_regimen (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pais (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS sac (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    fraccion    BIGINT NOT NULL UNIQUE,
    descripcion TEXT
);

CREATE TABLE IF NOT EXISTS tipo_unidad_medida (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS polizas (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    correlativo       BIGINT NOT NULL UNIQUE,
    fecha_declaracion DATE,
    tipo_cambio_dolar FLOAT,
    aduana_id         INTEGER,
    tipo_regimen_id   INTEGER,
    FOREIGN KEY (aduana_id) REFERENCES aduanas (id),
    FOREIGN KEY (tipo_regimen_id) REFERENCES tipo_regimen (id)
);

CREATE TABLE IF NOT EXISTS poliza_items (
    id                         INTEGER PRIMARY KEY AUTOINCREMENT,
    poliza_id                  INTEGER,
    sac_id                     INTEGER,
    pais_id                    INTEGER,
    tipo_unidad_medida_id      INTEGER,
    cantidad_fraccion          FLOAT,
    tasa_dai                   FLOAT,
    valor_dai                  FLOAT,
    valor_cif_uds              FLOAT,
    tasa_cif_cantidad_fraccion FLOAT,
    FOREIGN KEY (poliza_id) REFERENCES polizas (id) ON DELETE CASCADE,
    FOREIGN KEY (sac_id) REFERENCES sac (id),
    FOREIGN KEY (pais_id) REFERENCES pais (id),
    FOREIGN KEY (tipo_unidad_medida_id) REFERENCES tipo_unidad_medida (id)
);

CREATE INDEX IF NOT EXISTS idx_poliza_items_poliza_id ON poliza_items (poliza_id);
CREATE INDEX IF NOT EXISTS idx_poliza_items_sac_id ON poliza_items (sac_id);
CREATE INDEX IF NOT EXISTS idx_poliza_items_pais_id ON poliza_items (pais_id);

-- +goose Down
-- Intentionally empty. These tables hold import data we did not create, so a rollback must not drop them.
