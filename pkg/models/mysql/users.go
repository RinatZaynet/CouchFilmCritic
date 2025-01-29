package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models"
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

func (manager *ManagerDB) InsertUser(nickName, email, passwordHash string) (userID int, err error) {
	result, err := manager.Database.Exec(sqlInsertUser, nickName, email, passwordHash)
	if err != nil {
		// не отдавать err, написать свою ошибку
		return 0, fmt.Errorf("an error occurred while insert user %s in the InsertUser(). Error: %w", nickName, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		// не отдавать err, написать свою ошибку
		return 0, fmt.Errorf("an error occurred while getting %s user ID in the InsertUser(). Error: %w", nickName, err)
	}
	return int(id), nil
}

func (manager *ManagerDB) GetUserByNickName(nickName string) (user *models.User, err error) {
	row := manager.Database.QueryRow(sqlGetUserByNickName, nickName)
	user = &models.User{}
	err = row.Scan(&user.ID, &user.NickName, &user.Email, &user.PasswordHash, &user.SignUpDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("an error occurred while get user %s in the GetUserByNickName(). Error: %w", nickName, models.ErrNoRows)
		}
		// не отдавать err, написать свою ошибку
		return nil, fmt.Errorf("an error occurred while get user %s in the GetUserByNickName(). Error: %w", nickName, err)
	}
	return user, nil
}

func (manager *ManagerDB) IsNickNameUnique(nickName string) (unique bool, err error) {
	user := &models.User{}
	row := manager.Database.QueryRow(sqlIsNickNameUnique, nickName)
	err = row.Scan(&user.NickName)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (manager *ManagerDB) IsEmailUnique(email string) (unique bool, err error) {
	user := &models.User{}
	row := manager.Database.QueryRow(sqlIsEmailUnique, email)
	err = row.Scan(&user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
