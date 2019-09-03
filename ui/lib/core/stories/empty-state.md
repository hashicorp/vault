<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/empty-state.js. To make changes, first edit that file and run "yarn gen-story-md empty-state" to re-generate the content.-->

## EmptyState
`EmptyState` components are used to render a helpful message and any necessary content when a user
encounters a state that would usually be blank.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| title | <code>String</code> | <code></code> | A short label for the empty state |
| message | <code>String</code> | <code></code> | A description of why a user might be seeing the empty state and possibly instructions for actions they may take. |

**Example**
  
```js
<EmptyState @title="You don't have an secrets yet" @message="An explanation of why you don't have any secrets but also you maybe want to create one." />
```

**See**

- [Uses of EmptyState](https://github.com/hashicorp/vault/search?l=Handlebars&q=EmptyState)
- [EmptyState Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/empty-state.js)

---
