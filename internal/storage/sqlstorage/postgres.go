package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
		
	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email, username string, hashedPass []byte) error {
	const op = "storage.postgres.SaveUser"

	tx, _ := s.db.Begin()
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT 
		INTO users (email, username, slug, hashed_password, last_password_change, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now().UTC()
	slug := strings.ToLower(username)

	_, err = stmt.ExecContext(
		ctx,
		email,
		username,
		slug,
		hashedPass,
		now,
		now,
		now,
	) 
	if err != nil {
		// TODO: handle error with sql errors
		// var sqlErr sql.ErrNoRows
		// if errors.Is(err, &sqlErr) && sqlErr.ExtendedCode == sql.ErrConstraintUnique{
		// 	return fmt.Errorf("%s: %w", op, storage.ErrUserExist)
		// }
		return fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return nil
}

func (s *Storage) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	const op = "storage.postgres.GetUserByID"
	var user models.User

	tx, _ := s.db.Begin()
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`SELECT id, email username, slug, hashed_password, last_password_change, created_at, updated_at  
		FROM users
		WHERE id = $1`,
	)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID) 

	err = row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Slug, &user.HashedPass,
		&user.LastPasswordChange, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows{
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return user, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.GetUserByEmail"
	var user models.User

	tx, _ := s.db.Begin()
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`SELECT id, email, username, slug, hashed_password, last_password_change, created_at, updated_at  
		FROM users
		WHERE email like $1`,
	)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email) 

	err = row.Scan(
		&user.ID, &user.Email, &user.Username, &user.Slug, &user.HashedPass,
		&user.LastPasswordChange, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows{
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return user, nil
}

func (s *Storage) PatchUsername(ctx context.Context, user models.User, username string) error {
	const op = "storage.postgres.PatchUsername"
	// TODO: handle

	tx, _ := s.db.Begin()
	defer tx.Rollback()
		
	stmt, err := tx.Prepare(
		`UPDATE users 
		SET 
			username = $1,
			slug = $2,
			updated_at = $3
		WHERE id = $4`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	slug := strings.ToLower(username)
	now := time.Now().UTC()

	_, err = stmt.ExecContext(ctx, username, slug, now, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return nil
}

func (s *Storage) PatchPassword(ctx context.Context, user models.User, password []byte) error {
	const op = "storage.postgres.PatchPassword"
	// TODO: handle

	tx, _ := s.db.Begin()
	defer tx.Rollback()
	
	stmt, err := tx.Prepare(
		`UPDATE users 
		SET 
			hashed_password = $1,
			updated_at = $2,
			last_password_change = $3
		WHERE id = $4`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now().UTC()

	_, err = stmt.ExecContext(ctx, password, now, now, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	tx, _ := s.db.Begin()
	defer tx.Rollback()

	var app models.App

	stmt, err := tx.Prepare(
		`SELECT id, name, secret
		FROM apps
		WHERE id = $1`,
	)
	if err != nil {
		return app, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID) 

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		return app, fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	return app, nil
}
