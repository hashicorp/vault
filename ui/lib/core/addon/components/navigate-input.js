import { debounce } from '@ember/runloop';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';

//TODO MOVE THESE TO THE ADDON
// ARG TODO RETURN
import utils from 'vault/lib/key-utils';
import keys from 'vault/lib/keycodes';
import { encodePath } from 'vault/utils/path-encoding-helpers';

const routeFor = function (type, mode, urls) {
  const MODES = {
    secrets: 'vault.cluster.secrets.backend',
    'secrets-cert': 'vault.cluster.secrets.backend',
    'policy-show': 'vault.cluster.policy',
    'policy-list': 'vault.cluster.policies',
    leases: 'vault.cluster.access.leases',
  };
  // urls object should have create, list, show keys
  // so we'll return that here
  if (urls) {
    return urls[type.replace('-root', '')];
  }
  let useSuffix = true;
  const typeVal = mode === 'secrets' || mode === 'leases' ? type : type.replace('-root', '');
  const modeKey = mode + '-' + typeVal;
  const modeVal = MODES[modeKey] || MODES[mode];
  if (modeKey === 'policy-list') {
    useSuffix = false;
  }

  return useSuffix ? modeVal + '.' + typeVal : modeVal;
};

export default class NavigateInput extends Component {
  @service router;

  get focusFilter() {
    return this.args.filter ? true : false;
  }

  get mode() {
    return this.args.mode || 'secrets';
  }

  transitionToRoute(...args) {
    const params = args.map((param, index) => {
      if (index === 0 || typeof param !== 'string') {
        return param;
      }
      return encodePath(param);
    });

    this.router.transitionTo(...params);
  }

  keyForNav(key) {
    if (this.mode !== 'secrets-cert') {
      return key;
    }
    return `cert/${key}`;
  }

  onEnter(val) {
    const mode = this.mode;
    const baseKey = this.args.baseKey;
    const extraParams = this.args.extraNavParams;
    if (mode.startsWith('secrets') && (!val || val === baseKey)) {
      return;
    }
    if (this.args.filterMatchesKey && !utils.keyIsFolder(val)) {
      const params = [routeFor('show', mode, this.urls), extraParams, this.keyForNav(val)].compact();
      this.transitionToRoute(...params);
    } else {
      if (mode === 'policies') {
        return;
      }
      const route = routeFor('create', mode, this.urls);
      if (baseKey) {
        this.transitionToRoute(route, this.keyForNav(baseKey), {
          queryParams: {
            initialKey: val,
          },
        });
      } else if (this.urls) {
        this.transitionToRoute(route, {
          queryParams: {
            initialKey: this.keyForNav(val),
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
  }

  // pop to the nearest parentKey or to the root
  onEscape(val) {
    const key = utils.parentKeyForKey(val) || '';
    this.args.filterDidChange(key);
    this.filterUpdated(key);
  }

  onTab(event) {
    var firstPartialMatch = this.args.firstPartialMatch.id;
    if (!firstPartialMatch) {
      return;
    }
    event.preventDefault();
    this.args.filterDidChange(firstPartialMatch);
    this.filterUpdated(firstPartialMatch);
  }

  // as you type, navigates through the k/v tree
  filterUpdated(val) {
    var mode = this.mode;
    if (mode === 'policies' || !this.args.shouldNavigateTree) {
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
  }

  navigate(key, mode, pageFilter) {
    const route = routeFor(key ? 'list' : 'list-root', mode, this.urls);
    const args = [route];
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
  }

  filterUpdatedNoNav(val, mode) {
    var key = val ? val.trim() : null;
    this.transitionToRoute(routeFor('list-root', mode, this.urls), {
      queryParams: {
        pageFilter: key,
        page: 1,
      },
    });
  }

  @action
  handleInput(filter) {
    if (this.args.filterDidChange) {
      this.args.filterDidChange(filter.target.value);
    }
    debounce(this, 'filterUpdated', filter.target.value, 200);
  }
  @action
  setFilterFocused(isFocused) {
    if (this.args.filterFocusDidChange) {
      this.args.filterFocusDidChange(isFocused);
    }
  }
  @action
  handleKeyPress(event) {
    if (event.keyCode === keys.TAB) {
      this.onTab(event);
    }
  }
  @action
  handleKeyUp(event) {
    var keyCode = event.keyCode;
    const val = event.target.value;
    if (keyCode === keys.ENTER) {
      this.onEnter(val);
    }
    if (keyCode === keys.ESC) {
      this.onEscape(val);
    }
  }
}
