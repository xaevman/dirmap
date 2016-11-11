package dirmap

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
)

const (
    testDir = "test"
    dbName  = "test.db"
)

type Test struct {
    Data    []FileInfo
    Results []string
}

type FileInfo struct {
    Name string
    Data string
}

var (
    tests = []Test{
        {
            Data: []FileInfo{
                {Name: filepath.Join(testDir, "file1"), Data: "my test1"},
                {Name: filepath.Join(testDir, "file2"), Data: "my test2"},
                {Name: filepath.Join(testDir, "file3"), Data: "my test3"},
            },
            Results: []string{
                fmt.Sprintf("%c%s", filepath.Separator, "file1"),
                fmt.Sprintf("%c%s", filepath.Separator, "file2"),
                fmt.Sprintf("%c%s", filepath.Separator, "file3"),
            },
        },
        {
            Data: []FileInfo{
                {Name: filepath.Join(testDir, "file1"), Data: "my test1.2"},
            },
            Results: []string{
                fmt.Sprintf("%c%s", filepath.Separator, "file1"),
            },
        },
        {
            Data: []FileInfo{
                {Name: filepath.Join(testDir, "file1"), Data: "my test1.2.3"},
                {Name: filepath.Join(testDir, "file3"), Data: "my test3.2"},
            },
            Results: []string{
                fmt.Sprintf("%c%s", filepath.Separator, "file1"),
                fmt.Sprintf("%c%s", filepath.Separator, "file3"),
            },
        },
    }
)

func TestMain(m *testing.M) {
    err := cleanup()
    if err != nil {
        panic(err)
    }

    err = os.MkdirAll(testDir, 0770)
    if err != nil {
        panic(err)
    }

    os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
    for i := range tests {
        err := runTest(tests[i].Data, tests[i].Results)
        if err != nil {
            t.Fatalf("Test %d error: %v", i, err)
        }
    }
}

func runTest(data []FileInfo, result []string) error {
    for i := range data {
        err := ioutil.WriteFile(data[i].Name, []byte(data[i].Data), 0660)
        if err != nil {
            return err
        }
    }

    changes, err := GetChanges(testDir, dbName)
    if err != nil {
        return fmt.Errorf("%v", err)
    }

    if len(changes) != len(result) {
        return fmt.Errorf("initial change map size mismatch (%d != %d)", len(changes), len(result))
    }

    for i := range result {
        if result[i] != changes[i] {
            return fmt.Errorf("result[%d] != changes[%d] (%s != %s)", i, i, result[i], changes[i])
        }
    }

    return nil
}

func cleanup() error {
    err := os.RemoveAll(testDir)
    if err != nil {
        return err
    }
    err = os.RemoveAll(dbName)
    if err != nil {
        return err
    }

    return nil
}
