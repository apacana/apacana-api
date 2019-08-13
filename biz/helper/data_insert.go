package helper

import (
	"bufio"
	"fmt"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"io"
	"os"
	"strconv"
	"strings"
)

func insertHotelInfoAgoda(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	index := 0

	tmpStr := ""

	for {
		index += 1

		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if index == 1 {
			continue
		}
		if err != nil || io.EOF == err {
			break
		}
		line = tmpStr + line
		messages := strings.Split(line, ",")
		if messages[len(messages)-1] != "\"CNY\"\n" {
			tmpStr = line
			continue
		}
		tmpStr = ""

		prueLine := strings.Replace(line, "\n", " ", -1)
		//fmt.Println(prueLine)

		point := strings.Index(prueLine, ",")
		str := prueLine[0:point]
		hotelID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(hotelID, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		chainID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+2:]
		//fmt.Println(chainID, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		chainName := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(chainName, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		brandID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+2:]
		//fmt.Println(brandID, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		brandName := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(brandName, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		hotelName := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(hotelName, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		hotelFormerlyName := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(hotelFormerlyName, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		hotelTranslatedName := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(hotelTranslatedName, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		addressLine1 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(addressLine1, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		addressLine2 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(addressLine2, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		zipCode := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(zipCode, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		city := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(city, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		state := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(state, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		country := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(country, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		countryISOCode := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(countryISOCode, "-", prueLine)

		point = strings.Index(prueLine, ",")
		starRating := prueLine[0:point]
		prueLine = prueLine[point+1:]
		//fmt.Println(starRating, "-", prueLine)

		point = strings.Index(prueLine, ",")
		longitude := prueLine[0:point]
		prueLine = prueLine[point+1:]
		//fmt.Println(longitude, "-", prueLine)

		point = strings.Index(prueLine, ",")
		latitude := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(latitude, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		url := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(url, "-", prueLine)

		point = strings.Index(prueLine, ",")
		checkin := prueLine[0:point]
		prueLine = prueLine[point+1:]
		//fmt.Println(checkin, "-", prueLine)

		point = strings.Index(prueLine, ",")
		checkout := prueLine[0:point]
		prueLine = prueLine[point+1:]
		//fmt.Println(checkout, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		numberRooms, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(numberRooms, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		numberFloors, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(numberFloors, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		yearOpened, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(yearOpened, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		yearRenovated, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+2:]
		//fmt.Println(yearRenovated, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		photo1 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(photo1, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		photo2 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(photo2, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		photo3 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(photo3, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		photo4 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(photo4, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		photo5 := prueLine[0:point]
		prueLine = prueLine[point+3:]
		//fmt.Println(photo5, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		overView := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(overView, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		rateFrom, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(rateFrom, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		continentID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+2:]
		//fmt.Println(continentID, "-", prueLine)

		point = strings.Index(prueLine, "\",")
		continentName := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(continentName, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		cityID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(cityID, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		countryID, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(countryID, "-", prueLine)

		point = strings.Index(prueLine, ",")
		str = prueLine[0:point]
		numberOfReviews, _ := strconv.ParseInt(str, 10, 64)
		prueLine = prueLine[point+1:]
		//fmt.Println(numberOfReviews, "-", prueLine)

		point = strings.Index(prueLine, ",")
		ratingAverage := prueLine[0:point]
		prueLine = prueLine[point+2:]
		//fmt.Println(ratingAverage, "-", prueLine)

		point = strings.Index(prueLine, "\"")
		ratesCurrency := prueLine[0:point]
		//fmt.Println(ratesCurrency, "-", prueLine)

		hotelInfoAgoda := &mysql.HotelInfoAgoda{
			HotelID:             hotelID,
			ChainID:             chainID,
			ChainName:           chainName,
			BrandID:             brandID,
			BrandName:           brandName,
			HotelName:           hotelName,
			HotelFormerlyName:   hotelFormerlyName,
			HotelTranslatedName: hotelTranslatedName,
			AddressLine1:        addressLine1,
			AddressLine2:        addressLine2,
			ZipCode:             zipCode,
			City:                city,
			State:               state,
			Country:             country,
			CountryISOCode:      countryISOCode,
			StarRating:          starRating,
			Longitude:           longitude,
			Latitude:            latitude,
			Url:                 url,
			CheckIn:             checkin,
			CheckOut:            checkout,
			NumberRooms:         numberRooms,
			NumberFloors:        numberFloors,
			YearOpened:          yearOpened,
			YearRenovated:       yearRenovated,
			Photo1:              photo1,
			Photo2:              photo2,
			Photo3:              photo3,
			Photo4:              photo4,
			Photo5:              photo5,
			OverView:            overView,
			RatesFrom:           rateFrom,
			ContinentID:         continentID,
			ContinentName:       continentName,
			CityID:              cityID,
			CountryID:           countryID,
			NumberOfReviews:     numberOfReviews,
			RatingAverage:       ratingAverage,
			RatesCurrency:       ratesCurrency,
		}

		fmt.Println("index:", index)
		//fmt.Printf("\r\n %+v", *hotelInfoAgoda)
		err = mysql.Insert(hotelInfoAgoda)
		if err != nil {
			fmt.Println(err)
		}

		//break
	}
}
