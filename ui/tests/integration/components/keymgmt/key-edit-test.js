import { module, test } from 'qunit';
import EmberObject from '@ember/object';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | keymgmt/key-edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const now = new Date().toString();
    let model = EmberObject.create({
      name: 'Unicorns',
      id: 'Unicorns',
      minEnabledVersion: 1,
      versions: [
        {
          id: 1,
          creation_time: now,
        },
        {
          id: 2,
          creation_time: now,
        },
      ],
      canDelete: true,
    });
    this.model = model;
    this.tab = '';
  });

  // TODO: Add capabilities tests
  test('it renders show view as default', async function (assert) {
    assert.expect(8);
    await render(hbs`<Keymgmt::KeyEdit @model={{model}} @tab={{tab}} /><div id="modal-wormhole" />`);
    assert.dom('[data-test-secret-header]').hasText('Unicorns', 'Shows key name');
    assert.dom('[data-test-keymgmt-key-toolbar]').exists('Subnav toolbar exists');
    assert.dom('[data-test-tab="Details"]').exists('Details tab exists');
    assert.dom('[data-test-tab="Versions"]').exists('Versions tab exists');
    assert.dom('[data-test-keymgmt-key-destroy]').isDisabled('Destroy button is disabled');
    assert.dom('[data-test-keymgmt-dist-empty-state]').exists('Distribution empty state exists');

    this.set('tab', 'versions');
    assert.dom('[data-test-keymgmt-key-version]').exists({ count: 2 }, 'Renders two version list items');
    assert
      .dom('[data-test-keymgmt-key-current-min]')
      .exists({ count: 1 }, 'Checks only one as current minimum');
  });

  test('it renders the correct elements on edit view', async function (assert) {
    assert.expect(4);
    let model = EmberObject.create({
      name: 'Unicorns',
      id: 'Unicorns',
    });
    this.set('mode', 'edit');
    this.set('model', model);

    await render(hbs`<Keymgmt::KeyEdit @model={{model}} @mode={{mode}} /><div id="modal-wormhole" />`);
    assert.dom('[data-test-secret-header]').hasText('Edit key', 'Shows edit header');
    assert.dom('[data-test-keymgmt-key-toolbar]').doesNotExist('Subnav toolbar does not exist');
    assert.dom('[data-test-tab="Details"]').doesNotExist('Details tab does not exist');
    assert.dom('[data-test-tab="Versions"]').doesNotExist('Versions tab does not exist');
  });

  test('it renders the correct elements on create view', async function (assert) {
    assert.expect(4);
    let model = EmberObject.create({});
    this.set('mode', 'create');
    this.set('model', model);

    await render(hbs`<Keymgmt::KeyEdit @model={{model}} @mode={{mode}} /><div id="modal-wormhole" />`);
    assert.dom('[data-test-secret-header]').hasText('Create key', 'Shows edit header');
    assert.dom('[data-test-keymgmt-key-toolbar]').doesNotExist('Subnav toolbar does not exist');
    assert.dom('[data-test-tab="Details"]').doesNotExist('Details tab does not exist');
    assert.dom('[data-test-tab="Versions"]').doesNotExist('Versions tab does not exist');
  });
});
