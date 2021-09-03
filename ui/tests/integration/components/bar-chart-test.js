import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | bar-chart', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    let dataset = [
      {
        namespace_id: 'root',
        namespace_path: 'root',
        counts: {
          distinct_entities: 268,
          non_entity_tokens: 985,
          clients: 1253,
        },
      },
      {
        namespace_id: 'O0i4m',
        namespace_path: 'top-namespace',
        counts: {
          distinct_entities: 648,
          non_entity_tokens: 220,
          clients: 868,
        },
      },
      {
        namespace_id: '1oihz',
        namespace_path: 'anotherNamespace',
        counts: {
          distinct_entities: 547,
          non_entity_tokens: 337,
          clients: 884,
        },
      },
      {
        namespace_id: '1oihz',
        namespace_path: 'someOtherNamespaceawgagawegawgawgawgaweg',
        counts: {
          distinct_entities: 807,
          non_entity_tokens: 234,
          clients: 1041,
        },
      },
    ];

    let flattenData = () => {
      return dataset.map(d => {
        return {
          label: d['namespace_path'],
          non_entity_tokens: d['counts']['non_entity_tokens'],
          distinct_entities: d['counts']['distinct_entities'],
          total: d['counts']['clients'],
        };
      });
    };

    this.set('title', 'Top Namespaces');
    this.set('description', 'Each namespaces client count includes clients in child namespaces.');
    this.set('dataset', flattenData());

    await render(hbs`
      <BarChart 
        @title={{title}}
        @description={{description}}
        @dataset={{dataset}}
        @mapLegend={{array
        (hash key="non_entity_tokens" label="Active direct tokens")
        (hash key="distinct_entities" label="Unique Entities")}}
        >    
          <button type="button" class="link">
          Export all namespace data
          </button>
      </BarChart>
    `);

    assert.dom('.bar-chart-wrapper').exists('it renders');
  });
});
