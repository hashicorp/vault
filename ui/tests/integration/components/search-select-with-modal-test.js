import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { clickTrigger, typeInSearch } from 'ember-power-select/test-support/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import ss from 'vault/tests/pages/components/search-select';

const component = create(ss);

module('Integration | Component | search select with modal', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.get('identity/entity/id', () => {
      return {
        request_id: 'entity-list-id',
        data: {
          key_info: {
            'entity-1-id': {
              name: 'entity-1',
            },
            'entity-2-id': {
              name: 'entity-2',
            },
          },
          keys: ['entity-1-id', 'entity-2-id'],
        },
      };
    });
    this.server.get('sys/policies/acl', () => {
      return {
        request_id: 'acl-policy-list-id',
        data: {
          keys: ['default', 'root'],
        },
      };
    });
    this.server.get('sys/policies/rgp', () => {
      return {
        request_id: 'rgp-policy-list-id',
        data: {
          keys: ['default', 'root'],
        },
      };
    });
    this.server.get('/identity/entity/id/entity-1-id', () => {
      return {
        request_id: 'some-entity-id-1',
        data: {
          id: 'entity-1-id',
          name: 'entity-1',
          namespace_id: 'root',
          policies: ['default'],
        },
      };
    });
    this.server.get('/identity/entity/id/entity-2-id', () => {
      return {
        request_id: 'some-entity-id-2',
        data: {
          id: 'entity-2-id',
          name: 'entity-2',
          namespace_id: 'root',
          policies: ['default'],
        },
      };
    });
  });

  test('it renders passed in model', async function (assert) {
    await render(hbs`
    <SearchSelectWithModal
      @id="entity"
      @label="entity"
      @subText="Search for an existing entity, or type a new name to create it."
      @model="identity/entity"
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormComponent="identity/edit-form"
      @modalSubtext="Some modal subtext"
      />
      <div id="modal-wormhole"></div>
  `);

    assert.dom('[data-test-search-select-with-modal]').exists('the component renders');
    assert.equal(component.labelText, 'Entity name', 'label text is correct');
    assert.ok(component.hasTrigger, 'it renders the power select trigger');
    assert.equal(component.selectedOptions.length, 0, 'there are no selected options');

    await clickTrigger();
    assert.equal(component.options.length, 2, 'dropdown renders passed in models as options');
  });
  test('it filters options and adds option to create new item', async function (assert) {
    assert.expect(7);
    await render(hbs`
    <SearchSelectWithModal
      @id="entity"
      @label="entity"
      @subText="Search for an existing entity, or type a new name to create it."
      @model="identity/entity"
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormComponent="identity/edit-form"
      @modalSubtext="Some modal subtext"
      />
      <div id="modal-wormhole"></div>
  `);

    await clickTrigger();
    assert.equal(component.options.length, 2, 'dropdown renders all options');

    await typeInSearch('e');
    assert.equal(component.options.length, 3, 'dropdown renders all options plus add option');

    await typeInSearch('entity-1');
    assert.equal(component.options[0].text, 'entity-1-id', 'dropdown renders only matching option');

    await typeInSearch('entity-1-new');
    assert.equal(
      component.options[0].text,
      'Create new entity: entity-1-new',
      'dropdown gives option to create new option'
    );

    await component.selectOption();
    assert.dom('[data-test-modal-div]').hasAttribute('class', 'modal is-info is-active', 'modal is active');
    assert.dom('[data-test-modal-subtext]').hasText('Some modal subtext', 'renders modal text');
    assert.dom('[data-test-component="identity-edit-form"]').exists('renders identity form');
  });
});
