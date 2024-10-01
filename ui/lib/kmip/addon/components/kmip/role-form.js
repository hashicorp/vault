/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { removeManyFromArray } from 'vault/helpers/remove-from-array';
import { operationFieldsWithoutSpecial, tlsFields } from 'vault/utils/model-helpers/kmip-role-fields';

export default class KmipRoleFormComponent extends Component {
  @service flashMessages;
  @service store;

  // Actual attribute fields
  get tlsFormFields() {
    return tlsFields().map((attr) => this.args.model.allByKey[attr]);
  }
  get operationFormGroups() {
    const objects = [
      'operationCreate',
      'operationActivate',
      'operationGet',
      'operationLocate',
      'operationRekey',
      'operationRevoke',
      'operationDestroy',
    ];
    const attributes = ['operationAddAttribute', 'operationGetAttributes'];
    const server = ['operationDiscoverVersions'];
    const others = removeManyFromArray(operationFieldsWithoutSpecial(this.args.model.editableFields), [
      ...objects,
      ...attributes,
      ...server,
    ]);
    const groups = [
      { name: 'Managed Cryptographic Objects', fields: objects },
      { name: 'Object Attributes', fields: attributes },
      { name: 'Server', fields: server },
    ];
    if (others.length) {
      groups.push({
        name: 'Other',
        fields: others,
      });
    }
    // expand field names to attributes
    return groups.map((group) => ({
      ...group,
      fields: group.fields.map((attr) => this.args.model.allByKey[attr]),
    }));
  }

  placeholderOrModel = (model, attrName) => {
    return model.operationAll ? { [attrName]: true } : model;
  };

  preSave() {
    const opFieldsWithoutSpecial = operationFieldsWithoutSpecial(this.args.model.editableFields);
    // if we have operationAll or operationNone, we want to clear
    // out the others so that display shows the right data
    if (this.args.model.operationAll || this.args.model.operationNone) {
      opFieldsWithoutSpecial.forEach((field) => (this.args.model[field] = null));
    }
    // set operationNone if user unchecks 'operationAll' instead of toggling the 'operationNone' input
    // doing here instead of on the 'operationNone' input because a user might deselect all, then reselect some options
    // and immediately setting operationNone will hide all of the checkboxes in the UI
    this.args.model.operationNone =
      opFieldsWithoutSpecial.every((attr) => this.args.model[attr] !== true) && !this.args.model.operationAll;
    return this.args.model;
  }

  @action toggleOperationSpecial(evt) {
    const { checked } = evt.target;
    this.args.model.operationNone = !checked;
    this.args.model.operationAll = checked;
  }

  save = task(async (evt) => {
    evt.preventDefault();
    const model = this.preSave();
    try {
      await model.save();
      this.flashMessages.success(`Saved role ${model.role}`);
    } catch (err) {
      // err will display via model state
      // AdapterErrors are handled by the error-message component
      if (err instanceof AdapterError === false) {
        throw err;
      }
      return;
    }
    this.args.onSave();
  });

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && this.args?.model?.isDirty) {
      this.args.model.rollbackAttributes();
    }
    super.willDestroy();
  }
}
