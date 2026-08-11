/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { v4 as uuidv4 } from 'uuid';
import { subYears } from 'date-fns';

// Static OAuth resource server profiles keyed by config_id.
// The mount_accessor for an OAuth alias follows the format:
//   oauth-resource-server_<namespace_id>_<config_id>
const OAUTH_PROFILES = {
  'b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03': {
    config_id: 'b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03',
    profile_name: 'github-actions',
    issuer_id: 'https://token.actions.githubusercontent.com',
    use_jwks: true,
    jwks_uri: 'https://token.actions.githubusercontent.com/.well-known/jwks',
    audiences: ['https://github.com/my-org'],
    no_default_policy: false,
    optional_authorization_details: false,
    user_claim: 'sub',
    supported_algorithms: ['RS256'],
    jwt_type: 'access_token',
    clock_skew_leeway: 30,
    enabled: true,
  },
};

// Namespace ID used in the OAuth mount_accessor for mirage data.
const OAUTH_NAMESPACE_ID = 'root';
// Config ID of the OAuth profile linked to the mirage OAuth alias.
const OAUTH_CONFIG_ID = 'b2f5c891-3a7d-4e12-9f8a-1c6d4e7b2a03';

const generateDates = () => {
  const now = new Date();
  const start = subYears(now, 5).getTime();

  return [0, 0].map(() => new Date(start + Math.random() * (now - start))).sort((a, b) => a - b);
};

