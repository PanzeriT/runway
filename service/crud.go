package service

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

func (s *service) GetModelNames() []string {
	var ns []string

	for k := range s.models() {
		ns = append(ns, k)
	}

	return ns
}

func (s *service) CreateRow(model string, dbObj any) error {
	result := s.db.Create(dbObj)
	return result.Error
}

func (s *service) FindRows(model string, limit, offset int) (any, error) {
	for n := range s.GetModelNames() {
		fmt.Println("Registered model:", n)
	}
	ts, err := s.getTypeSlice(model)
	if err != nil {
		fmt.Println("Error getting type slice:", err)
		return nil, err
	}

	result := s.db.Find(&ts)

	return ts, result.Error
}

func (s *service) getType(model string) (any, error) {
	t, ok := s.models()[model]
	if !ok {
		return nil, fmt.Errorf("model %s not found", model)
	}

	obj := reflect.New(t).Interface()
	return obj, nil
}

func (s *service) getTypeSlice(model string) (any, error) {
	for k, v := range s.models() {
		fmt.Printf("registered models: %s: %v\n", k, v)
	}
	fmt.Println("getTypeSlice called for model:", model)
	t, ok := s.models()[model]
	if !ok {
		return nil, fmt.Errorf("model %s not found", model)
	}

	st := reflect.SliceOf(t)
	slice := reflect.MakeSlice(st, 0, 0)

	fmt.Println("Creating slice of type:", reflect.TypeOf(slice).String())

	return slice.Interface(), nil
}

func (s *service) DeleteRow(model string, id uuid.UUID) error {
	obj, err := s.getType(model)
	if err != nil {
		return err
	}

	result := s.db.Delete(obj, id)
	return result.Error
}
