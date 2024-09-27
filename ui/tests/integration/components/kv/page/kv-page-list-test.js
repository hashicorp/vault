/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const CREATE_RECORDS = (number, store, server) => {
  const mirageList = server.createList('kv-metadatum', number, 'withCustomPath');
  mirageList.forEach((record) => {
    record.data.path = record.path;
    record.id = kvMetadataPath(record.data.backend, record.data.path);
    store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...record,
    });
  });
};

const META = {
  currentPage: 1,
  lastPage: 2,
  nextPage: 2,
  prevPage: 1,
  total: 16,
  filteredTotal: 16,
  pageSize: 15,
};

module('Integration | Component | kv | Page::List', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    setRunOptions({
      rules: {
        // TODO: ConfirmAction renders modal within list when @isInDropdown
        list: { enabled: false },
      },
    });
  });

  test('it renders Pagination and allows you to delete a kv/metadata record', async function (assert) {
    assert.expect(20);
    CREATE_RECORDS(15, this.store, this.server);
    this.model = await this.store.peekAll('kv/metadata');
    this.model.meta = META;
    this.backend = 'kv-engine';
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
    ];
    this.failedDirectoryQuery = false;
    await render(
      hbs`<Page::List
      @secrets={{this.model}}
      @backend={{this.backend}}
      @failedDirectoryQuery={{this.failedDirectoryQuery}}
      @breadcrumbs={{this.breadcrumbs}}
      @meta={{this.model.meta}}
      @currentRouteParams={{array this.backend}}
    />`,
      {
        owner: this.engine,
      }
    );

    assert.dom(PAGE.list.pagination).exists('shows hds pagination component');
    assert.dom(PAGE.list.paginationInfo).hasText('1–15 of 16', 'shows correct page of pages');
    assert.dom(PAGE.title).includesText(this.backend, 'shows backend as title');

    this.model.forEach((record) => {
      assert.dom(PAGE.list.item(record.path)).exists('lists all records from 0-14 on the first page');
    });

    this.server.delete(kvMetadataPath('kv-engine', 'my-secret-0'), () => {
      assert.ok(true, 'request made to correct endpoint on delete metadata.');
    });

    const popupSelector = `${PAGE.list.item('my-secret-0')} ${PAGE.popup}`;
    await click(popupSelector);
    await click('[data-test-popup-metadata-delete]');
    await click('[data-test-confirm-button]');
    assert.dom(PAGE.list.item('my-secret-0')).doesNotExist('deleted the first record from the list');
  });
});
