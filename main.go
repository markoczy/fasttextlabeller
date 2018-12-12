package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

/**
FastTextLabeller:

Quick and dirty implementation of a labeller for fasttext
scans for files in input dir (specified by -in) and
writes the labelled content to the output file (specified by -out)
*/

// checks for error and panics
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// discard ignores variable
func discard(interface{}) {}

// task single entity of work
type task struct {
	file, label string
}

func initTasksRecursive(path, label string, tasks []task) []task {
	files, err := ioutil.ReadDir(path)
	check(err)

	for _, file := range files {
		switch {
		case file.IsDir():
			tasks = initTasksRecursive(path+"/"+file.Name(), label, tasks)
		case strings.HasSuffix(strings.ToLower(file.Name()), ".txt"):
			tasks = append(tasks, task{file: path + "/" + file.Name(), label: label})
		}

	}

	return tasks
}

func initTasks(rootDir string) []task {
	// ioutil.WriteFile()
	files, err := ioutil.ReadDir(rootDir)
	check(err)

	tasks := []task{}

	// top level = label
	for _, file := range files {
		if file.IsDir() {
			label := "__label__" + file.Name()
			fmt.Println("Found label:", label)
			tasks = initTasksRecursive(rootDir+"/"+file.Name(), label, tasks)
		}
	}
	return tasks
}

func shuffle(tasks []task) []task {
	count := len(tasks)

	ret := make([]task, count)

	perm := rand.Perm(count)
	for i := 0; i < count; i++ {
		ret[i] = tasks[perm[i]]
	}

	return ret
}

func appendText(filename, text string) {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func processTasks(tasks []task, path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		os.Remove(path)
	}
	os.Create(path)

	rx, err := regexp.Compile("(\\r?\\n)")
	check(err)

	for _, task := range tasks {
		data, err := ioutil.ReadFile(task.file)
		check(err)

		str := string(data)
		str = rx.ReplaceAllString(str, " ")

		appendText(path, task.label+"   "+str+"\n")
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())
	in := flag.String("in", "", "path to input directory")
	out := flag.String("out", "", "path to output directory")
	flag.Parse()

	if *in == "" || *out == "" {
		flag.Usage()
		return
	}

	tasks := initTasks(*in)
	fmt.Println("InitTasks returns:", tasks)

	tasks = shuffle(tasks)

	processTasks(tasks, *out)

}
