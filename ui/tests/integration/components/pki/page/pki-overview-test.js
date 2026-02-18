/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { PKI_OVERVIEW } from 'vault/tests/helpers/pki/pki-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { overviewCard } = GENERAL;
module('Integration | Component | Page::PkiOverview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.secretMountPath = this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.issuers = ['abcd-efgh', 'ijkl-mnop'];
    this.roles = ['role-0', 'role-1', 'role-2'];
    this.certificates = ['22:2222:22222:2222', '33:3333:33333:3333'];
    this.engineId = 'pki';
    this.canListCertificates = true;
    this.canListRoles = true;

    this.renderComponent = () =>
      render(
        hbs`<Page::PkiOverview
          @issuers={{this.issuers}}
          @roles={{this.roles}}
          @certificates={{this.certificates}}
          @engine={{this.engineId}}
          @canListCertificates={{this.canListCertificates}}
          @canListRoles={{this.canListRoles}}
        />`,
        { owner: this.engine }
      );
  });

  test('shows the correct information on issuer card', async function (assert) {
    await this.renderComponent();
    assert
      .dom(overviewCard.container('Issuers'))
      .hasText(
        'Issuers View issuers The total number of issuers in this PKI mount. Includes both root and intermediate certificates. 2'
      );

    this.issuers = [];
    await this.renderComponent();
    assert
      .dom(overviewCard.container('Issuers'))
      .hasText(
        'Issuers View issuers The total number of issuers in this PKI mount. Includes both root and intermediate certificates. 0'
      );
  });

  test('shows the correct information on roles card', async function (assert) {
    await this.renderComponent();
    assert
      .dom(overviewCard.container('Roles'))
      .hasText(
        'Roles View roles The total number of roles in this PKI mount that have been created to generate certificates. 3'
      );
    this.roles = [];
    await this.renderComponent();
    assert
      .dom(overviewCard.container('Roles'))
      .hasText(
        'Roles View roles The total number of roles in this PKI mount that have been created to generate certificates. 0'
      );
  });

  test('shows the search select dropdown for View Certificates card', async function (assert) {
    await this.renderComponent();
    assert.dom(overviewCard.title('View certificate')).hasText('View certificate');
    assert
      .dom(overviewCard.description('View certificate'))
      .hasText('Quickly view a certificate by looking up its serial number.');
    assert.dom(PKI_OVERVIEW.viewCertificateInput).exists();
    assert.dom(GENERAL.inputSearch('certificate')).doesNotExist('it does not render certificate input');
    assert.dom(PKI_OVERVIEW.viewCertificateButton).hasText('View');
  });

  test('shows the search select dropdown for Issue Certificates card', async function (assert) {
    await this.renderComponent();
    assert.dom(overviewCard.title('Issue certificate')).hasText('Issue certificate');
    assert
      .dom(overviewCard.description('Issue certificate'))
      .hasText('Begin issuing a certificate by choosing a role.');
    assert.dom(PKI_OVERVIEW.issueCertificateInput).exists();
    assert.dom(GENERAL.inputSearch('role')).doesNotExist('it does not render role input');
    assert.dom(PKI_OVERVIEW.issueCertificateButton).hasText('Issue');
  });

  test('it renders manual search inputs when no list permission', async function (assert) {
    this.canListCertificates = false;
    this.canListRoles = false;
    await this.renderComponent();
    assert.dom(overviewCard.container('Roles')).doesNotExist();
    assert
      .dom(overviewCard.description('View certificate'))
      .hasText('Quickly view a certificate by providing its serial number.');
    assert
      .dom(overviewCard.description('Issue certificate'))
      .hasText('Begin issuing a certificate by entering a role.');
    assert.dom(PKI_OVERVIEW.issueCertificateInput).doesNotExist('role search select does not render');
    assert.dom(GENERAL.inputSearch('role')).exists('it renders input instead of search select');
    assert.dom(PKI_OVERVIEW.viewCertificateInput).doesNotExist('certificate search select does not render');
    assert.dom(GENERAL.inputSearch('certificate')).exists('it renders input instead of search selects');
  });
});
