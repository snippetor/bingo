package module

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/snippetor/bingo/app"
	"github.com/snippetor/bingo/errors"
	"github.com/snippetor/bingo/mvc"
)

type MySqlModule interface {
	Module
	Dial(app app.Application, addr, username, pwd, defaultDb, tbPrefix string)
	DB() *gorm.DB
	TableName(tbName string) string
	AutoMigrate(model mvc.OrmModel)

	Create(model mvc.OrmModel) bool
	Find(model mvc.OrmModel) bool
	FindAll(models interface{}) bool
	FindMany(models interface{}, limit int, orderBy string, whereAndArgs ... interface{}) bool
	Begin() *gorm.DB
	Rollback()
	Commit()
}

type mysqlModule struct {
	app      app.Application
	db       *gorm.DB
	tbPrefix string
}

func (m *mysqlModule) Dial(app app.Application, addr, username, pwd, defaultDb, tbPrefix string) {
	m.app = app
	// db
	db, err := gorm.Open("mysql", username+":"+pwd+"@tcp("+addr+")/"+defaultDb+"?charset=utf8&parseTime=True&loc=Local")
	errors.Check(err)
	m.db = db
	m.tbPrefix = tbPrefix
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tbPrefix + "_" + defaultTableName
	}
	//db.LogMode(true)
}

func (m *mysqlModule) DB() *gorm.DB {
	if m.db == nil {
		panic("- DB not be initialized, invoke InitDB at first!")
	}
	return m.db
}

func (m *mysqlModule) TableName(tbName string) string {
	if m.tbPrefix != "" {
		return m.tbPrefix + "_" + tbName
	}
	return tbName
}

func (m *mysqlModule) AutoMigrate(model mvc.OrmModel) {
	model.Init(m.app, m.db, model)
	m.DB().AutoMigrate(model)
}

func (m *mysqlModule) Create(model mvc.OrmModel) bool {
	model.Init(m.app, m.db, model)
	res := m.DB().Create(model)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *mysqlModule) Find(model mvc.OrmModel) bool {
	model.Init(m.app, m.db, model)
	res := m.DB().Where(model).First(model)
	return res.Error == nil
}

func (m *mysqlModule) FindAll(models interface{}) bool {
	res := m.DB().Find(models)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *mysqlModule) FindMany(models interface{}, limit int, orderBy string, whereAndArgs ... interface{}) bool {
	db := m.DB()
	if limit > 0 {
		db = db.Limit(limit)
	}
	if orderBy != "" {
		db = db.Order(orderBy)
	}
	if len(whereAndArgs) > 0 && len(whereAndArgs)%2 == 0 {
		var args = make(map[string]interface{})
		for i := 0; i < len(whereAndArgs); i += 2 {
			args[whereAndArgs[i].(string)] = whereAndArgs[i+1]
		}
		db = db.Where(args)
	}
	db = db.Find(models)
	if db.Error != nil {
		panic(db.Error)
	}
	return true
}

func (m *mysqlModule) Begin() *gorm.DB {
	return m.DB().Begin()
}

func (m *mysqlModule) Rollback() {
	m.DB().Rollback()
}

func (m *mysqlModule) Commit() {
	m.DB().Rollback()
}

func (m *mysqlModule) Close() {
	if m.db != nil {
		m.db.Close()
	}
}
