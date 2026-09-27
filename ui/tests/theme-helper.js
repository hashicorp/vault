/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import * as QUnit from 'qunit';
import {
  setupGlobalA11yHooks,
  a11yAudit,
  setCustomReporter,
  _DEFAULT_REPORTER,
} from 'ember-a11y-testing/test-support';

export function setupDarkThemeAudit() {
  // axe-core's color-contrast rule flags comment nodes (<!-- -->) and empty
  // text nodes that are siblings of real elements as contrast violations. These
  // are false positives — there is no visible text to check. Filter them out
  // so genuine violations are still caught.
  function isGhostNode(html) {
    // A node whose serialized HTML is purely comment markers, empty labels,
    // whitespace, or a bare closing tag is a ghost node.
    // Handles patterns seen in practice:
    //   <!---->          — single Glimmer comment
    //   <!----><!---->   — double Glimmer comment (empty conditional)
    //   <!----></label>  — comment followed by bare closing tag
    //   </label>         — bare closing tag with no text
    //
    // Use two separate passes to handle all combinations:
    // 1. Strip HTML comments (including empty ones like <!---->)
    // 2. Strip bare closing tags (nothing between < and >)
    const step1 = html.replace(/<!--[\s\S]*?-->/g, ''); // strip comments (no `s` flag needed)
    const step2 = step1.replace(/<\/[a-z][a-z0-9]*>/gi, ''); // strip bare closing tags
    const step3 = step2.replace(/<[a-z][a-z0-9]*\s*\/>/gi, ''); // strip self-closing void tags
    return step3.trim() === '';
  }

  function filterGhostNodes(results) {
    results.violations = results.violations
      .map((violation) => ({
        ...violation,
        nodes: violation.nodes.filter((node) => !isGhostNode(node.html)),
      }))
      .filter((violation) => violation.nodes.length > 0);
    return results;
  }

  setCustomReporter((results) => _DEFAULT_REPORTER(filterGhostNodes(results)));

  async function colorContrastAudit() {
    await a11yAudit({
      runOnly: {
        type: 'rule',
        values: ['color-contrast'],
      },
    });
  }

  setupGlobalA11yHooks(() => true, colorContrastAudit, {
    helpers: [
      'visit',
      'render',
      'click',
      'doubleClick',
      'tap',
      'fillIn',
      'focus',
      'blur',
      'tab',
      'select',
      'triggerEvent',
      'triggerKeyEvent',
      'typeIn',
    ],
  });

  // Ensure dark mode is active for every test (integration and acceptance alike).
  // Two mechanisms are needed:
  //   1. data-theme="dark" on <html> — covers integration tests where ThemeService
  //      is never instantiated, so _applyTheme() never runs.
  //   2. vault:theme="dark" in localStorage — covers acceptance tests where the
  //      full app boots and ThemeService._resolveInitialTheme() runs before the
  //      constructor calls _applyTheme(). Without the stored value, _resolveInitialTheme
  //      returns "system", isDarkMode defers to matchMedia (which is light in CI),
  //      and _applyTheme() removes data-theme, undoing the attribute we just set.
  QUnit.hooks.beforeEach(function () {
    window.localStorage.setItem('vault:theme', 'dark');
    document.documentElement.setAttribute('data-theme', 'dark');
    // ember-qunit hard-codes background-color: #fff on #ember-testing-container.
    // In dark mode axe-core uses that white background when computing contrast for
    // elements rendered inside the container. Override to the dark page background
    // so axe sees the correct backdrop.
    const container = document.getElementById('ember-testing-container');
    if (container) {
      container.style.setProperty('background-color', '#030303');
    }
    // Portal-rendered elements (ember-power-select dropdowns, HDS dropdown menus,
    // tooltips) are inserted as children of <body>, outside #ember-testing-container.
    // axe-core computes their background by walking up to <body>, so <body> must also
    // be painted dark — otherwise these elements fail against the default white body.
    document.body.style.setProperty('background-color', '#030303');
  });
  QUnit.hooks.afterEach(function () {
    window.localStorage.removeItem('vault:theme');
    document.documentElement.removeAttribute('data-theme');
    const container = document.getElementById('ember-testing-container');
    if (container) {
      container.style.removeProperty('background-color');
    }
    document.body.style.removeProperty('background-color');
  });
}
