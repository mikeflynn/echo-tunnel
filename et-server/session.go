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
	Conn      *mgo.Collection
}

func GetSession(sessid string) *TunnelSession {
	mongodb, err := mgo.Dial("localhost")
	if err != nil {
		Debug(err.Error())
	}

	col := mongodb.DB("echo").C("echotunnel")
	defer mongodb.Close()

	user := &TunnelSession{}
	err = col.Find(bson.M{"awsid": sessid}).One(&user)
	if err != nil || user.AWSID == "" {
		user.AWSID = sessid
		user.Conn = col
		user.UpdatedAt = time.Now().Unix()
		err = col.Insert(&user)
		if err != nil {
			Debug(err.Error())
		}
	}

	return user
}

func (this *TunnelSession) Update() error {
	this.UpdatedAt = time.Now().Unix()
	err := this.Conn.Update(bson.M{"awsid": this.AWSID}, this)
	if err != nil {
		return err
	}

	return nil
}

func (this *TunnelSession) Delete() error {
	err := this.Conn.Remove(bson.M{"awsid": this.AWSID})
	if err != nil {
		return err
	}

	return nil
}
