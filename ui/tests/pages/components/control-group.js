import { isPresent, clickable, text } from 'ember-cli-page-object';

export default {
  showsAccessorCallout: isPresent('[data-test-accessor-callout]'),
  authorizationText: text('[data-test-authorizations]'),
  bannerPrefix: text('[data-test-banner-prefix]'),
  bannerText: text('[data-test-banner-text]'),
  requestorText: text('[data-test-requestor-text]'),
  showsTokenText: isPresent('[data-test-token]'),
  refresh: clickable('[data-test-refresh-button]'),
  authorize: clickable('[data-test-authorize-button]'),
  showsSuccessComponent: isPresent('[data-test-control-group-success]'),

  accessor: text('[data-test-accessor-value]'),
  token: text('[data-test-token-value]'),
  showsRefresh: isPresent('[data-test-refresh-button]'),
  showsAuthorize: isPresent('[data-test-authorize-button]'),
  showsBackLink: isPresent('[data-test-back-link]'),
};
