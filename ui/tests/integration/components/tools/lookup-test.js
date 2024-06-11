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
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import { formatRFC3339, getYear } from 'date-fns';
import sinon from 'sinon';

module('Integration | Component | tools/lookup', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Lookup />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Lookup Token', 'Title renders');
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/wrapping/lookup', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(TS.submit);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong', 'Error renders');
  });

  test('it submits', async function (assert) {
    const flashSuccessSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    // not stubbing the timestamp util here because this component uses the date-fns formatDistanceToNow method
    // so we need an actual now date for testing (which is why we don't assert the timestamp below, just the day of the month)
    const now = new Date();
    const data = {
      creation_path: 'sys/wrapping/wrap',
      creation_time: formatRFC3339(now),
      creation_ttl: 3.156e7, // one year in seconds
    };
    const token = 'token.OMZFbUurY0ppT2RTMGpRa0JOSUFqUzJUaGNqdWUQ6ooG';
    this.server.post('sys/wrapping/lookup', (schema, req) => {
      assert.propEqual(JSON.parse(req.requestBody), { token }, `payload trims token: ${req.requestBody}`);
      return { data };
    });
    await this.renderComponent();
    await fillIn(TS.toolsInput('wrapping-token'), `${token}   `);
    await click(TS.submit);

    await waitUntil(() => find(GENERAL.infoRowValue('Creation path')));
    assert.true(flashSuccessSpy.calledWith('Lookup was successful.'), 'it renders success flash');
    assert.dom(GENERAL.infoRowValue('Creation path')).hasText(data.creation_path);
    assert.dom(GENERAL.infoRowValue('Creation time')).hasText(data.creation_time);
    assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText(`${data.creation_ttl}`);
    assert.dom(GENERAL.infoRowValue('Expiration date')).hasTextContaining(`${getYear(now) + 1}`);
    assert.dom(GENERAL.infoRowValue('Expires in')).hasText('about 1 year');

    // clicking done resets form
    await click(TS.button('Done'));
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.button('Done')).doesNotExist();
  });
});
