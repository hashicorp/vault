import Controller from '@ember/controller';

export default class VaultClusterOidcProviderController extends Controller {
  queryParams = [
    'scope', // *
    'response_type', // *
    'client_id', // *
    'redirect_uri', // *
    'state', // *
    'nonce', // *
    'display',
    'prompt',
    'max_age',
    'code_challenge',
    'code_challenge_method',
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
  code_challenge = null;
  code_challenge_method = null;
}
