package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg                     = sync.WaitGroup{}
	ivestiesLaukimasSignal bool
	mu                     = sync.Mutex{}
)

func ivestiesLaukimas(kanalas chan int) {
	var inputas int
	fmt.Println("Kiek bus 5*5 ??? ")
	fmt.Scanf("%v", &inputas)

	mu.Lock()
	ivestiesLaukimasSignal = true
	mu.Unlock()

	kanalas <- inputas
	close(kanalas)
}

func skubintojas(signalas chan int) {
	time.Sleep(time.Second * 1)
	mu.Lock()
	if !ivestiesLaukimasSignal {

		signalas <- 1
	}
	mu.Unlock()

	close(signalas)
}

func main() {
	kanalas := make(chan int)
	signalas := make(chan int)
	wg.Add(1)

	go ivestiesLaukimas(kanalas)
	go skubintojas(signalas)

	go func() {
		select {
		case <-signalas:
			fmt.Println("tik.. tak..")
		}
	}()

	for {
		select {
		case gautasInputas, ok := <-kanalas:
			if !ok {
				wg.Done()
				return
			}

			if gautasInputas == 25 {
				fmt.Println("Šaunu! Skaičiuoti moki! :) ")
			} else {
				fmt.Println("Esi greitas, bet skaičiuoti nemoki :( ")
			}

			wg.Done()
			return

		case <-time.After(3 * time.Second):
			fmt.Println("Skaičiuoti nemoki :(")
			wg.Done()
			return
		}
	}
}
