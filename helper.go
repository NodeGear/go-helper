package main

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "os"
import "fmt"

type msg struct {
	Id    bson.ObjectId `bson:"_id"`
	Msg   string        `bson:"name"`
}

func main () {
	// say hi
	fmt.Printf("hello matej\n")
	
	// create new connection to mongodb
	session, err := mgo.Dial("mongodb://127.0.0.1")
	
	// check for errors
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}

	// get the collection
	db := session.DB("nodegear")
	users := db.C("users")

	var updatedmsg msg
	err = users.Find(bson.M{}).One(&updatedmsg)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
		os.Exit(1)
	}

	fmt.Printf("Found document: %+v\n", updatedmsg)
}