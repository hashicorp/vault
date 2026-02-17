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

module('Integration | Component | kmip | Page::Scopes::Create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(secrets, 'kmipCreateScope').resolves();
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.indexRoute = `vault.cluster.secrets.backend.kmip.scopes.index`;

    this.renderComponent = () => render(hbs`<Page::Scopes::Create />`, { owner: this.engine });
  });

  test('it should render error when name is not defined', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required', 'Shows validation error for name field');
    assert.dom(GENERAL.inlineError).hasText('There is an error with this form.');
  });

  test('it should transition to index route on cancel', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.cancelButton);

    assert.true(this.routerStub.calledWith(this.indexRoute));
  });

  test('it should save scope', async function (assert) {
    await this.renderComponent();

    this.name = 'new-scope';
    await fillIn(GENERAL.inputByAttr('name'), this.name);
    await click(GENERAL.submitButton);

    assert.true(this.apiStub.calledWith(this.name, this.backend, {}));
    assert.true(this.flashStub.calledWith(`Successfully created scope ${this.name}`));
    assert.true(this.routerStub.calledWith(this.indexRoute));
  });

  test('it should handle save error', async function (assert) {
    this.name = 'new-scope';
    const error = `scope "${this.name}" already exists`;
    this.apiStub.rejects(getErrorResponse({ errors: [error] }, 400));

    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), this.name);
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.messageDescription).hasText(error, 'Error message renders in alert banner');
    assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
  });
});
