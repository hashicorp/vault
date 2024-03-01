/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  // There are two list routes because Ember won't let a route param (e.g. *path_to_secret) be blank.
  // *path_to_secret is used when we're listing a secret directory.
  // Must use a wildcard for path-to-secret because the value can contain a forward slash if it's a secret directory. Ember's router decodes encoded forward slashes which leads to beep%2fboop becoming beep/boop and messing up routing after copying and pasting the URL.
  this.route('list');
  this.route('list-directory', { path: '/list/*path_to_secret' });
  this.route('create');
  this.route('secret', { path: '/:name' }, function () {
    this.route('paths');
    this.route('details', function () {
      this.route('edit'); // route to create new version of a secret
    });
    this.route('metadata', function () {
      this.route('edit');
      this.route('versions');
      this.route('diff');
    });
  });
  this.route('configuration');
});
