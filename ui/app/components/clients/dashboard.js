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
  // For startDate Modal
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked startDate = null;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;

  constructor() {
    super(...arguments);
    // ARG TODO will need to get startDate from license endpoint
    let date = new Date();
    date.setMonth(date.getMonth() - 12); // by default start date is 12 months from now
    this.startDate = format(date, 'MMMM yyyy');
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
    return dataList.map((d) => {
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
    return dataset.map((d) => {
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
    return [
      { month: '1/21', clients: 100, new: 100 },
      { month: '2/21', clients: 300, new: 200 },
      { month: '3/21', clients: 300, new: 0 },
      { month: '4/21', clients: 300, new: 0 },
      { month: '5/21', clients: 300, new: 0 },
      { month: '6/21', clients: 300, new: 0 },
      { month: '7/21', clients: 300, new: 0 },
      { month: '8/21', clients: 350, new: 50 },
      { month: '9/21', clients: 400, new: 50 },
      { month: '10/21', clients: 450, new: 50 },
      { month: '11/21', clients: 500, new: 50 },
      { month: '12/21', clients: 1000, new: 1000 },
    ];
  }

  // TODO: dataset for new monthly clients vertical bar chart (manage in serializer?)
  get newMonthlyClients() {
    return [
      { month: 'January', distinct_entities: 1000, non_entity_tokens: 322, total: 1322 },
      { month: 'February', distinct_entities: 1500, non_entity_tokens: 122, total: 1622 },
      { month: 'March', distinct_entities: 4300, non_entity_tokens: 700, total: 5000 },
      { month: 'April', distinct_entities: 1550, non_entity_tokens: 229, total: 1779 },
      { month: 'May', distinct_entities: 5560, non_entity_tokens: 124, total: 5684 },
      { month: 'June', distinct_entities: 1570, non_entity_tokens: 142, total: 1712 },
      { month: 'July', distinct_entities: 300, non_entity_tokens: 112, total: 412 },
      { month: 'August', distinct_entities: 1610, non_entity_tokens: 130, total: 1740 },
      { month: 'September', distinct_entities: 1900, non_entity_tokens: 222, total: 2122 },
      { month: 'October', distinct_entities: 500, non_entity_tokens: 166, total: 666 },
      { month: 'November', distinct_entities: 480, non_entity_tokens: 132, total: 612 },
      { month: 'December', distinct_entities: 980, non_entity_tokens: 202, total: 1182 },
    ];
  }

  // TODO: dataset for vault usage vertical bar chart (manage in serializer?)
  get monthlyUsage() {
    return [
      { month: 'January', distinct_entities: 1000, non_entity_tokens: 322, total: 1322 },
      { month: 'February', distinct_entities: 1500, non_entity_tokens: 122, total: 1622 },
      { month: 'March', distinct_entities: 4300, non_entity_tokens: 700, total: 5000 },
      { month: 'April', distinct_entities: 1550, non_entity_tokens: 229, total: 1779 },
      { month: 'May', distinct_entities: 5560, non_entity_tokens: 124, total: 5684 },
      { month: 'June', distinct_entities: 1570, non_entity_tokens: 142, total: 1712 },
      { month: 'July', distinct_entities: 300, non_entity_tokens: 112, total: 412 },
      { month: 'August', distinct_entities: 1610, non_entity_tokens: 130, total: 1740 },
      { month: 'September', distinct_entities: 1900, non_entity_tokens: 222, total: 2122 },
      { month: 'October', distinct_entities: 500, non_entity_tokens: 166, total: 666 },
      { month: 'November', distinct_entities: 480, non_entity_tokens: 132, total: 612 },
      { month: 'December', distinct_entities: 980, non_entity_tokens: 202, total: 1182 },
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

    namespaces.forEach(function (item) {
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

  async handleEndMonth() {
    // ARG TOOD consume param: endTime from calendar-widget fire off new network request
  }

  // Get the namespace by matching the path from the namespace list
  getNamespace(path) {
    return this.args.model.activity.byNamespace.find((ns) => {
      if (path === 'root') {
        return ns.namespace_path === '';
      }
      return ns.namespace_path === path;
    });
  }

  @action
  handleEditStartMonth() {
    // ARG TODO will need to handle the action when the click save on the modal and you are given a new start Month, via month and year in params,
    // if no endDate selected, default to 12 months range
    // then send range to queryData this.queryData(range)
  }

  // ARG TODO this might be a carry over from history, will need to confirm
  @action
  resetData() {
    this.barChartSelection = false;
    this.selectedNamespace = null;
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
  selectStartMonth(month) {
    this.startMonth = month;
  }

  @action
  selectStartYear(year) {
    this.startYear = year;
  }
}
