import ExternalLink from './external-link';

/**
 * @module DocLink
 * `DocLink` components are used to render anchor links to relevant Vault documentation at developer.hashicorp.com.
 *
 * @example
 * ```js
    <DocLink @path="/vault/docs/secrets/kv/kv-v2.html">Learn about KV v2</DocLink>
 * ```
 *
 * @param {string} path="/" - The path to documentation on developer.hashicorp.com that the component should link to.
 *
 */
export default class DocLinkComponent extends ExternalLink {
  host = 'https://developer.hashicorp.com';

  get href() {
    return `${this.host}${this.args.path}`;
  }
}
