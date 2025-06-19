package config

type AppConfig struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type JawtConfig struct {
	Project ProjectConfig `json:"project"`
	Server  ServerConfig  `json:"server"`
}

type ProjectConfig struct {
	Name string `json:"name"`
}

type ServerConfig struct {
	Port int `json:"port"`
}
