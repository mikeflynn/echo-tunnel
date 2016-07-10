package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TunnelSession struct {
	AWSID     string
	Target    string
	Cmd       string
	Payload   string
	UpdatedAt int64
}

func MongoConnect() (*mgo.Session, *mgo.Collection) {
	mongodb, err := mgo.Dial("localhost")
	if err != nil {
		Debug(err.Error())
	}

	col := mongodb.DB("echo").C("echotunnel")

	return mongodb, col
}

func GetSession(col *mgo.Collection, sessid string) *TunnelSession {
	user := &TunnelSession{}
	err := col.Find(bson.M{"awsid": sessid}).One(&user)
	if err != nil || user.AWSID == "" {
		user.AWSID = sessid
		user.UpdatedAt = time.Now().Unix()
		err = col.Insert(&user)
		if err != nil {
			Debug(err.Error())
		}
	}

	return user
}

func (this *TunnelSession) Update(col *mgo.Collection) error {
	this.UpdatedAt = time.Now().Unix()
	err := col.Update(bson.M{"awsid": this.AWSID}, this)
	if err != nil {
		return err
	}

	return nil
}

func (this *TunnelSession) Delete(col *mgo.Collection) error {
	err := col.Remove(bson.M{"awsid": this.AWSID})
	if err != nil {
		return err
	}

	return nil
}
