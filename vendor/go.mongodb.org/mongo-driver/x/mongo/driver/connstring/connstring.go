// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package connstring // import "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/internal"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/dns"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// Parse parses the provided uri and returns a URI object.
func Parse(s string) (ConnString, error) {
	p := parser{dnsResolver: dns.DefaultResolver}
	err := p.parse(s)
	if err != nil {
		err = internal.WrapErrorf(err, "error parsing uri")
	}
	return p.ConnString, err
}

// ConnString represents a connection string to mongodb.
type ConnString struct {
	Original                           string
	AppName                            string
	AuthMechanism                      string
	AuthMechanismProperties            map[string]string
	AuthSource                         string
	Compressors                        []string
	Connect                            ConnectMode
	ConnectSet                         bool
	ConnectTimeout                     time.Duration
	ConnectTimeoutSet                  bool
	Database                           string
	HeartbeatInterval                  time.Duration
	HeartbeatIntervalSet               bool
	Hosts                              []string
	J                                  bool
	JSet                               bool
	LocalThreshold                     time.Duration
	LocalThresholdSet                  bool
	MaxConnIdleTime                    time.Duration
	MaxConnIdleTimeSet                 bool
	MaxPoolSize                        uint64
	MaxPoolSizeSet                     bool
	MinPoolSize                        uint64
	MinPoolSizeSet                     bool
	Password                           string
	PasswordSet                        bool
	ReadConcernLevel                   string
	ReadPreference                     string
	ReadPreferenceTagSets              []map[string]string
	RetryWrites                        bool
	RetryWritesSet                     bool
	RetryReads                         bool
	RetryReadsSet                      bool
	MaxStaleness                       time.Duration
	MaxStalenessSet                    bool
	ReplicaSet                         string
	Scheme                             string
	ServerSelectionTimeout             time.Duration
	ServerSelectionTimeoutSet          bool
	SocketTimeout                      time.Duration
	SocketTimeoutSet                   bool
	SSL                                bool
	SSLSet                             bool
	SSLClientCertificateKeyFile        string
	SSLClientCertificateKeyFileSet     bool
	SSLClientCertificateKeyPassword    func() string
	SSLClientCertificateKeyPasswordSet bool
	SSLInsecure                        bool
	SSLInsecureSet                     bool
	SSLCaFile                          string
	SSLCaFileSet                       bool
	WString                            string
	WNumber                            int
	WNumberSet                         bool
	Username                           string
	ZlibLevel                          int
	ZlibLevelSet                       bool
	ZstdLevel                          int
	ZstdLevelSet                       bool

	WTimeout              time.Duration
	WTimeoutSet           bool
	WTimeoutSetFromOption bool

	Options        map[string][]string
	UnknownOptions map[string][]string
}

func (u *ConnString) String() string {
	return u.Original
}

// ConnectMode informs the driver on how to connect
// to the server.
type ConnectMode uint8

// ConnectMode constants.
const (
	AutoConnect ConnectMode = iota
	SingleConnect
)

// Scheme constants
const (
	SchemeMongoDB    = "mongodb"
	SchemeMongoDBSRV = "mongodb+srv"
)

type parser struct {
	ConnString

	dnsResolver *dns.Resolver
	tlsssl      *bool // used to determine if tls and ssl options are both specified and set differently.
}

