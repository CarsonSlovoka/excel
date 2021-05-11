package file

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestCSV2Json(t *testing.T) {
    byteData := []byte(`
Name, Msg, number
AA, Hello world, 1
BB, cc, 4.5
`)
    memoryReader := bytes.NewReader(byteData)
    jsonObj, err := CSV2Json(memoryReader)
    if err != nil {
        panic(err)
    }
    jsonByte, err := json.Marshal(jsonObj)
    if err != nil {
        panic(err)
    }
    assert.Equal(t, `[{"Msg":"Hello world","Name":"AA","number":1},{"Msg":"cc","Name":"BB","number":4.5}]`, string(jsonByte))

    beautifulJsonByte, err := json.MarshalIndent(jsonObj, "", "  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(beautifulJsonByte))
}
