kubebuilder init --domain lubenwei.io --repo lubenwei.io/autoservice
kubebuilder create api --group chore --version v1 --kind AutoService
kubebuilder create api --group testapp --version v1 --kind Foo

# Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
make manifests
