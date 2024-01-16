/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  // dynamic props
  dynamic: trait({
    backend: 'database',
    name: 'some-role',
    type: 'dynamic',
    default_ttl: '1h',
    max_ttl: '24h',
    rotation_period: '24h',
    path: 'roles',
    db_name: 'connection',
  }),
});
