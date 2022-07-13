package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"

	"transaction/internal/adapters/db/postgresql"
	user "transaction/internal/users"
)

type Repository struct {
	client postgresql.Client
	logger *log.Logger
}

func NewRepository(client postgresql.Client, logger *log.Logger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

func (r *Repository) Create(ctx context.Context, user *user.User) error {

	query := "INSERT INTO public.users (user_name, email, balance) VALUES ($1, $2, $3) RETURNING id;"
	err := r.client.QueryRow(ctx, query, user.Name, user.Email, user.Balance).Scan(&user.Id)

	if err != nil {

		if pgError, ok := err.(*pgconn.PgError); ok {
			sqlError := fmt.Errorf(
				fmt.Sprintf(
					"PostgreSQL Error: %s, Details: %s, Where: %s, SQL state: %s",
					pgError.Message,
					pgError.Detail,
					pgError.Where,
					pgError.SQLState(),
				),
			)
			r.logger.Println(sqlError)
			return sqlError
		}
		return err
	}
	return nil
}

func (r *Repository) FindAll(ctx context.Context) (users []user.User, err error) {

	query := "SELECT id, user_name, balance FROM public.users;"
	rows, err := r.client.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	users = make([]user.User, 0)
	for rows.Next() {

		var newUser user.User

		err := rows.Scan(&newUser.Id, &newUser.Name, &newUser.Balance)
		if err != nil {
			return nil, err
		}

		users = append(users, newUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) FindOne(ctx context.Context, id string) (user.User, error) {

	var newUser user.User

	query := "SELECT id, user_name, balance FROM public.users WHERE id = $1;"
	err := r.client.QueryRow(ctx, query, id).Scan(&newUser.Id, &newUser.Name, &newUser.Balance)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (r *Repository) Update(
	ctx context.Context, user user.User, transaction float64, sign string) error {

	userBalance := user.Balance
	newBalance, err := transactionOperation(userBalance, transaction, sign)

	if err != nil {
		return err
	}

	query := "UPDATE public.users SET balance = $1 WHERE user_balance = $2"
	_, err = r.client.Exec(ctx, query, newBalance, user.Id)

	if err != nil {
		return err
	}

	return nil
}

func transactionOperation(userBalance, transaction float64, sign string) (
	newBalance float64, err error) {

	err = nil

	if sign == "plus" {
		newBalance = userBalance + transaction
	} else {
		newBalance = userBalance - transaction
	}

	if newBalance < 0 {
		err = fmt.Errorf("the balance is not enough. Cancelling the operation")
	}
	return
}
