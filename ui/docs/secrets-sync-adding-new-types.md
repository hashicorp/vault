# Adding New Secret Types to Secrets Sync

This document explains how to add support for new secret types to the Secrets Sync feature.

## Overview

The Secrets Sync feature uses a mount-based type detection system where selecting a secrets engine mount automatically determines the secret type. The system is designed to be extensible by updating configuration files rather than modifying core component logic.

## Current Secret Types

- **KV** - KV version 2 secrets (requires KV v2, not v1)
- **Database** - Database static role credentials

## Architecture

The sync feature uses:
- **Mount-based detection**: Secret type is automatically determined from the selected mount's engine type
- **HDS SuperSelect**: Dropdown showing all supported mounts with icons and type indicators
- **Configuration-driven**: All secret type behavior defined in `SECRET_TYPE_CONFIGS` and `SECRET_TYPE_FETCHERS`
- **Helper functions**: `getSecretTypeFromMount()` and `getSecretTypeFromAccessor()` handle type detection

## Steps to Add a New Secret Type

### 1. Update Type Definition

In `/ui/types/vault/sync.d.ts`, add your new secret type to the `SecretType` union:

```typescript
export type SecretType = 'kv' | 'database' | 'your-new-type';
```

### 2. Update Type Detection Functions

There are two detection functions in `secret-type-config.ts`:

**`getSecretTypeFromMount(type, version)`** — used when a user selects a mount from the dropdown. If your type has special requirements (like KV requiring v2), handle them here:

```typescript
export function getSecretTypeFromMount(type: string, version?: number): SecretType | null {
  if (type === 'kv') {
    return version === 2 ? 'kv' : null; // KV requires version 2
  }
  if (type === 'database') {
    return 'database';
  }
  // Add your new type — add a version check if needed
  if (type === 'your-engine-type') {
    return 'your-new-type';
  }
  return null;
}
```

**`getSecretTypeFromAccessor(accessor)`** — used when displaying already-synced secrets (via the `accessor` field returned by the API, e.g. `kv_9b39cc0f`). Add a prefix match for your new type:

```typescript
export function getSecretTypeFromAccessor(accessor: string): SecretType | null {
  const prefix = accessor.split('_')[0];
  if (prefix === 'kv' || prefix === 'database' || prefix === 'your-engine-type') return prefix;
  return null;
}
```

**Note**: `fetchMounts()` in `sync.ts` uses `getSecretTypeFromMount()` to filter which mounts appear in the dropdown. Your engine type will automatically appear once this function returns a non-null value for it.

### 3. Add Secret Type Configuration

In `/ui/lib/sync/addon/utils/secret-type-config.ts`, add a new entry to `SECRET_TYPE_CONFIGS`:

```typescript
export const SECRET_TYPE_CONFIGS: Record<SecretType, SecretTypeConfig> = {
  // ... existing configs
  'your-new-type': {
    placeholder: 'Input placeholder text',
    noMatchesMessage: 'Message when no suggestions found',
    accessorType: 'Your Engine Display Name', // Shown in the type badge
    icon: 'your-hds-icon-name', // HDS icon name for this engine type
    route: 'externalRouteName', // External route name for viewing the secret
    supportsExternalLink: true, // Whether to show a link to the secret in the success message
    getModels: (mount, secretName) => {
      // Return array of route parameters
      // Example: [mount, secretName] or [mount, `role/${secretName}`]
      return [mount, secretName];
    },
    getQuery: () => {
      // Optional: Return query parameters for the route
      // Example: { type: 'static' }
      // Omit this field if no query params are needed
      return undefined;
    },
  },
};
```

### 4. Add Fetcher Configuration

In `/ui/lib/core/addon/utils/secret-type-fetchers.ts`, add a new entry to `SECRET_TYPE_FETCHERS`:

```typescript
export const SECRET_TYPE_FETCHERS: Record<SecretType, SecretTypeFetcher> = {
  // ... existing fetchers
  'your-new-type': {
    fetch: async (api, mountPath, value) => {
      // Fetch suggestions from API
      // Example:
      const backend = mountPath.endsWith('/') ? mountPath.slice(0, -1) : mountPath;
      const { keys } = await api.yourService.listItems(backend);
      return keys || [];
    },
    filter: (items, value, isDirectory) => {
      // Filter items based on current input value
      if (!value) return items;
      return items.filter((item) => item.toLowerCase().includes(value.toLowerCase()));
    },
    onSelect: (item, pathToSecret) => {
      // Return the final value when an item is selected
      // For simple cases: return item;
      // For hierarchical (like KV): return `${pathToSecret}${item}`;
      return item;
    },
  },
};
```

### 5. Register External Route

If `supportsExternalLink: true`, register the route in `/ui/app/app.js`:

```javascript
const externalRoutes = {
  kvSecretOverview: 'vault.cluster.secrets.backend.kv.secret.index',
  databaseStaticRoleOverview: 'vault.cluster.secrets.backend.show',
  yourNewTypeOverview: 'vault.cluster.your.route.path', // Add your route
};
```

### 6. Update Engine Dependencies (if needed)

In `/ui/lib/sync/addon/engine.js`, add the new external route to the array:

```javascript
const externalRoutes = [
  'kvSecretOverview',
  'databaseStaticRoleOverview',
  'yourNewTypeOverview', // Add your route
];
```

### 7. Update Mount Fetching Logic

In `/ui/lib/sync/addon/components/secrets/page/destinations/destination/sync.ts`, update the `fetchMounts()` method to include your new engine type:

