package main

import (
  "flag"
  "log"

  "bufio"
  "errors"
  "fmt"
  "github.com/howeyc/fsnotify"
  "io"
  "os"
  "os/signal"
  "time"
)

func tail(checkpointPath string, topic *LogTopic, watcher *fsnotify.Watcher) {
  checkpoint := &CheckPoint{}
  for {
    select {
    case ev := <-watcher.Event:
      if ev.IsCreate() || (ev.IsModify() && !ev.IsAttrib()) {
        f, err := os.Open(ev.Name)
        if err != nil {
          log.Fatal("Error")
        }
        defer f.Close()
        r := bufio.NewReader(f)
        line, isPrefix, err := r.ReadLine()
        for err == nil && !isPrefix {
          s := string(line)
          fmt.Println(s)
          line, isPrefix, err = r.ReadLine()
        }
        if isPrefix {
          fmt.Println(errors.New("buffer size to small"))
          return
        }
        if err != io.EOF {
          fmt.Println(err)
          return
        }
        time.Sleep(time.Second)
      }
    case err := <-watcher.Error:
      log.Println("error:", err)
    }
  }
}

/**
 * 1. read config file
 * 2. each group use one goroutine and one channel
 * 3. each coroutine do : watch dir event and consume the event, send
 */
func main() {
  // parse config file path from command line
  configFilePtr := flag.String("conf", "config.ini", "Config file name")
  flag.Parse()

  // read config file, the config file is ini format
  config, err := ParseFromIni(*configFilePtr)
  if err != nil {
    log.Fatal(err)
  }

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt)

  topics := config.topics
  for _, topic := range topics {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
      log.Fatal(err)
    }

    err = watcher.Watch(topic.logBasePath)
    if err != nil {
      log.Fatal(err)
    }

    go tail(config.checkpointPath, topic, watcher)
  }

  for _ = range interrupt {
    os.Exit(0)
  }
}
