package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

var (
	sqlInsertReview = `INSERT INTO reviews (work_title, genres, work_type, review, rating, create_date, author_user_id)
	VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	sqlGetLatestReviews = `SELECT 
    r.id AS review_id,
    r.work_title,
    r.genres,
    r.work_type,
    r.review,
    r.rating,
    r.create_date,
    u.nick_name
	FROM reviews r
	JOIN users u ON r.author_user_id = u.id 
	ORDER BY r.create_date DESC LIMIT 10;`

	sqlGetReviewsByAuthor = `SELECT 
    r.id AS review_id,
    r.work_title,
    r.genres,
    r.work_type,
    r.review,
    r.rating,
    r.create_date,
    u.nick_name 
	FROM reviews r 
	JOIN users u ON r.author_user_id = u.id 
	WHERE u.nick_name = ? 
	ORDER BY r.create_date DESC;`

	sqlDeleteReviewByID = `DELETE
	FROM reviews
	WHERE id = ?;`

	sqlGetReviewByID = `SELECT
	work_title,
	genres,
	work_type,
	review,
	rating
	FROM reviews
	WHERE id = ?;`
)

func (manager *ManagerDB) InsertReview(workTitle, genres, workType, review string, rating int, authorUserID int) (reviewID int, err error) {
	const fn = "storage.mysql.managerDB.InsertReview"
	result, err := manager.Database.Exec(sqlInsertReview, workTitle, genres, workType, review, rating, authorUserID)
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

func (manager *ManagerDB) GetLatestReviews() ([]*storage.Review, error) {
	const fn = "storage.mysql.managerDB.GetLatestReviews"
	rows, err := manager.Database.Query(sqlGetLatestReviews)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	reviews := make([]*storage.Review, 0, 10)
	defer rows.Close()

	for rows.Next() {
		s := &storage.Review{}

		err := rows.Scan(&s.ID, &s.WorkTitle, &s.Genres, &s.WorkType, &s.Review, &s.Rating, &s.CreateDate, &s.Author)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan: %w", fn, err)
		}
		reviews = append(reviews, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return reviews, nil
}

func (manager *ManagerDB) GetReviewsByAuthor(author string) ([]*storage.Review, error) {
	const fn = "storage.mysql.managerDB.GetReviewsByAuthor"
	rows, err := manager.Database.Query(sqlGetReviewsByAuthor, author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	reviews := make([]*storage.Review, 0, 20)
	defer rows.Close()

	for rows.Next() {
		s := &storage.Review{}

		err := rows.Scan(&s.ID, &s.WorkTitle, &s.Genres, &s.WorkType, &s.Review, &s.Rating, &s.CreateDate, &s.Author)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan: %w", fn, err)
		}
		reviews = append(reviews, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return reviews, nil
}

func (manager *ManagerDB) DeleteReviewByID(reviewID int) error {
	const fn = "storage.mysql.managerDB.DeleteReviewByID"

	_, err := manager.Database.Exec(sqlDeleteReviewByID, reviewID)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (manager *ManagerDB) GetReviewByID(id int) (*storage.Review, error) {
	const fn = "storage.mysql.managerDB.GetReviewByID"

	row := manager.Database.QueryRow(sqlGetReviewByID, id)

	review := &storage.Review{}
	err := row.Scan(&review.WorkTitle, &review.Genres, review.WorkType, &review.Review, &review.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return review, nil
}
