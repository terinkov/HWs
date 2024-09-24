package storage

import (
	"errors"
	// "sync"

	"github.com/terinkov_HW2/models"
)

// Реализация репозитория сессий в оперативной памяти
type RamSessionRepository struct {
	sessions map[string]models.Session
}

// Создание нового репозитория сессий в оперативной памяти
func NewRamSessionRepository() *RamSessionRepository {
	return &RamSessionRepository{
		sessions: make(map[string]models.Session),
	}
}

// // Создание новой сессии	
// func (rs *RamSessionRepository) CreateSession(session *models.Session) error {
// 	rs.sessions[session.UserID] = session
// 	return nil
// }

// Получение сессии по ID
func (rs *RamSessionRepository) GetSession(sessionID string) (*models.Session, error) {
	session, ok := rs.sessions[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}
	return &session, nil
}

// Создание сессиии по токену и самой сессии
func (rs *RamSessionRepository) PostSession(session models.Session) error{
	if _, alreadyExists := rs.sessions[session.SessionId]; alreadyExists {
		return errors.New("this user's session already exists")
	}
	rs.sessions[session.SessionId] = session
	return nil
}

// Удаление сессии по ID
func (rs *RamSessionRepository) DeleteSession(sessionID string) error {
	_, exists := rs.sessions[sessionID]
	if !exists {
		return errors.New("Session not found")
	}
	delete(rs.sessions, sessionID)

	return nil
}