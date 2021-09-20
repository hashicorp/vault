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
