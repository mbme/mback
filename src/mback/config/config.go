package config

const (
	CONF_ROOT = "/home/mbme/temp/conf"

	LOG_LEVEL = 1

	USER = "mbme"

	CONF_FILE_NAME = ".mback"
)

var root_uid string

func GetRepositoryRoot() string {
	return CONF_ROOT
}
