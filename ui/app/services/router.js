import Service from '@ember/service';
import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';

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
});
