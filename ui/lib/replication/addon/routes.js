/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('index', { path: '/' });
  this.route('mode', { path: '/:replication_mode' }, function () {
    //details
    this.route('index', { path: '/' });
    this.route('manage');
    this.route('secondaries', function () {
      this.route('add');
      this.route('revoke');
      this.route('config-show', { path: '/config/show/:secondary_id' });
      this.route('config-edit', { path: '/config/edit/:secondary_id' });
      this.route('config-create', { path: '/config/create/:secondary_id' });
    });
  });
});