```typescript
async fetchMounts() {
  const supportedMounts: MountOption[] = [];

  try {
    const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts();
    if (secret) {
      for (const path in secret) {
        const { type, options } = secret[path];
        const version = options?.version ? Number(options.version) : undefined;

        // Check if this mount type is supported for sync
        const secretType = getSecretTypeFromMount(type, version);
        if (!secretType) continue;

        supportedMounts.push({
          name: path,
          id: path,
          engineType: type,
          version,
        });
      }
    }
    this.allSupportedMounts = supportedMounts;
  } catch (error) {
    // Handle error
  }
}
```

**Note**: The `MountOption` interface is defined in `secret-type-config.ts`:

```typescript
export interface MountOption {
  name: string;
  id: string;
  engineType: string;
  version?: number;
}
```

The icon and display name for each type are read from the corresponding entry in `SECRET_TYPE_CONFIGS` (the `icon` and `accessorType` fields).

## Testing Your New Secret Type

After adding the configuration:

1. The new engine type will automatically appear in the HDS SuperSelect mount dropdown (if mounts of that type exist)
2. Selecting a mount of your engine type will automatically detect the secret type
3. The mount dropdown will display with the correct icon and type indicator
4. The SuggestionInput component will use your fetch/filter logic
5. The secrets list will use your route and model configuration
6. All getters and conditionals use the configuration automatically

**Important**: The system uses mount-based type detection. Ensure your engine type is handled in `getSecretTypeFromMount()` (for new sync) and `getSecretTypeFromAccessor()` (for existing synced secrets).

## Example: Adding PKI Certificate Type

Here's a complete example of adding support for PKI certificates:

### 1. Type Definition

```typescript
// In /ui/types/vault/sync.d.ts
export type SecretType = 'kv' | 'database' | 'pki';
```

### 2. Type Detection

```typescript
// getSecretTypeFromMount — for new sync
if (type === 'pki') return 'pki';

// getSecretTypeFromAccessor — for existing synced secrets
if (prefix === 'pki') return 'pki';
```

### 3. Configuration

```typescript
// Add config to SECRET_TYPE_CONFIGS
pki: {
  placeholder: 'Certificate serial number',
  noMatchesMessage: 'No certificates found',
  accessorType: 'PKI',
  icon: 'certificate',
  route: 'pkiCertificateOverview',
  supportsExternalLink: true,
  getModels: (mount, serialNumber) => [mount, serialNumber],
},
```

### 4. Fetcher

```typescript
// Add fetcher to SECRET_TYPE_FETCHERS
pki: {
  fetch: async (api, mountPath) => {
    const backend = mountPath.endsWith('/') ? mountPath.slice(0, -1) : mountPath;
    const { keys } = await api.secrets.pkiListCerts(backend);
    return keys || [];
  },
  filter: (certs, value) => {
    if (!value) return certs;
    return certs.filter(c => c.toLowerCase().includes(value.toLowerCase()));
  },
  onSelect: (item) => item,
},
```
## Benefits of This Architecture

- **Mount-Based Detection**: Users select a secrets engine mount; the system automatically determines the secret type
- **Single Source of Truth**: All secret type behavior defined in `SECRET_TYPE_CONFIGS` and `SECRET_TYPE_FETCHERS`
- **No Component Changes**: Core components (`SuggestionInput`, `sync.ts`) require no modification when adding new types
- **Clear Engine Boundaries**: `core` components receive configuration as arguments; they do not import from engine-specific packages like `sync`
- **Type Safe**: TypeScript ensures all required config fields are provided when adding a new `SecretType`
- **Consistent UI**: Icons and type indicators are driven by the `icon` and `accessorType` fields in `SECRET_TYPE_CONFIGS`, shown via `SecretTypeBadge`
- **Easy Testing**: Each secret type config can be tested independently
- **Extensible**: Adding a new type requires changes only to `secret-type-config.ts` and `secret-type-fetchers.ts`

1. User navigates to the sync page for a destination
2. User selects a secrets engine mount from the HDS SuperSelect dropdown
   - Dropdown shows all supported mounts with icons and type indicators (e.g., `secret/ (KV v2)`)
   - Only mounts supported for sync appear (determined by `getSecretTypeFromMount()`)
3. System automatically detects the secret type based on the mount's engine type
4. `SuggestionInput` appears with appropriate placeholder and help text for that secret type
5. User enters or selects a secret/role to sync
6. On successful sync, a link to the secret is shown (if `supportsExternalLink` is true)

## Key Files

| File | Purpose |
|------|---------|
| `/ui/types/vault/sync.d.ts` | `SecretType` union and shared type definitions |
| `/ui/lib/sync/addon/utils/secret-type-config.ts` | `SECRET_TYPE_CONFIGS`, `getSecretTypeFromMount()`, `getSecretTypeFromAccessor()` |
| `/ui/lib/core/addon/utils/secret-type-fetchers.ts` | `SECRET_TYPE_FETCHERS` — API fetch, filter, and select logic |
| `/ui/lib/sync/addon/components/secrets/page/destinations/destination/sync.ts` | Main sync component — mount fetching and type detection |
| `/ui/lib/sync/addon/components/secrets/page/destinations/destination/sync.hbs` | Template with HDS SuperSelect and `SuggestionInput` |
| `/ui/lib/core/addon/components/suggestion-input.ts` | Generic suggestion input component |
| `/ui/lib/sync/addon/components/secrets/secret-type-badge.ts` | Reusable badge showing engine type with icon |
