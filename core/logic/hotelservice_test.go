package logic

//
//import (
//	"hotel-engine/utils/array"
//	"reflect"
//	"testing"
//)
//
//func Test_chunks(t *testing.T) {
//	type args struct {
//		ids  []string
//		size int
//	}
//	tests := []struct {
//		name string
//		args args
//		want [][]string
//	}{
//		{
//			name: "default",
//			args: args{
//				ids:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
//				size: 4,
//			},
//			want: [][]string{{"1", "2", "3", "4"}, {"5", "6", "7", "8"}, {"9", "10", "11"}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := array.Chunks(tt.args.ids, tt.args.size); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("chunks() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
