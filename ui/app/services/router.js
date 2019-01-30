import Service from '@ember/service';
import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
export function extractRouteArgs(args) {
  args = args.slice();
  let possibleQueryParams = args[args.length - 1];

  let queryParams;
  if (possibleQueryParams && possibleQueryParams.hasOwnProperty('queryParams')) {
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
    if (a.hasOwnProperty(k)) {
      if (a[k] !== b[k]) {
        return false;
      }
      aCount++;
    }
  }

  for (k in b) {
    if (b.hasOwnProperty(k)) {
      bCount++;
    }
  }

  return aCount === bCount;
}

export default Service.extend({
  routing: service('-routing'),
  router: alias('routing.router'),
  transitionTo() {
    let r = this.router;
    return r.transitionTo.call(r, ...arguments);
  },
  replaceWith() {
    let r = this.router;
    return r.replaceWith.call(r, ...arguments);
  },
  urlFor() {
    let r = this.router;
    return r.generate.call(r, ...arguments);
  },
  currentURL: alias('router.currentURL'),
  currentRouteName: alias('router.currentRouteName'),
  rootURL: alias('router.rootURL'),
  location: alias('router.location'),

  //adapted from:
  // https://github.com/emberjs/ember.js/blob/abf753a3d494830dc9e95b1337b3654b671b11be/packages/ember-routing/lib/services/router.js#L220
  isActive(...args) {
    let { routeName, models, queryParams } = extractRouteArgs(args);
    let routerMicrolib = this.router._routerMicrolib;

    if (!routerMicrolib.isActiveIntent(routeName, models, null)) {
      return false;
    }
    let hasQueryParams = Object.keys(queryParams).length > 0;

    if (hasQueryParams) {
      this.router._prepareQueryParams(routeName, models, queryParams, true /* fromRouterService */);
      return shallowEqual(queryParams, routerMicrolib.state.queryParams);
    }

    return true;
  },
});
