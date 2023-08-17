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
    destination: '/vault/docs/agent-and-proxy/autoauth/:slug',
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
    destination: '/vault/docs/agent-and-proxy/agent/caching/:slug',
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
  }
]
