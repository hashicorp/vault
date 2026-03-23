/**
 * Copyright IBM Corp. 2016, 2025
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

  fetchConfig(type: string, section: string, path: string) {
    switch (type) {
      case 'aws': {
        switch (section) {
          case 'client':
            return this.api.auth.awsReadClientConfiguration(path);
          case 'identity-accesslist':
            return this.api.auth.awsReadIdentityAccessListTidySettings(path);
          case 'roletag-denylist':
            return this.api.auth.awsReadRoleTagDenyListTidySettings(path);
        }
        break;
      }
      case 'azure':
        return this.api.auth.azureReadAuthConfiguration(path);
      case 'github':
        return this.api.auth.githubReadConfiguration(path);
      case 'gcp':
        return this.api.auth.googleCloudReadAuthConfiguration(path);
      case 'jwt':
      case 'oidc':
        return this.api.auth.jwtReadConfiguration(path);
      case 'kubernetes':
        return this.api.auth.kubernetesReadAuthConfiguration(path);
      case 'ldap':
        return this.api.auth.ldapReadAuthConfiguration(path);
      case 'okta':
        return this.api.auth.oktaReadConfiguration(path);
      case 'radius':
        return this.api.auth.radiusReadConfiguration(path);
    }

    throw { httpStatus: 404 };
  }

  schemaForType(type: string, section?: string) {
    if (type === 'aws' && section) {
      return (
        {
          client: 'AwsConfigureClientRequest',
          'identity-accesslist': 'AwsConfigureIdentityAccessListTidyOperationRequest',
          'roletag-denylist': 'AwsConfigureRoleTagDenyListTidyOperationRequest',
        }[section] || ''
      );
    }
    return (
      {
        azure: 'AzureConfigureAuthRequest',
        github: 'GithubConfigureRequest',
        gcp: 'GoogleCloudConfigureAuthRequest',
        jwt: 'JwtConfigureRequest',
        oidc: 'JwtConfigureRequest',
        kubernetes: 'KubernetesConfigureRequest',
        ldap: 'LdapConfigureAuthRequest',
        okta: 'OktaConfigureRequest',
        radius: 'RadiusConfigureRequest',
      }[type] || ''
    );
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
    const form = new OpenApiForm(this.schemaForType(methodType, section), formData, formOptions);
    const defaultGroup = form.formFieldGroups[0]?.['default'] || [];
    // to improve UX, set kubernetes_ca_cert editType to file
    if (methodType === 'kubernetes') {
      const kubernetesCaCertField = defaultGroup.find((field) => field.name === 'kubernetes_ca_cert');
      if (kubernetesCaCertField) {
        kubernetesCaCertField.options.editType = 'file';
      }
    }
    // for jwt and oidc types, the jwks_pairs field is not deprecated but we do not render it in the UI
    // remove the field from the group before rendering the form
    if (['jwt', 'oidc'].includes(methodType)) {
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
