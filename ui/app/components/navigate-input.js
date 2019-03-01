import { schedule, debounce } from '@ember/runloop';
import { observer } from '@ember/object';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import utils from 'vault/lib/key-utils';
import keys from 'vault/lib/keycodes';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import { encodePath } from 'vault/utils/path-encoding-helpers';

const routeFor = function(type, mode) {
  const MODES = {
    secrets: 'vault.cluster.secrets.backend',
    'secrets-cert': 'vault.cluster.secrets.backend',
    'policy-show': 'vault.cluster.policy',
    'policy-list': 'vault.cluster.policies',
    leases: 'vault.cluster.access.leases',
  };
  let useSuffix = true;
  const typeVal = mode === 'secrets' || mode === 'leases' ? type : type.replace('-root', '');
  const modeKey = mode + '-' + typeVal;
  const modeVal = MODES[modeKey] || MODES[mode];
  if (modeKey === 'policy-list') {
    useSuffix = false;
  }

  return useSuffix ? modeVal + '.' + typeVal : modeVal;
};

export default Component.extend(FocusOnInsertMixin, {
  router: service(),

  classNames: ['navigate-filter'],

  // these get passed in from the outside
  // actions that get passed in
  filterFocusDidChange: null,
  filterDidChange: null,
  mode: 'secrets',
  shouldNavigateTree: false,
  extraNavParams: null,

  baseKey: null,
  filter: null,
  filterMatchesKey: null,
  firstPartialMatch: null,

  transitionToRoute(...args) {
    let params = args.map((param, index) => {
      if (index === 0 || typeof param !== 'string') {
        return param;
      }
      return encodePath(param);
    });

    this.get('router').transitionTo(...params);
  },

  shouldFocus: false,

  focusFilter: observer('filter', function() {
    if (!this.get('filter')) return;
    schedule('afterRender', this, 'forceFocus');
  }).on('didInsertElement'),

  keyForNav(key) {
    if (this.get('mode') !== 'secrets-cert') {
      return key;
    }
    return `cert/${key}`;
  },
  onEnter: function(val) {
    let { baseKey, mode } = this;
    let extraParams = this.get('extraNavParams');
    if (mode.startsWith('secrets') && (!val || val === baseKey)) {
      return;
    }
    if (this.get('filterMatchesKey') && !utils.keyIsFolder(val)) {
      let params = [routeFor('show', mode), extraParams, this.keyForNav(val)].compact();
      this.transitionToRoute(...params);
    } else {
      if (mode === 'policies') {
        return;
      }
      let route = routeFor('create', mode);
      if (baseKey) {
        this.transitionToRoute(route, this.keyForNav(baseKey), {
          queryParams: {
            initialKey: val,
          },
        });
      } else {
        this.transitionToRoute(route + '-root', {
          queryParams: {
            initialKey: this.keyForNav(val),
          },
        });
      }
    }
  },

  // pop to the nearest parentKey or to the root
  onEscape: function(val) {
    var key = utils.parentKeyForKey(val) || '';
    this.get('filterDidChange')(key);
    this.filterUpdated(key);
  },

  onTab: function(event) {
    var firstPartialMatch = this.get('firstPartialMatch.id');
    if (!firstPartialMatch) {
      return;
    }
    event.preventDefault();
    this.get('filterDidChange')(firstPartialMatch);
    this.filterUpdated(firstPartialMatch);
  },

  // as you type, navigates through the k/v tree
  filterUpdated: function(val) {
    var mode = this.get('mode');
    if (mode === 'policies' || !this.get('shouldNavigateTree')) {
      this.filterUpdatedNoNav(val, mode);
      return;
    }
    // select the key to nav to, assumed to be a folder
    var key = val ? val.trim() : '';
    var isFolder = utils.keyIsFolder(key);

    if (!isFolder) {
      // nav to the closest parentKey (or the root)
      key = utils.parentKeyForKey(val) || '';
    }

    const pageFilter = val.replace(key, '');
    this.navigate(this.keyForNav(key), mode, pageFilter);
  },

  navigate(key, mode, pageFilter) {
    const route = routeFor(key ? 'list' : 'list-root', mode);
    let args = [route];
    if (key) {
      args.push(key);
    }
    if (pageFilter && !utils.keyIsFolder(pageFilter)) {
      args.push({
        queryParams: {
          page: 1,
          pageFilter,
        },
      });
    } else {
      args.push({
        queryParams: {
          page: 1,
          pageFilter: null,
        },
      });
    }
    this.transitionToRoute(...args);
  },

  filterUpdatedNoNav: function(val, mode) {
    var key = val ? val.trim() : null;
    this.transitionToRoute(routeFor('list-root', mode), {
      queryParams: {
        pageFilter: key,
        page: 1,
      },
    });
  },

  actions: {
    handleInput: function(filter) {
      this.get('filterDidChange')(filter);
      debounce(this, 'filterUpdated', filter, 200);
    },

    setFilterFocused: function(isFocused) {
      this.get('filterFocusDidChange')(isFocused);
    },

    handleKeyPress: function(event) {
      if (event.keyCode === keys.TAB) {
        this.onTab(event);
      }
    },

    handleKeyUp: function(event) {
      var keyCode = event.keyCode;
      let val = event.target.value;
      if (keyCode === keys.ENTER) {
        this.onEnter(val);
      }
      if (keyCode === keys.ESC) {
        this.onEscape(val);
      }
    },
  },
});
