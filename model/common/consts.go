package common

const (
	ListingFields = "SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price, " +
		"listing.discount as discount, listing.discount_description as discount_description, listing.description as description, listing.start_date as start_date, " +
		"listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time, " +
		"listing.multiple_days as multiple_days, " +
		"listing.recurring as recurring, listing.recurring_end_date as recurring_date, listing.listing_type as listing_type, " +
		"listing.business_id as business_id, listing.listing_id as listing_id "

	ListingDateFields = "listing_date.listing_date_id as listing_date_id, listing_date.listing_date as listing_date"

	ListingBusinessFields = "business.name as bname"

	ListingBusinessAddressFields = "business_address.latitude as latitude, business_address.longitude as longitude "

	ListingImageFields = "listing_image.path as path"

	FavoriteFields = "favorites.favorite_id as favorite_id, favorites.favorite_add_date as favorite_add_date"

	FromClauseListing = "FROM listing " +
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id " +
		"INNER JOIN business ON listing.business_id = business.business_id " +
		"INNER JOIN business_cuisine ON listing.business_id = business_cuisine.business_id " +
		"INNER JOIN listing_image ON listing.listing_id = listing_image.listing_id"

	FromClauseListingWithAddress = "FROM listing " +
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id " +
		"INNER JOIN business ON listing.business_id = business.business_id " +
		"INNER JOIN business_cuisine ON listing.business_id = business_cuisine.business_id " +
		"INNER JOIN business_address ON listing.business_id = business_address.business_id " +
		"INNER JOIN listing_image ON listing.listing_id = listing_image.listing_id"

	FromClauseFavorites = "FROM favorites " +
		"INNER JOIN listing ON listing.listing_id = favorites.listing_id " +
		"INNER JOIN listing_image ON listing_image.listing_id = favorites.listing_id " +
		"INNER JOIN business ON listing.business_id = business.business_id " +
		"INNER JOIN business_address ON listing.business_id = business_address.business_id "

	FromClauseListingAdmin = "FROM listing " +
		"INNER JOIN business ON listing.business_id = business.business_id " +
		"INNER JOIN business_cuisine ON listing.business_id = business_cuisine.business_id " +
		"INNER JOIN listing_image ON listing.listing_id = listing_image.listing_id"

	MaxFilterDistance = "15"

	MaxFutureDays = 3

	MaxDistanceForTodaysDeals = 15.0

	MaxDistanceForFutureDeals = 25.0

	MaxRangeAroundSunnyvale = 100.0

	MaxDistanceToGroupNow = 7.5
)
