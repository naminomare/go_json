package mongoutil

import (
	"encoding/json"

	"../util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateSession(mongoip string) *mgo.Session {
	session, err := mgo.Dial(mongoip)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

func Find(c *mgo.Collection, query interface{}) []byte {
	//	queryBson := bson.M(query)
	var results []bson.M
	err := c.Find(query).All(&results)
	util.IfnilError(err)

	ret, err := json.Marshal(results)
	util.IfnilError(err)

	return ret
}

func Insert(c *mgo.Collection, query interface{}) {
	err := c.Insert(query)
	util.IfnilError(err)
}

func Update(c *mgo.Collection, findquery interface{}, updatequery interface{}) {
	err := c.Update(findquery, updatequery)
	util.IfnilError(err)
}
