# API Support

This document tracks LinodeGo support for the features of the [Linode API](https://developers.linode.com/changelog/api/).

Endpoints are implemented as needed, by need or user-request.  As new features are added (as reported in the [Linode API Changelog](https://developers.linode.com/changelog/api/)) this document should be updated to reflect any missing endpoints.  New or deprecated fields should also be indicated below the affected HTTP method, for example:

```markdown
- `/fake/endpoint`
  - [x] `GET`
        * 4.0.29 field X is not implemented <http://...>
  - [ ] `POST`
        * 4.0.30 added support to create Fake things <http://...>
```

See `template.go` and `template_test.go` for tips on adding new endpoints.

## Linodes

- `/linode/instances`
  - [x] `GET`
  - [X] `POST`
- `/linode/instances/$id`
  - [x] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/linode/instances/$id/boot`
  - [x] `POST`
- `/linode/instances/$id/clone`
  - [x] `POST`
- `/linode/instances/$id/mutate`
  - [X] `POST`
- `/linode/instances/$id/reboot`
  - [x] `POST`
- `/linode/instances/$id/rebuild`
  - [X] `POST`
- `/linode/instances/$id/rescue`
  - [X] `POST`
- `/linode/instances/$id/resize`
  - [x] `POST`
- `/linode/instances/$id/shutdown`
  - [x] `POST`
- `/linode/instances/$id/volumes`
  - [X] `GET`

### Backups

- `/linode/instances/$id/backups`
  - [X] `GET`
  - [ ] `POST`
- `/linode/instances/$id/backups/$id/restore`
  - [ ] `POST`
- `/linode/instances/$id/backups/cancel`
  - [ ] `POST`
- `/linode/instances/$id/backups/enable`
  - [ ] `POST`

### Configs

- `/linode/instances/$id/configs`
  - [X] `GET`
  - [X] `POST`
- `/linode/instances/$id/configs/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`

### Disks

- `/linode/instances/$id/disks`
  - [X] `GET`
  - [X] `POST`
- `/linode/instances/$id/disks/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `POST`
  - [X] `DELETE`
- `/linode/instances/$id/disks/$id/password`
  - [X] `POST`
- `/linode/instances/$id/disks/$id/resize`
  - [X] `POST`

### IPs

- `/linode/instances/$id/ips`
  - [X] `GET`
  - [X] `POST`
- `/linode/instances/$id/ips/$ip_address`
  - [X] `GET`
  - [X] `PUT`
  - [ ] `DELETE`
- `/linode/instances/$id/ips/sharing`
  - [ ] `POST`

### Kernels

- `/linode/kernels`
  - [X] `GET`
- `/linode/kernels/$id`
  - [X] `GET`

### StackScripts

- `/linode/stackscripts`
  - [x] `GET`
  - [X] `POST`
- `/linode/stackscripts/$id`
  - [x] `GET`
  - [X] `PUT`
  - [X] `DELETE`

### Stats

- `/linode/instances/$id/stats`
  - [X] `GET`
- `/linode/instances/$id/stats/$year/$month`
  - [X] `GET`

### Types

- `/linode/types`
  - [X] `GET`
- `/linode/types/$id`
  - [X] `GET`

## Domains

- `/domains`
  - [X] `GET`
  - [X] `POST`
- `/domains/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/domains/$id/clone`
  - [ ] `POST`
- `/domains/$id/records`
  - [X] `GET`
  - [X] `POST`
- `/domains/$id/records/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`

## LKE

- `/lke/clusters`
  - [X] `POST`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/lke/clusters/$id/pools`
  - [X] `POST`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/lke/clusters/$id/api-endpoint`
  - [X] `GET`
- `/lke/clusters/$id/kubeconfig`
  - [X] `GET`
- `/lke/clusters/$id/versions`
  - [X] `GET`
- `/lke/clusters/$id/versions/$id`
  - [X] `GET`

## Longview

- `/longview/clients`
  - [X] `GET`
  - [ ] `POST`
- `/longview/clients/$id`
  - [X] `GET`
  - [ ] `PUT`
  - [ ] `DELETE`

### Subscriptions

- `/longview/subscriptions`
  - [ ] `GET`
- `/longview/subscriptions/$id`
  - [ ] `GET`

### NodeBalancers

- `/nodebalancers`
  - [X] `GET`
  - [X] `POST`
- `/nodebalancers/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/nodebalancers/$id/stats`
  - [X] `GET`

### NodeBalancer Configs

- `/nodebalancers/$id/configs`
  - [X] `GET`
  - [X] `POST`
- `/nodebalancers/$id/configs/$id`
  - [X] `GET`
  - [X] `DELETE`
- `/nodebalancers/$id/configs/$id/nodes`
  - [X] `GET`
  - [X] `POST`
- `/nodebalancers/$id/configs/$id/nodes/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/nodebalancers/$id/configs/$id/rebuild`
  - [X] `POST`

## Networking

- `/networking/ip-assign`
  - [ ] `POST`
- `/networking/ips`
  - [X] `GET`
  - [ ] `POST`
- `/networking/ips/$address`
  - [X] `GET`
  - [X] `PUT`
  - [ ] `DELETE`

### IPv6

- `/networking/ips`
  - [X] `GET`
- `/networking/ips/$address`
  - [X] `GET`
  - [ ] `PUT`
- `/networking/ipv6/ranges`
  - [X] `GET`
- `/networking/ipv6/pools`
  - [X] `GET`

## Regions

- `/regions`
  - [x] `GET`
- `/regions/$id`
  - [x] `GET`

## Support

- `/support/tickets`
  - [X] `GET`
  - [ ] `POST`
- `/support/tickets/$id`
  - [X] `GET`
- `/support/tickets/$id/attachments`
  - [ ] `POST`
- `/support/tickets/$id/replies`
  - [ ] `GET`
  - [ ] `POST`

## Tags

- `/tags/`
  - [X] `GET`
  - [X] `POST`
- `/tags/$id`
  - [X] `GET`
  - [X] `DELETE`

## Account

### Events

- `/account/events`
  - [X] `GET`
- `/account/events/$id`
  - [X] `GET`
- `/account/events/$id/read`
  - [X] `POST`
- `/account/events/$id/seen`
  - [X] `POST`

### Invoices

- `/account/invoices/`
  - [X] `GET`
- `/account/invoices/$id`
  - [X] `GET`
- `/account/invoices/$id/items`
  - [X] `GET`

### Notifications

- `/account/notifications`
  - [X] `GET`

### OAuth Clients

- `/account/oauth-clients`
  - [X] `GET`
  - [X] `POST`
- `/account/oauth-clients/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/account/oauth-clients/$id/reset_secret`
  - [ ] `POST`
- `/account/oauth-clients/$id/thumbnail`
  - [ ] `GET`
  - [ ] `PUT`

### Object Storage Keys

- `/object-storage/keys`
  - [X] `GET`
  - [X] `POST`
- `/object-storage/keys/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`

### Object Storage Clusters
- `/object-storage/clusters`
  - [X] `GET`
- `/object-storage/clusters/$id`
  - [X] `GET`

### Object Storage Buckets

- `/object-storage/buckets`
  - [X] `GET`
  - [X] `POST`
- `/object-storage/buckets/$id/$id`
  - [X] `GET`
  - [X] `DELETE`

### Payments

- `/account/payments`
  - [X] `GET`
  - [X] `POST`
- `/account/payments/$id`
  - [X] `GET`
- `/account/payments/paypal`
  - [ ] `GET`
- `/account/payments/paypal/execute`
  - [ ] `POST`

### Settings

- `/account/settings`
  - [X] `GET`
  - [X] `PUT`

### Users

- `/account/users`
  - [X] `GET`
  - [X] `POST`
- `/account/users/$username`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/account/users/$username/grants`
  - [ ] `GET`
  - [ ] `PUT`
- `/account/users/$username/password`
  - [ ] `POST`

## Profile

### Personalized User Settings

- `/profile`
  - [X] `GET`
  - [X] `PUT`

### Granted OAuth Apps

- `/profile/apps`
  - [ ] `GET`
- `/profile/apps/$id`
  - [ ] `GET`
  - [ ] `DELETE`

### Grants to Linode Resources

- `/profile/grants`
  - [ ] `GET`

### SSH Keys

- `/profile/sshkeys`
  - [x] `GET`
  - [x] `POST`
- `/profile/sshkeys/$id`
  - [x] `GET`
  - [x] `PUT`
  - [x] `DELETE`

### Two-Factor

- `/profile/tfa-disable`
  - [ ] `POST`
- `/profile/tfa-enable`
  - [ ] `POST`
- `/profile/tfa-enable-confirm`
  - [ ] `POST`

### Personal Access API Tokens

- `/profile/tokens`
  - [X] `GET`
  - [X] `POST`
- `/profile/tokens/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`

## Images

- `/images`
  - [x] `GET`
- `/images/$id`
  - [x] `GET`
  - [X] `POST`
  - [X] `PUT`
  - [X] `DELETE`

## Volumes

- `/volumes`
  - [X] `GET`
  - [X] `POST`
- `/volumes/$id`
  - [X] `GET`
  - [X] `PUT`
  - [X] `DELETE`
- `/volumes/$id/attach`
  - [X] `POST`
- `/volumes/$id/clone`
  - [X] `POST`
- `/volumes/$id/detach`
  - [X] `POST`
- `/volumes/$id/resize`
  - [X] `POST`
