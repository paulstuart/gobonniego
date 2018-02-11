package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"math/rand"
	"math"
	"runtime"
	"github.com/cloudfoundry/gosigar"
	"io/ioutil"
	"path"
	"io"
	"bytes"
)

const Blocksize = 0x1 << 16 // 65,536, 2^16

func testWritePerformance(filename string, fileSize int, randomBlock []byte, bytesWrittenChannel chan<- int) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)

	bytesWritten := 0
	for i := 0; i < fileSize; i += len(randomBlock) {
		n, err := w.Write(randomBlock)
		check(err)
		bytesWritten += n
	}

	w.Flush()
	f.Close()
	bytesWrittenChannel <- bytesWritten
}

func main() {

	//var ip = flag.Int("flagname", 1234, "help message for flagname")
	//flag.Parse()
	//fmt.Printf("ip: %d\n", *ip)

	cores := runtime.NumCPU()
	fmt.Printf("cores: %d\n", cores)

	mem := sigar.Mem{}
	mem.Get()
	fmt.Printf("memory: %d MiB\n", mem.Total>>20)

	fileSize := int(mem.Total) * 2 / cores

	dir, err := ioutil.TempDir("", "bonniego")
	fmt.Printf("directory: %s\n", dir)
	check(err)
	defer os.RemoveAll(dir)

	randomBlock := make([]byte, Blocksize)
	_, err = rand.Read(randomBlock)
	check(err)

	start := time.Now()

	bytesWrittenChannel := make(chan int)
	go testWritePerformance(path.Join(dir,"bonniego"), fileSize, randomBlock, bytesWrittenChannel)
	bytesWritten :=  <-bytesWrittenChannel


	finish := time.Now()
	duration := finish.Sub(start)
	fmt.Printf("wrote %d MiB\n", bytesWritten)
	fmt.Printf("took %f seconds\n", duration.Seconds())
	fmt.Printf("throughput %0.2f MiB/s\n", float64(bytesWritten)/float64(duration.Seconds())/math.Exp2(20))

	f, err := os.Open(path.Join(dir, "bonniego"))
	check(err)
	defer f.Close()

	bytesRead := 0
	data := make([]byte, Blocksize)

	start = time.Now()

	for {
		n, err := f.Read(data)
		bytesRead += n
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
	}

	finish = time.Now()
	duration = finish.Sub(start)

	fmt.Printf("read %d MiB\n", bytesRead >> 20)
	fmt.Printf("took %f seconds\n", duration.Seconds())
	fmt.Printf("throughput %0.2f MiB/s\n", float64(bytesWritten)/float64(duration.Seconds())/math.Exp2(20))

	if ! bytes.Equal(randomBlock, data) {
		panic("last block didn't match")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
