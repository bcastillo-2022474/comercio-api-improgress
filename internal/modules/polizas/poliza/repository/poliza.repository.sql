-- name: InsertPoliza :execlastid
INSERT INTO polizas (correlativo, fecha_declaracion, tipo_cambio_dolar, aduana_id, tipo_regimen_id)
VALUES (@correlativo, CAST(@fecha_declaracion AS TEXT), @tipo_cambio_dolar, @aduana_id, @tipo_regimen_id);

-- name: InsertPolizaItem :exec
INSERT INTO poliza_items (poliza_id, sac_id, pais_id, tipo_unidad_medida_id, cantidad_fraccion, tasa_dai, valor_dai, valor_cif_uds, tasa_cif_cantidad_fraccion)
VALUES (@poliza_id, @sac_id, @pais_id, @tipo_unidad_medida_id, @cantidad_fraccion, @tasa_dai, @valor_dai, @valor_cif_uds, @tasa_cif_cantidad_fraccion);

-- name: GetPolizaById :one
-- The driver turns a DATE column into time.Time on read, so the date is cast to TEXT to stay YYYY-MM-DD.
SELECT id, correlativo, CAST(fecha_declaracion AS TEXT) AS fecha_declaracion, tipo_cambio_dolar, aduana_id, tipo_regimen_id
FROM polizas
WHERE id = @id;

-- name: ListItemsByPolizaId :many
SELECT sac_id, pais_id, tipo_unidad_medida_id, cantidad_fraccion, tasa_dai, valor_dai, valor_cif_uds, tasa_cif_cantidad_fraccion
FROM poliza_items
WHERE poliza_id = @poliza_id
ORDER BY id;
