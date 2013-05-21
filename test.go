// https://groups.google.com/forum/m/?fromgroups#!topic/golang-nuts/zzj8Kgt0TpQ
// Other:
// https://groups.google.com/forum/m/?fromgroups#!topic/golang-nuts/bPy_s3yQTMs
// https://groups.google.com/forum/m/?fromgroups#!msg/golang-nuts/_VoZfniBTZE/QfSXW-iPoCsJ
// https://groups.google.com/forum/m/?fromgroups#!topic/golang-nuts/1fOUy40iEJc
// https://groups.google.com/forum/m/?fromgroups#!topic/golang-nuts/q-XD__2s6fE
// https://groups.google.com/forum/m/?fromgroups#!topic/golang-nuts/ILODEL0ldpw
// http://play.golang.org/p/AZ4OUMQT2l
// http://play.golang.org/p/OGTdl1WrT9
// http://golang.org/doc/play/

package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var i map[string]interface{}
	b := []byte(`{"a":{"b":12,"c":22,"d":[1,2,3,4,5]}}`)
	json.Unmarshal(b, &i)
	j := i["a"]
	k := j.(map[string]interface{})["b"]
	fmt.Println(k)

	r := map[string]string{}
	r["login"] = "jan"
	v := map[interface{}]interface{}{}
	v["req"] = req
	req := v["req"]
	fmt.Println(req.(map[string]string)["login"])
}
