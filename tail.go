package main

import "flag"

/**
 * 1. read config file
 * 2. each group use one coroutine and one channel
 * 3. each coroutine do : watch dir event and consume the event, send
 */
func main() {
  // parse config file path from command line
  configFilePtr := flag.String("conf", "config.ini", "Config file name")
  flag.Parse()

  // read config file, the config file is ini format
  ParseFromIni(*configFilePtr)

  //
  //	watcher, err := fsnotify.NewWatcher()
  //	if err != nil {
  //		log.Fatal(err)
  //	}
  //
  //	done := make(chan bool)
  //
  //	// Process events
  //	go func() {
  //		for {
  //			select {
  //			case ev := <-watcher.Event:
  //				log.Println("event:", ev)
  //			case err := <-watcher.Error:
  //				log.Println("error:", err)
  //			}
  //		}
  //	}()
  //
  //	err = watcher.Watch("testDir")
  //	if err != nil {
  //		log.Fatal(err)
  //	}
  //
  //	<-done
  //
  //	/* ... do stuff ... */
  //	watcher.Close()
}
