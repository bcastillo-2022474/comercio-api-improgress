-- name: KpiTotals :one
SELECT CAST(COALESCE(SUM(i.valor_cif_uds), 0) AS REAL) AS total_valor_cif_uds,
       (SELECT COUNT(*) FROM polizas) AS polizas_count,
       COUNT(*) AS items_count
FROM poliza_items i;

-- name: TopCountries :many
SELECT i.pais_id, p.name AS pais_name, CAST(SUM(i.valor_cif_uds) AS REAL) AS valor_cif_uds
FROM poliza_items i
JOIN pais p ON p.id = i.pais_id
GROUP BY i.pais_id, p.name
ORDER BY valor_cif_uds DESC
LIMIT 5;

-- name: TopProducts :many
SELECT i.sac_id, s.fraccion, s.descripcion, CAST(SUM(i.valor_cif_uds) AS REAL) AS valor_cif_uds
FROM poliza_items i
JOIN sac s ON s.id = i.sac_id
GROUP BY i.sac_id, s.fraccion, s.descripcion
ORDER BY valor_cif_uds DESC
LIMIT 5;

-- name: ValorPorMes :many
SELECT CAST(substr(p.fecha_declaracion, 6, 2) AS INTEGER) AS mes,
       CAST(SUM(i.valor_cif_uds) AS REAL) AS valor_cif_uds
FROM poliza_items i
JOIN polizas p ON p.id = i.poliza_id
WHERE CAST(substr(p.fecha_declaracion, 1, 4) AS INTEGER) = CAST(@year AS INTEGER)
GROUP BY mes
ORDER BY mes;

-- name: LatestYear :one
SELECT CAST(COALESCE(MAX(substr(fecha_declaracion, 1, 4)), '0') AS INTEGER) AS year
FROM polizas;
