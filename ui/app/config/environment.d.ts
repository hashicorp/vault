/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Type declarations for
 *    import config from 'my-app/config/environment'
 */
declare const config: {
  environment: string;
  modulePrefix: string;
  podModulePrefix: string;
  locationType: 'history' | 'hash' | 'none';
  rootURL: string;
  APP: {
    POLLING_URLS: string[];
    NAMESPACE_ROOT_URLS: string[];
    DEFAULT_PAGE_SIZE: number;
    LOG_TRANSITIONS?: boolean;
    LOG_ACTIVE_GENERATION?: boolean;
    LOG_VIEW_LOOKUPS?: boolean;
    rootElement?: string;
    autoboot?: boolean;
  };
};

export default config;
