/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module SplitObject
 * SplitObject helper takes in a class of data as the first param and an array of keys that you want to split into another object as the second param. 
 * You will end up with an array of two objects. One no longer with the array of params, and the second with just the array of params.
 *
 * @example
 * ```js
 * splitObject(data, ['max_versions', 'delete_version_after', 'cas_required'])
 * ```
 
 * @param {object} - The object you want to split into two. This object will have all the keys from the second param (the array param).
 * @param {array} - An array of params that you want to split off the object and turn into its own object.

 */
import { helper as buildHelper } from '@ember/component/helper';

export function splitObject(originalObject, array) {
  const object1 = {};
  const object2 = {};
  // convert object to key's array
  const keys = Object.keys(originalObject);
  keys.forEach((key) => {
    if (array.includes(key)) {
      object1[key] = originalObject[key];
    } else {
      object2[key] = originalObject[key];
    }
  });
  return [object1, object2];
}

export default buildHelper(splitObject);
