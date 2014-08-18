package main

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "github.com/garyburd/redigo/redis"
import "os"
import "fmt"
import "encoding/json"

import "./utils"

func main () {
	// create new connection to mongodb
	mongo_session, err := mgo.Dial("mongodb://127.0.0.1")
	
	// check for errors
	if err != nil {
		fmt.Printf("MongoDB Connection Error: %v\n", err)
		os.Exit(1)
	}

	pubsub_session, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Redis Connection Error: %v\n", err)
		os.Exit(1)
	}

	redis_session, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Redis Connection Error: %v\n", err)
		os.Exit(1)
	}

	pub_sub := redis.PubSubConn{pubsub_session}
	pub_sub.Subscribe("git")

	// get the collection
	db := mongo_session.DB("nodegear")
	keys := db.C("rsakeys")

	fmt.Println("Listening to Redis commands (namespace git)")

	for {
		switch v := pub_sub.Receive().(type) {
			case redis.Message:
				go dispatch(v, redis_session, keys)
			case error:
				fmt.Printf("Redis PubSub Error: %v", v)
		}
	}
}

func dispatch (v redis.Message, redis_session redis.Conn, keys *mgo.Collection) {
	// Parse the message. it is JSON formatted
	var msg utils.RedisMessage

	err := json.Unmarshal(v.Data, &msg)
	if err != nil {
		fmt.Printf("JSON Parse Error %v", err)
		return
	}

	var key utils.RSAKey
	query := bson.M{"_id": bson.ObjectIdHex(msg.Key_id), "deleted": false}

	err = keys.Find(query).One(&key)
	if err != nil {
		fmt.Printf("Error finding rsakey document: %v\n", err)
		return
	}
	
	if key.Installing == true {
		redis_session.Do("PUBLISH", "git:install", key.Id.Hex()+"|Already Installing")
		return
	}

	err = keys.UpdateId(key.Id, bson.M{"$set": bson.M{"installing": true}})
	if err != nil {
		fmt.Printf("Update Failed: %v", err)
		return
	}

	status := ""
	switch msg.Action {
		case "createSystemKey":
			status = utils.CreateSystemKey(key, keys, msg)
		case "verifyKey":
			status = utils.VerifyKey(key, keys, msg)
	}

	err = keys.UpdateId(key.Id, bson.M{"$set": bson.M{"installing": false}})
	if err != nil {
		fmt.Printf("Could not save collection %v", err)
		status = "System Error"
		return
	}

	redis_session.Do("PUBLISH", "git:install", key.Id.Hex()+"|"+status)
}