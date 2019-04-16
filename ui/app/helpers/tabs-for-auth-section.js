import { helper as buildHelper } from '@ember/component/helper';
import { pluralize } from 'ember-inflector';

const TABS_FOR_SETTINGS = {
  aws: [
    {
      label: 'Client',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'client'],
    },
    {
      label: 'Identity Whitelist Tidy',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'identity-whitelist'],
    },
    {
      label: 'Role Tag Blacklist Tidy',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'roletag-blacklist'],
    },
  ],
  azure: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  github: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  gcp: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  jwt: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  oidc: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  kubernetes: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  ldap: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  okta: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
  radius: [
    {
      label: 'Configuration',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'configuration'],
    },
  ],
};

const TABS_FOR_SHOW = {};

export function tabsForAuthSection([model, sectionType = 'authSettings', paths]) {
  let tabs;
  debugger; // eslint-disable-line
  if (sectionType === 'authSettings') {
    tabs = (TABS_FOR_SETTINGS[model.type] || []).slice();
    tabs.push({
      label: 'Method Options',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'options'],
    });
    return tabs;
  }

  if (paths) {
    debugger; // eslint-disable-line
    tabs =
      paths.list.length > 0
        ? paths.list.map(path => {
            return {
              label:
                pluralize(path.slice(1))
                  .charAt(0)
                  .toUpperCase() + pluralize(path.slice(1)).slice(1),
              routeParams: ['vault.cluster.access.method.list', path.slice(1)],
            };
          })
        : [];
  } else {
    tabs = (TABS_FOR_SHOW[model.type] || []).slice();
  }
  tabs.push({
    label: 'Configuration',
    routeParams: ['vault.cluster.access.method.section', 'configuration'],
  });

  return tabs;
}

export default buildHelper(tabsForAuthSection);
