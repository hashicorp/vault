/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.mountData = {
      id: 'my-kv',
      accessor: 'kv_80616825',
      config: this.store.createRecord('mount-config', {
        defaultLeaseTtl: '72h',
        forceNoCache: false,
        maxLeaseTtl: '123h',
      }),
      options: {
        version: '2',
      },
      description: '',
      path: 'my-kv',
      sealWrap: false,
      type: 'kv',
      uuid: 'f1739f9d-dfc0-83c8-011f-ec17103a06a1',
      // TODO: remove when attrs aren't duplicated across models
      // these kv specific attrs exist on the secret-engine model (for POST request when mounting the engine)
      // we want to make sure we're rendering values from kv/config while duplicates exist
      maxVersions: 'this should never render',
      casRequired: 'test is failing if this shows',
      deleteVersionAfter: `definitely shouldn't render this`,
    };
    this.store.pushPayload('kv/config', {
      modelName: 'kv/config',
      id: 'my-config',
      data: { max_versions: 0, cas_required: false, delete_version_after: '0s' },
    });

    // this is the route model, not an ember data model
    this.model = {
      engineConfig: this.store.peekRecord('kv/config', 'my-config'),
      mountConfig: this.store.createRecord('secret-engine', this.mountData),
    };

    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.model.mountConfig.path, route: 'list' },
      { label: 'Configuration' },
    ];
  });

  test('it renders kv configuration details', async function (assert) {
    assert.expect(11);

    await render(
      hbs`
      <Page::Configuration
        @engineConfig={{this.model.engineConfig}}
        @mountConfig={{this.model.mountConfig}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).includesText(this.mountData.path, 'renders engine path as page title');
    assert.dom(PAGE.infoRowValue('Require check and set')).hasText('No');
    assert.dom(PAGE.infoRowValue('Automate secret deletion')).hasText('Never delete');
    assert.dom(PAGE.infoRowValue('Maximum number of versions')).hasText('0');
    assert.dom(PAGE.infoRowValue('Accessor')).hasText(this.mountData.accessor);
    assert.dom(PAGE.infoRowValue('Path')).hasText(this.mountData.path);
    assert.dom(PAGE.infoRowValue('Type')).hasText(this.mountData.type);
    assert.dom(PAGE.infoRowValue('Description')).doesNotExist();
    assert.dom(PAGE.infoRowValue('Seal wrap')).hasText('No');
    assert.dom(PAGE.infoRowValue('Default Lease TTL')).hasText('3 days');
    assert.dom(PAGE.infoRowValue('Max Lease TTL')).hasText('5 days 3 hours');
  });

  test('it renders non default kv engine config data', async function (assert) {
    assert.expect(3);
    this.model.engineConfig.maxVersions = 10;
    this.model.engineConfig.casRequired = true;
    this.model.engineConfig.deleteVersionAfter = '10d';

    await render(
      hbs`
      <Page::Configuration
        @engineConfig={{this.model.engineConfig}}
        @mountConfig={{this.model.mountConfig}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(PAGE.infoRowValue('Require check and set')).hasText('Yes');
    assert.dom(PAGE.infoRowValue('Automate secret deletion')).hasText('10 days');
    assert.dom(PAGE.infoRowValue('Maximum number of versions')).hasText('10');
  });
});
