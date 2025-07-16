/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class VaultClusterOidcProviderController extends Controller {
  queryParams = [
    'scope', // *
    'response_type', // *
    'client_id', // *
    'redirect_uri', // *
    'state', // *
    'nonce', // *
    'display',
    'prompt',
    'max_age',
    'code_challenge',
    'code_challenge_method',
    'request',
    'request_uri',
  ];
  scope = null;
  response_type = null;
  client_id = null;
  redirect_uri = null;
  state = null;
  nonce = null;
  display = null;
  prompt = null;
  max_age = null;
  code_challenge = null;
  code_challenge_method = null;
  request = null;
  request_uri = null;
}
