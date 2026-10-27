-- name: ListPaises :many
SELECT id, name
FROM pais
ORDER BY name;

-- name: ListAduanas :many
SELECT id, name
FROM aduanas
ORDER BY name;

-- name: ListTiposRegimen :many
SELECT id, name
FROM tipo_regimen
ORDER BY name;
