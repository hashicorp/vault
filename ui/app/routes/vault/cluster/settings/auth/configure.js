import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    const { method } = this.paramsFor(this.routeName);
    return this.store.findAll('auth-method').then(() => {
      return this.store.peekRecord('auth-method', method);
    });
  },
});