export default function (server) {
  server.get('agent-registry/registration/id', (schema, request) => {
    if (request.queryParams.list) {
      // Get all registrations from the database
      const registrations = schema.db.agentRegistryRegistrations.where({});

      // Build the response in the required format
      // keys array contains ID values (not display names)
      // key_info maps id -> { display_name }
      const data = registrations.reduce(
        (resp, record) => {
          resp.keys.push(record.id);
          resp.key_info[record.id] = {};
          // not all properties are returned in the list response
          const keys = [
            'display_name',
            'entity_id',
            'creation_time',
            'owner',
            'last_updated_time',
            'no_default_ceiling_policy',
          ];
          keys.forEach((key) => {
            resp.key_info[record.id][key] = record[key];
          });
          return resp;
        },
        {
          keys: [],
          key_info: {},
        }
      );

      return { data };
    }
  });

  server.get('agent-registry/registration/id/:id', (schema, request) => {
    const { id } = request.params;

    // Find registration by id
    const registration = schema.db.agentRegistryRegistrations.find(id);

    if (!registration) {
      return new Response(404, {}, { errors: [`Registration with id ${id} not found`] });
    }

    // Return the full registration details
    return {
      data: registration,
    };
  });

  server.delete('agent-registry/registration/id/:id', (schema, request) => {
    const { id } = request.params;

    // Find registration by id
    const registration = schema.db.agentRegistryRegistrations.find(id);

    if (!registration) {
      return new Response(404, {}, { errors: [`Registration with id ${id} not found`] });
    }

    // Delete the registration
    schema.db.agentRegistryRegistrations.remove(id);

    // Return empty response (204 No Content)
    return new Response(204);
  });

  server.post('agent-registry/registration/id/:id', (schema, request) => {
    const { id } = request.params;
    const data = JSON.parse(request.requestBody);

    // Check if registration already exists
    const existingRegistration = schema.db.agentRegistryRegistrations.find(id);

    if (existingRegistration) {
      // Update existing registration
      schema.db.agentRegistryRegistrations.update(id, {
        ...data,
        last_updated_time: new Date().toISOString(),
      });

      return {
        data: {
          id: existingRegistration.id,
          display_name: data.display_name || existingRegistration.display_name,
        },
      };
    } else {
      // Create new registration with specified ID
      const newRegistration = server.create('agent-registry-registration', {
        ...data,
        id,
      });

      return {
        data: {
          id: newRegistration.id,
          display_name: newRegistration.display_name,
        },
      };
    }
  });

  server.get('agent-registry/registration/entity-id/:entity_id', (schema, request) => {
    const { entity_id } = request.params;

    // Find registration by entity_id
    const registration = schema.db.agentRegistryRegistrations.findBy({ entity_id });

    if (!registration) {
      return new Response(404, {}, { errors: [`Registration with entity_id ${entity_id} not found`] });
    }

    // Return the full registration details
    return {
      data: {
        id: registration.id,
        display_name: registration.display_name,
        entity_id: registration.entity_id,
        description: registration.description,
        owner: registration.owner,
        ceiling_policy: registration.ceiling_policy,
        no_default_ceiling_policy: registration.no_default_ceiling_policy,
        creation_time: registration.creation_time,
        last_updated_time: registration.last_updated_time,
      },
    };
  });

  server.get('identity/entity/id', (schema, request) => {
    if (request.queryParams.list) {
      const registrations = schema.db.agentRegistryRegistrations.where({});

      const data = registrations.reduce(
        (resp, registration, index) => {
          const { entity_id: id } = registration;

          if (resp.key_info[id]) {
            return resp;
          }

          // select random number between 0 and 3 as alias count
          const aliasCount = Math.floor(Math.random() * 4);
          const aliases = Array.from({ length: aliasCount }, (_, index) => {
            const aliasNumber = index + 1;
            const isTokenAlias = aliasNumber % 2 === 1;
            const [creation_time, last_update_time] = generateDates();

            return {
              id: uuidv4(),
              name: `alias-${index}`,
              creation_time,
              last_update_time,
              mount_accessor: isTokenAlias ? 'auth_token_d2eccb3b' : 'auth_userpass_a212c204',
              mount_path: isTokenAlias ? 'auth/token/' : 'auth/userpass/',
              mount_type: isTokenAlias ? 'token' : 'userpass',
            };
          });

          const [creation_time, last_update_time] = generateDates();
          resp.key_info[id] = {
            aliases,
            name: `entity-${index}`,
            disabled: Math.random() < 0.5,
            creation_time,
            last_update_time,
          };
          resp.keys.push(id);

          return resp;
        },
        {
          keys: [],
          key_info: {},
        }
      );

      return { data };
    }
  });

  server.get('identity/entity/id/:id', (schema, request) => {
    const { id } = request.params;
    const registrations = schema.db.agentRegistryRegistrations.where({});
    const registrationIndex = registrations.findIndex((registration) => registration.entity_id === id);

    if (registrationIndex === -1) {
      return new Response(404, {}, { errors: [`Entity with id ${id} not found`] });
    }

    const regularAliasCount = Math.floor(Math.random() * 3);
    const aliases = Array.from({ length: regularAliasCount }, (_, index) => {
      const aliasNumber = index + 1;
      const isTokenAlias = aliasNumber % 2 === 1;
      const [creation_time, last_update_time] = generateDates();

      return {
        id: uuidv4(),
        name: `alias-${index}`,
        creation_time,
        last_update_time,
        custom_metadata: {
          contact_email: 'james@example.com',
          alias_metadata: `alias-metadata-${index}`,
        },
        mount_accessor: isTokenAlias ? 'auth_token_d2eccb3b' : 'auth_userpass_a212c204',
        mount_path: isTokenAlias ? 'auth/token/' : 'auth/userpass/',
        mount_type: isTokenAlias ? 'token' : 'userpass',
      };
    });

    // Include one OAuth resource server alias so the profile enrichment code path is exercised.
    const [oauthCreationTime, oauthLastUpdateTime] = generateDates();
    aliases.push({
      id: uuidv4(),
      name: `repo:my-org/my-repo:ref:refs/heads/main`,
      creation_time: oauthCreationTime,
      last_update_time: oauthLastUpdateTime,
      custom_metadata: null,
      mount_accessor: `oauth-resource-server_${OAUTH_NAMESPACE_ID}_${OAUTH_CONFIG_ID}`,
      mount_path: `auth/oauth-resource-server/`,
      mount_type: 'oauth-resource-server',
    });

    const [creation_time, last_update_time] = generateDates();
    const mergedEntityIds =
      Math.random() < 0.5 ? null : [`entity-${registrationIndex + 1}`, `entity-${registrationIndex + 2}`];
    const groupCount = Math.floor(Math.random() * 4);
    const group_ids = Array.from({ length: groupCount }, () => uuidv4());
    const policyCount = Math.floor(Math.random() * 4);
    const policies = Array.from(
      { length: policyCount },
      (_, index) => `policy-${registrationIndex}-${index + 1}`
    );

    return {
      data: {
        id,
        aliases,
        name: `entity-${registrationIndex}`,
        disabled: Math.random() < 0.5,
        creation_time,
        last_update_time,
        group_ids,
        policies,
        metadata: {
          organization: 'hashicorp',
          team: `team-${registrationIndex}`,
        },
        merged_entity_ids: mergedEntityIds,
      },
    };
  });

  server.get('/identity/group/id/:id', (_schema, request) => {
    const { id } = request.params;

    return {
      data: {
        alias: {},
        creation_time: '2017-11-13T19:36:47.102945Z',
        id,
        last_update_time: '2017-11-13T19:36:47.102945Z',
        member_entity_ids: [],
        member_group_ids: null,
        metadata: {
          hello: 'world',
        },
        modify_index: 1,
        name: 'group_ab813d63',
        policies: ['grouppolicy1', 'grouppolicy2'],
        type: 'internal',
      },
    };
  });

  server.get('sys/policies/acl/:name', (_schema, request) => {
    const { name } = request.params;
    const policyString = `path "secret/data/svcY/*" {
      capabilities = ["delete", "read"]
    }

    path "kv/data/appB/*" {
      capabilities = ["update", "create", "list", "delete", "read"]
    }

    path "identity/*" {
      capabilities = ["read"]
    }`;

    return {
      data: {
        name,
        policy: policyString,
      },
    };
  });

  server.get('agent-registry/registration/display-name', (schema, request) => {
    if (request.queryParams.list) {
      // Get all registrations from the database
      const registrations = schema.db.agentRegistryRegistrations.where({});

      // Build the response in the required format
      // keys array contains display_name values
      // key_info maps id -> { display_name }
      const data = registrations.reduce(
        (resp, record) => {
          resp.keys.push(record.display_name);
          resp.key_info[record.id] = {
            display_name: record.display_name,
            entity_id: record.entity_id,
          };
          return resp;
        },
        {
          keys: [],
          key_info: {},
        }
      );

      return { data };
    }
  });

  server.get('agent-registry/registration/display-name/:name', (schema, request) => {
    const { name } = request.params;

    // Find registration by display_name
    const registration = schema.db.agentRegistryRegistrations.findBy({ display_name: name });

    if (!registration) {
      return new Response(404, {}, { errors: [`Registration with name ${name} not found`] });
    }

    // Return the full registration details
    return {
      data: {
        id: registration.id,
        display_name: registration.display_name,
        entity_id: registration.entity_id,
        description: registration.description,
        owner: registration.owner,
        ceiling_policy: registration.ceiling_policy,
        no_default_ceiling_policy: registration.no_default_ceiling_policy,
        creation_time: registration.creation_time,
        last_updated_time: registration.last_updated_time,
      },
    };
  });

  server.post('agent-registry/registration/display-name/:name', (schema, request) => {
    const { name } = request.params;
    const data = JSON.parse(request.requestBody);

    // Check if registration already exists
    const existingRegistration = schema.db.agentRegistryRegistrations.findBy({ display_name: name });

    if (existingRegistration) {
      // Update existing registration
      schema.db.agentRegistryRegistrations.update(existingRegistration.id, {
        ...data,
        display_name: name,
        last_updated_time: new Date().toISOString(),
      });

      return {
        data: {
          id: existingRegistration.id,
          display_name: name,
        },
      };
    } else {
      // Create new registration
      const newRegistration = server.create('agent-registry-registration', {
        ...data,
        display_name: name,
      });

      return {
        data: {
          id: newRegistration.id,
          display_name: newRegistration.display_name,
        },
      };
    }
  });

  server.delete('agent-registry/registration/display-name/:name', (schema, request) => {
    const { name } = request.params;

    // Find registration by display_name
    const registration = schema.db.agentRegistryRegistrations.findBy({ display_name: name });

    if (!registration) {
      return new Response(404, {}, { errors: [`Registration with name ${name} not found`] });
    }

    // Delete the registration
    schema.db.agentRegistryRegistrations.remove(registration.id);

    // Return empty response (204 No Content)
    return new Response(204);
  });
  server.get('sys/config/oauth-resource-server/id/:id', (_schema, request) => {
    const { id } = request.params;
    const profile = OAUTH_PROFILES[id];

    if (!profile) {
      return new Response(404, {}, { errors: [`OAuth resource server profile with id ${id} not found`] });
    }

    return { data: profile };
  });
}
