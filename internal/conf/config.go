package conf

type Web struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Log struct {
	// Level: debug, info, warning, error
	Level string `json:"level"`
	Path  string `json:"path"`
}

type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASS"`
	Name        string `json:"name" env:"NAME"`
	DBFile      string `json:"db_file" env:"FILE"`
	TablePrefix string `json:"table_prefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"ssl_mode" env:"SSL_MODE"`
	DSN         string `json:"dsn" env:"DSN"`
}

type Config struct {
	Web      Web      `json:"web"`
	Log      Log      `json:"log"`
	Database Database `json:"database"`
}

func GetDefaultConfig() *Config {
	return &Config{
		Web: Web{
			Host: "localhost",
			Port: 8080,
		},
		Log: Log{
			Level: "debug",
			Path:  "backend.log",
		},
		Database: Database{
			Type:   "sqlite3",
			DBFile: "database.db",
		},
	}
}
