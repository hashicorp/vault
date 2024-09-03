/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { dateFormat } from 'core/helpers/date-format';
import { dateFromNow } from 'core/helpers/date-from-now';
import { baseSetup } from 'vault/tests/helpers/kv/kv-run-commands';

const { overviewCard } = GENERAL;

// subkeys access is enterprise only (in the GUI) but we don't have any version testing here because the @subkeys arg is null for non-enterprise versions
module('Integration | Component | kv-v2 | Page::Secret::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    baseSetup(this);
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];
    this.subkeys = {
      subkeys: {
        foo: null,
        bar: {
          baz: null,
        },
        quux: null,
      },
      metadata: {
        created_time: '2021-12-14T20:28:00.773477Z',
        custom_metadata: null,
        deletion_time: '',
        destroyed: false,
        version: 1,
      },
    };
    this.canReadMetadata = true;
    this.canUpdateData = true;

    this.format = (time) => dateFormat([time, 'MMM d yyyy, h:mm:ss aa'], {});
    this.renderComponent = async () => {
      return render(
        hbs`
        <Page::Secret::Overview
          @backend={{this.backend}}
          @breadcrumbs={{this.breadcrumbs}}
          @canReadMetadata={{this.canReadMetadata}}
          @canUpdateData={{this.canUpdateData}}
          @metadata={{this.metadata}}
          @path={{this.path}}
          @subkeys={{this.subkeys}}
        />`,
        { owner: this.engine }
      );
    };
  });

  module('active secret (version not deleted or destroyed)', function () {
    test('it renders tabs', async function (assert) {
      await this.renderComponent();
      const tabs = ['Overview', 'Secret', 'Metadata', 'Paths', 'Version History'];
      for (const tab of tabs) {
        assert.dom(PAGE.secretTab(tab)).hasText(tab);
      }
    });

    test('it renders header', async function (assert) {
      await this.renderComponent();
      assert.dom(PAGE.breadcrumbs).hasText(`Secrets ${this.backend} ${this.path}`);
      assert.dom(PAGE.title).hasText(this.path);
    });

    test('it renders with full permissions', async function (assert) {
      await this.renderComponent();
      const fromNow = dateFromNow([this.metadata.updatedTime]); // uses date-fns so can't stub timestamp util
      assert.dom(`${overviewCard.container('Current version')} .hds-badge`).doesNotExist();
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Create new The current version of this secret. ${this.metadata.currentVersion}`
        );
      assert
        .dom(overviewCard.container('Secret age'))
        .hasText(
          `Secret age View metadata Current secret version age. Last updated on ${this.format(
            this.metadata.updatedTime
          )}. ${fromNow}`
        );
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
      assert
        .dom(overviewCard.container('Subkeys'))
        .hasText(
          `Subkeys JSON The table is displaying the top level subkeys. Toggle on the JSON view to see the full depth. Keys ${Object.keys(
            this.subkeys.subkeys
          ).join(' ')}`
        );
    });

    test('it hides link when no secret update permissions', async function (assert) {
      // creating a new version of a secret is updating a secret
      // the overview only exists after an initial version is created
      // which is why we just check for update and not also create
      this.canUpdateData = false;
      await this.renderComponent();
      assert
        .dom(`${overviewCard.container('Current version')} a`)
        .doesNotExist('create link does not render');
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(`Current version The current version of this secret. ${this.metadata.currentVersion}`);
    });

    test('it renders with no metadata permissions', async function (assert) {
      this.metadata = null;
      this.canReadMetadata = false;
      // all secret metadata instead comes from subkeys endpoint
      const subkeyMeta = this.subkeys.metadata;
      await this.renderComponent();
      const fromNow = dateFromNow([subkeyMeta.created_time]); // uses date-fns so can't stub timestamp util
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(`Current version Create new The current version of this secret. ${subkeyMeta.version}`);
      assert
        .dom(overviewCard.container('Secret age'))
        .hasText(
          `Secret age Current secret version age. Last updated on ${this.format(
            subkeyMeta.created_time
          )}. ${fromNow}`
        );
      assert.dom(`${overviewCard.container('Secret age')} a`).doesNotExist('metadata link does not render');
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
      assert
        .dom(overviewCard.container('Subkeys'))
        .hasText(
          `Subkeys JSON The table is displaying the top level subkeys. Toggle on the JSON view to see the full depth. Keys ${Object.keys(
            this.subkeys.subkeys
          ).join(' ')}`
        );
    });

    test('it renders with no subkeys permissions', async function (assert) {
      this.subkeys = null;
      await this.renderComponent();
      const fromNow = dateFromNow([this.metadata.updatedTime]); // uses date-fns so can't stub timestamp util
      const expectedTime = this.format(this.metadata.updatedTime);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Create new The current version of this secret. ${this.metadata.currentVersion}`
        );
      assert
        .dom(overviewCard.container('Secret age'))
        .hasText(
          `Secret age View metadata Current secret version age. Last updated on ${expectedTime}. ${fromNow}`
        );
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
    });

    test('it renders with no subkey or metadata permissions', async function (assert) {
      this.subkeys = null;
      this.metadata = null;
      await this.renderComponent();
      assert.dom(overviewCard.container('Current version')).doesNotExist();
      assert.dom(overviewCard.container('Secret age')).doesNotExist();
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
    });
  });

  module('deleted version', function (hooks) {
    hooks.beforeEach(async function () {
      // subkeys is null but metadata still has data
      this.subkeys = {
        subkeys: null,
        metadata: {
          created_time: '2021-12-14T20:28:00.773477Z',
          custom_metadata: null,
          deletion_time: '2022-02-14T20:28:00.773477Z',
          destroyed: false,
          version: 1,
        },
      };
      this.metadata.versions[4].deletion_time = '2024-08-15T23:01:08.312332Z';
      this.assertBadge = (assert) => {
        assert
          .dom(`${overviewCard.container('Current version')} .hds-badge`)
          .hasClass('hds-badge--color-neutral');
        assert
          .dom(`${overviewCard.container('Current version')} .hds-badge`)
          .hasClass('hds-badge--type-inverted');
        assert.dom(`${overviewCard.container('Current version')} .hds-badge`).hasText('Deleted');
      };
    });

    test('with full permissions', async function (assert) {
      const expectedTime = this.format(this.metadata.versions[4].deletion_time);
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Deleted Create new The current version of this secret was deleted ${expectedTime}. ${this.metadata.currentVersion}`
        );
      assert.dom(overviewCard.container('Secret age')).doesNotExist();
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
    });

    test('with no metadata permissions', async function (assert) {
      this.metadata = null;
      const expectedTime = this.format(this.subkeys.metadata.deletion_time);
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Deleted Create new The current version of this secret was deleted ${expectedTime}. ${this.subkeys.metadata.version}`
        );
    });

    test('with no subkey permissions', async function (assert) {
      this.subkeys = null;
      const expectedTime = this.format(this.metadata.versions[4].deletion_time);
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Deleted Create new The current version of this secret was deleted ${expectedTime}. ${this.metadata.currentVersion}`
        );
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
    });

    test('with no permissions', async function (assert) {
      this.subkeys = null;
      this.metadata = null;
      await this.renderComponent();
      assert.dom(overviewCard.container('Current version')).doesNotExist();
    });
  });

  module('destroyed version', function (hooks) {
    hooks.beforeEach(async function () {
      // subkeys is null but metadata still has data
      this.subkeys = {
        subkeys: null,
        metadata: {
          created_time: '2024-08-15T01:24:43.658478Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: true,
          version: 1,
        },
      };
      this.metadata.versions[4].destroyed = true;
      this.assertBadge = (assert) => {
        assert
          .dom(`${overviewCard.container('Current version')} .hds-badge`)
          .hasClass('hds-badge--color-critical');
        assert
          .dom(`${overviewCard.container('Current version')} .hds-badge`)
          .hasClass('hds-badge--type-outlined');
        assert.dom(`${overviewCard.container('Current version')} .hds-badge`).hasText('Destroyed');
      };
    });

    test('with full permissions', async function (assert) {
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Destroyed Create new The current version of this secret has been permanently deleted and cannot be restored. ${this.metadata.currentVersion}`
        );
      assert.dom(overviewCard.container('Secret age')).doesNotExist();
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
      assert
        .dom(overviewCard.container('Paths'))
        .hasText(
          `Paths The paths to use when referring to this secret in API or CLI. API path /v1/${this.backend}/data/${this.path} CLI path -mount="${this.backend}" "${this.path}"`
        );
    });

    test('with no metadata permissions', async function (assert) {
      this.metadata = null;
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Destroyed Create new The current version of this secret has been permanently deleted and cannot be restored. ${this.subkeys.metadata.version}`
        );
    });

    test('with no subkeys permissions', async function (assert) {
      this.subkeys = null;
      await this.renderComponent();
      this.assertBadge(assert);
      assert
        .dom(overviewCard.container('Current version'))
        .hasText(
          `Current version Destroyed Create new The current version of this secret has been permanently deleted and cannot be restored. ${this.metadata.currentVersion}`
        );
      assert.dom(overviewCard.container('Subkeys')).doesNotExist();
    });

    test('with no permissions', async function (assert) {
      this.subkeys = null;
      this.metadata = null;
      await this.renderComponent();
      assert.dom(overviewCard.container('Current version')).doesNotExist();
    });
  });
});
