package main

import (
  "flag"
  "github.com/golang/glog"

  "bufio"
  "errors"
  "fmt"
  "github.com/howeyc/fsnotify"
  "io"
  "os"
  "os/signal"
  "path"
  "time"
)

func tail(checkpointPath string, topic *LogTopic, watcher *fsnotify.Watcher) {
  cpFile, err := OpenCheckPointFile(path.Join(checkpointPath, topic.checkpointName))
  if err != nil {
    glog.Fatal(err)
    os.Exit(4)
  }
  defer cpFile.Close()
  checkpoint, err := cpFile.Retrieve()
  if err != nil {
    glog.Fatal("checkpoint Error")
  }

  for {
    select {
    case ev := <-watcher.Event:
      if ev.IsCreate() || (ev.IsModify() && !ev.IsAttrib()) {
        file, err := os.Open(ev.Name)

        if err != nil {
          glog.Fatal(err)
        }
        defer file.Close()

        if ev.Name == checkpoint.FileName {
          file.Seek(checkpoint.Position, 0)
        }
        r := bufio.NewReader(file)
        line, isPrefix, err := r.ReadLine()
        for err == nil && !isPrefix {
          s := string(line)
          fmt.Println(s)
          if checkpoint.Position, err = file.Seek(0, os.SEEK_CUR); err != nil {
            glog.Error(err)
          }
          checkpoint.FileName = ev.Name
          cpFile.Save(checkpoint)
          line, isPrefix, err = r.ReadLine()
        }
        if isPrefix {
          glog.Error(errors.New("buffer size to small"))
          return
        }
        if err != io.EOF {
          glog.Error(err)
          return
        }
        time.Sleep(time.Second)
      }
    case err := <-watcher.Error:
      glog.Fatal(err)
    }
  }
}

/**
 * Gotail works following the steps below
 * 1. read config file
 * 2. assign one watcher to each topic
 * 3. each goroutine do : consume dir event and then send it to kafka
 */
func main() {
  // parse config file path from command line
  configFilePtr := flag.String("conf", "config.ini", "Config file name")
  flag.Parse()

  // read config file, the config file is ini format
  config, err := ParseFromIni(*configFilePtr)
  if err != nil {
    glog.Fatal(err)
    os.Exit(1)
  }

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt)

  topics := config.topics
  for _, topic := range topics {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
      glog.Fatal(err)
      os.Exit(2)
    }

    err = watcher.Watch(topic.logBasePath)
    if err != nil {
      glog.Fatal(err)
      os.Exit(3)
    }

    go tail(config.checkpointPath, topic, watcher)
  }

  for _ = range interrupt {
    os.Exit(0)
  }
}
