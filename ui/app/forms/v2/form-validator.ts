/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { get } from '@ember/object';
import type { FieldValue, FormField } from './form-config';
import { validators } from './form-validators';

/**
 * Options that can be passed to validators.
 *
 * @property {boolean} [nullable] - For 'length' and 'number': allow null/undefined values
 * @property {number} [min] - For 'min': minimum numeric value
 * @property {number} [max] - For 'max': maximum numeric value
 * @property {number} [minLength] - For 'minLength': minimum string length
 * @property {number} [maxLength] - For 'maxLength': maximum string length
 * @property {string | RegExp} [pattern] - For 'pattern': regex pattern string or RegExp object
 * @property {string} [flags] - For 'pattern': regex flags (e.g., 'i', 'g', 'gi')
 * @property {string | number | boolean} [value] - For 'isNot': value to compare against
 */
export interface ValidatorOptions {
  nullable?: boolean;
  min?: number;
  max?: number;
  minLength?: number;
  maxLength?: number;
  pattern?: string | RegExp;
  flags?: string;
  value?: string | number | boolean;
}

/**
 * Named validator using a predefined validator type.
 * References a validator in validators.js (e.g., 'required', 'email', 'url').
 *
 * @property {ValidatorType} type - Reference to a validator in validators.js (e.g., 'required', 'email', 'url')
 * @property {string | ((formData: Record<string, unknown>) => string)} message - Error message shown when validation fails
 * @property {ValidatorOptions} [options] - Options passed to the validator function (e.g., { min: 3, max: 10 })
 */
interface NamedValidationRule {
  type: ValidatorType;
  message: string | ((formData: Record<string, unknown>) => string);
  options?: ValidatorOptions;
}

/**
 * Custom validation rule with a custom validator function.
 * The validator function receives the entire form data and returns true if valid.
 *
 * @property {(formData: Record<string, unknown>, options?: ValidatorOptions) => boolean} validator - Custom validator function that receives entire form data
 * @property {string | ((formData: Record<string, unknown>) => string)} message - Error message shown when validation fails
 * @property {ValidatorOptions} [options] - Optional configuration for the validator
 */
interface CustomValidationRule {
  validator: (formData: Record<string, unknown>, options?: ValidatorOptions) => boolean;
  message: string | ((formData: Record<string, unknown>) => string);
  options?: ValidatorOptions;
}

/**
 * Validation rule for a form field.
 * Can be either a named validator (with type) or a custom validator function.
 */
export type ValidationRule = NamedValidationRule | CustomValidationRule;

/**
 * HTML5 standard validator types.
 * These correspond to built-in validators in form-validators.ts.
 */
export type ValidatorType =
  | 'required'
  | 'email'
  | 'url'
  | 'pattern'
  | 'minLength'
  | 'maxLength'
  | 'min'
  | 'max';

/**
 * Default error messages for built-in validators.
 * Used as fallback when validation rule doesn't provide a message.
 */
const DEFAULT_MESSAGES: Record<ValidatorType, string> = {
  required: 'This field is required',
  email: 'Please enter a valid email address',
  url: 'Please enter a valid URL',
  pattern: 'Invalid format',
  minLength: 'Value is too short',
  maxLength: 'Value is too long',
  min: 'Value is too small',
  max: 'Value is too large',
};

/**
 * Validate a single field and return validation errors.
 *
 * @param field - The field configuration to validate
 * @param value - The current value of the field
 * @param payload - The entire form payload for cross-field validation
 * @returns Array of error messages (empty if valid)
 */
export function validateField<TPayload extends object>(
  field: FormField,
  value: FieldValue,
  payload: TPayload
): string[] {
  if (!field.validations || field.validations.length === 0) {
    return [];
  }

  return field.validations
    .filter((rule) => !runValidator(rule, value, payload))
    .map((rule) => {
      // Function message (dynamic based on form data)
      if (typeof rule.message === 'function') {
        return rule.message(payload as Record<string, unknown>);
      }
      // Explicit message provided
      if (rule.message) {
        return rule.message;
      }
      // Fallback to default message for named validators
      if ('type' in rule && rule.type in DEFAULT_MESSAGES) {
        return DEFAULT_MESSAGES[rule.type];
      }
      // Last resort fallback
      return 'Validation failed';
    });
}

/**
 * Validate all fields in a form and return a map of field names to errors.
 *
 * @param fields - Array of all field configurations in the form
 * @param payload - The entire form payload
 * @returns Map of field names to their validation errors
 */
export function validateAllFields<TPayload extends object>(
  fields: FormField[],
  payload: TPayload
): Map<string, string[]> {
  const errors = new Map<string, string[]>();

  for (const field of fields) {
    const fieldErrors = validateField(field, get(payload, field.name) as FieldValue, payload);
    if (fieldErrors.length > 0) {
      errors.set(field.name, fieldErrors);
    }
  }

  return errors;
}

/**
 * Execute a single validation rule on a field value.
 */
function runValidator<TPayload extends object>(
  rule: ValidationRule,
  value: FieldValue,
  payload: TPayload
): boolean {
  // Named validator (type-based)
  if ('type' in rule) {
    const validatorFn = validators[rule.type];
    if (!validatorFn) {
      console.warn(`Unknown validator type: ${rule.type}`);
      return true;
    }
    return validatorFn(value, rule.options || {});
  }

  // Custom validator function
  if ('validator' in rule) {
    return rule.validator(payload as Record<string, unknown>, rule.options);
  }

  return true;
}
