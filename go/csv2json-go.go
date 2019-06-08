package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type myTime time.Time

var layout = "2006-01-02 15:04:05"

func (m myTime) MarshalJSON() ([]byte, error) {
	s := time.Time(m).Format(layout)
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

// Trip holds each line of the CSV
type Trip struct {
	Id                    string   `csv:"id" json:"id"`
	Vendor_Id             string   `csv:"vendor_id" json:"vendor_id"`
	Pickup_Datetime       myTime   `csv:"pickup_datetime" json:"pickup_datetime"`
	Dropoff_Datetime      myTime   `csv:"dropoff_datetime" json:"dropoff_datetime"`
	Store_And_Fwd_Flag    string   `csv:"store_and_fwd_flag" json:"store_and_fwd_flag"`
	Rate_Code_Id          string   `csv:"rate_code_id" json:"rate_code_id"`
	Pickup_Longitude      *float64 `csv:"pickup_longitude" json:"pickup_longitude"`
	Pickup_Latitude       *float64 `csv:"pickup_latitude" json:"pickup_latitude"`
	Dropoff_Longitude     *float64 `csv:"dropoff_longitude" json:"dropoff_longitude"`
	Dropoff_Latitude      *float64 `csv:"dropoff_latitude" json:"dropoff_latitude"`
	Passenger_Count       int64    `csv:"passenger_count" json:"passenger_count"`
	Trip_Distance         float64  `csv:"trip_distance" json:"trip_distance"`
	Fare_Amount           float64  `csv:"fare_amount" json:"fare_amount"`
	Extra                 float64  `csv:"extra" json:"extra"`
	Mta_Tax               float64  `csv:"mta_tax" json:"mta_tax"`
	Tip_Amount            float64  `csv:"tip_amount" json:"tip_amount"`
	Tolls_Amount          float64  `csv:"tolls_amount" json:"tolls_amount"`
	Ehail_Fee             *float64 `csv:"ehail_fee" json:"ehail_fee"`
	Improvement_Surcharge *float64 `csv:"improvement_surcharge" json:"improvement_surcharge"`
	Total_Amount          float64  `csv:"total_amount" json:"total_amount"`
	Payment_Type          int64    `csv:"payment_type" json:"payment_type"`
	Trip_Type             *int64   `csv:"trip_type" json:"trip_type"`
	Pickup_Location_Id    string   `csv:"pickup_location_id" json:"pickup_location_id"`
	Dropoff_Location_Id   string   `csv:"dropoff_location_id" json:"dropoff_location_id"`
	Cab_Type              string   `csv:"cab_type" json:"cab_type"`
	Precipitation         float64  `csv:"precipitation" json:"precipitation"`
	Snow_Depth            float64  `csv:"snow_depth" json:"snow_depth"`
	Snowfall              float64  `csv:"snowfall" json:"snowfall"`
	Max_Temp              float64  `csv:"max_temp" json:"max_temp"`
	Min_Temp              float64  `csv:"min_temp" json:"min_temp"`
	Wind                  float64  `csv:"wind" json:"wind"`
	Pickup_Nyct2010_Gid   string   `csv:"pickup_nyct2010_gid" json:"pickup_nyct2010_gid"`
	Pickup_Ctlabel        string   `csv:"pickup_ctlabel" json:"pickup_ctlabel"`
	Pickup_Borocode       string   `csv:"pickup_borocode" json:"pickup_borocode"`
	Pickup_Boroname       string   `csv:"pickup_boroname" json:"pickup_boroname"`
	Pickup_Ct2010         string   `csv:"pickup_ct2010" json:"pickup_ct2010"`
	Pickup_Boroct2010     string   `csv:"pickup_boroct2010" json:"pickup_boroct2010"`
	Pickup_Cdeligibil     string   `csv:"pickup_cdeligibil" json:"pickup_cdeligibil"`
	Pickup_Ntacode        string   `csv:"pickup_ntacode" json:"pickup_ntacode"`
	Pickup_Ntaname        string   `csv:"pickup_ntaname" json:"pickup_ntaname"`
	Pickup_Puma           string   `csv:"pickup_puma" json:"pickup_puma"`
	Dropoff_Nyct2010_Gid  string   `csv:"dropoff_nyct2010_gid" json:"dropoff_nyct2010_gid"`
	Dropoff_Ctlabel       string   `csv:"dropoff_ctlabel" json:"dropoff_ctlabel"`
	Dropoff_Borocode      string   `csv:"dropoff_borocode" json:"dropoff_borocode"`
	Dropoff_Boroname      string   `csv:"dropoff_boroname" json:"dropoff_boroname"`
	Dropoff_Ct2010        string   `csv:"dropoff_ct2010" json:"dropoff_ct2010"`
	Dropoff_Boroct2010    string   `csv:"dropoff_boroct2010" json:"dropoff_boroct2010"`
	Dropoff_Cdeligibil    string   `csv:"dropoff_cdeligibil" json:"dropoff_cdeligibil"`
	Dropoff_Ntacode       string   `csv:"dropoff_ntacode" json:"dropoff_ntacode"`
	Dropoff_Ntaname       string   `csv:"dropoff_ntaname" json:"dropoff_ntaname"`
	Dropoff_Puma          string   `csv:"dropoff_puma" json:"dropoff_puma"`
}

var (
	fieldType []string
	fieldName []string
)

func parseTime(i int, r []string) myTime {
	t, e := time.Parse(layout, r[i])
	if e != nil {
		panic(e)
	}
	return myTime(t)
}

func parseFloatPtr(i int, r []string) *float64 {
	if r[i] == "" {
		return nil
	}
	v := parseFloat(i, r)
	return &v
}

func parseFloat(i int, r []string) float64 {
	v, e := strconv.ParseFloat(r[i], 64)
	if e != nil {
		panic(e)
	}
	return v
}

func parseInt(i int, r []string) int64 {
	v, e := strconv.ParseInt(r[i], 10, 64)
	if e != nil {
		panic(e)
	}
	return v
}

func parseIntPtr(i int, r []string) *int64 {
	if r[i] == "" {
		return nil
	}
	v := parseInt(i, r)
	return &v
}

func parseRecord(record []string) *Trip {
	return &Trip{
		Id:                    record[0],
		Vendor_Id:             record[1],
		Pickup_Datetime:       parseTime(2, record),
		Dropoff_Datetime:      parseTime(3, record),
		Store_And_Fwd_Flag:    record[4],
		Rate_Code_Id:          record[5],
		Pickup_Longitude:      parseFloatPtr(6, record),
		Pickup_Latitude:       parseFloatPtr(7, record),
		Dropoff_Longitude:     parseFloatPtr(8, record),
		Dropoff_Latitude:      parseFloatPtr(9, record),
		Passenger_Count:       parseInt(10, record),
		Trip_Distance:         parseFloat(11, record),
		Fare_Amount:           parseFloat(12, record),
		Extra:                 parseFloat(13, record),
		Mta_Tax:               parseFloat(14, record),
		Tip_Amount:            parseFloat(15, record),
		Tolls_Amount:          parseFloat(16, record),
		Ehail_Fee:             parseFloatPtr(17, record),
		Improvement_Surcharge: parseFloatPtr(18, record),
		Total_Amount:          parseFloat(19, record),
		Payment_Type:          parseInt(20, record),
		Trip_Type:             parseIntPtr(21, record),
		Pickup_Location_Id:    record[22],
		Dropoff_Location_Id:   record[23],
		Cab_Type:              record[24],
		Precipitation:         parseFloat(25, record),
		Snow_Depth:            parseFloat(26, record),
		Snowfall:              parseFloat(27, record),
		Max_Temp:              parseFloat(28, record),
		Min_Temp:              parseFloat(29, record),
		Wind:                  parseFloat(30, record),
		Pickup_Nyct2010_Gid:   record[31],
		Pickup_Ctlabel:        record[32],
		Pickup_Borocode:       record[33],
		Pickup_Boroname:       record[34],
		Pickup_Ct2010:         record[35],
		Pickup_Boroct2010:     record[36],
		Pickup_Cdeligibil:     record[37],
		Pickup_Ntacode:        record[38],
		Pickup_Ntaname:        record[39],
		Pickup_Puma:           record[40],
		Dropoff_Nyct2010_Gid:  record[41],
		Dropoff_Ctlabel:       record[42],
		Dropoff_Borocode:      record[43],
		Dropoff_Boroname:      record[44],
		Dropoff_Ct2010:        record[45],
		Dropoff_Boroct2010:    record[46],
		Dropoff_Cdeligibil:    record[47],
		Dropoff_Ntacode:       record[48],
		Dropoff_Ntaname:       record[49],
		Dropoff_Puma:          record[50],
	}
}

// takes a stream reader
func process_reader(data io.Reader) {
	// We are going to reflect on the Trip struct once to discover all it's fields
	t := reflect.TypeOf(Trip{})
	// track the type of the field as a string by position
	fieldType = make([]string, t.NumField())
	// track the field position to it's name
	fieldName = make([]string, t.NumField())
	// iterate struct fields
	for i := 0; i < t.NumField(); i++ {
		fieldType[i] = t.Field(i).Type.String()
		// parse the tag and extract the json notation we have right there.
		// this is pretty dirty and prone to error if the tag changes
		parts := strings.Split(string(t.Field(i).Tag), " ")
		secondParts := strings.Split(parts[1], "\"")
		fieldName[i] = secondParts[1]
	}
	r := csv.NewReader(data)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		trip := parseRecord(record)
		b, e := json.Marshal(trip)
		check(e)

		fmt.Println(string(b))
	}

}

func main() {
	// fmt.Println("Hellooo")
	// Parse the flags to see wether we read from the stdin or a file
	flag.Parse()

	switch flag.NArg() {
	case 0: // Read from stdin
		// procress stdin
		process_reader(os.Stdin)
		break

	case 1: //read from file
		//open file and process
		data, err := os.Open(flag.Arg(0))
		check(err)
		process_reader(data)
		break
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}
}