func (p *parser) parse(original string) error {
	p.Original = original
	uri := original

	var err error
	if strings.HasPrefix(uri, SchemeMongoDBSRV+"://") {
		p.Scheme = SchemeMongoDBSRV
		// remove the scheme
		uri = uri[len(SchemeMongoDBSRV)+3:]
	} else if strings.HasPrefix(uri, SchemeMongoDB+"://") {
		p.Scheme = SchemeMongoDB
		// remove the scheme
		uri = uri[len(SchemeMongoDB)+3:]
	} else {
		return fmt.Errorf("scheme must be \"mongodb\" or \"mongodb+srv\"")
	}

	if idx := strings.Index(uri, "@"); idx != -1 {
		userInfo := uri[:idx]
		uri = uri[idx+1:]

		username := userInfo
		var password string

		if idx := strings.Index(userInfo, ":"); idx != -1 {
			username = userInfo[:idx]
			password = userInfo[idx+1:]
			p.PasswordSet = true
		}

		if len(username) > 1 {
			if strings.Contains(username, "/") {
				return fmt.Errorf("unescaped slash in username")
			}
		}

		p.Username, err = url.QueryUnescape(username)
		if err != nil {
			return internal.WrapErrorf(err, "invalid username")
		}
		if len(password) > 1 {
			if strings.Contains(password, ":") {
				return fmt.Errorf("unescaped colon in password")
			}
			if strings.Contains(password, "/") {
				return fmt.Errorf("unescaped slash in password")
			}
			p.Password, err = url.QueryUnescape(password)
			if err != nil {
				return internal.WrapErrorf(err, "invalid password")
			}
		}
	}

	// fetch the hosts field
	hosts := uri
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 {
		if uri[idx] == '@' {
			return fmt.Errorf("unescaped @ sign in user info")
		}
		if uri[idx] == '?' {
			return fmt.Errorf("must have a / before the query ?")
		}
		hosts = uri[:idx]
	}

	var connectionArgsFromTXT []string
	parsedHosts := strings.Split(hosts, ",")

	if p.Scheme == SchemeMongoDBSRV {
		parsedHosts, err = p.dnsResolver.ParseHosts(hosts, true)
		if err != nil {
			return err
		}
		connectionArgsFromTXT, err = p.dnsResolver.GetConnectionArgsFromTXT(hosts)
		if err != nil {
			return err
		}

		// SSL is enabled by default for SRV, but can be manually disabled with "ssl=false".
		p.SSL = true
		p.SSLSet = true
	}

	for _, host := range parsedHosts {
		err = p.addHost(host)
		if err != nil {
			return internal.WrapErrorf(err, "invalid host \"%s\"", host)
		}
	}
	if len(p.Hosts) == 0 {
		return fmt.Errorf("must have at least 1 host")
	}

	uri = uri[len(hosts):]

	extractedDatabase, err := extractDatabaseFromURI(uri)
	if err != nil {
		return err
	}

	uri = extractedDatabase.uri
	p.Database = extractedDatabase.db

	connectionArgsFromQueryString, err := extractQueryArgsFromURI(uri)
	connectionArgPairs := append(connectionArgsFromTXT, connectionArgsFromQueryString...)

	for _, pair := range connectionArgPairs {
		err = p.addOption(pair)
		if err != nil {
			return err
		}
	}

	err = p.setDefaultAuthParams(extractedDatabase.db)
	if err != nil {
		return err
	}

	err = p.validateAuth()
	if err != nil {
		return err
	}

	// Check for invalid write concern (i.e. w=0 and j=true)
	if p.WNumberSet && p.WNumber == 0 && p.JSet && p.J {
		return writeconcern.ErrInconsistent
	}

	// If WTimeout was set from manual options passed in, set WTImeoutSet to true.
	if p.WTimeoutSetFromOption {
		p.WTimeoutSet = true
	}

	return nil
}

func (p *parser) setDefaultAuthParams(dbName string) error {
	switch strings.ToLower(p.AuthMechanism) {
	case "plain":
		if p.AuthSource == "" {
			p.AuthSource = dbName
			if p.AuthSource == "" {
				p.AuthSource = "$external"
			}
		}
	case "gssapi":
		if p.AuthMechanismProperties == nil {
			p.AuthMechanismProperties = map[string]string{
				"SERVICE_NAME": "mongodb",
			}
		} else if v, ok := p.AuthMechanismProperties["SERVICE_NAME"]; !ok || v == "" {
			p.AuthMechanismProperties["SERVICE_NAME"] = "mongodb"
		}
		fallthrough
	case "mongodb-x509":
		if p.AuthSource == "" {
			p.AuthSource = "$external"
		} else if p.AuthSource != "$external" {
			return fmt.Errorf("auth source must be $external")
		}
	case "mongodb-cr":
		fallthrough
	case "scram-sha-1":
		fallthrough
	case "scram-sha-256":
		if p.AuthSource == "" {
			p.AuthSource = dbName
			if p.AuthSource == "" {
				p.AuthSource = "admin"
			}
		}
	case "":
		if p.AuthSource == "" && (p.AuthMechanismProperties != nil || p.Username != "" || p.PasswordSet) {
			p.AuthSource = dbName
			if p.AuthSource == "" {
				p.AuthSource = "admin"
			}
		}
	default:
		return fmt.Errorf("invalid auth mechanism")
	}
	return nil
}

