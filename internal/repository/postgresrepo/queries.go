package postgresrepo

const (
	// User management queries.
	queryAddUser = "INSERT INTO users (login, pass_hash) VALUES ($1, $2) RETURNING id"
	queryGetUser = "SELECT id, login, pass_hash FROM users WHERE login = $1"

	// Metadata management queries.
	queryUploadMetadata = `INSERT INTO metadata 
(name, is_file, public, mime, owner_id, json_data, file_size) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	queryGetUsersMetadata = `SELECT 
    m.id,
    m.name,
    m.mime,
    m.is_file,
    m.public,
    m.created,
    m.owner_id,
    m.json_data,
    m.file_size,
    COALESCE(string_agg(u.username, ','), '') AS grant
FROM 
    metadata m
LEFT JOIN 
    meta_access ma ON m.id = ma.meta_id
LEFT JOIN 
    users u ON ma.user_id = u.id
WHERE 
    m.id = $1 AND m.deleted = false
GROUP BY 
    m.id, m.name, m.mime, m.is_file, m.public, m.created
ORDER BY 
    m.name ASC, 
    m.created DESC;
`
	queryDeleteMetadata = `UPDATE metadata SET deleted = true WHERE id = $1 AND owner_id = $2`

	// Metadata Access queries.
	queryGrantMetadataAcsess = `INSERT INTO meta_access (meta_id, user_id)
VALUES (
    $1,
    (SELECT id FROM users WHERE username = $2)
)`
)
