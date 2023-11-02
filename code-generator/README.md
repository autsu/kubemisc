（🤮好恶心的 code-generator）

运行：
```shell
go mod vendor
chmod +x ./vendor/k8s.io/code-generator/generate-groups.sh
chmod +x ./vendor/k8s.io/code-generator/generate-internal-groups.sh
bash hack/update-codegen.sh
```

遇到的问题：
只生成了 client 的代码，deepcopy,informer,lister 的没有生成

解决：
重点是脚本
```shell
"${CODEGEN_PKG}/generate-groups.sh" "deepcopy,client,informer,lister" \
  fuck.codegenerator.com/gen/generated \
  fuck.codegenerator.com/gen/pkg/apis \
  samplecrd:v1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
```
中的 `fuck.codegenerator.com/gen/pkg/apis \`
这一行，必须要写对地方，不然可能是 code-generator 读取不到 apis 里面的 struct，导致无法生成代码
格式是 <gomodule name>/<apis 所在路径>

然后运行完脚本以后，它会在项目根目录下生成一个 fuck.codegenerator.com/gen 目录，里面又有两个目录 generated ，这个目录里面是生成好的
clientset, informers, listers 代码，然后还有一个 pkg/apis/samplecrd/v1 目录，里面有一个 zz_generated.deepcopy.go 文件

我不知道如何像网上的博客一样，运行后直接生成在项目目录下，而不是生成在 fuck.codegenerator.com/gen 目录下，所以没有办法，只能手动操作，
把 pkg/apis/samplecrd/v1/zz_generated.deepcopy.go 移到 pkg/samplecrd/v1 下，然后把 fuck.codegenerator.com/gen/generated
移到 pkg 下

然后还要修改 fuck.codegenerator.com/gen/generated/informers 里面的包引用，将其中所有的 
internalclientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
改为
internalclientset "fuck.codegenerator.com/gen/pkg/generated/clientset/versioned"

然后 go mod tidy
最后删掉 vendor，不然 goland 会依赖报红

再次运行脚本，会删掉 pkg/apis/samplecrd/v1/zz_generated.deepcopy.go，真 sb
狗屎一样的东西真尼玛难用
不搞了，在这上面纯属浪费时间