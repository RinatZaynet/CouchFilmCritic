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
	sqlGetUserByNickName = `SELECT id, nick_name, email, password_hash, signup_date FROM users
	WHERE nick_name = ?`
	sqlIsNickNameUnique = `SELECT nick_name FROM users
	WHERE nick_name = ?`
	sqlIsEmailUnique = `SELECT email FROM users
	WHERE email = ?`
)

func (manager *ManagerDB) InsertUser(nickName, email, passwordHash string) (int, error) {
	const fn = "storage.mysql.managerDB.InsertUser"

	result, err := manager.Database.Exec(sqlInsertUser, nickName, email, passwordHash)
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

func (manager *ManagerDB) GetUserByNickName(nickName string) (*storage.User, error) {
	const fn = "storage.mysql.managerDB.GetUserByNickName"

	row := manager.Database.QueryRow(sqlGetUserByNickName, nickName)

	user := &storage.User{}
	err := row.Scan(&user.ID, &user.NickName, &user.Email, &user.PasswordHash, &user.SignUpDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, storage.ErrNoRows)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (manager *ManagerDB) IsNickNameUnique(nickName string) (unique bool, err error) {
	const fn = "storage.mysql.managerDB.IsNickNameUnique"
	user := &storage.User{}

	row := manager.Database.QueryRow(sqlIsNickNameUnique, nickName)

	err = row.Scan(&user.NickName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return false, nil
}

func (manager *ManagerDB) IsEmailUnique(email string) (unique bool, err error) {
	const fn = "storage.mysql.managerDB.IsEmailUnique"
	user := &storage.User{}

	row := manager.Database.QueryRow(sqlIsEmailUnique, email)

	err = row.Scan(&user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return false, nil
}
