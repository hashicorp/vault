/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

// see: https://helios.hashicorp.design/components/form/key-value-inputs?tab=code#input-types
export interface KeyValueField {
  name: string;
  label?: string;
  type?: 'text' | 'textarea' | 'select' | 'masked' | 'file';
  placeholder?: string;
  possibleValues?: unknown[];
  // prepends a "Select one" placeholder option when type is "select"
  noDefault?: boolean;
  width?: string;
  isRequired?: boolean;
  isOptional?: boolean;
  helpText?: string;
}

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
  keyInputType?: string;
  valueInputType?: string;
  keyPossibleValues?: unknown[];
  valuePossibleValues?: unknown[];
  // fields rendered in each row for editType "keyValueInputs"
  keyValueFields?: KeyValueField[];
  isRequired?: boolean;
  isOptional?: boolean;
  maxRows?: number;
  maxRowsText?: string;
  addRowButtonText?: string;
  editDisabled?: boolean;
  sensitive?: boolean;
  noCopy?: boolean;
  docLink?: string;
  helpText?: string;
  helperTextDisabled?: string;
  helperTextEnabled?: string;
  placeholder?: string;
  noDefault?: boolean;
  isSectionHeader?: boolean;
  hideToggle?: boolean;
  labelDisabled?: string;
  mapToBoolean?: string;
  isOppositeValue?: boolean;
  defaultSubText?: string;
  defaultShown?: string;
  example?: string;
  mode?: string;
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
