package main

import (
	"context"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type FileWatch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	p           string
	f           *os.File
	q           workqueue.RateLimitingInterface
	checkTime   time.Duration
	lastModTime time.Time
}

func (f *FileWatch) DeepCopyObject() runtime.Object {
	deepCopy := NewFileWatch(f.p, f.checkTime)
	deepCopy.lastModTime = f.lastModTime
	deepCopy.q = f.q
	return deepCopy
}

func NewFileWatch(filepath string, checkTime time.Duration) *FileWatch {
	fw := new(FileWatch)
	fw.p = filepath
	fw.checkTime = checkTime

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	fw.f = f

	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fw.lastModTime = stat.ModTime()

	return fw
}

func (f *FileWatch) Sync() {
	ticket := time.NewTicker(f.checkTime)
	go func() {
		for {
			select {
			case <-ticket.C:
				modify, modTime, err := f.FileIsModify()
				if err != nil {
					klog.Error(err)
					return
				}
				if modify {
					// Add 的参数必须是 ctrl.Request 类型的，其他类型会直接被 controller 丢弃
					f.q.Add(ctrl.Request{NamespacedName: types.NamespacedName{Name: f.p}})
					f.lastModTime = modTime
				}
			}
		}
	}()
}

func (f *FileWatch) FileIsModify() (bool, time.Time, error) {
	//klog.Info("Checking file ", f.p)
	stat, err := os.Stat(f.p)
	if err != nil {
		return false, time.Time{}, err
	}
	if stat.ModTime().After(f.lastModTime) {
		return true, stat.ModTime(), nil
	}
	return false, stat.ModTime(), nil
}

func (f *FileWatch) Start(ctx context.Context, h handler.EventHandler, queue workqueue.RateLimitingInterface, p ...predicate.Predicate) error {
	klog.Info("fileWatch start...")
	f.q = queue
	go f.Sync()
	return nil
}
