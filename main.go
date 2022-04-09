package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	// CPU pprof
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	dir := "testdata"
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create("out.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	c := make(chan File, 1)
	go func() {
		defer close(c)
		for _, f := range files {
			c <- NewFile(dir, f)
		}
	}()
	err = readWriteAsync(out, c)
	if err != nil {
		fmt.Println(err)
	}

	// Memory pprof
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	os.Exit(0)
}

type File struct {
	filename string
}

func NewFile(dir string, e os.DirEntry) File {
	fmt.Println("NewFile", e.Name())
	return File{
		filename: filepath.Join(dir, e.Name()),
	}
}

func (f File) Open() (*os.File, error) {
	return os.Open(f.filename)
}

func (f File) Name() string {
	return f.filename
}

func readWriteAsync(output io.Writer, files <-chan File) error {
	var i int
	for file := range files {
		if err := openAndCopyFile(file, output); err != nil {
			return err
		}
		i++
		fmt.Printf("%d: %s\n", i, file.Name())
	}

	return nil
}

func openAndCopyFile(file File, w io.Writer) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()
	_, err = io.Copy(w, fileReader)
	return err
}
