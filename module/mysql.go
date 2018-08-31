package module

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/snippetor/bingo/errors"
)

// mysql
type MySqlModule interface {
	Module
	DefaultDB() MySqlDB
	DB(name string) MySqlDB
}

func NewMysqlModule(dbs map[string]MySqlDB) MySqlModule {
	return &mysqlModule{dbs}
}

type mysqlModule struct {
	dbs map[string]MySqlDB
}

func (m *mysqlModule) DefaultDB() MySqlDB {
	if db, ok := m.dbs["default"]; ok {
		return db
	}
	for _, v := range m.dbs {
		return v
	}
	return nil
}

func (m *mysqlModule) DB(name string) MySqlDB {
	if db, ok := m.dbs[name]; ok {
		return db
	}
	return nil
}

func (m *mysqlModule) Close() {
	for _, v := range m.dbs {
		v.Close()
	}
}

type MySqlDB interface {
	OrmDB() *gorm.DB
	TableName(tbName string) string
	AutoMigrate(model interface{})

	Create(model interface{}) bool
	Find(model interface{}) bool
	FindAll(models interface{}) bool
	FindMany(models interface{}, limit int, orderBy string, whereAndArgs ... interface{}) bool
	Begin() *gorm.DB
	Rollback()
	Commit()
	Close()
}

type mysqlDB struct {
	tbPrefix string
	db       *gorm.DB
}

func NewMysqlDB(addr, username, pwd, defaultDb, tbPrefix string) MySqlDB {
	m := &mysqlDB{tbPrefix: tbPrefix}
	m.dial(addr, username, pwd, defaultDb)
	return m
}

func (m *mysqlDB) dial(addr, username, pwd, defaultDb string) {
	// db
	db, err := gorm.Open("mysql", username+":"+pwd+"@tcp("+addr+")/"+defaultDb+"?charset=utf8&parseTime=True&loc=Local")
	errors.Check(err)
	m.db = db
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return m.tbPrefix + "_" + defaultTableName
	}
	//db.LogMode(true)
}

func (m *mysqlDB) OrmDB() *gorm.DB {
	if m.db == nil {
		panic("- DB not be initialized, invoke InitDB at first!")
	}
	return m.db
}

func (m *mysqlDB) TableName(tbName string) string {
	if m.tbPrefix != "" {
		return m.tbPrefix + "_" + tbName
	}
	return tbName
}

func (m *mysqlDB) AutoMigrate(model interface{}) {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	m.OrmDB().AutoMigrate(model)
}

func (m *mysqlDB) Create(model interface{}) bool {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	res := m.OrmDB().Create(model)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *mysqlDB) Find(model interface{}) bool {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	res := m.OrmDB().Where(model).First(model)
	return res.Error == nil
}

func (m *mysqlDB) FindAll(models interface{}) bool {
	res := m.OrmDB().Find(models)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *mysqlDB) FindMany(models interface{}, limit int, orderBy string, whereAndArgs ... interface{}) bool {
	db := m.OrmDB()
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

func (m *mysqlDB) Begin() *gorm.DB {
	return m.OrmDB().Begin()
}

func (m *mysqlDB) Rollback() {
	m.OrmDB().Rollback()
}

func (m *mysqlDB) Commit() {
	m.OrmDB().Rollback()
}

func (m *mysqlDB) Close() {
	if m.db != nil {
		m.db.Close()
	}
}
