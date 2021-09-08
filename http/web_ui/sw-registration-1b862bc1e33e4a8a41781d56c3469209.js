!function(){"use strict"
var e,n=[],r=[]
"serviceWorker"in navigator&&navigator.serviceWorker.register("/ui/sw.js",{scope:"/v1/sys/storage/raft/snapshot"}).then((function(e){for(var r=Promise.resolve(),o=function(o,t){r=r.then((function(){return n[o](e)}))},t=0,s=n.length;t<s;t++)o(t)
return r.then((function(){console.log("Service Worker registration succeeded. Scope is "+e.scope)}))})).catch((function(e){for(var n=Promise.resolve(),o=function(o,t){n=n.then((function(){return r[o](e)}))},t=0,s=r.length;t<s;t++)o(t)
return n.then((function(){console.log("Service Worker registration failed with "+e)}))})),e=function(e){navigator.serviceWorker.addEventListener("message",(function(e){var n=e.data.action,r=e.ports[0]
if("getToken"===n){var o=Ember.Namespace.NAMESPACES_BY_ID.vault.__container__.lookup("service:auth").currentToken
o||console.error("Unable to retrieve Vault tokent"),r.postMessage({token:o})}else console.error("Unknown event",e),r.postMessage({error:"Unknown request"})})),window.addEventListener("unload",(function(){e.unregister()}))},n.push(e)}()
