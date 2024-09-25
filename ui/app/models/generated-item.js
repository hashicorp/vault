/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import { tracked } from '@glimmer/tracking';

// This model is used for OpenApi-generated models in path-help service's getNewModel method
export default class GeneratedItemModel extends Model {
  allFields = [];

  @tracked _id;
  get mutableId() {
    return this._id || this.id;
  }
  set mutableId(value) {
    this._id = value;
  }

  get fieldGroups() {
    const groups = {
      default: [],
    };
    const fieldGroups = [];
    this.constructor.eachAttribute((name, attr) => {
      // if the attr comes in with a fieldGroup from OpenAPI,
      if (attr.options.fieldGroup) {
        if (groups[attr.options.fieldGroup]) {
          groups[attr.options.fieldGroup].push(attr);
        } else {
          groups[attr.options.fieldGroup] = [attr];
        }
      } else {
        // otherwise just add that attr to the default group
        groups.default.push(attr);
      }
    });
    for (const group in groups) {
      fieldGroups.push({ [group]: groups[group] });
    }
    return fieldGroups;
  }
}
