/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';

module('Integration | Component | kv | Page::Secret::MetadataEdit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const router = this.owner.lookup('service:router');
    const routerStub = sinon.stub(router, 'transitionTo');
    this.onCancel = sinon.spy();
    this.transitionCalledWith = (routeName, path) => {
      const args = ['vault.cluster.secrets.backend.kv.secret.metadata', path];
      return routerStub.calledWith(...args);
    };

    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const data = this.server.create('kv-metadatum');
    data.id = kvMetadataPath('kv-engine', 'my-secret');
    store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...data,
    });
    this.model = store.peekRecord('kv/metadata', data.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path, route: 'secret', model: this.model.path },
      { label: 'metadata' },
    ];
  });

  test('meep it renders empty state when no custom_metadata is present', async function (assert) {
    assert.expect(2);
    await render(
      hbs`
      <Page::Secret::MetadataDetails
        @model={{this.model}}
        @breadcrumbs={{this.breadcrumbs}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}} />`,
      {
        owner: this.engine,
      }
    );
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No custom metadata', 'renders the correct empty state');
    assert.dom('[data-test-value-div="Maximum versions"]').hasText('0', 'renders maximum versions');
  });
});
