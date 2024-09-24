package storage

import (
	"errors"
	"github.com/terinkov_HW2/models"
)

// RamUserRepository - реализация UserRepository в оперативной памяти
type RamUserRepository struct {
	users map[string]models.User // key - это userLogin а не userID
}

// NewRamUserRepository - создание нового репозитория пользователей в оперативной памяти
func NewRamUserRepository() *RamUserRepository {
	return &RamUserRepository{
		users: make(map[string]models.User),
	}
}

// GetUserByUserLogin - получение пользователя по Login
func (rur *RamUserRepository) GetUserByUserLogin(UserLogin string) (*models.User, error) {
	user, ok := rur.users[UserLogin]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// PostUserByNameAndPassword - создание нового пользователя
func (rur *RamUserRepository) PostUserByNameAndPassword(user models.User) error {
	if _, ok := rur.users[user.Login]; ok {
		return errors.New("user already exists")
	}
	rur.users[user.Login] = user
	return nil
}

// UpdateUser - обновление информации о пользователе
func (rur *RamUserRepository) UpdateUser(user models.User) error {
	if _, ok := rur.users[user.Login]; !ok {
		return errors.New("user not found")
	}
	rur.users[user.Login] = user
	return nil
}

// DeleteUserByUserLogin - удаление пользователя
func (rur *RamUserRepository) DeleteUserByUserLogin(userLogin string) error {
	if _, ok := rur.users[userLogin]; !ok {
		return errors.New("user not found")
	}
	delete(rur.users, userLogin)
	return nil
}

// LoginUser - проверка аутентификации пользователя
func (rur *RamUserRepository) LoginUser(login string, password string) (*models.User, error) {
	user, err := rur.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
   
	// Проверка пароля
	if user.Password != password {
		return nil, errors.New("invalid password")
	}
   
	return user, nil
}

// GetUserByLogin - получение пользователя по логину
func (rur *RamUserRepository) GetUserByLogin(login string) (*models.User, error) {
	user, ok := rur.users[login]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}