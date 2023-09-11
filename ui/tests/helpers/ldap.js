import { visit, currentURL } from '@ember/test-helpers';

export const createSecretsEngine = (store) => {
  store.pushPayload('secret-engine', {
    modelName: 'secret-engine',
    data: {
      accessor: 'ldap_7e838627',
      path: 'ldap-test/',
      type: 'ldap',
    },
  });
  return store.peekRecord('secret-engine', 'ldap-test');
};

export const generateBreadcrumbs = (backend, childRoute) => {
  const breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }];
  const root = { label: backend };
  if (childRoute) {
    root.route = 'overview';
    breadcrumbs.push({ label: childRoute });
  }
  breadcrumbs.splice(1, 0, root);
  return breadcrumbs;
};

const baseURL = '/vault/secrets/ldap-test/ldap/';
const stripLeadingSlash = (uri) => (uri.charAt(0) === '/' ? uri.slice(1) : uri);

export const isURL = (uri) => {
  return currentURL() === `${baseURL}${stripLeadingSlash(uri)}`;
};

export const visitURL = (uri) => {
  return visit(`${baseURL}${stripLeadingSlash(uri)}`);
};
