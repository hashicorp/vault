import { computed } from '@ember/object';
import Component from '@ember/component';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default Component.extend({
  onLinkClick() {},
  tagName: '',
  // so that ember-test-selectors doesn't log a warning
  supportsDataTestProperties: true,
  mode: 'list',

  secret: null,
  queryParams: null,
  ariaLabel: null,

  link: computed('mode', 'secret', function () {
    const route = `vault.cluster.secrets.backend.${this.mode}`;
    if ((this.mode !== 'versions' && !this.secret) || this.secret === ' ') {
      return { route: `${route}-root`, models: [] };
    } else {
      return { route, models: [encodePath(this.secret)] };
    }
  }),
  query: computed('queryParams', function () {
    const qp = this.queryParams || {};
    return qp.isQueryParams ? qp.values : qp;
  }),
});
