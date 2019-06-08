use chrono::prelude::{DateTime, Utc};
use csv::ReaderBuilder;
use serde::{Deserialize, Serialize};
use serde_json;
use std::io;

mod my_date_format {
    use chrono::{DateTime, TimeZone, Utc};
    use serde::{self, Deserialize, Deserializer, Serializer};

    const FORMAT: &'static str = "%Y-%m-%d %H:%M:%S";

    // The signature of a serialize_with function must follow the pattern:
    //
    //    fn serialize<S>(&T, S) -> Result<S::Ok, S::Error>
    //    where
    //        S: Serializer
    //
    // although it may also be generic over the input types T.
    pub fn serialize<S>(date: &DateTime<Utc>, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let s = format!("{}", date.format(FORMAT));
        serializer.serialize_str(&s)
    }

    // The signature of a deserialize_with function must follow the pattern:
    //
    //    fn deserialize<'de, D>(D) -> Result<T, D::Error>
    //    where
    //        D: Deserializer<'de>
    //
    // although it may also be generic over the output types T.
    pub fn deserialize<'de, D>(deserializer: D) -> Result<DateTime<Utc>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        Utc.datetime_from_str(&s, FORMAT)
            .map_err(serde::de::Error::custom)
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct Trip {
    // field num: 0
    id: String,        // 3389224,
    vendor_id: String, // 2,

    #[serde(with = "my_date_format")]
    pickup_datetime: DateTime<Utc>, // 2014-03-26 00:26:15,

    #[serde(with = "my_date_format")]
    dropoff_datetime: DateTime<Utc>, // 2014-03-26 00:28:38,

    store_and_fwd_flag: String,     // N,
    rate_code_id: String,           // 1,
    pickup_longitude: Option<f64>,  // -73.950431823730469,
    pickup_latitude: Option<f64>,   // 40.792251586914063,
    dropoff_longitude: Option<f64>, // -73.938949584960937,
    dropoff_latitude: Option<f64>,  // 40.794425964355469,

    // field num: 10
    passenger_count: i64,               // 1,
    trip_distance: f64,                 // 0.84,
    fare_amount: f64,                   // 4.5,
    extra: f64,                         // 0.5,
    mta_tax: f64,                       // 0.5,
    tip_amount: f64,                    // 1,
    tolls_amount: f64,                  // 0,
    ehail_fee: Option<f64>,             // ,
    improvement_surcharge: Option<f64>, // ,
    total_amount: f64,                  // 6.5,

    // field num: 20
    payment_type: i32,            // 1,
    trip_type: Option<i32>,       // 1,
    pickup_location_id: String,   // 75,
    dropoff_location_id: String,  // 74,
    cab_type: String,             // green,
    precipitation: f64,           // 0.00,
    snow_depth: f64,              // 0.0,
    snowfall: f64,                // 0.0,
    max_temp: f64,                // 36,
    min_temp: f64,                // 24,
    wind: f64,                    // 11.86,
    pickup_nyct2010_gid: String,  // 1267,
    pickup_ctlabel: String,       // 168,
    pickup_borocode: String,      // 1,
    pickup_boroname: String,      // Manhattan,
    pickup_ct2010: String,        // 016800,
    pickup_boroct2010: String,    // 1016800,
    pickup_cdeligibil: String,    // E,
    pickup_ntacode: String,       // MN33,
    pickup_ntaname: String,       // East Harlem South,
    pickup_puma: String,          // 3804,
    dropoff_nyct2010_gid: String, // 1828,
    dropoff_ctlabel: String,      // 180,
    dropoff_borocode: String,     // 1,
    dropoff_boroname: String,     // Manhattan,
    dropoff_ct2010: String,       // 018000,
    dropoff_boroct2010: String,   // 1018000,
    dropoff_cdeligibil: String,   // E,
    dropoff_ntacode: String,      // MN34,
    dropoff_ntaname: String,      // East Harlem North,
    dropoff_puma: String,         // 3804,
}

fn main() {
    let mut rdr = ReaderBuilder::new()
        .has_headers(false)
        .from_reader(io::stdin());
    rdr.deserialize::<Trip>().for_each(|result| {
        println!(
            "{}",
            serde_json::to_string(&result.expect("expected csv rec"))
                .expect("failed to encode to json")
        )
    })
}
