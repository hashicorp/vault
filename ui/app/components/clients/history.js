import Component from '@glimmer/component';

export default class HistoryComponent extends Component {
  max_namespaces = 10;

  get hasClientData() {
    if (this.args.tab === 'current') {
      return this.args.model.activity && this.args.model.activity.clients;
    }
    return this.args.model.activity && this.args.model.activity.total;
  }

  get barChartDataset() {
    // let dataset = [
    // {label: 'namespace2/', non_entity_tokens: 1548, distinct_entities: 1269, total: 2817},
    // {label: 'namespace8/', non_entity_tokens: 1119, distinct_entities: 1591, total: 2710},
    // {label: 'namespace5/', non_entity_tokens: 1579, distinct_entities: 943, total: 2522},
    // {label: 'namespace1/', non_entity_tokens: 939, distinct_entities: 1460, total: 2399},
    // {label: 'namespace6/', non_entity_tokens: 1042, distinct_entities: 1290, total: 2332},
    // {label: 'namespace7/', non_entity_tokens: 1372, distinct_entities: 892, total: 2264},
    // {label: 'namespace3/', non_entity_tokens: 1179, distinct_entities: 993, total: 2172},
    // {label: 'namespace9/', non_entity_tokens: 501, distinct_entities: 1453, total: 1954},
    // {label: 'ns1/', non_entity_tokens: 1160, distinct_entities: 660, total: 1820},
    // {label: 'namespacereallyreallylong/', non_entity_tokens: 847, distinct_entities: 813, total: 1660}]
    // return dataset
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    let dataset = this.args.model.activity.byNamespace;
    // Filter out root data
    dataset = dataset.filter(item => {
      return item.namespace_id !== 'root';
    });
    // Show only top 10 namespaces
    dataset = dataset.slice(0, this.max_namespaces);
    return dataset.map(d => {
      return {
        label: d['namespace_path'],
        non_entity_tokens: d['counts']['non_entity_tokens'],
        distinct_entities: d['counts']['distinct_entities'],
        total: d['counts']['clients'],
      };
    });
  }

  get showGraphs() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    return this.args.model.activity.byNamespace.length > 1;
  }
}
