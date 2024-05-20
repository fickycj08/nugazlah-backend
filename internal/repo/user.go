package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/nugazlah-backend/internal/entity"
	"github.com/vandenbill/nugazlah-backend/internal/ierr"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) error {
	q := `INSERT INTO users (fullname, email, password)
	VALUES ($1, $2, $3)`

	_, err := u.conn.Exec(ctx, q, user.FullName, user.Email, user.Password)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	q := `SELECT id, fullname, email, password FROM users WHERE email = $1`

	var user entity.User
	err := u.conn.QueryRow(ctx, q, email).Scan(&user.ID, &user.FullName, &user.Email, &user.Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return entity.User{}, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return entity.User{}, ierr.ErrNotFound
			}
		}
		return entity.User{}, err
	}

	return user, nil
}
