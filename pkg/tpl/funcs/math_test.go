package funcs

import (
    "errors"
    "fmt"
    "github.com/CarsonSlovoka/excel/pkg/utils"
    "github.com/spf13/cast"
    "github.com/stretchr/testify/assert"
    "math"
    "reflect"
    "testing"
)

var anError error

func init() {
    anError = errors.New("an error")
}

func check(t *testing.T, testStruct interface{}, expected, actual interface{}, err error, extraMsg ...interface{}) {
    var ok bool
    if err != nil {
        if reflect.ValueOf(expected).Kind() == reflect.Ptr {
            ok = assert.Error(t, expected.(error), extraMsg...)
        } else {
            ok = assert.Equal(t, expected, err, extraMsg...)
        }
    } else {
        ok = assert.Equal(t, expected, actual, extraMsg...)
    }
    if !ok {
        fmt.Printf("%+v\n", testStruct)
    }
}

func roundDown(a interface{}, n int) float64 {
    // we compare only "n" digits behind point if its a real float
    // otherwise we usually get different float values on the last positions

    af, err := cast.ToFloat64E(a)
    if err != nil {
        panic(err)
    }
    return float64(int(af*math.Pow10(n))) / 10000
}

func TestBasicArithmetic(t *testing.T) {
    t.Parallel()

    for _, test := range []struct {
        fn             func(a, b interface{}) (interface{}, error)
        a              interface{}
        b              interface{}
        expect         interface{}
        extraMsgForErr interface{}
    }{
        {Add, 5, 3, float64(8), "Add"},
        {Add, 1.0, "foo", anError, "Add"},
        {Sub, 4, 2, float64(2), "Sub"},
        {Sub, 1.0, "foo", anError, "Sub"},
        {Mul, 4, 2, float64(8), "Mul"},
        {Mul, 1.0, "foo", anError, "Mul"},
        {Div, 4, 2, float64(2), "Div"},
        {Div, 1.0, "foo", anError, "Div"},
    } {
        actual, err := test.fn(test.a, test.b)
        check(t, test, test.expect, actual, err, test.extraMsgForErr)
    }
}

func TestCeil(t *testing.T) {
    t.Parallel()
    for _, test := range []struct {
        x      interface{}
        expect interface{}
    }{
        {0.1, 1.0},
        {0.5, 1.0},
        {1.1, 2.0},
        {1.5, 2.0},
        {-0.1, 0.0},
        {-0.5, 0.0},
        {-1.1, -1.0},
        {-1.5, -1.0},
        {"abc", errors.New("ceil operator can't be used with non-float value")},
    } {
        actual, err := Ceil(test.x)
        check(t, test, test.expect, actual, err)
    }
}

func TestFloor(t *testing.T) {
    t.Parallel()
    for _, test := range []struct {
        x      interface{}
        expect interface{}
    }{
        {0.1, 0.0},
        {0.5, 0.0},
        {1.1, 1.0},
        {1.5, 1.0},
        {-0.1, -1.0},
        {-0.5, -1.0},
        {-1.1, -2.0},
        {-1.5, -2.0},
        {"abc", errors.New("floor operator can't be used with non-float value")},
    } {

        actual, err := Floor(test.x)
        check(t, test, test.expect, actual, err)
    }
}

func TestLog(t *testing.T) {
    t.Parallel()
    for _, test := range []struct {
        a      interface{}
        expect interface{}
    }{
        {1, float64(0)},
        {3, 1.0986}, // default float64
        {1.0, float64(0)},
        {3.1, 1.1314},
        {"abc", errors.New("log operator can't be used with non integer or float value")},

        // 👇 special test
        {math.Inf(1), math.Inf(1)}, // Log(+Inf) = +Inf
        {0, math.Inf(-1)},          // Log(0) = -Inf
        {-1, math.NaN()},           // Log(x < 0) = NaN
        {math.NaN(), math.NaN()},   // Log(NaN) = NaN
    } {

        actual, err := Log(test.a)

        if !utils.In(actual, math.Inf(-1), math.Inf(1), math.NaN()) {
            actual = roundDown(actual, 4)
        }

        expectValue := reflect.ValueOf(test.expect)
        actualValue := reflect.ValueOf(actual)
        if actualValue.Kind() == reflect.Float64 && math.IsNaN(actualValue.Float()) &&
            expectValue.Kind() == reflect.Float64 && math.IsNaN(expectValue.Float()) {
            assert.Equal(t, true, math.IsNaN(actual))
            continue
        }
        check(t, test, test.expect, actual, err)
    }
}

