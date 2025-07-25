package model

import (
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type Field struct {
	Name   string
	HtmlId string
}

func NewField(name string) Field {
	return Field{
		Name:   name,
		HtmlId: fieldNameToHTMLId(name),
	}
}

type Model struct {
	ModelType          reflect.Type
	CreationFieldNames []Field
	HiddenFieldNames   []Field
}

var (
	registeredModels map[string]Model = make(map[string]Model)

	ErrUnableToRegisterType = errors.New("unable to register type, must be a struct")
	ErrUnknownAnnotation    = errors.New("unknown runway annotation, only 'hidden' and 'create' are supported")
)

func Register(model any) error {
	t := reflect.TypeOf(model)

	// check if the model is a struct
	if t.Kind() != reflect.Struct {
		return ErrUnableToRegisterType
	}

	m := Model{
		ModelType: t,
	}

	for i := range t.NumField() {
		tag, ok := t.Field(i).Tag.Lookup("runway")
		if ok {
			err := checkAnnotation(tag)
			if err != nil {
				return err
			}

			if isCreationField(tag) {
				m.CreationFieldNames = append(m.CreationFieldNames, NewField(t.Field(i).Name))
			}

			if isHiddenField(tag) {
				m.HiddenFieldNames = append(m.HiddenFieldNames, NewField(t.Field(i).Name))
			}
		}
	}

	name := strings.ToLower(t.Name())
	registeredModels[name] = m

	return nil
}

func GetRegisteredModels() map[string]Model {
	return registeredModels
}

func checkAnnotation(annotation string) error {
	allowed_parts := []string{
		"hidden",
		"create",
	}
	parts := strings.SplitSeq(annotation, ",")

	for p := range parts {
		if !slices.Contains(allowed_parts, p) {
			return ErrUnknownAnnotation
		}
	}

	return nil
}

func isHiddenField(annotation string) bool {
	parts := strings.Split(annotation, ",")
	return slices.Contains(parts, "hidden")
}

func isCreationField(annotation string) bool {
	parts := strings.Split(annotation, ",")
	return slices.Contains(parts, "create")
}

var (
	replaceRegexp     = regexp.MustCompile(`([a-z-]*)([A-Z]+)(.*)`)
	nonAlphaNumRegexp = regexp.MustCompile(`[^a-zA-Z0-9-]+`)
)

func fieldNameToHTMLId(name string) string {
	if nonAlphaNumRegexp.MatchString(name) {
		panic("filed names can just contain letters, numbers and hyphens")
	}

	var newName string
	parts := replaceRegexp.FindStringSubmatch(name)

	if len(parts) == 0 {
		return name
	}

	newName = parts[1] + "-" + strings.ToLower(parts[2]) + parts[3]
	return fieldNameToHTMLId(newName)
}
