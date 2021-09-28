import Controller from '@ember/controller';

export default class VaultClusterOidcProviderController extends Controller {
  queryParams = [
    'scope',
    'response_type',
    'client_id',
    'redirect_uri',
    'state',
    'nonce',
    'display',
    'prompt',
    'max_age',
  ];
  scope = null;
  response_type = null;
  client_id = null;
  redirect_uri = null;
  state = null;
  nonce = null;
  display = null;
  prompt = null;
  max_age = null;
}
