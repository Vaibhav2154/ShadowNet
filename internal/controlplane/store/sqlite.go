package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/model"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository implements PeerRepository using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &SQLiteRepository{db: db}

	// Initialize schema
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the peers table if it doesn't exist
func (r *SQLiteRepository) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS peers (
		id TEXT PRIMARY KEY,
		wg_public_key TEXT NOT NULL,
		endpoint_ip TEXT NOT NULL,
		endpoint_port INTEGER NOT NULL,
		last_seen DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_last_seen ON peers(last_seen);
	`

	_, err := r.db.Exec(query)
	return err
}

// CreateOrUpdate creates a new peer or updates existing one
func (r *SQLiteRepository) CreateOrUpdate(peer *model.Peer) error {
	query := `
	INSERT INTO peers (id, wg_public_key, endpoint_ip, endpoint_port, last_seen)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		wg_public_key = excluded.wg_public_key,
		endpoint_ip = excluded.endpoint_ip,
		endpoint_port = excluded.endpoint_port,
		last_seen = excluded.last_seen
	`

	_, err := r.db.Exec(query,
		peer.ID,
		peer.WGPublicKey,
		peer.EndpointIP,
		peer.EndpointPort,
		peer.LastSeen,
	)

	if err != nil {
		return fmt.Errorf("failed to create/update peer: %w", err)
	}

	return nil
}

// GetByID retrieves a peer by ID
func (r *SQLiteRepository) GetByID(id string) (*model.Peer, error) {
	query := `
	SELECT id, wg_public_key, endpoint_ip, endpoint_port, last_seen
	FROM peers
	WHERE id = ?
	`

	var peer model.Peer
	err := r.db.QueryRow(query, id).Scan(
		&peer.ID,
		&peer.WGPublicKey,
		&peer.EndpointIP,
		&peer.EndpointPort,
		&peer.LastSeen,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get peer: %w", err)
	}

	return &peer, nil
}

// GetAllActive retrieves all peers active within the timeout duration
func (r *SQLiteRepository) GetAllActive(timeout time.Duration) ([]*model.Peer, error) {
	cutoff := time.Now().Add(-timeout)

	query := `
	SELECT id, wg_public_key, endpoint_ip, endpoint_port, last_seen
	FROM peers
	WHERE last_seen > ?
	ORDER BY last_seen DESC
	`

	rows, err := r.db.Query(query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query active peers: %w", err)
	}
	defer rows.Close()

	var peers []*model.Peer
	for rows.Next() {
		var peer model.Peer
		err := rows.Scan(
			&peer.ID,
			&peer.WGPublicKey,
			&peer.EndpointIP,
			&peer.EndpointPort,
			&peer.LastSeen,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan peer: %w", err)
		}
		peers = append(peers, &peer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating peers: %w", err)
	}

	return peers, nil
}

// UpdateLastSeen updates the last seen timestamp for a peer
func (r *SQLiteRepository) UpdateLastSeen(id string) error {
	query := `
	UPDATE peers
	SET last_seen = ?
	WHERE id = ?
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("peer not found: %s", id)
	}

	return nil
}

// Delete removes a peer from storage
func (r *SQLiteRepository) Delete(id string) error {
	query := `DELETE FROM peers WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete peer: %w", err)
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
