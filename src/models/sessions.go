package models

import "time"

// this map stores the users Sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}

// each Session contains the username of the user and the time at which it expires
type Session struct {
	Username string
	Expiry   time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
