/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This cannot be called kv-metadata because mirage checks for plural factory names, and metadata and data are considered plural. It will throw an error.
import { Factory } from 'ember-cli-mirage';
/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
export default Factory.extend({
  openapi: '3.0.2',
  info: {
    title: 'HashiCorp Vault API',
    description: 'HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.',
    version: '1.0.0',
    license: {
      name: 'Mozilla Public License 2.0',
      url: 'https://www.mozilla.org/en-US/MPL/2.0',
    },
  },
  paths: {
    '/auth/token/create': {
      description: 'The token create path is used to create new tokens.',
      post: {
        summary: 'The token create path is used to create new tokens.',
        tags: ['auth'],
        responses: {
          200: {
            description: 'OK',
          },
        },
      },
    },
  },
});
