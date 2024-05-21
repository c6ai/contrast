"use strict";(self.webpackChunkcontrast_docs=self.webpackChunkcontrast_docs||[]).push([[1226],{5103:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>a,contentTitle:()=>o,default:()=>l,frontMatter:()=>i,metadata:()=>c,toc:()=>d});var n=r(4848),s=r(8453);const i={},o="Observability",c={id:"architecture/observability",title:"Observability",description:"The Contrast Coordinator exposes metrics in the",source:"@site/docs/architecture/observability.md",sourceDirName:"architecture",slug:"/architecture/observability",permalink:"/contrast/pr-preview/pr-460/next/architecture/observability",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Certificate authority",permalink:"/contrast/pr-preview/pr-460/next/architecture/certificates"},next:{title:"Known limitations",permalink:"/contrast/pr-preview/pr-460/next/known-limitations"}},a={},d=[{value:"Exposed metrics",id:"exposed-metrics",level:2}];function h(e){const t={a:"a",code:"code",h1:"h1",h2:"h2",p:"p",...(0,s.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.h1,{id:"observability",children:"Observability"}),"\n",(0,n.jsxs)(t.p,{children:["The Contrast Coordinator exposes metrics in the\n",(0,n.jsx)(t.a,{href:"https://prometheus.io/",children:"Prometheus"})," format. These can be monitored to quickly\nidentify problems in the gRPC layer or attestation errors. Prometheus metrics\nare numerical values associated with a name and additional key/values pairs,\ncalled labels."]}),"\n",(0,n.jsx)(t.h2,{id:"exposed-metrics",children:"Exposed metrics"}),"\n",(0,n.jsxs)(t.p,{children:["The Coordinator pod has the annotation ",(0,n.jsx)(t.code,{children:"prometheus.io/scrape"})," set to ",(0,n.jsx)(t.code,{children:"true"})," so\nit can be found by the ",(0,n.jsx)(t.a,{href:"https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config",children:"service discovery of\nPrometheus"}),".\nThe metrics can be accessed at the Coordinator pod at port ",(0,n.jsx)(t.code,{children:"9102"})," under the\n",(0,n.jsx)(t.code,{children:"/metrics"})," endpoint."]}),"\n",(0,n.jsxs)(t.p,{children:["The Coordinator starts two gRPC servers, one for the user API on port ",(0,n.jsx)(t.code,{children:"1313"})," and\none for the mesh API on port ",(0,n.jsx)(t.code,{children:"7777"}),". Metrics for both servers can be accessed\nusing different prefixes."]}),"\n",(0,n.jsxs)(t.p,{children:["All metric names for the user API are prefixed with ",(0,n.jsx)(t.code,{children:"userapi_grpc_server_"}),".\nExposed metrics include the number of  handled requests of the methods\n",(0,n.jsx)(t.code,{children:"SetManifest"})," and ",(0,n.jsx)(t.code,{children:"GetManifest"}),", which get called when ",(0,n.jsx)(t.a,{href:"../deployment#set-the-manifest",children:"setting the\nmanifest"})," and ",(0,n.jsx)(t.a,{href:"../deployment#verify-the-coordinator",children:"verifying the\nCoordinator"})," respectively. For each method\nyou can see the gRPC status code indicating whether the request succeeded or\nnot."]}),"\n",(0,n.jsxs)(t.p,{children:["For the mesh API, the metric names are prefixed with ",(0,n.jsx)(t.code,{children:"meshapi_grpc_server_"}),". The\nmetrics include similar data to the user API for the method ",(0,n.jsx)(t.code,{children:"NewMeshCert"})," which\ngets called by the ",(0,n.jsx)(t.a,{href:"../components#the-initializer",children:"Initializer"})," when starting a\nnew workload."]})]})}function l(e={}){const{wrapper:t}={...(0,s.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(h,{...e})}):h(e)}},8453:(e,t,r)=>{r.d(t,{R:()=>o,x:()=>c});var n=r(6540);const s={},i=n.createContext(s);function o(e){const t=n.useContext(i);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function c(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),n.createElement(i.Provider,{value:t},e.children)}}}]);