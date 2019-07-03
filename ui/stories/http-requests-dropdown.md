<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/http-requests-dropdown.js. To make changes, first edit that file and run "yarn gen-story-md http-requests-dropdown" to re-generate the content.-->

## HttpRequestsDropdown
HttpRequestsDropdown components are used to render a dropdown that filters the HttpRequestsBarChart.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint. |

**Example**
  
```js
<HttpRequestsDropdown @counters={counters} />
```

**See**

- [Uses of HttpRequestsDropdown](https://github.com/hashicorp/vault/search?l=Handlebars&q=HttpRequestsDropdown+OR+http-requests-dropdown)
- [HttpRequestsDropdown Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/http-requests-dropdown.js)

---
