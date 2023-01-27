import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

const selectors = {
  keys: '[data-test-toggle-group="Key parameters"]',
  sanOptions: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  subjectFields: '[data-test-toggle-group="Additional subject fields"]',
  toggleByName: (name) => `[data-test-toggle-group="${name}"]`,
};

module('Integration | Component | PkiGenerateToggleGroups', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.model = this.owner
      .lookup('service:store')
      .createRecord('pki/action', { actionType: 'generate-root' });
  });

  test('it should render key parameters', async function (assert) {
    await render(hbs`<PkiGenerateToggleGroups @model={{this.model}} />`, { owner: this.engine });

    assert.dom(selectors.keys).hasText('Key parameters', 'Key parameters group renders');

    await click(selectors.keys);

    assert
      .dom('[data-test-toggle-group-description]')
      .hasText(
        'Please choose a type to see key parameter options.',
        'Placeholder renders for key params when type is not selected'
      );
    const fields = {
      exported: ['keyName', 'keyType', 'keyBits'],
      internal: ['keyName', 'keyType', 'keyBits'],
      existing: ['keyRef'],
      kms: ['keyName', 'managedKeyName', 'managedKeyId'],
    };
    for (const type in fields) {
      this.model.type = type;
      await settled();
      assert
        .dom('[data-test-field]')
        .exists({ count: fields[type].length }, `Correct number of fields render for ${type} type`);
      fields[type].forEach((key) => {
        assert.dom(`[data-test-input="${key}"]`).exists(`${key} input renders for ${type} type`);
      });
    }
  });

  test('it should render SAN options', async function (assert) {
    await render(hbs`<PkiGenerateToggleGroups @model={{this.model}} />`, { owner: this.engine });

    assert
      .dom(selectors.sanOptions)
      .hasText('Subject Alternative Name (SAN) Options', 'SAN options group renders');

    await click(selectors.sanOptions);

    const fields = ['excludeCnFromSans', 'serialNumber', 'altNames', 'ipSans', 'uriSans', 'otherSans'];
    assert.dom('[data-test-field]').exists({ count: 6 }, `Correct number of fields render`);
    fields.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} input renders for generate-root actionType`);
    });

    this.model.actionType = 'generate-csr';
    await settled();

    assert
      .dom('[data-test-field]')
      .exists({ count: 4 }, 'Correct number of fields render for generate-csr actionType');

    assert
      .dom('[data-test-input="excludeCnFromSans"]')
      .doesNotExist('excludeCnFromSans field hidden for generate-csr actionType');
    assert
      .dom('[data-test-input="serialNumber"]')
      .doesNotExist('serialNumber field hidden for generate-csr actionType');
  });

  test('it should render additional subject fields', async function (assert) {
    await render(hbs`<PkiGenerateToggleGroups @model={{this.model}} />`, { owner: this.engine });

    assert.dom(selectors.subjectFields).hasText('Additional subject fields', 'SAN options group renders');

    await click(selectors.subjectFields);

    const fields = ['ou', 'organization', 'country', 'locality', 'province', 'streetAddress', 'postalCode'];
    assert.dom('[data-test-field]').exists({ count: fields.length }, 'Correct number of fields render');
    fields.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} input renders`);
    });
  });

  test('it should render groups according to the passed @groups', async function (assert) {
    assert.expect(11);
    const fieldsA = ['ou', 'organization'];
    const fieldsZ = ['country', 'locality', 'province', 'streetAddress', 'postalCode'];
    this.set('groups', {
      'Group A': fieldsA,
      'Group Z': fieldsZ,
    });
    await render(hbs`<PkiGenerateToggleGroups @model={{this.model}} @groups={{this.groups}} />`, {
      owner: this.engine,
    });

    assert.dom(selectors.toggleByName('Group A')).hasText('Group A', 'First group renders');
    assert.dom(selectors.toggleByName('Group Z')).hasText('Group Z', 'Second group renders');

    await click(selectors.toggleByName('Group A'));
    assert.dom('[data-test-field]').exists({ count: fieldsA.length }, 'Correct number of fields render');
    fieldsA.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} input renders`);
    });

    await click(selectors.toggleByName('Group Z'));
    assert.dom('[data-test-field]').exists({ count: fieldsZ.length }, 'Correct number of fields render');
    fieldsZ.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} input renders`);
    });
  });
});
