<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/http-requests-bar-chart.js. To make changes, first edit that file and run "yarn gen-story-md http-requests-bar-chart" to re-generate the content.-->

## HttpRequestsBarChart
`HttpRequestsBarChart` components are used to render a bar chart with the total number of HTTP Requests to a Vault server per month.


| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests` endpoint. |

**Example**
  
```js
<HttpRequestsBarChart @counters={{counters}} />
```

**See**

- [Uses of HttpRequestsBarChart](https://github.com/hashicorp/vault/search?l=Handlebars&q=HttpRequestsBarChart+OR+http-requests-bar-chart)
- [HttpRequestsBarChart Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/http-requests-bar-chart.js)

---
