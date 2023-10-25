package main

import (
	"flag"
	"fmt"
	"github.com/murakmii/gj"
	"github.com/murakmii/gj/vm"
	_ "github.com/murakmii/gj/vm/native"
	"os"
	"strings"
	"time"
)

var (
	configPath string
	mainClass  string
	print      bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "path of configuration file")
	flag.StringVar(&mainClass, "main", "", "main class name")
	flag.BoolVar(&print, "print", false, "print disassembled class file")
}

func main() {
	flag.Parse()
	if len(configPath) == 0 || len(mainClass) == 0 {
		flag.Usage()
		return
	}
	mainClass = strings.ReplaceAll(mainClass, ".", "/")

	config, err := readConfig()
	if err != nil {
		fmt.Printf("failed to read config: %s", err)
		return
	}

	classPaths, err := gj.InitClassPaths(config.ClassPath)
	if err != nil {
		fmt.Printf("failed to init class path: %s", err)
		return
	}
	defer func() {
		for _, classPath := range classPaths {
			classPath.Close()
		}
	}()

	if print {
		execPrint(classPaths)
	} else {
		execVM(config)
	}
}

func readConfig() (*gj.Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return gj.ReadConfig(f)
}

func execPrint(classPaths []gj.ClassPath) {
	for _, classPath := range classPaths {
		classFile, err := classPath.SearchClass(mainClass + ".class")
		if err != nil {
			fmt.Printf("failed to search class: %s", err)
			return
		}

		if classFile == nil {
			continue
		}

		fmt.Printf(classFile.String())
		return
	}

	fmt.Println("class not found")
}

func execVM(config *gj.Config) {
	start := time.Now().UnixMilli()
	vmInstance, err := vm.InitVM(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("-> VM initialized!(%d ms)\n", time.Now().UnixMilli()-start)
	fmt.Printf("-> Loaded classes: %d\n", vmInstance.ClassCacheNum())
	fmt.Printf("-> Execute main method...\n")
	fmt.Println("--------------------------------------")

	if err := vmInstance.ExecMain(mainClass, []string{}); err != nil {
		panic(err)
	}

	for {
		result, ok := <-vmInstance.Executor().Wait()
		if !ok {
			break
		}

		if result.Err != nil {
			fmt.Printf("[VM] occurred error in thread '%s': %s\n", result.Thread.Name(), result.Err)
		}
	}

	fmt.Println("--------------------------------------")
	fmt.Println("Finished all non-daemon threads")
}
