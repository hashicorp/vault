import { singularize } from 'ember-inflector';

export const adminPolicy = (mountPath, intMount) => {
  return `
    path "${mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
    path "${intMount}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};
export const readerPolicy = (mountPath, resource) => {
  // keys require singularized paths for GET
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
export const updatePolicy = (mountPath, resource) => {
  return `
    path "${mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${mountPath}/${resource}/*" {
      capabilities = ["read", "update"]
    },
    path "${mountPath}/${singularize(resource)}" {
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
