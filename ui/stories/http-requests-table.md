<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/http-requests-table.js. To make changes, first edit that file and run "yarn gen-story-md http-requests-table" to re-generate the content.-->

## HttpRequestsTable
`HttpRequestsTable` components render a table with the total number of HTTP Requests to a Vault server per month.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint. |

**Example**
  
```js
<HttpRequestsTable @counters={{counters}} />
```

**See**

- [Uses of HttpRequestsTable](https://github.com/hashicorp/vault/search?l=Handlebars&q=HttpRequestsTable+OR+http-requests-table)
- [HttpRequestsTable Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/http-requests-table.js)

---
