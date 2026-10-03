/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kmip | Page::Credentials::Generate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    const data = { serial_number: '12345' };
    this.apiStub = sinon.stub(secrets, 'kmipGenerateClientCertificate').resolves({ data });
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.roleName = 'role-1';
    this.scopeName = 'scope-1';
    this.capabilities = { canDelete: true };

    this.renderComponent = () =>
      render(
        hbs`<Page::Credentials::Generate @scopeName={{this.scopeName}} @roleName={{this.roleName}} @capabilities={{this.capabilities}} />`,
        { owner: this.engine }
      );
  });

  test('it should render format field with correct options', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.fieldByAttr('format')).hasValue('pem', 'pem option is selected by default');
    assert.dom('[data-test-format="der"]').exists('Renders der option');
    assert.dom('[data-test-format="pem_bundle"]').exists('Renders pem_bundle option');
  });

  test('it should transition to index route on cancel', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.cancelButton);

    assert.true(this.routerStub.calledWith('vault.cluster.secrets.backend.kmip.credentials.index'));
  });

  test('it should generate credentials', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.true(
      this.apiStub.calledWith(this.roleName, this.scopeName, this.backend, { format: 'pem' }),
      'API called with default format'
    );
    assert.true(this.flashStub.calledWith(`Successfully generated credentials from role ${this.roleName}.`));

    await this.renderComponent();
    await fillIn(GENERAL.fieldByAttr('format'), 'der');
    await click(GENERAL.submitButton);
    assert.true(
      this.apiStub.calledWith(this.roleName, this.scopeName, this.backend, { format: 'pem' }),
      'API called with selected format'
    );
  });

  test('it should render details on generate success', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.fieldByAttr('format')).doesNotExist('Hides form on success');
    assert.dom(GENERAL.infoRowLabel('Serial number')).exists('Renders credentials details on success');
  });

  test('it should handle save error', async function (assert) {
    const error = 'Cannot generate credentials at this time';
    this.apiStub.rejects(getErrorResponse({ errors: [error] }, 400));

    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageDescription).hasText(error, 'Error message renders in alert banner');
    assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
  });
});
