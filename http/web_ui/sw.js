!function(){"use strict"
self.CACHE_BUSTER="1631132269789|0.9141803921560516",self.addEventListener("install",(function(e){return self.skipWaiting()})),self.addEventListener("message",(function(e){if("skipWaiting"===e.data)return self.skipWaiting()})),self.addEventListener("activate",(function(e){return self.clients.claim()}))
var e=["/v1/sys/storage/raft/snapshot"].map((function(e){var n=function(e){var n=arguments.length>1&&void 0!==arguments[1]?arguments[1]:self.location
return decodeURI(new URL(encodeURI(e),n).toString())}(e)
return new RegExp("^".concat(n,"$"))}))
self.addEventListener("fetch",(function(n){var t=n.request
return function(e,n){return!!n.find((function(n){return n.test(decodeURI(e))}))}(t.url,e)&&"GET"===t.method?n.respondWith(function(e){var n,t=new Headers(e.headers)
return(n={action:"getToken"},self.clients.matchAll({includeUncontrolled:!0,type:"window"}).then((function(e){var t=e[0]
return new Promise((function(e,r){var s=new MessageChannel
s.port2.onmessage=function(n){n.data.error?r(n.data.error):e(n.data.token)},t.postMessage(n,[s.port1])}))}))).then((function(n){return t.set("X-Vault-Token",n),fetch(new Request(e.url,{method:e.method,headers:t}))}))}(t)):n.respondWith(fetch(t))}))}()
