import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { format } from 'date-fns';

export default class Dashboard extends Component {
  maxNamespaces = 10;
  chartLegend = [
    { key: 'distinct_entities', label: 'Direct entities' },
    { key: 'non_entity_tokens', label: 'Active direct tokens' },
  ];
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @tracked isEditStartMonthOpen = false;
  @tracked barChartSelection = false;
  @tracked selectedNamespace = null;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked startDate = null;
  @tracked endDate = null;

  constructor() {
    super(...arguments);
    // these will come in from the endpoint or are passed in for now I'm hardcoding?
    this.startDate = 'Jan-2021';
    this.endDate = 'April-2021';
  }

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

  // Create namespaces data for csv format
  // get getCsvData() {
  //   if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
  //     return null;
  //   }
  //   let results = '',
  //     namespaces = this.args.model.activity.byNamespace,
  //     fields = ['Namespace path', 'Active clients', 'Unique entities', 'Non-entity tokens'];

  //   results = fields.join(',') + '\n';

  //   namespaces.forEach(function(item) {
  //     let path = item.namespace_path !== '' ? item.namespace_path : 'root',
  //       total = item.counts.clients,
  //       unique = item.counts.distinct_entities,
  //       non_entity = item.counts.non_entity_tokens;

  //     results += path + ',' + total + ',' + unique + ',' + non_entity + '\n';
  //   });
  //   return results;
  // }

  // Return csv filename with start and end dates
  // get getCsvFileName() {
  //   let defaultFileName = `clients-by-namespace`,
  //     startDate =
  //       this.args.model.queryStart || `${format(new Date(this.args.model.activity.startTime), 'MM-yyyy')}`,
  //     endDate =
  //       this.args.model.queryEnd || `${format(new Date(this.args.model.activity.endTime), 'MM-yyyy')}`;
  //   if (startDate && endDate) {
  //     defaultFileName += `-${startDate}-${endDate}`;
  //   }
  //   return defaultFileName;
  // }

  // Get the namespace by matching the path from the namespace list
  getNamespace(path) {
    return this.args.model.activity.byNamespace.find(ns => {
      if (path === 'root') {
        return ns.namespace_path === '';
      }
      return ns.namespace_path === path;
    });
  }

  // query Data functions
  async handleQueryData(range) {
    // todo figure out what range should look like.
    // this fires off method on the adapter to query data and return it. await the data's return
  }

  // Edit Start Date modal actions
  @action
  handleEditStartMonth(month, year) {
    if (!month && !year) {
      // reset and do nothing. They pressed cancel
      this.startMonth = this.startYear = null;
      return;
    }
    // if no endDate selected, default to 12 months range
    // create range.
    // then send to queryData this.queryData(range)
  }
  @action
  toggleEditStartMonth() {
    this.isEditStartMonthOpen = !this.isEditStartMonthOpen;
  }
  @action
  selectStartMonth(month) {
    this.startMonth = month;
  }
  @action
  selectStartYear(year) {
    this.startYear = year;
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
