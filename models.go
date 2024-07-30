package main

type (
	PlaceDetails struct {
		InputId          string
		Link             string
		Title            string
		Category         string
		Address          string
		OpenHours        string
		PopularTimes     string
		Website          string
		Phone            string
		PlusCode         string
		ReviewCount      string
		ReviewRating     string
		ReviewsPerRating string
		Latitude         string
		Longitude        string
		Cid              string
		Status           string
		Descriptions     string
		ReviewsLink      string
		Thumbnail        string
		Timezone         string
		PriceRange       string
		DataId           string
		Images           string
		Reservations     string
		OrderOnline      string
		Menu             string
		Owner            string
		CompleteAddress  string
		About            string
		UserReviews      string
		Emails           string
	}

	QrDataDetails struct {
		MerchantName            string
		TotalTransaction        string
		TotalNominalTransaction string
		Last1MonthTransaction   string
		Last1MothNominal        string
	}

	EndResult struct {
		PlaceDetails
		QrDataDetails
	}
)
