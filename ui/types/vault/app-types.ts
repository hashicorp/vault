/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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

export interface Breadcrumb {
  label: string;
  route?: string;
  linkExternal?: boolean;
}
