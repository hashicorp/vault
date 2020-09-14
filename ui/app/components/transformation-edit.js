import TransformBase from './transform-edit-base';

export default TransformBase.extend({
  initialRoles: null,

  init() {
    this._super(...arguments);
    this.set('initialRoles', this.get('model.allowed_roles'));
  },

  handleUpdateRoles(updateRoles, transformationId, type = 'update') {
    console.log({ updateRoles });
    if (!updateRoles) return;
    const backend = this.get('model.backend');
    const promises = updateRoles.map(r => {
      return this.store
        .queryRecord('transform/role', {
          backend,
          id: r.id,
        })
        .then(role => {
          console.log('UPDATE ROLES -- ROLE FETCHEd', transformationId);
          let transformations = role.transformations;
          if (r.action === 'ADD') {
            transformations = this.addToList(transformations, transformationId);
          } else if (r.action === 'REMOVE') {
            transformations = this.removeFromList(transformations, transformationId);
          }
          role.setProperties({
            backend,
            transformations,
          });
          console.log(`Saving role ${r.id} with transforms:`, transformations);
          return role.save().catch(e => {
            console.log('ERror on save', e.message);
            return {
              errorStatus: e.httpStatus,
              ...r,
            };
          });
        })
        .catch(e => {
          if (e.httpStatus !== 403 && r.action === 'ADD') {
            // If role doesn't yet exist, create it with this transformation attached
            var newRole = this.store.createRecord('transform/role', {
              id: r.id,
              name: r.id,
              transformations: [transformationId],
              backend,
            });
            return newRole.save().catch(e => {
              return {
                errorStatus: e.httpStatus,
                ...r,
                action: 'CREATE',
              };
            });
          }

          // TODO: create role if httpStatus == 403 and r.action == 'ADD'
          return {
            ...r,
            errorStatus: e.httpStatus,
          };
        });
    });

    Promise.all(promises).then(results => {
      let hasError = results.find(r => !!r.errorStatus);

      if (hasError) {
        let message =
          'The edits to this transformation were successful, but transformations for its roles was not edited due to a lack of permissions.';
        if (hasError.find(e => e.errorStatus !== 403)) {
          // if the errors weren't all due to permissions show generic message
          // eg. trying to update a role with empty array as transformations
          message =
            'The edits to this transformation were successful, but some updates to roles resulted in an error.';
        }
        this.get('flashMessages').stickyInfo(message);
      }
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges('save', () => {
        const transformationId = this.get('model.id');
        const newModelRoles = this.get('model.allowed_roles');
        if (!this.initialRoles) {
          this.handleUpdateRoles(
            newModelRoles.map(r => ({
              action: 'ADD',
              id: r,
            })),
            transformationId,
            type
          );
        }
        // TODO: deal with wildcards
        const updateRoles = [...newModelRoles, ...this.initialRoles]
          .map(role => {
            if (this.initialRoles.indexOf(role) < 0) {
              return {
                id: role,
                action: 'ADD',
              };
            }
            if (newModelRoles.indexOf(role) < 0) {
              return {
                id: role,
                action: 'REMOVE',
              };
            }
            return null;
          })
          .filter(r => !!r);
        this.handleUpdateRoles(updateRoles, transformationId);
      });
    },
  },
});
