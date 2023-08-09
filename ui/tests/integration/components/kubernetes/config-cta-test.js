/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | kubernetes | ConfigCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  test('it should render message and action', async function (assert) {
    await render(hbs`<ConfigCta />`, { owner: this.engine });
    assert.dom('[data-test-empty-state-title]').hasText('Kubernetes not configured', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Get started by establishing the URL of the Kubernetes API to connect to, along with some additional options.',
        'Message renders'
      );
    assert.dom('[data-test-config-cta] a').hasText('Configure Kubernetes', 'Action renders');
  });
});
