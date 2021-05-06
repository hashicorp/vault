export const ALERT_BANNER_ACTIVE = false

// https://github.com/hashicorp/react-components/tree/master/packages/alert-banner
export default {
  tag: 'New',
  url:
    'https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_campaign=hcp_vault_ga',
  text: 'HCP Vault is now Generally Available',
  linkText: 'Sign up Today',
  // Set the `expirationDate prop with a datetime string (e.g. `2020-01-31T12:00:00-07:00`)
  // if you'd like the component to stop showing at or after a certain date
  expirationDate: '2021-04-30T11:59:00-05:00',
}
