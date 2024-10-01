/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, find, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const SHARED_STYLES = {
  success: {
    icon: 'check-circle-fill',
    class: 'hds-alert--color-success',
  },
  warning: {
    icon: 'alert-triangle-fill',
    class: 'hds-alert--color-warning',
  },
};
module('Integration | Component | alert-inline', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders alert message for each @color arg', async function (assert) {
    const COLORS = {
      ...SHARED_STYLES,
      neutral: {
        icon: 'info-fill',
        class: 'hds-alert--color-neutral',
      },
      highlight: {
        icon: 'info-fill',
        class: 'hds-alert--color-highlight',
      },
      critical: {
        icon: 'alert-diamond-fill',
        class: 'hds-alert--color-critical',
      },
    };

    const { neutral } = COLORS; // default color
    await render(hbs`<AlertInline @message="some very important alert" />`);
    assert.dom('[data-test-inline-error-message]').hasText('some very important alert');
    assert.dom(`[data-test-icon="${neutral.icon}"]`).exists('renders default icon');
    assert.dom('[data-test-inline-alert]').hasClass(neutral.class, 'renders default class');

    // assert deprecated @type arg values map to expected color
    for (const type in COLORS) {
      this.color = type;
      const color = COLORS[type];
      await render(hbs`<AlertInline @color={{this.color}}  @message="some very important alert" />`);
      assert.dom(`[data-test-icon="${color.icon}"]`).exists(`@color="${type}" renders icon: ${color.icon}`);
      assert
        .dom('[data-test-inline-alert]')
        .hasClass(color.class, `@color="${type}" renders class: ${color.class}`);
    }
  });

  test('it renders alert color for each deprecated @type arg', async function (assert) {
    const OLD_TYPES = {
      ...SHARED_STYLES,
      info: {
        icon: 'info-fill',
        class: 'hds-alert--color-highlight',
      },
      danger: {
        icon: 'alert-diamond-fill',
        class: 'hds-alert--color-critical',
      },
    };
    // assert deprecated @type arg values map to expected color
    for (const type in OLD_TYPES) {
      this.type = type;
      const color = OLD_TYPES[type];
      await render(hbs`<AlertInline @type={{this.type}}  @message="some very important alert" />`);
      assert
        .dom(`[data-test-icon="${color.icon}"]`)
        .exists(`deprecated @type="${type}" renders icon: ${color.icon}`);
      assert
        .dom('[data-test-inline-alert]')
        .hasClass(color.class, `deprecated @type="${type}" renders class: ${color.class}`);
    }
  });

  test('it mimics loading when message changes', async function (assert) {
    this.message = 'some very important alert';
    await render(hbs`
    <AlertInline @message={{this.message}}/>
    `);
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('some very important alert', 'it renders original message');

    this.set('message', 'some changed alert!!!');
    await waitUntil(() => find('[data-test-icon="loading"]'));
    assert.ok(find('[data-test-icon="loading"]'), 'it shows loading icon when message changes');
    await settled();
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('some changed alert!!!', 'it shows updated message');
  });
});
