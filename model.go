package runway

import (
	"errors"
	"reflect"
	"strings"
)

var (
	registeredModels map[string]reflect.Type

	ErrUnableToRegisterType = errors.New("unable to register type, must be a struct")
	ErrUnknownAnnotation    = errors.New("unknown runway annotation, only 'hidden' is supported")
)

func MustRegisterModel(model any) {
	err := RegisterModel(model)
	if err != nil {
		Terminate(NewAppError(ErrCannotRegisterModel, err))
	}
}

func RegisterModel(model any) error {
	t := reflect.TypeOf(model)

	// check if the model is a struct
	if t.Kind() != reflect.Struct {
		return ErrUnableToRegisterType
	}

	for i := 0; i < t.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("runway")
		if ok {
			err := checkAnnotation(tag)
			if err != nil {
				return err
			}
		}
	}

	name := strings.ToLower(t.Name())
	registeredModels[name] = t

	logger.Debug("registering model", "name", name)

	return nil
}

func GetRegisteredModels() map[string]reflect.Type {
	return registeredModels
}

func checkAnnotation(annotation string) error {
	if annotation == "hidden" {
		return nil
	}

	return ErrUnknownAnnotation
}
