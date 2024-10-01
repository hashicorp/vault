/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * combineFieldGroups takes the newFields returned from OpenAPI and adds them to the default field group
 * if they are not already accounted for in other field groups
 * @param {Record<string,string[]>[]} currentGroups Field groups, as an array of objects like: [{ default: [] }, { 'TLS options': [] }]
 * @param {string[]} newFields
 * @param {string[]} excludedFields
 * @returns modified currentGroups
 */
export const combineFieldGroups = function (currentGroups, newFields, excludedFields) {
  let allFields = [];
  for (const group of currentGroups) {
    const fieldName = Object.keys(group)[0];
    allFields = allFields.concat(group[fieldName]);
  }
  const otherFields = newFields.filter((field) => {
    return !allFields.includes(field) && !excludedFields.includes(field);
  });
  if (otherFields.length) {
    currentGroups[0].default = currentGroups[0].default.concat(otherFields);
  }

  return currentGroups;
};
