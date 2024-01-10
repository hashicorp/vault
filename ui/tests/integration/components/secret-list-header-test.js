/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { setupMirage } from 'ember-cli-mirage/test-support';
import mirageToModels from 'vault/tests/helpers/mirage-to-models';

module('Integration | Component | secret-list-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it should render version badge for kv and generic engine types', async function (assert) {
    const backends = supportedSecretBackends();
    assert.expect(backends.length);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
        />
      `);
      const selector = '[data-test-kv-version-badge]';

      if (['kv', 'generic'].includes(type)) {
        assert
          .dom(selector)
          .hasText(
            `version ${this.model.version || 1}`,
            `Badge renders with correct version for ${type} engine type`
          );
      } else {
        assert.dom(selector).doesNotExist(`Version badge does not render for ${type} engine type`);
      }
    }
  });
});
