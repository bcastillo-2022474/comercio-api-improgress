-- The driver reads a DATE column as time.Time, so dates are cast to TEXT and stay YYYY-MM-DD.

-- name: ListPolizas :many
SELECT p.id, p.correlativo, CAST(p.fecha_declaracion AS TEXT) AS fecha_declaracion, p.tipo_cambio_dolar,
       p.aduana_id, a.name AS aduana_name, p.tipo_regimen_id, r.name AS tipo_regimen_name,
       (SELECT COUNT(*) FROM poliza_items i WHERE i.poliza_id = p.id) AS items_count
FROM polizas p
LEFT JOIN aduanas a ON a.id = p.aduana_id
LEFT JOIN tipo_regimen r ON r.id = p.tipo_regimen_id
WHERE (sqlc.narg('aduana_id') IS NULL OR p.aduana_id = sqlc.narg('aduana_id'))
  AND (sqlc.narg('tipo_regimen_id') IS NULL OR p.tipo_regimen_id = sqlc.narg('tipo_regimen_id'))
  AND (sqlc.narg('fecha_desde') IS NULL OR p.fecha_declaracion >= sqlc.narg('fecha_desde'))
  AND (sqlc.narg('fecha_hasta') IS NULL OR p.fecha_declaracion <= sqlc.narg('fecha_hasta'))
ORDER BY p.fecha_declaracion DESC, p.id DESC
LIMIT @page_size OFFSET @page_offset;

-- name: CountPolizas :one
SELECT COUNT(*)
FROM polizas p
WHERE (sqlc.narg('aduana_id') IS NULL OR p.aduana_id = sqlc.narg('aduana_id'))
  AND (sqlc.narg('tipo_regimen_id') IS NULL OR p.tipo_regimen_id = sqlc.narg('tipo_regimen_id'))
  AND (sqlc.narg('fecha_desde') IS NULL OR p.fecha_declaracion >= sqlc.narg('fecha_desde'))
  AND (sqlc.narg('fecha_hasta') IS NULL OR p.fecha_declaracion <= sqlc.narg('fecha_hasta'));

-- name: GetPoliza :one
SELECT p.id, p.correlativo, CAST(p.fecha_declaracion AS TEXT) AS fecha_declaracion, p.tipo_cambio_dolar,
       p.aduana_id, a.name AS aduana_name, p.tipo_regimen_id, r.name AS tipo_regimen_name,
       (SELECT COUNT(*) FROM poliza_items i WHERE i.poliza_id = p.id) AS items_count
FROM polizas p
LEFT JOIN aduanas a ON a.id = p.aduana_id
LEFT JOIN tipo_regimen r ON r.id = p.tipo_regimen_id
WHERE p.id = @id;

-- name: ListPolizaItems :many
SELECT i.id, i.sac_id, s.fraccion AS sac_fraccion, s.descripcion AS sac_descripcion,
       i.pais_id, pa.name AS pais_name, i.tipo_unidad_medida_id, u.name AS tipo_unidad_medida_name,
       i.cantidad_fraccion, i.tasa_dai, i.valor_dai, i.valor_cif_uds, i.tasa_cif_cantidad_fraccion
FROM poliza_items i
JOIN sac s ON s.id = i.sac_id
JOIN pais pa ON pa.id = i.pais_id
JOIN tipo_unidad_medida u ON u.id = i.tipo_unidad_medida_id
WHERE i.poliza_id = @poliza_id
ORDER BY i.id;
