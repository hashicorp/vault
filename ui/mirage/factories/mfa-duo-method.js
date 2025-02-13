/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  api_hostname: 'api-foobar.duosecurity.com',
  mount_accessor: '',
  name: '', // returned but cannot be set at this time
  // namespace_path: 'root', commented out for testing behavior when this is not present
  pushinfo: '',
  type: 'duo',
  use_passcode: false,
  username_template: '',

  afterCreate(record) {
    if (record.name) {
      console.warn('Endpoint ignored these unrecognized parameters: [name]'); // eslint-disable-line
      record.name = '';
    }
  },
});
