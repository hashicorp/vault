import Evented from '@ember/object/evented';
import Service from '@ember/service';

import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';

let hasOwn = (obj, prop) => {
  return Object.prototype.hasOwnProperty.call(obj, prop);
};

export function extractRouteArgs(args) {
  args = args.slice();
  let possibleQueryParams = args[args.length - 1];

  let queryParams;
  if (possibleQueryParams && hasOwn(possibleQueryParams, 'queryParams')) {
    queryParams = args.pop().queryParams;
  } else {
    queryParams = {};
  }

  let routeName = args.shift();

  return { routeName, models: args, queryParams };
}
//https://github.com/emberjs/ember.js/blob/abf753a3d494830dc9e95b1337b3654b671b11be/packages/ember-routing/lib/utils.js#L210
export function shallowEqual(a, b) {
  let k;
  let aCount = 0;
  let bCount = 0;
  for (k in a) {
    if (hasOwn(a, k)) {
      if (a[k] !== b[k]) {
        return false;
      }
      aCount++;
    }
  }

  for (k in b) {
    if (hasOwn(b, k)) {
      bCount++;
    }
  }

  return aCount === bCount;
}

export default Service.extend(Evented, {
  init() {
    this._super(...arguments);

    this._router.on('routeWillChange', transition => {
      this.trigger('routeWillChange', transition);
    });

    this._router.on('routeDidChange', transition => {
      this.trigger('routeDidChange', transition);
    });
  },

  routing: service('-routing'),
  _router: alias('routing.router'),
  transitionTo() {
    let r = this._router;
    return r.transitionTo.call(r, ...arguments);
  },
  replaceWith() {
    let r = this._router;
    return r.replaceWith.call(r, ...arguments);
  },
  urlFor() {
    let r = this._router;
    return r.generate.call(r, ...arguments);
  },
  currentURL: alias('_router.currentURL'),
  currentRouteName: alias('_router.currentRouteName'),
  rootURL: alias('_router.rootURL'),
  location: alias('_router.location'),

  //adapted from:
  // https://github.com/emberjs/ember.js/blob/abf753a3d494830dc9e95b1337b3654b671b11be/packages/ember-routing/lib/services/router.js#L220
  isActive(...args) {
    let { routeName, models, queryParams } = extractRouteArgs(args);
    let routerMicrolib = this._router._routerMicrolib;

    if (!routerMicrolib.isActiveIntent(routeName, models, null)) {
      return false;
    }
    let hasQueryParams = Object.keys(queryParams).length > 0;

    if (hasQueryParams) {
      this._router._prepareQueryParams(routeName, models, queryParams, true /* fromRouterService */);
      return shallowEqual(queryParams, routerMicrolib.state.queryParams);
    }

    return true;
  },
});
