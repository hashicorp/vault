/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { SELECTORS } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth-config-form options', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.router = this.owner.lookup('service:router');
    this.store = this.owner.lookup('service:store');
    this.path = 'my-auth-method/';
    this.model = this.store.createRecord('auth-method', { path: this.path, type: 'approle' });
    this.model.set('config', this.store.createRecord('mount-config'));
  });

  test('it submits data correctly', async function (assert) {
    assert.expect(2);
    this.router.reopen({
      transitionTo() {
        return {
          followRedirects() {
            assert.ok('calls transitionTo on save');
          },
        };
      },
    });

    this.server.post(`sys/mounts/auth/${this.path}/tune`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      const expected = {
        default_lease_ttl: '30s',
        listing_visibility: 'unauth',
        user_lockout_config: {
          lockout_threshold: '7',
          lockout_duration: '600s',
          lockout_counter_reset: '5s',
          lockout_disable: true,
        },
      };
      assert.propEqual(payload, expected, 'payload contains tune parameters');
      return { payload };
    });
    await render(hbs`<AuthConfigForm::Options @model={{this.model}} />`);

    await click(SELECTORS.inputByAttr('config.listingVisibility'));

    await click(SELECTORS.ttl.toggle('Default Lease TTL'));
    await fillIn(SELECTORS.ttl.input('Default Lease TTL'), '30');

    await fillIn(SELECTORS.inputByAttr('config.lockoutThreshold'), '7');

    await click(SELECTORS.ttl.toggle('Lockout duration'));
    await fillIn(SELECTORS.ttl.input('Lockout duration'), '10');
    await fillIn(
      `${SELECTORS.inputByAttr('config.lockoutDuration')} ${SELECTORS.selectByAttr('ttl-unit')}`,
      'm'
    );
    await click(SELECTORS.ttl.toggle('Lockout counter reset'));
    await fillIn(SELECTORS.ttl.input('Lockout counter reset'), '5');

    await click(SELECTORS.inputByAttr('config.lockoutDisable'));

    await click('[data-test-save-config]');
  });
});
