/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export default function (server) {
  // ENTITY SEARCH SELECT
  server.get('/identity/entity/id', () => ({
    data: {
      key_info: { '1234-12345': { name: 'test-entity' } },
      keys: ['1234-12345'],
    },
  }));

  // GROUP SEARCH SELECT
  server.get('/identity/group/id', () => ({
    data: {
      key_info: { 'abcdef-123': { name: 'test-group' } },
      keys: ['abcdef-123'],
    },
  }));
}
