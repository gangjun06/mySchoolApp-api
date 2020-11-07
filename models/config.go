package models

type Config struct {
	Server struct {
		Port      int
		Debug     bool
		JWTSecret string
	}
	DB struct {
		UseSqlite bool
		Server    struct {
			Hostname string
			Port     int
			Username string
			Password string
			DBName   string
		}
	}
}
