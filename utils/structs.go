package utils

import "gopkg.in/mgo.v2/bson"
import "time"

type RedisMessage struct {
	Action	string
	Key_id	string
}

type RSAKey struct {
	Id				bson.ObjectId `bson:"_id"`
	Created		time.Time `bson:"created"`
	Deleted		bool `bson:"deleted"`
	User			bson.ObjectId `bson:"user"`
	Name			string `bson:"name"`
	nameLowercase string `bson:"nameLowercase"`
	System_key	bool `bson:"system_key"`
	Private_key string `bson:"private_key"`
	Public_key	string `bson:"public_key"`
	Installed	bool `bson:"installed"`
	Installing 	bool `bson:"installing"`
}