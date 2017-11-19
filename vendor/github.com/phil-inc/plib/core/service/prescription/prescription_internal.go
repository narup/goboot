package prescription

import "github.com/narup/gmgo"

// rx query based on config. Config has 3 variables leading up to 8 possible query combinations
func rxQuery(config *RxLoadConfig, orderIds []string) gmgo.Q {
	q := gmgo.Q{}
	if !config.FilterArchived && !config.FilterPaused && !config.FilterPending {
		//000
		q = gmgo.Q{
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if !config.FilterArchived && !config.FilterPaused && config.FilterPending {
		//001
		q = gmgo.Q{
			"pending": false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if !config.FilterArchived && config.FilterPaused && !config.FilterPending {
		//010
		q = gmgo.Q{
			"suspended": false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if !config.FilterArchived && config.FilterPaused && config.FilterPending {
		//011
		q = gmgo.Q{
			"suspended": false,
			"pending":   false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if config.FilterArchived && !config.FilterPaused && !config.FilterPending {
		//100
		q = gmgo.Q{
			"archived": false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if config.FilterArchived && !config.FilterPaused && config.FilterPending {
		//101
		q = gmgo.Q{
			"archived": false,
			"pending":  false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if config.FilterArchived && config.FilterPaused && !config.FilterPending {
		//110
		q = gmgo.Q{
			"archived":  false,
			"suspended": false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else if config.FilterArchived && config.FilterPaused && config.FilterPending {
		//111
		q = gmgo.Q{
			"archived":  false,
			"pending":   false,
			"suspended": false,
			"currentOrderId": gmgo.Q{
				"$in": orderIds,
			},
		}
	} else {
		//
	}
	return q
}
