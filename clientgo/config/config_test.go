package config

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestKubeConfig1(t *testing.T) {
	cfg := getKubeConfig1()
	if wrapErr := testConfig(context.TODO(), cfg); wrapErr != nil {
		t.Fatal(wrapErr.err)
	}
}

func TestKubeConfig2(t *testing.T) {
	cfgs := getKubeConfig2()
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancelFunc()
	for name, clusterCfg := range cfgs {
		if wrapErr := testConfig(ctx, clusterCfg); wrapErr != nil {
			switch wrapErr.kind {
			case newClientSetError:
				fmt.Fprintf(os.Stderr, "connect to cluster [%v] error: %v\n", name, wrapErr.err)
			case listPodsError:
				fmt.Fprintf(os.Stderr, "list pods from context [%v] error: %v\n", name, wrapErr.err)
			}
			continue
		}
	}
}

func TestKubeConfigContent(t *testing.T) {
	t.Logf("%+v\n", getKubeConfig1())
}
