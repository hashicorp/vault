/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  base_url: 'okta.com',
  mount_accessor: '',
  name: '', // returned but cannot be set at this time
  namespace_id: 'root',
  org_name: 'dev-foobar',
  type: 'okta',
  username_template: '', // returned but cannot be set at this time

  afterCreate(record) {
    if (record.name) {
      console.warn('Endpoint ignored these unrecognized parameters: [name]'); // eslint-disable-line
      record.name = '';
    }
  },
});
