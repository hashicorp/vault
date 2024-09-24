/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Model } from 'miragejs';

export default Model.extend({
  account: '', // should match ID
  library: '',
  available: false,
  borrower_client_token: undefined,
});
