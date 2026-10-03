/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { PKI_ROLE_DETAILS } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | pki role details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = 'pki';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.role = {
      name: 'Foobar',
      no_store: false,
      no_store_metadata: true,
      key_usage: [],
      ext_key_usage: ['bar', 'baz'],
      ttl: 600,
      issuer_ref: 'issuer-1',
    };
    this.capabilities = {
      canEdit: true,
      canDelete: true,
      canGenerateCert: true,
      canSign: true,
    };

    this.renderComponent = () =>
      render(
        hbs`
          <Page::PkiRoleDetails @role={{this.role}} @capabilities={{this.capabilities}} />
        `,
        { owner: this.engine }
      );
  });

  test('it should render the page component', async function (assert) {
    await this.renderComponent();

    assert.dom(PKI_ROLE_DETAILS.issuerLabel).hasText('Issuer', 'Label is');
    assert
      .dom(`${PKI_ROLE_DETAILS.keyUsageValue} [data-test-icon="minus"]`)
      .exists('Key usage shows dash when array is empty');
    assert
      .dom(PKI_ROLE_DETAILS.extKeyUsageValue)
      .hasText('bar,baz', 'Key usage shows comma-joined values when array has items');
    assert
      .dom(PKI_ROLE_DETAILS.noStoreValue)
      .containsText('Yes', 'noStore shows opposite of what the value is');
    assert
      .dom(PKI_ROLE_DETAILS.noStoreMetadataValue)
      .doesNotExist('does not render value for enterprise-only field');
    assert.dom(PKI_ROLE_DETAILS.customTtlValue).containsText('10 minutes', 'TTL shown as duration');
  });

  test('it should render the enterprise-only values in enterprise edition', async function (assert) {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';

    await this.renderComponent();
    assert
      .dom(PKI_ROLE_DETAILS.noStoreMetadataValue)
      .containsText('No', 'noStoreMetadata shows opposite of what the value is');
  });

  test('it should render the notAfter date if present', async function (assert) {
    assert.expect(1);

    this.role = {
      name: 'Foobar',
      no_store: false,
      key_usage: [],
      ext_key_usage: ['bar', 'baz'],
      not_after: '2030-05-04T12:00:00.000Z',
      issuer_ref: 'issuer-1',
    };

    await this.renderComponent();
    assert
      .dom(PKI_ROLE_DETAILS.customTtlValue)
      .containsText('May', 'Formats the notAfter date instead of TTL');
  });

  test('it should hide actions when user does not have capabilities', async function (assert) {
    this.capabilities = {
      canEdit: false,
      canDelete: false,
      canGenerateCert: false,
      canSign: false,
    };

    await this.renderComponent();

    assert.dom(PKI_ROLE_DETAILS.editRoleLink).doesNotExist('Edit link is not rendered');
    assert.dom(PKI_ROLE_DETAILS.deleteRoleButton).doesNotExist('Delete button is not rendered');
    assert.dom(PKI_ROLE_DETAILS.generateCertLink).doesNotExist('Generate Cert link is not rendered');
    assert.dom(PKI_ROLE_DETAILS.signCertLink).doesNotExist('Sign link is not rendered');
  });

  test('it should render actions when user has capabilities and delete role', async function (assert) {
    const deleteStub = sinon.stub(this.owner.lookup('service:api').secrets, 'pkiDeleteRole');

    await this.renderComponent();

    assert.dom(PKI_ROLE_DETAILS.editRoleLink).exists('Edit link renders');
    assert.dom(PKI_ROLE_DETAILS.deleteRoleButton).exists('Delete button renders');
    assert.dom(PKI_ROLE_DETAILS.generateCertLink).exists('Generate Cert link renders');
    assert.dom(PKI_ROLE_DETAILS.signCertLink).exists('Sign link renders');

    await click(PKI_ROLE_DETAILS.deleteRoleButton);
    await click(GENERAL.confirmButton);

    assert.true(
      deleteStub.calledWith(this.role.name, this.backend),
      'Delete API called with correct parameters'
    );
  });
});
