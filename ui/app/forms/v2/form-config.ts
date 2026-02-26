/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import CONFIG_REGISTRY from './generated/index';
import type ApiService from 'vault/services/api';

export type FormConfigKey = keyof typeof CONFIG_REGISTRY;

/**
 * Form element types matching HDS (HashiCorp Design System) components.
 * Currently only TextInput is supported.
 *
 * TODO: Add support for additional field types:
 * - 'Toggle' | 'Checkbox'
 * - 'Select' | 'SuperSelect' | 'Radio' | 'RadioCard'
 * - 'TextArea' | 'MaskedInput'
 * - 'FileInput' | 'KeyValueInput'
 */
export type FormElement = 'TextInput';

/**
 * Union type for all possible field values.
 */
export type FieldValue = string | number | boolean | string[] | null | undefined;

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
  /** Optional placeholder text */
  placeholder?: string;
  /** Default value for the field */
  defaultValue?: FieldValue;
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
  /** Fields belonging to this section */
  fields: FormField[];
}

export interface FormConfig<Request extends object = object, Response = unknown> {
  /** Unique identifier for the form, typically matching the API method name */
  name: string;
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
  /** Organized groups of fields with type-safe field names */
  sections: FormSection[];
}
