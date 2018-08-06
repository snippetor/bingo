package mvc

import (
	"github.com/snippetor/bingo/app"
	"gopkg.in/mgo.v2"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/errors"
	"gopkg.in/mgo.v2/bson"
)

type NoSqlModel interface {
	Init(app app.Application, self NoSqlModel)
	App() app.Application
	Session() *mgo.Session
	Sync()
	Del()
}

type MongoModel struct {
	Id   bson.ObjectId `bson:"_id"`
	app  app.Application
	self NoSqlModel
}

func (m *MongoModel) Init(app app.Application, self NoSqlModel) {
	m.app = app
	m.self = self
}

func (m *MongoModel) App() app.Application {
	return m.app
}

// must close session on transaction finish
func (m *MongoModel) Session() *mgo.Session {
	return m.App().Mongo().Session()
}

func (m *MongoModel) Sync() {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(m.self))
	err := c.UpdateId(m.Id, m.self)
	errors.Check(err)
}

func (m *MongoModel) Del() {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(m.self))
	err := c.RemoveId(m.Id)
	errors.Check(err)
}
