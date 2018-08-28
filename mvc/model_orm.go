package mvc

import (
	"github.com/jinzhu/gorm"
	"github.com/snippetor/bingo/errors"
	"github.com/snippetor/bingo/codec"
	"reflect"
)

var ormModelType = reflect.TypeOf((*MysqlOrmModel)(nil)).Elem()

func IsOrmModel(i interface{}) bool {
	return reflect.TypeOf(i).Implements(ormModelType)
}

type MysqlOrmModel interface {
	Init(db *gorm.DB, self interface{})
	DB() *gorm.DB
	Sync(cols ...interface{})
	SyncInTx(tx *gorm.DB, cols ...interface{})
	Del()
	DelInTx(tx *gorm.DB)
	FieldToString(f interface{}) string
	FieldFromString(s string, f interface{})
}

type BaseMysqlOrmModel struct {
	db   *gorm.DB
	self interface{}
	Id   uint32 `gorm:"primary_key"`
}

func (m *BaseMysqlOrmModel) Init(db *gorm.DB, self interface{}) {
	m.db = db
	m.self = self
}

func (m *BaseMysqlOrmModel) DB() *gorm.DB {
	return m.db
}

// 更新到数据库
func (m *BaseMysqlOrmModel) Sync(cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		errors.Check(m.DB().Model(m.self).UpdateColumn(cols).Error)
	} else {
		errors.Check(m.DB().Model(m.self).Updates(m.self).Error)
	}
}

// 更新到数据库
func (m *BaseMysqlOrmModel) SyncInTx(tx *gorm.DB, cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		errors.Check(tx.Model(m.self).UpdateColumn(cols).Error)
	} else {
		errors.Check(tx.Model(m.self).Updates(m.self).Error)
	}
}

// 从数据库移除，ID必须存在
func (m *BaseMysqlOrmModel) Del() {
	errors.Check(m.DB().Delete(m.self).Error)
}

// 从数据库移除，ID必须存在
func (m *BaseMysqlOrmModel) DelInTx(tx *gorm.DB) {
	errors.Check(tx.Delete(m.self).Error)
}

func (m *BaseMysqlOrmModel) FieldToString(f interface{}) string {
	bs, err := codec.JsonCodec.Marshal(f)
	errors.Check(err)
	return string(bs)
}

func (m *BaseMysqlOrmModel) FieldFromString(s string, f interface{}) {
	codec.JsonCodec.Unmarshal([]byte(s), f)
}
