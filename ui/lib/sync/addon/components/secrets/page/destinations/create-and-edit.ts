/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
import { assert } from '@ember/debug';
import { next } from '@ember/runloop';
import { findDestination } from 'core/helpers/sync-destinations';
import apiMethodResolver from 'sync/utils/api-method-resolver';
import { DEFAULT_IDENTITY_TOKEN_TTL } from 'vault/forms/sync/create-destination';

import type { ValidationMap } from 'vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import VersionService from 'vault/services/version';
import {
  DestinationType,
  CredentialType,
  CLOUD_DESTINATION_TYPES,
  WIF_CREDENTIAL_FIELDS,
  ACCOUNT_CREDENTIAL_FIELDS,
  type CloudDestinationType,
} from 'sync/utils/constants';
import type { DestinationForm, DestinationRoleTypeOption } from 'vault/sync';
import type Owner from '@ember/owner';
import AwsSmForm from 'vault/forms/sync/aws-sm';
import AzureKvForm from 'vault/forms/sync/azure-kv';
import GcpSmForm from 'vault/forms/sync/gcp-sm';

type CloudDestinationForm = AwsSmForm | AzureKvForm | GcpSmForm;

interface Args {
  type: DestinationType;
  form: DestinationForm;
}

function isCloudDestinationForm(
  type: DestinationType,
  _form: DestinationForm
): _form is CloudDestinationForm {
  return CLOUD_DESTINATION_TYPES.includes(type as CloudDestinationType);
}

