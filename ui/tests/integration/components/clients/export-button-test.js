/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import Sinon from 'sinon';
import { Response } from 'miragejs';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Integration | Component | clients/export-button', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.downloadStub = Sinon.stub(this.owner.lookup('service:download'), 'download');
    this.startTimestamp = '2022-06-01T23:00:11.050Z';
    this.endTimestamp = '2022-12-01T23:00:11.050Z';
    this.selectedNamespace = undefined;

    this.renderComponent = async () => {
      return render(hbs`
        <Clients::ExportButton
          @startTimestamp={{this.startTimestamp}}
          @endTimestamp={{this.endTimestamp}}
          @selectedNamespace={{this.selectedNamespace}}
        />`);
    };
  });

  test('it renders modal with yielded alert', async function (assert) {
    await render(hbs`
      <Clients::ExportButton
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
      >
        <:alert>
          <Hds::Alert class="has-top-padding-m" @type="compact" @color="warning" as |A|>
            <A.Description data-test-custom-alert>Yielded alert!</A.Description>
          </Hds::Alert>
        </:alert>
      </Clients::ExportButton>
    `);

    await click('[data-test-attribution-export-button]');
    assert.dom('[data-test-custom-alert]').hasText('Yielded alert!');
  });

  test('shows the API error on the modal', async function (assert) {
    this.server.get('/sys/internal/counters/activity/export', function () {
      return new Response(
        403,
        { 'Content-Type': 'application/json' },
        { errors: ['this is an error from the API'] }
      );
    });

    await this.renderComponent();

    await click('[data-test-attribution-export-button]');
    await click(GENERAL.confirmButton);
    assert.dom('[data-test-export-error]').hasText('this is an error from the API');
  });

  test('it works for json format', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.deepEqual(req.queryParams, {
        format: 'json',
        start_time: '2022-06-01T23:00:11.050Z',
        end_time: '2022-12-01T23:00:11.050Z',
      });
      return new Response(200, { 'Content-Type': 'application/json' }, { example: 'data' });
    });

    await this.renderComponent();

    await click('[data-test-attribution-export-button]');
    await fillIn('[data-test-download-format]', 'jsonl');
    await click(GENERAL.confirmButton);
    const extension = this.downloadStub.lastCall.args[2];
    assert.strictEqual(extension, 'jsonl');
  });

  test('it works for csv format', async function (assert) {
    assert.expect(2);

    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.deepEqual(req.queryParams, {
        format: 'csv',
        start_time: '2022-06-01T23:00:11.050Z',
        end_time: '2022-12-01T23:00:11.050Z',
      });
      return new Response(200, { 'Content-Type': 'text/csv' }, 'example,data');
    });

    await this.renderComponent();

    await click('[data-test-attribution-export-button]');
    await fillIn('[data-test-download-format]', 'csv');
    await click(GENERAL.confirmButton);
    const extension = this.downloadStub.lastCall.args[2];
    assert.strictEqual(extension, 'csv');
  });

  test('it sends the current namespace in export request', async function (assert) {
    assert.expect(2);
    const namespaceSvc = this.owner.lookup('service:namespace');
    namespaceSvc.path = 'foo';
    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foo');
      return new Response(200, { 'Content-Type': 'text/csv' }, '');
    });

    await render(hbs`
      <Clients::ExportButton
        @totalClientAttribution={{this.totalClientAttribution}}
        @responseTimestamp={{this.timestamp}}
        />
    `);
    assert.dom('[data-test-attribution-export-button]').exists();
    await click('[data-test-attribution-export-button]');
    await click(GENERAL.confirmButton);
  });
  test('it sends the selected namespace in export request', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foobar');
      return new Response(200, { 'Content-Type': 'text/csv' }, '');
    });
    this.selectedNamespace = 'foobar/';

    await render(hbs`
      <Clients::ExportButton
        @totalClientAttribution={{this.totalClientAttribution}}
        @responseTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        />
    `);
    assert.dom('[data-test-attribution-export-button]').exists();
    await click('[data-test-attribution-export-button]');
    await click(GENERAL.confirmButton);
  });

  test('it sends the current + selected namespace in export request', async function (assert) {
    assert.expect(2);
    const namespaceSvc = this.owner.lookup('service:namespace');
    namespaceSvc.path = 'foo';
    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foo/bar');
      return new Response(200, { 'Content-Type': 'text/csv' }, '');
    });
    this.selectedNamespace = 'bar/';

    await render(hbs`
      <Clients::ExportButton
        @totalClientAttribution={{this.totalClientAttribution}}
        @responseTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        />
    `);
    assert.dom('[data-test-attribution-export-button]').exists();
    await click('[data-test-attribution-export-button]');
    await click(GENERAL.confirmButton);
  });

  test('it shows a no data message if endpoint returns 204', async function (assert) {
    this.server.get('/sys/internal/counters/activity/export', () => overrideResponse(204));

    await render(hbs`
      <Clients::ExportButton
        @totalClientAttribution={{this.totalClientAttribution}}
        @responseTimestamp={{this.timestamp}}
        />
    `);
    await click('[data-test-attribution-export-button]');
    await click(GENERAL.confirmButton);
    assert.dom('[data-test-export-error]').hasText('no data to export in provided time range.');
  });

  module('download naming', function () {
    test('is correct for date range', async function (assert) {
      assert.expect(2);
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
          start_time: '2022-06-01T23:00:11.050Z',
          end_time: '2022-12-01T23:00:11.050Z',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await this.renderComponent();
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const args = this.downloadStub.lastCall.args;
      const [filename] = args;
      assert.strictEqual(filename, 'clients_export_June 2022-December 2022', 'csv has expected filename');
    });

    test('is correct for a single month', async function (assert) {
      assert.expect(2);
      this.endTimestamp = '2022-06-21T23:00:11.050Z';
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
          start_time: '2022-06-01T23:00:11.050Z',
          end_time: '2022-06-21T23:00:11.050Z',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });
      await this.renderComponent();

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.downloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_export_June 2022', 'csv has single month in filename');
    });
    test('omits date if no start/end timestamp', async function (assert) {
      assert.expect(2);
      this.startTimestamp = undefined;
      this.endTimestamp = undefined;

      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await this.renderComponent();

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.downloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_export');
    });

    test('includes current namespace', async function (assert) {
      assert.expect(2);
      this.startTimestamp = undefined;
      this.endTimestamp = undefined;
      const namespace = this.owner.lookup('service:namespace');
      namespace.path = 'bar/';

      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await this.renderComponent();

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.downloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_export_bar');
    });

    test('includes selectedNamespace', async function (assert) {
      assert.expect(2);
      this.startTimestamp = undefined;
      this.endTimestamp = undefined;
      this.selectedNamespace = 'foo/';

      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await this.renderComponent();

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.downloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_export_foo');
    });
  });
});
