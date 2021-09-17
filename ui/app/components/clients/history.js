import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class HistoryComponent extends Component {
  max_namespaces = 10;

  @tracked selectedNamespace = null;

  // Determine if we have client count data based on the current tab,
  // since model is slightly different for current month vs history api
  get hasClientData() {
    if (this.args.tab === 'current') {
      return this.args.model.activity && this.args.model.activity.clients;
    }
    return this.args.model.activity && this.args.model.activity.total;
  }

  // Show namespace graph only if we have more than 1
  get showGraphs() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    return this.args.model.activity.byNamespace.length > 1;
  }

  // Construct the namespace model for the search select component
  get searchDataset() {
    let dataList = this.cleanUpNamespaces();
    if (!dataList) {
      return null;
    }
    return dataList.map(d => {
      return {
        name: d['namespace_id'],
        id: d['namespace_path'] || 'root',
      };
    });
  }

  // Construct the namespace model for the car chart component
  get barChartDataset() {
    let dataset = this.cleanUpNamespaces();
    if (!dataset) {
      return null;
    }
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

  // Filter out root data
  cleanUpNamespaces() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    let namespaces = this.args.model.activity.byNamespace;
    namespaces = namespaces.filter(item => {
      return item.namespace_id !== 'root';
    });
    return namespaces;
  }

  // Set the namespace from the search select picker
  setNamespace(path) {
    this.selectedNamespace = this.args.model.activity.byNamespace.find(ns => {
      return ns.namespace_path === path;
    });
  }

  @action
  initNamespace(value) {
    if (value && value.length) {
      this.setNamespace(value[0]);
    } else {
      this.selectedNamespace = null;
    }
  }
}
