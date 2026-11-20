/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';
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

  // Defined before the loops below: the `return` in the unsupported-types loop ends this module's setup.
  module('saved lockout settings', function () {
    test('it shows the saved values and saves them back unchanged', async function (assert) {
      assert.expect(8);
      const path = 'my-userpass-auth';
      const requests = [];
      this.server.post(`sys/mounts/auth/${path}/tune`, (schema, req) => {
        requests.push(JSON.parse(req.requestBody));
        return {};
      });
      this.form = new AuthMethodForm({
        path,
        config: { listing_visibility: false },
        user_lockout_config: {
          lockout_threshold: '7',
          lockout_duration: '600s',
          lockout_counter_reset_duration: '1800s',
          lockout_disable: true,
        },
      });
      this.form.type = 'userpass';
      await render(hbs`<AuthConfigForm::Options @form={{this.form}} />`);

      assert.dom(GENERAL.inputByAttr('user_lockout_config.lockout_threshold')).hasValue('7');
      assert.dom(GENERAL.inputByAttr('user_lockout_config.lockout_disable')).isChecked();
      assert.dom(GENERAL.ttl.input('Lockout duration')).hasValue('10');
      assert.dom(TTL_PICKER.unitInputByLabel('Lockout duration')).hasValue('m');
      assert.dom(GENERAL.ttl.input('Lockout counter reset')).hasValue('30');
      assert.dom(TTL_PICKER.unitInputByLabel('Lockout counter reset')).hasValue('m');

      await click(GENERAL.submitButton);
      assert.strictEqual(requests.length, 1, 'it saves once');
      assert.propEqual(
        requests[0].user_lockout_config,
        {
          lockout_threshold: '7',
          lockout_duration: '600s',
          lockout_counter_reset_duration: '1800s',
          lockout_disable: true,
        },
        'the saved values are sent back with the keys Vault reads'
      );
    });
  });

  // Defined before the loops below: the `return` in the unsupported-types loop ends this module's setup.
  module('durations longer than Vault supports', function (hooks) {
    hooks.beforeEach(function () {
      this.path = 'my-userpass-auth';
      this.tuneRequests = [];
      this.server.post(`sys/mounts/auth/${this.path}/tune`, (schema, req) => {
        this.tuneRequests.push(JSON.parse(req.requestBody));
        return {};
      });
      this.enterDuration = async (label, value, unit = 'd') => {
        await click(GENERAL.ttl.toggle(label));
        await fillIn(GENERAL.ttl.input(label), value);
        await fillIn(TTL_PICKER.unitInputByLabel(label), unit);
      };
    });

    const FIELDS = [
      ['Default Lease TTL', 'config.default_lease_ttl'],
      ['Max Lease TTL', 'config.max_lease_ttl'],
      ['Lockout duration', 'user_lockout_config.lockout_duration'],
      ['Lockout counter reset', 'user_lockout_config.lockout_counter_reset_duration'],
    ];

    for (const [label, attr] of FIELDS) {
      test(`it blocks ${label}`, async function (assert) {
        assert.expect(6);
        await this.renderComponent(this.path, 'userpass');

        await this.enterDuration(label, '9999999999');
        await click(GENERAL.submitButton);

        assert
          .dom(GENERAL.validationErrorByAttr(attr))
          .hasText(`${label} is too large. Maximum allowed value is 106,751 days.`);
        assert.dom(GENERAL.invalidFormMessage).hasText('There is an error with this form.');
        assert.strictEqual(this.tuneRequests.length, 0, 'nothing is saved');
        assert.false(this.transitionStub.called, 'the form stays open');

        await fillIn(GENERAL.ttl.input(label), '106751');
        await click(GENERAL.submitButton);
        assert.strictEqual(this.tuneRequests.length, 1, 'it saves once the duration is valid');
        assert.dom(GENERAL.invalidFormMessage).doesNotExist('the form alert clears');
      });
    }

    test('it blocks a duration too long to be stored as a number', async function (assert) {
      assert.expect(2);
      await this.renderComponent(this.path, 'userpass');

      await this.enterDuration('Default Lease TTL', '9'.repeat(400));
      await click(GENERAL.submitButton);

      assert.dom(GENERAL.invalidFormMessage).exists('the form alert shows');
      assert.strictEqual(this.tuneRequests.length, 0, 'nothing is saved');
    });

    test('it saves when an over-limit lockout duration is toggled off', async function (assert) {
      assert.expect(3);
      await this.renderComponent(this.path, 'userpass');

      await this.enterDuration('Lockout duration', '9999999999');
      await click(GENERAL.ttl.toggle('Lockout duration'));
      await click(GENERAL.submitButton);

      assert.strictEqual(this.tuneRequests.length, 1, 'it saves');
      assert.strictEqual(
        this.tuneRequests[0].user_lockout_config.lockout_duration,
        '0',
        'the hidden value is replaced'
      );
      assert.dom(GENERAL.invalidFormMessage).doesNotExist('no error blocks the save');
    });

    test('it only reports durations that are too large', async function (assert) {
      assert.expect(2);
      await this.renderComponent(this.path, 'userpass');

      await this.enterDuration('Default Lease TTL', '106751');
      await click(GENERAL.submitButton);

      assert.strictEqual(this.tuneRequests.length, 1, 'the maximum whole days saves');
      assert.strictEqual(this.tuneRequests[0].default_lease_ttl, '9223286400s');
    });
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
            lockout_counter_reset_duration: '5s',
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
