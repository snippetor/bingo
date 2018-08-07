package mvc

import (
	"github.com/snippetor/bingo/app"
	"gopkg.in/mgo.v2"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/errors"
	"gopkg.in/mgo.v2/bson"
)

type MongoModel interface {
	Init(app app.Application, self MongoModel)
	App() app.Application
	Session() *mgo.Session
	Sync()
	Del()
}

type BaseMongoModel struct {
	Id   bson.ObjectId `bson:"_id"`
	app  app.Application
	self MongoModel
}

func (m *BaseMongoModel) Init(app app.Application, self MongoModel) {
	m.app = app
	m.self = self
}

func (m *BaseMongoModel) App() app.Application {
	return m.app
}

// must close session on transaction finish
func (m *BaseMongoModel) Session() *mgo.Session {
	return m.App().Mongo().Session()
}

func (m *BaseMongoModel) Sync() {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(m.self))
	err := c.UpdateId(m.Id, m.self)
	errors.Check(err)
}

func (m *BaseMongoModel) Del() {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(m.self))
	err := c.RemoveId(m.Id)
	errors.Check(err)
}
