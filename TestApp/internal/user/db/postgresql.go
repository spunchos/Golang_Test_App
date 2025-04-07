package db

import (
	"TestApp/internal/user"
	"TestApp/pkg/client/postgresql"
	"TestApp/pkg/logging"
	"context"
	"errors"
	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, user *user.User) (string, error) {

	//TODO implement passwordHashing somewhere
	q := `INSERT INTO public.user (name, email) 
			VALUES ($1, $2) 
			RETURNING id`
	r.logger.Tracef("query: %s, args: %v", q, user)
	err := r.client.QueryRow(ctx, q, user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Errorf("failed to create user. message: %s, detail: %s, where: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where)
		}
		return "", err
	}
	r.logger.Tracef("inserted user with id %s", user.ID)
	return user.ID, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `SELECT id, name, email FROM public.user WHERE id = $1`
	r.logger.Tracef("query: %s, args: %v", q, id)
	var u user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		return user.User{}, err
	}

	if u.ID == "" {
		return user.User{}, errors.New("user not found")
	}

	r.logger.Tracef("found user with id %s", u.ID)
	return u, nil
}

func (r *repository) FindAll(ctx context.Context) (u []user.User, err error) {
	q := `SELECT id, name, email FROM public.user`
	r.logger.Tracef("query: %s", q)
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	users := make([]user.User, 0)

	for rows.Next() {
		var u user.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Update(ctx context.Context, user user.User) error {
	q := `UPDATE public.user SET name = $1, email = $2 WHERE id = $3`
	r.logger.Tracef("query: %s, args: %v", q, user)
	_, err := r.client.Exec(ctx, q, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}

	r.logger.Tracef("updated user with id %s", user.ID)
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM public.user WHERE id = $1`
	r.logger.Tracef("query: %s, args: %v", q, id)
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	r.logger.Tracef("deleted user with id %s", id)
	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Storage {
	return &repository{
		client: client,
		logger: logger,
	}
}
