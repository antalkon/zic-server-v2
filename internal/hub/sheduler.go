package hub

import (
	"backend/pkg/licenze"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

func InitScheduler() {
	scheduler = cron.New(cron.WithSeconds())

	_, err := scheduler.AddFunc("0 19 10 * * *", func() {
		runNightlyJob()
	})
	if err != nil {
		panic(fmt.Sprintf("❌ Не удалось добавить задачу в планировщик: %v", err))
	}

	scheduler.Start()
	fmt.Println("🛠️ Планировщик задач запущен")
}

func runNightlyJob() {
	fmt.Println("⏰ Задача проверки лицензии выполнена в:", time.Now().Format(time.RFC3339))
	licenze.CheckLicenze()
}

func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
		fmt.Println("🛑 Планировщик остановлен")
	}
}
