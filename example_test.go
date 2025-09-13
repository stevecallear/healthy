package healthy_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stevecallear/healthy"
)

func ExampleNew() {
	addr := fmt.Sprintf(":%d", getFreePort())
	url := "http://localhost" + addr

	close := startHTTPDelayed(addr, 100*time.Millisecond)
	defer close()

	err := healthy.Wait(
		healthy.HTTP(url).Timeout(time.Millisecond).Expect(http.StatusOK),
		healthy.WithTimeout(time.Second),
		healthy.WithDelay(10*time.Millisecond),
		healthy.WithJitter(0),
	)

	fmt.Println(err)
	//output: <nil>
}
