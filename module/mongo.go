package module

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/snippetor/bingo/errors"
)

type MongoModule interface {
	Dial(addr, user, pwd, db string)
	Session() *mgo.Session
}

type mongoModule struct {
	session *mgo.Session
}

func (m *mongoModule) Dial(addr, user, pwd, db string) {
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	//mongodb://myuser:mypass@localhost:40001,otherhost:40001/
	var session *mgo.Session
	var err error
	if user == "" {
		session, err = mgo.Dial(fmt.Sprintf("mongodb://@%s/%s", addr, db))
	} else {
		session, err = mgo.Dial(fmt.Sprintf("mongodb://%s:%s@%s/%s", addr, user, pwd, db))
	}
	errors.Check(err)
	m.session = session
}

// must close session on transaction finish
func (m *mongoModule) Session() *mgo.Session {
	return m.session.Copy()
}
