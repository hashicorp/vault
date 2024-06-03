/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('overview');
  this.route('roles', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('role', { path: '/:role' }, function () {
      this.route('details');
      this.route('edit');
      this.route('generate');
      this.route('sign');
    });
  });
  this.route('issuers', function () {
    this.route('index', { path: '/' });
    this.route('import');
    this.route('generate-root');
    this.route('generate-intermediate');
    this.route('issuer', { path: '/:issuer_ref' }, function () {
      this.route('details');
      this.route('edit');
      this.route('sign');
      this.route('cross-sign');
      this.route('rotate-root');
    });
  });
  this.route('keys', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('import');
    this.route('key', { path: '/:key_id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('certificates', function () {
    this.route('index', { path: '/' });
    this.route('certificate', { path: '/:serial' }, function () {
      this.route('details');
      this.route('edit');
    });
  });
  this.route('tidy', function () {
    this.route('index', { path: '/' });
    this.route('auto', function () {
      this.route('configure');
    });
    this.route('manual');
  });
  this.route('configuration', function () {
    this.route('index', { path: '/' });
    this.route('create');
    this.route('edit');
  });
});
