package auth

import (
	"encoding/json"

	"github.com/TerrexTech/uuuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/pkg/errors"
)

// User represents a system-user.
type User struct {
	ID        objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UUID      uuuid.UUID        `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email     string            `bson:"email,omitempty" json:"email,omitempty"`
	FirstName string            `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string            `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Username  string            `bson:"username,omitempty" json:"username,omitempty"`
	Password  string            `bson:"password,omitempty" json:"password,omitempty"`
	Role      string            `bson:"role,omitempty" json:"role,omitempty"`
	Version   int64             `bson:"version,omitempty" json:"version,omitempty"`
}

// marshalUser is a simplified User, for convenient marshalling/unmarshalling operations
type marshalUser struct {
	ID        objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UUID      string            `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email     string            `bson:"email,omitempty" json:"email,omitempty"`
	FirstName string            `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string            `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Username  string            `bson:"username,omitempty" json:"username,omitempty"`
	Password  string            `bson:"password,omitempty" json:"password,omitempty"`
	Role      string            `bson:"role,omitempty" json:"role,omitempty"`
	Version   int64             `bson:"version,omitempty" json:"version,omitempty"`
}

// MarshalBSON converts the User to its BSON representation.
func (u User) MarshalBSON() ([]byte, error) {
	mu := &marshalUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		Role:      u.Role,
	}

	if u.UUID.String() != (uuuid.UUID{}).String() {
		mu.UUID = u.UUID.String()
	}
	return bson.Marshal(mu)
}

// MarshalJSON converts the User to its JSON representation.
func (u *User) MarshalJSON() ([]byte, error) {
	// No password here since JSON is for external use, while BSON is used internally
	mu := &map[string]interface{}{
		"_id":        u.ID.Hex(),
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"username":   u.Username,
		"role":       u.Role,
		"uuid":       u.UUID.String(),
	}
	return json.Marshal(mu)
}

// UnmarshalBSON converts the BSON representation of User back to User-struct.
func (u *User) UnmarshalBSON(in []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		return err
	}
	u.ID = m["_id"].(objectid.ObjectID)

	u.UUID, err = uuuid.FromString(m["uuid"].(string))
	if err != nil {
		err = errors.Wrap(err, "Error parsing UUID for user")
		return err
	}
	if m["email"] != nil {
		u.Email = m["email"].(string)
	}
	if m["first_name"] != nil {
		u.FirstName = m["first_name"].(string)
	}
	if m["last_name"] != nil {
		u.LastName = m["last_name"].(string)
	}
	if m["username"] != nil {
		u.Username = m["username"].(string)
	}
	if m["password"] != nil {
		u.Password = m["password"].(string)
	}
	if m["role"] != nil {
		u.Role = m["role"].(string)
	}

	return nil
}
