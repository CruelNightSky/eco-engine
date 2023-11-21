package table

import cl "eco-engine/customlog"

type CostTable struct {
	UpgradesCost struct {
		Damage struct {
			Value        []int  `json:"value"`
			ResourceType string `json:"resourceType"`
		} `json:"damage"`
		Attack struct {
			Value        []int  `json:"value"`
			ResourceType string `json:"resourceType"`
		} `json:"attack"`
		Health struct {
			Value        []int  `json:"value"`
			ResourceType string `json:"resourceType"`
		} `json:"health"`
		Defence struct {
			Value        []int  `json:"value"`
			ResourceType string `json:"resourceType"`
		} `json:"defence"`
	} `json:"upgradesCost"`
	UpgradeMultiplier struct {
		Damage  []float64 `json:"damage"`
		Attack  []float64 `json:"attack"`
		Health  []float64 `json:"health"`
		Defence []float64 `json:"defence"`
	} `json:"upgradeMultiplier"`
	UpgradeBaseStats struct {
		Damage struct {
			Min []int `json:"min"`
			Max []int `json:"max"`
		} `json:"damage"`
		Attack  []float64 `json:"attack"`
		Health  []int     `json:"health"`
		Defence []float64 `json:"defence"`
	} `json:"upgradeBaseStats"`
	Bonuses struct {
		StrongerMinions struct {
			MaxLevel     int       `json:"maxLevel"`
			Cost         []int     `json:"cost"`
			ResourceType string    `json:"resourceType"`
			Value        []float64 `json:"value"`
		} `json:"strongerMinions"`
		TowerMultiAttack struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"towerMultiAttack"`
		TowerAura struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"towerAura"`
		TowerVolley struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"towerVolley"`
		XpSeeking struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"xpSeeking"`
		TomeSeeking struct {
			MaxLevel     int       `json:"maxLevel"`
			Cost         []int     `json:"cost"`
			ResourceType string    `json:"resourceType"`
			Value        []float64 `json:"value"`
		} `json:"tomeSeeking"`
		EmeraldsSeeking struct {
			MaxLevel     int       `json:"maxLevel"`
			Cost         []int     `json:"cost"`
			ResourceType string    `json:"resourceType"`
			Value        []float64 `json:"value"`
		} `json:"emeraldsSeeking"`
		LargerResourceStorage struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"largerResourceStorage"`
		LargerEmeraldsStorage struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"largerEmeraldsStorage"`
		EfficientResource struct {
			MaxLevel     int       `json:"maxLevel"`
			Cost         []int     `json:"cost"`
			ResourceType string    `json:"resourceType"`
			Value        []float64 `json:"value"`
		} `json:"efficientResource"`
		EfficientEmeralds struct {
			MaxLevel     int       `json:"maxLevel"`
			Cost         []int     `json:"cost"`
			ResourceType string    `json:"resourceType"`
			Value        []float64 `json:"value"`
		} `json:"efficientEmeralds"`
		ResourceRate struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"resourceRate"`
		EmeraldsRate struct {
			MaxLevel     int    `json:"maxLevel"`
			Cost         []int  `json:"cost"`
			ResourceType string `json:"resourceType"`
			Value        []int  `json:"value"`
		} `json:"emeraldsRate"`
	} `json:"bonuses"`
}

type TerritoryProperty struct {
	TargetUpgrades  TerritoryPropertyUpgradeData `json:"targetUpgrades"`
	TargetBonuses   TerritoryPropertyBonusesData `json:"targetBonuses"`
	CurrentUpgrades TerritoryPropertyUpgradeData `json:"currentUpgrades"`
	CurrentBonuses  TerritoryPropertyBonusesData `json:"currentBonuses"`
	Tax             Tax                          `json:"tax"`
	Border          string                       `json:"border"`
	TradingStyle    string                       `json:"tradingStyle"`
	HQ              bool                         `json:"hq"`
}

type TerritoryPropertyUpgradeData struct {
	Damage  int `json:"damage"`
	Attack  int `json:"attack"`
	Health  int `json:"health"`
	Defence int `json:"defence"`
}

type TerritoryPropertyBonusesData struct {
	StrongerMinions       int `json:"strongerMinions"`
	TowerMultiAttack      int `json:"towerMultiAttacks"`
	TowerAura             int `json:"towerAura"`
	TowerVolley           int `json:"towerVolley"`
	LargerResourceStorage int `json:"largerResourceStorage"`
	LargerEmeraldStorage  int `json:"largerEmeraldsStorage"`
	EfficientResource     int `json:"efficientResource"`
	EfficientEmerald      int `json:"efficientEmeralds"`
	ResourceRate          int `json:"resourceRate"`
	EmeraldRate           int `json:"emeraldsRate"`
}

type Tax struct {
	Ally   int `json:"ally"`
	Others int `json:"others"`
}

type TerritoryResource struct {
	Emerald float64 `json:"emeralds"`
	Ore     float64 `json:"ore"`
	Wood    float64 `json:"wood"`
	Fish    float64 `json:"fish"`
	Crop    float64 `json:"crop"`
}

type Territory struct {
	ID                         int                      `json:"id"`
	Name                       string                   `json:"name"`
	Type                       string                   `json:"type"`
	Level                      string                   `json:"level"`
	RawLevel                   int                      `json:"rawLevel"`
	Claim                      bool                     `json:"claim"`
	Ally                       bool                     `json:"ally"`
	BaseResourceProduction     TerritoryResource        `json:"baseResourceProduction"`
	ResourceProduction         TerritoryResource        `json:"resourceProduction"`
	TerritoryUsage             TerritoryResource        `json:"territoryUsage"`
	Property                   TerritoryProperty        `json:"property"`
	Stats                      TowerStats               `json:"stats"`
	Storage                    TerritoryResourceStorage `json:"storage"`
	TraversingResourceToHQ   []TraversingResource    `json:"trasversingResourceToHQ"`
	TraversingResourceFromHQ []TraversingResource    `json:"trasversingResourceFromHQ"`
	TradingRoutes              []string                 `json:"tradingRoutes"`
	RouteToHQ                  []string                 `json:"routeToHQ"`
	RouteFromHQ                []string                 `json:"routeFromHQ"`
	RouteTax                   float64                  `json:"routeTax"`
	Warnings									 []string                 `json:"warnings"`
}

type TowerStats struct {
	Health            uint64  `json:"health"`
	Attack            float32 `json:"attack"`
	Defence           float32 `json:"defence"`
	StrongerMinions   int     `json:"strongerMinions"`
	TowerMultiAttacks int     `json:"towerMultiAttacks"`
	TowerAura         int     `json:"towerAura"`
	TowerVolley       int     `json:"towerVolley"`
	Damage            struct {
		Min uint64 `json:"min"`
		Max uint64 `json:"max"`
	} `json:"damage"`
}

type TraversingResource struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	RouteToDest []string `json:"routeToDest"`
	Emerald     float64  `json:"emeralds"`
	Ore         float64  `json:"ore"`
	Crop        float64  `json:"crop"`
	Wood        float64  `json:"wood"`
	Fish        float64  `json:"fish"`
}

type TerritoryResourceStorage struct {
	Capacity TerritoryResourceStorageValue `json:"capacity"`
	Current  TerritoryResourceStorageValue `json:"current"`
}

type TerritoryResourceStorageValue struct {
	Emerald float64 `json:"emeralds"`
	Wood    float64 `json:"wood"`
	Ore     float64 `json:"ore"`
	Fish    float64 `json:"fish"`
	Crop    float64 `json:"crop"`
}

type RawTerritoryData struct {
	Type          string
	TradingRoutes []string `json:"Trading Routes"`
	Resources     struct {
		Emeralds float64 `json:"emeralds"`
		Ore      float64 `json:"ore"`
		Crops    float64 `json:"crops"`
		Fish     float64 `json:"fish"`
		Wood     float64 `json:"wood"`
	} `json:"resources"`
}

type TerritoryUpdateData struct {
	Property TerritoryProperty `json:"territoryUpdateDataProperty"`
}

func (t *Territory) Set(d TerritoryUpdateData) *Territory {
	cl.Log("Updating territor data for", t.Name, "[ ID:", t.ID, "] with data", d)
	// validate the territory before setting
	// Damage, Attack, Health and Defence must be between 0 and 11
	// StrongerMinions 0 - 3, Tower MultiAttack 0 - 1, Tower Aura and Volley 0 - 3
	// Larger Emerald and Resource storage 0 - 6, Efficient Resource 0 - 6, Efficient Emerald 0 - 3 and Resource and Emerald Rate 0 - 3
	if d.Property.TargetUpgrades.Damage < 0 || d.Property.TargetUpgrades.Damage > 11 {
		return t
	} else if d.Property.TargetUpgrades.Attack < 0 || d.Property.TargetUpgrades.Attack > 11 {
		return t
	} else if d.Property.TargetUpgrades.Health < 0 || d.Property.TargetUpgrades.Health > 11 {
		return t
	} else if d.Property.TargetUpgrades.Defence < 0 || d.Property.TargetUpgrades.Defence > 11 {
		return t
	} else if d.Property.TargetBonuses.StrongerMinions < 0 || d.Property.TargetBonuses.StrongerMinions > 3 {
		return t
	} else if d.Property.TargetBonuses.TowerMultiAttack < 0 || d.Property.TargetBonuses.TowerMultiAttack > 1 {
		return t
	} else if d.Property.TargetBonuses.TowerAura < 0 || d.Property.TargetBonuses.TowerAura > 3 {
		return t
	} else if d.Property.TargetBonuses.TowerVolley < 0 || d.Property.TargetBonuses.TowerVolley > 3 {
		return t
	} else if d.Property.TargetBonuses.LargerResourceStorage < 0 || d.Property.TargetBonuses.LargerResourceStorage > 6 {
		return t
	} else if d.Property.TargetBonuses.LargerEmeraldStorage < 0 || d.Property.TargetBonuses.LargerEmeraldStorage > 6 {
		return t
	} else if d.Property.TargetBonuses.EfficientResource < 0 || d.Property.TargetBonuses.EfficientResource > 6 {
		return t
	} else if d.Property.TargetBonuses.EfficientEmerald < 0 || d.Property.TargetBonuses.EfficientEmerald > 3 {
		return t
	} else if d.Property.TargetBonuses.ResourceRate < 0 || d.Property.TargetBonuses.ResourceRate > 3 {
		return t
	} else if d.Property.TargetBonuses.EmeraldRate < 0 || d.Property.TargetBonuses.EmeraldRate > 3 {
		return t
	}

	t.Property = d.Property
	return t
}

func (t *Territory) Undefend() *Territory {
	var zero = TerritoryProperty{
		TargetUpgrades: TerritoryPropertyUpgradeData{
			Damage:  0,
			Attack:  0,
			Health:  0,
			Defence: 0,
		},
		TargetBonuses: TerritoryPropertyBonusesData{
			StrongerMinions:       0,
			TowerMultiAttack:      0,
			TowerAura:             0,
			TowerVolley:           0,
			LargerResourceStorage: 0,
			LargerEmeraldStorage:  0,
			EfficientResource:     0,
			EfficientEmerald:      0,
			ResourceRate:          0,
			EmeraldRate:           0,
		},
		Tax: Tax{
			Ally:   5,
			Others: 5,
		},
	}

	t.Property = zero

	return t
}

func (t *Territory) CloseBorder() *Territory {
	t.Property.Border = "Closed"
	cl.Debug("Border closed for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) OpenBorder() *Territory {
	t.Property.Border = "Open"
	cl.Debug("Border opened for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) Fastest() *Territory {
	t.Property.TradingStyle = "Fastest"
	cl.Debug("Trading style set to Fastest for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) Cheapest() *Territory {
	t.Property.TradingStyle = "Cheapest"
	cl.Debug("Trading style set to Cheapest for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) ToggleBorder() *Territory {
	if t.Property.Border == "Closed" {
		t.Property.Border = "Open"
		cl.Debug("Border opened for", t.Name, "[ ID:", t.ID, "]")
	} else {
		t.Property.Border = "Closed"
		cl.Debug("Border closed for", t.Name, "[ ID:", t.ID, "]")
	}
	return t
}

func (t *Territory) SetAllyTax(n int) *Territory {
	// tax has to be within 5 and 60
	if n < 5 || n > 60 {
		cl.Warn("Ally tax out of range, ignoring...")
		return t
	}
	t.Property.Tax.Ally = n
	cl.Debug("Ally tax updated to", n, "for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) SetOthersTax(n int) *Territory {
	// tax has to be within 5 and 60
	if n < 5 || n > 60 {
		cl.Warn("Others tax out of range, ignoring...")
		return t
	}
	t.Property.Tax.Others = n
	cl.Debug("Others tax updated to", n, "for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) AddTradingRoute(r string) *Territory {
	t.TradingRoutes = append(t.TradingRoutes, r)
	cl.Log("Trading route added for", t.Name, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) SetHQ() *Territory {
	t.Property.HQ = true
	cl.Debug("HQ set for", t.Name, "[ ID:", t.ID, "]")
	//hq have 5x storage capacity
	t.Storage.Capacity.Emerald = t.Storage.Capacity.Emerald * 5
	t.Storage.Capacity.Ore = t.Storage.Capacity.Ore * 5
	t.Storage.Capacity.Wood = t.Storage.Capacity.Wood * 5
	t.Storage.Capacity.Fish = t.Storage.Capacity.Fish * 5
	t.Storage.Capacity.Crop = t.Storage.Capacity.Crop * 5
	cl.Debug("New HQ capacity", t.Storage.Capacity, "[ ID:", t.ID, "]")
	return t
}

func (t *Territory) UnsetHQ() *Territory {
	t.Property.HQ = false
	cl.Debug("HQ unset for", t.Name, "[ ID:", t.ID, "]")
	// divide the storage by 5
	t.Storage.Capacity.Emerald /= 5
	t.Storage.Capacity.Ore /= 5
	t.Storage.Capacity.Wood /= 5
	t.Storage.Capacity.Fish /= 5
	t.Storage.Capacity.Crop /= 5
	cl.Debug("New HQ capacity", t.Storage.Capacity, "[ ID:", t.ID, "]")

	return t
}

func (t *Territory) ToggleAlly() *Territory {
	if t.Ally {
		t.Ally = false
		cl.Debug("Ally status removed for", t.Name, "[ ID:", t.ID, "]")
	} else {
		t.Ally = true
		cl.Debug("Ally status added for", t.Name, "[ ID:", t.ID, "]")
	}
	return t
}

func (t *Territory) SetArbitraryStorage(rs *TerritoryResourceStorageValue) *Territory {

	t.Storage.Current.Emerald = rs.Emerald
	t.Storage.Current.Ore = rs.Ore
	t.Storage.Current.Wood = rs.Wood
	t.Storage.Current.Fish = rs.Fish
	cl.Debug("Storage updated to", t.Storage.Current, "for", t.Name, "[ ID:", t.ID, "]")

	return t
}
