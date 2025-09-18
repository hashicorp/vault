/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { capitalize } from '@ember/string';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';
import { isPresent } from '@ember/utils';

const METHOD_PROPS = {
  common: [],
  duo: ['username_format', 'secret_key', 'integration_key', 'api_hostname', 'push_info', 'use_passcode'],
  okta: ['username_format', 'mount_accessor', 'org_name', 'api_token', 'base_url', 'primary_email'],
  totp: [
    'issuer',
    'period',
    'key_size',
    'qr_size',
    'algorithm',
    'digits',
    'skew',
    'max_validation_attempts',
    'enable_self_enrollment',
  ],
  pingid: [
    'username_format',
    'settings_file_base64',
    'use_signature',
    'idp_url',
    'admin_url',
    'authenticator_url',
    'org_alias',
  ],
};

const REQUIRED_PROPS = {
  duo: ['secret_key', 'integration_key', 'api_hostname'],
  okta: ['org_name', 'api_token'],
  totp: ['issuer'],
  pingid: ['settings_file_base64'],
};

const validators = Object.keys(REQUIRED_PROPS).reduce((obj, type) => {
  REQUIRED_PROPS[type].forEach((prop) => {
    obj[`${prop}`] = [
      {
        message: `${prop.replace(/_/g, ' ')} is required`,
        validator(model) {
          return model.type === type ? isPresent(model[prop]) : true;
        },
      },
    ];
  });
  return obj;
}, {});

@withModelValidations(validators)
export default class MfaMethod extends Model {
  // common
  @attr('string') type;
  @attr('string', {
    label: 'Username format',
    subText: 'How to map identity names to MFA method names. ',
  })
  username_format;
  @attr('string', {
    label: 'Namespace',
  })
  namespace_path;
  @attr('string') mount_accessor;

  // PING ID
  @attr('string', {
    label: 'Settings file',
    subText: 'A base-64 encoded third party setting file retrieved from the PingIDs configuration page.',
  })
  settings_file_base64;
  @attr('boolean') use_signature;
  @attr('string') idp_url;
  @attr('string') admin_url;
  @attr('string') authenticator_url;
  @attr('string') org_alias;

  // OKTA
  @attr('string', {
    label: 'Organization name',
    subText: 'Name of the organization to be used in the Okta API.',
  })
  org_name;
  @attr('string', {
    label: 'Okta API key',
  })
  api_token;
  @attr('string', {
    label: 'Base URL',
    subText:
      'If set, will be used as the base domain for API requests. Example are okta.com, oktapreview.com and okta-emea.com.',
  })
  base_url;
  @attr('boolean') primary_email;

  // DUO
  @attr('string', {
    label: 'Duo secret key',
    sensitive: true,
  })
  secret_key;
  @attr('string', {
    label: 'Duo integration key',
    sensitive: true,
  })
  integration_key;
  @attr('string', {
    label: 'Duo API hostname',
  })
  api_hostname;
  @attr('string', {
    label: 'Duo push information',
    subText: 'Additional information displayed to the user when the push is presented to them.',
  })
  push_info;
  @attr('boolean', {
    label: 'Passcode reminder',
    subText: 'If this is turned on, the user is reminded to use the passcode upon MFA validation.',
  })
  use_passcode;

  // TOTP
  @attr('string', {
    label: 'Issuer',
    subText: 'The human-readable name of the keys issuing organization.',
  })
  issuer;
  @attr({
    label: 'Period',
    editType: 'ttl',
    helperTextEnabled: 'How long each generated TOTP is valid.',
    hideToggle: true,
    defaultValue: 30, // API accepts both an integer as seconds and sting with unit e.g 30 || '30s'
  })
  period;
  @attr('number', {
    label: 'Key size',
    subText: 'The size in bytes of the Vault generated key.',
  })
  key_size;
  @attr('number', {
    label: 'QR size',
    subText: 'The pixel size of the generated square QR code.',
  })
  qr_size;
  @attr('string', {
    label: 'Algorithm',
    editType: 'radio',
    possibleValues: ['SHA1', 'SHA256', 'SHA512'],
    subText: 'The hashing algorithm used to generate the TOTP code.',
  })
  algorithm;
  @attr('number', {
    label: 'Digits',
    editType: 'radio',
    possibleValues: [6, 8],
    subText: 'The number digits in the generated TOTP code.',
  })
  digits;
  @attr('number', {
    label: 'Skew',
    editType: 'radio',
    possibleValues: [0, 1],
    subText: 'The number of delay periods allowed when validating a TOTP token.',
  })
  skew;
  @attr('number') max_validation_attempts;
  @attr('boolean', {
    label: 'Enable self-enrollment',
    editType: 'toggleButton',
    helperTextEnabled:
      'Let end users enroll in this MFA method on their own. You still control which auth mounts, groups, or entities it applies to.',
    helperTextDisabled:
      'Let end users enroll in this MFA method on their own. You still control which auth mounts, groups, or entities it applies to.',
  })
  enable_self_enrollment;

  get name() {
    return this.type === 'totp' ? this.type.toUpperCase() : capitalize(this.type);
  }

  get icon() {
    switch (this.type) {
      case 'totp':
        return 'history';
      case 'pingid':
        return 'ping-identity-color';
      case 'duo':
        return 'duo-color';
      default:
        return this.type;
    }
  }

  get formFields() {
    return [...METHOD_PROPS.common, ...METHOD_PROPS[this.type]];
  }

  get attrs() {
    return expandAttributeMeta(this, this.formFields);
  }
}
