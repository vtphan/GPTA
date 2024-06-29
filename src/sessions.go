package main

import "time"

// each session contains the username of the user and the time at which it expires
type session struct {
	user   *User
	expiry time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
