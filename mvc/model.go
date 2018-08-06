package mvc

import (
	"github.com/snippetor/bingo/app"
	"github.com/jinzhu/gorm"
	"github.com/snippetor/bingo/errors"
	"github.com/snippetor/bingo/codec"
)

type OrmModel interface {
	Init(app app.Application, db *gorm.DB, self OrmModel)
	App() app.Application
	DB() *gorm.DB
	Sync(cols ...interface{})
	SyncInTx(tx *gorm.DB, cols ...interface{})
	Del()
	DelInTx(tx *gorm.DB)
	FieldToString(f interface{}) string
	FieldFromString(s string, f interface{})
}

type MysqlOrmModel struct {
	app  app.Application
	db   *gorm.DB
	self OrmModel
}

func (m *MysqlOrmModel) Init(app app.Application, db *gorm.DB, self OrmModel) {
	m.app = app
	m.db = db
	m.self = self
}

func (m *MysqlOrmModel) App() app.Application {
	return m.app
}

func (m *MysqlOrmModel) DB() *gorm.DB {
	if m.db == nil {
		panic("-- db is nil in OrmModel --")
	}
	return m.db
}

// 更新到数据库
func (m *MysqlOrmModel) Sync(cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		errors.Check(m.DB().Model(m.self).UpdateColumn(cols).Error)
	} else {
		errors.Check(m.DB().Model(m.self).Updates(m.self).Error)
	}
}

// 更新到数据库
func (m *MysqlOrmModel) SyncInTx(tx *gorm.DB, cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		errors.Check(tx.Model(m.self).UpdateColumn(cols).Error)
	} else {
		errors.Check(tx.Model(m.self).Updates(m.self).Error)
	}
}

// 从数据库移除，ID必须存在
func (m *MysqlOrmModel) Del() {
	errors.Check(m.DB().Delete(m.self).Error)
}

// 从数据库移除，ID必须存在
func (m *MysqlOrmModel) DelInTx(tx *gorm.DB) {
	errors.Check(tx.Delete(m.self).Error)
}

func (m *MysqlOrmModel) FieldToString(f interface{}) string {
	return string(codec.JsonCodec.Marshal(f))
}

func (m *MysqlOrmModel) FieldFromString(s string, f interface{}) {
	codec.JsonCodec.Unmarshal([]byte(s), f)
}
