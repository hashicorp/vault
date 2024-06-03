# Test Helpers Organization

Our test are constantly evolving, but here's a general overview of how we set up and organize tests.

## Folder organization

### /acceptance

Acceptance tests should test the overall workflows and navigation within Vault. When possible, they should use the real API instead of mocked so that breaking changes from the backend can be caught. Reasons you may opt to use a mocked backend instead of the real one:

- Using the real backend would cause instability in concurrently-running tests (eg. seal/unseal flow)
- There isn't a way to set up a 3rd party dependency that the backend needs to run correctly (Database Secrets Engine, Sync Secrets)

### /helpers

Shared helpers such as selectors, common interactions, WebREPL commands, and API response stubs live in this folder. When selectors are only used for a single test, they should be defined on the same file where they are used for the test. Once the selectors are being used for multiple tests, they should be moved to this folder so they can be defined in a single place and shared to wherever they are needed.

Often we will need a set of selectors for "workflow" tests, or acceptance tests that navigate through an area of the app. For these, the helpers should be organized as such:

- `/helpers/<area>/<area>-selectors.ts` - exports selector consts (never default) for each page -- eg. for PKI we would have PKI_OVERVIEW, PKI_ROLE, etc.
- `/helpers/<area>/<area>-helpers.js` - exports methods and consts which are otherwise helpful in the tests -- eg. example API responses, common interactions (eg. writeVersionedSecret for KV v2)

Whenever possible we should try to use the general selectors exported from `/helpers/general-selectors.ts`.

### integration

Integration tests are most often used to test specific components out of context from the rest of the app. Be sure to mock anything that the component needs to work correctly -- for example, if the component has a certain behavior in enterprise than community edition, in your tests for each scenario it should not assume that the underlying Vault binary is in one state or the other, and mock the enterprise/community state in all the scenarios. The exports in `helpers/stubs.js` might be helpful for these tests, particularly when the component uses a model which fetches capabilities.

### pages

[DEPRECATED] This file should be removed in favor of selectors within the "helpers" folder. We are moving away from ember-page-object selectors toward simple strings

### unit

Unit tests are most often used to test utils, adapters, serializers, routes, and services functionality.
