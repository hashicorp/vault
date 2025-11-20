/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { render } from '@ember/test-helpers';
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

  hooks.beforeEach(function () {
    this.backend = 'pki-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.engineId = 'pki';

    this.renderComponent = () =>
      render(
        hbs`<Page::PkiIssuerList @backend={{this.backend}} @issuers={{this.issuers}} @mountPoint={{this.engineId}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it renders when issuer metadata tags when info is provided', async function (assert) {
    this.issuers = [
      {
        issuer_id: 'abcd-efgh',
        issuer_name: 'issuer-0',
        common_name: 'common-name-issuer-0',
        isRoot: true,
        serial_number: '74:2d:ed:f2:c4:3b:76:5e:6e:0d:f1:6a:c0:8b:6f:e3:3c:62:f9:03',
      },
      {
        issuer_id: 'ijkl-mnop',
        issuer_name: 'issuer-1',
        isRoot: false,
        parsedCertificate: {
          common_name: 'common-name-issuer-1',
        },
        serial_number: '74:2d:ed:f2:c4:3b:76:5e:6e:0d:f1:6a:c0:8b:6f:e3:3c:62:f9:03',
      },
    ];
    this.issuers.meta = STANDARD_META;

    await this.renderComponent();

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
    this.issuers = [
      {
        issuer_id: 'abcd-efgh',
        issuer_name: 'issuer-0',
        is_default: false,
      },
      {
        issuer_id: 'ijkl-mnop',
        issuer_name: 'issuer-1',
        is_default: true,
      },
    ];
    this.issuers.meta = STANDARD_META;

    await this.renderComponent();
    assert.dom(`[data-test-is-default="1"]`).hasText('default issuer');
  });
});
