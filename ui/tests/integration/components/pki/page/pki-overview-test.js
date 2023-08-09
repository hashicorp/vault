/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/overview';

module('Integration | Component | Page::PkiOverview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';

    this.store.createRecord('pki/issuer', { issuerId: 'abcd-efgh' });
    this.store.createRecord('pki/issuer', { issuerId: 'ijkl-mnop' });
    this.store.createRecord('pki/role', { name: 'role-0' });
    this.store.createRecord('pki/role', { name: 'role-1' });
    this.store.createRecord('pki/role', { name: 'role-2' });
    this.store.createRecord('pki/certificate/base', { serialNumber: '22:2222:22222:2222' });
    this.store.createRecord('pki/certificate/base', { serialNumber: '33:3333:33333:3333' });

    this.issuers = this.store.peekAll('pki/issuer');
    this.roles = this.store.peekAll('pki/role');
    this.engineId = 'pki';
  });

  test('shows the correct information on issuer card', async function (assert) {
    await render(
      hbs`<Page::PkiOverview @issuers={{this.issuers}} @roles={{this.roles}} @engine={{this.engineId}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.issuersCardTitle).hasText('Issuers');
    assert.dom(SELECTORS.issuersCardOverviewNum).hasText('2');
    assert.dom(SELECTORS.issuersCardLink).hasText('View issuers');
  });

  test('shows the correct information on roles card', async function (assert) {
    await render(
      hbs`<Page::PkiOverview @issuers={{this.issuers}} @roles={{this.roles}} @engine={{this.engineId}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rolesCardTitle).hasText('Roles');
    assert.dom(SELECTORS.rolesCardOverviewNum).hasText('3');
    assert.dom(SELECTORS.rolesCardLink).hasText('View roles');
    this.roles = 404;
    await render(
      hbs`<Page::PkiOverview @issuers={{this.issuers}} @roles={{this.roles}} @engine={{this.engineId}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rolesCardOverviewNum).hasText('0');
  });

  test('shows the input search fields for View Certificates card', async function (assert) {
    await render(
      hbs`<Page::PkiOverview @issuers={{this.issuers}} @roles={{this.roles}} @engine={{this.engineId}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.issueCertificate).hasText('Issue certificate');
    assert.dom(SELECTORS.issueCertificateInput).exists();
    assert.dom(SELECTORS.issueCertificateButton).hasText('Issue');
  });

  test('shows the input search fields for Issue Certificates card', async function (assert) {
    await render(
      hbs`<Page::PkiOverview @issuers={{this.issuers}} @roles={{this.roles}} @engine={{this.engineId}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.viewCertificate).hasText('View certificate');
    assert.dom(SELECTORS.viewCertificateInput).exists();
    assert.dom(SELECTORS.viewCertificateButton).hasText('View');
  });
});
