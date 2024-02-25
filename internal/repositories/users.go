package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wildanfaz/store-management/internal/models"
)

type ImplementUsers struct {
	dbPostgreSQL *pgx.Conn
}

type Users interface {
	Register(ctx context.Context, payload models.User) error
	Profile(ctx context.Context, email string) (*models.User, error)
	ResetPassword(ctx context.Context, email string, password string) error
	UpdateIsLogin(ctx context.Context, email string, isLogin bool) error
	IsLogin(ctx context.Context, email string) (bool, error)
}

func NewUsersRepository(dbPostgreSQL *pgx.Conn) Users {
	return &ImplementUsers{
		dbPostgreSQL: dbPostgreSQL,
	}
}

func (r *ImplementUsers) Register(ctx context.Context, payload models.User) error {
	var query = `
	INSERT INTO users
	(full_name, email, password)
	VALUES
	($1, $2, $3)
	`

	_, err := r.dbPostgreSQL.Exec(ctx, query,
		payload.FullName, payload.Email, payload.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementUsers) Profile(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	var query = `
	SELECT id, full_name, email, password, is_login, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := r.dbPostgreSQL.QueryRow(ctx, query, email).Scan(
		&user.Id, &user.FullName, &user.Email,
		&user.Password, &user.IsLogin, &user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ImplementUsers) ResetPassword(ctx context.Context, email string, password string) error {
	var query = `
	UPDATE users
	SET password = $1
	WHERE email = $2
	`

	_, err := r.dbPostgreSQL.Exec(ctx, query, password, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementUsers) UpdateIsLogin(ctx context.Context, email string, isLogin bool) error {
	var query = `
	UPDATE users
	SET is_login = $1
	WHERE email = $2
	`

	_, err := r.dbPostgreSQL.Exec(ctx, query, isLogin, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementUsers) IsLogin(ctx context.Context, email string) (bool, error) {
	var isLogin bool

	var query = `
	SELECT is_login
	FROM users
	WHERE email = $1
	`

	err := r.dbPostgreSQL.QueryRow(ctx, query, email).Scan(&isLogin)
	if err == pgx.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return isLogin, nil
}
