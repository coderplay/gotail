package main

import (
  "encoding/gob"
  "os"
)

type CheckPoint struct {
  FileName string
  Position int64
  // ModTime  uint64
}

type CheckPointFile struct {
  *os.File
  encoder *gob.Encoder
  decoder *gob.Decoder
}

func OpenCheckPointFile(path string) (file *CheckPointFile, err error) {
  f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
  if err != nil {
    return nil, err
  }

  return &CheckPointFile{f, gob.NewEncoder(f), gob.NewDecoder(f)}, nil
}

func (cpFile *CheckPointFile) Save(checkpoint *CheckPoint) error {
  cpFile.Truncate(0)
  if err := cpFile.encoder.Encode(checkpoint); err != nil {
    return err
  }
  return cpFile.Sync()
}

func (cpFile *CheckPointFile) Retrieve() (checkpoint *CheckPoint, err error) {
  checkpoint = &CheckPoint{}
  fi, err := cpFile.Stat()
  if err != nil {
    return nil, err
  } else if fi.Size() == 0 {
    return checkpoint, nil
  }

  if err := cpFile.decoder.Decode(checkpoint); err != nil {
    return nil, err
  }
  return checkpoint, nil
}
