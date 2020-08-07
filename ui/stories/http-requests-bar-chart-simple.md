<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in app/components/http-requests-bar-chart-simple.js. To make changes, first edit that file and run "yarn gen-story-md http-requests-bar-chart-simple" to re-generate the content.-->

## HttpRequestsBarChartSimple
The HttpRequestsBarChartSimple is a simplified version of the HttpRequestsBarChart component.

**Params**

| Param | Type | Default | Description |
| --- | --- | --- | --- |
| counters | <code>Array</code> | <code></code> | A list of objects containing the total number of HTTP Requests for each month. `counters` should be the response from the `/internal/counters/requests`. The response is then filtered showing only the 12 most recent months of data.  This property is called filteredHttpsRequests. |

**Example**
  
```js
<HttpRequestsBarChartSimple @counters={counters}/>
```

**See**

- [Uses of HttpRequestsBarChartSimple](https://github.com/hashicorp/vault/search?l=Handlebars&q=HttpRequestsBarChartSimple+OR+http-requests-bar-chart-simple)
- [HttpRequestsBarChartSimple Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/http-requests-bar-chart-simple.js)

---
