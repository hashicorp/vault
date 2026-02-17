/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import PkiRoleForm from 'vault/forms/secrets/pki/role';

module('Integration | Component | pki-role-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.backend = 'pki';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'pkiWriteRole');

    this.issuers = [
      { issuer_name: 'issuer-0', issuer_id: 'abcd-efgh' },
      { issuer_name: 'issuer-1', issuer_id: 'ijkl-mnop' },
    ];
    this.onCancel = sinon.spy();
    this.onSave = sinon.spy();

    this.formDefaults = {
      allow_ip_sans: true,
      allow_localhost: true,
      client_flag: true,
      enforce_hostnames: true,
      key_usage: ['DigitalSignature', 'KeyAgreement', 'KeyEncipherment'],
      not_before_duration: 30,
      serial_number_source: 'json-csr',
      server_flag: true,
      use_csr_common_name: true,
      use_csr_sans: true,
    };

    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });

    this.renderComponent = () => {
      this.form = new PkiRoleForm(this.role, { isNew: !this.role });
      return render(
        hbs`<PkiRoleForm @form={{this.form}} @issuers={{this.issuers}} @onCancel={{this.onCancel}} @onSave={{this.onSave}} />`,
        { owner: this.engine }
      );
    };
  });

  test('it should render default fields and toggle groups', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.ttl.toggle('issuer_ref-toggle')).exists();
    assert.dom(GENERAL.ttl.input('Backdate validity')).exists();
    assert.dom(GENERAL.fieldByAttr('customTtl')).exists();
    assert.dom(GENERAL.ttl.toggle('Max TTL')).exists();
    assert.dom(GENERAL.fieldByAttr('generate_lease')).exists();
    assert.dom(GENERAL.fieldByAttr('no_store')).exists();
    assert
      .dom(GENERAL.fieldByAttr('no_store_metadata'))
      .doesNotExist('no_store_metadata is not shown for community edition');
    assert.dom(GENERAL.inputByAttr('basic_constraints_valid_for_non_ca')).exists();
    assert.dom(GENERAL.button('Domain handling')).exists('shows form-field group add domain handling');
    assert.dom(GENERAL.button('Key parameters')).exists('shows form-field group key params');
    assert.dom(GENERAL.button('Key usage')).exists('shows form-field group key usage');
    assert.dom(GENERAL.button('Policy identifiers')).exists('shows form-field group policy identifiers');
    assert.dom(GENERAL.button('Subject Alternative Name (SAN) Options')).exists('shows form-field group SAN');
    assert
      .dom(GENERAL.button('Additional subject fields'))
      .exists('shows form-field group additional subject fields');
  });

  test('it renders enterprise-only values in enterprise edition', async function (assert) {
    this.owner.lookup('service:version').type = 'enterprise';
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('no_store_metadata')).exists();
  });

  test('it should save a new pki role with various options selected', async function (assert) {
    // Key usage, Key params and Not valid after options are tested in their respective component tests
    assert.expect(3);

    await this.renderComponent();

    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .includesText('Name is required.', 'show correct error message');

    const name = 'test-role';
    await fillIn(GENERAL.inputByAttr('name'), name);
    await click('[data-test-input="basic_constraints_valid_for_non_ca"]');
    await click(GENERAL.button('Domain handling'));
    await click('[data-test-input="allowed_domains_template"]');
    await click(GENERAL.button('Policy identifiers'));
    await fillIn('[data-test-input="policy_identifiers"] [data-test-string-list-input="0"]', 'some-oid');
    await click(GENERAL.button('Subject Alternative Name (SAN) Options'));
    await click('[data-test-input="allowed_uri_sans_template"]');
    await click(GENERAL.button('Additional subject fields'));
    await fillIn(
      '[data-test-input="allowed_serial_numbers"] [data-test-string-list-input="0"]',
      'some-serial-number'
    );

    await click(GENERAL.submitButton);
    const payload = {
      ...this.formDefaults,
      allowed_domains_template: true,
      basic_constraints_valid_for_non_ca: true,
      policy_identifiers: ['some-oid'],
      allowed_uri_sans_template: true,
      allowed_serial_numbers: ['some-serial-number'],
    };
    assert.true(
      this.writeStub.calledWith(name, this.backend, payload),
      'Correct endpoint is called to save role'
    );
    assert.true(this.onSave.calledWith(name), 'onSave called with role name after successful save');
  });

  test('it should edit a role', async function (assert) {
    assert.expect(7);

    this.role = {
      name: 'test-role',
      issuer_ref: 'default',
      key_type: 'rsa',
      key_bits: 3072, // string type in dropdown, API returns as numbers
      signature_bits: 512, // string type in dropdown, API returns as numbers
    };
    await this.renderComponent();

    await click(GENERAL.ttl.toggle('issuer_ref-toggle'));
    await fillIn(GENERAL.selectByAttr('issuer_ref'), 'issuer-1');

    await click(GENERAL.button('Key parameters'));
    assert.dom(GENERAL.inputByAttr('key_type')).hasValue('rsa');
    assert
      .dom(GENERAL.inputByAttr('key_bits'))
      .hasValue('3072', 'dropdown has model value, not default value (2048)');
    assert
      .dom(GENERAL.inputByAttr('signature_bits'))
      .hasValue('512', 'dropdown has model value, not default value (0)');

    await fillIn(GENERAL.inputByAttr('key_type'), 'ec');
    await fillIn(GENERAL.inputByAttr('key_bits'), '224');
    assert
      .dom(GENERAL.inputByAttr('key_bits'))
      .hasValue('224', 'dropdown has selected value, not default value (256)');
    await fillIn(GENERAL.inputByAttr('signature_bits'), '384');
    assert
      .dom(GENERAL.inputByAttr('signature_bits'))
      .hasValue('384', 'dropdown has selected value, not default value (0)');

    await click(GENERAL.submitButton);

    const payload = {
      ...this.formDefaults,
      issuer_ref: 'issuer-1',
      key_bits: '224',
      key_type: 'ec',
      signature_bits: '384',
    };
    assert.true(
      this.writeStub.calledWith('test-role', this.backend, payload),
      'Correct endpoint is called to save role'
    );
    assert.true(this.onSave.calledWith('test-role'), 'onSave called with role name after successful save');
  });
});
