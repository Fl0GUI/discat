package custom

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path"
)

const dir = "assets/custom"

func Breeds() []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	res := make([]string, len(entries))
	var i int
	var v os.DirEntry
	for i, v = range entries {
		res[i] = v.Name()
	}

	return res[:i+1]
}

func GetCat(breed string) fs.FileInfo {
	entries, err := os.ReadDir(path.Join(dir, breed))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	i := rand.Intn(len(entries))
	file, err := entries[i].Info()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return file
}

func Open(breed string, fi fs.FileInfo) *os.File {
  f, err := os.Open(path.Join(dir, breed, fi.Name()))
  if err != nil {
    fmt.Println(err)
    return nil
  }
  return f
}
