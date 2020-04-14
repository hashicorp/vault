<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/replication/addon/components/replication-primary-card.js. To make changes, first edit that file and run "yarn gen-story-md replication-primary-card" to re-generate the content.-->

## ReplicationPrimaryCard
ReplicationPrimaryCard components

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| [title] | <code>string</code> | <code>null</code> | The title to be displayed on the top left corner of the card. |
| [description] | <code>string</code> | <code>null</code> | Helper text to describe the metric on the card. |
| metric | <code>string</code> | <code>null</code> | The main metric to highlight on the card. |

**Example**
  
```js
<ReplicationPrimaryCard
    @title='Last WAL entry'
    @description='Index of last Write Ahead Logs entry written on local storage.'
    @metric={{replicationAttrs.lastWAL}}
    />
```
    

**See**

- [Uses of ReplicationPrimaryCard](https://github.com/hashicorp/vault/search?l=Handlebars&q=ReplicationPrimaryCard+OR+replication-primary-card)
- [ReplicationPrimaryCard Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/replication/addon/components/replication-primary-card.js)

---
