package service

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const paraCount int = 15

var p parameter_t

type parameter_t struct {
	min     [15]float64
	max     [15]float64
	weights [15]float64
	bias    float64
}

func LoadParameters() error {
	//Reads all trained parameters returns error as string if any
	bytes, err := ioutil.ReadFile("ml/model.txt")
	if err != nil {
		return err
	}

	data := string(bytes)
	eachArr := strings.Split(data, "\n\n")

	c := 1
	for i, e := range strings.Split(eachArr[0], " ") {
		if c > 15 {
			break
		}
		c++
		p.max[i], err = strconv.ParseFloat(e, 64)
	}
	c = 1
	for i, e := range strings.Split(eachArr[1], " ") {
		if c > 15 {
			break
		}
		c++
		p.min[i], err = strconv.ParseFloat(e, 64)
	}
	c = 1
	for i, e := range strings.Split(eachArr[2], " ") {
		if c > 15 {
			break
		}
		c++
		p.weights[i], err = strconv.ParseFloat(e, 64)
	}
	p.bias, err = strconv.ParseFloat(eachArr[3], 64)

	return nil
}

func formatData(r *http.Request) [15]float64 {
	var arr [15]float64
	if r.FormValue("sex") == "male" {
		arr[0] = 1
	}
	arr[1], _ = strconv.ParseFloat(r.FormValue("age"), 64)
	arr[2], _ = strconv.ParseFloat(r.FormValue("education"), 64)
	if r.FormValue("smoker") == "yes" {
		arr[3] = 1
		arr[4], _ = strconv.ParseFloat(r.FormValue("cigsPerDay"), 64)
	}
	if r.FormValue("bpmed") == "yes" {
		arr[5] = 1
	}
	if r.FormValue("stroke") == "yes" {
		arr[6] = 1
	}
	if r.FormValue("hyp") == "yes" {
		arr[7] = 1
	}
	if r.FormValue("diabetes") == "yes" {
		arr[8] = 1
	}

	arr[9], _ = strconv.ParseFloat(r.FormValue("chol"), 64)
	arr[10], _ = strconv.ParseFloat(r.FormValue("sysBP"), 64)

	arr[11], _ = strconv.ParseFloat(r.FormValue("diaBP"), 64)

	h, _ := strconv.ParseFloat(r.FormValue("height"), 64)
	h = h / 100
	w, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
	arr[12] = w / (h * h)

	arr[13], _ = strconv.ParseFloat(r.FormValue("hrate"), 64)
	arr[14], _ = strconv.ParseFloat(r.FormValue("glucose"), 64)
	return arr
}

func GetPrediction(r *http.Request) float64 {
	arr := formatData(r)
	//Scale data
	for i := 0; i < paraCount; i++ {
		d := 1.0
		if p.min[i] != p.max[i] {
			d = (p.max[i] - p.min[i])
		}
		arr[i] = (arr[i] - p.min[i]) / d
	}

	//Predict probabilty a=W*X+b
	output := 0.0
	for i := 0; i < paraCount; i++ {
		output += arr[i] * p.weights[i]
	}
	return output + p.bias
}
