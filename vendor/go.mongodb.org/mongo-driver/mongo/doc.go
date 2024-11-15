// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// NOTE: This documentation should be kept in line with the Example* test functions.

// Package mongo provides a MongoDB Driver API for Go.
//
// Basic usage of the driver starts with creating a Client from a connection
// string. To do so, call Connect:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
//	if err != nil { return err }
//
// This will create a new client and start monitoring the MongoDB server on localhost.
// The Database and Collection types can be used to access the database:
//
//	collection := client.Database("baz").Collection("qux")
//
// A Collection can be used to query the database or insert documents:
//
//	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
//	if err != nil { return err }
//	id := res.InsertedID
//
// Several methods return a cursor, which can be used like this:
//
//	cur, err := collection.Find(context.Background(), bson.D{})
//	if err != nil { log.Fatal(err) }
//	defer cur.Close(context.Background())
//	for cur.Next(context.Background()) {
//	  // To decode into a struct, use cursor.Decode()
//	  result := struct{
//	    Foo string
//	    Bar int32
//	  }{}
//	  err := cur.Decode(&result)
//	  if err != nil { log.Fatal(err) }
//	  // do something with result...
//
//	  // To get the raw bson bytes use cursor.Current
//	  raw := cur.Current
//	  // do something with raw...
//	}
//	if err := cur.Err(); err != nil {
//	  return err
//	}
//
// Cursor.All will decode all of the returned elements at once:
//
//	var results []struct{
//	  Foo string
//	  Bar int32
//	}
//	if err = cur.All(context.Background(), &results); err != nil {
//	  log.Fatal(err)
//	}
//	// do something with results...
//
// Methods that only return a single document will return a *SingleResult, which works
// like a *sql.Row:
//
//	result := struct{
//	  Foo string
//	  Bar int32
//	}{}
//	filter := bson.D{{"hello", "world"}}
//	err := collection.FindOne(context.Background(), filter).Decode(&result)
//	if err != nil { return err }
//	// do something with result...
//
// All Client, Collection, and Database methods that take parameters of type interface{}
// will return ErrNilDocument if nil is passed in for an interface{}.
//
// Additional examples can be found under the examples directory in the driver's repository and
// on the MongoDB website.
//
// # Error Handling
//
// Errors from the MongoDB server will implement the ServerError interface, which has functions to check for specific
// error codes, labels, and message substrings. These can be used to check for and handle specific errors. Some methods,
// like InsertMany and BulkWrite, can return an error representing multiple errors, and in those cases the ServerError
// functions will return true if any of the contained errors satisfy the check.
//
// There are also helper functions to check for certain specific types of errors:
//
//	IsDuplicateKeyError(error)
//	IsNetworkError(error)
//	IsTimeout(error)
//
// # Potential DNS Issues
//
// Building with Go 1.11+ and using connection strings with the "mongodb+srv"[1] scheme is unfortunately
// incompatible with some DNS servers in the wild due to the change introduced in
// https://github.com/golang/go/issues/10622. You may receive an error with the message "cannot unmarshal DNS message"
// while running an operation when using DNS servers that non-compliantly compress SRV records. Old versions of kube-dns
// and the native DNS resolver (systemd-resolver) on Ubuntu 18.04 are known to be non-compliant in this manner. We suggest
// using a different DNS server (8.8.8.8 is the common default), and, if that's not possible, avoiding the "mongodb+srv"
// scheme.
//
// # Client Side Encryption
//
// Client-side encryption is a new feature in MongoDB 4.2 that allows specific data fields to be encrypted. Using this
// feature requires specifying the "cse" build tag during compilation:
//
//	go build -tags cse
//
// Note: Auto encryption is an enterprise- and Atlas-only feature.
//
// The libmongocrypt C library is required when using client-side encryption. Specific versions of libmongocrypt
// are required for different versions of the Go Driver:
//
// - Go Driver v1.2.0 requires libmongocrypt v1.0.0 or higher
//
// - Go Driver v1.5.0 requires libmongocrypt v1.1.0 or higher
//
// - Go Driver v1.8.0 requires libmongocrypt v1.3.0 or higher
//
// - Go Driver v1.10.0 requires libmongocrypt v1.5.0 or higher.
// There is a severe bug when calling RewrapManyDataKey with libmongocrypt versions less than 1.5.2.
// This bug may result in data corruption.
// Please use libmongocrypt 1.5.2 or higher when calling RewrapManyDataKey.
//
// - Go Driver v1.12.0 requires libmongocrypt v1.8.0 or higher.
//
// To install libmongocrypt, follow the instructions for your
// operating system:
//
// 1. Linux: follow the instructions listed at
// https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-from-distribution-packages to install the correct
// deb/rpm package.
//
// 2. Mac: Follow the instructions listed at https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-on-macos
// to install packages via brew and compile the libmongocrypt source code.
//
// 3. Windows:
//
//	mkdir -p c:/libmongocrypt/bin
//	mkdir -p c:/libmongocrypt/include
//
//	// Run the curl command in an empty directory as it will create new directories when unpacked.
//	curl https://s3.amazonaws.com/mciuploads/libmongocrypt/windows/latest_release/libmongocrypt.tar.gz --output libmongocrypt.tar.gz
//	tar -xvzf libmongocrypt.tar.gz
//
//	cp ./bin/mongocrypt.dll c:/libmongocrypt/bin
//	cp ./include/mongocrypt/*.h c:/libmongocrypt/include
//	export PATH=$PATH:/cygdrive/c/libmongocrypt/bin
//
// libmongocrypt communicates with the mongocryptd process or mongo_crypt shared library for automatic encryption.
// See AutoEncryptionOpts.SetExtraOptions for options to configure use of mongocryptd or mongo_crypt.
//
// [1] See https://www.mongodb.com/docs/manual/reference/connection-string/#dns-seedlist-connection-format
package mongo
