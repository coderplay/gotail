package main

import (
  "errors"
  "fmt"

  ini "github.com/coderplay/goini"
  "strconv"
  "strings"
)

/**
 *
 */
type LogTopic struct {
  topic          string
  logBasePath    string
  logFilePattern string
  encoding       string
  compressed     bool
  timeout        uint64
  checkpointName string
}

type Config struct {
  checkpointPath string
  topics         []LogTopic
}

const (
  formatStr = `topic=%s
log_base_path=%s
log_file_regex=%s
encoding=%s
compressed=%t
timeout=%d
checkpoint_name=%s`
)

func (item *LogTopic) String() string {
  return fmt.Sprintf(
    formatStr,
    item.topic,
    item.logBasePath,
    item.logFilePattern,
    item.encoding,
    item.compressed,
    item.timeout,
    item.checkpointName)
}

/**
 * Parse a Config instance from ini file
 */
func ParseFromIni(path string) (*Config, error) {
  dict, err := ini.Load(path)
  if err != nil {
    return nil, err
  }

  checkpointPath, ok := dict.GetString("system", "checkpoint_path")
  if !ok {
    return nil, errors.New("Invalid config format, misses checkpoint path")
  }

  // get topics from dict
  topics := make([]LogTopic, 1)
  for k, v := range dict {
    if strings.HasPrefix(k, "topic_") {
      topic, present := v["topic"]
      if !present {
        topic = k
      }

      logBasePath, present := v["log_base_path"]
      if !present {
        return nil, errors.New("Invalid config format, misses log_base_path in " + k)
      }

      logFilePattern := v["log_file_regex"]
      encoding := v["encoding"]

      c, present := v["compressed"]
      var compressed bool
      if !present || c == "n" || c == "N" ||
        c == "0" || c == "f" || c == "F" || c == "false" {
        compressed = false
      } else if c == "y" || c == "Y" ||
        c == "1" || c == "t" || c == "T" || c == "true" {
        compressed = true
      } else {
        return nil,
          errors.New("Invalid config format, invalid compressed value " + c)
      }

      timeoutStr, present := v["timeout"]
      var timeout uint64
      if !present {
        timeout = 10000
      } else {
        timeout, err = strconv.ParseUint(timeoutStr, 10, 0)
        if err != nil {
          return nil,
            errors.New("Invalid config format, invalid timeout value " + timeoutStr)
        }
      }

      checkpointName, present := v["checkpoint_name"]
      if !present {
        checkpointName = k
      }

      logTopic := &LogTopic{topic, logBasePath, logFilePattern,
        encoding, compressed, timeout, checkpointName}
      fmt.Println(logTopic.String())
      topics = append(topics, *logTopic)
    }
  }

  return &Config{checkpointPath, topics}, nil
}
