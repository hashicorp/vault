/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import config from '../config/environment';

export default Controller.extend({
  queryParams: [
    {
      wrappedToken: 'wrapped_token',
      redirectTo: 'redirect_to',
    },
  ],
  wrappedToken: '',
  redirectTo: '',
  env: config.environment,
  auth: service(),
  store: service(),
});
