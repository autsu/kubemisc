package remotecommand

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func TestExecExportDir(t *testing.T) {
	ctx := context.TODO()
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	remoteCommand, err := NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	reader := remoteCommand.
		BuildExecOptions(ctx, metav1.NamespaceDefault, "busybox", "").
		ExportDir("/bin")

	if _, err := io.Copy(io.Discard, reader); err != nil {
		panic(err)
	}
}

// 对容器中的某个文件/文件夹执行 tar + gzip 命令，然后将压缩后的文件下载到本地
// TODO: 用不着这么麻烦，tar -z 参数就是执行 gzip 压缩
func TestExecTarAndGzip(t *testing.T) {
	ctx := context.TODO()
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	remoteCommand, err := NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	exector := remoteCommand.BuildExecOptions(ctx, metav1.NamespaceDefault, "busybox", "")

	path := "/bin"
	tmpFileName := fmt.Sprintf("tmp-%s", strconv.Itoa(int(time.Now().Unix())))
	tmpFileNameTar := tmpFileName + ".tar"
	tmpFileNameTarGiz := tmpFileNameTar + ".gz"

	defer func() {
		exector.RemoveFile(tmpFileNameTarGiz)
	}()

	err = exector.TarFile(tmpFileNameTar, path)
	if err != nil {
		slog.Error("tar file error", err)
		return
	}

	err = exector.GzipFile(tmpFileNameTar)
	if err != nil {
		slog.Error("gzip file error", err)
		return
	}

	b, err := exector.ReadFile(tmpFileNameTarGiz)
	if err != nil {
		slog.Error("read file error", err)
		return
	}

	// 将容器中压缩后的文件下载到本地
	f, err := os.Create("/simplectrl/bin.tar.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		panic(err)
	}
}
