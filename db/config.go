package db

import "fmt"

type Config struct {
	Dialect  string `json:"dialect"` // mysql,postgres,sqlite3
	Host     string `json:"host"`    // if Dialect is `sqlite3`, host should be db file path
	ReadHost string `json:"read_host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Location string `json:"location,omitempty"`
	SSLMode  string `json:"sslmode,omitempty"` // postgres has this property, which could be "disable"
	Debug    bool   `json:"debug,omitempty"`
}

func (c Config) ReadHostDSN() (string, error) {
	if c.ReadHost == "" {
		return "", fmt.Errorf("read_host not set")
	}

	c.Host = c.ReadHost
	return c.DSN()
}

func (c Config) DSN() (string, error) {
	var dsn string
	switch c.Dialect {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&charset=utf8mb4,utf8&loc=%s",
			c.User,
			c.Password,
			"tcp",
			c.Host,
			c.Port,
			c.Database,
			c.Location,
		)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			c.Host,
			c.Port,
			c.User,
			c.Database,
			c.Password,
			c.SSLMode,
		)
	case "sqlite3":
		dsn = c.Host
	default:
		return "", fmt.Errorf("unkonow db dialect: %s", c.Dialect)
	}

	return dsn, nil
}

func SqliteInMemory() Config {
	return Config{
		Dialect: "sqlite3",
		Host:    ":memory:",
	}
}