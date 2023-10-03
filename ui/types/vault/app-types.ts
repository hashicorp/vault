/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import type EmberDataModel from '@ember-data/model';
import type Owner from '@ember/owner';

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

export interface Model extends Omit<EmberDataModel, 'isNew'> {
  // override isNew which is a computed prop and ts will complain since it sees it as a function
  isNew: boolean;
}

export interface WithFormFieldsModel extends Model {
  formFields: Array<FormField>;
  formFieldGroups: FormFieldGroups;
  allFields: Array<FormField>;
}

export interface WithValidationsModel extends Model {
  validate(): ModelValidations;
}

export interface WithFormFieldsAndValidationsModel extends WithFormFieldsModel, WithValidationsModel {}

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

export interface Breadcrumb {
  label: string;
  route?: string;
  linkExternal?: boolean;
}

export interface EngineOwner extends Owner {
  mountPoint: string;
}

// Generic interfaces
export interface StringMap {
  [key: string]: string;
}
