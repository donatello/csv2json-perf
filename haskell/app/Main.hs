{-# LANGUAGE DeriveGeneric, OverloadedStrings, BangPatterns #-}
module Main where

import Data.Text (Text)
import Data.Time (UTCTime)
import Data.Int
import qualified Data.CSV.Conduit.Conversion as Csv
import qualified Data.CSV.Conduit as Csv
import qualified Data.Aeson as Json
import Data.Aeson.Casing (aesonDrop, snakeCase)
import GHC.Generics (Generic)
import Data.Time (parseTimeM, defaultTimeLocale)
import qualified Data.ByteString.Char8 as BC
import Data.ByteString (ByteString)
import Data.ByteString.Lazy (toStrict)
import qualified Conduit as C
import Conduit ((.|))

data Trip =
    Trip
    { id                   :: !Text
    , vendorId             :: !Text
    , pickupDatetime       :: !UTCTime
    , dropoffDatetime      :: !UTCTime
    , storeAndFwdFlag      :: !Text
    , rateCodeId           :: !Text
    , pickupLongitude      :: Maybe Double
    , pickupLatitude       :: Maybe Double
    , dropoffLongitude     :: Maybe Double
    , dropoffLatitude      :: Maybe Double
    , passengerCount       :: !Int64
    , tripDistance         :: !Double
    , fareAmount           :: !Double
    , extra                :: !Double
    , mtaTax               :: !Double
    , tipAmount            :: !Double
    , tollsAmount          :: !Double
    , ehailFee             :: Maybe Double
    , improvementSurcharge :: Maybe Double
    , totalAmount          :: !Double
    , paymentType          :: !Int32
    , tripType             :: Maybe Int32
    , pickupLocationId     :: !Text
    , dropoffLocationId    :: !Text
    , cabType              :: !Text
    , precipitation        :: !Double
    , snowDepth            :: !Double
    , snowfall             :: !Double
    , maxTemp              :: !Double
    , minTemp              :: !Double
    , wind                 :: !Double
    , pickupNyct2010Gid    :: !Text
    , pickupCtlabel        :: !Text
    , pickupBorocode       :: !Text
    , pickupBoroname       :: !Text
    , pickupCt2010         :: !Text
    , pickupBoroct2010     :: !Text
    , pickupCdeligibil     :: !Text
    , pickupNtacode        :: !Text
    , pickupNtaname        :: !Text
    , pickupPuma           :: !Text
    , dropoffNyct2010Gid   :: !Text
    , dropoffCtlabel       :: !Text
    , dropoffBorocode      :: !Text
    , dropoffBoroname      :: !Text
    , dropoffCt2010        :: !Text
    , dropoffBoroct2010    :: !Text
    , dropoffCdeligibil    :: !Text
    , dropoffNtacode       :: !Text
    , dropoffNtaname       :: !Text
    , dropoffPuma          :: !Text
    }
    deriving (Show, Generic)

instance Csv.FromField UTCTime where
    parseField b = parseTimeM False defaultTimeLocale
                   "%Y-%m-%d %H:%M:%S" $ BC.unpack b

instance Csv.FromRecord Trip where

jsonOpts :: Json.Options
jsonOpts = aesonDrop 0 snakeCase

instance Json.ToJSON Trip where
    toJSON = Json.genericToJSON jsonOpts
    toEncoding = Json.genericToEncoding jsonOpts

parseCsv :: Csv.Record -> Trip
parseCsv = either undefined Prelude.id
         . Csv.runParser
         . Csv.parseRecord

jsonEncode :: Trip -> ByteString
jsonEncode t = (toStrict $ Json.encode t) <> "\n"

main = do
    C.runConduitRes
        $ C.stdinC
        .| Csv.intoCSV Csv.defCSVSettings
        .| C.mapC parseCsv
        .| C.mapC jsonEncode
        .| C.stdoutC
