package db

import (
	"app_server/pkg/idgen"
	"log"
	"os"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

// SetDB 设置数据库连接（主要用于测试）
func SetDB(database *gorm.DB) {
	db = database
}

func Init(dsn string, debug bool) error {
	var err error
	
	config := &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{SingularTable: true},
	}
	
	// 如果开启debug模式，配置日志记录器
	if debug {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // 慢 SQL 阈值
				LogLevel:                  logger.Info,   // 日志级别
				IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,          // 彩色打印
			},
		)
		config.Logger = newLogger
	}
	
	db, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return err
	}

	db.Callback().Create().Before("gorm:create").Register("set_custom_id", setCustomID)

	return nil
}

func setCustomID(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	idField, ok := db.Statement.Schema.FieldsByName["ID"]
	if !ok {
		return
	}

	if idField.FieldType.Kind() != reflect.Int && idField.FieldType.Kind() != reflect.Uint {
		return
	}

	switch db.Statement.ReflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
			elem := db.Statement.ReflectValue.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			if elem.FieldByName("ID").IsZero() {
				id := idgen.NewID().Int()
				idField := elem.FieldByName("ID")

				// 根据ID字段类型设置适当的值
				switch idField.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					idField.SetInt(int64(id))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					idField.SetUint(uint64(id))
				}
			}
		}
	case reflect.Struct:
		if db.Statement.ReflectValue.FieldByName("ID").IsZero() {
			id := idgen.NewID().Int()
			db.Statement.SetColumn("id", id)
		}
	}
}
