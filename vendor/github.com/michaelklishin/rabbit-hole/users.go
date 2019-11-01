package rabbithole

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
)

// HashingAlgorithm represents a hashing algorithm used
// by RabbitMQ's an internal authentication backend.
type HashingAlgorithm string

func (algo HashingAlgorithm) String() string {
	return string(algo)
}

const (
	// HashingAlgorithmSHA256 sets password hashing algorithm to SHA-256.
	HashingAlgorithmSHA256 HashingAlgorithm = "rabbit_password_hashing_sha256"
	// HashingAlgorithmSHA512 sets password hashing algorithm to SHA-512.
	HashingAlgorithmSHA512 HashingAlgorithm = "rabbit_password_hashing_sha512"

	// HashingAlgorithmMD5 provided to support responses that include users created
	// before RabbitMQ 3.6 and other legacy scenarios. This algorithm is
	// deprecated.
	HashingAlgorithmMD5 HashingAlgorithm = "rabbit_password_hashing_md5"
)

// UserInfo represents a user record. Only relevant when internal authentication
// backend is used.
type UserInfo struct {
	Name             string           `json:"name"`
	PasswordHash     string           `json:"password_hash"`
	HashingAlgorithm HashingAlgorithm `json:"hashing_algorithm,omitempty"`
	// Tags control permissions. Built-in tags: administrator, management, policymaker.
	Tags string `json:"tags"`
}

// UserSettings represents properties of a user. Used to create users.
// Tags must be comma-separated.
type UserSettings struct {
	Name string `json:"name"`
	// Tags control permissions. Administrator grants full
	// permissions, management grants management UI and HTTP API
	// access, policymaker grants policy management permissions.
	Tags string `json:"tags"`

	// *never* returned by RabbitMQ. Set by the client
	// to create/update a user. MK.
	Password         string           `json:"password,omitempty"`
	PasswordHash     string           `json:"password_hash,omitempty"`
	HashingAlgorithm HashingAlgorithm `json:"hashing_algorithm,omitempty"`
}

//
// GET /api/users
//

// Example response:
// [{"name":"guest","password_hash":"8LYTIFbVUwi8HuV/dGlp2BYsD1I=","tags":"administrator"}]

// ListUsers returns a list of all users in a cluster.
func (c *Client) ListUsers() (rec []UserInfo, err error) {
	req, err := newGETRequest(c, "users/")
	if err != nil {
		return []UserInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []UserInfo{}, err
	}

	return rec, nil
}

//
// GET /api/users/{name}
//

// GetUser returns information about individual user.
func (c *Client) GetUser(username string) (rec *UserInfo, err error) {
	req, err := newGETRequest(c, "users/"+url.PathEscape(username))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/users/{name}
//

// PutUser updates information about an individual user.
func (c *Client) PutUser(username string, info UserSettings) (res *http.Response, err error) {
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "users/"+url.PathEscape(username), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

// PutUserWithoutPassword creates a passwordless user. Such users can only authenticate
// using an X.509 certificate or another authentication mechanism (or backend) that does not
// use passwords..
func (c *Client) PutUserWithoutPassword(username string, info UserSettings) (res *http.Response, err error) {
	body, err := json.Marshal(UserInfo{Tags: info.Tags})
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "users/"+url.PathEscape(username), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/users/{name}
//

// DeleteUser deletes a user by name.
func (c *Client) DeleteUser(username string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "users/"+url.PathEscape(username), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// Password Hash generation
//

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateSalt generates a password salt. Used to compute password hashes
// when creating or updating user information.
// See https://www.rabbitmq.com/passwords.html#computing-password-hash
// for details.
func GenerateSalt(n int) string {
	bs := make([]byte, n)
	for i := range bs {
		bs[i] = characters[rand.Intn(len(characters))]
	}
	return string(bs)
}

// SaltedPasswordHashSHA256 is used to compute SHA-256 password hashes
// when creating or updating user information.
// See https://www.rabbitmq.com/passwords.html#computing-password-hash
// for details.
func SaltedPasswordHashSHA256(password string) (string, string) {
	salt := GenerateSalt(4)
	hashed := sha256.Sum256([]byte(salt + password))
	return salt, string(hashed[:])
}

// Base64EncodedSaltedPasswordHashSHA256 produces a salted hash value expected by the HTTP API.
// See https://www.rabbitmq.com/passwords.html#computing-password-hash
// for details.
func Base64EncodedSaltedPasswordHashSHA256(password string) string {
	salt, saltedHash := SaltedPasswordHashSHA256(password)
	return base64.URLEncoding.EncodeToString([]byte(salt + saltedHash))
}

// SaltedPasswordHashSHA512 is used to compute SHA-512 password hashes
// when creating or updating user information.
// See https://www.rabbitmq.com/passwords.html#computing-password-hash
// for details.
func SaltedPasswordHashSHA512(password string) (string, string) {
	salt := GenerateSalt(4)
	hashed := sha512.Sum512([]byte(salt + password))
	return salt, string(hashed[:])
}

// Base64EncodedSaltedPasswordHashSHA512 produces a salted hash value expected by the HTTP API.
// See https://www.rabbitmq.com/passwords.html#computing-password-hash
// for details.
func Base64EncodedSaltedPasswordHashSHA512(password string) string {
	salt, saltedHash := SaltedPasswordHashSHA512(password)
	return base64.URLEncoding.EncodeToString([]byte(salt + saltedHash))
}
