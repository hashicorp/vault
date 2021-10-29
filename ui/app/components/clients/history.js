import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class HistoryComponent extends Component {
  max_namespaces = 10;

  @tracked selectedNamespace = null;

  @tracked barChartSelection = false;

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
    return (
      this.args.model.activity &&
      this.args.model.activity.byNamespace &&
      this.args.model.activity.byNamespace.length > 1
    );
  }

  // Construct the namespace model for the search select component
  get searchDataset() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    let dataList = this.args.model.activity.byNamespace;
    return dataList.map(d => {
      return {
        name: d['namespace_id'],
        id: d['namespace_path'] === '' ? 'root' : d['namespace_path'],
      };
    });
  }

  // Construct the namespace model for the bar chart component
  get barChartDataset() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    let dataset = this.args.model.activity.byNamespace.slice(0, this.max_namespaces);
    return dataset.map(d => {
      return {
        label: d['namespace_path'] === '' ? 'root' : d['namespace_path'],
        non_entity_tokens: d['counts']['non_entity_tokens'],
        distinct_entities: d['counts']['distinct_entities'],
        total: d['counts']['clients'],
      };
    });
  }

  // Create namespaces data for csv format
  get getCsvData() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    let results = '',
      namespaces = this.args.model.activity.byNamespace,
      fields = ['Namespace path', 'Active clients', 'Unique entities', 'Non-entity tokens'];

    results = fields.join(',') + '\n';

    namespaces.forEach(function(item) {
      let path = item.namespace_path !== '' ? item.namespace_path : 'root',
        total = item.counts.clients,
        unique = item.counts.distinct_entities,
        non_entity = item.counts.non_entity_tokens;

      results += path + ',' + total + ',' + unique + ',' + non_entity + '\n';
    });
    return results;
  }

  // Get the namespace by matching the path from the namespace list
  getNamespace(path) {
    return this.args.model.activity.byNamespace.find(ns => {
      if (path === 'root') {
        return ns.namespace_path === '';
      }
      return ns.namespace_path === path;
    });
  }

  @action
  selectNamespace(value) {
    // In case of search select component, value returned is an array
    if (Array.isArray(value)) {
      this.selectedNamespace = this.getNamespace(value[0]);
      this.barChartSelection = false;
    } else if (typeof value === 'object') {
      // While D3 bar selection returns an object
      this.selectedNamespace = this.getNamespace(value.label);
      this.barChartSelection = true;
    }
  }

  @action
  resetData() {
    this.barChartSelection = false;
    this.selectedNamespace = null;
  }
}
