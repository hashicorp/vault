import ExternalLink from './external-link';

/**
 * @module DocLink
 * `DocLink` components are used to render anchor links to relevant Vault documentation at vaultproject.io.
 *
 * @example
 * ```js
    <DocLink @path="/docs/secrets/kv/kv-v2.html">Learn about KV v2</DocLink>
 * ```
 *
 * @param path="/"{String} - The path to documentation on vaultproject.io that the component should link to.
 *
 */
export default class DocLinkComponent extends ExternalLink {
  host = 'https://developer.hashicorp.com';

  get href() {
    return `${this.host}${this.args.path}`;
  }
}