export default class DestinationsCreateForm extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';
  isAccessTypeDisabled = false;

  declare readonly initialCustomTags?: Record<string, string>;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    // cache initial custom tags value to compare against updates
    // tags that are removed when editing need to be added to the payload
    // cast type here since not all types have custom_tags
    const { custom_tags } = args.form.data as unknown as Record<string, unknown>;
    if (custom_tags) {
      this.initialCustomTags = { ...custom_tags };
    }

    // the following checks are only relevant to existing cloud destination configurations with WIF support
    if (this.version.isEnterprise && !args.form.isNew && isCloudDestinationForm(args.type, args.form)) {
      const cloudForm = args.form;
      const { isWifPluginConfigured, isAccountPluginConfigured } = cloudForm;

      assert(
        `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`,
        isWifPluginConfigured !== undefined
      );
      const credentialType = isWifPluginConfigured ? CredentialType.WIF : CredentialType.ACCOUNT;
      next(() => {
        cloudForm.credentialType = credentialType;
        cloudForm.data.credential_type = credentialType;
      });
      // if wif or account only attributes are defined, disable the user's ability to change the access type
      this.isAccessTypeDisabled = isWifPluginConfigured || isAccountPluginConfigured;
    }
  }

  get roleTypeOptions(): Array<DestinationRoleTypeOption> {
    const { type } = this.args;
    const destination = findDestination(type);
    return destination.roleTypeOptions ?? [];
  }

  get header() {
    const { type, form } = this.args;
    const { name: typeDisplayName } = findDestination(type);
    const { name } = form.data;

    return form.isNew
      ? {
          title: `Create Destination for ${typeDisplayName}`,
          breadcrumbs: [
            { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
            { label: 'Secrets sync', route: 'secrets.overview' },
            { label: 'Select destination', route: 'secrets.destinations.create' },
            { label: 'Create destination' },
          ],
        }
      : {
          title: `Edit ${name}`,
          breadcrumbs: [
            { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
            { label: 'Secrets sync', route: 'secrets.overview' },
            { label: 'Destinations', route: 'secrets.destinations' },
            {
              label: 'Destination',
              route: 'secrets.destinations.destination.secrets',
              models: [type, name],
            },
            { label: 'Edit destination' },
          ],
        };
  }

  groupSubtext(group: string, isNew: boolean) {
    const dynamicText = isNew
      ? 'used to authenticate with the destination'
      : 'and the value cannot be read. Enable the input to update';
    switch (group) {
      case 'Advanced configuration':
        return 'Configuration options for the destination.';
      case 'Credentials':
      case 'IAM credentials':
      case 'Client secret':
      case 'JSON credentials':
        return `Connection credentials are sensitive information ${dynamicText}.`;
      default:
        return '';
    }
  }

  isCredentialTypeGroup = (group: string): boolean => {
    const { type } = this.args;
    const credentialGroups = ['WIF credentials', 'IAM credentials', 'Client secret', 'JSON credentials'];

    return CLOUD_DESTINATION_TYPES.includes(type as CloudDestinationType) && credentialGroups.includes(group);
  };

  diffCustomTags(payload: Record<string, unknown>) {
    // if tags were removed we need to add them to the payload
    const { isNew } = this.args.form;
    const { custom_tags } = payload;
    if (!isNew && custom_tags && this.initialCustomTags) {
      // compare the new and old keys of custom_tags object to determine which need to be removed
      const oldKeys = Object.keys(this.initialCustomTags).filter(
        (k) => !Object.keys(custom_tags).includes(k)
      );
      // add tags_to_remove to the payload if there is a diff
      if (oldKeys.length > 0) {
        payload['tags_to_remove'] = oldKeys;
      }
    }
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      this.error = '';
      // clear out validation warnings
      this.modelValidations = null;

      const { form, type } = this.args;
      const { name } = form.data;
      const { isValid, state, invalidFormMessage, data } = form.toJSON();

      this.modelValidations = isValid ? null : state;
      this.invalidFormMessage = isValid ? '' : invalidFormMessage;

      if (isValid) {
        try {
          const payload = data as unknown as Record<string, unknown>;
          // remove credential_type since it's not an actual field on the API payload, it's only used for form validation
          delete payload['credential_type'];
          this.diffCustomTags(payload);
          const method = apiMethodResolver(form.isNew ? 'write' : 'patch', type);
          await this.api.sys[method](name, payload);

          this.router.transitionTo('vault.cluster.sync.secrets.destinations.destination.details', type, name);
          const successMessage = form.isNew
            ? 'You have successfully created a sync destination.'
            : 'You have successfully updated the sync destination.';
          const successTitle = form.isNew ? 'Connection successful' : 'Destination updated';
          this.flashMessages.success(successMessage, {
            title: successTitle,
          });
        } catch (error) {
          const { message } = await this.api.parseError(
            error,
            'Error saving destination. Please try again or contact support.'
          );
          this.error = message;
        }
      }
    })
  );

  @action
  updateWarningValidation() {
    if (this.args.form.isNew) return;
    const { state } = this.args.form.toJSON();
    this.modelValidations = state;
  }

  private resetWifFields(form: CloudDestinationForm, type: CloudDestinationType) {
    const fields = WIF_CREDENTIAL_FIELDS[type];
    fields.forEach((field) => {
      if (field in form.data) {
        (form.data as unknown as Record<string, unknown>)[field] = undefined;
      }
    });
  }

  private resetAccountFields(form: CloudDestinationForm, type: CloudDestinationType) {
    const fields = ACCOUNT_CREDENTIAL_FIELDS[type];
    fields.forEach((field) => {
      if (field in form.data) {
        (form.data as unknown as Record<string, unknown>)[field] = undefined;
      }
    });
  }

  @action
  onTypeChange(option: DestinationRoleTypeOption) {
    const { type, form } = this.args;

    if (!isCloudDestinationForm(type, form)) {
      return;
    }

    form.credentialType = option.value;
    form.data.credential_type = option.value;

    if (option.value === CredentialType.ACCOUNT) {
      this.resetWifFields(form, type as CloudDestinationType);
    } else if (option.value === CredentialType.WIF) {
      this.resetAccountFields(form, type as CloudDestinationType);
      form.data.identity_token_ttl = DEFAULT_IDENTITY_TOKEN_TTL;
    }

    this.modelValidations = null;
    this.invalidFormMessage = '';
  }

  @action
  cancel() {
    const route = this.args.form.isNew ? 'create' : 'destination';
    this.router.transitionTo(`vault.cluster.sync.secrets.destinations.${route}`);
  }
}
