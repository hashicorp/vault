/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

// For the purposes of testing, we only use a subset of fields relevant to mysql
export default Factory.extend({
  backend: 'database',
  name: 'connection',
  plugin_name: 'mysql-database-plugin',
  verify_connection: true,
  connection_url: '{{username}}:{{password}}@tcp(127.0.0.1:33060)/',
  username: 'admin',
  max_open_connections: 4,
  max_idle_connections: 0,
  max_connection_lifetime: '0s',
  allowed_roles: () => [],
  root_rotation_statements: () => [
    'SELECT user from mysql.user',
    "GRANT ALL PRIVILEGES ON *.* to 'sudo'@'%'",
  ],
});
