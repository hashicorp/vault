package mongodb

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"gopkg.in/mgo.v2"
)

const mongoDBTypeName = "mongodb"

// MongoDB is an implementation of Database interface
type MongoDB struct {
	connutil.ConnectionProducer
	credsutil.CredentialsProducer
}

// New returns a new MongoDB instance
func New() (interface{}, error) {
	connProducer := &connutil.MongoDBConnectionProducer{}
	connProducer.Type = mongoDBTypeName

	credsProducer := &credsutil.MongoDBCredentialsProducer{}

	dbType := &MongoDB{
		ConnectionProducer:  connProducer,
		CredentialsProducer: credsProducer,
	}
	return dbType, nil
}

// Run instantiates a MongoDB object, and runs the RPC server for the plugin
func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}

	plugins.Serve(dbType.(*MongoDB), apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (m *MongoDB) Type() (string, error) {
	return mongoDBTypeName, nil
}

func (m *MongoDB) getConnection() (*mgo.Session, error) {
	session, err := m.Connection()
	if err != nil {
		return nil, err
	}

	return session.(*mgo.Session), nil
}
