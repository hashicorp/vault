/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

module.exports = [
  // example redirect:
  // {
  //   source: '/vault/docs/some/path',
  //   destination: '/vault/docs/some/other/path',
  //   permanent: true,
  // },
  {
    source: '/vault/docs/plugins/plugin-portal',
    destination: '/vault/integrations',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/autoauth',
    destination: '/vault/docs/agent-and-proxy/autoauth',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/autoauth/:slug',
    destination: '/vault/docs/agent-and-proxy/autoauth/:slug*',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/template',
    destination: '/vault/docs/agent-and-proxy/agent/template',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/winsvc',
    destination: '/vault/docs/agent-and-proxy/agent/winsvc',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/versions',
    destination: '/vault/docs/agent-and-proxy/agent/versions',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/apiproxy',
    destination: '/vault/docs/agent-and-proxy/agent/apiproxy',
    permanent: true,
  },
  {
    source: '/vault/docs/agent',
    destination: '/vault/docs/agent-and-proxy/agent',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching',
    destination: '/vault/docs/agent-and-proxy/agent/caching',
    permanent: true,
  },
  {
    source: '/vault/docs/agent/caching/:slug',
    destination: '/vault/docs/agent-and-proxy/agent/caching/:slug*',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:4|5|6|7|8|9|10|11|12|13)\.x)/agent-and-proxy/agent',
    destination: '/vault/docs/:version/agent/',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:4|5|6|7|8|9|10|11|12|13)\.x)/agent-and-proxy/agent/template',
    destination: '/vault/docs/:version/agent/template',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:4|5|6|7|8|9|10|11|12|13)\.x)/agent-and-proxy/agent/caching',
    destination: '/vault/docs/:version/agent/caching',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:4|5|6|7|8|9|10|11|12|13)\.x)/agent-and-proxy/autoauth/:slug*',
    destination: '/vault/docs/:version/agent/autoauth/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:8|9|10|11|12|13)\.x)/agent-and-proxy/agent/caching/:slug*',
    destination: '/vault/docs/:version/agent/caching/:slug',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:7|8|9|10|11|12|13)\.x)/agent-and-proxy/agent/winsvc',
    destination: '/vault/docs/:version/agent/winsvc',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:8|9)\.x)/agent-and-proxy/agent/generate-config',
    destination: '/vault/docs/:version/agent/template-config',
    permanent: true,
  },
  {
    source: '/vault/docs/v1.13.x/agent-and-proxy/agent/versions',
    destination: '/vault/docs/v1.13.x/agent/versions',
    permanent: true,
  },
  {
    source: '/vault/docs/v1.13.x/agent-and-proxy/agent/apiproxy',
    destination: '/vault/docs/v1.13.x/agent/apiproxy',
    permanent: true,
  },
  {
    source: '/vault/api-docs/system/plugins-reload-backend',
    destination: '/vault/api-docs/system/plugins-reload',
    permanent: true,
  },
  {
    source: '/vault/docs/deprecation/faq',
    destination: '/vault/docs/deprecation',
    permanent: true,
  },
  {
    source: '/vault/docs/concepts/lease-explosions',
    destination: '/vault/docs/configuration/prevent-lease-explosions',
    permanent: true,
  },
  {
    source: '/vault/docs/troubleshoot/lease-explosions',
    destination: '/vault/docs/configuration/prevent-lease-explosions',
    permanent: true,
  },
  {
    source: '/vault/docs/concepts/lease-count-quota-exceeded',
    destination: '/vault/docs/troubleshoot/lease-count-quota-exceeded',
    permanent: true,
  },
  {
    source: '/vault/docs/command/web',
    destination: '/vault/docs/ui/web-cli',
    permanent: true,
  },
  {
    source: '/vault/docs/deprecation',
    destination: '/vault/docs/updates/deprecation',
    permanent: true,
  },
  {
    source: '/vault/docs/:version(v1\.(?:8|9)\.x)/agent-and-proxy/agent/generate-config',
    destination: '/vault/docs/:version/agent/template-config',
    permanent: true,
  },
  {
    source: '/vault/docs/faq/ssct',
    destination: '/vault/docs/v1.10.x/faq/ssct',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading',
    destination: '/vault/docs/upgrade',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/raft-wal',
    destination: '/vault/docs/upgrade/raft-wal',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/vault-ha-upgrade',
    destination: '/vault/docs/upgrade/vault-ha-upgrade',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/plugins',
    destination: '/vault/docs/plugins/upgrade',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/upgrade-to-1.19.x',
    destination: '/vault/docs/v1.19.x/updates/important-changes',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/deduplication/:slug*',
    destination: '/vault/docs/secrets/identity/deduplication/:slug*',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/upgrade-to-:version(1\.(?:12|13|14|15|16|17|18)\.x)',
    destination: '/vault/docs/v:version/upgrading/upgrade-to-:version',
    permanent: true,
  },
  {
    source: '/vault/docs/release-notes/1.19.0',
    destination: '/vault/docs/v1.19.x/updates/release-notes',
    permanent: true,
  },
  {
    source: '/vault/docs/v:version(1\.(?:4|5|6|7|8|9|10|11|12|13|14|15|16|17|18)\.x)/updates/important-changes',
    destination: '/vault/docs/v:version/upgrading/upgrade-to-:version',
    permanent: true,
  },
  {
    source: '/vault/docs/v:version(1\.(?:4|5|6|7|8|9|10|11|12|13|14|15|17|18)).x/updates/release-notes',
    destination: '/vault/docs/v:version.x/release-notes/:version.0',
    permanent: true,
  },
  {
    source: '/vault/docs/v1.16.x/updates/release-notes',
    destination: '/vault/docs/v1.16.x/release-notes/1.16.1',
    permanent: true,
  },
  {
    source: '/vault/docs/release-notes/:version(1\.(?:4|5|6|7|8|9|10|11|12|13|14|15|17|18)).0',
    destination: '/vault/docs/v:version.x/release-notes/:version.0',
    permanent: true,
  },
  {
    source: '/vault/docs/release-notes/1.16.1',
    destination: '/vault/docs/v1.16.x/release-notes/1.16.1',
    permanent: true,
  },
  {
    source: '/vault/docs/what-is-vault',
    destination: '/vault/docs/about-vault/what-is-vault',
    permanent: true,
  },
  {
    source: '/vault/docs/use-cases',
    destination: '/vault/docs/about-vault/why-use-vault',
    permanent: true,
  },
  {
    source: '/vault/docs/interoperability-matrix',
    destination: '/vault/docs/partners',
    permanent: true,
  },
  {
    source: '/vault/docs/v:version(1\.(?:4|5|6|7|8|9|10|11|12|13|14|15|16|17|18)\.x)/partners',
    destination: '/vault/docs/:version/interoperability-matrix',
    permanent: true,
  },
  {
    source: '/vault/docs/partnerships',
    destination: '/vault/docs/partners/program',
    permanent: true,
  },
  {
    source: '/vault/docs/run-as-service',
    destination: '/vault/docs/deploy/run-as-service',
    permanent: true,
  },
  {
    source: '/vault/docs/install/:slug*',
    destination: '/vault/docs/get-vault/:slug*',
    permanent: true,
  },
  {
    source: '/vault/docs/platform/aws/:slug*',
    destination: '/vault/docs/deploy/aws/:slug*',
    permanent: true,
  },
  {
    source: '/vault/docs/platform/k8s/:slug*',
    destination: '/vault/docs/deploy/kubernetes/:slug*',
    permanent: true,
  },
  {
    source: '/vault/api-docs/secret/ad',
    destination: '/vault/api-docs/secret/ldap',
    permanent: true,
  },
  {
    source: '/vault/docs/secrets/ad',
    destination: '/vault/docs/secrets/ldap',
    permanent: true,
  },
  {
    source: '/vault/docs/secrets/ad/migration-guide',
    destination: '/vault/docs/v1.18.x/secrets/ad/migration-guide',
    permanent: true,
  },
  {
    source: '/vault/docs/upgrading/vault-ha-upgrade',
    destination: '/vault/docs/v1.10.x/upgrading/vault-ha-upgrade',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license',
    destination: '/vault/docs/license',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license/autoloading',
    destination: '/vault/docs/license/autoloading',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license/utilization-reporting',
    destination: '/vault/docs/license/utilization/auto-reporting',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license/manual-reporting',
    destination: '/vault/docs/license/utilization/manual-reporting',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license/product-usage-reporting',
    destination: '/vault/docs/license/product-usage-reporting',
    permanent: true,
  },
  {
    source: '/vault/docs/enterprise/license/faq',
    destination: '/vault/docs/license',
    permanent: true,
  }
]