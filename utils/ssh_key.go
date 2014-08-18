package utils

import "fmt"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "os"
import "time"

type msg struct {
	Id    bson.ObjectId `bson:"_id"`
	Msg   string        `bson:"name"`
}

func DeleteKey (db *mgo.Database) {
	fmt.Printf("hello matej\n")

	users := db.C("users")

	var updatedmsg msg
	err := users.Find(bson.M{}).One(&updatedmsg)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n", err)
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)

	fmt.Printf("Found document: %+v\n", updatedmsg)
}