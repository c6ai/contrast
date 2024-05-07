"use strict";(self.webpackChunkcontrast_docs=self.webpackChunkcontrast_docs||[]).push([[5945],{6022:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>o,default:()=>d,frontMatter:()=>r,metadata:()=>c,toc:()=>h});var i=t(4848),s=t(8453);const r={},o="Service Mesh",c={id:"components/service-mesh",title:"Service Mesh",description:"The Contrast service mesh secures the communication of the workload by automatically",source:"@site/docs/components/service-mesh.md",sourceDirName:"components",slug:"/components/service-mesh",permalink:"/contrast/pr-preview/pr-428/next/components/service-mesh",draft:!1,unlisted:!1,tags:[],version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Policies",permalink:"/contrast/pr-preview/pr-428/next/components/policies"},next:{title:"Architecture",permalink:"/contrast/pr-preview/pr-428/next/architecture/"}},a={},h=[{value:"Configuring the Proxy",id:"configuring-the-proxy",level:2},{value:"Ingress",id:"ingress",level:3},{value:"Egress",id:"egress",level:3}];function l(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,s.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"service-mesh",children:"Service Mesh"}),"\n",(0,i.jsxs)(n.p,{children:["The Contrast service mesh secures the communication of the workload by automatically\nwrapping the network traffic inside mutual TLS (mTLS) connections. The\nverification of the endpoints in the connection establishment is based on\ncertificates that are part of the\n",(0,i.jsx)(n.a,{href:"/contrast/pr-preview/pr-428/next/architecture/certificates",children:"PKI of the Coordinator"}),"."]}),"\n",(0,i.jsxs)(n.p,{children:["The service mesh can be enabled on a per-pod basis by adding the ",(0,i.jsx)(n.code,{children:"service-mesh"}),"\ncontainer as a ",(0,i.jsx)(n.a,{href:"https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/",children:"sidecar container"}),".\nThe service mesh container first sets up ",(0,i.jsx)(n.code,{children:"iptables"}),"\nrules based on its configuration and then starts ",(0,i.jsx)(n.a,{href:"https://www.envoyproxy.io/",children:"Envoy"}),"\nfor TLS origination and termination."]}),"\n",(0,i.jsx)(n.h2,{id:"configuring-the-proxy",children:"Configuring the Proxy"}),"\n",(0,i.jsxs)(n.p,{children:["The service mesh container can be configured using the ",(0,i.jsx)(n.code,{children:"EDG_INGRESS_PROXY_CONFIG"}),"\nand ",(0,i.jsx)(n.code,{children:"EDG_EGRESS_PROXY_CONFIG"})," environment variables."]}),"\n",(0,i.jsx)(n.h3,{id:"ingress",children:"Ingress"}),"\n",(0,i.jsxs)(n.p,{children:["All TCP ingress traffic is routed over Envoy by default. Since we use\n",(0,i.jsx)(n.a,{href:"https://docs.kernel.org/networking/tproxy.html",children:"TPROXY"}),", the destination address\nremains the same throughout the packet handling."]}),"\n",(0,i.jsxs)(n.p,{children:["Any incoming connection is required to present a client certificate signed by the\n",(0,i.jsx)(n.a,{href:"/contrast/pr-preview/pr-428/next/architecture/certificates#usage-of-the-different-certificates",children:"mesh CA certificate"}),".\nEnvoy presents a certificate chain of the mesh\ncertificate of the workload and the intermediate CA certificate as the server certificate."]}),"\n",(0,i.jsxs)(n.p,{children:["If the deployment contains workloads which should be reachable from outside the\nService Mesh, while still handing out the certificate chain, disable client\nauthentication by setting the environment variable ",(0,i.jsx)(n.code,{children:"EDG_INGRESS_PROXY_CONFIG"})," as\n",(0,i.jsx)(n.code,{children:"<name>#<port>#false"}),". Separate multiple entries with ",(0,i.jsx)(n.code,{children:"##"}),". You can choose any\ndescriptive string identifying the service on the given port for the\ninformational-only field ",(0,i.jsx)(n.code,{children:"<name>"}),"."]}),"\n",(0,i.jsxs)(n.p,{children:["Disable redirection and TLS termination altogether by specifying\n",(0,i.jsx)(n.code,{children:"<name>#<port>#true"}),". This can be beneficial if the workload itself handles TLS\non that port or if the information exposed on this port is non-sensitive."]}),"\n",(0,i.jsx)(n.p,{children:"The following example workload exposes a web service on port 8080 and metrics on\nport 7890. The web server is exposed to a 3rd party end-user which wants to\nverify the deployment, therefore it's still required that the server hands out\nit certificate chain signed by the mesh CA certificate. The metrics should be\nexposed via TCP without TLS."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-yaml",children:'apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web\nspec:\n  replicas: 1\n  template:\n    spec:\n      runtimeClassName: contrast-cc\n      initContainers:\n        - name: initializer\n          image: "ghcr.io/edgelesssys/contrast/initializer@sha256:..."\n          env:\n            - name: COORDINATOR_HOST\n              value: coordinator\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n        - name: sidecar\n          image: "ghcr.io/edgelesssys/contrast/service-mesh-proxy@sha256:..."\n          restartPolicy: Always\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n          env:\n            - name: EDG_INGRESS_PROXY_CONFIG\n              value: "web#8080#false##metrics#7890#true"\n          securityContext:\n            privileged: true\n            capabilities:\n              add:\n                - NET_ADMIN\n      containers:\n        - name: web-svc\n          image: ghcr.io/edgelesssys/frontend:v1.2.3@...\n          ports:\n            - containerPort: 8080\n              name: web\n            - containerPort: 7890\n              name: metrics\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n      volumes:\n        - name: tls-certs\n          emptyDir: {}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"egress",children:"Egress"}),"\n",(0,i.jsx)(n.p,{children:"To be able to route the egress traffic of the workload through Envoy, the remote\nendpoints' IP address and port must be configurable."}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Choose an IP address inside the ",(0,i.jsx)(n.code,{children:"127.0.0.0/8"})," CIDR and a port not yet in use\nby the pod."]}),"\n",(0,i.jsx)(n.li,{children:"Configure the workload to connect to this IP address and port."}),"\n",(0,i.jsxs)(n.li,{children:["Set ",(0,i.jsx)(n.code,{children:"<name>#<chosen IP>:<chosen port>#<original-hostname-or-ip>:<original-port>"}),"\nas ",(0,i.jsx)(n.code,{children:"EDG_EGRESS_PROXY_CONFIG"}),". Separate multiple entries with ",(0,i.jsx)(n.code,{children:"##"}),". Choose any\nstring identifying the service on the given port as ",(0,i.jsx)(n.code,{children:"<name>"}),"."]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["This redirects the traffic over Envoy. The endpoint must present a valid\ncertificate chain which must be verifiable with the\n",(0,i.jsx)(n.a,{href:"/contrast/pr-preview/pr-428/next/architecture/certificates#usage-of-the-different-certificates",children:"mesh CA certificate"}),".\nFurthermore, Envoy uses a certificate chain with the mesh certificate of the workload\nand the intermediate CA certificate as the client certificate."]}),"\n",(0,i.jsxs)(n.p,{children:["The following example workload has no ingress connections and two egress\nconnection to different microservices. The microservices are themselves part\nof the confidential deployment. One is reachable under ",(0,i.jsx)(n.code,{children:"billing-svc:8080"})," and\nthe other under ",(0,i.jsx)(n.code,{children:"cart-svc:8080"}),"."]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-yaml",children:'apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web\nspec:\n  replicas: 1\n  template:\n    spec:\n      runtimeClassName: contrast-cc\n      initContainers:\n        - name: initializer\n          image: "ghcr.io/edgelesssys/contrast/initializer@sha256:..."\n          env:\n            - name: COORDINATOR_HOST\n              value: coordinator\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n        - name: sidecar\n          image: "ghcr.io/edgelesssys/contrast/service-mesh-proxy@sha256:..."\n          restartPolicy: Always\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n          env:\n            - name: EDG_EGRESS_PROXY_CONFIG\n              value: "billing#127.137.0.1:8081#billing-svc:8080##cart#127.137.0.2:8081#cart-svc:8080"\n          securityContext:\n            privileged: true\n            capabilities:\n              add:\n                - NET_ADMIN\n      containers:\n        - name: currency-conversion\n          image: ghcr.io/edgelesssys/conversion:v1.2.3@...\n          volumeMounts:\n            - name: tls-certs\n              mountPath: /tls-config\n      volumes:\n        - name: tls-certs\n          emptyDir: {}\n'})})]})}function d(e={}){const{wrapper:n}={...(0,s.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(l,{...e})}):l(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>o,x:()=>c});var i=t(6540);const s={},r=i.createContext(s);function o(e){const n=i.useContext(r);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function c(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),i.createElement(r.Provider,{value:n},e.children)}}}]);