package dirmap

import (
    "compress/zlib"
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    ReadBufferSize = 64 * 1024
)

type FileData struct {
    Size    int64
    ModTime time.Time
}

func GetChanges(mapPath, dbPath string) ([]string, error) {
    prevMap, err := loadDb(dbPath)
    if err != nil {
        return nil, err
    }

    changeList := make([]string, 0)

    filepath.Walk(mapPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return err
        }

        sPath := strings.Replace(path, mapPath, "", -1)

        prevInfo, ok := prevMap[sPath]
        if !ok {
            // new
            changeList = append(changeList, sPath)
            prevMap[sPath] = &FileData{info.Size(), info.ModTime()}
        } else if prevInfo.Size != info.Size() || prevInfo.ModTime != info.ModTime() {
            // changed
            changeList = append(changeList, sPath)
            prevMap[sPath] = &FileData{info.Size(), info.ModTime()}
        }

        return nil
    })

    saveDb(dbPath, prevMap)

    return changeList, nil
}

func loadDb(dbPath string) (map[string]*FileData, error) {
    data := make(map[string]*FileData)

    f, err := os.Open(dbPath)
    if err != nil {
        return data, saveDb(dbPath, data)
    }
    defer f.Close()

    decomp, err := zlib.NewReader(f)
    if err != nil {
        return nil, err
    }
    defer decomp.Close()

    js := json.NewDecoder(decomp)

    err = js.Decode(&data)
    return data, err
}

func saveDb(dbPath string, data map[string]*FileData) error {
    f, err := os.Create(dbPath)
    if err != nil {
        return err
    }
    defer f.Close()

    comp := zlib.NewWriter(f)
    defer comp.Close()

    js := json.NewEncoder(comp)
    return js.Encode(data)
}
