package utils

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "fmt"
import "os/exec"
import "os"
import "time"
import "strconv"
import "io/ioutil"

func CreateSystemKey (key RSAKey, keys *mgo.Collection, msg RedisMessage) string {
	key_file := "/tmp/ng_key_"
	key_file += key.Id.Hex()
	key_file += "-"
	key_file += strconv.Itoa(int(time.Now().Unix()))
	
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-C", "nodegear", "-q", "-f", key_file, "-N", "")

	err := cmd.Start()
	if err != nil {
		return "Creation Failed"
	}

	err = cmd.Wait()
	if err != nil {
		return "Creation Failed"
	}

	private_key, err := ioutil.ReadFile(key_file)
	if err != nil {
		return "Creation Failed"
	}

	defer unlink(key_file)
	
	public_key, err := ioutil.ReadFile(key_file+".pub")
	if err != nil {
		return "Creation Failed"
	}

	defer unlink(key_file+".pub")

	updateFieldsSet := bson.M{"installed": true, "public_key": string(public_key), "private_key": string(private_key)}

	updateFields := bson.M{"$set": updateFieldsSet}

	err = keys.UpdateId(key.Id, updateFields)
	if err != nil {
		return "Server Error"
	}

	return "Installation Finished"
}

func unlink (key_file string) {
	os.Remove(key_file)
}

func VerifyKey (key RSAKey, keys *mgo.Collection, msg RedisMessage) string {
	pkey := []byte(key.Public_key)

	key_file := "/tmp/ng_pub_valid_"
	key_file += key.Id.Hex()
	key_file += "-"
	key_file += strconv.Itoa(int(time.Now().Unix()))
	key_file += ".pub"

	err := ioutil.WriteFile(key_file, pkey, 0640)
	if err != nil {
		return "Verification Failed"
	}

	defer unlink(key_file)

	cmd := exec.Command("ssh-keygen", "-lf", key_file)

	err = cmd.Start()
	if err != nil {
		return "Verification Failed"
	}

	err = cmd.Wait()
	if err != nil {
		return "Verification Failed"
	}

	err = keys.UpdateId(key.Id, bson.M{"$set": bson.M{"installed": true}})
	if err != nil {
		return "Server Error"
	}

	return "Installation Finished"
}