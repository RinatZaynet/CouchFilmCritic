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
)

func (manager *ManagerDB) InsertUser(nickName string, email string, passwordHash string) (userID int, err error) {
	result, err := manager.Database.Exec(sqlInsertUser, nickName, email, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("an error occurred while insert user %s in the InsertUser(). Error: %w", nickName, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("an error occurred while getting %s user ID in the InsertUser(). Error: %w", nickName, err)
	}
	return int(id), nil
}

func (manager *ManagerDB) GetUserByNickName(nickName string) (user *models.User, err error) {
	user = &models.User{}
	row := manager.Database.QueryRow(sqlGetUserByNickName, nickName)
	err = row.Scan(&user.ID, &user.NickName, &user.Email, &user.PasswordHash, &user.SignUp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("an error occurred while get user %s in the GetUserByNickName(). Error: %w", nickName, models.ErrNoRows)
		}
		return nil, fmt.Errorf("an error occurred while get user %s in the GetUserByNickName(). Error: %w", nickName, err)
	}
	return user, nil
}
