import Route from '@ember/routing/route';
// import EditBase from '../secret-edit';
import { inject as service } from '@ember/service';
// import { action } from '@ember/object';

export default class MetadataShow extends Route {
  @service store;
  // ARG TODO unloadmodelroute
  model(params) {
    console.log('model hook', params);
    let { secret } = params; // of dynamic route /*secret
    let parentModel = this.modelFor('vault.cluster.secrets.backend');
    return this.store
      .queryRecord('secret-v2', {
        backend: 'kv',
        id: secret,
      })
      .then(record => {
        console.log(record, 'RECORD!!!');
        return record;
      });
    // make an API request that uses the id
  }
}
// params because its a dynamic segment
// constructor() {
//   super(...arguments);
//   let test = this.paramsFor(this.routeName);
//   console.log(test, 'l;ajdf');
//   // this.fetchCapabilities();
// }

// controller.setProperties({
//   model: model.lease,
//   capabilities: model.capabilities,
//   baseKey: { id: secret },
// });
// @service store;
// // …
// @action
// async visitUserProfile(id) {
//   this.store.findRecord('user', id).then(function (user) {
//     // Success callback
//     this.transitionTo('user.profile', user);
//   }).catch(function () {
//     // Error callback
//     this.transitionTo('not-found', 404);
//   }
// }
