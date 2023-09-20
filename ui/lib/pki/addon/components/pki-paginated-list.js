import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class PkiPaginatedListComponent extends Component {
  @service router;
  @tracked pageFilter = '';

  constructor() {
    super(...arguments);
    this.pageFilter = this.args.pageFilter;
  }

  @action
  onFilterChange(filterText) {
    const pageFilter = !filterText ? undefined : filterText;
    const fullRouteName = `vault.cluster.secrets.backend.pki.${this.args.listRoute}`;
    this.router.transitionTo(fullRouteName, { queryParams: { pageFilter, currentPage: 1 } });
  }

  get paginationQueryParams() {
    return (page) => ({
      currentPage: page,
    });
  }
}