func (p *parser) validateAuth() error {
	switch strings.ToLower(p.AuthMechanism) {
	case "mongodb-cr":
		if p.Username == "" {
			return fmt.Errorf("username required for MONGO-CR")
		}
		if p.Password == "" {
			return fmt.Errorf("password required for MONGO-CR")
		}
		if p.AuthMechanismProperties != nil {
			return fmt.Errorf("MONGO-CR cannot have mechanism properties")
		}
	case "mongodb-x509":
		if p.Password != "" {
			return fmt.Errorf("password cannot be specified for MONGO-X509")
		}
		if p.AuthMechanismProperties != nil {
			return fmt.Errorf("MONGO-X509 cannot have mechanism properties")
		}
	case "gssapi":
		if p.Username == "" {
			return fmt.Errorf("username required for GSSAPI")
		}
		for k := range p.AuthMechanismProperties {
			if k != "SERVICE_NAME" && k != "CANONICALIZE_HOST_NAME" && k != "SERVICE_REALM" {
				return fmt.Errorf("invalid auth property for GSSAPI")
			}
		}
	case "plain":
		if p.Username == "" {
			return fmt.Errorf("username required for PLAIN")
		}
		if p.Password == "" {
			return fmt.Errorf("password required for PLAIN")
		}
		if p.AuthMechanismProperties != nil {
			return fmt.Errorf("PLAIN cannot have mechanism properties")
		}
	case "scram-sha-1":
		if p.Username == "" {
			return fmt.Errorf("username required for SCRAM-SHA-1")
		}
		if p.Password == "" {
			return fmt.Errorf("password required for SCRAM-SHA-1")
		}
		if p.AuthMechanismProperties != nil {
			return fmt.Errorf("SCRAM-SHA-1 cannot have mechanism properties")
		}
	case "scram-sha-256":
		if p.Username == "" {
			return fmt.Errorf("username required for SCRAM-SHA-256")
		}
		if p.Password == "" {
			return fmt.Errorf("password required for SCRAM-SHA-256")
		}
		if p.AuthMechanismProperties != nil {
			return fmt.Errorf("SCRAM-SHA-256 cannot have mechanism properties")
		}
	case "":
		if p.Username == "" && p.AuthSource != "" {
			return fmt.Errorf("authsource without username is invalid")
		}
	default:
		return fmt.Errorf("invalid auth mechanism")
	}
	return nil
}

func (p *parser) addHost(host string) error {
	if host == "" {
		return nil
	}
	host, err := url.QueryUnescape(host)
	if err != nil {
		return internal.WrapErrorf(err, "invalid host \"%s\"", host)
	}

	_, port, err := net.SplitHostPort(host)
	// this is unfortunate that SplitHostPort actually requires
	// a port to exist.
	if err != nil {
		if addrError, ok := err.(*net.AddrError); !ok || addrError.Err != "missing port in address" {
			return err
		}
	}

	if port != "" {
		d, err := strconv.Atoi(port)
		if err != nil {
			return internal.WrapErrorf(err, "port must be an integer")
		}
		if d <= 0 || d >= 65536 {
			return fmt.Errorf("port must be in the range [1, 65535]")
		}
	}
	p.Hosts = append(p.Hosts, host)
	return nil
}

