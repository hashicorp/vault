/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import TransformBase, { addToList, removeFromList } from './transform-edit-base';
import { service } from '@ember/service';

export default TransformBase.extend({
  flashMessages: service(),
  store: service(),
  initialRoles: null,

  init() {
    this._super(...arguments);
    if (!this.model) return;
    this.set('initialRoles', this.model.allowed_roles);
  },

  updateOrCreateRole(role, transformationId, backend) {
    return this.store
      .queryRecord('transform/role', {
        backend,
        id: role.id,
      })
      .then((roleStore) => {
        let transformations = roleStore.transformations;
        if (role.action === 'ADD') {
          transformations = addToList(transformations, transformationId);
        } else if (role.action === 'REMOVE') {
          transformations = removeFromList(transformations, transformationId);
        }
        roleStore.setProperties({
          backend,
          transformations,
        });
        return roleStore.save().catch((e) => {
          return {
            errorStatus: e.httpStatus,
            ...role,
          };
        });
      })
      .catch((e) => {
        if (e.httpStatus !== 403 && role.action === 'ADD') {
          // If role doesn't yet exist, create it with this transformation attached
          var newRole = this.store.createRecord('transform/role', {
            id: role.id,
            name: role.id,
            transformations: [transformationId],
            backend,
          });
          return newRole.save().catch((e) => {
            return {
              errorStatus: e.httpStatus,
              ...role,
              action: 'CREATE',
            };
          });
        }

        return {
          ...role,
          errorStatus: e.httpStatus,
        };
      });
  },

  handleUpdateRoles(updateRoles, transformationId) {
    if (!updateRoles) return;
    const backend = this.model.backend;
    const promises = updateRoles.map((r) => this.updateOrCreateRole(r, transformationId, backend));

    Promise.all(promises).then((results) => {
      const hasError = results.find((role) => !!role.errorStatus);

      if (hasError) {
        let message =
          'The edits to this transformation were successful, but transformations for its roles was not edited due to a lack of permissions.';
        if (results.find((e) => !!e.errorStatus && e.errorStatus !== 403)) {
          // if the errors weren't all due to permissions show generic message
          // eg. trying to update a role with empty array as transformations
          message = `You've edited the allowed_roles for this transformation. However, the corresponding edits to some roles' transformations were not made`;
        }
        this.flashMessages.info(message, {
          sticky: true,
          priority: 300,
        });
      }
    });
  },

  isWildcard(role) {
    if (typeof role === 'string') {
      return role.indexOf('*') >= 0;
    }
    if (role && role.id) {
      return role.id.indexOf('*') >= 0;
    }
    return false;
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges('save', () => {
        const transformationId = this.model.id || this.model.name;
        const newModelRoles = this.model.allowed_roles || [];
        const initialRoles = this.initialRoles || [];

        const updateRoles = [...newModelRoles, ...initialRoles]
          .filter((r) => !this.isWildcard(r)) // CBS TODO: expand wildcards into included roles instead
          .map((role) => {
            if (initialRoles.indexOf(role) < 0) {
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
          .filter((r) => !!r);
        this.handleUpdateRoles(updateRoles, transformationId);
      });
    },
  },
});
