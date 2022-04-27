import Model, { attr } from '@ember-data/model';

export default class MfaMethod extends Model {
  // common
  @attr('string') type;
  @attr('string') username_template;
  @attr('string') namespace_id;
  @attr('string') mount_accessor;
  // ping id properties
  @attr('string') settings_file_base64;
  @attr('boolean') use_signature;
  @attr('string') idp_url;
  @attr('string') admin_url;
  @attr('string') authenticator_url;
  @attr('string') org_alias;
  // okta properties
  @attr('string') org_name;
  @attr('string') api_token;
  @attr('string') base_url;
  @attr('boolean') primary_email;
  // duo props
  @attr('string') secret_key;
  @attr('string') api_hostname;
  @attr('string') integration_key;
  @attr('string') push_info;
  @attr('boolean') use_passcode;
  @attr('string') pushinfo;
  // totp props
  @attr('string') issuer;
  @attr('int') period;
  @attr('int') key_size;
  @attr('int') qr_size;
  @attr('string') algorithm;
  @attr('int') digits;
  @attr('int') skew;
  @attr('int') max_validation_attempts;
}
