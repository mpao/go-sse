package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func generator(ctx context.Context, ch chan<- int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ticker := time.NewTicker(time.Second)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		case <-ticker.C:
			p := r.Intn(100)
			ch <- p
		}
	}
	ticker.Stop()
	close(ch)
}

func formatter(event string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}
	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(m)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}
