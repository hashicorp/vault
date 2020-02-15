// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package scram

import "sync"

// Server implements the server side of SCRAM authentication.  It holds
// configuration values needed to initialize new server-side conversations.
// Generally, this can be persistent within an application.
type Server struct {
	sync.RWMutex
	credentialCB CredentialLookup
	nonceGen     NonceGeneratorFcn
	hashGen      HashGeneratorFcn
}

func newServer(cl CredentialLookup, fcn HashGeneratorFcn) (*Server, error) {
	return &Server{
		credentialCB: cl,
		nonceGen:     defaultNonceGenerator,
		hashGen:      fcn,
	}, nil
}

// WithNonceGenerator replaces the default nonce generator (base64 encoding of
// 24 bytes from crypto/rand) with a custom generator.  This is provided for
// testing or for users with custom nonce requirements.
func (s *Server) WithNonceGenerator(ng NonceGeneratorFcn) *Server {
	s.Lock()
	defer s.Unlock()
	s.nonceGen = ng
	return s
}

// NewConversation constructs a server-side authentication conversation.
// Conversations cannot be reused, so this must be called for each new
// authentication attempt.
func (s *Server) NewConversation() *ServerConversation {
	s.RLock()
	defer s.RUnlock()
	return &ServerConversation{
		nonceGen:     s.nonceGen,
		hashGen:      s.hashGen,
		credentialCB: s.credentialCB,
	}
}
