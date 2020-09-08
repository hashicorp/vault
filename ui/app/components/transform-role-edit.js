import TransformBase from './transform-edit-base';
import { inject as service } from '@ember/service';

export default TransformBase.extend({
  store: service(),
  initialTransformations: [],

  init() {
    this._super(...arguments);
    this.set('initialTransformations', this.get('model.transformations'));
  },

  handleAddedTransformations(added, roleId) {
    if (!added) return;
    const backend = this.get('model.backend');

    // For each added transformation, add current role id to allowed_roles
    added.forEach(transformName => {
      this.store
        .queryRecord('transform', {
          backend,
          id: transformName,
        })
        .then(function(transformation) {
          let roles = transformation.allowed_roles;
          roles.push(roleId);
          transformation.setProperties({
            backend,
            allowed_roles: roles.uniq(),
          });
          transformation.save();
        });
      // TODO: Handle errors if no read or write access
    });
  },

  handleRemovedTransformations(removed, roleId) {
    if (!removed) return;
    const backend = this.get('model.backend');

    // For each removed transformation, remove current role id from allowed_roles
    removed.forEach(transformName => {
      this.store
        .queryRecord('transform', {
          backend,
          id: transformName,
        })
        .then(function(transformation) {
          let roles = transformation.allowed_roles;
          const index = roles.indexOf(roleId);
          if (index < 0) return;
          roles.removeAt(index);
          transformation.setProperties({
            backend,
            allowed_roles: roles.uniq(),
          });
          transformation.save();
        });
      // TODO: Handle errors if no read or write access
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges('save', () => {
        const roleId = this.get('model.id');
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
        this.handleAddedTransformations(addedTransformations, roleId);
        this.handleRemovedTransformations(removedTransformations, roleId);
      });
    },
  },
});
