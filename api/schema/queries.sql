-- name: CreateURL :execresult
INSERT INTO urls (short_code, original_url) VALUES (?, ?);

-- name: GetURLByCode :one
SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = ? LIMIT 1;

-- name: RecordClick :exec
INSERT INTO clicks (url_id, ip_address, user_agent) VALUES (?, ?, ?);

-- name: GetAnalyticsByCode :many
SELECT c.clicked_at, c.ip_address, c.user_agent
FROM clicks c
JOIN urls u ON c.url_id = u.id
WHERE u.short_code = ?
ORDER BY c.clicked_at DESC;