/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('overview');
  this.route('roles', function () {
    this.route('create');
    this.route('role', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
      this.route('credentials');
    });
  });
  this.route('configure');
  this.route('configuration', function () {
    // index here is the route that shows the tabs for "General settings" and "Aws settings"
    this.route('plugin-settings', function () {
      this.route('edit');
    });
    this.route('mount', function () {
      this.route('edit');
    });
  });
});
