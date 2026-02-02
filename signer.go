package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Тип job для pipeline
type job func(in, out chan interface{})

// Счётчик одновременных вызовов md5
var md5Counter int32

// DataSignerMd5 - мок функция (10мс, только 1 одновременный вызов)
func DataSignerMd5(data string) string {
	atomic.AddInt32(&md5Counter, 1)
	if atomic.LoadInt32(&md5Counter) > 1 {
		panic("DataSignerMd5 вызвана более 1 раза одновременно!")
	}
	time.Sleep(10 * time.Millisecond)
	hash := md5.Sum([]byte(data))
	atomic.AddInt32(&md5Counter, -1)
	return fmt.Sprintf("%x", hash)
}

// DataSignerCrc32 - мок функция (1 сек)
func DataSignerCrc32(data string) string {
	time.Sleep(time.Second)
	crcH := crc32.ChecksumIEEE([]byte(data))
	return strconv.FormatUint(uint64(crcH), 10)
}

// ExecutePipeline обеспечивает конвейерную обработку функций-воркеров
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	for _, j := range jobs {
		out := make(chan interface{})
		go func(j job, in, out chan interface{}) {
			j(in, out)
			close(out)
		}(j, in, out)
		in = out
	}

	for range in {
	}
}

// SingleHash считает crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			dataStr := strconv.Itoa(data.(int))
			fmt.Println(dataStr, "SingleHash data", dataStr)

			mu.Lock()
			md5Data := DataSignerMd5(dataStr)
			mu.Unlock()
			fmt.Println(dataStr, "SingleHash md5(data)", md5Data)

			crc32Chan := make(chan string)
			crc32Md5Chan := make(chan string)

			go func() {
				crc32Chan <- DataSignerCrc32(dataStr)
			}()

			go func() {
				crc32Md5Chan <- DataSignerCrc32(md5Data)
			}()

			crc32Data := <-crc32Chan
			crc32Md5Data := <-crc32Md5Chan

			fmt.Println(dataStr, "SingleHash crc32(data)", crc32Data)
			fmt.Println(dataStr, "SingleHash crc32(md5(data))", crc32Md5Data)

			result := crc32Data + "~" + crc32Md5Data
			fmt.Println(dataStr, "SingleHash result", result)
			out <- result
		}(data)
	}

	wg.Wait()
}

// MultiHash считает конкатенацию crc32(th+data) для th=0..5
func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			dataStr := data.(string)
			results := make([]string, 6)
			var innerWg sync.WaitGroup

			for th := 0; th < 6; th++ {
				innerWg.Add(1)
				go func(th int) {
					defer innerWg.Done()
					results[th] = DataSignerCrc32(strconv.Itoa(th) + dataStr)
					fmt.Println(dataStr, "MultiHash: crc32(th+step1))", th, results[th])
				}(th)
			}

			innerWg.Wait()
			result := strings.Join(results, "")
			fmt.Println(dataStr, "MultiHash result:", result)
			out <- result
		}(data)
	}

	wg.Wait()
}

// CombineResults собирает все результаты, сортирует и объединяет через _
func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, data.(string))
	}

	sort.Strings(results)
	result := strings.Join(results, "_")
	fmt.Println("CombineResults", result)
	out <- result
}

func main() {
	inputData := []int{0, 1}

	// Функция, которая отправляет данные в pipeline
	inputJob := func(in, out chan interface{}) {
		for _, v := range inputData {
			out <- v
		}
	}

	start := time.Now()

	ExecutePipeline(inputJob, SingleHash, MultiHash, CombineResults)

	elapsed := time.Since(start)
	fmt.Printf("\nВремя выполнения: %v\n", elapsed)

	expected := "29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542"
	fmt.Printf("Ожидаемый результат:\n%s\n", expected)
}
