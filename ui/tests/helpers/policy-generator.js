import { singularize } from 'ember-inflector';

export const adminPolicy = (mountPath) => {
  return `
    path "${mountPath}/*" {
      capabilities = ["create", "read", "update", "delete", "list"]
    },
  `;
};
export const readerPolicy = (mountPath, resource, includeSingular = false) => {
  if (includeSingular) {
    return `
    path "${this.mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${this.mountPath}/${resource}/*" {
      capabilities = ["read", "list"]
    },
    path "${this.mountPath}/${singularize(resource)}" {
      capabilities = ["read", "list"]
    },
    path "${this.mountPath}/${singularize(resource)}/*" {
      capabilities = ["read", "list"]
    },
  `;
  }
  return `
    path "${this.mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${this.mountPath}/${resource}/*" {
      capabilities = ["read", "list"]
    },
  `;
};

export const updatePolicy = (mountPath, resource) => {
  return `
    path "${this.mountPath}/${resource}" {
      capabilities = ["read", "list"]
    },
    path "${this.mountPath}/${resource}/*" {
      capabilities = ["read", "update"]
    },
    path "${this.mountPath}/issue/*" {
      capabilities = ["update"]
    },
    path "${this.mountPath}/generate/*" {
      capabilities = ["update"]
    },
    path "${this.mountPath}/import" {
      capabilities = ["update"]
    },
    path "${this.mountPath}/sign/*" {
      capabilities = ["update"]
    },
  `;
};
