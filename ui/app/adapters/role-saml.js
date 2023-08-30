/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { inject as service } from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { v4 as uuidv4 } from 'uuid';

export default ApplicationAdapter.extend({
  router: service(),

  // generateClientChallenge generates a client challenge from a verifier.
  // The client challenge is the base64(sha256(verifier)). The verifier is
  // later presented to the server to obtain the resulting Vault token.
  async generateClientChallenge(verifier) {
    const encoder = new TextEncoder();
    const data = encoder.encode(verifier);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = new Uint8Array(hashBuffer);
    return btoa(String.fromCharCode.apply(null, hashArray));
  },

  async findRecord(store, type, id) {
    let [path, role] = JSON.parse(id);
    path = encodePath(path);

    // Create the client verifier and challenge
    const verifier = uuidv4();
    const challenge = await this.generateClientChallenge(verifier);
    // Kick off the authentication flow by generating the SSO service URL
    // It requires the client challenge generated from the verifier. We'll
    // later provide the verifier to match up with the challenge on the server
    // when we poll for the Vault token by its returned token poll ID.
    const response = await this.ajax(`/v1/auth/${path}/sso_service_url`, 'PUT', {
      data: {
        role,
        client_challenge: challenge,
        client_type: 'browser',
      },
    });
    return {
      ...response.data,
      client_verifier: verifier,
    };
  },
});
