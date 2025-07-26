package service

import (
	"encoding/json"
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

func (s *service) CreateRow(dbObj any) error {
	result := s.db.Create(dbObj)
	return result.Error
}

func (s *service) CreateRowFromJSON(model string, jsonData []byte) error {
	t, err := s.getTypeWithData(model)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, t)
	if err != nil {
		return err
	}

	result := s.db.Create(t)
	return result.Error
}

func (s *service) FindRows(model string, limit, page int) (any, error) {
	ts, err := s.getTypeSlice(model)
	if err != nil {
		return nil, err
	}

	result := s.db.Limit(limit).Offset(limit * (page - 1)).Find(&ts)

	return ts, result.Error
}

func (s *service) FindRowCount(model string) (int64, error) {
	ts, err := s.getTypeSlice(model)
	if err != nil {
		return 0, err
	}

	var count int64
	result := s.db.Model(&ts).Count(&count)

	return count, result.Error
}

func (s *service) getType(model string) (any, error) {
	t, ok := s.models()[model]
	if !ok {
		return nil, fmt.Errorf("model %s not found", model)
	}

	obj := reflect.New(t.ModelType).Interface()
	return obj, nil
}

func (s *service) getTypeWithData(model string) (any, error) {
	obj, err := s.getType(model)
	if err != nil {
		return nil, err
	}

	reflect.ValueOf(obj).Elem().FieldByName("ID").Set(reflect.ValueOf(uuid.New()))
	return obj, nil
}

func (s *service) getTypeSlice(model string) (any, error) {
	t, ok := s.models()[model]
	if !ok {
		return nil, fmt.Errorf("model %s not found", model)
	}

	st := reflect.SliceOf(t.ModelType)
	slice := reflect.MakeSlice(st, 0, 0)

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
