export const ALERT_BANNER_ACTIVE = true

// https://github.com/hashicorp/web-components/tree/master/packages/alert-banner
export default {
  tag: 'Blog post',
  url: 'https://www.hashicorp.com/blog/a-new-chapter-for-hashicorp',
  text:
    'HashiCorp shares have begun trading on the Nasdaq. Read the blog from our founders, Mitchell Hashimoto and Armon Dadgar.',
  linkText: 'Read the post',
  // Set the expirationDate prop with a datetime string (e.g. '2020-01-31T12:00:00-07:00')
  // if you'd like the component to stop showing at or after a certain date
  expirationDate: '2021-12-17T23:00:00-07:00',
}
