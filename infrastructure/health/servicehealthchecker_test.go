package health

// import (
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/joho/godotenv"
// )

// func Test_serviceHealthChecker_Check(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		c        Checker
// 		want     HealthResultDto
// 		attempts int
// 	}{
// 		{
// 			name: "Healthy10_3",
// 			c:    NewServiceHealthChecker("hotelBaseService", "http://client-api.k8s.indra.local/api/v1/alive2", getThresholdInSecond()),
// 			want: HealthResultDto{
// 				Status: Healthy,
// 			},
// 			attempts: 7,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := check(tt.c, tt.attempts)
// 			if !reflect.DeepEqual(got.Status, tt.want.Status) {
// 				t.Errorf("serviceHealthChecker.Check() = %v, want %v", got.Status, tt.want.Status)
// 			}
// 		})
// 	}
// }

// func check(c Checker, attempts int) HealthResultDto {
// 	var got = HealthResultDto{}
// 	for i := 1; i <= attempts; i++ {
// 		got = c.Check()
// 		time.Sleep(time.Duration(10) * time.Second)
// 	}
// 	return got
// }

// func getThresholdInSecond() int64 {
// 	godotenv.Load("../../dev.env")
// 	thresholdInSecond, _ := strconv.ParseInt(os.Getenv("HOTEL_ENGINE_HealthCheck_ThresholdInSecond"), 10, 64)
// 	return thresholdInSecond
// }
