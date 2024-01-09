/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | transform-list-item', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders un-clickable item if no read capability', async function (assert) {
    const item = EmberObject.create({
      id: 'foo',
      updatePath: {
        canRead: false,
        canDelete: true,
        canUpdate: true,
      },
    });
    this.set('itemPath', 'role/foo');
    this.set('itemType', 'role');
    this.set('item', item);
    await render(hbs`<SecretList::TransformListItem
      @item={{this.item}}
      @itemPath={{this.itemPath}}
      @itemType={{this.itemType}}
    />`);

    assert.dom('[data-test-view-only-list-item]').exists('shows view only list item');
    assert.dom('[data-test-view-only-list-item]').hasText(item.id, 'has correct label');
  });

  test('it is clickable with details menu item if read capability', async function (assert) {
    const item = EmberObject.create({
      id: 'foo',
      updatePath: {
        canRead: true,
        canDelete: false,
        canUpdate: false,
      },
    });
    this.set('itemPath', 'template/foo');
    this.set('itemType', 'template');
    this.set('item', item);
    await render(hbs`<SecretList::TransformListItem
      @item={{this.item}}
      @itemPath={{this.itemPath}}
      @itemType={{this.itemType}}
    />`);

    assert.dom('[data-test-secret-link="template/foo"]').exists('shows clickable list item');
    await click('button.popup-menu-trigger');
    assert.dom('.popup-menu-content li').exists({ count: 1 }, 'has one option');
  });

  test('it has details and edit menu item if read & edit capabilities', async function (assert) {
    const item = EmberObject.create({
      id: 'foo',
      updatePath: {
        canRead: true,
        canDelete: true,
        canUpdate: true,
      },
    });
    this.set('itemPath', 'alphabet/foo');
    this.set('itemType', 'alphabet');
    this.set('item', item);
    await render(hbs`<SecretList::TransformListItem
      @item={{this.item}}
      @itemPath={{this.itemPath}}
      @itemType={{this.itemType}}
    />`);

    assert.dom('[data-test-secret-link="alphabet/foo"]').exists('shows clickable list item');
    await click('button.popup-menu-trigger');
    assert.dom('.popup-menu-content li').exists({ count: 2 }, 'has both options');
  });

  test('it is not clickable if built-in template with all capabilities', async function (assert) {
    const item = EmberObject.create({
      id: 'builtin/foo',
      updatePath: {
        canRead: true,
        canDelete: true,
        canUpdate: true,
      },
    });
    this.set('itemPath', 'template/builtin/foo');
    this.set('itemType', 'template');
    this.set('item', item);
    await render(hbs`<SecretList::TransformListItem
      @item={{this.item}}
      @itemPath={{this.itemPath}}
      @itemType={{this.itemType}}
    />`);

    assert.dom('[data-test-view-only-list-item]').exists('shows view only list item');
    assert.dom('[data-test-view-only-list-item]').hasText(item.id, 'has correct label');
  });

  test('it is not clickable if built-in alphabet', async function (assert) {
    const item = EmberObject.create({
      id: 'builtin/foo',
      updatePath: {
        canRead: true,
        canDelete: true,
        canUpdate: true,
      },
    });
    this.set('itemPath', 'alphabet/builtin/foo');
    this.set('itemType', 'alphabet');
    this.set('item', item);
    await render(hbs`<SecretList::TransformListItem
      @item={{this.item}}
      @itemPath={{this.itemPath}}
      @itemType={{this.itemType}}
    />`);

    assert.dom('[data-test-view-only-list-item]').exists('shows view only list item');
    assert.dom('[data-test-view-only-list-item]').hasText(item.id, 'has correct label');
  });
});
