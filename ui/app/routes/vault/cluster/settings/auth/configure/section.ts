/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import AuthMethodForm from 'vault/forms/auth/method';
import OpenApiForm from 'vault/forms/open-api';

import type ApiService from 'vault/services/api';
import type PathHelpService from 'vault/services/path-help';
import type Store from '@ember-data/store';
import type { ClusterSettingsAuthConfigureRouteModel } from '../configure';
import type { MountConfig } from 'vault/mount';
import type { HTTPRequestInit, RequestOpts } from '@hashicorp/vault-client-typescript';
import type { OpenApiHelpResponse } from 'vault/utils/openapi-helpers';

export default class ClusterSettingsAuthConfigureRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly pathHelp: PathHelpService;
  @service declare readonly store: Store;

  get configRouteModel() {
    return this.modelFor('vault.cluster.settings.auth.configure') as ClusterSettingsAuthConfigureRouteModel;
  }

  modelForOptions() {
    const { methodOptions, method } = this.configRouteModel;
    const config = methodOptions.config as MountConfig;
    const listing_visibility = config.listing_visibility === 'unauth' ? true : false;

    const form = new AuthMethodForm({
      ...methodOptions,
      path: method.id,
      config: { ...config, listing_visibility },
      user_lockout_config: {},
    });
    // `type` is a tracked property on the form class and not a field param
    // so we need to set it here even though it is spread in methodOptions above.
    form.type = methodOptions.type;

    return {
      form,
      section: 'options',
    };
  }

  get configFieldGroupsMap() {
    const { method } = this.configRouteModel;
    return {
      kubernetes: {
        default: ['kubernetes_host', 'kubernetes_ca_cert', 'disable_local_ca_jwt'],
        'Kubernetes Options': ['token_reviewer_jwt', 'pem_keys', 'use_annotations_as_alias_metadata'],
      },
    }[method.methodType];
  }

  fetchConfig(type: string, section: string, path: string, help = false) {
    const initOverride = help
      ? (context: { init: HTTPRequestInit; context: RequestOpts }) =>
          this.api.addQueryParams(context, { help: 1 })
      : undefined;

    switch (type) {
      case 'aws': {
        switch (section) {
          case 'client':
            return this.api.auth.awsReadClientConfiguration(path, initOverride);
          case 'identity-accesslist':
            return this.api.auth.awsReadIdentityAccessListTidySettings(path, initOverride);
          case 'roletag-denylist':
            return this.api.auth.awsReadRoleTagDenyListTidySettings(path, initOverride);
        }
        break;
      }
      case 'azure':
        return this.api.auth.azureReadAuthConfiguration(path, initOverride);
      case 'github':
        return this.api.auth.githubReadConfiguration(path, initOverride);
      case 'gcp':
        return this.api.auth.googleCloudReadAuthConfiguration(path, initOverride);
      case 'jwt':
      case 'oidc':
        return this.api.auth.jwtReadConfiguration(path, initOverride);
      case 'kubernetes':
        return this.api.auth.kubernetesReadAuthConfiguration(path, initOverride);
      case 'ldap':
        return this.api.auth.ldapReadAuthConfiguration(path, initOverride);
      case 'okta':
        return this.api.auth.oktaReadConfiguration(path, initOverride);
      case 'radius':
        return this.api.auth.radiusReadConfiguration(path, initOverride);
    }

    throw { httpStatus: 404 };
  }

  async modelForConfiguration(section: string) {
    const { path, methodType } = this.configRouteModel.method;

    const formOptions = { isNew: false };
    let formData;
    // make request to fetch configuration data for method
    try {
      const { data } = await this.fetchConfig(methodType, section, path);
      formData = data as object;
    } catch (e) {
      const { message, status } = await this.api.parseError(e);
      if (status === 404) {
        formOptions.isNew = true;
      } else {
        throw { message, httpsStatus: status };
      }
    }
    // make request to fetch OpenAPI properties with help query param
    const helpResponse = (await this.fetchConfig(
      methodType,
      section,
      path,
      true
    )) as unknown as OpenApiHelpResponse;
    const form = new OpenApiForm(helpResponse, formData, formOptions);
    // for jwt and oidc types, the jwks_pairs field is not deprecated but we do not render it in the UI
    // remove the field from the group before rendering the form
    if (['jwt', 'oidc'].includes(methodType)) {
      const defaultGroup = form.formFieldGroups[0]?.['default'] || [];
      const index = defaultGroup.findIndex((field) => field.name === 'jwks_pairs');
      if (index !== undefined && index >= 0) {
        defaultGroup.splice(index, 1);
      }
    }

    return {
      form,
      section,
      method: this.configRouteModel.method,
    };
  }

  model(params: { section_name: 'options' | 'configuration' }) {
    const { section_name: section } = params;
    return section === 'options' ? this.modelForOptions() : this.modelForConfiguration(section);
  }
}
