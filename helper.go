package main

import "gopkg.in/mgo.v2"
import "github.com/garyburd/redigo/redis"
import "os"
import "fmt"
import "encoding/json"

import "./utils"

type RedisMessage struct {
	Action string
	Key_id string
}

func main () {
	// say hi
	fmt.Printf("hello matej\n")
	
	// create new connection to mongodb
	mongo_session, err := mgo.Dial("mongodb://127.0.0.1")
	
	// check for errors
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}

	redis_session, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	psc := redis.PubSubConn{redis_session}
	psc.Subscribe("git")

	// get the collection
	db := mongo_session.DB("nodegear")

	for {
		switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
				var msg RedisMessage
				err := json.Unmarshal(v.Data, &msg)

				if err != nil {
					fmt.Printf("JSON Parse Error %v", err)
					os.Exit(1)
				}

				fmt.Printf("Hello %s", msg.Action)

				go utils.DeleteKey(db)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Printf("%v", v)
		}
	}
}