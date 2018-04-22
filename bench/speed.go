package bench

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type SpeedTest struct {
	Name     string
	URL      string
	Size     int64
	Duration time.Duration
	Speed    float64

	finished bool
}

func NewSpeedTest(name, url string, size int64) (st *SpeedTest) {
	return &SpeedTest{
		Name: name,
		URL:  url,
		Size: size,
	}
}

func (st *SpeedTest) Do() (err error) {
	t := time.Now()
	resp, err := http.Get(st.URL)
	if err != nil {
		return
	}

	_, err = io.Copy(&ZeroReadWriter{}, &io.LimitedReader{
		R: resp.Body,
		N: st.Size,
	})

	if err != nil {
		return
	}

	st.Duration = time.Since(t)
	st.Speed = float64(st.Size/1048576) / time.Since(t).Seconds()
	st.finished = true
	return
}

func (st *SpeedTest) Result() (result string) {
	return fmt.Sprintf("%s %dMB: time %.2fs, speed %.2fMB/s", st.Name, st.Size/1048576, st.Duration.Seconds(), st.Speed)
}
