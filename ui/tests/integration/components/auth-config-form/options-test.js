/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { methods } from 'vault/helpers/mountable-auth-methods';

const userLockoutSupported = ['approle', 'ldap', 'userpass'];
const userLockoutUnsupported = methods()
  .map((m) => m.type)
  .filter((m) => !userLockoutSupported.includes(m));

module('Integration | Component | auth-config-form options', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.router = this.owner.lookup('service:router');
    this.store = this.owner.lookup('service:store');
    this.createModel = (path, type) => {
      this.model = this.store.createRecord('auth-method', { path, type });
      this.model.set('config', this.store.createRecord('mount-config'));
    };
  });

  for (const type of userLockoutSupported) {
    test(`it submits data correctly for ${type} method (supports user_lockout_config)`, async function (assert) {
      assert.expect(3);
      const path = `my-${type}-auth/`;
      this.createModel(path, type);

      this.router.reopen({
        transitionTo() {
          return {
            followRedirects() {
              assert.ok(true, `saving ${type} calls transitionTo on save`);
            },
          };
        },
      });

      this.server.post(`sys/mounts/auth/${path}/tune`, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          default_lease_ttl: '30s',
          listing_visibility: 'unauth',
          token_type: 'default-batch',
          user_lockout_config: {
            lockout_threshold: '7',
            lockout_duration: '600s',
            lockout_counter_reset: '5s',
            lockout_disable: true,
          },
        };
        assert.propEqual(payload, expected, `${type} method payload contains tune parameters`);
        return { payload };
      });
      await render(hbs`<AuthConfigForm::Options @model={{this.model}} />`);

      assert.dom('[data-test-user-lockout-section]').hasText('User lockout configuration');

      await click(GENERAL.inputByAttr('config.listingVisibility'));
      await fillIn(GENERAL.inputByAttr('config.tokenType'), 'default-batch');

      await click(GENERAL.ttl.toggle('Default Lease TTL'));
      await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

      await fillIn(GENERAL.inputByAttr('config.lockoutThreshold'), '7');

      await click(GENERAL.ttl.toggle('Lockout duration'));
      await fillIn(GENERAL.ttl.input('Lockout duration'), '10');
      await fillIn(
        `${GENERAL.inputByAttr('config.lockoutDuration')} ${GENERAL.selectByAttr('ttl-unit')}`,
        'm'
      );
      await click(GENERAL.ttl.toggle('Lockout counter reset'));
      await fillIn(GENERAL.ttl.input('Lockout counter reset'), '5');

      await click(GENERAL.inputByAttr('config.lockoutDisable'));

      await click('[data-test-save-config]');
    });
  }

  for (const type of userLockoutUnsupported) {
    if (type === 'token') return; // separate test below because does not include tokenType field

    test(`it submits data correctly for ${type} auth method`, async function (assert) {
      assert.expect(7);

      const path = `my-${type}-auth/`;
      this.createModel(path, type);

      this.router.reopen({
        transitionTo() {
          return {
            followRedirects() {
              assert.ok(true, `saving ${type} calls transitionTo on save`);
            },
          };
        },
      });

      this.server.post(`sys/mounts/auth/${path}/tune`, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          default_lease_ttl: '30s',
          listing_visibility: 'unauth',
          token_type: 'default-batch',
        };
        assert.propEqual(payload, expected, `${type} method payload contains tune parameters`);
        return { payload };
      });
      await render(hbs`<AuthConfigForm::Options @model={{this.model}} />`);

      assert
        .dom('[data-test-user-lockout-section]')
        .doesNotExist(`${type} method does not render user lockout section`);

      await click(GENERAL.inputByAttr('config.listingVisibility'));
      await fillIn(GENERAL.inputByAttr('config.tokenType'), 'default-batch');

      await click(GENERAL.ttl.toggle('Default Lease TTL'));
      await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

      assert
        .dom(GENERAL.inputByAttr('config.lockoutThreshold'))
        .doesNotExist(`${type} method does not render lockout threshold`);
      assert
        .dom(GENERAL.ttl.toggle('Lockout duration'))
        .doesNotExist(`${type} method does not render lockout duration `);
      assert
        .dom(GENERAL.ttl.toggle('Lockout counter reset'))
        .doesNotExist(`${type} method does not render lockout counter reset`);
      assert
        .dom(GENERAL.inputByAttr('config.lockoutDisable'))
        .doesNotExist(`${type} method does not render lockout disable`);

      await click('[data-test-save-config]');
    });
  }

  test('it submits data correctly for token auth method', async function (assert) {
    assert.expect(8);
    const type = 'token';
    const path = `my-${type}-auth/`;
    this.createModel(path, type);

    this.router.reopen({
      transitionTo() {
        return {
          followRedirects() {
            assert.ok(true, `saving token calls transitionTo on save`);
          },
        };
      },
    });

    this.server.post(`sys/mounts/auth/${path}/tune`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      const expected = {
        default_lease_ttl: '30s',
        listing_visibility: 'unauth',
      };
      assert.propEqual(payload, expected, `${type} method payload contains tune parameters`);
      return { payload };
    });
    await render(hbs`<AuthConfigForm::Options @model={{this.model}} />`);

    assert
      .dom(GENERAL.inputByAttr('config.tokenType'))
      .doesNotExist('does not render tokenType for token auth method');

    await click(GENERAL.inputByAttr('config.listingVisibility'));
    await click(GENERAL.ttl.toggle('Default Lease TTL'));
    await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

    assert.dom('[data-test-user-lockout-section]').doesNotExist('token does not render user lockout section');
    assert
      .dom(GENERAL.inputByAttr('config.lockoutThreshold'))
      .doesNotExist('token method does not render lockout threshold');
    assert
      .dom(GENERAL.ttl.toggle('Lockout duration'))
      .doesNotExist('token method does not render lockout duration ');
    assert
      .dom(GENERAL.ttl.toggle('Lockout counter reset'))
      .doesNotExist('token method does not render lockout counter reset');
    assert
      .dom(GENERAL.inputByAttr('config.lockoutDisable'))
      .doesNotExist('token method does not render lockout disable');

    await click('[data-test-save-config]');
  });
});
