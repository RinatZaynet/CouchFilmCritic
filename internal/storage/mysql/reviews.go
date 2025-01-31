package mysql

import (
	"fmt"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models"
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
)

func (manager *ManagerDB) InsertReview(workTitle, genres, workType, review string, rating float64, authorUserID int) (reviewID int, err error) {
	result, err := manager.Database.Exec(sqlInsertReview, workTitle, genres, workType, review, rating, authorUserID)
	if err != nil {
		// не отдавать err, написать свою ошибку
		return 0, fmt.Errorf("an error occurred while insert review in the InsertReview(). Error: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		// не отдавать err, написать свою ошибку
		return 0, fmt.Errorf("an error occurred while getting review ID in the InsertReview(). Error: %w", err)
	}
	return int(id), nil
}

func (manager *ManagerDB) GetLatestReviews() ([]*models.Review, error) {

	rows, err := manager.Database.Query(sqlGetLatestReviews)
	if err != nil {
		// не отдавать err, написать свою ошибку
		return nil, fmt.Errorf("an error occurred while getting last 10 review in the GetLatestReviews(). Error: %w", err)
	}
	reviews := make([]*models.Review, 0, 10)
	defer rows.Close()
	for rows.Next() {
		s := &models.Review{}

		err := rows.Scan(&s.ID, &s.WorkTitle, &s.Genres, &s.WorkType, &s.Review, &s.Rating, &s.CreateDate, &s.Author)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while parsing last 10 review in the GetLatestReviews(). Error: %w", err)
		}
		reviews = append(reviews, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred while parsing last 10 review in the GetLatestReviews(). Error: %w", err)
	}
	return reviews, nil
}

func (manager *ManagerDB) GetReviewsByAuthor(author string) ([]*models.Review, error) {
	rows, err := manager.Database.Query(sqlGetReviewsByAuthor, author)
	if err != nil {
		// не отдавать err, написать свою ошибку
		return nil, fmt.Errorf("an error occurred while getting reviews by Author:%s, in the GetReviewsByAuthor(). Error: %w", author, err)
	}
	reviews := make([]*models.Review, 0, 20)
	defer rows.Close()
	for rows.Next() {
		s := &models.Review{}

		err := rows.Scan(&s.ID, &s.WorkTitle, &s.Genres, &s.WorkType, &s.Review, &s.Rating, &s.CreateDate, &s.Author)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while parsing reviews by Author:%s, in the GetReviewsByAuthor(). Error: %w", author, err)
		}
		reviews = append(reviews, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred while parsing reviews by Author:%s, in the GetReviewsByAuthor(). Error: %w", author, err)
	}
	return reviews, nil
}
