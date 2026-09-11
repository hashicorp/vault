/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const SELECTORS = {
  container: '[data-test-about-this-engine]',
  heading: '[data-test-about-engine-heading]',
  content: '[data-test-about-engine-content]',
  toggle: '[data-test-button="About engine toggle"]',
  row: (type) => `[data-test-about-engine-row="${type}"]`,
};

module('Integration | Component | mount/about-this-engine', function (hooks) {
  setupRenderingTest(hooks);

  module('common engines', function () {
    test('it renders the about-this-engine section for kv', async function (assert) {
      // The KV engine is one of the five common engines and should render the card
      // with all three rows open by default.
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      assert.dom(SELECTORS.container).exists('renders the about-this-engine wrapper for kv');
      assert.dom(SELECTORS.toggle).exists('renders the toggle button');
      assert.dom(SELECTORS.content).exists('content is visible (open by default)');
      assert.dom(SELECTORS.row('preferredFor')).exists('renders the "Preferred for" row');
      assert.dom(SELECTORS.row('keyFeatures')).exists('renders the "Key features" row');
      assert.dom(SELECTORS.row('notSuitedFor')).exists('renders the "Not suited for" row');
    });

    test('it renders correct label text for each row', async function (assert) {
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      assert
        .dom(SELECTORS.row('preferredFor'))
        .containsText('Preferred for', 'preferredFor row has correct label');
      assert
        .dom(SELECTORS.row('keyFeatures'))
        .containsText('Key features', 'keyFeatures row has correct label');
      assert
        .dom(SELECTORS.row('notSuitedFor'))
        .containsText('Not suited for', 'notSuitedFor row has correct label');
    });

    test('it renders engine-specific content for kv', async function (assert) {
      // Verifies that KV-specific bullet copy is rendered in the expected rows
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      assert
        .dom(SELECTORS.row('preferredFor'))
        .containsText('Bearer tokens', 'kv preferredFor row contains kv-specific text');
      assert
        .dom(SELECTORS.row('keyFeatures'))
        .containsText('key-value format', 'kv keyFeatures row contains kv-specific text');
      assert
        .dom(SELECTORS.row('notSuitedFor'))
        .containsText('short-lived secrets', 'kv notSuitedFor row contains kv-specific text');
    });

    test('it renders engine-specific content for aws', async function (assert) {
      // Verifies that AWS-specific bullet copy is rendered when the engine type is aws
      await render(hbs`<Mount::AboutThisEngine @engineType="aws" />`);

      assert.dom(SELECTORS.container).exists('renders for aws engine');
      assert
        .dom(SELECTORS.row('preferredFor'))
        .containsText('scoped AWS credentials', 'aws preferredFor row contains aws-specific text');
    });

    test('it renders the section open by default', async function (assert) {
      // The "About this engine" section should be expanded on initial render
      // so users see the information without having to click anything.
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      assert.dom(SELECTORS.content).exists('content is visible on initial render (open by default)');
    });

    test('it can be toggled closed and re-opened', async function (assert) {
      // Users should be able to collapse and expand the section freely
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      // Initially open
      assert.dom(SELECTORS.content).exists('content is visible initially');

      // Click the toggle button to close
      await click(SELECTORS.toggle);

      assert.dom(SELECTORS.content).doesNotExist('content is hidden after closing');

      // Click again to re-open
      await click(SELECTORS.toggle);

      assert.dom(SELECTORS.content).exists('content is visible again after re-opening');
    });

    test('it renders inline links in bullet points', async function (assert) {
      // KV key features include inline links (e.g., Custom metadata, Version management)
      await render(hbs`<Mount::AboutThisEngine @engineType="kv" />`);

      assert
        .dom(`${SELECTORS.row('keyFeatures')} a`)
        .exists('at least one inline link renders inside the key features row');
    });
  });

  module('non-common engines', function () {
    test('it does not render for engines without about info', async function (assert) {
      // Non-common engines (e.g., transit) have no entry in the about-engine-info helper
      // and should render nothing at all.
      await render(hbs`<Mount::AboutThisEngine @engineType="transit" />`);

      assert.dom(SELECTORS.container).doesNotExist('does not render for transit (non-common engine)');
    });

    test('it does not render when engine type is empty', async function (assert) {
      await render(hbs`<Mount::AboutThisEngine @engineType="" />`);

      assert.dom(SELECTORS.container).doesNotExist('does not render when engineType is empty');
    });
  });
});
