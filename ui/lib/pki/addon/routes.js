/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';
import ENV from 'vault/config/environment';

export default buildRoutes(function () {
  // Internal PKI
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

  // Public/External PKI
  if (ENV.environment !== 'production') {
    this.route('external', function () {
      this.route('configuration');
      this.route('overview');
      this.route('roles', function () {
        this.route('role', { path: '/:role_name' }, function () {
          this.route('details');
          this.route('active-orders'); // Specific order details route to orders.order.details
        });
      });
      this.route('orders', function () {
        this.route('order', { path: '/:order_id' });
      });
      this.route('certificates', function () {
        this.route('certificate', { path: '/:serial_number' });
      });
      this.route('dns-providers');
      this.route('acme-accounts');
    });
  }
});
