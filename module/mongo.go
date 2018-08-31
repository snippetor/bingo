package module

import (
	"gopkg.in/mgo.v2"
	"github.com/snippetor/bingo/errors"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/snippetor/bingo/utils"
)

// mongo
type MongoModule interface {
	Module
	DefaultDB() MongoDB
	DB(name string) MongoDB
}

func NewMongoModule(dbs map[string]MongoDB) MongoModule {
	return &mongoModule{dbs}
}

type mongoModule struct {
	dbs map[string]MongoDB
}

func (m *mongoModule) DefaultDB() MongoDB {
	if db, ok := m.dbs["default"]; ok {
		return db
	}
	for _, v := range m.dbs {
		return v
	}
	return nil
}

func (m *mongoModule) DB(name string) MongoDB {
	if db, ok := m.dbs[name]; ok {
		return db
	}
	return nil
}

func (m *mongoModule) Close() {
	for _, v := range m.dbs {
		v.Close()
	}
}

type MongoDB interface {
	Session() *mgo.Session
	Create(model interface{})
	CreateMany(models interface{})
	FindAll(models interface{})
	FindMany(bson bson.M, models interface{})
	Find(bson bson.M, model interface{})
	Close()
}

func NewMongoDB(addr, username, pwd, defaultDb string) MongoDB {
	m := &mongoDB{}
	m.dial(addr, username, pwd, defaultDb)
	return m
}

type mongoDB struct {
	session *mgo.Session
}

func (m *mongoDB) dial(addr, user, pwd, defaultDb string) {
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	//mongodb://myuser:mypass@localhost:40001,otherhost:40001/
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{addr},
		Timeout:  20 * time.Second,
		Database: defaultDb,
		Username: user,
		Password: pwd,
	})
	errors.Check(err)
	//session.SetMode(mgo.Monotonic, true)
	m.session = session
}

// must close session on transaction finish
func (m *mongoDB) Session() *mgo.Session {
	return m.session.Copy()
}

func (m *mongoDB) Create(model interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(model))
	err := c.Insert(model)
	errors.Check(err)
}

func (m *mongoDB) CreateMany(models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Insert(models)
	errors.Check(err)
}

func (m *mongoDB) FindAll(models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Find(bson.M{}).All(models)
	errors.Check(err)
}

func (m *mongoDB) FindMany(bson bson.M, models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Find(bson).All(models)
	errors.Check(err)
}

func (m *mongoDB) Find(bson bson.M, model interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(model))
	err := c.Find(bson).One(model)
	errors.Check(err)
}

func (m *mongoDB) Close() {
	if m.session != nil {
		m.session.Close()
	}
}
