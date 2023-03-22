/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { debounce } from '@ember/runloop';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';

// TODO MOVE THESE TO THE ADDON
import utils from 'vault/lib/key-utils';
import keys from 'vault/lib/keycodes';
import { encodePath } from 'vault/utils/path-encoding-helpers';

/**
 * @module NavigateInput
 * `NavigateInput` components are used to filter list data.
 *
 * @example
 * ```js
 * <NavigateInput @filter={@roleFiltered} @placeholder="placeholder text" urls="{{hash list="vault.cluster.secrets.backend.kubernetes.roles"}}"/>
 * ```
 *
 * @param {String} filter=null  - The filtered string.
 * @param {String} [placeholder="Filter items"] - The message inside the input to indicate what the user should enter into the space.
 * @param {Object} [urls=null] - An object containing list=route url.
 * @param {Function} [filterFocusDidChange=null] - A function called when the focus changes.
 * @param {Function} [filterDidChange=null] - A function called when the filter string changes.
 * @param {Function} [filterMatchesKey=null] - A function used to match to a specific key, such as an Id.
 * @param {Function} [filterPartialMatch=null] - A function used to filter through a partial match. Such as "oo" of "root".
 * @param {String} [baseKey=""] - A string to transition by Id.
 * @param {Boolean} [shouldNavigateTree=false] - If true, navigate a larger tree, such as when you're navigating leases under access.
 * @param {String} [mode="secrets"] - Mode which plays into navigation type.
 * @param {String} [extraNavParams=""] - A string used in route transition when necessary.
 */

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
      const params = [routeFor('show', mode, this.args.urls), extraParams, this.keyForNav(val)].compact();
      this.transitionToRoute(...params);
    } else {
      if (mode === 'policies') {
        return;
      }
      const route = routeFor('create', mode, this.args.urls);
      if (baseKey) {
        this.transitionToRoute(route, this.keyForNav(baseKey), {
          queryParams: {
            initialKey: val,
          },
        });
      } else if (this.args.urls) {
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
    const firstPartialMatch = this.args.firstPartialMatch.id;
    if (!firstPartialMatch) {
      return;
    }
    event.preventDefault();
    this.args.filterDidChange(firstPartialMatch);
    this.filterUpdated(firstPartialMatch);
  }

  // as you type, navigates through the k/v tree
  filterUpdated(val) {
    const mode = this.mode;
    if (mode === 'policies' || !this.args.shouldNavigateTree) {
      this.filterUpdatedNoNav(val, mode);
      return;
    }
    // select the key to nav to, assumed to be a folder
    let key = val ? val.trim() : '';
    const isFolder = utils.keyIsFolder(key);

    if (!isFolder) {
      // nav to the closest parentKey (or the root)
      key = utils.parentKeyForKey(val) || '';
    }

    const pageFilter = val.replace(key, '');
    this.navigate(this.keyForNav(key), mode, pageFilter);
  }

  navigate(key, mode, pageFilter) {
    const route = routeFor(key ? 'list' : 'list-root', mode, this.args.urls);
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
    const key = val ? val.trim() : null;
    this.transitionToRoute(routeFor('list-root', mode, this.args.urls), {
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
    debounce(this, this.filterUpdated, filter.target.value, 200);
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
    const keyCode = event.keyCode;
    const val = event.target.value;
    if (keyCode === keys.ENTER) {
      this.onEnter(val);
    }
    if (keyCode === keys.ESC) {
      this.onEscape(val);
    }
  }
}
