/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import ControlGroupConfigForm from 'vault/forms/control-group-config';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | control-group-config-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const api = this.owner.lookup('service:api');
    this.writeStub = sinon.stub(api.sys, 'systemWriteConfigControlGroup').resolves();
    this.flash = sinon.stub(this.owner.lookup('service:flash-messages'), 'success');

    this.renderComponent = (data = { max_ttl: 0 }) => {
      this.form = new ControlGroupConfigForm(data);
      return render(hbs`<ControlGroupConfigForm @form={{this.form}} />`);
    };
  });

  test('it saves the max ttl and flashes success', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.toggleInput('Maximum TTL'));
    await fillIn(GENERAL.ttl.input('Maximum TTL'), '2');
    await fillIn(GENERAL.selectByAttr('ttl-unit'), 'h');
    await click(GENERAL.submitButton);

    assert.true(this.writeStub.calledOnce, 'write is called once');
    assert.strictEqual(this.writeStub.firstCall.args[0].max_ttl, '7200s', 'payload has max_ttl');
    assert.true(this.flash.calledWith('The configuration was saved successfully.'), 'success flash is shown');
  });

  test('it renders the api error and does not flash', async function (assert) {
    this.writeStub.rejects(getErrorResponse({ errors: ['permission denied'] }, 403));
    await this.renderComponent({ max_ttl: '1h' });
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.messageError).hasText('Error permission denied');
    assert.true(this.flash.notCalled, 'success flash is not shown');
  });
});
