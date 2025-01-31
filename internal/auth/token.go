package auth

import "time"

type Claims struct {
	Sub string
	Exp *time.Time
}
