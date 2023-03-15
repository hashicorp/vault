/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { FormField, FormFieldGroups, FormFieldGroupOptions } from 'vault/app-types';

export default function _default(modelClass: Model, fieldGroups: FormFieldGroupOptions): FormFieldGroups;

export function expandAttributeMeta(modelClass: Model, attributeNames: Array<string>): Array<FormField>;
