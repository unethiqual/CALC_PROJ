package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"DistributedArithmeticExpressionCalculator/agent/client"
)

func worker() {
	for {
		task, err := client.GetTask()
		if err != nil {
			log.Println("Ошибка получения задачи:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Получена задача: %+v\n", task)
		result := client.Compute(task)
		log.Printf("Задача %d выполнена, результат: %f\n", task.ID, result)
		if err := client.PostResult(task.ID, result); err != nil {
			log.Println("Ошибка отправки результата:", err)
		}
	}
}

func main() {
	cpStr := os.Getenv("COMPUTING_POWER")
	cp, err := strconv.Atoi(cpStr)
	if err != nil || cp < 1 {
		cp = 2
	}
	for i := 0; i < cp; i++ {
		go worker()
	}
	select {}
}