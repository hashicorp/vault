import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const testCases = [
  {
    // default case should show all possible fields for each type
    pluginType: '',
    staticRoleFields: ['name', 'username', 'rotation_period', 'rotation_statements'],
    dynamicRoleFields: [
      'name',
      'ttl',
      'max_ttl',
      'creation_statements',
      'revocation_statements',
      'rollback_statements',
      'renew_statements',
    ],
  },
  {
    pluginType: 'mongodb-database-plugin',
    staticRoleFields: ['username', 'rotation_period'],
    dynamicRoleFields: ['creation_statement', 'revocation_statement', 'ttl', 'max_ttl'],
    statementsHidden: true,
  },
  {
    pluginType: 'mssql-database-plugin',
    staticRoleFields: ['username', 'rotation_period'],
    dynamicRoleFields: ['creation_statements', 'revocation_statements', 'ttl', 'max_ttl'],
  },
];

// used to calculate checks that fields do NOT show up
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
  { name: 'rollback_statements', type: 'string', options: {} },
  { name: 'renew_statements', type: 'string', options: {} },
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

  test('it shows appropriate fields based on roleType and db plugin', async function(assert) {
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
    for (let testCase of testCases) {
      let staticFields = getFields(testCase.staticRoleFields);
      let dynamicFields = getFields(testCase.dynamicRoleFields);
      this.set('dbType', testCase.pluginType);
      this.set('roleType', 'static');
      staticFields.show.forEach(attr => {
        assert
          .dom(`[data-test-input="${attr.name}"]`)
          .exists(
            `${attr.name} attribute exists on static role for ${testCase.pluginType || 'default'} db type`
          );
      });
      staticFields.hide.forEach(attr => {
        assert
          .dom(`[data-test-input="${attr.name}"]`)
          .doesNotExist(
            `${attr.name} attribute does not exist on static role for ${testCase.pluginType ||
              'default'} db type`
          );
      });
      if (testCase.statementsHidden) {
        assert
          .dom('[data-test-statements-section]')
          .doesNotExist(`Statements section is hidden for static ${testCase.pluginType} role`);
      }
      this.set('roleType', 'dynamic');
      dynamicFields.show.forEach(attr => {
        assert
          .dom(`[data-test-input="${attr.name}"]`)
          .exists(
            `${attr.name} attribute exists on dynamic role for ${testCase.pluginType || 'default'} db type`
          );
      });
      dynamicFields.hide.forEach(attr => {
        assert
          .dom(`[data-test-input="${attr.name}"]`)
          .doesNotExist(
            `${attr.name} attribute does not exist on dynamic role for ${testCase.pluginType ||
              'default'} db type`
          );
      });
    }
  });
});
