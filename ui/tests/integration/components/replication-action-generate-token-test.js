/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | replication-action-generate-token', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', () => {});
  });

  test('it renders with the expected elements', async function (assert) {
    await render(hbs`
      
      {{replication-action-generate-token}}
    `);
    assert.dom('h4.title').hasText('Generate operation token', 'renders default title');
    assert.dom('[data-test-replication-action-trigger]').hasText('Generate token', 'renders default CTA');
  });
});
