/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

/**
 * this test is for the page component only. A separate test is written for the form rendered
 */
module('Integration | Component | page/pki-issuer-generate-intermediate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.breadcrumbs = [{ label: 'something' }];
    this.backend = 'pki-component';
    this.secretMountPath = this.owner.lookup('service:secret-mount-path').update(this.backend);

    this.generateStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'pkiIssuersGenerateIntermediate')
      .resolves();
    this.capabilitiesStub = sinon
      .stub(this.owner.lookup('service:capabilities'), 'for')
      .resolves({ canCreate: true });

    this.renderComponent = () =>
      render(hbs`<Page::PkiIssuerGenerateIntermediate @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
  });

  test('it renders correct title before and after submit', async function (assert) {
    assert.expect(4);

    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Generate intermediate CSR');

    const { backend } = this;
    const type = 'internal';
    await fillIn(GENERAL.inputByAttr('type'), type);
    await fillIn(GENERAL.inputByAttr('common_name'), 'foobar');
    await click('[data-test-submit]');

    const payload = {
      common_name: 'foobar',
      format: 'pem',
      key_type: 'rsa',
      not_before_duration: 30,
      private_key_format: 'der',
    };
    assert.true(
      this.capabilitiesStub.calledWith('pkiIssuersGenerateIntermediate', { backend, type }),
      'Capabilities checked for api path'
    );
    assert.true(this.generateStub.calledWith(type, backend, payload), 'API called with correct params');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('View Generated CSR');
  });

  test('it does not update title if API response is an error', async function (assert) {
    assert.expect(2);

    this.generateStub.rejects(getErrorResponse({ errors: ['API returns this error'] }, 403));

    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Generate intermediate CSR');
    // Fill in
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('common_name'), 'foobar');
    await click('[data-test-submit]');
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('Generate intermediate CSR', 'title does not change if response is unsuccessful');
  });
});
