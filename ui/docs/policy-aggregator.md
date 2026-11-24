# Policy Aggregator Utility

A utility for aggregating multiple Vault policies by combining capabilities for duplicate paths without duplication.

## Usage

```typescript
import { aggregatePolicies } from 'vault/utils/policy-aggregator';

const policy1 = `
# Comment example
path "/foo" {
    capabilities = ["create"]
}
path "/foo/bar" {
    capabilities = ["create", "read", "update"]
}
`;

const policy2 = `
path "/foo" {
    capabilities = ["delete"]
}
path "/bar" {
    capabilities = ["patch", "delete"]
}
`;

const result = aggregatePolicies([policy1, policy2]);
```

## Output

The function returns an object with two properties:

### `policy` Object

An object with paths as keys and arrays of capabilities as values:

```typescript
{
  "/foo": ["create", "delete"],
  "/foo/bar": ["create", "read", "update"],
  "/bar": ["delete", "patch"]
}
```

### `policyString` String

A formatted HCL string matching the input format:

```hcl
path "/foo" {
    capabilities = ["create", "delete"]
}
path "/foo/bar" {
    capabilities = ["create", "read", "update"]
}
path "/bar" {
    capabilities = ["delete", "patch"]
}
```

## Features

- **Deduplication**: Automatically removes duplicate paths and capabilities
- **Merging**: Combines capabilities when the same path exists across multiple policies
- **Sorting**: Capabilities are sorted alphabetically for consistent output
- **Comment Handling**: Ignores comments in policy strings
- **Format Preservation**: Output string matches the input HCL format

## API

### `aggregatePolicies(policyStrings: string[]): AggregatePolicy`

**Parameters:**

- `policyStrings`: Array of HCL-formatted Vault policy strings

**Returns:**

```typescript
interface AggregatePolicy {
  policy: {
    [path: string]: string[];
  };
  policyString: string;
}
```

## Example Use Cases

### Combining User Policies

When a user has multiple policies assigned, aggregate them to see their effective permissions:

```typescript
const userPolicies = getUserPolicies(userId);
const effective = aggregatePolicies(userPolicies);
console.log('Effective permissions:', effective.policy);
```

### Policy Analysis

Analyze which paths have overlapping capabilities across different policies:

```typescript
const result = aggregatePolicies(allPolicies);
Object.entries(result.policy).forEach(([path, caps]) => {
  if (caps.length > 3) {
    console.log(`Path ${path} has many capabilities:`, caps);
  }
});
```

### Generating Merged Policy

Create a single merged policy from multiple sources:

```typescript
const merged = aggregatePolicies([adminPolicy, devPolicy, readPolicy]);
// Use merged.policyString to create a new policy
```
