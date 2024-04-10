/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { singularize } from 'ember-inflector';

export const adminPolicy = (mountPath: string) => {
  return `
    path "${mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};

// keys require singularized paths for GET
export const readerPolicy = (mountPath: string, resource: string) => {
  return `
    path "${mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${mountPath}/${resource}/*" {
      capabilities = ["read", "list"]
    },
    path "${mountPath}/${singularize(resource)}" {
      capabilities = ["read", "list"]
    },
    path "${mountPath}/${singularize(resource)}/*" {
      capabilities = ["read", "list"]
    },
  `;
};
export const updatePolicy = (mountPath: string, resource: string) => {
  return `
    path "${mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${mountPath}/${resource}/*" {
      capabilities = ["read", "update"]
    },
    path "${mountPath}/${singularize(resource)}/*" {
      capabilities = ["read", "update"]
    },
    path "${mountPath}/issue/*" {
      capabilities = ["update"]
    },
    path "${mountPath}/generate/*" {
      capabilities = ["update"]
    },
    path "${mountPath}/import" {
      capabilities = ["update"]
    },
    path "${mountPath}/sign/*" {
      capabilities = ["update"]
    },
  `;
};
