package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
	Log    Log    `yaml:"log"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Log struct {
	Level       string `yaml:"level"`
	Output      string `yaml:"output"`
	Development bool   `yaml:"development"`
	Rotate      Rotate `yaml:"rotate"`
}

type Rotate struct {
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"maxage"`
	MaxBackups int    `yaml:"maxbackups"`
	LocalTime  bool   `yaml:"localtime"`
	Compress   bool   `yaml:"compress"`
}

type Data struct {
	DB    DB `yaml:"database"`
	Redis `yaml:"redis"`
}

type DB struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	ConnTimeout     time.Duration `yaml:"conn_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_dile_time"`
	ConnMaxLifeTime time.Duration `yaml:"conn_max_life_time"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
}

type Redis struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
}

func Parse(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
