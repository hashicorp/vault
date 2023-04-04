/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import { camelize } from '@ember/string';

export default function camelizeKeys(object) {
  const newObject = {};
  Object.entries(object).forEach(([key, value]) => {
    newObject[camelize(key)] = value;
  });
  return newObject;
}
