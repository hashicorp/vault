import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { format } from 'date-fns';

export default class Dashboard extends Component {
  maxNamespaces = 10;
  chartLegend = [
    { key: 'distinct_entities', label: 'unique entities' },
    { key: 'non_entity_tokens', label: 'non-entity tokens' },
  ];
  @tracked selectedNamespace = null;

  @tracked barChartSelection = false;

  // Determine if we have client count data based on the current tab
  get hasClientData() {
    if (this.args.tab === 'current') {
      // Show the current numbers as long as config is on
      return this.args.model.config?.enabled !== 'Off';
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
    let dataset = this.args.model.activity.byNamespace.slice(0, this.maxNamespaces);
    return dataset.map(d => {
      return {
        label: d['namespace_path'] === '' ? 'root' : d['namespace_path'],
        // the order here determines which data is the left bar and which is the right
        distinct_entities: d['counts']['distinct_entities'],
        non_entity_tokens: d['counts']['non_entity_tokens'],
        total: d['counts']['clients'],
      };
    });
  }

  // TODO: dataset for line chart
  get lineChartData() {
    return 'hi';
  }

  // TODO: dataset for new monthly clients vertical bar chart (manage in serializer?)
  get newMonthlyClients() {
    return [
      { month: 'January', directEntities: 1000, nonEntityTokens: 322, total: 1322 },
      { month: 'February', directEntities: 1500, nonEntityTokens: 122, total: 1622 },
      { month: 'March', directEntities: 4300, nonEntityTokens: 700, total: 5000 },
      { month: 'April', directEntities: 1550, nonEntityTokens: 229, total: 1779 },
      { month: 'May', directEntities: 5560, nonEntityTokens: 124, total: 5684 },
      { month: 'June', directEntities: 1570, nonEntityTokens: 142, total: 1712 },
      { month: 'July', directEntities: 300, nonEntityTokens: 112, total: 412 },
      { month: 'August', directEntities: 1610, nonEntityTokens: 130, total: 1740 },
      { month: 'September', directEntities: 1900, nonEntityTokens: 222, total: 2122 },
      { month: 'October', directEntities: 500, nonEntityTokens: 166, total: 666 },
      { month: 'November', directEntities: 480, nonEntityTokens: 132, total: 612 },
      { month: 'December', directEntities: 980, nonEntityTokens: 202, total: 1182 },
    ];
  }

  // TODO: dataset for vault usage vertical bar chart (manage in serializer?)
  get totalClientUsage() {
    return [
      { month: 'January', directEntities: 1000, nonEntityTokens: 322, total: 1322 },
      { month: 'February', directEntities: 1500, nonEntityTokens: 122, total: 1622 },
      { month: 'March', directEntities: 4300, nonEntityTokens: 700, total: 5000 },
      { month: 'April', directEntities: 1550, nonEntityTokens: 229, total: 1779 },
      { month: 'May', directEntities: 5560, nonEntityTokens: 124, total: 5684 },
      { month: 'June', directEntities: 1570, nonEntityTokens: 142, total: 1712 },
      { month: 'July', directEntities: 300, nonEntityTokens: 112, total: 412 },
      { month: 'August', directEntities: 1610, nonEntityTokens: 130, total: 1740 },
      { month: 'September', directEntities: 1900, nonEntityTokens: 222, total: 2122 },
      { month: 'October', directEntities: 500, nonEntityTokens: 166, total: 666 },
      { month: 'November', directEntities: 480, nonEntityTokens: 132, total: 612 },
      { month: 'December', directEntities: 980, nonEntityTokens: 202, total: 1182 },
    ];
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

  // Return csv filename with start and end dates
  get getCsvFileName() {
    let defaultFileName = `clients-by-namespace`,
      startDate =
        this.args.model.queryStart || `${format(new Date(this.args.model.activity.startTime), 'MM-yyyy')}`,
      endDate =
        this.args.model.queryEnd || `${format(new Date(this.args.model.activity.endTime), 'MM-yyyy')}`;
    if (startDate && endDate) {
      defaultFileName += `-${startDate}-${endDate}`;
    }
    return defaultFileName;
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
