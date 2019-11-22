package cvdh

import (
	"image"
	"reflect"
	"testing"
)

func TestCreateBatchTensorFromImageandGrayedEdgeKernel(t *testing.T) {
	type args struct {
		original      []image.Image
		edgedetection []image.Image
		NCHW          bool
	}
	tests := []struct {
		name string
		args args
		want *Tensor4d
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateBatchTensorFromImageandGrayedEdgeKernel(tt.args.original, tt.args.edgedetection, tt.args.NCHW); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBatchTensorFromImageandGrayedEdgeKernel() = %v, want %v", got, tt.want)
			}
		})
	}
}
