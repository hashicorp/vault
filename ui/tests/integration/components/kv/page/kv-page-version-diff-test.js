/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::VersionDiff', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.metadata = {
      current_version: 4,
      updated_time: '2023-07-21T03:11:58.095971Z',
      versions: {
        1: {
          created_time: '2018-03-22T02:24:06.945319214Z',
          deletion_time: '',
          destroyed: false,
        },
        2: {
          created_time: '2023-07-20T02:15:35.86465Z',
          deletion_time: '2023-07-25T00:36:19.950545Z',
          destroyed: false,
        },
        3: {
          created_time: '2023-07-20T02:15:40.164549Z',
          deletion_time: '',
          destroyed: true,
        },
        4: {
          created_time: '2023-07-21T03:11:58.095971Z',
          deletion_time: '',
          destroyed: false,
        },
      },
    };
    this.breadcrumbs = [{ label: 'Version History', route: 'secret.metadata.versions' }, { label: 'Diff' }];

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Metadata::VersionDiff
            @metadata={{this.metadata}}
            @path={{this.path}}
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
          />
        `,
        { owner: this.engine }
      );

    this.api = this.owner.lookup('service:api');
    this.queryParamsStub = sinon.stub(this.api, 'addQueryParams');
    this.fetchStub = sinon.stub(this.api.secrets, 'kvV2Read').callsFake((path, backend, initOverride) => {
      initOverride();
      return Promise.resolve({});
    });
  });

  test('it renders empty states when current version is deleted or destroyed', async function (assert) {
    assert.expect(7);

    const { current_version } = this.metadata;
    // destroyed
    this.metadata.versions[current_version].destroyed = true;
    await this.renderComponent();

    assert.true(this.fetchStub.calledWith(this.path, this.backend));
    assert.deepEqual(
      this.queryParamsStub.firstCall.args[1],
      { version: 1 },
      'correct version passed as query param to first request'
    );
    assert.deepEqual(
      this.queryParamsStub.secondCall.args[1],
      { version: 4 },
      'correct version passed as query param to second request'
    );
    assert.dom(PAGE.emptyStateTitle).hasText(`Version ${current_version} has been destroyed`);
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText('The current version of this secret has been destroyed. Select another version to compare.');

    // deleted
    this.metadata.versions[current_version].destroyed = false;
    this.metadata.versions[current_version].deletion_time = '2023-07-25T00:36:19.950545Z';
    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.path}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.emptyStateTitle).hasText(`Version ${current_version} has been deleted`);
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText('The current version of this secret has been deleted. Select another version to compare.');
  });

  test('it renders compared data of the two versions and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(12);

    this.fetchStub.onFirstCall().resolves({ data: { hello: 'world' } }); // version 1
    this.fetchStub.onSecondCall().resolves({ data: { foo: 'bar' } }); // version 4 (current version)

    await this.renderComponent();

    const [left, right] = findAll(PAGE.detail.versionDropdown);
    assert.dom(PAGE.diff.visualDiff).hasText(
      `foo\"bar\"hello\"world\"`, // eslint-disable-line no-useless-escape
      'correctly pull in the data from version 4 and compared to version 1.'
    );
    assert.dom(PAGE.diff.deleted).hasText(`hello"world"`);
    assert.dom(PAGE.diff.added).hasText(`foo"bar"`);
    assert.dom(right).hasText('Version 4', 'shows the current version for the left side default version.');
    assert.dom(left).hasText('Version 1', 'shows the latest active version on init.');

    await click(left);

    for (const num in this.metadata.versions) {
      const data = this.metadata.versions[num];
      assert.dom(PAGE.detail.version(num)).exists('renders the button for each version.');

      if (data.destroyed || data.deletion_time) {
        assert
          .dom(`${PAGE.detail.version(num)} [data-test-icon="x-square-fill"]`)
          .hasClass(`${data.destroyed ? 'has-text-danger' : 'has-text-grey'}`);
      }
    }
    assert
      .dom(`${PAGE.detail.version('1')} button`)
      .hasClass('is-active', 'correctly shows the selected version 1 as active.');
  });
});
