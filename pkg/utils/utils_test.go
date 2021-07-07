package utils

import (
    "github.com/stretchr/testify/assert"
    "math"
    "path/filepath"
    "strings"
    "testing"
)

func TestIn(t *testing.T) {
    assert.Equal(t, false, In(25, 37, 18))
    assert.Equal(t, true, In(25, 25, 18))
    assert.Equal(t, true, In(25, 18, 25))
    assert.Equal(t, true, In(25, 25))
    assert.Equal(t, true, In(filepath.Ext("my.csv"), ".json", ".csv"))

    assert.Equal(t, true, In(math.Inf(1), math.Inf(-1), math.Inf(1), math.NaN()))
    assert.Equal(t, true, In(math.Inf(-1), math.Inf(-1), math.Inf(1), math.NaN()))
    assert.Equal(t, true, In(math.NaN(), math.Inf(-1), math.Inf(1), math.NaN()))

    var mySlice []interface{}
    mySlice = append(mySlice, "A", "B")
    assert.Equal(t, false, In("a", mySlice...))
    assert.Equal(t, false, In("a", mySlice))
    assert.Equal(t, true, In("B", mySlice...))
    assert.Equal(t, false, In("B", mySlice)) // 👈 be careful

    type Point struct {
        x int
        y int
    }
    assert.Equal(t, true, In(Point{3, 5}, Point{1, 2}, Point{3, 5}))
    assert.Equal(t, false, In(Point{3, 5}, []Point{Point{1, 2}, Point{3, 5}})) // 👈
    assert.Equal(t, true, In(Point{3, 5}, []interface{}{Point{1, 2}, Point{3, 5}}...))

    assert.Equal(t, false, In("c", strings.Split("a,b,c", ","))) // 👈
    var mySlice2 []interface{}
    for _, v := range strings.Split("a,b,c", ",") {
        mySlice2 = append(mySlice2, v)
    }
    assert.Equal(t, true, In("c", mySlice2...))
    assert.Equal(t, false, In("c", []interface{}{strings.Split("a,b,c", ",")}...)) // 👈
}
