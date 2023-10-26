package remotecommand

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"k8s.io/klog/v2"
)

// ExportDir exports file content.
// 将 path 执行 tar 压缩，并将压缩后的内容返回（返回为 io.Reader）
func (o *ExecOptions) ExportDir(path string) io.Reader {
	reader, outStream := io.Pipe()
	o.Out = outStream
	// 这里在容器中执行压缩命令，因为这里直接执行标准输出，所以就算 path 不存在也会输出一个空 path
	// 完整的命令是 tar cf - [path]，其中 - 代表将结果输出到 stdout 而不是写入到磁盘
	o.command = []string{"tar", "cf", "-", path}
	go func() {
		defer outStream.Close()
		if err := o.execute(); err != nil {
			klog.Errorf("failed to execute command in container: %v", err)
		}
	}()
	return reader
}

func (o *ExecOptions) TarFile(dst, src string) error {
	o.Out = os.Stdout
	o.command = []string{"tar", "cf", dst, src}
	if err := o.execute(); err != nil {
		klog.Errorf("failed to execute command in container: %v", err)
		return fmt.Errorf("execute command error, command=%v, error=%v", o.command, err)
	}
	return nil
}

func (o *ExecOptions) GzipFile(src string) error {
	o.Out = os.Stdout
	o.command = []string{"gzip", src}
	if err := o.execute(); err != nil {
		klog.Errorf("failed to execute command in container: %v", err)
		return fmt.Errorf("execute command error, command=%v, error=%v", o.command, err)
	}
	return nil
}

func (o *ExecOptions) ReadFile(filename string) ([]byte, error) {
	if o.err != nil {
		return nil, o.err
	}
	outStream := bytes.NewBuffer([]byte{})
	o.Out = outStream
	o.command = []string{"cat", filename}
	if err := o.execute(); err != nil {
		return nil, err
	}
	return io.ReadAll(outStream)
}

func (o *ExecOptions) RemoveFile(src string) error {
	o.Out = os.Stdout
	o.command = []string{"rm", "-d", src}
	if err := o.execute(); err != nil {
		klog.Errorf("failed to execute command in container: %v", err)
		return fmt.Errorf("execute command error, command=%v, error=%v", o.command, err)
	}
	return nil
}

func (o *ExecOptions) Error() error {
	return o.err
}
