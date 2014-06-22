package main

type CheckPoint struct {
  FileName string
  Position uint64
  ModTime  uint64
}

func (cp *CheckPoint) Save() error {
  return nil
}
