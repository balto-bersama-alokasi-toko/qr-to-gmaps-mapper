package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/schollz/progressbar/v3"
)

func main() {
	gmapsDataMap := make(map[string][]PlaceDetails)
	var endResultData []EndResult

	gmapsSource, err := os.Open("gmaps.csv")
	if err != nil {
		panic(err)
	}
	defer func(gmapsSource *os.File) {
		err := gmapsSource.Close()
		if err != nil {
			panic(err)
		}
	}(gmapsSource)

	csvReader := csv.NewReader(gmapsSource)
	records, _ := csvReader.ReadAll()

	for _, record := range records {
		gmapsDataMap[record[0]] = append(gmapsDataMap[record[0]], PlaceDetails{
			InputId:          record[0],
			Link:             record[1],
			Title:            record[2],
			Category:         record[3],
			Address:          record[4],
			OpenHours:        record[5],
			PopularTimes:     record[6],
			Website:          record[7],
			Phone:            record[8],
			PlusCode:         record[9],
			ReviewCount:      record[10],
			ReviewRating:     record[11],
			ReviewsPerRating: record[12],
			Latitude:         record[13],
			Longitude:        record[14],
			Cid:              record[15],
			Status:           record[16],
			Descriptions:     record[17],
			ReviewsLink:      record[18],
			Thumbnail:        record[19],
			Timezone:         record[20],
			PriceRange:       record[21],
			DataId:           record[22],
			Images:           record[23],
			Reservations:     record[24],
			OrderOnline:      record[25],
			Menu:             record[26],
			Owner:            record[27],
			CompleteAddress:  record[28],
			About:            record[29],
			UserReviews:      record[30],
			Emails:           record[31],
		})
	}

	// get source QR Data
	qrFile, err := os.Open("cleaned.csv")
	if err != nil {
		panic(err)
	}
	defer func(qrFile *os.File) {
		err := qrFile.Close()
		if err != nil {
			panic(err)
		}
	}(qrFile)

	csvReader = csv.NewReader(qrFile)
	records, _ = csvReader.ReadAll()
	var qrDataDetails []QrDataDetails
	for _, record := range records {
		qrDataDetails = append(qrDataDetails, QrDataDetails{
			MerchantName:            record[0],
			TotalTransaction:        record[1],
			TotalNominalTransaction: record[2],
			Last1MonthTransaction:   record[3],
			Last1MothNominal:        record[4],
		})
	}

	fmt.Println("PROCESSING DATA...")
	fmt.Println("")

	bar := progressbar.Default(int64(len(qrDataDetails)))
	// fuzzy search
	var sliceStr []string
	var stringLen int
	for index, place := range qrDataDetails {
		bar.Add(1)
		if index == 0 {
			continue
		}
		sliceStr = strings.Split(strings.ToLower(strings.TrimSpace(place.MerchantName[:25])), " ")
		stringLen = len(sliceStr)
		var match fuzzy.Rank
		var placeDetailMatch []PlaceDetails
		var ranks fuzzy.Ranks
		var gmapsScrapKey string
		for iterator := 0; iterator < stringLen-1; iterator++ {
			if len(sliceStr) < 2 || len(ranks) != 0 || sliceStr[len(sliceStr)-1] == "bang" {
				break
			}
			for keyId, placeDetails := range gmapsDataMap {
				var slicesOfPlaceName []string
				for _, place := range placeDetails {
					slicesOfPlaceName = append(slicesOfPlaceName, strings.ToLower(place.Title))
				}
				joinedStr := strings.Join(sliceStr, " ")
				ranks = fuzzy.RankFind(joinedStr, slicesOfPlaceName)
				for i, rank := range ranks {
					placeDetailMatch = placeDetails
					gmapsScrapKey = keyId
					if i == 0 {
						match = rank
					} else {
						if rank.Distance < match.Distance {
							match = rank
						}
					}
				}
			}
			if match.Target != "" && match.Distance == 0 {
				delete(gmapsDataMap, gmapsScrapKey)
				endResultData = append(endResultData, EndResult{
					PlaceDetails:  placeDetailMatch[match.OriginalIndex],
					QrDataDetails: place,
				})
				break
			}
			sliceStr = sliceStr[:stringLen-1]
		}
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("WRITING DATA...")

	writeOutputBar := progressbar.Default(int64(len(endResultData)))
	endresultFile, err := os.OpenFile("endresult.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer endresultFile.Close()

	writer := csv.NewWriter(endresultFile)

	writer.Write([]string{
		"input_id",
		"link",
		"title",
		"category",
		"address",
		"open_hours",
		"popular_times",
		"website",
		"phone",
		"plus_code",
		"review_count",
		"review_rating",
		"reviews_per_rating",
		"latitude",
		"longitude",
		"cid",
		"status",
		"descriptions",
		"reviews_link",
		"thumbnail",
		"timezone",
		"price_range",
		"data_id",
		"images",
		"reservations",
		"order_online",
		"menu",
		"owner",
		"complete_address",
		"about",
		"user_reviews",
		"emails",
		"MerchantName",
		"TotalTransaction",
		"TotalNominalTransaction",
		"Last1MonthTransaction",
		"Last1MothNominal",
	})
	for _, endResult := range endResultData {
		writeOutputBar.Add(1)
		sliceOfStr := convertEndResultToStringSlice(endResult)
		err = writer.Write(sliceOfStr)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

// gmaps-to-qr.go
func convertEndResultToStringSlice(er EndResult) []string {
	var result []string

	// Process PlaceDetails
	v := reflect.ValueOf(er.PlaceDetails)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				result = append(result, fmt.Sprintf("%v", field.Elem().Interface()))
			} else {
				result = append(result, "") // Add empty string for nil fields
			}
		} else {
			result = append(result, fmt.Sprintf("%v", field.Interface()))
		}
	}

	// Process QrDataDetails
	v = reflect.ValueOf(er.QrDataDetails)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				result = append(result, fmt.Sprintf("%v", field.Elem().Interface()))
			} else {
				result = append(result, "") // Add empty string for nil fields
			}
		} else {
			result = append(result, fmt.Sprintf("%v", field.Interface()))
		}
	}

	return result
}
