/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { addMinutes } from 'date-fns';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import license from '../../pages/components/license-info';
import { allFeatures } from 'vault/helpers/all-features';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const FEATURES = allFeatures();

const component = create(license);

module('Integration | Component | license info', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it renders feature status properly for features associated with license', async function (assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}}/>`
    );
    assert.strictEqual(
      component.detailRows.length,
      3,
      'Shows License ID, Valid from, and License State rows'
    );
    assert.strictEqual(component.featureRows.length, FEATURES.length, 'it renders all of the features');
    const activeFeatures = component.featureRows.filter((f) => f.featureStatus === 'Active');
    assert.strictEqual(activeFeatures.length, 2, 'Has two features listed as active');
  });

  test('it renders properly for autoloaded license', async function (assert) {
    const now = Date.now();
    this.set('licenseId', 'test');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('autoloaded', true);
    this.set('startTime', now);
    this.set('features', ['HSM', 'Namespaces']);
    await render(
      hbs`<LicenseInfo
        @licenseId={{this.licenseId}}
        @expirationTime={{this.expirationTime}}
        @startTime={{this.startTime}}
        @features={{this.features}}
        @autoloaded={{true}}
      />`
    );
    const row = component.detailRows.filterBy('rowName', 'License state')[0];
    assert.strictEqual(row.rowValue, 'Autoloaded', 'Shows autoloaded status');
  });

  test('it renders Performance Standby as inactive if count is 0', async function (assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('model', { performanceStandbyCount: 0 });
    this.set('features', ['Performance Standby', 'Namespaces']);

    await render(
      hbs`<LicenseInfo @licenseId={{this.licenseId}} @expirationTime={{this.expirationTime}} @startTime={{this.startTime}} @features={{this.features}} @model={{this.model}}/>`
    );

    const row = component.featureRows.filterBy('featureName', 'Performance Standby')[0];
    assert.strictEqual(
      row.featureStatus,
      'Not Active',
      'renders feature as inactive because when count is 0'
    );
  });

  test('it renders Performance Standby as active and shows count', async function (assert) {
    const now = Date.now();
    this.set('licenseId', 'temporary');
    this.set('expirationTime', addMinutes(now, 30));
    this.set('startTime', now);
    this.set('model', { performanceStandbyCount: 4 });
    this.set('features', ['Performance Standby', 'Namespaces']);

    await render(
      hbs`<LicenseInfo
        @licenseId={{this.licenseId}}
        @expirationTime={{this.expirationTime}}
        @startTime={{this.startTime}}
        @features={{this.features}}
        @performanceStandbyCount={{this.model.performanceStandbyCount}}
      />`
    );

    const row = component.featureRows.filterBy('featureName', 'Performance Standby')[0];
    assert.strictEqual(
      row.featureStatus,
      'Active — 4 standby nodes allotted',
      'renders active and displays count'
    );
  });
});
