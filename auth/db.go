package auth

import (
	"log"

	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/TerrexTech/uuuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// DBIConfig is the configuration for the authDB.
type DBIConfig struct {
	Hosts               []string
	Username            string
	Password            string
	TimeoutMilliseconds uint32
	Database            string
	Collection          string
}

// DBI is the Database-interface for authentication.
// This fetches/writes data to/from database for auth-actions such as
// login, registeration etc.
type DBI interface {
	Collection() *mongo.Collection
	UserByUUID(uid uuuid.UUID) (*User, error)
	Login(user *User) (*User, error)
}

// DB is the implementation for dbI.
// dbI is the Database-interface for authentication.
// It fetches/writes data to/from database for auth-actions such as
// login, registeration etc.
type DB struct {
	collection *mongo.Collection
}

// EnsureAuthDB exists ensures that the required Database and Collection exists before
// auth-operations can be done on them. It creates Database/Collection if they don't exist.
func EnsureAuthDB(dbConfig DBIConfig) (*DB, error) {
	config := mongo.ClientConfig{
		Hosts:               dbConfig.Hosts,
		Username:            dbConfig.Username,
		Password:            dbConfig.Password,
		TimeoutMilliseconds: dbConfig.TimeoutMilliseconds,
	}

	client, err := mongo.NewClient(config)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}

	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}

	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name: "username",
				},
			},
			IsUnique: true,
			Name:     "username_index",
		},
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name:        "version",
					IsDescOrder: true,
				},
			},
			IsUnique: true,
			Name:     "version_index",
		},
	}

	// ====> Create New Collection
	collConfig := &mongo.Collection{
		Connection:   conn,
		Database:     dbConfig.Database,
		Name:         dbConfig.Collection,
		SchemaStruct: &User{},
		Indexes:      indexConfigs,
	}
	c, err := mongo.EnsureCollection(collConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}
	return &DB{
		collection: c,
	}, nil
}

// UserByUUID gets the User from DB using specified UUID.
// An error is returned if no user is found.
func (d *DB) UserByUUID(uid uuuid.UUID) (*User, error) {
	user := &User{
		UUID: uid,
	}

	findResults, err := d.collection.Find(user)
	if err != nil {
		err = errors.Wrap(err, "UserByUUID: Error getting user from Database")
		return nil, err
	}
	if len(findResults) == 0 {
		return nil, errors.New("UserByUUID: User not found")
	}

	resultUser := findResults[0].(*User)
	return resultUser, nil
}

// Login authenticates the provided User.
// An error is returned if Authentication fails.
func (d *DB) Login(user *User) (*User, error) {
	authUser := &User{
		Username: user.Username,
	}

	findResults, err := d.collection.Find(authUser)
	if err != nil {
		err = errors.Wrap(err, "Login: Error getting user from Database")
		return nil, err
	}
	if len(findResults) == 0 {
		log.Println("===========================")
		return nil, errors.New("Login: Invalid Credentials")
	}

	newUser := findResults[0].(*User)
	passErr := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(user.Password))
	// log.Println("%+v", user)
	if passErr != nil {
		log.Println("++++++++++++++++++++++++++++")
		log.Println(passErr)
		return nil, errors.New("Login: Invalid Credentials")
	}

	return newUser, nil
}

// Collection returns the currrent MongoDB collection being used for user-auth operations.
func (d *DB) Collection() *mongo.Collection {
	return d.collection
}
