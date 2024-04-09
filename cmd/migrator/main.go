package migrator

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string         `yaml:"env"`
	Database      DatabaseConfig `yaml:"database" env_required:"true"`
	MigrationPath string         `yaml:"migration_path" env_required:"true"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

const (
	postgreSQL = "postgres"
	// sqLite = "sqlite"
)

func main() {
	cfg := MustLoadConfig()

	var dbURI string

	switch cfg.Database.Driver{
	case postgreSQL:
		dbURI = fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.Driver,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
	default:
		// TODO add sqlite3 support
		panic("use postgres")
	}

	m, err := migrate.New(
		"file://"+cfg.MigrationPath,
		dbURI, 
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations apply successfully")
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("can not parse config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
