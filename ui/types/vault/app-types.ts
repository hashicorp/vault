/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import type Owner from '@ember/owner';
import type { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

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

export type Validator =
  | 'presence'
  | 'length'
  | 'number'
  | 'containsWhiteSpace'
  | 'endsInSlash'
  | 'hasWhitespace'
  | 'isNonString'
  | 'isNot'
  | 'maxDuration';

export type ValidatorOption = {
  nullable?: boolean;
  min?: number;
  max?: number;
  value?: string | number | boolean;
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
export interface FormValidations {
  isValid: boolean;
  state: ValidationMap;
  invalidFormMessage: string;
}

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
  /** False when an enabled value is longer than Vault supports. */
  isValid: boolean;
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

export type WizardId = (typeof WIZARD_ID_MAP)[keyof typeof WIZARD_ID_MAP];
