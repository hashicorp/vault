/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import type EmberDataModel from 'ember-data/model'; // eslint-disable-line ember/use-ember-data-rfc-395-imports
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

export type Model = Omit<EmberDataModel, 'isNew'> & {
  // override isNew which is a computed prop and ts will complain since it sees it as a function
  isNew: boolean;
};

export type WithFormFieldsModel = Model & {
  formFields: Array<FormField>;
  formFieldGroups: FormFieldGroups;
  allFields: Array<FormField>;
};

export type WithValidationsModel = Model & {
  validate(): ModelValidations;
};

export type WithFormFieldsAndValidationsModel = WithFormFieldsModel & {
  validate(): ModelValidations;
};

export interface Breadcrumb {
  label: string;
  route?: string;
  icon?: string;
  model?: string;
  models?: string[];
  linkExternal?: boolean;
}

export interface TtlEvent {
  enabled: boolean;
  seconds: number;
  timeString: string;
  goSafeTimeString: string;
}

export interface EngineOwner extends Owner {
  mountPoint: string;
}

export interface SearchSelectOption {
  name: string;
  id: string;
}

// Generic interfaces
export interface StringMap {
  [key: string]: string;
}
