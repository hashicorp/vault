import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const STANDARD_FIELDS = [
  'name',
  'ttl',
  'max_ttl',
  'username',
  'rotation_period',
  'creation_statements',
  'revocation_statements',
  'rotation_statements',
];
const MONGODB_STATIC_FIELDS = ['username', 'rotation_period'];
const MONGODB_DYNAMIC_FIELDS = ['creation_statement', 'revocation_statement', 'ttl', 'max_ttl'];
const ALL_ATTRS = [
  { name: 'ttl', type: 'string', options: {} },
  { name: 'max_ttl', type: 'string', options: {} },
  { name: 'username', type: 'string', options: {} },
  { name: 'rotation_period', type: 'string', options: {} },
  { name: 'creation_statements', type: 'string', options: {} },
  { name: 'creation_statement', type: 'string', options: {} },
  { name: 'revocation_statements', type: 'string', options: {} },
  { name: 'revocation_statement', type: 'string', options: {} },
  { name: 'rotation_statements', type: 'string', options: {} },
];
const getFields = nameArray => {
  const show = ALL_ATTRS.filter(attr => nameArray.indexOf(attr.name) >= 0);
  const hide = ALL_ATTRS.filter(attr => nameArray.indexOf(attr.name) < 0);
  return { show, hide };
};

module('Integration | Component | database-role-setting-form', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set(
      'model',
      EmberObject.create({
        // attrs is not its own set value b/c ember hates arrays as args
        attrs: ALL_ATTRS,
      })
    );
  });

  test('it shows empty states when no roleType passed in', async function(assert) {
    await render(hbs`<DatabaseRoleSettingForm @attrs={{model.attrs}} @model={{model}}/>`);
    assert.dom('[data-test-component="empty-state"]').exists({ count: 2 }, 'Two empty states exist');
  });

  test('it shows appropriate fields based on roleType with default db type', async function(assert) {
    this.set('roleType', 'static');
    this.set('dbType', '');
    await render(hbs`
      <DatabaseRoleSettingForm
        @attrs={{model.attrs}}
        @model={{model}}
        @roleType={{roleType}}
        @dbType={{dbType}}
      />
    `);
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Does not show empty states');
    const defaultFields = getFields(STANDARD_FIELDS);
    defaultFields.show.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .exists(`${attr.name} attribute exists for default db type`);
    });
    defaultFields.hide.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .doesNotExist(`${attr.name} attribute does not exist for default db type`);
    });

    this.set('roleType', 'dynamic');
    defaultFields.show.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .exists(`${attr.name} attribute exists for default db with dynamic type`);
    });
    defaultFields.hide.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .doesNotExist(`${attr.name} attribute does not exist for default db with dynamic type`);
    });
  });
  test('it shows appropriate fields based on roleType with mongodb', async function(assert) {
    this.set('roleType', 'static');
    this.set('dbType', 'mongodb-database-plugin');
    await render(hbs`
      <DatabaseRoleSettingForm
        @attrs={{model.attrs}}
        @model={{model}}
        @roleType={{roleType}}
        @dbType={{dbType}}
      />
    `);
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Does not show empty states');
    const staticFields = getFields(MONGODB_STATIC_FIELDS);
    staticFields.show.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .exists(`${attr.name} attribute exists for mongodb static role`);
    });
    staticFields.hide.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .doesNotExist(`${attr.name} attribute does not exist for mongodb static role`);
    });
    assert
      .dom('[data-test-statements-section]')
      .doesNotExist('Statements section is hidden for dynamic mongodb role');

    this.set('roleType', 'dynamic');
    const dynamicFields = getFields(MONGODB_DYNAMIC_FIELDS);
    dynamicFields.show.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .exists(`${attr.name} attribute exists for mongodb dynamic role`);
    });
    dynamicFields.hide.forEach(attr => {
      assert
        .dom(`[data-test-input="${attr.name}"]`)
        .doesNotExist(`${attr.name} attribute does not exist for mongodb dynamic role`);
    });
  });
});
