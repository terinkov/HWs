package storage

import (
	"errors"
	"github.com/terinkov_HW2/models"
)


type RamStorage struct {
	// uuid (string) ->  
	tasks map[string]models.Task
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		tasks: make(map[string]models.Task),
	}
}

func (rs *RamStorage) GetTaskById(key string) (*models.Task, error) {
	value, exists := rs.tasks[key]
	if !exists {
		return nil, errors.New("uuid not found")
	}
	return &value, nil
}

func (rs *RamStorage) UpdateTaskById(value models.Task) error {
	rs.tasks[value.UUID] = value
	return nil	
}

func (rs *RamStorage) PostTaskById(value models.Task) error {
	if _, exists := rs.tasks[value.UUID]; exists {
		return errors.New("uuid already exists")
	}
	rs.tasks[value.UUID] = value
	return nil
}

func (rs *RamStorage) DeleteTaskById(key string) error {
	if _, exists := rs.tasks[key]; !exists {
		return errors.New("uuid not found")
	}
	delete(rs.tasks, key)
	return nil
}
