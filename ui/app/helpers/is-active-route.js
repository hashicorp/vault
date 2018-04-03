import Ember from 'ember';

const { Helper, inject, observer } = Ember;

const exact = (a, b) => a === b;
const startsWith = (a, b) => a.indexOf(b) === 0;

export default Helper.extend({
  routing: inject.service('-routing'),

  onRouteChange: observer('routing.router.currentURL', 'routing.router.currentRouteName', function() {
    this.recompute();
  }),

  compute([routeName, model], { isExact }) {
    const router = this.get('routing.router');
    const currentRoute = router.get('currentRouteName');
    let currentURL = router.get('currentURL');
    // if we have any query params we want to discard them
    currentURL = currentURL.split('?')[0];
    const comparator = isExact ? exact : startsWith;
    if (!currentRoute) {
      return false;
    }
    if (Ember.isArray(routeName)) {
      return routeName.some(name => comparator(currentRoute, name));
    } else if (model) {
      // slice off the rootURL from the generated route
      return comparator(currentURL, router.generate(routeName, model).slice(router.rootURL.length - 1));
    } else {
      return comparator(currentRoute, routeName);
    }
  },
});
