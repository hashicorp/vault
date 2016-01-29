package physical

import (
	"fmt"
	"sort"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

// MongoBackend is a physical backend that stores data
// within MongoDB database.
type MongoBackend struct {
	database        string
	collection      string
	session         *mgo.Session
}

// newMongoBackend constructs a MongoDB backend using the given API client and
// server address and credential for accessing MongoDB database.
func newMongoBackend(conf map[string]string) (Backend, error) {
	// Get the MongoDB credentials to perform read/write operations.
	username, ok := conf["username"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing username")
	}
	password, ok := conf["password"]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get or set MongoDB server address. Defaults to localhost and default port(27017)
	address, ok := conf["address"]
	if !ok {
		address = "127.0.0.1:27017"
	}

	// Get the MongoDB database and collection details.
	database, ok := conf["database"]
	if !ok {
		database = "vault"
	}
	collection, ok := conf["collection"]
	if !ok {
		collection = "vault_collection"
	}

	// Create MongoDB handle for the database.
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	credential := mgo.Credential{
		Username:   username,
		Password: 	password,
	}

	err = session.Login(&credential)
	if err != nil {
		return nil, fmt.Errorf("failed to login to mongodb: %v", err)
	}

	// Setup the backend.
	m := &MongoBackend{
		database:    database,
		collection:	 collection,
		session:     session,
	}

	return m, nil
}

// Put is used to insert or update an entry.
func (m *MongoBackend) Put(entry *Entry) error {

	_, err := m.session.DB(m.database).C(m.collection).Upsert(bson.M{"key": entry.Key}, bson.M{"key": entry.Key, "value": entry.Value})
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *MongoBackend) Get(key string) (*Entry, error) {

	type keyValue struct {
	Key string
	Value []byte
	}

	result := keyValue{}
	err := m.session.DB(m.database).C(m.collection).Find(bson.M{"key": key}).One(&result)
	if err != nil && err.Error() != "not found" {
		fmt.Println("result.key:")
		fmt.Println(result.Key)
		return nil, err
	}

	if err != nil && err.Error() == "not found" {
		return nil, nil
	}

	ent := &Entry{
		Key:   key,
		Value: result.Value,
	}

	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *MongoBackend) Delete(key string) error {
	err := m.session.DB(m.database).C(m.collection).Remove(bson.M{"key": key})
	if err != nil && err.Error() != "not found" {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *MongoBackend) List(prefix string) ([]string, error) {

	likePrefix := "^"+prefix
	iter := m.session.DB(m.database).C(m.collection).Find(bson.M{ "key": bson.M{"$regex": bson.RegEx{likePrefix,""}}}).Iter()

	var keys []string

	type keyValue struct {
	Key 	string
	Value []byte
	}

	result := keyValue{}
	var insertedFolder bool = false
	regCheck, _ := regexp.Compile("(.)*/$")
	folderCheck, _ := regexp.Compile("/")

	if regCheck.MatchString(prefix){
		// ending with /
		for iter.Next(&result) {
			noPrefregex,_ := regexp.Compile("^"+prefix)
			noPrefixKey := noPrefregex.ReplaceAllString(result.Key ,"")
			noPostregex,_ := regexp.Compile("/(.)*")
			noPostfixKey := noPostregex.ReplaceAllString(noPrefixKey,"/")
			found := false
			for i := range keys {
			 if(keys[i] == noPostfixKey){found = true}
	 		}
			if found == false{
				keys = append(keys, noPostfixKey)
			}

			if err := iter.Close(); err != nil {
				return nil, err
			}
		}
	}else{
		// not ending with /
		for iter.Next(&result) {
			noPrefregex,_ := regexp.Compile("^"+prefix)
			noPrefixKey := noPrefregex.ReplaceAllString(result.Key ,"")
				if insertedFolder == false && folderCheck.MatchString(noPrefixKey){
					noPostregex,_ := regexp.Compile("/(.)*")
					noPostfixKey := noPostregex.ReplaceAllString(noPrefixKey,"/")
					keys = append(keys, noPostfixKey)
					insertedFolder = true
				}else{
						if folderCheck.MatchString(noPrefixKey) == false{
							keys = append(keys, noPrefixKey)
						}
				}
				if err := iter.Close(); err != nil {
					return nil, err
				}
		}
	}

	sort.Strings(keys)
	return keys, nil
}
