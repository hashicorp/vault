/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-role-form';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | pki-role-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.role = this.store.createRecord('pki/role');
    this.store.createRecord('pki/issuer', { issuerName: 'issuer-0', issuerId: 'abcd-efgh' });
    this.store.createRecord('pki/issuer', { issuerName: 'issuer-1', issuerId: 'ijkl-mnop' });
    this.issuers = this.store.peekAll('pki/issuer');
    this.role.backend = 'pki';
    this.onCancel = sinon.spy();
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should render default fields and toggle groups', async function (assert) {
    assert.expect(13);
    await render(
      hbs`
      <PkiRoleForm
         @role={{this.role}}
         @issuers={{this.issuers}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.issuerRefToggle).exists('shows issuer ref toggle');
    assert.dom(SELECTORS.backdateValidity).exists('shows form-field backdate validity');
    assert.dom(SELECTORS.customTtl).exists('shows custom yielded form field');
    assert.dom(SELECTORS.maxTtl).exists('shows form-field max ttl');
    assert.dom(SELECTORS.generateLease).exists('shows form-field generateLease');
    assert.dom(SELECTORS.noStore).exists('shows form-field no store');
    assert.dom(SELECTORS.addBasicConstraints).exists('shows form-field add basic constraints');
    assert.dom(SELECTORS.domainHandling).exists('shows form-field group add domain handling');
    assert.dom(SELECTORS.keyParams).exists('shows form-field group key params');
    assert.dom(SELECTORS.keyUsage).exists('shows form-field group key usage');
    assert.dom(SELECTORS.policyIdentifiers).exists('shows form-field group policy identifiers');
    assert.dom(SELECTORS.san).exists('shows form-field group SAN');
    assert.dom(SELECTORS.additionalSubjectFields).exists('shows form-field group additional subject fields');
  });

  test('it should save a new pki role with various options selected', async function (assert) {
    // Key usage, Key params and Not valid after options are tested in their respective component tests
    assert.expect(8);
    this.server.post(`/${this.role.backend}/roles/test-role`, (schema, req) => {
      assert.ok(true, 'Request made to save role');
      const request = JSON.parse(req.requestBody);
      const allowedDomainsTemplate = request.allowed_domains_template;
      const policyIdentifiers = request.policy_identifiers;
      const allowedUriSansTemplate = request.allow_uri_sans_template;
      const allowedSerialNumbers = request.allowed_serial_numbers;

      assert.true(allowedDomainsTemplate, 'correctly sends allowed_domains_template');
      assert.strictEqual(policyIdentifiers[0], 'some-oid', 'correctly sends policy_identifiers');
      assert.true(allowedUriSansTemplate, 'correctly sends allowed_uri_sans_template');
      assert.strictEqual(
        allowedSerialNumbers[0],
        'some-serial-number',
        'correctly sends allowed_serial_numbers'
      );
      return {};
    });

    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(
      hbs`
      <PkiRoleForm
         @role={{this.role}}
         @issuers={{this.issuers}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
  `,
      { owner: this.engine }
    );

    await click(SELECTORS.roleCreateButton);

    assert
      .dom(SELECTORS.roleName)
      .hasClass('has-error-border', 'shows border error on role name field when no role name is submitted');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('Name is required.', 'show correct error message');

    await fillIn(SELECTORS.roleName, 'test-role');
    await click('[data-test-input="addBasicConstraints"]');
    await click(SELECTORS.domainHandling);
    await click('[data-test-input="allowedDomainsTemplate"]');
    await click(SELECTORS.policyIdentifiers);
    await fillIn('[data-test-input="policyIdentifiers"] [data-test-string-list-input="0"]', 'some-oid');
    await click(SELECTORS.san);
    await click('[data-test-input="allowUriSansTemplate"]');
    await click(SELECTORS.additionalSubjectFields);
    await fillIn(
      '[data-test-input="allowedSerialNumbers"] [data-test-string-list-input="0"]',
      'some-serial-number'
    );

    await click(SELECTORS.roleCreateButton);
  });

  test('it should update attributes on the model on update', async function (assert) {
    assert.expect(1);

    this.store.pushPayload('pki/role', {
      modelName: 'pki/role',
      name: 'test-role',
      backend: 'pki-test',
      id: 'role-id',
    });

    this.role = this.store.peekRecord('pki/role', 'role-id');

    await render(
      hbs`
      <PkiRoleForm
        @role={{this.role}}
        @issuers={{this.issuers}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      `,
      { owner: this.engine }
    );
    await click(SELECTORS.issuerRefToggle);
    await fillIn(SELECTORS.issuerRefSelect, 'issuer-1');
    await click(SELECTORS.roleCreateButton);
    assert.strictEqual(this.role.issuerRef, 'issuer-1', 'Issuer Ref correctly saved on create');
  });

  test('it should edit a role', async function (assert) {
    assert.expect(8);
    this.server.post(`/pki-test/roles/test-role`, (schema, req) => {
      assert.ok(true, 'Request made to correct endpoint to update role');
      const request = JSON.parse(req.requestBody);
      assert.propEqual(
        request,
        {
          allow_ip_sans: true,
          issuer_ref: 'issuer-1',
          key_bits: '224',
          key_type: 'ec',
          key_usage: ['DigitalSignature', 'KeyAgreement', 'KeyEncipherment'],
          not_before_duration: '30s',
          require_cn: true,
          signature_bits: '384',
          use_csr_common_name: true,
          use_csr_sans: true,
        },
        'sends role params in correct type'
      );
      return {};
    });

    this.store.pushPayload('pki/role', {
      modelName: 'pki/role',
      name: 'test-role',
      backend: 'pki-test',
      id: 'role-id',
      key_type: 'rsa',
      key_bits: 3072, // string type in dropdown, API returns as numbers
      signature_bits: 512, // string type in dropdown, API returns as numbers
    });

    this.role = this.store.peekRecord('pki/role', 'role-id');

    await render(
      hbs`
      <PkiRoleForm
        @role={{this.role}}
        @issuers={{this.issuers}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      `,
      { owner: this.engine }
    );

    await click(SELECTORS.issuerRefToggle);
    await fillIn(SELECTORS.issuerRefSelect, 'issuer-1');

    await click(SELECTORS.keyParams);
    assert.dom(SELECTORS.keyType).hasValue('rsa');
    assert.dom(SELECTORS.keyBits).hasValue('3072', 'dropdown has model value, not default value (2048)');
    assert.dom(SELECTORS.signatureBits).hasValue('512', 'dropdown has model value, not default value (0)');

    await fillIn(SELECTORS.keyType, 'ec');
    await fillIn(SELECTORS.keyBits, '224');
    assert.dom(SELECTORS.keyBits).hasValue('224', 'dropdown has selected value, not default value (256)');
    await fillIn(SELECTORS.signatureBits, '384');
    assert.dom(SELECTORS.signatureBits).hasValue('384', 'dropdown has selected value, not default value (0)');

    await click(SELECTORS.roleCreateButton);
    assert.strictEqual(this.role.issuerRef, 'issuer-1', 'Issuer Ref correctly saved on create');
  });
});
