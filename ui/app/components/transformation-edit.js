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

  async updateOrCreateRole(role, transformationId, backend) {
    const roleRecord = await this.store
      .queryRecord('transform/role', {
        backend,
        id: role.id,
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
    // if an error occurs while querying the role, exit function and return the error
    if (roleRecord.errorStatus) return roleRecord;
    // otherwise update the role with the transformation and save
    let transformations = roleRecord.transformations;
    if (role.action === 'ADD') {
      transformations = addToList(transformations, transformationId);
    } else if (role.action === 'REMOVE') {
      transformations = removeFromList(transformations, transformationId);
    }
    roleRecord.setProperties({
      backend,
      transformations,
    });
    return roleRecord.save().catch((e) => {
      return {
        errorStatus: e.httpStatus,
        ...role,
      };
    });
  },

  handleUpdateRoles(updateRoles, transformationId) {
    if (!updateRoles) return;
    const { backend } = this.model;
    updateRoles.forEach(async (record) => {
      // For each role that needs to be updated, update the role with the transformation.
      const updateOrCreateResponse = await this.updateOrCreateRole(record, transformationId, backend);
      // If an error was returned, check error type and show a message.
      const errorStatus = updateOrCreateResponse?.errorStatus;
      let message;
      if (errorStatus == 403) {
        message = `The edits to this transformation were successful, but transformations for the role ${record.id} were not edited due to a lack of permissions.`;
      } else if (errorStatus) {
        message = `You've edited the allowed_roles for this transformation. However, there was a problem updating the role: ${record.id}.`;
      }
      this.flashMessages.info(message, {
        sticky: true,
        priority: 300,
      });
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
