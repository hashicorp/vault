/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  accessor: 'kv_eb4acbae', // mount will be added to API response for use in the ui but leaving since it is a documented property
  mount: 'my-kv',
  secret_name: 'my-path/my-secret-1',
  sync_status: 'SYNCED',
  updated_at: '2023-09-20T10:51:53.961861096-04:00',
  // set on create for lookup by destination
  type: null,
  name: null,
});
