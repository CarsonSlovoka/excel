package file

import (
    "encoding/csv"
    "encoding/json"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "io"
    "strconv"
    "strings"
)

func CSV2Json(file io.Reader) ([]map[string]interface{}, error) {
    csvReader := csv.NewReader(file)

    var jsonData []string
    headers := make([]string, 0)
    formatStr := &utils.FString{}
    formatStr.MustCompile(`"{{.Field}}":{{.Value}},`, nil)
    for i := 0; ; i++ {
        row, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if i == 0 {
            for _, val := range row {
                headers = append(headers, strings.Trim(val, " ")) // remove the left and right spaces so that you can use `Name, Msg , Other`
            }
            continue
        }

        var value interface{}
        for idx, valString := range row {
            valString = strings.Trim(valString, " ")
            if number, err := strconv.ParseFloat(valString, 64); err == nil {
                value = number
            } else {
                value = `"` + valString + `"`
            }
            if err = formatStr.Render(utils.Context{
                "Field": headers[idx], "Value": value,
            }); err != nil {
                panic(err)
            }
        }
        formatStr.Data = utils.RFindAndReplace(formatStr.Data, ",", "") // remove the last comma
        formatStr.Data = "{" + formatStr.Data + "}"
        jsonData = append(jsonData, formatStr.Clear())
    }
    jsonString := strings.Join(jsonData, ",")
    jsonString = "[" + jsonString + "]"

    var dataArray []map[string]interface{}
    if err := json.Unmarshal([]byte(jsonString), &dataArray); err != nil {
        return nil, err
    }
    return dataArray, nil

}
