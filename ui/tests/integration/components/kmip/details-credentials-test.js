/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, render, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kmip | DetailsCredentials', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.roleName = 'role-1';
    this.scopeName = 'scope-1';
    this.capabilities = { canDelete: true };
    this.credentials = {
      ca_chain: [
        '-----BEGIN CERTIFICATE-----ca1-----END CERTIFICATE-----',
        '-----BEGIN CERTIFICATE-----ca2-----END CERTIFICATE-----',
      ],
      certificate: '-----BEGIN CERTIFICATE-----certificate-----END CERTIFICATE-----',
      private_key: '-----BEGIN EC PRIVATE KEY-----private key-----END EC PRIVATE KEY-----',
      serial_number: '609451918007712020381976412167930834908821189113',
    };

    this.apiStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'kmipRevokeClientCertificate')
      .resolves();
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.renderComponent = () =>
      render(
        hbs`<DetailsCredentials @roleName={{this.roleName}} @scopeName={{this.scopeName}} @credentials={{this.credentials}} @capabilities={{this.capabilities}} />`,
        { owner: this.engine }
      );
  });

  test('it should render/hide toolbar actions based on capabilities', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.confirmTrigger)
      .hasText('Revoke credentials', 'Revoke credentials action renders in toolbar');
    assert
      .dom('[data-test-copy-button]')
      .containsText('Copy certificate', 'Copy certificate action renders in toolbar');
    assert.dom('[data-test-back-to-role]').hasText('Back to role', 'Back to role action renders in toolbar');

    this.capabilities = { canDelete: false };
    await this.renderComponent();

    assert.dom(GENERAL.confirmTrigger).doesNotExist('Revoke credentials action is hidden without capability');
  });

  test('it should render credentials info', async function (assert) {
    await this.renderComponent();

    await click(`${GENERAL.infoRowValue('Serial number')} ${GENERAL.button('toggle-masked')}`);
    assert
      .dom(GENERAL.infoRowValue('Serial number'))
      .containsText(this.credentials.serial_number, 'Serial Number renders');

    assert
      .dom(`${GENERAL.infoRowValue('Private key')} [data-test-warning]`)
      .containsText(
        'You will not be able to access the private key later, so please copy the information below.',
        'Private key warning message renders'
      );
    await click(`${GENERAL.infoRowValue('Private key')} ${GENERAL.button('toggle-masked')}`);
    assert
      .dom(GENERAL.infoRowValue('Private key'))
      .containsText(this.credentials.private_key, 'Private key renders');

    await click(`${GENERAL.infoRowValue('Certificate')} ${GENERAL.button('toggle-masked')}`);
    assert
      .dom(GENERAL.infoRowValue('Certificate'))
      .containsText(this.credentials.certificate, 'Certificate renders');

    const caChainMaskBtns = findAll(`${GENERAL.infoRowValue('CA Chain')} ${GENERAL.button('toggle-masked')}`);
    await click(caChainMaskBtns[0]);
    assert
      .dom(GENERAL.infoRowValue('CA Chain'))
      .containsText(this.credentials.ca_chain[0], 'First CA Chain renders');
    await click(caChainMaskBtns[0]);
    await click(caChainMaskBtns[1]);
    assert
      .dom(GENERAL.infoRowValue('CA Chain'))
      .containsText(this.credentials.ca_chain[1], 'Second CA Chain renders');
  });

  test('it should hide private_key when not available', async function (assert) {
    this.credentials.private_key = undefined;
    await this.renderComponent();
    assert.dom(GENERAL.infoRowLabel('Private key')).doesNotExist('Private key is hidden');
  });

  test('it should revoke credentials', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    const { serial_number } = this.credentials;
    assert.true(
      this.apiStub.calledWith(this.roleName, this.scopeName, this.backend, { serial_number }),
      'API called to revoke credentials'
    );
    assert.true(
      this.flashStub.calledWith(`Successfully revoked credentials.`),
      'Success flash message shown'
    );
    assert.true(
      this.routerStub.calledWith(
        'vault.cluster.secrets.backend.kmip.credentials.index',
        this.scopeName,
        this.roleName
      ),
      'Transitions to credentials list on revoke success'
    );
  });
});
