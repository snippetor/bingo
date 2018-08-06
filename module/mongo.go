package module

import (
	"gopkg.in/mgo.v2"
	"github.com/snippetor/bingo/errors"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/app"
)

type MongoModule interface {
	Module
	Session() *mgo.Session
	Create(model interface{})
	CreateMany(models interface{})
	FindAll(models interface{})
	FindMany(bson bson.M, models interface{})
	Find(bson bson.M, model interface{})
}

type mongoModule struct {
	app     app.Application
	session *mgo.Session
}

func NewMongoModule(app app.Application, addr, username, pwd, defaultDb string) MongoModule {
	m := &mongoModule{app: app}
	m.dial(addr, username, pwd, defaultDb)
	return m
}

func (m *mongoModule) dial(addr, user, pwd, defaultDb string) {
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
func (m *mongoModule) Session() *mgo.Session {
	return m.session.Copy()
}

func (m *mongoModule) Create(model interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(model))
	err := c.Insert(model)
	errors.Check(err)
}

func (m *mongoModule) CreateMany(models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Insert(models)
	errors.Check(err)
}

func (m *mongoModule) FindAll(models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Find(bson.M{}).All(models)
	errors.Check(err)
}

func (m *mongoModule) FindMany(bson bson.M, models interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.ElementName(models))
	err := c.Find(bson).All(models)
	errors.Check(err)
}

func (m *mongoModule) Find(bson bson.M, model interface{}) {
	session := m.Session()
	defer session.Close()
	c := session.DB("").C(utils.StructName(model))
	err := c.Find(bson).One(model)
	errors.Check(err)
}

func (m *mongoModule) Close() {
	if m.session != nil {
		m.session.Close()
	}
}
