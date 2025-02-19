package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

var (
	sqlInsertUser = `INSERT INTO users (nick_name, email, password_hash, signup_date)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`
	sqlGetUserByNickname = `SELECT id, nick_name, email, password_hash, signup_date FROM users
	WHERE nick_name = ?`
	sqlIsUniqueNickname = `SELECT nick_name FROM users
	WHERE nick_name = ?`
	sqlIsUniqueEmail = `SELECT email FROM users
	WHERE email = ?`
)

func (manager *ManagerDB) InsertUser(nickname, email, passwordHash string) (int, error) {
	const fn = "storage.mysql.managerDB.InsertUser"

	result, err := manager.Database.Exec(sqlInsertUser, nickname, email, passwordHash)
	if err != nil {
		// if err == duplicate ...
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return int(id), nil
}

func (manager *ManagerDB) GetUserByNickname(nickname string) (*storage.User, error) {
	const fn = "storage.mysql.managerDB.GetUserByNickname"

	row := manager.Database.QueryRow(sqlGetUserByNickname, nickname)

	user := &storage.User{}
	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.PasswordHash, &user.SignUpDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, storage.ErrNoRows)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (manager *ManagerDB) IsUniqueNickname(nickname string) (unique bool, err error) {
	const fn = "storage.mysql.managerDB.IsUniqueNickname"
	user := &storage.User{}

	row := manager.Database.QueryRow(sqlIsUniqueNickname, nickname)

	err = row.Scan(&user.Nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return false, nil
}

func (manager *ManagerDB) IsUniqueEmail(email string) (unique bool, err error) {
	const fn = "storage.mysql.managerDB.IsUniqueEmail"
	user := &storage.User{}

	row := manager.Database.QueryRow(sqlIsUniqueEmail, email)

	err = row.Scan(&user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return false, nil
}
