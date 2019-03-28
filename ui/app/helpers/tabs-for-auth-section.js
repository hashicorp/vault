import { helper as buildHelper } from '@ember/component/helper';

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

export function tabsForAuthSection([methodType, sectionType = 'authSettings']) {
  let tabs;

  if (sectionType === 'authSettings') {
    tabs = (TABS_FOR_SETTINGS[methodType] || []).slice();
    tabs.push({
      label: 'Method Options',
      routeParams: ['vault.cluster.settings.auth.configure.section', 'options'],
    });
    return tabs;
  }

  tabs = (TABS_FOR_SHOW[methodType] || []).slice();
  tabs.push({
    label: 'Configuration',
    routeParams: ['vault.cluster.access.method.section', 'configuration'],
  });

  return tabs;
}

export default buildHelper(tabsForAuthSection);
