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

export type Validator =
  | 'presence'
  | 'length'
  | 'number'
  | 'containsWhiteSpace'
  | 'endsInSlash'
  | 'hasWhitespace'
  | 'isNonString';

export type ValidatorOption =
  | {
      nullable?: boolean;
    }
  | {
      nullable?: boolean;
      min?: number;
      max?: number;
    };

export type Validation =
  | {
      message: string | ((data: unknown) => string);
      type: Validator;
      options?: ValidatorOption;
      level?: 'warn';
      validator?: never;
    }
  | {
      message: string | ((data: unknown) => string);
      type?: never;
      options?: ValidatorOption;
      level?: 'warn';
      validator(data: unknown, options?: unknown): boolean;
    };
export interface Validations {
  [key: string]: Validation[];
}
export interface ValidationMap {
  [key: string]: {
    isValid: boolean;
    errors: Array<string>;
    warnings: Array<string>;
  };
}
export interface ModelValidations {
  isValid: boolean;
  state: ValidationMap;
  invalidFormMessage: string;
}
// TODO: [Ember Data] - ModelValidations can be renamed to FormValidations and this can be removed
export type FormValidations = ModelValidations;

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

// capabilities
export interface Capabilities {
  canCreate: boolean;
  canDelete: boolean;
  canList: boolean;
  canPatch: boolean;
  canRead: boolean;
  canSudo: boolean;
  canUpdate: boolean;
}

export interface CapabilitiesMap {
  [key: string]: Capabilities;
}

export type CapabilityTypes =
  | 'root'
  | 'sudo'
  | 'deny'
  | 'create'
  | 'read'
  | 'update'
  | 'delete'
  | 'list'
  | 'patch';
export interface CapabilitiesData {
  [key: string]: CapabilityTypes[];
}

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
