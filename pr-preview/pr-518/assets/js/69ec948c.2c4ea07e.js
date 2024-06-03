"use strict";(self.webpackChunkcontrast_docs=self.webpackChunkcontrast_docs||[]).push([[3712],{4654:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>a,default:()=>h,frontMatter:()=>s,metadata:()=>i,toc:()=>l});var o=n(4848),r=n(8453);const s={},a="Workload deployment",i={id:"deployment",title:"Workload deployment",description:"The following instructions will guide you through the process of making an existing Kubernetes deployment",source:"@site/versioned_docs/version-0.5/deployment.md",sourceDirName:".",slug:"/deployment",permalink:"/contrast/pr-preview/pr-518/0.5/deployment",draft:!1,unlisted:!1,editUrl:"https://github.com/edgelesssys/contrast/edit/main/docs/versioned_docs/version-0.5/deployment.md",tags:[],version:"0.5",frontMatter:{},sidebar:"docs",previous:{title:"Confidential emoji voting",permalink:"/contrast/pr-preview/pr-518/0.5/examples/emojivoto"},next:{title:"Architecture",permalink:"/contrast/pr-preview/pr-518/0.5/architecture/"}},c={},l=[{value:"Deploy the Contrast Coordinator",id:"deploy-the-contrast-coordinator",level:2},{value:"Prepare your Kubernetes resources",id:"prepare-your-kubernetes-resources",level:2},{value:"Generate policy annotations and manifest",id:"generate-policy-annotations-and-manifest",level:2},{value:"Apply the resources",id:"apply-the-resources",level:2},{value:"Connect to the Contrast Coordinator",id:"connect-to-the-contrast-coordinator",level:2},{value:"Set the manifest",id:"set-the-manifest",level:2},{value:"Verify the Coordinator",id:"verify-the-coordinator",level:2},{value:"Communicate with workloads",id:"communicate-with-workloads",level:2}];function d(e){const t={a:"a",admonition:"admonition",code:"code",h1:"h1",h2:"h2",p:"p",pre:"pre",...(0,r.R)(),...e.components},{TabItem:n,Tabs:s}=t;return n||p("TabItem",!0),s||p("Tabs",!0),(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(t.h1,{id:"workload-deployment",children:"Workload deployment"}),"\n",(0,o.jsx)(t.p,{children:"The following instructions will guide you through the process of making an existing Kubernetes deployment\nconfidential and deploying it together with Contrast."}),"\n",(0,o.jsxs)(t.p,{children:["A running CoCo-enabled cluster is required for these steps, see the ",(0,o.jsx)(t.a,{href:"/contrast/pr-preview/pr-518/0.5/getting-started/cluster-setup",children:"setup guide"})," on how to set it up."]}),"\n",(0,o.jsx)(t.h2,{id:"deploy-the-contrast-coordinator",children:"Deploy the Contrast Coordinator"}),"\n",(0,o.jsx)(t.p,{children:"Install the latest Contrast Coordinator release, comprising a single replica deployment and a\nLoadBalancer service, into your cluster."}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"kubectl apply -f https://github.com/edgelesssys/contrast/releases/latest/download/coordinator.yml\n"})}),"\n",(0,o.jsx)(t.h2,{id:"prepare-your-kubernetes-resources",children:"Prepare your Kubernetes resources"}),"\n",(0,o.jsx)(t.p,{children:"Contrast will add annotations to your Kubernetes YAML files. If you want to keep the original files\nunchanged, you can copy the files into a separate local directory.\nYou can also generate files from a Helm chart or from a Kustomization."}),"\n",(0,o.jsxs)(s,{groupId:"yaml-source",children:[(0,o.jsx)(n,{value:"kustomize",label:"kustomize",children:(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"mkdir resources\nkustomize build $MY_RESOURCE_DIR > resources/all.yml\n"})})}),(0,o.jsx)(n,{value:"helm",label:"helm",children:(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"mkdir resources\nhelm template $RELEASE_NAME $CHART_NAME > resources/all.yml\n"})})}),(0,o.jsx)(n,{value:"copy",label:"copy",children:(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"cp -R $MY_RESOURCE_DIR resources/\n"})})})]}),"\n",(0,o.jsxs)(t.p,{children:["To specify that a workload (pod, deployment, etc.) should be deployed as confidential containers,\nadd ",(0,o.jsx)(t.code,{children:"runtimeClassName: kata-cc-isolation"})," to the pod spec (pod definition or template).\nIn addition, add the Contrast Initializer as ",(0,o.jsx)(t.code,{children:"initContainers"})," to these workloads and configure the\nworkload to use the certificates written to a ",(0,o.jsx)(t.code,{children:"volumeMount"})," named ",(0,o.jsx)(t.code,{children:"tls-certs"}),"."]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-yaml",children:'spec: # v1.PodSpec\n  runtimeClassName: kata-cc-isolation\n  initContainers:\n  - name: initializer\n    image: "ghcr.io/edgelesssys/contrast/initializer:latest"\n    env:\n    - name: COORDINATOR_HOST\n      value: coordinator\n    volumeMounts:\n    - name: tls-certs\n      mountPath: /tls-config\n  volumes:\n  - name: tls-certs\n    emptyDir: {}\n'})}),"\n",(0,o.jsx)(t.h2,{id:"generate-policy-annotations-and-manifest",children:"Generate policy annotations and manifest"}),"\n",(0,o.jsxs)(t.p,{children:["Run the ",(0,o.jsx)(t.code,{children:"generate"})," command to generate the execution policies and add them as annotations to your\ndeployment files. A ",(0,o.jsx)(t.code,{children:"manifest.json"})," with the reference values of your deployment will be created."]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"contrast generate resources/\n"})}),"\n",(0,o.jsx)(t.h2,{id:"apply-the-resources",children:"Apply the resources"}),"\n",(0,o.jsx)(t.p,{children:"Apply the resources to the cluster. Your workloads will block in the initialization phase until a\nmanifest is set at the Coordinator."}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"kubectl apply -f resources/\n"})}),"\n",(0,o.jsx)(t.h2,{id:"connect-to-the-contrast-coordinator",children:"Connect to the Contrast Coordinator"}),"\n",(0,o.jsx)(t.p,{children:"For the next steps, we will need to connect to the Coordinator. The released Coordinator resource\nincludes a LoadBalancer definition we can use."}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"coordinator=$(kubectl get svc coordinator -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')\n"})}),"\n",(0,o.jsxs)(t.admonition,{title:"Port-forwarding of Confidential Containers",type:"info",children:[(0,o.jsxs)(t.p,{children:[(0,o.jsx)(t.code,{children:"kubectl port-forward"})," uses a Container Runtime Interface (CRI) method that isn't supported by the Kata shim.\nIf you can't use a public load balancer, you can deploy a ",(0,o.jsx)(t.a,{href:"https://github.com/edgelesssys/contrast/blob/main/deployments/emojivoto/portforwarder.yml",children:"port-forwarder"}),".\nThe port-forwarder relays traffic from a CoCo pod and can be accessed via ",(0,o.jsx)(t.code,{children:"kubectl port-forward"}),"."]}),(0,o.jsxs)(t.p,{children:["Upstream tracking issue: ",(0,o.jsx)(t.a,{href:"https://github.com/kata-containers/kata-containers/issues/1693",children:"https://github.com/kata-containers/kata-containers/issues/1693"}),"."]})]}),"\n",(0,o.jsx)(t.h2,{id:"set-the-manifest",children:"Set the manifest"}),"\n",(0,o.jsx)(t.p,{children:"Attest the Coordinator and set the manifest:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:'contrast set -c "${coordinator}:1313" resources/\n'})}),"\n",(0,o.jsx)(t.p,{children:"After this step, the Coordinator will start issuing TLS certs to the workloads. The init container\nwill fetch a certificate for the workload and the workload is started."}),"\n",(0,o.jsx)(t.h2,{id:"verify-the-coordinator",children:"Verify the Coordinator"}),"\n",(0,o.jsxs)(t.p,{children:["An end user (data owner) can verify the Contrast deployment using the ",(0,o.jsx)(t.code,{children:"verify"})," command."]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:'contrast verify -c "${coordinator}:1313"\n'})}),"\n",(0,o.jsxs)(t.p,{children:["The CLI will attest the Coordinator using embedded reference values. The CLI will write the service mesh\nroot certificate and the history of manifests into the ",(0,o.jsx)(t.code,{children:"verify/"})," directory. In addition, the policies referenced\nin the manifest are also written to the directory."]}),"\n",(0,o.jsx)(t.h2,{id:"communicate-with-workloads",children:"Communicate with workloads"}),"\n",(0,o.jsxs)(t.p,{children:["You can securely connect to the workloads using the Coordinator's ",(0,o.jsx)(t.code,{children:"mesh-root.pem"})," as a trusted CA certificate.\nFirst, expose the service on a public IP address via a LoadBalancer service:"]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"kubectl patch svc ${MY_SERVICE} -p '{\"spec\": {\"type\": \"LoadBalancer\"}}'\ntimeout 30s bash -c 'until kubectl get service/${MY_SERVICE} --output=jsonpath='{.status.loadBalancer}' | grep \"ingress\"; do sleep 2 ; done'\nlbip=$(kubectl get svc ${MY_SERVICE} -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')\necho $lbip\n"})}),"\n",(0,o.jsxs)(t.admonition,{title:"Subject alternative names and LoadBalancer IP",type:"info",children:[(0,o.jsx)(t.p,{children:"By default, mesh certificates are issued with a wildcard DNS entry. The web frontend is accessed\nvia load balancer IP in this demo. Tools like curl check the certificate for IP entries in the SAN field.\nValidation fails since the certificate contains no IP entries as a subject alternative name (SAN).\nFor example, a connection attempt using the curl and the mesh root certificate with throw the following error:"}),(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"$ curl --cacert ./verify/mesh-root.pem \"https://${frontendIP}:443\"\ncurl: (60) SSL: no alternative certificate subject name matches target host name '203.0.113.34'\n"})})]}),"\n",(0,o.jsxs)(t.p,{children:["Using ",(0,o.jsx)(t.code,{children:"openssl"}),", the certificate of the service can be validated with the ",(0,o.jsx)(t.code,{children:"mesh-root.pem"}),":"]}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-sh",children:"openssl s_client -CAfile verify/mesh-root.pem -verify_return_error -connect ${frontendIP}:443 < /dev/null\n"})})]})}function h(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,o.jsx)(t,{...e,children:(0,o.jsx)(d,{...e})}):d(e)}function p(e,t){throw new Error("Expected "+(t?"component":"object")+" `"+e+"` to be defined: you likely forgot to import, pass, or provide it.")}},8453:(e,t,n)=>{n.d(t,{R:()=>a,x:()=>i});var o=n(6540);const r={},s=o.createContext(r);function a(e){const t=o.useContext(s);return o.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function i(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:a(e.components),o.createElement(s.Provider,{value:t},e.children)}}}]);