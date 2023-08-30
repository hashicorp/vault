/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  // There are two list routes because Ember won't let a route param (e.g. *path_to_secret) be blank.
  // :path_to_secret is used when we're listing a secret directory. Example { path: '/:beep%2Fboop%2F/directory' });
  this.route('list');
  this.route('list-directory', { path: '/:path_to_secret/directory' });
  this.route('create');
  this.route('secret', { path: '/:name' }, function () {
    this.route('paths');
    this.route('details', function () {
      this.route('edit'); // route to create new version of a secret
    });
    this.route('metadata', function () {
      this.route('edit');
      this.route('versions');
    });
  });
  this.route('configuration');
});
