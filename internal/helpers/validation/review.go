package validation

import (
	"strconv"
	"strings"
)

var (
	validGenres = map[string]struct{}{
		"биография":       {},
		"боевик":          {},
		"вестерн":         {},
		"военный":         {},
		"детектив":        {},
		"детский":         {},
		"документальный":  {},
		"драма":           {},
		"история":         {},
		"комедия":         {},
		"короткометражка": {},
		"криминал":        {},
		"мелодрама":       {},
		"музыка":          {},
		"мюзикл":          {},
		"приключения":     {},
		"семейный":        {},
		"спорт":           {},
		"триллер":         {},
		"ужасы":           {},
		"фантастика":      {},
		"фэнтези":         {},
	}

	validWorkTypes = map[string]struct{}{
		"Фильм":      {},
		"Сериал":     {},
		"Аниме":      {},
		"Мультфильм": {},
	}
)

const (
	minLenWorkTitle = 8
	maxLenWorkTitle = 200

	maxLenReview = 500
)

func IsValidWorkTitle(workTitle string) bool {
	tWorkTitle := strings.Trim(workTitle, " ")

	if len(tWorkTitle) < minLenWorkTitle ||
		len(tWorkTitle) > maxLenWorkTitle {
		return false
	}

	return true
}

func IsValidGenres(genres []string) bool {
	for _, genre := range genres {
		if _, ok := validGenres[genre]; !ok {
			return false
		}
	}

	return true
}

func IsValidWorkType(workType string) bool {
	if _, ok := validWorkTypes[workType]; !ok {
		return false
	}

	return true
}

func IsValidReview(review string) bool {
	tReview := strings.Trim(review, " ")

	return len(tReview) < maxLenReview
}

func IsValidRating(rating string) bool {
	iRating, err := strconv.Atoi(rating)

	if err != nil {
		return false
	}

	if iRating > 10 {
		return false
	}

	return true
}
