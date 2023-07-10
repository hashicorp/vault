/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('secrets', function () {
    this.route('create');
    this.route('secret', { path: '/:name' }, function () {
      this.route('details');
      this.route('edit');
      this.route('metadata', function () {
        this.route('edit');
        this.route('versions');
        this.route('diff');
      });
    });
  });
  this.route('configuration');
});
