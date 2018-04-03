import { decodeString } from 'vault/utils/b64';

/*
 * @param token - Replication Secondary Activation Token
 * @returns config Object if successful | undefined if not
 *
 */
export default function(token) {
  if (!token) {
    return;
  }
  const tokenParts = token.split('.');
  // config is the second item in the JWT
  let [, configB64] = tokenParts;
  let config;

  if (tokenParts.length !== 3) {
    return;
  }

  // JWTs strip padding from their b64 parts.
  // since we're converting to a typed array before
  // decoding back to utf-8, we need to add any padding back
  while (configB64.length % 4 !== 0) {
    configB64 = configB64 + '=';
  }
  try {
    config = JSON.parse(decodeString(configB64));
  } catch (e) {
    // swallow error
  }

  return config;
}
