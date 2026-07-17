/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { ValidatorOptions } from './form-validator';

/**
 * Built-in validator functions
 * All validators return true if valid, false if invalid
 */
export const validators = {
  /**
   * Required - value must be present
   * Rejects: null, undefined, '', [], {}
   */
  required: (value: unknown, _options?: ValidatorOptions): boolean => {
    if (value === null || value === undefined) return false;
    if (typeof value === 'string') return value.trim().length > 0;
    if (Array.isArray(value)) return value.length > 0;
    if (typeof value === 'object') return Object.keys(value).length > 0;
    return true;
  },

  /**
   * Email - validates email format
   */
  email: (value: unknown, _options?: ValidatorOptions): boolean => {
    if (!value) return true; // Use with 'required' for mandatory emails
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(String(value));
  },

  /**
   * URL - validates URL format
   */
  url: (value: unknown, _options?: ValidatorOptions): boolean => {
    if (!value) return true;
    try {
      new URL(String(value));
      return true;
    } catch {
      return false;
    }
  },

  /**
   * Pattern - validates against regex pattern
   */
  pattern: (value: unknown, { pattern, flags = '' }: ValidatorOptions): boolean => {
    if (!value || !pattern) return true;
    const regex = typeof pattern === 'string' ? new RegExp(pattern, flags) : pattern;
    return regex.test(String(value));
  },

  /**
   * MinLength - validates minimum string length
   */
  minLength: (value: unknown, { minLength }: ValidatorOptions): boolean => {
    if (minLength === undefined) return true;
    if (!value) return false;
    return String(value).length >= minLength;
  },

  /**
   * MaxLength - validates maximum string length
   */
  maxLength: (value: unknown, { maxLength }: ValidatorOptions): boolean => {
    if (maxLength === undefined) return true;
    if (!value) return true;
    return String(value).length <= maxLength;
  },

  /**
   * Min - validates minimum numeric value
   */
  min: (value: unknown, { min }: ValidatorOptions): boolean => {
    if (min === undefined) return true;
    if (value === null || value === undefined || value === '') return true;
    const num = Number(value);
    if (isNaN(num)) return false;
    return num >= min;
  },

  /**
   * Max - validates maximum numeric value
   */
  max: (value: unknown, { max }: ValidatorOptions): boolean => {
    if (max === undefined) return true;
    if (value === null || value === undefined || value === '') return true;
    const num = Number(value);
    if (isNaN(num)) return false;
    return num <= max;
  },
};