func (p *parser) addOption(pair string) error {
	kv := strings.SplitN(pair, "=", 2)
	if len(kv) != 2 || kv[0] == "" {
		return fmt.Errorf("invalid option")
	}

	key, err := url.QueryUnescape(kv[0])
	if err != nil {
		return internal.WrapErrorf(err, "invalid option key \"%s\"", kv[0])
	}

	value, err := url.QueryUnescape(kv[1])
	if err != nil {
		return internal.WrapErrorf(err, "invalid option value \"%s\"", kv[1])
	}

	lowerKey := strings.ToLower(key)
	switch lowerKey {
	case "appname":
		p.AppName = value
	case "authmechanism":
		p.AuthMechanism = value
	case "authmechanismproperties":
		p.AuthMechanismProperties = make(map[string]string)
		pairs := strings.Split(value, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) != 2 || kv[0] == "" {
				return fmt.Errorf("invalid authMechanism property")
			}
			p.AuthMechanismProperties[kv[0]] = kv[1]
		}
	case "authsource":
		p.AuthSource = value
	case "compressors":
		compressors := strings.Split(value, ",")
		if len(compressors) < 1 {
			return fmt.Errorf("must have at least 1 compressor")
		}
		p.Compressors = compressors
	case "connect":
		switch strings.ToLower(value) {
		case "automatic":
		case "direct":
			p.Connect = SingleConnect
		default:
			return fmt.Errorf("invalid 'connect' value: %s", value)
		}

		p.ConnectSet = true
	case "connecttimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.ConnectTimeout = time.Duration(n) * time.Millisecond
		p.ConnectTimeoutSet = true
	case "heartbeatintervalms", "heartbeatfrequencyms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.HeartbeatInterval = time.Duration(n) * time.Millisecond
		p.HeartbeatIntervalSet = true
	case "journal":
		switch value {
		case "true":
			p.J = true
		case "false":
			p.J = false
		default:
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		p.JSet = true
	case "localthresholdms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.LocalThreshold = time.Duration(n) * time.Millisecond
		p.LocalThresholdSet = true
	case "maxidletimems":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.MaxConnIdleTime = time.Duration(n) * time.Millisecond
		p.MaxConnIdleTimeSet = true
	case "maxpoolsize":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.MaxPoolSize = uint64(n)
		p.MaxPoolSizeSet = true
	case "minpoolsize":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.MinPoolSize = uint64(n)
		p.MinPoolSizeSet = true
	case "readconcernlevel":
		p.ReadConcernLevel = value
	case "readpreference":
		p.ReadPreference = value
	case "readpreferencetags":
		if value == "" {
			// for when readPreferenceTags= at end of URI
			break
		}

		tags := make(map[string]string)
		items := strings.Split(value, ",")
		for _, item := range items {
			parts := strings.Split(item, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid value for %s: %s", key, value)
			}
			tags[parts[0]] = parts[1]
		}
		p.ReadPreferenceTagSets = append(p.ReadPreferenceTagSets, tags)
	case "maxstaleness", "maxstalenessseconds":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.MaxStaleness = time.Duration(n) * time.Second
		p.MaxStalenessSet = true
	case "replicaset":
		p.ReplicaSet = value
	case "retrywrites":
		switch value {
		case "true":
			p.RetryWrites = true
		case "false":
			p.RetryWrites = false
		default:
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		p.RetryWritesSet = true
	case "retryreads":
		switch value {
		case "true":
			p.RetryReads = true
		case "false":
			p.RetryReads = false
		default:
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		p.RetryReadsSet = true
	case "serverselectiontimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.ServerSelectionTimeout = time.Duration(n) * time.Millisecond
		p.ServerSelectionTimeoutSet = true
	case "sockettimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.SocketTimeout = time.Duration(n) * time.Millisecond
		p.SocketTimeoutSet = true
	case "ssl", "tls":
		switch value {
		case "true":
			p.SSL = true
		case "false":
			p.SSL = false
		default:
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		if p.tlsssl != nil && *p.tlsssl != p.SSL {
			return errors.New("tls and ssl options, when both specified, must be equivalent")
		}

		p.tlsssl = new(bool)
		*p.tlsssl = p.SSL

		p.SSLSet = true
	case "sslclientcertificatekeyfile", "tlscertificatekeyfile":
		p.SSL = true
		p.SSLSet = true
		p.SSLClientCertificateKeyFile = value
		p.SSLClientCertificateKeyFileSet = true
	case "sslclientcertificatekeypassword", "tlscertificatekeyfilepassword":
		p.SSLClientCertificateKeyPassword = func() string { return value }
		p.SSLClientCertificateKeyPasswordSet = true
	case "sslinsecure", "tlsinsecure":
		switch value {
		case "true":
			p.SSLInsecure = true
		case "false":
			p.SSLInsecure = false
		default:
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		p.SSLInsecureSet = true
	case "sslcertificateauthorityfile", "tlscafile":
		p.SSL = true
		p.SSLSet = true
		p.SSLCaFile = value
		p.SSLCaFileSet = true
	case "w":
		if w, err := strconv.Atoi(value); err == nil {
			if w < 0 {
				return fmt.Errorf("invalid value for %s: %s", key, value)
			}

			p.WNumber = w
			p.WNumberSet = true
			p.WString = ""
			break
		}

		p.WString = value
		p.WNumberSet = false

	case "wtimeoutms":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.WTimeout = time.Duration(n) * time.Millisecond
		p.WTimeoutSet = true
	case "wtimeout":
		// Defer to wtimeoutms, but not to a manually-set option.
		if p.WTimeoutSet {
			break
		}
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		p.WTimeout = time.Duration(n) * time.Millisecond
	case "zlibcompressionlevel":
		level, err := strconv.Atoi(value)
		if err != nil || (level < -1 || level > 9) {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		if level == -1 {
			level = wiremessage.DefaultZlibLevel
		}
		p.ZlibLevel = level
		p.ZlibLevelSet = true
	case "zstdcompressionlevel":
		const maxZstdLevel = 22 // https://github.com/facebook/zstd/blob/a880ca239b447968493dd2fed3850e766d6305cc/contrib/linux-kernel/lib/zstd/compress.c#L3291
		level, err := strconv.Atoi(value)
		if err != nil || (level < -1 || level > maxZstdLevel) {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}

		if level == -1 {
			level = wiremessage.DefaultZstdLevel
		}
		p.ZstdLevel = level
		p.ZstdLevelSet = true
	default:
		if p.UnknownOptions == nil {
			p.UnknownOptions = make(map[string][]string)
		}
		p.UnknownOptions[lowerKey] = append(p.UnknownOptions[lowerKey], value)
	}

	if p.Options == nil {
		p.Options = make(map[string][]string)
	}
	p.Options[lowerKey] = append(p.Options[lowerKey], value)

	return nil
}

func extractQueryArgsFromURI(uri string) ([]string, error) {
	if len(uri) == 0 {
		return nil, nil
	}

	if uri[0] != '?' {
		return nil, errors.New("must have a ? separator between path and query")
	}

	uri = uri[1:]
	if len(uri) == 0 {
		return nil, nil
	}
	return strings.FieldsFunc(uri, func(r rune) bool { return r == ';' || r == '&' }), nil

}

type extractedDatabase struct {
	uri string
	db  string
}

// extractDatabaseFromURI is a helper function to retrieve information about
// the database from the passed in URI. It accepts as an argument the currently
// parsed URI and returns the remainder of the uri, the database it found,
// and any error it encounters while parsing.
func extractDatabaseFromURI(uri string) (extractedDatabase, error) {
	if len(uri) == 0 {
		return extractedDatabase{}, nil
	}

	if uri[0] != '/' {
		return extractedDatabase{}, errors.New("must have a / separator between hosts and path")
	}

	uri = uri[1:]
	if len(uri) == 0 {
		return extractedDatabase{}, nil
	}

	database := uri
	if idx := strings.IndexRune(uri, '?'); idx != -1 {
		database = uri[:idx]
	}

	escapedDatabase, err := url.QueryUnescape(database)
	if err != nil {
		return extractedDatabase{}, internal.WrapErrorf(err, "invalid database \"%s\"", database)
	}

	uri = uri[len(database):]

	return extractedDatabase{
		uri: uri,
		db:  escapedDatabase,
	}, nil
}
