/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, waitFor, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import Sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { capabilitiesStub, overrideResponse } from 'vault/tests/helpers/stubs';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

// this test coverage mostly is around the export button functionality
// since everything else is static
module('Integration | Component | clients/page-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.downloadStub = Sinon.stub(this.owner.lookup('service:download'), 'download');
    this.startTimestamp = '2022-06-01T23:00:11.050Z';
    this.endTimestamp = '2022-12-01T23:00:11.050Z';
    this.selectedNamespace = undefined;
    this.upgradesDuringActivity = [];
    this.noData = undefined;
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub('sys/internal/counters/activity/export', ['sudo'])
    );

    this.renderComponent = async () => {
      return render(hbs`
        <Clients::PageHeader
          @startTimestamp={{this.startTimestamp}}
          @endTimestamp={{this.endTimestamp}}
          @namespace={{this.selectedNamespace}}
          @upgradesDuringActivity={{this.upgradesDuringActivity}}
          @noData={{this.noData}}
        />`);
    };
  });

  test('it shows the export button if user does has SUDO capabilities', async function (assert) {
    await this.renderComponent();
    assert.dom(CLIENT_COUNT.exportButton).exists();
  });

  test('it hides the export button if user does has SUDO capabilities but there is no data', async function (assert) {
    this.noData = true;
    await this.renderComponent();
    assert.dom(CLIENT_COUNT.exportButton).doesNotExist();
  });

  test('it hides the export button if user does not have SUDO capabilities', async function (assert) {
    this.server.post('/sys/capabilities-self', () =>
      capabilitiesStub('sys/internal/counters/activity/export', ['read'])
    );

    await this.renderComponent();
    assert.dom(CLIENT_COUNT.exportButton).doesNotExist();
  });

  test('defaults to show the export button if capabilities cannot be read', async function (assert) {
    this.server.post('/sys/capabilities-self', () => overrideResponse(403));

    await this.renderComponent();
    assert.dom(CLIENT_COUNT.exportButton).exists();
  });

  test('it shows the export API error on the modal', async function (assert) {
    this.server.get('/sys/internal/counters/activity/export', function () {
      return overrideResponse(403);
    });

    await this.renderComponent();

    await click(CLIENT_COUNT.exportButton);
    await click(GENERAL.confirmButton);
    await waitFor('[data-test-export-error]');
    assert.dom('[data-test-export-error]').hasText('permission denied');
  });

  test('it exports when json format', async function (assert) {
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

    await click(CLIENT_COUNT.exportButton);
    await fillIn('[data-test-download-format]', 'jsonl');
    await click(GENERAL.confirmButton);
    await waitUntil(() => this.downloadStub.calledOnce);
    const extension = this.downloadStub.lastCall.args[2];
    assert.strictEqual(extension, 'jsonl');
  });

  test('it exports when csv format', async function (assert) {
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

    await click(CLIENT_COUNT.exportButton);
    await fillIn('[data-test-download-format]', 'csv');
    await click(GENERAL.confirmButton);
    await waitUntil(() => this.downloadStub.calledOnce);
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

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.exportButton).exists();
    await click(CLIENT_COUNT.exportButton);
    await click(GENERAL.confirmButton);
  });
  test('it sends the selected namespace in export request', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/internal/counters/activity/export', function (_, req) {
      assert.strictEqual(req.requestHeaders['X-Vault-Namespace'], 'foobar');
      return new Response(200, { 'Content-Type': 'text/csv' }, '');
    });
    this.selectedNamespace = 'foobar/';

    await this.renderComponent();
    assert.dom(CLIENT_COUNT.exportButton).exists();
    await click(CLIENT_COUNT.exportButton);
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

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.exportButton).exists();
    await click(CLIENT_COUNT.exportButton);
    await click(GENERAL.confirmButton);
  });

  test('it shows a no data message if export returns 204', async function (assert) {
    this.server.get('/sys/internal/counters/activity/export', () => overrideResponse(204));
    await this.renderComponent();

    await click(CLIENT_COUNT.exportButton);
    await click(GENERAL.confirmButton);
    await waitFor('[data-test-export-error]');
    assert.dom('[data-test-export-error]').hasText('no data to export in provided time range.');
  });

  test('it shows upgrade data in export modal', async function (assert) {
    this.upgradesDuringActivity = [
      { version: '1.10.1', previousVersion: '1.9.9', timestampInstalled: '2021-11-18T10:23:16Z' },
    ];
    await this.renderComponent();
    await click(CLIENT_COUNT.exportButton);
    await waitFor('[data-test-export-upgrade-warning]');
    assert.dom('[data-test-export-upgrade-warning]').includesText('1.10.1 (Nov 18, 2021)');
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
      await click(CLIENT_COUNT.exportButton);
      await click(GENERAL.confirmButton);
      await waitUntil(() => this.downloadStub.calledOnce);
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

      await click(CLIENT_COUNT.exportButton);
      await click(GENERAL.confirmButton);
      await waitUntil(() => this.downloadStub.calledOnce);
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

      await click(CLIENT_COUNT.exportButton);
      await click(GENERAL.confirmButton);
      await waitUntil(() => this.downloadStub.calledOnce);
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

      await click(CLIENT_COUNT.exportButton);
      await click(GENERAL.confirmButton);
      await waitUntil(() => this.downloadStub.calledOnce);
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

      await click(CLIENT_COUNT.exportButton);
      await click(GENERAL.confirmButton);
      await waitUntil(() => this.downloadStub.calledOnce);
      const [filename] = this.downloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_export_foo');
    });
  });
});
