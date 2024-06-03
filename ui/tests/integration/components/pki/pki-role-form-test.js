/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_ROLE_FORM } from 'vault/tests/helpers/pki/pki-selectors';
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
    assert.dom(GENERAL.ttl.toggle('issuerRef-toggle')).exists();
    assert.dom(GENERAL.ttl.input('Backdate validity')).exists();
    assert.dom(GENERAL.fieldByAttr('customTtl')).exists();
    assert.dom(GENERAL.ttl.toggle('Max TTL')).exists();
    assert.dom(GENERAL.fieldByAttr('generateLease')).exists();
    assert.dom(GENERAL.fieldByAttr('noStore')).exists();
    assert
      .dom(GENERAL.fieldByAttr('noStoreMetadata'))
      .doesNotExist('noStoreMetadata is not shown b/c not enterprise');
    assert.dom(GENERAL.inputByAttr('addBasicConstraints')).exists();
    assert.dom(PKI_ROLE_FORM.domainHandling).exists('shows form-field group add domain handling');
    assert.dom(PKI_ROLE_FORM.keyParams).exists('shows form-field group key params');
    assert.dom(PKI_ROLE_FORM.keyUsage).exists('shows form-field group key usage');
    assert.dom(PKI_ROLE_FORM.policyIdentifiers).exists('shows form-field group policy identifiers');
    assert.dom(PKI_ROLE_FORM.san).exists('shows form-field group SAN');
    assert
      .dom(PKI_ROLE_FORM.additionalSubjectFields)
      .exists('shows form-field group additional subject fields');
  });

  test('it renders enterprise-only values in enterprise edition', async function (assert) {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';
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
    assert.dom(GENERAL.fieldByAttr('noStoreMetadata')).exists();
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

    await click(GENERAL.saveButton);

    assert
      .dom(GENERAL.inputByAttr('name'))
      .hasClass('has-error-border', 'shows border error on role name field when no role name is submitted');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('Name is required.', 'show correct error message');

    await fillIn(GENERAL.inputByAttr('name'), 'test-role');
    await click('[data-test-input="addBasicConstraints"]');
    await click(PKI_ROLE_FORM.domainHandling);
    await click('[data-test-input="allowedDomainsTemplate"]');
    await click(PKI_ROLE_FORM.policyIdentifiers);
    await fillIn('[data-test-input="policyIdentifiers"] [data-test-string-list-input="0"]', 'some-oid');
    await click(PKI_ROLE_FORM.san);
    await click('[data-test-input="allowUriSansTemplate"]');
    await click(PKI_ROLE_FORM.additionalSubjectFields);
    await fillIn(
      '[data-test-input="allowedSerialNumbers"] [data-test-string-list-input="0"]',
      'some-serial-number'
    );

    await click(GENERAL.saveButton);
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
    await click(GENERAL.ttl.toggle('issuerRef-toggle'));
    await fillIn(GENERAL.selectByAttr('issuerRef'), 'issuer-1');
    await click(GENERAL.saveButton);
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

    await click(GENERAL.ttl.toggle('issuerRef-toggle'));
    await fillIn(GENERAL.selectByAttr('issuerRef'), 'issuer-1');

    await click(PKI_ROLE_FORM.keyParams);
    assert.dom(GENERAL.inputByAttr('keyType')).hasValue('rsa');
    assert
      .dom(GENERAL.inputByAttr('keyBits'))
      .hasValue('3072', 'dropdown has model value, not default value (2048)');
    assert
      .dom(GENERAL.inputByAttr('signatureBits'))
      .hasValue('512', 'dropdown has model value, not default value (0)');

    await fillIn(GENERAL.inputByAttr('keyType'), 'ec');
    await fillIn(GENERAL.inputByAttr('keyBits'), '224');
    assert
      .dom(GENERAL.inputByAttr('keyBits'))
      .hasValue('224', 'dropdown has selected value, not default value (256)');
    await fillIn(GENERAL.inputByAttr('signatureBits'), '384');
    assert
      .dom(GENERAL.inputByAttr('signatureBits'))
      .hasValue('384', 'dropdown has selected value, not default value (0)');

    await click(GENERAL.saveButton);
    assert.strictEqual(this.role.issuerRef, 'issuer-1', 'Issuer Ref correctly saved on create');
  });
});
