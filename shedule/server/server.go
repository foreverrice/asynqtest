package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"asynqtest/tpl"

	"github.com/hibiken/asynq"
)

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:63779", Password: "G62m50oigInC30sf"},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	//关闭民宿订单任务
	mux.HandleFunc(tpl.EMAIL_TPL, emailMqHandler)
	mux.HandleFunc(tpl.DELAY_TPL, emailMqHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func emailMqHandler(ctx context.Context, t *asynq.Task) error {

	var p tpl.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("emailMqHandler err:%+v", err)
	}

	fmt.Printf("p : %+v \n", p)

	return nil
}
