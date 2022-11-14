import Component from '@glimmer/component';

/**
 * @module ExternalLinkComponent
 * `ExternalLink` components are used to render anchor links to non-cluster links. Automatically opens in a new tab with noopener noreferrer.
 * To link to vaultproject.io, use DocLink. To link to learn.hashicorp.com, use LearnLink.
 *
 * @example
 * ```js
    <ExternalLink @href="https://hashicorp.com">Arbitrary Link</ExternalLink>
 * ```
 *
 * @param href="https://example.com/"{String} - The full href with protocol
 * @param sameTab=false {Boolean} - by default, these links open in new tab. To override, pass @sameTab={{true}}
 *
 */
export default class ExternalLinkComponent extends Component {
  get href() {
    return this.args.href;
  }
}
