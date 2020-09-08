import TransformBase from './transform-edit-base';
import { inject as service } from '@ember/service';
import transform from '../adapters/transform';

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
    const testTransform = added[0];
    console.log({ testTransform });
    const backend = this.get('model.backend');

    this.store
      .queryRecord('transform', {
        backend,
        id: testTransform,
      })
      .then(function(transformation) {
        let roles = transformation.allowed_roles;
        roles.push(roleId);
        transformation.set('allowed_roles', roles);
        // transformation.allowed_roles = [roleId];
        // debugger;
        transformation.save({ backend }); // => PATCH to '/posts/1'
      });
  },

  handleRemovedTransformations(removed) {
    console.log({ removed });
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
