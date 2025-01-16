package logger

// GetLoggerPaths returns all configured logging paths for the application
func GetLoggerPaths() []*LogPath {
	return []*LogPath{
		{
			// Restaurant Summary Component Logs
			Path:          formatPath("quananqr1/app/manage/admin/orders/restaurant-summary/restaurant-summary"),
			Enabled:       false,
			Description:   "Restaurant Summary Component Logs",
			EnabledLogIDs: []int{1, 2, 3},
			DisabledLogIDs: []int{},
			LogIDs:        []int{1, 2, 3},
			LogDescriptions: map[int]LogDescription{
				1: {
					Description: "Log initial dish aggregation state",
					Location:    "aggregateDishes function - initialization",
					Status:      "enabled",
				},
				2: {
					Description: "Log aggregated dishes for order groups",
					Location:    "RestaurantSummary component - groupedOrders processing",
					Status:      "enabled",
				},
				3: {
					Description: "Log aggregation completion state",
					Location:    "RestaurantSummary component - final state",
					Status:      "enabled",
				},
			},
		},
		// Add other paths here...
	}
}