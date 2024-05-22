/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupModels } from 'vault/tests/helpers/sync/setup-models';
import hbs from 'htmlbars-inline-precompile';
import { render, click, fillIn, settled } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { selectChoose } from 'ember-power-select/test-support';
import sinon from 'sinon';
import { Response } from 'miragejs';

const { destinations, searchSelect, messageError, kvSuggestion } = PAGE;
const { mountSelect, mountInput, submit, cancel, successMessage } = destinations.sync;

module('Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupModels(hooks);

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.server.get('/sys/internal/ui/mounts', () => ({
      data: { secret: { 'my-kv/': { type: 'kv', options: { version: '2' } } } },
    }));
    this.server.get('/my-kv/metadata', () => ({
      data: { keys: ['my-path/', 'my-secret'] },
    }));
    this.server.get('/my-kv/metadata/my-path', () => ({
      data: { keys: ['nested-secret'] },
    }));

    await render(hbs`<Secrets::Page::Destinations::Destination::Sync @destination={{this.destination}} />`, {
      owner: this.engine,
    });
  });

  test('it should fetch and render kv mounts', async function (assert) {
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    assert
      .dom(searchSelect.selectedOption())
      .hasText('my-kv/', 'kv mounts are fetched and render in search select');
  });

  test('it should render secret suggestions for selected mount', async function (assert) {
    assert.dom(kvSuggestion.input).isDisabled('Secret input disabled when mount has not been selected');
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(kvSuggestion.input);
    assert.dom(searchSelect.option()).hasText('my-path/', 'Nested secret path renders');
    assert.dom(searchSelect.option(1)).hasText('my-secret', 'Secret renders');
    await click(searchSelect.removeSelected);
    assert.dom(kvSuggestion.input).hasValue('', 'secret path value is cleared when mount is unset');
  });

  test('it should render secret suggestions for nested paths', async function (assert) {
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(kvSuggestion.input);
    await click(searchSelect.option());
    assert
      .dom(searchSelect.option())
      .hasText('nested-secret', 'Suggestions render for secret at nested path');
  });

  test('it should sync secret', async function (assert) {
    assert.expect(9);

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-kv', secret_name: 'my-secret' };
      assert.deepEqual(data, expected, 'Sync request made with mount and secret name');
      return { data: { associated_secrets: { 'my-kv_12345': data } } };
    });

    assert.dom(submit).isDisabled('Submit button is disabled when mount is not selected');
    assert.dom(cancel).hasText('Back', 'back button renders');
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    assert.dom(submit).isDisabled('Submit button is disabled when secret is not selected');
    await click(kvSuggestion.input);
    await click(searchSelect.option(1));
    await click(submit);
    assert.dom(cancel).hasText('View synced secrets', 'view secrets tertiary renders');
    assert.dom(kvSuggestion.input).hasNoValue('Secret path is unset after submit success');
    assert.dom(submit).isDisabled('Submit button is disabled');
    assert
      .dom(successMessage)
      .includesText('Sync operation successfully initiated for my-secret.', 'Success banner renders');
    await click(searchSelect.removeSelected);
    assert.dom(successMessage).doesNotExist('clearing kv v2 mount path clears success banner');
  });

  test('it should allow manual mount path input if kv mounts are not returned', async function (assert) {
    assert.expect(1);

    this.server.get('/sys/internal/ui/mounts', () => ({
      data: { secret: { 'cubbyhole/': { type: 'cubbyhole' } } },
    }));

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-kv', secret_name: 'my-secret' };
      assert.deepEqual(data, expected, 'Sync request made with mount and secret name');
      return {};
    });

    await render(hbs`<Secrets::Page::Destinations::Destination::Sync @destination={{this.destination}} />`, {
      owner: this.engine,
    });

    await fillIn(mountInput, 'my-kv');
    await click(kvSuggestion.input);
    await click(searchSelect.option(1));
    await click(submit);
  });

  test('it should render alert banner on sync error', async function (assert) {
    assert.expect(1);

    const { type, name } = this.destination;
    const error = 'Secret not found. Please provide full path to existing secret';
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, () => {
      return new Response(400, {}, { errors: [error] });
    });

    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(kvSuggestion.input);
    await click(searchSelect.option(1));
    await click(submit);

    assert.dom(messageError).hasTextContaining(error, 'Error renders in alert banner');
  });

  test('it should clear sync associations from store in willDestroy hook', async function (assert) {
    const clearDatasetStub = sinon.stub(this.store, 'clearDataset');

    this.renderComponent = true;
    await render(
      hbs`
      {{#if this.renderComponent}}
        <Secrets::Page::Destinations::Destination::Sync @destination={{this.destination}} />
      {{/if}}
    `,
      { owner: this.engine }
    );
    this.set('renderComponent', false);
    await settled();

    assert.true(
      clearDatasetStub.calledWith('sync/association'),
      'Sync associations are cleared from store on component teardown'
    );
  });
});
