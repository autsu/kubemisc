package main

import (
	"context"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
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
					// Add 的参数必须是 ctrl.Request 类型的，其他类型会直接被 controller 丢弃
					f.q.Add(ctrl.Request{NamespacedName: types.NamespacedName{Name: f.p}})
					f.lastModTime = modTime
				}
			}
		}
	}()
}

func (f *FileWatch) FileIsModify() (bool, time.Time, error) {
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
