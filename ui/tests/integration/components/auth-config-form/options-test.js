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
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import AuthMethodForm from 'vault/forms/auth/method';
import sinon from 'sinon';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const userLockoutSupported = ['approle', 'ldap', 'userpass'];
const userLockoutUnsupported = filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: false })
  .map((m) => m.type)
  .filter((m) => !userLockoutSupported.includes(m));

module('Integration | Component | auth-config-form options', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flashSuccessSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon
      .stub(this.router, 'transitionTo')
      .returns({ followRedirects: () => Promise.resolve() });

    this.renderComponent = (path, type) => {
      this.form = new AuthMethodForm({
        path,
        config: { listing_visibility: false },
        user_lockout_config: {},
      });
      this.form.type = type;
      return render(hbs`<AuthConfigForm::Options @form={{this.form}} />`);
    };
  });

  test('it submits data correctly for token auth method', async function (assert) {
    assert.expect(8);

    const type = 'token';
    const path = `my-${type}-auth`;

    this.server.post(`sys/mounts/auth/${path}/tune`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      const expected = {
        default_lease_ttl: '30s',
        listing_visibility: 'unauth',
      };
      assert.propEqual(payload, expected, `${type} method payload contains tune parameters`);
      return { payload };
    });

    await this.renderComponent(path, type);
    assert
      .dom(GENERAL.inputByAttr('config.token_type'))
      .doesNotExist('does not render token_type for token auth method');

    await click(GENERAL.toggleInput('toggle-config.listing_visibility'));
    await click(GENERAL.ttl.toggle('Default Lease TTL'));
    await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

    assert.dom('[data-test-user-lockout-section]').doesNotExist('token does not render user lockout section');
    assert
      .dom(GENERAL.inputByAttr('user_lockout_config.lockout_threshold'))
      .doesNotExist('token method does not render lockout threshold');
    assert
      .dom(GENERAL.ttl.toggle('Lockout duration'))
      .doesNotExist('token method does not render lockout duration ');
    assert
      .dom(GENERAL.ttl.toggle('Lockout counter reset'))
      .doesNotExist('token method does not render lockout counter reset');
    assert
      .dom(GENERAL.inputByAttr('user_lockout_config.lockout_disable'))
      .doesNotExist('token method does not render lockout disable');

    await click(GENERAL.submitButton);

    assert.true(
      this.transitionStub.calledWith('vault.cluster.access.methods'),
      'transitions to access methods list on save'
    );
  });

  for (const type of userLockoutSupported) {
    test(`it submits data correctly for ${type} method (supports user_lockout_config)`, async function (assert) {
      assert.expect(3);

      const path = `my-${type}-auth`;

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

      await this.renderComponent(path, type);

      assert.dom('[data-test-user-lockout-section]').hasText('User lockout configuration');

      await click(GENERAL.toggleInput('toggle-config.listing_visibility'));
      await fillIn(GENERAL.inputByAttr('config.token_type'), 'default-batch');

      await click(GENERAL.ttl.toggle('Default Lease TTL'));
      await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

      await fillIn(GENERAL.inputByAttr('user_lockout_config.lockout_threshold'), '7');

      await click(GENERAL.ttl.toggle('Lockout duration'));
      await fillIn(GENERAL.ttl.input('Lockout duration'), '10');
      await fillIn(
        `${GENERAL.inputByAttr('user_lockout_config.lockout_duration')} ${GENERAL.selectByAttr('ttl-unit')}`,
        'm'
      );
      await click(GENERAL.ttl.toggle('Lockout counter reset'));
      await fillIn(GENERAL.ttl.input('Lockout counter reset'), '5');

      await click(GENERAL.inputByAttr('user_lockout_config.lockout_disable'));

      await click(GENERAL.submitButton);

      assert.true(
        this.transitionStub.calledWith('vault.cluster.access.methods'),
        'transitions to access methods list on save'
      );
    });

    test(`${type}: it renders error banner`, async function (assert) {
      assert.expect(3);
      const path = `my-${type}-auth`;
      this.server.post(`sys/mounts/auth/${path}/tune`, () => {
        return overrideResponse(400, { errors: ['uh oh'] });
      });
      await this.renderComponent(path, type);
      await click(GENERAL.submitButton);
      assert.false(this.flashSuccessSpy.calledOnce, 'flash success is NOT called');
      assert.false(this.transitionStub.calledOnce, 'transitionTo is NOT called');
      assert.dom(GENERAL.messageError).hasText('Error uh oh');
    });
  }

  for (const type of userLockoutUnsupported) {
    if (type === 'token') return; // separate test below because does not include tokenType field

    test(`it submits data correctly for ${type} auth method`, async function (assert) {
      assert.expect(7);

      const path = `my-${type}-auth`;

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

      await this.renderComponent(path, type);

      assert
        .dom('[data-test-user-lockout-section]')
        .doesNotExist(`${type} method does not render user lockout section`);

      await click(GENERAL.toggleInput('toggle-config.listing_visibility'));
      await fillIn(GENERAL.inputByAttr('config.token_type'), 'default-batch');

      await click(GENERAL.ttl.toggle('Default Lease TTL'));
      await fillIn(GENERAL.ttl.input('Default Lease TTL'), '30');

      assert
        .dom(GENERAL.inputByAttr('user_lockout_config.lockout_threshold'))
        .doesNotExist(`${type} method does not render lockout threshold`);
      assert
        .dom(GENERAL.ttl.toggle('Lockout duration'))
        .doesNotExist(`${type} method does not render lockout duration `);
      assert
        .dom(GENERAL.ttl.toggle('Lockout counter reset'))
        .doesNotExist(`${type} method does not render lockout counter reset`);
      assert
        .dom(GENERAL.inputByAttr('user_lockout_config.lockout_disable'))
        .doesNotExist(`${type} method does not render lockout disable`);

      await click(GENERAL.submitButton);

      assert.true(
        this.transitionStub.calledWith('vault.cluster.access.methods'),
        'transitions to access methods list on save'
      );
    });

    test(`${type}: it renders error banner`, async function (assert) {
      assert.expect(3);
      const path = `my-${type}-auth`;
      this.server.post(`sys/mounts/auth/${path}/tune`, () => {
        return overrideResponse(400, { errors: ['uh oh'] });
      });
      await this.renderComponent(path, type);
      await click(GENERAL.submitButton);
      assert.false(this.flashSuccessSpy.calledOnce, 'flash success is NOT called');
      assert.false(this.transitionStub.calledOnce, 'transitionTo is NOT called');
      assert.dom(GENERAL.messageError).hasText('Error uh oh');
    });
  }
});
