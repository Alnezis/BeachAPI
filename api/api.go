package api

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for s, v := range m {
			result[s] = v
		}
	}
	return result
}

func CheckErrInfo(err error, info string) {
	if err != nil {
		fmt.Println("------------------------")
		log.Println(err)
		color.Yellow(info)
		fmt.Println("------------------------")
	}
}

// Round return rounded version of x with prec precision.
//
// Special cases are:
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func FormatFloat(v float64) string {
	return ToString(Round(v, 2))
}

func Print(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	CheckErrInfo(err, "MarshalIndent api.Print")
	text := string(s)
	fmt.Println(text)
	return text
}

func FormatInt(v int) string {
	sign := ""

	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(int64(v%1000), 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ".")
}

func Tips() string {
	s := "\n\n 💡 "
	r := RandInt(0, 100)
	fmt.Println(" - ", r)
	if 50 > r && r > 0 {
		return s + " "

	} else if 76 > r && r > 51 {
		return s + "2"

	} else if 100 > r && r > 77 {
		return s + "3"
	}
	return s + "err"
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func IsNum(a interface{}) bool {
	var count int
	if _, err2 := strconv.Atoi(ToString(a)); err2 != nil {
		return false
	}
	count, _ = strconv.Atoi(ToString(a))
	if count <= 0 {
		return false
	}
	return true
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Delete the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// ToString converts interface{} to string
func ToString(i1 interface{}) string {
	if i1 == nil {
		return ""
	}
	switch i2 := i1.(type) {
	default:
		return fmt.Sprint(i2)
	case bool:
		if i2 {
			return "true"
		} else {
			return "false"
		}
	case string:
		return i2
	case *bool:
		if i2 == nil {
			return ""
		}
		if *i2 {
			return "true"
		} else {
			return "false"
		}
	case *string:
		if i2 == nil {
			return ""
		}
		return *i2
	case *json.Number:
		return i2.String()
	case json.Number:
		return i2.String()
	}
	return ""
}

//125
func DeclOfNum(number int, titles []string) string {
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return ToString(number) + " " + titles[currentCase]
}

func FormatTime(sec int64) string {
	if sec < 0 {
		return "∞"
	}
	if sec < 60 {
		return ToString(sec) + " сек."
	} else if sec < 60*60 {
		return ToString(sec/60) + " мин."
	} else if sec < 60*60*24 {
		return ToString(sec/3600) + " ч."
	} else {
		return ToString(sec/3600/24) + " дн."
	}
}

func RandString(i int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, i)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Id struct {
	User int
	Chat int
}

type EditID struct {
	PeerID int
	MsgID  int
}
