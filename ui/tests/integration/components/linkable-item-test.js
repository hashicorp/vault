/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | linkable-item', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders anything passed in', async function (assert) {
    await render(hbs`<LinkableItem />`);
    assert.dom(this.element).hasText('', 'No content rendered');

    await render(hbs`
      <LinkableItem as |Li|>
        <Li.content>
          stuff here
        </Li.content>
        <Li.menu>
          menu
        </Li.menu>
      </LinkableItem>
    `);
    assert.dom('[data-test-linkable-item-content]').hasText('stuff here');
    assert.dom('[data-test-linkable-item-menu]').hasText('menu');
  });

  test('it is not wrapped in a linked block if disabled is true', async function (assert) {
    await render(hbs`
      <LinkableItem @disabled={{true}} as |Li|>
        <Li.content>
          stuff here
        </Li.content>
      </LinkableItem>
    `);
    assert.dom('.list-item-row').exists('List item row exists');
    assert.dom('.list-item-row.linked-block').doesNotExist('Does not render linked block');
    assert.dom('[data-test-secret-path]').doesNotExist('Title is not rendered');
    assert.dom('[data-test-linkable-item-accessor]').doesNotExist('Accessor is not rendered');
    assert.dom('[data-test-linkable-item-accessor]').doesNotExist('Accessor is not rendered');
    assert.dom('[data-test-linkable-item-glyph]').doesNotExist('Glyph is not rendered');
  });

  test('it is wrapped in a linked block if a link is passed', async function (assert) {
    await render(hbs`
      <LinkableItem @link={{hash route="vault" model="modelId"}} as |Li|>
        <Li.content
          @title="A title"
          @link={{hash route="vault" model="modelId"}}
        >
          stuff here
        </Li.content>
      </LinkableItem>
    `);

    assert.dom('.list-item-row.linked-block').exists('Renders linked block');
  });

  test('it renders standard attributes on content', async function (assert) {
    this.set('title', 'A Title');
    this.set('accessor', 'my accessor');
    this.set('description', 'my description');
    this.set('glyph', 'key');
    this.set('glyphText', 'Here is some extra info');

    // Template block usage:
    await render(hbs`
      <LinkableItem data-test-example as |Li|>
        <Li.content
          @accessor={{this.accessor}}
          @description={{this.description}}
          @glyph={{this.glyph}}
          @glyphText={{this.glyphText}}
          @title={{this.title}}
        />
      </LinkableItem>
    `);
    assert.dom('.list-item-row').exists('List item row exists');
    assert.dom('[data-test-secret-path]').hasText(this.title, 'Title is rendered');
    assert.dom('[data-test-linkable-item-accessor]').hasText(this.accessor, 'Accessor is rendered');
    assert.dom('[data-test-linkable-item-description]').hasText(this.description, 'Description is rendered');
    assert.dom('[data-test-linkable-item-glyph]').exists('Glyph is rendered');
  });
});
