package codec

import "testing"

func Test_decodeParameters(t *testing.T) {
	type args struct {
		u string
		t *testing.T
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				u: "http://127.0.0.1:8080/apis/clusterpedia.io/v1beta1/collectionresources/workloads?labelSelector=search.clusterpedia.io/namespaces=kube-system",
				t: t,
			},
		}, {
			name: "",
			args: args{
				u: "http://127.0.0.1:8080/apis/clusterpedia.io/v1beta1/collectionresources/workloads?labelSelector=search.clusterpedia.io/namespaces in (defualt,kube-system)",
				t: t,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeParameters(tt.args.u, tt.args.t)
		})
	}
}
