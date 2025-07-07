package config

type AdminConf struct {
	JWTSecret string
	err       error
}

func (ac AdminConf) IsValid() bool {
	return ac.err == nil
}

var AdminConfig = AdminConf{}
