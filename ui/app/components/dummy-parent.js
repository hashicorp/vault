/**
 * @module DummyParentComponent
 * DummyParentComponent components are used to...
 *
 * @example
 * ```js
 * <DummyParent @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/dummy-parent';
import { setComponentTemplate } from '@ember/component';

class DummyParentComponent extends Component {
  dataset = [
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

  get formattedData() {
    return this.flattenData(this.dataset);
  }

  // data sent to bar chart component must be flattened
  flattenData(data) {
    return data.map(d => {
      return {
        label: d['namespace_path'],
        non_entity_tokens: d['counts']['non_entity_tokens'],
        distinct_entities: d['counts']['distinct_entities'],
        total: d['counts']['clients'],
      };
    });
  }
}

export default setComponentTemplate(layout, DummyParentComponent);
