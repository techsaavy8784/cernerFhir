package common

import (
	//"errors"
	"fmt"
	"math"
	"strings"
	//"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// func LinesPerPage() int64 {
// 	pageSize, err := strconv.ParseInt( m.ActiveConfig().Env("page_size"), 10, 64)
// 	if err != nil {
// 		pageSize = 20
// 	}
// 	return pageSize
// }
func StringDateToFhirDate(date string)( time.Time, error) {
	Layout := "2006-01-02T15:04:05.000Z"
	//d, err := time.Parse(time.RFC3339,date)
	d, err := time.Parse(Layout,date)
	if err != nil {
		return time.Time{}, err
	}
	return d, nil
}

func FhirDateToString(d time.Time, format string) string {
	var date string
	switch format {
	case "full":
		date = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.000Z", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
	case "mdy":
		date = fmt.Sprintf("%02d-%02d-%4d", d.Month(), d.Day(), d.Year())
	case "ymd":
		date = fmt.Sprintf("%4d-%02d-%02d", d.Year(), d.Month(), d.Day())
	}
	return date
}

func MDYToFhir(mdy string) (time.Time, error){
	d, err := time.Parse("01-02-2006", mdy)
	//fmt.Printf("MDY Date: %s\n", d)
	return d, err

}


func CalcPages(docs, pgSize int64) (int, error){

	if docs < 1 || pgSize  < 1 {
		return 0, fmt.Errorf("both Docs and PageSize must be > 0")
	}
	var pages int
	x  :=  float64(docs)/float64(pgSize)
	xt := math.Trunc(x)
	it := int(x)
	fmt.Printf("CalcPages:63 -  Pages needed: %f,  config Page size: %f,  xt: %f, it: %d\n", x, float64(pgSize), xt, it)
	rndPages := math.Round(x)
	if x == float64(pgSize) {
		fmt.Printf("Exact fit\n")
		pages = int(rndPages)
	} else {
		fmt.Printf("Need extra partial page\n")
		pages = it + 1
	 }
	 return pages, nil

}



func FilterFromString(query string) (bson.M, error) {
	q, _ := MapFromString(query)
	return FilterFromMap(q)
	// m_q := []bson.M{}
	// //filter := bson.D{}
	// for k := range q {
	// 	//fmt.Printf("k: %s,  v: %s\n", k, q[k])
	// 	if k == "id" {
	// 		m_q = append(m_q, bson.M{"enterpriseid": q[k]})
	// 	} else {
	// 		val := q[k]
	// 		m_q = append(m_q, bson.M{k: val})
	// 	}
	// }
	// m_query := bson.M{}
	// if len(m_q) > 0 {
	// 	m_query = bson.M{"$and": m_q}
	// }
	// //filter := bson.D{m_query}
	// //fmt.Printf("Generated Query: %T,  %v\n", m_query, m_query)
	// //spew.Dump(m_query)
	// //err = fmt.Errorf("Invalid query: [%s]", query)
	// err = nil
	// return m_query, err
}

func MapFromString(query string) (map[string]string, error) {
	var m = make(map[string]string)
	items := strings.Split(query, "&")
	for _, v := range items {
		detail := strings.Split(v, "=")
		cnt := len(detail)
		//fmt.Printf("Detail: %v   query: %s\n", detail, query)
		if cnt == 0 || cnt > 2 {
			return nil, fmt.Errorf("Invalid query: [%s]", v)
		}
		m[strings.Trim(detail[0], " ")] = strings.Trim(detail[1], " ")
	}
	return m, nil
}

func FilterFromMap(m map[string]string) (bson.M, error) {
	m_q := []bson.M{}
	//filter := bson.D{}
	for k := range m {
		val := m[k]
		//fmt.Printf("k: %s,  v: %s\n", k, q[k])
		if k == "id" {
			fmt.Printf("Converting search for id %s to search for enterpriseid\n", val)
			m_q = append(m_q, bson.M{"enterpriseid": val})
		} else if k == "given" {
			m_q = append(m_q, bson.M{"given": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "family" {
			m_q = append(m_q, bson.M{"family": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "email" {
			m_q = append(m_q, bson.M{"email": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else {
			m_q = append(m_q, bson.M{k: val})
		}
	}
	m_query := bson.M{}
	if len(m_q) > 0 {
		m_query = bson.M{"$and": m_q}
	}
	return m_query, nil
}

func StringFromMap(m map[string]string) string {
	//fmt.Printf("StringFromMap\n")
	q := ""
	for k := range m {

		if q == "" {
			q = fmt.Sprintf("%s=%s", k, m[k])
		} else {
			q = fmt.Sprintf("%s&%s=%s", q, k, m[k])
		}
	}

	return q
}
