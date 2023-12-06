/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Type that comes back from expandAttributeMeta
export interface FormField {
  name: string;
  type: string;
  options: unknown;
}

export interface FormFieldGroups {
  [key: string]: Array<FormField>;
}

export interface FormFieldGroupOptions {
  [key: string]: Array<string>;
}

export interface ModelValidation {
  isValid: boolean;
  state: {
    [key: string]: {
      isValid: boolean;
      errors: Array<string>;
    };
  };
  invalidFormMessage: string;
}
