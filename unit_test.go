package main_test

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	t = atomic.Int64{}
)

func BenchmarkMain(b *testing.B) {
	defer b.Log(t.Load())
	client := http.NewSimpleClient()
	for i := 0; i < b.N; i++ {
		t.Add(1)
		client.ExecuteRequest(http.NewRequestBuilder().
			WithPath("http://localhost:10000/akasha-whisper/v1/models").
			WithHeader("Authorization", "Bearer aw_1ac5e86ef34946448b2a1a9d94024ee6").
			WithMethod(http.GET).WithJsonBody(nil),
		)
	}
}

func TestMainF(t *testing.T) {
	client := http.NewSimpleClient()
	st := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(100)

	for range 100 {
		go func() {
			res, e := client.ExecuteRequest(http.NewRequestBuilder().
				WithPath("http://localhost:10000/akasha-whisper/v1/chat/completions").
				WithHeader("Authorization", "Bearer aw_1ac5e86ef34946448b2a1a9d94024ee6").
				WithMethod(http.POST).WithJsonBody(&openai.CompleteChatRequestBody{Model: "gpt-3.5-turbo", Messages: []openai.ChatMessageObject{{Role: "user", Content: "hello"}}}),
			)
			resMap := map[string]any{}
			res.BindJson(&resMap)
			t.Log(resMap, e)
			wg.Done()
		}()
	}

	wg.Wait()
	t.Log(time.Since(st))
}
