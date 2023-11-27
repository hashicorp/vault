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
import { render, click, fillIn } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { selectChoose } from 'ember-power-select/test-support/helpers';
import sinon from 'sinon';
import { Response } from 'miragejs';

const { destinations, searchSelect, messageError } = PAGE;
const { mountSelect, mountInput, secretInput, submit, cancel } = destinations.sync;

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
    assert.dom(secretInput).isDisabled('Secret input disabled when mount has not been selected');
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(secretInput);
    assert.dom(searchSelect.option()).hasText('my-path/', 'Nested secret path renders');
    assert.dom(searchSelect.option(1)).hasText('my-secret', 'Secret renders');
  });

  test('it should render secret suggestions for nested paths', async function (assert) {
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(secretInput);
    await click(searchSelect.option());
    assert
      .dom(searchSelect.option())
      .hasText('nested-secret', 'Suggestions render for secret at nested path');
    await click(searchSelect.option());
    assert
      .dom(searchSelect.noMatch)
      .hasText('No suggestions for this path', 'No match message renders when secret is selected');
  });

  test('it should sync secret', async function (assert) {
    assert.expect(3);

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-kv', secret_name: 'my-secret' };
      assert.deepEqual(data, expected, 'Sync request made with mount and secret name');
      return {};
    });

    assert.dom(submit).isDisabled('Submit button is disabled when mount is not selected');
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    assert.dom(submit).isDisabled('Submit button is disabled when secret is not selected');
    await click(secretInput);
    await click(searchSelect.option(1));
    await click(submit);
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
    await click(secretInput);
    await click(searchSelect.option(1));
    await click(submit);
  });

  test('it should transition to destination secrets route on cancel', async function (assert) {
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    await click(cancel);
    assert.propEqual(
      transitionStub.lastCall.args,
      ['vault.cluster.sync.secrets.destinations.destination.secrets'],
      'Transitions to destination secrets route on cancel'
    );
  });

  test('it should render alert banner on sync error', async function (assert) {
    assert.expect(1);

    const { type, name } = this.destination;
    const error = 'Secret not found. Please provide full path to existing secret';
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, () => {
      return new Response(400, {}, { errors: [error] });
    });

    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(secretInput);
    await click(searchSelect.option(1));
    await click(submit);

    assert.dom(messageError).hasTextContaining(error, 'Error renders in alert banner');
  });
});
