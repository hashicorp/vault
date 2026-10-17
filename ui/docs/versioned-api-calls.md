# Versioned API Calls

When the Vault UI runs against an older Vault binary, some API methods introduced in newer versions won't exist in the generated client. Because the UI is a single shared codebase across all Vault versions, this can cause TypeScript compile errors in CI for older release branches. This document describes the system for writing version-gated API calls that compile cleanly on all supported branches while remaining fully type-safe on the latest branch.

---

## How It Works

The system has four pieces:

### 1. `FEATURE_VERSIONS` — the feature registry

`ui/app/utils/feature-versions.ts` is the single source of truth for all version-gated features. Each entry maps a stable key to a minimum Vault version and a human-readable name:

```typescript
// ui/app/utils/feature-versions.ts
export const FEATURE_VERSIONS: VersionedFeature[] = [
  { key: 'agents', name: 'Agentic Security', version: '2.1.0' },
];
```

If a feature's minimum version changes (e.g. it is backported), update the `version` field here. Call sites automatically pick up the change.

### 2. `this.version.hasFeature(key)` — the runtime gate

`VersionService.hasFeature(key)` looks up the key in `FEATURE_VERSIONS` and returns `true` when the running Vault binary meets the minimum version. Use this to guard both UI display logic and API calls:

```typescript
if (this.version.hasFeature('agents')) {
  // safe to run — Vault version is sufficient
}
```

### 3. `this.api.versioned()` — the type escape hatch

Inside a `hasFeature` gate you may call API methods that only exist in newer Vault versions. Because the TypeScript client is generated from the current binary, those methods won't exist in the type definitions on older branches. `ApiService.versioned()` returns the four API instances (`auth`, `identity`, `secrets`, `sys`) typed as `any`, silencing the missing-method error on older branches without affecting the real types on the latest branch.

```typescript
if (this.version.hasFeature('agents')) {
  // On the latest branch: newAgentMethod is real and fully typed on the generated client.
  // On an older branch: the gate is false at runtime, so this line is never reached;
  // the `any` type silences the TypeScript error at compile time.
  const result = await this.api.versioned().sys.newAgentMethod(params);
}
```

> **Important:** `versioned()` must **only** be called inside a `hasFeature` guard. The ESLint rule `vault/require-version-guard` enforces this and will fail CI if it is used elsewhere.

### 4. ESLint rule — `vault/require-version-guard`

A custom lint rule that errors if `versioned()` is called outside a `hasFeature` guard. This makes the pattern machine-enforced rather than just a convention.

---

## Step-by-Step Guide

Follow these steps every time you add a new version-gated feature that calls a newer API method.

### Step 1 — Register the feature

Add one entry to `FEATURE_VERSIONS` in `ui/app/utils/feature-versions.ts`:

```typescript
{ key: 'my-feature', name: 'My Feature Display Name', version: '2.1.0' },
```

- `key`: Stable kebab-case identifier. **Never rename it** — existing call sites reference it by string.
- `name`: Human-readable name for error messages and UI copy.
- `version`: The minimum Vault version required. No `v` prefix, no `+ent` suffix.

### Step 2 — Gate the UI and API call

In your component, route, or service:

```typescript
// ✅ Correct — gate the entire feature block
if (this.version.hasFeature('my-feature')) {
  const result = await this.api.versioned().sys.myNewMethod(params);
  // ... handle result
}
```

```typescript
// ❌ Wrong — versioned() outside a hasFeature guard (lint error)
const result = await this.api.versioned().sys.myNewMethod(params);
```

```typescript
// ❌ Wrong — wrong guard function (lint error)
if (this.version.isEnterprise) {
  const result = await this.api.versioned().sys.myNewMethod(params);
}
```

### Step 3 — Work normally on the latest branch

On the branch where the feature is being developed, `myNewMethod` exists in the real generated client. You get full TypeScript type checking, autocomplete, and parameter validation as normal. The `versioned()` call is just a pass-through to `any` — it does not widen the real types on the branch where the method exists.

---

## Displaying Version-Specific Error Messages

When you need to tell a user their Vault version doesn't support a feature, use `getFeatureVersion` to pull the metadata:

```typescript
import { getFeatureVersion } from 'vault/utils/feature-versions';

// In a component or route
if (!this.version.hasFeature('my-feature')) {
  const feature = getFeatureVersion('my-feature');
  // feature.name    → "My Feature Display Name"
  // feature.version → "2.1.0"
  this.flashMessages.danger(`${feature.name} requires Vault ${feature.version} or later.`);
  return;
}
```

---

## FAQ

**Q: I'm on an older branch and the TypeScript error is not inside a `versioned()` call — it's a type mismatch on a parameter or return value.**

A: The API method signature may have changed between versions (not just added). In this case use `versioned()` to call the method and handle the return value with an explicit type assertion. Keep the assertion as narrow as possible and document the version difference in a comment.

**Q: Can I call `versioned()` more than once in the same gate?**

A: Yes. You can also capture it once:

```typescript
if (this.version.hasFeature('my-feature')) {
  const api = this.api.versioned();
  const a = await api.sys.methodA();
  const b = await api.secrets.methodB();
}
```

**Q: What if the feature is available in enterprise but not community on the same version?**

A: Combine `hasFeature` with the existing `isEnterprise` check:

```typescript
if (this.version.isEnterprise && this.version.hasFeature('my-feature')) {
  await this.api.versioned().sys.myNewMethod();
}
```

The ESLint rule recognises `hasFeature` anywhere in the `if` condition, including inside `&&` expressions.

**Q: A feature was backported to an older version. What do I change?**

A: Update the `version` field in `FEATURE_VERSIONS`. That's it. No call sites need to change.

**Q: My new method doesn't exist on the current branch yet — how do I write and test the code?**

A: Write the feature on the branch where the method exists (the latest branch). The method will be present in the generated client there, so you get full type safety. When the code is merged back to older branches, the `versioned()` call silences the compile error and the `hasFeature` gate prevents the code from running at runtime.
