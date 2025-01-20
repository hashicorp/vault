/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, currentURL, fillIn, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import mfaConfigHandler from 'vault/mirage/handlers/mfa-config';
import { Response } from 'miragejs';
import { underscore } from '@ember/string';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | mfa-method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    mfaConfigHandler(this.server);
    this.store = this.owner.lookup('service:store');
    this.getMethods = () =>
      ['Totp', 'Duo', 'Okta', 'Pingid'].reduce((methods, type) => {
        methods = [...methods, ...this.server.db[`mfa${type}Methods`].where({})];
        return methods;
      }, []);
    return authPage.login();
  });

  test('it should display landing page when no methods exist', async function (assert) {
    this.server.get('/identity/mfa/method/', () => new Response(404, {}, { errors: [] }));
    await visit('/vault/access/mfa/methods');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.index',
      'Route redirects to mfa index when no methods exist'
    );
    await click('[data-test-mfa-configure]');
    assert.strictEqual(currentRouteName(), 'vault.cluster.access.mfa.methods.create');
  });

  test('it should list methods', async function (assert) {
    await visit('/vault/access/mfa');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.index',
      'Parent route redirects to methods when some exist'
    );
    assert.dom('[data-test-tab="methods"]').hasClass('active', 'Methods tab is active');
    assert.dom('.toolbar-link').exists({ count: 1 }, 'Correct number of toolbar links render');
    assert.dom('[data-test-mfa-method-create]').includesText('New MFA method', 'New mfa link renders');

    await click('[data-test-mfa-method-create]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.create',
      'New method link transitions to create route'
    );
    await click('.hds-breadcrumb a');

    const methods = this.getMethods();
    const model = this.store.peekRecord('mfa-method', methods[0].id);
    assert.dom('[data-test-mfa-method-list-item]').exists({ count: methods.length }, 'Methods list renders');
    assert.dom(`[data-test-mfa-method-list-icon="${model.type}"]`).exists('Icon renders for list item');
    assert
      .dom(`[data-test-mfa-method-list-item="${model.id}"]`)
      .includesText(
        `${model.name} ${model.id} Namespace: ${model.namespace_id}`,
        'Copy renders for list item'
      );

    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-mfa-method-menu-link="details"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.method.index',
      'Details more menu action transitions to method route'
    );
    await click('.hds-breadcrumb a');
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-mfa-method-menu-link="edit"]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.method.edit',
      'Edit more menu action transitions to method edit route'
    );
  });

  test('it should display method details', async function (assert) {
    // ensure methods are tied to an enforcement
    this.server.get('/identity/mfa/login-enforcement', () => {
      const record = this.server.create('mfa-login-enforcement', {
        mfa_method_ids: this.getMethods().map((m) => m.id),
      });
      return {
        data: {
          key_info: { [record.name]: record },
          keys: [record.name],
        },
      };
    });
    await visit('/vault/access/mfa/methods');
    await click('[data-test-mfa-method-list-item]');
    assert.dom('[data-test-tab="config"]').hasClass('active', 'Configuration tab is active by default');
    await click('[data-test-delete-mfa-config]');

    assert
      .dom('[data-test-confirm-action-message]')
      .hasText(
        "This method cannot be deleted until its enforcements are deleted. This can be done from the 'Enforcements' tab."
      );

    const fields = [
      ['Issuer', 'Period', 'Key size', 'QR size', 'Algorithm', 'Digits', 'Skew', 'Max validation attempts'],
      ['Duo API hostname', 'Passcode reminder'],
      ['Organization name', 'Base URL'],
      ['Use signature', 'Idp url', 'Admin url', 'Authenticator url', 'Org alias'],
    ];
    for (const [index, labels] of fields.entries()) {
      if (index) {
        await click(`[data-test-mfa-method-list-item]:nth-of-type(${index + 2})`);
      }
      const url = currentURL();
      const id = url.slice(url.lastIndexOf('/') + 1);
      const model = this.store.peekRecord('mfa-method', id);

      labels.forEach((label) => {
        assert.dom(`[data-test-row-label="${label}"]`).hasText(label, `${label} field label renders`);
        const key =
          {
            'Duo API hostname': 'api_hostname',
            'Passcode reminder': 'use_passcode',
            'Organization name': 'org_name',
          }[label] || underscore(label);
        const value = typeof model[key] === 'boolean' ? (model[key] ? 'Yes' : 'No') : model[key].toString();
        assert.dom(`[data-test-value-div="${label}"]`).hasText(value, `${label} value renders`);
      });
      await click('.hds-breadcrumb a');
    }

    await click('[data-test-mfa-method-list-item]');
    await click('[data-test-mfa-method-edit]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.method.edit',
      'Toolbar action transitions to edit route'
    );
  });

  test('it should delete method that is not associated with any login enforcements', async function (assert) {
    this.server.get('/identity/mfa/login-enforcement', () => new Response(404, {}, { errors: [] }));

    await visit('/vault/access/mfa/methods');
    const methodCount = this.element.querySelectorAll('[data-test-mfa-method-list-item]').length;
    await click('[data-test-mfa-method-list-item]');
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');
    assert.dom('[data-test-mfa-method-list-item]').exists({ count: methodCount - 1 }, 'Method was deleted');
  });

  test('it should create methods', async function (assert) {
    assert.expect(12);

    await visit('/vault/access/mfa/methods');
    const methodCount = this.element.querySelectorAll('[data-test-mfa-method-list-item]').length;

    const methods = [
      { type: 'totp', required: ['issuer'] },
      { type: 'duo', required: ['secret_key', 'integration_key', 'api_hostname'] },
      { type: 'okta', required: ['org_name', 'api_token'] },
      { type: 'pingid', required: ['settings_file_base64'] },
    ];
    for (const [index, method] of methods.entries()) {
      const { type, required } = method;
      await click('[data-test-mfa-method-create]');
      await click(`[data-test-radio-card="${method.type}"]`);
      await click('[data-test-mfa-create-next]');
      await click('[data-test-mleh-radio="skip"]');
      await click('[data-test-mfa-create-save]');
      assert
        .dom('[data-test-inline-error-message]')
        .exists({ count: required.length }, `Required field validations display for ${type}`);

      for (const field of required) {
        await fillIn(GENERAL.inputByAttr(field), 'foo');
      }

      await click('[data-test-mfa-create-save]');
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.access.mfa.methods.method.index',
        `${type} method is displayed on save`
      );
      await click('.hds-breadcrumb a');
      assert
        .dom('[data-test-mfa-method-list-item]')
        .exists({ count: methodCount + index + 1 }, `List updates with new ${type} method`);
    }
  });

  test('it should create method with new enforcement', async function (assert) {
    await visit('/vault/access/mfa/methods/create');
    await click('[data-test-radio-card="totp"]');
    await click('[data-test-mfa-create-next]');
    await fillIn('[data-test-input="issuer"]', 'foo');
    await fillIn('[data-test-mlef-input="name"]', 'bar');
    await fillIn('[data-test-mount-accessor-select]', 'auth_userpass_bb95c2b1');
    await click('[data-test-mlef-add-target]');
    await click('[data-test-mfa-create-save]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.method.index',
      'Route transitions to method on save'
    );
    await click('[data-test-tab="enforcements"]');
    assert.dom('[data-test-list-item]').hasTextContaining('bar', 'Enforcement is listed in method view');
    await click('[data-test-sidebar-nav-link="Multi-Factor Authentication"]');
    await click('[data-test-tab="enforcements"]');
    assert
      .dom('[data-test-list-item="bar"]')
      .hasTextContaining('bar', 'Enforcement is listed in enforcements view');
    await click('[data-test-list-item="bar"]');
    await click('[data-test-tab="methods"]');
    assert
      .dom('[data-test-mfa-method-list-item]')
      .includesText('TOTP', 'TOTP method is listed in enforcement view');
  });

  test('it should create method and add it to existing enforcement', async function (assert) {
    await visit('/vault/access/mfa/methods/create');
    await click('[data-test-radio-card="totp"]');
    await click('[data-test-mfa-create-next]');
    await fillIn('[data-test-input="issuer"]', 'foo');
    await click('[data-test-mleh-radio="existing"]');
    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    const enforcement = this.element.querySelector('.ember-power-select-option');
    const name = enforcement.children[0].textContent.trim();
    await click(enforcement);
    await click('[data-test-mfa-create-save]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.methods.method.index',
      'Route transitions to method on save'
    );
    await click('[data-test-tab="enforcements"]');
    assert.dom('[data-test-list-item]').hasTextContaining(name, 'Enforcement is listed in method view');
  });

  test('it should edit methods', async function (assert) {
    await visit('/vault/access/mfa/methods');
    const id = this.element.querySelector('[data-test-mfa-method-list-item] .tag').textContent.trim();
    const model = this.store.peekRecord('mfa-method', id);
    await click('[data-test-mfa-method-list-item] [data-test-popup-menu-trigger]');
    await click('[data-test-mfa-method-menu-link="edit"]');

    const keys = ['issuer', 'period', 'key_size', 'qr_size', 'algorithm', 'digits', 'skew'];
    keys.forEach((key) => {
      if (key === 'period') {
        assert
          .dom('[data-test-ttl-value="Period"]')
          .hasValue(model.period.toString(), 'Period form field is populated with model value');
        assert.dom('[data-test-select="ttl-unit"]').hasValue('s', 'Correct time unit is shown for period');
      } else if (key === 'algorithm' || key === 'digits' || key === 'skew') {
        const radioElem = this.element.querySelector(`input[name=${key}]:checked`);
        assert
          .dom(radioElem)
          .hasValue(model[key].toString(), `${key} form field is populated with model value`);
      } else {
        assert
          .dom(`[data-test-input="${key}"]`)
          .hasValue(model[key].toString(), `${key} form field is populated with model value`);
      }
    });

    await fillIn('[data-test-input="issuer"]', 'foo');
    const SHA1radioBtn = this.element.querySelectorAll('input[name=algorithm]')[0];
    await click(SHA1radioBtn);
    await fillIn('[data-test-input="max_validation_attempts"]', 10);
    await click('[data-test-mfa-save]');
    await fillIn('[data-test-confirmation-modal-input]', model.type);
    await click('[data-test-confirm-button]');

    assert.dom('[data-test-row-value="Issuer"]').hasText('foo', 'Issuer field is updated');
    assert.dom('[data-test-row-value="Algorithm"]').hasText('SHA1', 'Algorithm field is updated');
    assert
      .dom('[data-test-row-value="Max validation attempts"]')
      .hasText('10', 'Max validation attempts field is updated');
  });

  test('it should navigate to enforcements create route from method enforcement tab', async function (assert) {
    await visit('/vault/access/mfa/methods');
    await click('[data-test-mfa-method-list-item]');
    await click('[data-test-tab="enforcements"]');
    await click('[data-test-enforcement-create]');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.mfa.enforcements.create',
      'Navigates to enforcements create route from toolbar action'
    );
  });
});
