package main

import (
	"asynqtest/tpl"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:63779"
const redisPwd = "G62m50oigInC30sf"

func main() {
	// 周期性任务
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPwd,
		}, nil)

	payload, err := json.Marshal(tpl.EmailPayload{Email: "344795118@qq.com", Content: "发邮件呀"})
	if err != nil {
		log.Fatal(err)
	}
	// taskid保持唯一性
	task := asynq.NewTask(tpl.EMAIL_TPL, payload, asynq.TaskID(tpl.EMAIL_TPL))
	// 每隔1分钟同步一次
	entryID, err := scheduler.Register("*/1 * * * *", task)

	// region 异步10s
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPwd,
	})
	payloadDelay1, err := json.Marshal(tpl.EmailPayload{Email: "344795118@qq.com", Content: "立马发邮件呀"})

	t1 := asynq.NewTask(tpl.DELAY_TPL, payloadDelay1)

	if err != nil {
		log.Fatal(err)
	}

	payloadDelay2, err := json.Marshal(tpl.EmailPayload{Email: "344795118@qq.com", Content: "延迟发邮件呀"})
	_, err = client.Enqueue(t1) // imme
	if err != nil {
		fmt.Println("立马失败")
	}

	t2 := asynq.NewTask(tpl.DELAY_TPL, payloadDelay2)
	_, err = client.Enqueue(t2, asynq.ProcessIn(10*time.Second))
	if err != nil {
		fmt.Println("延迟失败")
	}
	// endregion

	log.Printf("registered an entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
