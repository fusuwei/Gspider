package mysql

import (
	"fmt"
	"github.com/fusuwei/gspider/pkg/logger"
	"github.com/fusuwei/gspider/pkg/storage"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLoggger "gorm.io/gorm/logger"
	"time"
)

type Mysql struct {
	Username string
	Password string
	Addr     string
	DBName   string
	Port     int
	DB       *gorm.DB
	logger   *logger.Logger
}

func NewMysql(user, pwd, addr, dbName string, Port int, logger *logger.Logger) *Mysql {
	return &Mysql{
		Username: user,
		Password: pwd,
		Addr:     addr,
		DBName:   dbName,
		Port:     Port,
		logger:   logger,
	}
}

func (m *Mysql) OpenConnection() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Addr, m.Port, m.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLoggger.New(m.logger.GetLogger(), gormLoggger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormLoggger.Silent,
		}),
	})
	if err != nil {
		m.logger.Error(err.Error())
		return
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		m.logger.Error(err.Error())
		return
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	m.DB = db
	return
}

func (m *Mysql) Save(item *storage.Item) error {
	err := m.DB.Create(item.Body).Error
	return err
}
