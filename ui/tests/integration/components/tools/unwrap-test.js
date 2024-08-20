/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { click, fillIn, find, render, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror from 'vault/tests/helpers/codemirror';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import sinon from 'sinon';

module('Integration | Component | tools/unwrap', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Unwrap />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Unwrap Data', 'Title renders');
    assert.dom(TS.submit).hasText('Unwrap data');
    assert.dom(TS.toolsInput('unwrap-token')).hasValue('');
    assert.dom(GENERAL.hdsTab('data')).doesNotExist();
    assert.dom(GENERAL.hdsTab('details')).doesNotExist();
    assert.dom('.CodeMirror').doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/wrapping/unwrap', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong', 'Error renders');
  });

  test('it submits and renders falsy values', async function (assert) {
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const unwrapData = { foo: 'bar' };
    const data = { token: 'token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG' };
    const expectedDetails = {
      'Request ID': '291290a6-5602-e49a-389b-5870e6c02976',
      'Lease ID': 'None',
      Renewable: 'No',
      'Lease Duration': 'None',
    };
    this.server.post('sys/wrapping/unwrap', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, data, `payload contains token: ${req.requestBody}`);
      return {
        data: unwrapData,
        lease_duration: 0,
        lease_id: '',
        renewable: false,
        request_id: '291290a6-5602-e49a-389b-5870e6c02976',
      };
    });

    await this.renderComponent();

    // test submit
    await fillIn(TS.toolsInput('unwrap-token'), data.token);
    await click(TS.submit);

    await waitUntil(() => find('.CodeMirror'));
    assert.true(flashSpy.calledWith('Unwrap was successful.'), 'it renders success flash');
    assert.dom('label').hasText('Unwrapped Data');
    assert.strictEqual(codemirror().getValue(' '), '{   "foo": "bar" }', 'it renders unwrapped data');
    assert.dom(GENERAL.hdsTab('data')).hasAttribute('aria-selected', 'true');

    await click(GENERAL.hdsTab('details'));
    assert.dom(GENERAL.hdsTab('details')).hasAttribute('aria-selected', 'true');
    assert
      .dom(`${GENERAL.infoRowValue('Renewable')} ${GENERAL.icon('x-square')}`)
      .exists('renders falsy icon for renewable');
    for (const detail in expectedDetails) {
      assert.dom(GENERAL.infoRowValue(detail)).hasText(expectedDetails[detail]);
    }

    // form resets clicking 'Done'
    await click(TS.button('Done'));
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('unwrap-token')).hasValue('', 'token input resets');
  });

  test('it submits and renders truthy values', async function (assert) {
    const unwrapData = { foo: 'bar' };
    const data = { token: 'token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG' };
    const expectedDetails = {
      'Request ID': '291290a6-5602-e49a-389b-5870e6c02976',
      'Lease ID': '123',
      Renewable: 'Yes',
      'Lease Duration': '1800',
    };
    this.server.post('sys/wrapping/unwrap', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, data, `payload contains token: ${req.requestBody}`);
      return {
        data: unwrapData,
        lease_duration: 1800,
        lease_id: '123',
        renewable: true,
        request_id: '291290a6-5602-e49a-389b-5870e6c02976',
      };
    });

    await this.renderComponent();

    await fillIn(TS.toolsInput('unwrap-token'), data.token);
    await click(TS.submit);

    await waitUntil(() => find('.CodeMirror'));
    await click(GENERAL.hdsTab('details'));
    assert
      .dom(`${GENERAL.infoRowValue('Renewable')} ${GENERAL.icon('check-circle')}`)
      .exists('renders truthy icon for renewable');
    for (const detail in expectedDetails) {
      assert.dom(GENERAL.infoRowValue(detail)).hasText(expectedDetails[detail]);
    }
  });

  test('it trims token whitespace', async function (assert) {
    const data = { token: 'token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG' };
    this.server.post('sys/wrapping/unwrap', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, data, `token does not include whitespace: "${req.requestBody}"`);
      return {};
    });

    await this.renderComponent();

    await fillIn(TS.toolsInput('unwrap-token'), `${data.token}  `);
    await click(TS.submit);
  });
});
