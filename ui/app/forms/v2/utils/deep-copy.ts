/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

type CloneableValue = unknown;

const isPlainObject = (value: unknown): value is Record<string, unknown> => {
  if (value === null || typeof value !== 'object') {
    return false;
  }
  const prototype = Object.getPrototypeOf(value);
  return prototype === Object.prototype || prototype === null;
};

/*
 * Deep-copy arrays and plain objects while preserving function references.
 * This avoids mutating config state without stripping function-valued
 * properties such as visibility predicates and custom validators.
 */
export const deepCopyValue = (value: CloneableValue, seen = new WeakMap<object, unknown>()): unknown => {
  if (value === null || value === undefined) {
    return value;
  }

  if (typeof value === 'function') {
    return value;
  }

  if (Array.isArray(value)) {
    if (seen.has(value)) {
      return seen.get(value);
    }
    const clonedArray: unknown[] = [];
    seen.set(value, clonedArray);
    for (const item of value) {
      clonedArray.push(deepCopyValue(item, seen));
    }
    return clonedArray;
  }

  if (value instanceof Date) {
    return new Date(value.getTime());
  }

  if (value instanceof RegExp) {
    return new RegExp(value);
  }

  if (isPlainObject(value)) {
    if (seen.has(value)) {
      return seen.get(value);
    }
    const clonedObject: Record<string, unknown> = {};
    seen.set(value, clonedObject);
    for (const [key, nestedValue] of Object.entries(value)) {
      clonedObject[key] = deepCopyValue(nestedValue, seen);
    }
    return clonedObject;
  }

  return value;
};
