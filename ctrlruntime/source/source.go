package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"

	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type FileWatch struct {
	p           string
	f           *os.File
	q           workqueue.RateLimitingInterface
	lastModTime time.Time
}

func NewFileWatch(filepath string) *FileWatch {
	fw := new(FileWatch)
	fw.p = filepath

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
	ticket := time.NewTicker(time.Second * 3)
	go func() {
		for {
			select {
			case <-ticket.C:
				modify, modTime, err := f.FileIsModify()
				if err != nil {
					return
				}
				if modify {
					klog.Info("a modify event")
					// Add 的对象必须是 ctrl.Request ?
					f.q.Add(ctrl.Request{NamespacedName: types.NamespacedName{Name: f.p}})
					//klog.Info("queue len: ", f.q.Len())
					f.lastModTime = modTime
				}
			}
		}
	}()
}

func (f *FileWatch) FileIsModify() (bool, time.Time, error) {
	//klog.Info("in")
	stat, err := os.Stat(f.p)
	if err != nil {
		return false, time.Time{}, err
	}
	//klog.Info("last mod time: ", f.lastModTime, "cur mod time: ", stat.ModTime())
	if stat.ModTime().After(f.lastModTime) {
		return true, stat.ModTime(), nil
	}
	return false, stat.ModTime(), nil
}

func (f *FileWatch) Start(ctx context.Context, h handler.EventHandler, queue workqueue.RateLimitingInterface, p ...predicate.Predicate) error {
	klog.Info("fileWatch start!!!")
	klog.Info("queue info", queue.Len())
	f.q = queue
	f.Sync()
	return nil
}
