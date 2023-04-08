kubebuilder init --domain lubenwei.io --repo lubenwei.io/autoservice
kubebuilder create api --group chore --version v1 --kind AutoService

# Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
make manifests
