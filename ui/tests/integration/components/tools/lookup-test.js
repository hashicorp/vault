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
import { format } from 'date-fns';

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
    // not stubbing the timestamp util here because this component uses the date-fns formatDistanceToNow method
    // so we need an actual now date for testing (which is why we don't assert the timestamp below, just the day of the month)
    const now = new Date();
    const data = {
      creation_path: 'sys/wrapping/wrap',
      creation_time: now.toISOString(),
      creation_ttl: 1800,
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
    assert.dom(GENERAL.infoRowValue('Creation path')).hasText(data.creation_path);
    assert.dom(GENERAL.infoRowValue('Creation time')).hasText(data.creation_time);
    assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText(`${data.creation_ttl}`);
    assert.dom(GENERAL.infoRowValue('Expiration date')).hasTextContaining(format(now, 'MMM dd yyyy')); // intentionally exclude time to avoid race conditions
    // remove below assertion if flaky (but unlikely this test would take longer than a minute..)
    assert.dom(GENERAL.infoRowValue('Expires in')).hasText('30 minutes'); // from 1800s ttl

    // clicking done resets form
    await click(TS.button('Done'));
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.button('Done')).doesNotExist();
  });
});
