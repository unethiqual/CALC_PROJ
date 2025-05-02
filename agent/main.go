package main

import (
    "log"
    "os"
    "strconv"
    "sync"
    "time"

    "github.com/unethiqual/CALC_PROJ/agent/client"
)

func worker(wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        task, err := client.GetTask()
        if err != nil {
            log.Println("Error fetching task:", err)
            time.Sleep(1 * time.Second)
            continue
        }
        if task == nil {
            time.Sleep(1 * time.Second)
            continue
        }

        result := compute(task.Arg1, task.Arg2, task.Operation)
        if err := client.SubmitResult(task.ID, result); err != nil {
            log.Println("Error submitting result:", err)
        }
    }
}

func compute(arg1, arg2 float64, operation string) float64 {
    switch operation {
    case "+":
        return arg1 + arg2
    case "-":
        return arg1 - arg2
    case "*":
        return arg1 * arg2
    case "/":
        if arg2 != 0 {
            return arg1 / arg2
        }
    }
    return 0
}

func main() {
    client.InitGRPCClient()
    defer client.CloseGRPCClient()

    computingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
    var wg sync.WaitGroup
    for i := 0; i < computingPower; i++ {
        wg.Add(1)
        go worker(&wg)
    }
    wg.Wait()
}