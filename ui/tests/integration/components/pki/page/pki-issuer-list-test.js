/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { render } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupRenderingTest } from 'vault/tests/helpers';
import { STANDARD_META } from 'vault/tests/helpers/pagination';

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-list', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.engineId = 'pki';
  });

  test('it renders when issuer metadata tags when info is provided', async function (assert) {
    this.store.createRecord('pki/issuer', {
      issuerId: 'abcd-efgh',
      issuerName: 'issuer-0',
      common_name: 'common-name-issuer-0',
      isRoot: true,
      serialNumber: '74:2d:ed:f2:c4:3b:76:5e:6e:0d:f1:6a:c0:8b:6f:e3:3c:62:f9:03',
    });
    this.store.createRecord('pki/issuer', {
      issuerId: 'ijkl-mnop',
      issuerName: 'issuer-1',
      isRoot: false,
      parsedCertificate: {
        common_name: 'common-name-issuer-1',
      },
      serialNumber: '74:2d:ed:f2:c4:3b:76:5e:6e:0d:f1:6a:c0:8b:6f:e3:3c:62:f9:03',
    });
    const issuers = this.store.peekAll('pki/issuer');
    issuers.meta = STANDARD_META;
    this.issuers = issuers;

    await render(
      hbs`<Page::PkiIssuerList @backend="pki-mount" @issuers={{this.issuers}} @mountPoint={{this.engineId}} />`,
      {
        owner: this.engine,
      }
    );

    this.issuers.forEach(async (issuer, idx) => {
      assert
        .dom(`[data-test-serial-number="${idx}"]`)
        .hasText('74:2d:ed:f2:c4:3b:76:5e:6e:0d:f1:6a:c0:8b:6f:e3:3c:62:f9:03');
      if (idx === 1) {
        assert.dom(`[data-test-is-root-tag="${idx}"]`).hasText('intermediate');
      } else {
        assert.dom(`[data-test-is-root-tag="${idx}"]`).hasText('root');
      }
    });
  });
  test('it renders when issuer data even though issuer metadata isnt provided', async function (assert) {
    this.store.createRecord('pki/issuer', {
      issuerId: 'abcd-efgh',
      issuerName: 'issuer-0',
      isDefault: false,
    });
    this.store.createRecord('pki/issuer', {
      issuerId: 'ijkl-mnop',
      issuerName: 'issuer-1',
      isDefault: true,
    });
    const issuers = this.store.peekAll('pki/issuer');
    issuers.meta = STANDARD_META;
    this.issuers = issuers;
    await render(
      hbs`<Page::PkiIssuerList @backend="pki-mount" @issuers={{this.issuers}} @mountPoint={{this.engineId}} />`,
      {
        owner: this.engine,
      }
    );
    assert.dom(`[data-test-is-default="1"]`).hasText('default issuer');
  });
});
