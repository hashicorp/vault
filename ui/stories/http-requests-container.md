<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/http-requests-container.js. To make changes, first edit that file and run "yarn gen-story-md http-requests-container" to re-generate the content.-->

## HttpRequestsContainer
The HttpRequestsContainer component is the parent component of the HttpRequestsDropdown, HttpRequestsBarChart, and HttpRequestsTable components. It is used to handle filtering the bar chart and table according to selected time window from the dropdown.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint. |

**Example**
  
```js
<HttpRequestsContainer @counters={counters}/>
```

**See**

- [Uses of HttpRequestsContainer](https://github.com/hashicorp/vault/search?l=Handlebars&q=HttpRequestsContainer+OR+http-requests-container)
- [HttpRequestsContainer Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/http-requests-container.js)

---
