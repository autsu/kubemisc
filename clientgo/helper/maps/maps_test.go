package maps

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		name    string
		m1      map[string]string
		m2      map[string]string
		wantRes bool
	}{
		{
			name:    "test1",
			m1:      map[string]string{"k": "v", "k1": "v1"},
			m2:      map[string]string{"k": "v"},
			wantRes: true,
		},
		{
			name:    "test2",
			m1:      map[string]string{"k": "v", "k1": "v1"},
			m2:      map[string]string{"k": "v1"},
			wantRes: false,
		},
		{
			name:    "test3",
			m1:      map[string]string{"k": "v", "k1": "v1"},
			m2:      map[string]string{"k": "v", "k1": "v1", "k2": "v2"},
			wantRes: false,
		}, {
			name:    "test4",
			m1:      map[string]string{"k": "v"},
			m2:      map[string]string{},
			wantRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Contains(tt.m1, tt.m2)
			if res != tt.wantRes {
				t.Errorf("%v test failed, want %v but got %v \n", tt.name, res, tt.wantRes)
			}
		})
	}
}
