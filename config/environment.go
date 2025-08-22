package config

import "os"

type environment interface {
	GetEnv(key string) string
	GetHome() (string, error)
	GetPwd() (string, error)
}

type realEnvironment struct{}

func (r *realEnvironment) GetEnv(key string) string { return os.Getenv(key) }
func (r *realEnvironment) GetHome() (string, error) { return os.UserHomeDir() }
func (r *realEnvironment) GetPwd() (string, error)  { return os.Getwd() }
