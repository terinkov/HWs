package models

import (
	"time"
)

type Session struct {
	UserLogin 		string
	SessionId 		string
 	ExpirityTime 	time.Time `json:"expirity time"`
}

func (s Session) getSessionExpired() bool{ 
	return s.ExpirityTime.Before(time.Now())
} // 1 - expired; 2 - not expired


