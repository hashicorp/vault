/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type ApiService from 'vault/services/api';
import type { ValidationRule } from './form-validator';
import CONFIG_REGISTRY from './generated/index';

export type FormConfigKey = keyof typeof CONFIG_REGISTRY;

/**
 * Form element types matching HDS (HashiCorp Design System) components.
 */
export type FormElement =
  | 'TextInput'
  | 'TextArea'
  | 'Select'
  | 'Toggle'
  | 'Checkbox'
  | 'Radio'
  | 'RadioCard'
  | 'MaskedInput';

/**
 * HDS-supported text input variants for TextInput fields.
 */
export type TextInputType =
  | 'text'
  | 'email'
  | 'password'
  | 'url'
  | 'search'
  | 'date'
  | 'time'
  | 'datetime-local'
  | 'month'
  | 'week'
  | 'tel';

/**
 * Union type for all possible field values.
 */
export type FieldValue = string | number | boolean | string[] | null | undefined;

/**
 * Visibility predicate evaluated against the form payload.
 */
export type VisibilityRule<Payload extends object = object> = boolean | ((payload: Payload) => boolean);

/**
 * Option for Select, Radio, and RadioCard fields
 */
export interface FieldOption {
  label: string;
  value: string | number | boolean;
  /** Optional description for RadioCard options */
  description?: string;
}

/**
 * Form field definition with common properties.
 */
export type FormField = {
  /** Field name supporting dotted-path notation for nested properties */
  name: string;
  /** Form element type */
  type: FormElement;
  /** Display label for the field */
  label: string;
  /** Optional helper text shown below the field */
  helperText?: string;
  /** Optional HDS TextInput variant when type is TextInput */
  inputType?: TextInputType;
  /** Optional placeholder text */
  placeholder?: string;
  /** Default value for the field */
  defaultValue?: FieldValue;
  /** Validation rules for this field */
  validations?: ValidationRule[];
  /** Options for Select, Radio, and RadioCard fields */
  options?: FieldOption[];
  /** Whether the field is required */
  isRequired?: boolean;
  /** Whether the field is disabled */
  isDisabled?: boolean;
  /** Optional conditional visibility for this field */
  isVisible?: VisibilityRule;
};

/**
 * Form section grouping related fields together
 */
export interface FormSection {
  /** Section identifier (used as key) */
  name: string;
  /** Optional display title for the section */
  title?: string;
  /** Optional description for the section */
  description?: string;
  /** Optional conditional visibility for this section */
  isVisible?: VisibilityRule;
  /** Fields belonging to this section */
  fields: FormField[];
}

/**
 * Wizard state tracking data for each completed step.
 * Stores only data - execution state is derived from ember-concurrency task properties.
 * Keyed by step name, contains the submitted payload, API response, and error message.
 */
export interface WizardStepState {
  /** The payload that was submitted for this step */
  payload: any;
  /** The API response received for this step (present = step succeeded) */
  response: any;
  /** Error message if the step failed */
  error?: string;
}

/**
 * Wizard state accumulator passed to payload resolver functions.
 * Maps step names to their completed state.
 */
export type WizardState = {
  [stepName: string]: WizardStepState;
};

/**
 * A single step in a multi-step wizard.
 * References a FormConfig and provides step-specific metadata.
 */
export interface WizardStep {
  /** Unique identifier for this step (used to key wizard state) */
  name: string;
  /** Display title for the step */
  title: string;
  /** Optional heading for the step in the panel */
  heading?: string;
  /** Optional description for the step in the panel */
  description?: string;
  /** The form configuration for this step */
  formConfig: FormConfig<any, any>;
}

/**
 * Multi-step wizard configuration.
 * Defines a sequential flow of forms with cross-step data sharing.
 */
export interface WizardConfig {
  /** Display title for the wizard */
  title: string;
  /** Optional description for the wizard */
  description?: string;
  /** Optional flag to indicate if the wizard has a final apply changes step */
  applyChanges?: boolean;
  /** Sequential steps in the wizard */
  steps: WizardStep[];
}

export interface FormConfig<Request extends object = object, Response = unknown> {
  /** Unique identifier for the form, typically matching the API method name */
  name: string;
  /** API endpoint path associated with this form -- useful for generating CURL request snippets */
  path: string;
  /** Title or description for the form */
  title?: string;
  description?: string;
  /**
   * Initial payload structure matching the API request shape.
   */
  payload: Request;
  /**
   * Submit handler that receives the API service and typed payload,
   * returning the typed API response
   */
  submit: (api: ApiService, payload: Request) => Promise<Response>;
  /**
   * Optional callback invoked after successful submission.
   * Receives the API response for post-submission handling (e.g., redirects, notifications).
   */
  onSuccess?: (response: Response) => void;
  /**
   * Optional callback invoked when submission fails.
   * Receives the extracted error message for custom error handling.
   */
  onError?: (error: string) => void;
  /** Organized groups of fields with type-safe field names */
  sections: FormSection[];
}
