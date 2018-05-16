package main

import (
	"fmt"
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"os"
	"bufio"
)
const FILENAME = "small.in"

func ArraySource(num ...int) <-chan int{
	out := make(chan int)

	go func() {
		for _, v := range num{
			out <- v
		}
		close(out)
	}()

	return out
}

func InMemorySort(in <-chan int) <-chan int{
	out := make(chan int)

	go func() {

		// build arr in memory
		arr := 	[]int{}
		for num := range in {
			arr = append(arr, num)
		}

		// sort
		sort.Ints(arr)

		// output
		for _, n := range arr{
			out <- n
		}

		close(out)

	}()

	return out

}

func Merge(in1, in2 <- chan int) <-chan int {
	out := make(chan int)

	// 开启 routine
	go func() {

		v1, ok1 := <-in1
		v2, ok2 := <-in2
		// ok1 or ok2 is
		// true
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2){
				out <- v1
				v1, ok1 = <-in1
			}else{
				out <- v2
				v2, ok2 = <-in2
			}
		}

		close(out)
	}()

	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <- chan int {
	out := make(chan int)

	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			// read buffer for each 8
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}

			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}

		close(out)
	}()


	return out

}

func WriteSink(writer io.Writer, in <-chan int){

	// write buffer
	for v := range in{
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}

}

func RandomSource (count int) <- chan int {
	out := make(chan int)

	go func() {
		// generate random number
		for i := 0; i < count; i++ {
			out <- rand.Intn(10000)
		}

		close(out)
	}()

	return out
}


// 归并排序
func memoryMergeSortDemo() {

	p := Merge(InMemorySort(ArraySource(6, 4, 3, 7, 5)), InMemorySort(ArraySource(1, 9, 2, 7, 8)))

	/*for {
		if num, ok := <- p; ok{
			fmt.Println(num)
		}else {
			break
		}
	}*/

	for num := range p {
		fmt.Println(num)
	}

}

func fileMergeSortDemo(){

	const COUNT = 64

	file, err := os.Create(FILENAME)
	if err != nil {
		panic(err)
	}

	//
	defer file.Close()

	p := RandomSource(COUNT)
	writer := bufio.NewWriter(file)
	WriteSink(writer, p)
	defer writer.Flush()

	file, err = os.Open(FILENAME)
	if err != nil {
		panic(err)
	}

	p = ReaderSource(bufio.NewReader(file), -1)
	count := 0

	for v := range p {
		fmt.Println(v)
		count ++
		if count > 100 {
			break
		}
	}

}

func MergeN(in ... <- chan int) <- chan int{

	if len(in) == 1{
		return in[0]
	}

	m := len(in) / 2

	return Merge(MergeN(in[:m]...), MergeN(in[m:]...))

}

func externalSort(){

	p := createPipeline(FILENAME, 512, 4)
	writeToFile(p, "small.out")
	printFile("small.out")

}

func createPipeline(fileName string, fileSize, chunkCount int) <- chan int{
	chunkSize := fileSize / chunkCount
	sortResult := []<-chan int{}
	for i := 0; i<chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize), 0)
		source := ReaderSource(bufio.NewReader(file), chunkSize)
		sortResult = append(sortResult, InMemorySort(source))

	}
	return MergeN(sortResult...)
}

func writeToFile(in <-chan int, fileName string){

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer  writer.Flush()

	// 先进后出，所以先close 再flush

	WriteSink(writer, in)

}

func printFile(fileName string){

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	p := ReaderSource(bufio.NewReader(file), -1)

	for v := range p {
		fmt.Println(v)
	}


}
func main_()  {
	fileMergeSortDemo()
	externalSort()
}