package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	. "go-bank-v2/internal/types"
)

func (s *PostgresqlStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO Account (Balance, OwnerID, Created)
              VALUES ($1, $2, $3)
              RETURNING AccountNumber`

	err := s.db.QueryRow(
		query,
		acc.Balance,
		acc.OwnerID,
		acc.Created,
	).Scan(&acc.AccountNumber)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresqlStore) CreateUser(user *User) error {
	query := `INSERT INTO "User" (FirstName, LastName, MemberSince, EncryptedPassword)
              VALUES ($1, $2, $3, $4)
              RETURNING ID`

	// Execute the query and scan the result into user.ID
	err := s.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.MemberSince,
		user.EncryptedPassword,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresqlStore) GetUserByID(userID int) (*User, error) {
	query := `SELECT ID, FirstName, LastName, MemberSince, EncryptedPassword FROM "User" WHERE ID = $1`

	user := &User{}

	err := s.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.MemberSince,
		&user.EncryptedPassword,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User not found
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *PostgresqlStore) UpdateAccountBalance(*Account) error {
	return nil
}

func (s *PostgresqlStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresqlStore) GetAccounts(ownerID int) ([]*Account, error) {
	query := `SELECT * FROM Account WHERE ownerId = $1`

	rows, err := s.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var accounts []*Account

	for rows.Next() {
		account := &Account{}
		err := rows.Scan(
			&account.AccountNumber,
			&account.Balance,
			&account.OwnerID,
			&account.Created,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *PostgresqlStore) GetAccountByID(accountNumber int) (*Account, error) {
	return nil, nil
}

func (s *PostgresqlStore) GetAccountByNumber(accountNumber int) (*Account, error) {
	query := `SELECT * FROM Account WHERE accountNumber = $1`

	account := &Account{}

	err := s.db.QueryRow(query, accountNumber).Scan(
		&account.AccountNumber,
		&account.Balance,
		&account.OwnerID,
		&account.Created,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Account not found
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

func (s *PostgresqlStore) DeleteUser(i int) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresqlStore) GetUserAccount(accounts []*Account, err error) {
	//TODO implement me
	panic("implement me")
}
