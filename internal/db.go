package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS images (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size_bytes INTEGER NOT NULL,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		uploaded_at DATETIME NOT NULL,
		original_path TEXT NOT NULL,
		thumbnail_path TEXT NOT NULL,
		INDEX idx_user_id (user_id)
	);
	`
	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) SaveImage(ctx context.Context, img *Image) error {
	query := `
	INSERT INTO images (id, user_id, filename, content_type, size_bytes, width, height, uploaded_at, original_path, thumbnail_path)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		img.ID, img.UserID, img.Filename, img.ContentType,
		img.SizeBytes, img.Width, img.Height, img.UploadedAt,
		img.OriginalPath, img.ThumbnailPath,
	)
	return err
}

func (db *DB) GetImage(ctx context.Context, imageID string) (*Image, error) {
	query := `
	SELECT id, user_id, filename, content_type, size_bytes, width, height, uploaded_at, original_path, thumbnail_path
	FROM images WHERE id = ?
	`
	var img Image
	var uploadedAt string

	err := db.conn.QueryRowContext(ctx, query, imageID).Scan(
		&img.ID, &img.UserID, &img.Filename, &img.ContentType,
		&img.SizeBytes, &img.Width, &img.Height, &uploadedAt,
		&img.OriginalPath, &img.ThumbnailPath,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image not found")
	}
	if err != nil {
		return nil, err
	}

	img.UploadedAt, _ = time.Parse("2006-01-02 15:04:05", uploadedAt)
	return &img, nil
}

func (db *DB) ListImages(ctx context.Context, userID string, limit, offset int) ([]*Image, error) {
	query := `
	SELECT id, user_id, filename, content_type, size_bytes, width, height, uploaded_at, original_path, thumbnail_path
	FROM images WHERE user_id = ?
	ORDER BY uploaded_at DESC
	LIMIT ? OFFSET ?
	`
	rows, err := db.conn.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*Image
	for rows.Next() {
		var img Image
		var uploadedAt string

		err := rows.Scan(
			&img.ID, &img.UserID, &img.Filename, &img.ContentType,
			&img.SizeBytes, &img.Width, &img.Height, &uploadedAt,
			&img.OriginalPath, &img.ThumbnailPath,
		)
		if err != nil {
			return nil, err
		}

		img.UploadedAt, _ = time.Parse("2006-01-02 15:04:05", uploadedAt)
		images = append(images, &img)
	}

	return images, nil
}

func (db *DB) CountImages(ctx context.Context, userID string) (int, error) {
	var count int
	err := db.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM images WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

func (db *DB) DeleteImage(ctx context.Context, imageID string) error {
	result, err := db.conn.ExecContext(ctx, "DELETE FROM images WHERE id = ?", imageID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
