import ExternalLink from './external-link';

/**
 * @module LearnLink
 * `LearnLink` components are used to render anchor links to relevant Vault learn documentation at learn.hashicorp.com.
 *
 * @example
 * ```js
    <LearnLink @path="/docs/secrets/kv/kv-v2.html">Learn about KV v2</LearnLink>
 * ```
 *
 * @param path="/"{String} - The path to documentation on learn.hashicorp.com that the component should link to.
 *
 */

// TODO update host to 'https://developer.hashicorp.com' once updated paths are established
export default class LearnLinkComponent extends ExternalLink {
  host = 'https://learn.hashicorp.com';

  get href() {
    return `${this.host}${this.args.path}`;
  }
}
