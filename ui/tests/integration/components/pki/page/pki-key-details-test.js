/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_KEYS } from 'vault/tests/helpers/pki/pki-selectors';
import sinon from 'sinon';

module('Integration | Component | pki key details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'danger']);

    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.deleteStub = sinon.stub(this.owner.lookup('service:api').secrets, 'pkiDeleteKey').resolves();

    this.key = { key_id: '724862ff-6438-bad0-b598-77a6c7f4e934', key_type: 'ec', key_name: 'test-key' };
    this.canDelete = true;
    this.canEdit = true;

    this.renderComponent = () =>
      render(
        hbs`
        <Page::PkiKeyDetails
          @key={{this.key}}
          @canDelete={{this.canDelete}}
          @canEdit={{this.canEdit}}
        />
      `,
        { owner: this.engine }
      );
  });

  test('it renders the page component and deletes a key', async function (assert) {
    assert.expect(7);

    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowValue('Key ID'))
      .hasText(' 724862ff-6438-bad0-b598-77a6c7f4e934', 'key id renders');
    assert.dom(GENERAL.infoRowValue('Key name')).hasText('test-key', 'key name renders');
    assert.dom(GENERAL.infoRowValue('Key type')).hasText('ec', 'key type renders');
    assert.dom(GENERAL.infoRowLabel('Key bits')).doesNotExist('does not render empty value');
    assert.dom(PKI_KEYS.keyEditLink).exists('renders edit link');
    assert.dom(PKI_KEYS.keyDeleteButton).exists('renders delete button');
    await click(PKI_KEYS.keyDeleteButton);
    await click(GENERAL.confirmButton);
    assert.true(
      this.deleteStub.calledWith(this.key.key_id, this.backend),
      'pkiDeleteKey called with correct args'
    );
  });

  test('it does not render actions when capabilities are false', async function (assert) {
    assert.expect(2);

    this.canDelete = false;
    this.canEdit = false;

    await this.renderComponent();

    assert.dom(PKI_KEYS.keyDeleteButton).doesNotExist('does not render delete button if no permission');
    assert.dom(PKI_KEYS.keyEditLink).doesNotExist('does not render edit button if no permission');
  });

  test('it renders the private key as a <CertificateCard> component when there is a private key', async function (assert) {
    this.key.private_key = 'private-key-value';

    await this.renderComponent();
    assert.dom('[data-test-certificate-card]').exists('Certificate card renders for the private key');
  });
});
