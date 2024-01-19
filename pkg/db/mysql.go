package db

import (
	"time"

	l "github.com/dungnh3/skool-mn/pkg/log"
	"github.com/go-sql-driver/mysql"
	mysqlgorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	zapgorm "moul.io/zapgorm2"
)

var ll = l.New()

// MySQL is settings of a MySQL server. It contains almost same fields as mysql.Config,
// but with some different field names and tags.
type MySQL struct {
	Username  string            `yaml:"username" mapstructure:"username"`
	Password  string            `yaml:"password" mapstructure:"password"`
	Protocol  string            `yaml:"protocol" mapstructure:"protocol"`
	Address   string            `yaml:"address" mapstructure:"address"`
	Database  string            `yaml:"database" mapstructure:"database"`
	Params    map[string]string `yaml:"params" mapstructure:"params"`
	Collation string            `yaml:"collation" mapstructure:"collation"`
	Loc       *time.Location    `yaml:"location" mapstructure:"loc"`
	TLSConfig string            `yaml:"tls_config" mapstructure:"tls_config"`

	Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`

	AllowAllFiles           bool   `yaml:"allow_all_files" mapstructure:"allow_all_files"`
	AllowCleartextPasswords bool   `yaml:"allow_cleartext_passwords" mapstructure:"allow_cleartext_passwords"`
	AllowOldPasswords       bool   `yaml:"allow_old_passwords" mapstructure:"allow_old_passwords"`
	ClientFoundRows         bool   `yaml:"client_found_rows" mapstructure:"client_found_rows"`
	ColumnsWithAlias        bool   `yaml:"columns_with_alias" mapstructure:"columns_with_alias"`
	InterpolateParams       bool   `yaml:"interpolate_params" mapstructure:"interpolate_params"`
	MultiStatements         bool   `yaml:"multi_statements" mapstructure:"multi_statements"`
	ParseTime               bool   `yaml:"parse_time" mapstructure:"parse_time"`
	GoogleAuthFile          string `yaml:"google_auth_file" mapstructure:"google_auth_file"`

	LogLevel string `yaml:"log_level" mapstructure:"log_level"`
}

// FormatDSN returns MySQL DSN from settings.
func (m *MySQL) FormatDSN() string {
	um := &mysql.Config{
		User:                    m.Username,
		Passwd:                  m.Password,
		Net:                     m.Protocol,
		Addr:                    m.Address,
		DBName:                  m.Database,
		Params:                  m.Params,
		Collation:               m.Collation,
		Loc:                     m.Loc,
		TLSConfig:               m.TLSConfig,
		Timeout:                 m.Timeout,
		ReadTimeout:             m.ReadTimeout,
		WriteTimeout:            m.WriteTimeout,
		AllowAllFiles:           m.AllowAllFiles,
		AllowCleartextPasswords: m.AllowCleartextPasswords,
		AllowOldPasswords:       m.AllowOldPasswords,
		ClientFoundRows:         m.ClientFoundRows,
		ColumnsWithAlias:        m.ColumnsWithAlias,
		InterpolateParams:       m.InterpolateParams,
		MultiStatements:         m.MultiStatements,
		ParseTime:               m.ParseTime,
		AllowNativePasswords:    true,
	}
	return um.FormatDSN()
}

func (m *MySQL) GormLogLevel() logger.LogLevel {
	switch m.LogLevel {
	case "SILENT":
		return logger.Silent
	case "INFO":
		return logger.Info
	case "WARN":
		return logger.Warn
	case "ERROR":
		return logger.Error
	default:
		return logger.Silent
	}
}

const (
	maxDBIdleConns  = 10
	maxDBOpenConns  = 100
	maxConnLifeTime = 30 * time.Minute
)

func ConnectMySQL(cfg *MySQL) *gorm.DB {
	if cfg.Address == "." {
		return nil
	}

	ll.Info("Connecting mysql")
	gormDB, err := gorm.Open(mysqlgorm.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: zapgorm.New(ll.Logger).LogMode(cfg.GormLogLevel()),
	})
	if err != nil {
		ll.Fatal("Error open mysql", l.Error(err))
	}

	err = gormDB.Raw("SELECT 1").Error
	if err != nil {
		ll.Fatal("Error querying SELECT 1", l.Error(err))
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		ll.Fatal("Error get sql DB", l.Error(err))
	}

	sqlDB.SetMaxIdleConns(maxDBIdleConns)
	sqlDB.SetMaxOpenConns(maxDBOpenConns)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	ll.Info("Connected mysql successfully")
	return gormDB
}
