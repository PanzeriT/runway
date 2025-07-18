package runway

import (
	"strings"

	"github.com/panzerit/runway/model"
)

var registeredModels map[string]any

func init() {
	registeredModels = make(map[string]any)
	registeredModels["user"] = model.User{}
}

func RegisterModel(name string, model any) {
	registeredModels[strings.ToLower(name)] = model
}
