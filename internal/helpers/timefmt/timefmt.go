package timefmt

import (
	"errors"
	"fmt"
	"time"

	"github.com/RinatZaynet/CouchFilmCritic/internal/storage"
)

var (
	ErrReviewsNotFound   = errors.New("reviews not found")
	ErrLocationNotFound  = errors.New("location not found")
	ErrLocationIsInvalid = errors.New("location is invalid")
)

func TimeReviewsFmt(reviews []*storage.Review, location string) (err error) {
	const fn = "timefmt.TimeReviewFmt"

	if reviews == nil {
		return fmt.Errorf("%s: %w", fn, ErrReviewsNotFound)
	}
	if location == "" {
		return fmt.Errorf("%s: %w", fn, ErrLocationNotFound)
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return fmt.Errorf("%s: %w. Value: %s", fn, ErrLocationIsInvalid, location)
	}

	for _, review := range reviews {
		formatTime := review.CreateDate.In(loc).Format("15:04:05 01.2.2006")
		review.FormatCreateDate = formatTime
	}

	return nil
}
