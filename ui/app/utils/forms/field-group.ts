/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type FormField from './field';

// very simple util class that accepts a key and an array of form fields
// returns an object in the expected shape for use in the formFieldsGroup array of forms
export default class FormFieldGroup {
  [key: string]: FormField[]; // Add an index signature to allow dynamic property assignment

  constructor(groupName: string, fields: FormField[]) {
    this[groupName] = fields;
  }
}