func TestSqrt(t *testing.T) {
    t.Parallel()

    for _, test := range []struct {
        a      interface{}
        expect interface{}
    }{
        {81, float64(9)},
        {0.25, 0.5},
        {0, float64(0)},
        {"abc", errors.New("sqrt operator can't be used with non integer or float value")},
        {-1, math.NaN()},
    } {

        actual, err := Sqrt(test.a)

        if test.a == -1 {
            // Separate test for Sqrt(-1) -- returns NaN
            assert.True(t, math.IsNaN(actual))
            continue
        }

        if actual != math.Inf(-1) {
            actual = roundDown(actual, 4)
        }

        check(t, test, test.expect, actual, err)
    }
}

func TestMod(t *testing.T) {
    t.Parallel()
    for _, test := range []struct {
        a      interface{}
        b      interface{}
        expect interface{}
    }{

        {3, 2, int64(1)},
        {3, 1, int64(0)},
        {0, 3, int64(0)},
        {3.1, 2, int64(1)},
        {3, 2.1, int64(1)},
        {3.1, 2.1, int64(1)},
        {int8(3), int8(2), int64(1)},
        {int16(3), int16(2), int64(1)},
        {int32(3), int32(2), int64(1)},
        {int64(3), int64(2), int64(1)},
        {"3", "2", int64(1)},

        {3, 0, errors.New("the number can't be divided by zero at modulo operation")},
        {"3.1", "2", errors.New("modulo operator can't be used with non integer value")},
        {"aaa", "0", errors.New("modulo operator can't be used with non integer value")},
        {"3", "aaa", errors.New("modulo operator can't be used with non integer value")},
    } {

        actual, err := Mod(test.a, test.b)
        check(t, test, test.expect, actual, err)
    }
}

func TestModBool(t *testing.T) {
    t.Parallel()

    for _, test := range []struct {
        a      interface{}
        b      interface{}
        expect interface{}
    }{
        {3, 3, true},
        {3, 2, false},
        {3, 1, true},
        {0, 3, true},
        {3.1, 2, false},
        {3, 2.1, false},
        {3.1, 2.1, false},
        {int8(3), int8(3), true},
        {int8(3), int8(2), false},
        {int16(3), int16(3), true},
        {int16(3), int16(2), false},
        {int32(3), int32(3), true},
        {int32(3), int32(2), false},
        {int64(3), int64(3), true},
        {int64(3), int64(2), false},
        {"3", "3", true},
        {"3", "2", false},

        {3, 0, errors.New("the number can't be divided by zero at modulo operation")},
        {"3.1", "2", errors.New("modulo operator can't be used with non integer value")},
        {"aaa", "0", errors.New("modulo operator can't be used with non integer value")},
        {"3", "aaa", errors.New("modulo operator can't be used with non integer value")},
    } {

        actual, err := ModBool(test.a, test.b)
        check(t, test, test.expect, actual, err)
    }
}

func TestRound(t *testing.T) {
    t.Parallel()
    for _, test := range []struct {
        x      interface{}
        expect interface{}
    }{
        {0.1, 0.0},
        {0.5, 1.0},
        {1.1, 1.0},
        {1.5, 2.0},
        {-0.1, -0.0},
        {-0.5, -1.0},
        {-1.1, -1.0},
        {-1.5, -2.0},
        {"abc", errors.New("round operator can't be used with non-float value")},
    } {

        actual, err := Round(test.x)
        check(t, test, test.expect, actual, err)
    }
}

func TestPow(t *testing.T) {
    t.Parallel()

    for _, test := range []struct {
        a      interface{}
        b      interface{}
        expect interface{}
    }{
        {0, 0, float64(1)},
        {2, 0, float64(1)},
        {2, 3, float64(8)},
        {-2, 3, float64(-8)},
        {2, -3, 0.125},
        {-2, -3, -0.125},
        {0.2, 3, 0.008},
        {2, 0.3, 1.2311},
        {0.2, 0.3, 0.617},
        {"aaa", "3", errors.New("pow operator can't be used with non-float value")},
        {"2", "aaa", errors.New("pow operator can't be used with non-float value")},
    } {
        actual, err := Pow(test.a, test.b)
        actual = roundDown(actual, 4)
        check(t, test, test.expect, actual, err)
    }
}
