/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import { FormField, FormFieldGroups, FormFieldGroupOptions } from 'vault/app-types';

export default function _default(modelClass: Model, fieldGroups: FormFieldGroupOptions): FormFieldGroups;

export function expandAttributeMeta(modelClass: Model, attributeNames: Array<string>): Array<FormField>;
