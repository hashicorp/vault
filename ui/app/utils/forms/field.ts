/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface FieldOptions {
  label?: string;
  subText?: string;
  fieldValue?: string;
  editType?: string;
  defaultValue?: unknown;
  possibleValues?: unknown[];
  allowWhiteSpace?: boolean;
  isSingleRow?: boolean;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  editDisabled?: boolean;
  sensitive?: boolean;
  noCopy?: boolean;
  docLink?: string;
  helpText?: string;
  helperTextDisabled?: string;
  helperTextEnabled?: string;
  placeholder?: string;
}

export default class FormField {
  name = '';
  type: string | undefined;
  options: FieldOptions = {};

  constructor(key: string, type?: string, options: FieldOptions = {}) {
    this.name = key;
    this.type = type;
    this.options = {
      ...options,
      fieldValue: options.fieldValue || key,
    };
  }
}
