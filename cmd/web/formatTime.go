package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/pkg/models"
)

var (
	ErrReviewsIsInvalid  = errors.New("reviews is invalid")
	ErrLocationIsInvalid = errors.New("location is invalid")
)

func formatTimeReviews(reviews []*models.Review, location string) (err error) {
	if reviews == nil {
		return fmt.Errorf("an error occurred in convertReviewsTimeZone(). Error: %w. Witch value: %#v", ErrReviewsIsInvalid, reviews)
	}
	if location == "" {
		return fmt.Errorf("an error occurred in convertReviewsTimeZone(). Error: %w. Witch value: %s", ErrLocationIsInvalid, location)
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		return fmt.Errorf("an error occurred in while Parse location in convertReviewsTimeZone(). Error: %w. Witch value: %s", ErrLocationIsInvalid, location)
	}
	for _, review := range reviews {
		formatTime := review.CreateDate.In(loc).Format("15:04:05 01.2.2006")
		review.FormatCreateDate = formatTime
	}
	return nil
}
