/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('secrets', function () {
    this.route('overview');
    this.route('destinations', function () {
      this.route('create', function () {
        this.route('destination', { path: '/:type' });
      });
      this.route('destination', { path: '/:type/:name' }, function () {
        this.route('edit');
        this.route('details');
        this.route('secrets', function () {
          this.route('sync');
        });
      });
    });
  });
});
