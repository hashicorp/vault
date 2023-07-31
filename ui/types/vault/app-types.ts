/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// Type that comes back from expandAttributeMeta
export interface FormField {
  name: string;
  type: string;
  options: AttributeOptions;
}

interface AttributeOptions {
  label: string;
  mapToBoolean: string;
  isOppositeValue: boolean;
}

export interface FormFieldGroups {
  [key: string]: Array<FormField>;
}

export interface FormFieldGroupOptions {
  [key: string]: Array<string>;
}

export interface ValidationMap {
  [key: string]: {
    isValid: boolean;
    errors: Array<string>;
  };
}
export interface ModelValidations {
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

export interface TtlEvent {
  enabled: boolean;
  seconds: number;
  timeString: string;
  goSafeTimeString: string;
}

// Generic interfaces
export interface StringMap {
  [key: string]: string;
}
