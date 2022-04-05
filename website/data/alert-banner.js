export const ALERT_BANNER_ACTIVE = false

// https://github.com/hashicorp/web-components/tree/master/packages/alert-banner
export default {
  tag: 'New',
  url: 'https://www.hashicorp.com/blog/login-mfa-support-added-to-vault-open-source-and-hcp-vault',
  text: 'Vault OSS Now Includes Multi-factor Authentication!',
  linkText: '',
  // Set the expirationDate prop with a datetime string (e.g. '2020-01-31T12:00:00-07:00')
  // if you'd like the component to stop showing at or after a certain date
  expirationDate: '2022-05-31T23:00:00-07:00',
}
