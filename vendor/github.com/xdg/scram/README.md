[![GoDoc](https://godoc.org/github.com/xdg/scram?status.svg)](https://godoc.org/github.com/xdg/scram)
[![Build Status](https://travis-ci.org/xdg/scram.svg?branch=master)](https://travis-ci.org/xdg/scram)

# scram – Go implementation of RFC-5802

## Description

Package scram provides client and server implementations of the Salted
Challenge Response Authentication Mechanism (SCRAM) described in
[RFC-5802](https://tools.ietf.org/html/rfc5802) and
[RFC-7677](https://tools.ietf.org/html/rfc7677).

It includes both client and server side support.

Channel binding and extensions are not (yet) supported.

## Examples

### Client side

    package main

    import "github.com/xdg/scram"

    func main() {
        // Get Client with username, password and (optional) authorization ID.
        clientSHA1, err := scram.SHA1.NewClient("mulder", "trustno1", "")
        if err != nil {
            panic(err)
        }

        // Prepare the authentication conversation. Use the empty string as the
        // initial server message argument to start the conversation.
        conv := clientSHA1.NewConversation()
        var serverMsg string

        // Get the first message, send it and read the response.
        firstMsg, err := conv.Step(serverMsg)
        if err != nil {
            panic(err)
        }
        serverMsg = sendClientMsg(firstMsg)

        // Get the second message, send it, and read the response.
        secondMsg, err := conv.Step(serverMsg)
        if err != nil {
            panic(err)
        }
        serverMsg = sendClientMsg(secondMsg)

        // Validate the server's final message.  We have no further message to
        // send so ignore that return value.
        _, err = conv.Step(serverMsg)
        if err != nil {
            panic(err)
        }

        return
    }

    func sendClientMsg(s string) string {
        // A real implementation would send this to a server and read a reply.
        return ""
    }

## Copyright and License

Copyright 2018 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
