/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';

module('Integration | Component | SecretEngine::PageHeader', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = keyMgmtMockModel;
  });

  test('it shows page header title, description, and general settings tab', async function (assert) {
    assert.expect(4);
    await render(hbs`
      <SecretEngine::PageHeader @model={{this.model}}/>
    `);
    assert.dom(GENERAL.tab('general-settings')).exists('contains general settings tab');
    assert.dom(GENERAL.tab('plugin-settings')).doesNotExist('does not contain plugin settings tab');
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(`${this.model.secretsEngine.id} configuration`, 'displays page header title');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .hasText(
        engineDisplayData(this.model.secretsEngine.type).displayName,
        'displays page header description'
      );
  });

  test('it shows page header title, description, and general and plugin settings tab for configurable secret engines', async function (assert) {
    assert.expect(4);
    this.model.secretsEngine = {
      type: 'aws',
      id: 'aws',
      config: {
        region: 'us-west-2',
        access_key: '123-key',
        iam_endpoint: 'iam-endpoint',
        sts_endpoint: 'sts-endpoint',
        max_retries: 1,
      },
    };

    await render(hbs`
      <SecretEngine::PageHeader @model={{this.model}}/>
    `);
    assert.dom(GENERAL.tab('general-settings')).exists('contains general settings tab');
    assert.dom(GENERAL.tab('plugin-settings')).exists('contains plugin settings tab');
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText(`${this.model.secretsEngine.id} configuration`, 'displays page header title');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .hasText(
        engineDisplayData(this.model.secretsEngine.type).displayName,
        'displays page header description'
      );
  });
});
