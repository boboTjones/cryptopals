package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func qsParse(qs string) []uint8 {
	out := make(map[string]interface{})
	for _, v := range strings.Split(qs, "&") {
		tmp := strings.Split(v, "=")
		out[tmp[0]] = tmp[1]
	}
	ret, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	return ret

}

// Assuming this is leading up to using JSON?
func profile_for(email string) string {
	var ret []string
	// eating & and =
	// will cause a lookup to fail.
	nopes := []string{"&", "="}
	for _, v := range nopes {
		email = strings.Replace(email, v, "", -1)
	}

	// ew
	tmp := make(map[string]interface{})
	tmp["email"] = email
	tmp["uid"] = 10
	tmp["role"] = "user"

	// ew
	for k, v := range tmp {
		switch x := v.(type) {
		case int:
			ret = append(ret, k+"="+strconv.Itoa(x))
		case float64:
			ret = append(ret, k+"="+strconv.FormatFloat(x, 'f', 2, 32))
		default:
			ret = append(ret, k+"="+v.(string))
		}
	}
	return strings.Join(ret, "&")
}

func main() {
	qp := qsParse("foo=bar&baz=qux&zap=zazzle")
	fmt.Printf("%q\n", qp)
	pq := profile_for("foo@bar.com&role=admin")
	fmt.Printf("%s\n", pq)
}
