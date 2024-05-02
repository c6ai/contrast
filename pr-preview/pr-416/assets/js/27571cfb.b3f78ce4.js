"use strict";(self.webpackChunkcontrast_docs=self.webpackChunkcontrast_docs||[]).push([[9234],{3183:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>l,contentTitle:()=>o,default:()=>d,frontMatter:()=>s,metadata:()=>a,toc:()=>c});var i=t(4848),r=t(8453);const s={},o="Known Limitations",a={id:"known-limitations",title:"Known Limitations",description:"As Contrast is currently in an early preview stage, it's built on several projects that are also under active development.",source:"@site/docs/known-limitations.md",sourceDirName:".",slug:"/known-limitations",permalink:"/contrast/pr-preview/pr-416/next/known-limitations",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Certificate authority",permalink:"/contrast/pr-preview/pr-416/next/architecture/certificates"},next:{title:"About",permalink:"/contrast/pr-preview/pr-416/next/about/"}},l={},c=[{value:"Availability",id:"availability",level:2},{value:"Kubernetes Features",id:"kubernetes-features",level:2},{value:"Runtime Policies",id:"runtime-policies",level:2},{value:"Tooling Integration",id:"tooling-integration",level:2}];function u(e){const n={a:"a",em:"em",h1:"h1",h2:"h2",li:"li",p:"p",strong:"strong",ul:"ul",...(0,r.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"known-limitations",children:"Known Limitations"}),"\n",(0,i.jsx)(n.p,{children:"As Contrast is currently in an early preview stage, it's built on several projects that are also under active development.\nThis section outlines the most significant known limitations, providing stakeholders with clear expectations and understanding of the current state."}),"\n",(0,i.jsx)(n.h2,{id:"availability",children:"Availability"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Platform Support"}),": At present, Contrast is exclusively available on Azure AKS, supported by the ",(0,i.jsx)(n.a,{href:"https://learn.microsoft.com/en-us/azure/confidential-computing/confidential-containers-on-aks-preview",children:"Confidential Container preview for AMD SEV-SNP"}),". Expansion to other cloud platforms is planned, pending the availability of necessary infrastructure enhancements."]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"kubernetes-features",children:"Kubernetes Features"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Persistent Volumes"}),": Not currently supported within Confidential Containers."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Port-Forwarding"}),": This feature isn't yet supported by Kata Containers."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Resource Limits"}),": There is an existing bug on AKS where container memory limits are incorrectly applied. The current workaround involves using only memory requests instead of limits."]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"runtime-policies",children:"Runtime Policies"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Coverage"}),": While the enforcement of workload policies generally functions well, ",(0,i.jsx)(n.a,{href:"https://github.com/microsoft/kata-containers/releases/tag/genpolicy-0.6.2-5",children:"there are scenarios not yet fully covered"}),". It's crucial to review deployments specifically for these edge cases."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"Policy Evaluation"}),": The current policy evaluation mechanism on API requests isn't stateful, which means it can't ensure a prescribed order of events. Consequently, there's no guaranteed enforcement that the ",(0,i.jsx)(n.a,{href:"/contrast/pr-preview/pr-416/next/components/#the-initializer",children:"initializer"})," container runs ",(0,i.jsx)(n.em,{children:"before"})," the workload container. This order is vital for ensuring that all traffic between pods is securely encapsulated within TLS connections. TODO: Consequences"]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"tooling-integration",children:"Tooling Integration"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"CLI Availability"}),": The CLI tool is currently only available for Linux. This limitation arises because certain upstream dependencies haven't yet been ported to other platforms."]}),"\n"]})]})}function d(e={}){const{wrapper:n}={...(0,r.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(u,{...e})}):u(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>o,x:()=>a});var i=t(6540);const r={},s=i.createContext(r);function o(e){const n=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:o(e.components),i.createElement(s.Provider,{value:n},e.children)}}}]);