/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setRunOptions, shouldForceAudit } from 'ember-a11y-testing/test-support';
import { module, test } from 'qunit';
import config from 'vault/config/environment';
import { setupRenderingTest } from 'vault/tests/helpers';

const auditModule = config.APP.DARK_THEME_AUDIT || !shouldForceAudit() ? module.skip : module;

auditModule('Integration | Accessibility audit', function (hooks) {
  setupRenderingTest(hooks);

  test('rendering enforces WCAG A rules as well as AA rules', async function (assert) {
    await assert.rejects(
      render(hbs`<button type="button"></button>`),
      /button-name/,
      'The global audit rejects a button without an accessible name'
    );
  });

  test('rendering enforces WCAG AA contrast requirements', async function (assert) {
    await assert.rejects(
      render(hbs`<p style="color: #aaa; background-color: #fff;">Low contrast text</p>`),
      /color-contrast/,
      'The global audit rejects text below the AA contrast threshold'
    );
  });

  test('per-test options preserve rule exceptions without replacing the A+AA policy', async function (assert) {
    setRunOptions({
      runOnly: { type: 'tag', values: ['wcag2a'] },
      rules: { label: { enabled: false } },
    });

    await render(hbs`<input data-test-unlabelled-input type="text" />`);
    assert.dom('[data-test-unlabelled-input]').exists('The explicit label exception is preserved');

    await assert.rejects(
      render(hbs`<p style="color: #aaa; background-color: #fff;">Low contrast text</p>`),
      /color-contrast/,
      'Per-test options cannot silently remove AA coverage'
    );
  });

  test('explicit contrast exceptions are preserved', async function (assert) {
    setRunOptions({
      rules: { 'color-contrast': { enabled: false } },
    });

    await render(
      hbs`<p data-test-low-contrast style="color: #aaa; background-color: #fff;">Low contrast text</p>`
    );
    assert
      .dom('[data-test-low-contrast]')
      .hasText('Low contrast text', 'The contrast exception is preserved');
  });

  test('interactions audit the resulting state', async function (assert) {
    this.set('lowContrast', false);
    this.set('showLowContrast', () => this.set('lowContrast', true));

    await render(hbs`
      <button type="button" data-test-show-low-contrast {{on "click" this.showLowContrast}}>
        Show low contrast text
      </button>
      {{#if this.lowContrast}}
        <p style="color: #aaa; background-color: #fff;">Low contrast text</p>
      {{/if}}
    `);

    await assert.rejects(
      click('[data-test-show-low-contrast]'),
      /color-contrast/,
      'The global audit rejects an AA violation introduced after rendering'
    );
  });
});
