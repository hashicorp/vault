import TransformBase from './transform-edit-base';
import { inject as service } from '@ember/service';
// import transform from '../adapters/transform';

export default TransformBase.extend({
  store: service(),
  initialTransformations: [],

  init() {
    this._super(...arguments);
    // CBS TODO: set initialTransformations for compare on createOrUpdate
    this.set('initialTransformations', this.get('model.transformations'));
  },

  handleAddedTransformations(added, roleId) {
    if (!added) return;

    // get just one for now
    // const testTransform = added[0];
    // console.log({ testTransform });
    const backend = this.get('model.backend');
    // console.log(this.get('model'));
    // For each added transformation, add current role id to allowed_roles
    added.forEach(transformName => {
      this.store
        .queryRecord('transform', {
          backend,
          id: transformName,
        })
        .then(function(transformation) {
          transformation.set('backend', backend);
          console.log('get successful', transformation);
          let roles = transformation.allowed_roles;
          roles.push(roleId);
          transformation.set('allowed_roles', roles);
          // transformation.allowed_roles = [roleId];
          transformation.save();
        });
    });
    // this.store
    //   .queryRecord('transform', {
    //     backend,
    //     id: testTransform,
    //   })
    //   .then(function(transformation) {
    //     transformation.set('backend', backend);
    //     console.log('get successful', transformation);
    //     let roles = transformation.allowed_roles;
    //     roles.push(roleId);
    //     transformation.set('allowed_roles', roles);
    //     // transformation.allowed_roles = [roleId];
    //     transformation.save();
    //   });
  },

  handleRemovedTransformations(removed) {
    if (!removed) return;
    const backend = this.get('model.backend');

    // For each added transformation, remove current role id from allowed_roles
    removed.forEach(transformName => {
      this.store
        .queryRecord('transform', {
          backend,
          id: transformName,
        })
        .then(function(transformation) {
          transformation.set('backend', backend);
          console.log('get successful', transformation);
          let roles = transformation.allowed_roles;
          roles.push(roleId);
          transformation.set('allowed_roles', roles);
          // transformation.allowed_roles = [roleId];
          transformation.save();
        });
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges('save', () => {
        console.log('-> -> -> This is where we compare initialTransformations to appliedTransformations');
        const newModelTransformations = this.get('model.transformations');
        if (!this.initialTransformations) {
          this.handleAddedTransformations(newModelTransformations);
          return;
        }
        const addedTransformations = newModelTransformations.filter(
          t => this.initialTransformations.indexOf(t) < 0
        );
        const removedTransformations = this.initialTransformations.filter(
          t => newModelTransformations.indexOf(t) < 0
        );
        console.log(this.initialTransformations, this.get('model.transformations'));
        this.handleAddedTransformations(addedTransformations, this.get('model.id'));
        this.handleRemovedTransformations(removedTransformations);
      });
    },
  },
});
