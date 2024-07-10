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

type Oauth struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
}

type GoogleOauth Oauth

type LineOauth Oauth

type Config struct {
	Web         Web         `json:"web"`
	Log         Log         `json:"log"`
	Database    Database    `json:"database"`
	GoogleOauth GoogleOauth `json:"google"`
	LineOauth   LineOauth   `json:"line"`
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
		GoogleOauth: GoogleOauth{
			ClientID:     "372889357683-xxxxxxxxxx.apps.googleusercontent.com",
			ClientSecret: "GOCSPX-xxxxxxxxxx-fmXr0Dc",
			RedirectURL:  "http://localhost:8080/api/google/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		},
		LineOauth: LineOauth{
			ClientID:     "20xxxxxx94",
			ClientSecret: "4bxxxxxxxxxxxxxxxxxxxxxxxxxxxx64",
			RedirectURL:  "http://localhost:8080/api/line/callback",
			Scopes:       []string{"profile"},
		},
	}
}
