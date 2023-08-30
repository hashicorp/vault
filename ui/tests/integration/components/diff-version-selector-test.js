/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// TODO kv engine cleanup
import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const VERSIONS = [
  {
    version: 2,
  },
  {
    version: 1,
  },
];

module('Integration | Component | diff-version-selector', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.set(
      'model',
      EmberObject.create({
        currentVersion: 2,
        versions: VERSIONS,
      })
    );
    await render(hbs`<DiffVersionSelector @model={{this.model}} />`);
    const leftSideVersion = document
      .querySelector('[data-test-popup-menu-trigger="left-version"]')
      .innerText.trim();
    assert.strictEqual(leftSideVersion, 'Version 2', 'left side toolbar defaults to currentVersion');

    await click('[data-test-popup-menu-trigger="left-version"]');

    assert.dom('[data-test-leftSide-version="1"]').exists('leftside shows both versions');
    assert.dom('[data-test-leftSide-version="2"]').exists('leftside shows both versions');
  });
});
