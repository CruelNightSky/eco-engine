// TODO
// fix the goddamn transversing resource fr

package main

import (
	"eco-engine/table"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/RyanCarrier/dijkstra"
	"github.com/gookit/goutil/arrutil"
	"github.com/gorilla/websocket"
)

const (
	VERSION = "0.0.1a"
)

const (
	_ = iota
	ore
	crop
	wood
	fish
)

var (
	t                 map[string]*table.Territory // loaded
	loadedTerritories = make(map[string]*table.Territory)
	upgrades          *table.CostTable
	initialised       = false
	hq                string
	/*
		zeroData          = &table.TerritoryUpdateData{
			Property: table.TerritoryProperty{
				TargetUpgrades: table.TerritoryPropertyUpgradeData{
					Damage:  0,
					Attack:  0,
					Health:  0,
					Defence: 0,
				},
				TargetBonuses: table.TerritoryPropertyBonusesData{
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
				Tax: table.Tax{
					Ally:   5,
					Others: 5,
				},
				Border:       "Open",
				TradingStyle: "Cheapest",
				HQ:           false,
			},
		}
	*/
)


func init() {

	// load all upgrades data
	var bytes, err = os.ReadFile("./upgrades.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &upgrades)
	if err != nil {
		panic(err)
	}

	var uninitTerritories map[string]table.RawTerritoryData
	bytes, err = os.ReadFile("./baseProperty.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &uninitTerritories)
	if err != nil {
		panic(err)
	}

	// initialise territory
	var territories = make(map[string]*table.Territory, len(uninitTerritories))
	var counter = 0
	for name, data := range uninitTerritories {
		territories[name] = &table.Territory{
			Name: name,
			BaseResourceProduction: table.TerritoryResource{
				Emerald: data.Resources.Emeralds,
				Ore:     data.Resources.Ore,
				Wood:    data.Resources.Wood,
				Fish:    data.Resources.Fish,
				Crop:    data.Resources.Crops,
			},
			Property: table.TerritoryProperty{
				TargetUpgrades: table.TerritoryPropertyUpgradeData{
					Damage:  0,
					Attack:  0,
					Health:  0,
					Defence: 0,
				},
				TargetBonuses: table.TerritoryPropertyBonusesData{
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
				CurrentUpgrades: table.TerritoryPropertyUpgradeData{
					Damage:  0,
					Attack:  0,
					Health:  0,
					Defence: 0,
				},
				CurrentBonuses: table.TerritoryPropertyBonusesData{
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
				Tax: table.Tax{
					Ally:   5,
					Others: 5,
				},
				Border:       "Open",
				TradingStyle: "Cheapest",
				HQ:           false,
			},
			Storage: table.TerritoryResourceStorage{
				Capacity: table.TerritoryResourceStorageValue{
					Emerald: 3000,
					Ore:     300,
					Wood:    300,
					Fish:    300,
					Crop:    300,
				},
				Current: table.TerritoryResourceStorageValue{
					Emerald: 0,
					Ore:     0,
					Wood:    0,
					Fish:    0,
					Crop:    0,
				},
			},
			ResourceProduction: table.TerritoryResource{
				Emerald: data.Resources.Emeralds,
				Ore:     data.Resources.Ore,
				Wood:    data.Resources.Wood,
				Fish:    data.Resources.Fish,
				Crop:    data.Resources.Crops,
			},
			TerritoryUsage: table.TerritoryResource{
				Emerald: 0,
				Ore:     0,
				Wood:    0,
				Fish:    0,
				Crop:    0,
			},
			TradingRoutes: data.TradingRoutes,
			ID:            counter,
		}
		counter++
	}
	t = territories
}

func main() {
	var port string
	if len(os.Args) < 2 {
		log.Println("port not specified using default 8080")
		port = "8080"
	} else {
		port = os.Args[1]
	}

	// if port == "" {
	//	log.Panicln("$PORT must be set")
	// }

	// start http server

	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"code":405,"error":"method not allowed"}`))
			return
		}

		if initialised {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"code":403,"error":"already initialised"}`))
			return
		}
		// init territories specified in the request body
		var bytes, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var territories struct {
			Territories []string `json:"territories"`
			HQ          string   `json:"hq"`
		}
		json.Unmarshal(bytes, &territories)

		log.Printf("Received %v", territories)
		log.Println("initialising territories")

		hq = territories.HQ

		// if no hq provided
		if hq == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, name := range territories.Territories {

			// set all the territory properties to 0 or default
			t[name].Storage.Capacity = table.TerritoryResourceStorageValue{
				Emerald: 3000,
				Ore:     300,
				Wood:    300,
				Fish:    300,
				Crop:    300,
			}
			loadedTerritories[name] = t[name]
			if name == hq {
				log.Println(t[name].Storage.Capacity)
				t[name].SetHQ()
			}
		}

		for _, terr := range territories.Territories {
			(t)[terr].Claim = true
		}

		var terrList = make(map[string]*table.Territory)
		for _, name := range territories.Territories {
			terrList[name] = t[name]
		}

		GetPathToHQCheapest(&t, hq)
		CalculateRouteToHQTax(&t, hq)

		log.Println("initialised territories")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"message":"initialised"}`))

		log.Println("HQ Territory :", hq)
		initialised = true

		// for testing
		t[hq].Set(table.TerritoryUpdateData{
			Property: table.TerritoryProperty{
				TargetUpgrades: table.TerritoryPropertyUpgradeData{
					Damage:  11,
					Attack:  11,
					Defence: 11,
					Health:  11,
				},
				TargetBonuses: table.TerritoryPropertyBonusesData{
					LargerResourceStorage: 6,
					LargerEmeraldStorage:  6,
				},
			},
		})
		t[hq].Storage.Current = table.TerritoryResourceStorageValue{
			Emerald: 300000.0,
			Ore:     120000.0,
			Wood:    120000.0,
			Fish:    120000.0,
			Crop:    120000.0,
		}

		t["Ahmsord"].Set(table.TerritoryUpdateData{
			Property: table.TerritoryProperty{
				TargetUpgrades: table.TerritoryPropertyUpgradeData{
					Damage:  6,
					Attack:  6,
					Defence: 6,
					Health:  6,
				},
			},
		})

		t["Central Islands"].SetHQ()

		startTimer(&t, hq)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		// upgrade connection to websocket
		var upgrader = websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// close connection when function returns
		defer c.Close()

		// we will only send messages to the client
		var data struct {
			Territories []*table.Territory `json:"territories"`
		}

		// send data to client every 1s using goroutine
		go func() {
			for {

				data.Territories = make([]*table.Territory, 0)
				for _, terr := range loadedTerritories {
					data.Territories = append(data.Territories, terr)
				}
				bytes, err := json.Marshal(data)
				if err != nil {
					log.Println(err)
				}
				err = c.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(time.Second * 1)
			}
		}()
	})

	http.HandleFunc("/modifyTerritory", func(w http.ResponseWriter, r *http.Request) {

		if !initialised {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"not initialised"}`))
			return
		}

		// if not POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code":400,"message":"method not allowed"}`))
			return
		}

		// add territory to the map
		var bytes, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data struct {
			Method      string   `json:"method"`
			Territories []string `json:"territories"`
		}

		err = json.Unmarshal(bytes, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"invalid json","rawError":"` + err.Error() + `"}`))
			log.Println(err)
			return
		}

		for _, name := range data.Territories {
			// set all the territory properties to 0 or default
			if data.Method == "add" {
				loadedTerritories[name] = t[name]
				loadedTerritories[name].Claim = true
			} else if data.Method == "remove" {
				// unload the territory, set all the territory properties to 0 or default, remove from loadedTerritories and mark as unclaimed
				t[name].SetAllyTax(5).SetOthersTax(60).OpenBorder().Cheapest()
				delete(loadedTerritories, name)
				t[name].Claim = false
			}
		}

		GetPathToHQCheapest(&t, hq)
		CalculateRouteToHQTax(&t, hq)

	})

	http.HandleFunc("/setArbitraryStorage", func(w http.ResponseWriter, r *http.Request) {
		if !initialised {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"not initialised"}`))
			return
		}

		// if not POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code":400,"message":"method not allowed"}`))
			return
		}

		var requestData struct {
			Territory string                              `json:"territory"`
			Value     *table.TerritoryResourceStorageValue `json:"value"`
		}
		json.Unmarshal([]byte(r.FormValue("data")), &requestData)

		t[requestData.Territory].SetArbitraryStorage(requestData.Value)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"message":"resource storage set"}`))

	})

	http.HandleFunc("/ally", func(w http.ResponseWriter, r *http.Request) {

		if !initialised {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"not initialised"}`))
			return
		}

		// if not POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code":400,"message":"method not allowed"}`))
			return
		}

		// add territory to the map
		var bytes, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data struct {
			Method      string
			Territories []string `json:"territories"`
		}

		err = json.Unmarshal(bytes, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"invalid json","rawError":"` + err.Error() + `"}`))
			log.Println(err)
			return
		}

		for _, name := range data.Territories {
			// mark the territory as ally
			t[name].SetAllyTax(5).SetOthersTax(60).OpenBorder().Cheapest()
			if data.Method == "add" {
				t[name].Ally = true
			} else if data.Method == "remove" {
				t[name].Ally = false
			}
		}

	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if !initialised {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"not initialised"}`))
			return
		}

		// if not POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code":400,"message":"method not allowed"}`))
			return
		}

		// add territory to the map
		var bytes, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data struct {
			Method      string
			Territories []string `json:"territories"`
		}

		err = json.Unmarshal(bytes, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code":400,"message":"invalid json","rawError":"` + err.Error() + `"}`))
			log.Println(err)
			return
		}
	})

	http.ListenAndServe(":"+port, nil)

}

func startTimer(t *map[string]*table.Territory, HQ string) {

	// run generateResource every 1s and resTick every 60s using goroutine
	go func(t *map[string]*table.Territory, HQ string) {
		// log.Println("New goroutine started")
		var counter = 0
		for {
			time.Sleep(time.Second * 1)
			GenerateResorce(t)
			CalculateTowerStats(t, HQ)
			CalculateTerritoryUsageCost(t)
			CalculateTerritoryLevel(t)
			SetStorageCapacity(t)

			counter++
			log.Println("tick")
			log.Println("HQ", (*t)[HQ].Property, (*t)[HQ].Storage.Current)
			log.Println("Ahmsord", (*t)["Ahmsord"].TerritoryUsage, (*t)["Ahmsord"].Storage.Current)
			// every 60s
			if counter%5 == 0 {
				ResourceTick(t, HQ)
				ResourceTickFromHQ(t, HQ)
				RequestResourceFromHQ(t)
				log.Println("resource tick")
				counter = 0
			}
		}
	}(t, HQ)
}

func CalculateTowerStats(t *map[string]*table.Territory, HQ string) {
	for _, territory := range *t {

		// if not claimed then it will be default value
		if (*territory).Claim {

			// calculate tower stats based on the upgrade value and nearby territories

			if (*territory).Property.HQ {

				var conns, ext = 0, 0

				for _, terr := range *t {
					if len((*terr).RouteToHQ) <= 3 && (*terr).Claim {
						ext++
					}
				}
				for _, terr := range (*t)[HQ].TradingRoutes {
					if (*t)[terr].Claim {
						conns++
					}
				}

				var tr = (*territory)

				// 50% damage boost for hq
				var towerDmgLevel = tr.Property.TargetUpgrades.Damage
				var towerAtkLevel = tr.Property.TargetUpgrades.Attack
				var towerDefLevel = tr.Property.TargetUpgrades.Defence
				var towerHpLevel = tr.Property.TargetUpgrades.Health

				var baseDamageMin = upgrades.UpgradeBaseStats.Damage.Min[towerDmgLevel]
				var baseDamageMax = upgrades.UpgradeBaseStats.Damage.Max[towerDmgLevel]
				var baseAttack = upgrades.UpgradeBaseStats.Attack[towerAtkLevel]
				var baseHp = upgrades.UpgradeBaseStats.Health[towerHpLevel]
				var baseDefence = upgrades.UpgradeBaseStats.Defence[towerDefLevel]

				// hq conns and ext buff
				var dmgMin float64 = float64(baseDamageMin) * (1.5) * (1 + (0.3 * float64(conns))) * (1 + (0.25 * float64(ext)))
				var dmgMax float64 = float64(baseDamageMax) * (1.5) * (1 + (0.3 * float64(conns))) * (1 + (0.25 * float64(ext)))
				// hp
				var hp float64 = float64(baseHp) * (1 + (0.3 * float64(conns))) * (1 + (0.25 * float64(ext)))

				tr.Stats.Damage.Min = uint64(math.Round(dmgMin))
				tr.Stats.Damage.Max = uint64(math.Round(dmgMax))
				tr.Stats.Attack = float32(baseAttack)
				tr.Stats.Health = uint64(hp)
				tr.Stats.Defence = float32(baseDefence)

			}

		} else if !(*territory).Property.HQ && (*territory).Claim {

			var conns = 0

			for _, terr := range (*territory).TradingRoutes {
				if (*t)[terr].Claim {
					conns++
				}
			}

			var tr = (*territory)

			// conns gives 30% damage and hp boost to normal terrs
			var towerDmgLevel = tr.Property.CurrentUpgrades.Damage
			var towerAtkLevel = tr.Property.CurrentUpgrades.Attack
			var towerDefLevel = tr.Property.CurrentUpgrades.Defence
			var towerHpLevel = tr.Property.CurrentUpgrades.Health

			var baseDamageMin = upgrades.UpgradeBaseStats.Damage.Min[towerDmgLevel]
			var baseDamageMax = upgrades.UpgradeBaseStats.Damage.Max[towerDmgLevel]
			var baseAttack = upgrades.UpgradeBaseStats.Attack[towerAtkLevel]
			var baseHp = upgrades.UpgradeBaseStats.Health[towerHpLevel]
			var baseDefence = upgrades.UpgradeBaseStats.Defence[towerDefLevel]

			var dmgMin float64 = float64(baseDamageMin) * (1.5) * (1 + (0.3 * float64(conns)))
			var dmgMax float64 = float64(baseDamageMax) * (1.5) * (1 + (0.3 * float64(conns)))

			tr.Stats.Damage.Min = uint64(math.Round(dmgMin))
			tr.Stats.Damage.Max = uint64(math.Round(dmgMax))
			tr.Stats.Attack = float32(baseAttack)
			tr.Stats.Health = uint64(math.Round(float64(baseHp) * (1 + (0.3 * float64(conns)))))
			tr.Stats.Defence = float32(baseDefence)

		}

		(*territory).Stats.Damage.Min = 1_000
		(*territory).Stats.Damage.Max = 1_500
		(*territory).Stats.Attack = 0.5
		(*territory).Stats.Health = 300_000
		(*territory).Stats.Defence = 10

		(*territory).Stats.StrongerMinions = 0
		(*territory).Stats.TowerMultiAttacks = 0
		(*territory).Stats.TowerAura = 0
		(*territory).Stats.TowerVolley = 0
	}

	UseResource(t)
}

func UseResource(t *map[string]*table.Territory) {

	log.Println("UseResource()")
	for _, territory := range *t {

		// offload the calculation to another goroutine
		go func(t *table.Territory) {
			// log.Println("New goroutine started for terriotry: ", (*t).Name)

			var towerDmgLevel = (*t).Property.TargetUpgrades.Damage
			var towerAtkLevel = (*t).Property.TargetUpgrades.Attack
			var towerDefLevel = (*t).Property.TargetUpgrades.Defence
			var towerHpLevel = (*t).Property.TargetUpgrades.Health

			var strongerMinionsLevel = (*t).Property.TargetBonuses.StrongerMinions
			var towerMultiAttackLevel = (*t).Property.TargetBonuses.TowerMultiAttack
			var towerAuraLevel = (*t).Property.TargetBonuses.TowerAura
			var towerVolleyLevel = (*t).Property.TargetBonuses.TowerVolley
			var largerEmeraldsStorageLevel = (*t).Property.TargetBonuses.LargerEmeraldStorage
			var largerResourceStorageLevel = (*t).Property.TargetBonuses.LargerResourceStorage
			var efficientResourceLevel = (*t).Property.TargetBonuses.EfficientResource
			var efficientEmeraldLevel = (*t).Property.TargetBonuses.EfficientEmerald
			var resourceRateLevel = (*t).Property.TargetBonuses.ResourceRate
			var emeraldRateLevel = (*t).Property.TargetBonuses.EmeraldRate

			var damageCost = upgrades.UpgradesCost.Damage.Value[towerDmgLevel]
			var attackCost = upgrades.UpgradesCost.Attack.Value[towerAtkLevel]
			var defenceCost = upgrades.UpgradesCost.Defence.Value[towerDefLevel]
			var healthCost = upgrades.UpgradesCost.Health.Value[towerHpLevel]

			var strongerMinionsCost = upgrades.Bonuses.StrongerMinions.Cost[strongerMinionsLevel]
			var towerMultiAttackCost = upgrades.Bonuses.TowerMultiAttack.Cost[towerMultiAttackLevel]
			var towerAuraCost = upgrades.Bonuses.TowerAura.Cost[towerAuraLevel]
			var towerVolleyCost = upgrades.Bonuses.TowerVolley.Cost[towerVolleyLevel]

			var largerEmeraldsStorageCost = upgrades.Bonuses.LargerEmeraldsStorage.Cost[largerEmeraldsStorageLevel]
			var largerResourceStorageCost = upgrades.Bonuses.LargerResourceStorage.Cost[largerResourceStorageLevel]
			var efficientResourceCost = upgrades.Bonuses.EfficientResource.Cost[efficientResourceLevel]
			var efficientEmeraldCost = upgrades.Bonuses.EfficientEmeralds.Cost[efficientEmeraldLevel]
			var resourceRateCost = upgrades.Bonuses.ResourceRate.Cost[resourceRateLevel]
			var emeraldRateCost = upgrades.Bonuses.EmeraldsRate.Cost[emeraldRateLevel]

			var emeraldCostPerSec = float64(largerResourceStorageCost+efficientResourceCost+resourceRateCost) / 3600 // 1 hour
			var oreCostPerSec = float64(damageCost+towerVolleyCost+efficientEmeraldCost) / 3600
			var cropCostPerSec = float64(attackCost+towerAuraCost+emeraldRateCost) / 3600
			var woodCostPerSec = float64(healthCost+strongerMinionsCost+largerEmeraldsStorageCost) / 3600
			var fishCostPerSec = float64(defenceCost+towerMultiAttackCost) / 3600

			if (*t).Storage.Current.Emerald < emeraldCostPerSec {
				(*t).Property.CurrentBonuses.LargerResourceStorage = 0
				(*t).Property.CurrentBonuses.EfficientResource = 0
				(*t).Property.CurrentBonuses.ResourceRate = 0
			} else if (*t).Storage.Current.Ore >= oreCostPerSec {
				(*t).Property.CurrentBonuses.LargerResourceStorage = (*t).Property.TargetBonuses.LargerResourceStorage
				(*t).Property.CurrentBonuses.EfficientResource = (*t).Property.TargetBonuses.EfficientResource
				(*t).Property.CurrentBonuses.ResourceRate = (*t).Property.TargetBonuses.ResourceRate
			}

			if (*t).Storage.Current.Ore < oreCostPerSec {
				(*t).Property.CurrentUpgrades.Damage = 0
				(*t).Property.CurrentBonuses.TowerVolley = 0
				(*t).Property.CurrentBonuses.EfficientEmerald = 0
			} else if (*t).Storage.Current.Crop >= cropCostPerSec {
				(*t).Property.CurrentUpgrades.Damage = (*t).Property.TargetUpgrades.Damage
				(*t).Property.CurrentBonuses.TowerVolley = (*t).Property.TargetBonuses.TowerVolley
				(*t).Property.CurrentBonuses.EfficientEmerald = (*t).Property.TargetBonuses.EfficientEmerald
			}

			if (*t).Storage.Current.Crop < cropCostPerSec {
				(*t).Property.CurrentUpgrades.Attack = 0
				(*t).Property.CurrentBonuses.TowerAura = 0
				(*t).Property.CurrentBonuses.EmeraldRate = 0
			} else if (*t).Storage.Current.Crop >= cropCostPerSec {
				(*t).Property.CurrentUpgrades.Attack = (*t).Property.TargetUpgrades.Attack
				(*t).Property.CurrentBonuses.TowerAura = (*t).Property.TargetBonuses.TowerAura
				(*t).Property.CurrentBonuses.EmeraldRate = (*t).Property.TargetBonuses.EmeraldRate
			}

			if (*t).Storage.Current.Wood < woodCostPerSec {
				(*t).Property.CurrentUpgrades.Health = 0
				(*t).Property.CurrentBonuses.StrongerMinions = 0
				(*t).Property.CurrentBonuses.LargerEmeraldStorage = 0
			} else if (*t).Storage.Current.Wood >= woodCostPerSec {
				(*t).Property.CurrentUpgrades.Health = (*t).Property.TargetUpgrades.Health
				(*t).Property.CurrentBonuses.StrongerMinions = (*t).Property.TargetBonuses.StrongerMinions
				(*t).Property.CurrentBonuses.LargerEmeraldStorage = (*t).Property.TargetBonuses.LargerEmeraldStorage
			}

			if (*t).Storage.Current.Fish < fishCostPerSec {
				(*t).Property.CurrentUpgrades.Defence = 0
				(*t).Property.CurrentBonuses.TowerMultiAttack = 0
			} else if (*t).Storage.Current.Fish >= fishCostPerSec {
				(*t).Property.CurrentUpgrades.Defence = (*t).Property.TargetUpgrades.Defence
				(*t).Property.CurrentBonuses.TowerMultiAttack = (*t).Property.TargetBonuses.TowerMultiAttack
			}

		}(territory)
	}
}

func SetStorageCapacity(territories *map[string]*table.Territory) {
	for _, territory := range *territories {
		var largerResourceStorageLevel = (*territory).Property.CurrentBonuses.LargerResourceStorage
		var largerEmeraldsStorageLevel = (*territory).Property.CurrentBonuses.LargerEmeraldStorage

		if !(*territory).Property.HQ {
			(*territory).Storage.Capacity.Emerald = float64(3000 * upgrades.Bonuses.LargerEmeraldsStorage.Value[largerEmeraldsStorageLevel])
			(*territory).Storage.Capacity.Ore = float64(300 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Crop = float64(300 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Wood = float64(300 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Fish = float64(300 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
		} else {
			(*territory).Storage.Capacity.Emerald = float64(5000 * upgrades.Bonuses.LargerEmeraldsStorage.Value[largerEmeraldsStorageLevel])
			(*territory).Storage.Capacity.Ore = float64(1500 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Crop = float64(1500 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Wood = float64(1500 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
			(*territory).Storage.Capacity.Fish = float64(1500 * upgrades.Bonuses.LargerResourceStorage.Value[largerResourceStorageLevel])
		}
	}
}

func GenerateResorce(territories *map[string]*table.Territory) {
	// rate means how many seconds it takes to generate n resource
	// n resource is calculated like this
	// nres = base res prod * efficient resource
	// so if rate is level 0 then it takes 1 second to generate 1/4 of the resource
	// and if resource stored in the storage excees the capacity then the excess resource will be lost
	// stoarge capacity is calculated like this
	// cap = base cap * larger storage

	for _, territory := range *territories {

		// emerald generation
		var baseEmeraldGeneration = (*territory).BaseResourceProduction.Emerald
		var efficientEmerald = (*territory).Property.TargetBonuses.EfficientEmerald
		var emeraldRate = (*territory).Property.TargetBonuses.EmeraldRate

		var emeraldMultiplier = float64(upgrades.Bonuses.EfficientEmeralds.Value[efficientEmerald]) * (4 / float64(upgrades.Bonuses.EmeraldsRate.Value[emeraldRate]))

		var emeraldGenerationPerSec = (float64(baseEmeraldGeneration) * emeraldMultiplier) / 3600

		// add the emerald to storage
		if (*territory).Storage.Current.Emerald < (*territory).Storage.Capacity.Emerald {
			(*territory).Storage.Current.Emerald += emeraldGenerationPerSec
		} else {
			(*territory).Storage.Current.Emerald = (*territory).Storage.Capacity.Emerald
		}

		// resource generation
		// check what kind of territory is this first
		var territoryType string

		if (*territory).BaseResourceProduction.Crop != 0 {
			territoryType = "Crop"
		} else if (*territory).BaseResourceProduction.Wood != 0 {
			territoryType = "Wood"
		} else if (*territory).BaseResourceProduction.Ore != 0 {
			territoryType = "Ore"
		} else if (*territory).BaseResourceProduction.Fish != 0 {
			territoryType = "Fish"
		} else if (*territory).BaseResourceProduction.Fish != 0 && (*territory).BaseResourceProduction.Crop != 0 {

			// gotta accomodate for that one stupid terr in ragni area
			territoryType = "FishCrop"
		}

		if territory.Name != "Maltic Coast" {

			// get struct field by string
			var baseResourceGeneration = reflect.ValueOf((*territory).BaseResourceProduction).FieldByName(territoryType).Float()
			var efficientResource = (*territory).Property.CurrentBonuses.EfficientResource
			var resourceRate = (*territory).Property.CurrentBonuses.ResourceRate

			var resourceMultiplier = float64(upgrades.Bonuses.EfficientResource.Value[efficientResource]) * (4 / float64(upgrades.Bonuses.ResourceRate.Value[resourceRate]))
			var resourceGenerationPerSec = (float64(baseResourceGeneration) * resourceMultiplier) / 3600

			// add the resource to storage using runtime reflection since we dont know what kind of resource it is
			var resourceStorage = reflect.ValueOf((*territory).Storage.Current).FieldByName(territoryType).Float()
			var resourceStorageCapacity = reflect.ValueOf((*territory).Storage.Capacity).FieldByName(territoryType).Float()

			// use reflection to set the value of the field
			var v = reflect.ValueOf(&territory.Storage.Current).Elem()
			var f = v.FieldByName(territoryType)

			// check if the field is valid and can be set
			if f.IsValid() && f.CanSet() {
				if resourceStorage < resourceStorageCapacity {
					f.SetFloat(resourceStorage + float64(resourceGenerationPerSec))
				} else {
					f.SetFloat(resourceStorageCapacity)
				}
			} else {
				_ = fmt.Errorf("field is not valid or cannot be set for %s", territory.Name)
			}
		} else {

			// for the maltic plains
			var baseResourceGenerationCrop = (*territory).BaseResourceProduction.Crop
			var baseResourceGenerationFish = (*territory).BaseResourceProduction.Fish

			var efficientResource = (*territory).Property.CurrentBonuses.EfficientResource
			var resourceRate = (*territory).Property.CurrentBonuses.ResourceRate

			var resourceMultiplier = float64(upgrades.Bonuses.EfficientResource.Value[efficientResource]) * (4 / float64(upgrades.Bonuses.ResourceRate.Value[resourceRate]))

			var resourceGenerationPerSecCrop = (float64(baseResourceGenerationCrop) * resourceMultiplier) / 3600
			var resourceGenerationPerSecFish = (float64(baseResourceGenerationFish) * resourceMultiplier) / 3600

			// for crops
			if (*territory).Storage.Current.Crop <= (*territory).Storage.Capacity.Crop {
				(*territory).Storage.Current.Crop += resourceGenerationPerSecCrop
			} else {
				(*territory).Storage.Current.Crop = (*territory).Storage.Capacity.Crop
			}

			// for fish
			if (*territory).Storage.Current.Fish <= (*territory).Storage.Capacity.Fish {
				(*territory).Storage.Current.Fish += resourceGenerationPerSecFish
			} else {
				(*territory).Storage.Current.Fish = (*territory).Storage.Capacity.Fish
			}
		}
	}
}

func CalculateTerritoryUsageCost(territories *map[string]*table.Territory) {
	for _, territory := range *territories {

		//ignore if the territory is not claimed
		if !territory.Claim {
			continue
		}

		var usage = table.TerritoryResource{
			Emerald: 0,
			Ore:     0,
			Wood:    0,
			Crop:    0,
			Fish:    0,
		}

		var damage, attack, hp, defence int
		damage = territory.Property.TargetUpgrades.Damage
		attack = territory.Property.TargetUpgrades.Attack
		hp = territory.Property.TargetUpgrades.Health
		defence = territory.Property.TargetUpgrades.Defence

		// calculate the upgrade usage cost of the territory
		usage.Ore += float64(upgrades.UpgradesCost.Damage.Value[damage])
		usage.Crop += float64(upgrades.UpgradesCost.Attack.Value[attack])
		usage.Wood += float64(upgrades.UpgradesCost.Health.Value[hp])
		usage.Fish += float64(upgrades.UpgradesCost.Defence.Value[defence])

		// bonuses
		var strongerMinions = territory.Property.TargetBonuses.StrongerMinions
		var towerMultiAttacks = territory.Property.TargetBonuses.TowerMultiAttack
		var aura = territory.Property.TargetBonuses.TowerAura
		var volley = territory.Property.TargetBonuses.TowerVolley

		var efficientResource = territory.Property.TargetBonuses.EfficientResource
		var resourceRate = territory.Property.TargetBonuses.ResourceRate
		var efficientEmerald = territory.Property.TargetBonuses.EfficientEmerald
		var emeraldRate = territory.Property.TargetBonuses.EmeraldRate

		var emStorage = territory.Property.TargetBonuses.LargerEmeraldStorage
		var resStorage = territory.Property.TargetBonuses.LargerResourceStorage

		usage.Ore += float64(upgrades.Bonuses.StrongerMinions.Cost[volley] + upgrades.Bonuses.EfficientEmeralds.Cost[efficientEmerald])
		usage.Crop += float64(upgrades.Bonuses.TowerAura.Cost[aura] + upgrades.Bonuses.EmeraldsRate.Cost[emeraldRate])
		usage.Wood += float64(upgrades.Bonuses.StrongerMinions.Cost[strongerMinions] + upgrades.Bonuses.LargerEmeraldsStorage.Cost[emStorage] + upgrades.Bonuses.StrongerMinions.Cost[strongerMinions])
		usage.Fish += float64(upgrades.Bonuses.TowerMultiAttack.Cost[towerMultiAttacks])
		usage.Emerald += float64(upgrades.Bonuses.LargerResourceStorage.Cost[resStorage] + upgrades.Bonuses.ResourceRate.Cost[resourceRate] + upgrades.Bonuses.EfficientResource.Cost[efficientResource])

		// set the usage cost of the territory
		(*territory).TerritoryUsage = usage
	}
}

func RequestResourceFromHQ(territories *map[string]*table.Territory) {

	var HQ = (*territories)[hq]

	for _, territory := range *territories {

		// if HQ then ignore
		if territory.Name == hq {
			runtime.Breakpoint()
			continue
		}

		// if the territory is not claimed then ignore
		if !territory.Claim {
			continue
		}

		var resourceUsage = territory.TerritoryUsage

		if resourceUsage.Emerald != 0 || resourceUsage.Ore != 0 || resourceUsage.Wood != 0 || resourceUsage.Crop != 0 || resourceUsage.Fish != 0 {

			(*HQ).TransversingResourceFromHQ = append((*HQ).TransversingResourceFromHQ, table.TransveringResource{
				Source:      (*territories)[hq].Name,
				Emerald:     math.Max(float64(resourceUsage.Emerald-(*territory).ResourceProduction.Emerald), 0),
				Ore:         math.Max(float64(resourceUsage.Ore-(*territory).ResourceProduction.Ore), 0),
				Wood:        math.Max(float64(resourceUsage.Wood-(*territory).ResourceProduction.Wood), 0),
				Crop:        math.Max(float64(resourceUsage.Crop-(*territory).ResourceProduction.Crop), 0),
				Fish:        math.Max(float64(resourceUsage.Fish-(*territory).ResourceProduction.Fish), 0),
				Destination: (*territory).Name,
				RouteToDest: (*territory).RouteFromHQ,
			})

			log.Println((*HQ).TransversingResourceFromHQ)

		}

	}
}

func ResourceTickFromHQ(t *map[string]*table.Territory, HQ string) {

	var visited = make(map[string]bool)

	for _, territory := range *t {

		if visited[territory.Name] {
			continue
		}

		for _, transversingResource := range (*territory).TransversingResourceFromHQ {
			var dest = transversingResource.Destination
			var route = transversingResource.RouteToDest

			// move the transversing resource to the next territory's transversing resource
			// if the territory is the destination then move it onto the territory's storage
			// then remove the current transversing resource from territory to not cause memory leak

			// if its the destination then move it onto the territory's storage
			for _, territoryRoute := range route {

				// next terr is PathFromHQ[1]
				if (*territory).Property.HQ {

					// problematic line
					if (*territory).Storage.Current.Emerald < transversingResource.Emerald ||
						(*territory).Storage.Current.Ore < transversingResource.Ore ||
						(*territory).Storage.Current.Wood < transversingResource.Wood ||
						(*territory).Storage.Current.Crop < transversingResource.Crop ||
						(*territory).Storage.Current.Fish < transversingResource.Fish {

						// only push what we have
						(*t)[territoryRoute].TransversingResourceFromHQ = append((*t)[territoryRoute].TransversingResourceFromHQ, table.TransveringResource{
							Emerald:     math.Min((*t)[dest].TerritoryUsage.Emerald, (*territory).Storage.Current.Emerald),
							Ore:         math.Min((*t)[dest].TerritoryUsage.Ore, (*territory).Storage.Current.Ore),
							Wood:        math.Min((*t)[dest].TerritoryUsage.Wood, (*territory).Storage.Current.Wood),
							Crop:        math.Min((*t)[dest].TerritoryUsage.Crop, (*territory).Storage.Current.Crop),
							Fish:        math.Min((*t)[dest].TerritoryUsage.Fish, (*territory).Storage.Current.Fish),
							Destination: dest,
							RouteToDest: route,
						})
					}

					// and set the current storage to 0
					(*territory).Storage.Current.Emerald = 0
					(*territory).Storage.Current.Ore = 0
					(*territory).Storage.Current.Wood = 0
					(*territory).Storage.Current.Crop = 0
					(*territory).Storage.Current.Fish = 0

				} else if territoryRoute == dest {

					// move onto real storage
					(*t)[territoryRoute].Storage.Current.Emerald += transversingResource.Emerald
					(*t)[territoryRoute].Storage.Current.Ore += transversingResource.Ore
					(*t)[territoryRoute].Storage.Current.Wood += transversingResource.Wood
					(*t)[territoryRoute].Storage.Current.Crop += transversingResource.Crop
					(*t)[territoryRoute].Storage.Current.Fish += transversingResource.Fish

				} else {

					// move onto transversing resource
					(*t)[territoryRoute].TransversingResourceFromHQ = append((*t)[territoryRoute].TransversingResourceFromHQ, transversingResource)

					// TODO: remove from the current transversing queue

				}

				// mark as visit
				visited[territoryRoute] = true
			}
		}
	}
}

func ResourceTick(territories *map[string]*table.Territory, HQ string) {

	var visited = make(map[string]bool)

	for _, t := range *territories {
		if visited[(*t).Name] || (*t).Name == HQ {
			continue
		}

		// just follow the calculated path
		// if the next terr is hq, then add the resources to the hq
		// else add the resources to the next terr's transversing resource
		for _, transversingResource := range (*t).TransversingResourceToHQ {
			var dest = transversingResource.Destination
			var route = transversingResource.RouteToDest

			if route[1] == dest {

				// move onto real storage
				(*territories)[route[1]].Storage.Current.Emerald += transversingResource.Emerald
				(*territories)[route[1]].Storage.Current.Ore += transversingResource.Ore
				(*territories)[route[1]].Storage.Current.Crop += transversingResource.Crop
				(*territories)[route[1]].Storage.Current.Wood += transversingResource.Wood
				(*territories)[route[1]].Storage.Current.Fish += transversingResource.Fish

				// shift the array to the left to remove current terr
				route = route[1:]
				transversingResource.RouteToDest = route

				// then remove
				(*t).TransversingResourceToHQ = append((*t).TransversingResourceToHQ[:0], (*t).TransversingResourceToHQ[1:]...)

				visited[(*t).Name] = true

			} else {

				// shift the array to the left to remove current terr
				route = route[1:]
				transversingResource.RouteToDest = route

				// copy the transversing resource to the next territory
				(*(*territories)[route[1]]).TransversingResourceToHQ = append((*(*territories)[route[1]]).TransversingResourceToHQ, transversingResource)

				// remove the current transversing resource
				(*t).TransversingResourceToHQ = (*t).TransversingResourceToHQ[1:]

				visited[(*t).Name] = true

			}
		}

	}
}

func GetPathToHQCheapest(territories *map[string]*table.Territory, HQ string) {

	// name is the HQ territory name
	// get path to hq using dijkstra, depending on the trading style
	// fastest  int
	// will find the shortest path while cheapest will find the shortest path with the least GLOBAL tax
	// if the territory is the hq then return empty array

	// connected nodes (territory) can be found at territories[name].TradingRoutes
	var graph = dijkstra.NewGraph()
	var HQID int

	// find the id of hq
	for _, territory := range *territories {
		if territory.Property.HQ {
			HQID = territory.ID
			break
		}
	}

	var vertexAdded = make(map[int]bool)

	for name := range *territories {
		// Add logic to compute the shortest path to HQ using Dijkstra's algorithm
		// add current node
		if vertexAdded[(*territories)[name].ID] {
			log.Println("Vertex ID:", (*territories)[name].ID, "already added")
			continue
		} else {
			graph.AddVertex((*territories)[name].ID)
			// log.Println("Added vertex ID:", territories[name].ID)
		}
	}

	// now add arc
	for _, territory := range *territories {

		var currTerr = territory.ID
		// log.Println("Current territory ID:", currTerr, " ", territory.Name)
		for _, route := range territory.TradingRoutes {

			var currConn = (*territories)[route].ID

			// log.Println("Connection ID", currConn, " ", route)

			// distance is the tax value
			var distance = float64((*territories)[route].Property.Tax.Others)
			var err = graph.AddArc(currTerr, currConn, int64(distance))
			if err != nil {
				log.Println(err)
			}

		}
	}

	// get terr id
	for _, territory := range *territories {
		if territory.Property.HQ {
			continue
		}

		var terrID = territory.ID
		var pathToHQRaw, err = graph.ShortestSafe(terrID, HQID)
		if err != nil {
			log.Println(err)
		}

		// Assign path to HQ to the territory
		// Convert terr ID to terr name and store in pathToHQ
		var pathList = pathToHQRaw.Path
		var path = make([]string, len(pathList))
		for i, id := range pathList {
			for _, terr := range *territories {
				if terr.ID == id {
					path[i] = terr.Name
					break
				}
			}
		}

		(*territory).RouteToHQ = path
		var reversePath = make([]string, len(path))
		copy(reversePath, path)
		arrutil.Reverse(reversePath)
		(*territory).RouteFromHQ = reversePath
	}
}

func GetPathToHQFastest(t *map[string]*table.Territory, HQ string) {

	// var dist int64 = 1
	// var path []string
	var graph = dijkstra.NewGraph()
	var HQID int

	// find the id of hq
	for _, territory := range *t {
		if territory.Property.HQ {
			// fmt.Println(territory.ID)
			HQID = territory.ID
			break
		}
	}

	var vertexAdded = make(map[int]bool)

	for name := range *t {
		// Add logic to compute the shortest path to HQ using Dijkstra's algorithm
		// add current node
		if vertexAdded[(*t)[name].ID] {
			log.Println("Vertex ID:", (*t)[name].ID, "already added")
			continue
		} else {
			graph.AddVertex((*t)[name].ID)
			// log.Println("Added vertex ID:", territories[name].ID)
		}
	}

	// now add arc
	for _, territory := range *t {

		var currTerr = territory.ID
		// log.Println("Current territory ID:", currTerr, " ", territory.Name)
		for _, route := range territory.TradingRoutes {

			var currConn = (*t)[route].ID

			// log.Println("Connection ID", currConn, " ", route)

			// distance is always 1
			var err = graph.AddArc(currTerr, currConn, 1)
			if err != nil {
				log.Println(err)
			}

		}
	}

	log.Println(HQID)

	// get terr id
	for _, territory := range *t {

		if territory.Property.HQ {
			continue
		}

		var terrID = territory.ID
		var pathToHQRaw, err = graph.ShortestSafe(terrID, HQID)

		if err != nil {
			log.Println(err)
		}

		// assign path to hq to the territory
		// convert terr id to terr name and store in pathToHQ
		var pathList = pathToHQRaw.Path
		var path = make([]string, len(pathList))
		var counter = 0
		for _, id := range pathList {
			for _, terr := range *t {
				if terr.ID == id {
					path[counter] = terr.Name
					counter += 1
				}
			}
		}
		counter = 0
		territory.RouteToHQ = path
		var reversePath = make([]string, len(path))
		copy(reversePath, path)
		arrutil.Reverse(reversePath)
		territory.RouteFromHQ = reversePath
	}
}

func CalculateRouteToHQTax(territories *map[string]*table.Territory, HQ string) {
	// Iterate through all territories and calculate the route tax to the HQ
	for _, territory := range *territories {
		// Skip the HQ territory
		if territory.Property.HQ {
			continue
		}

		// Calculate the route to the HQ
		routeToHQ := (*territory).RouteToHQ

		// Initialize the taxList
		taxList := []float64{}

		// Iterate through the route to HQ and get the tax of each territory
		for _, terr := range routeToHQ {
			if (*territories)[terr].Property.HQ {
				continue
			}

			if (*territories)[terr].Ally {

				// Use ally tax instead of others tax
				taxList = append(taxList, 1-float64((*territories)[terr].Property.Tax.Ally)/100)

			} else if (*territories)[terr].Claim {

				// if we claimed the territory, the tax should be 0
				taxList = append(taxList, 1)

			} else {

				// Use others tax
				taxList = append(taxList, 1-float64((*territories)[terr].Property.Tax.Others)/100)

			}

			// Calculate the route tax
			routeTax := 1.0
			for _, tax := range taxList {
				routeTax *= tax
			}

			routeTax *= 100

			// round it to 2 decimal places
			territory.RouteTax = math.Round((100-routeTax)*100) / 100
		}
	}
}

func CalculateTerritoryLevel(territories *map[string]*table.Territory) {
	for _, territory := range *territories {

		// add basic attack up
		var td = (*territory).Property.CurrentUpgrades
		var terrDef = td.Damage + td.Attack + td.Health + td.Defence

		// if aura is present
		if (*territory).Property.CurrentBonuses.TowerAura > 0 {
			terrDef += (*territory).Property.CurrentBonuses.TowerAura + 5
		}

		// if volley present
		if (*territory).Property.CurrentBonuses.TowerVolley > 0 {
			terrDef += (*territory).Property.CurrentBonuses.TowerVolley + 3
		}

		(*territory).RawLevel = terrDef

		// vlow = 0 - 5, low = 6 - 18, med = 19 - 30, high = 31 - 48, vhigh >= 49
		switch {
		case terrDef >= 0 && terrDef <= 5:
			(*territory).Level = "Very Low"
		case terrDef >= 6 && terrDef <= 18:
			(*territory).Level = "Low"
		case terrDef >= 19 && terrDef <= 30:
			(*territory).Level = "Medium"
		case terrDef >= 31 && terrDef <= 48:
			(*territory).Level = "High"
		case terrDef >= 49:
			(*territory).Level = "Very High"
		default:
			(*territory).Level = "An error has occured"
		}

		if terrDef > 0 {
			log.Println(td.Damage, td.Attack, td.Health, td.Defence, (territory.Level), terrDef, (*territory).Name)
		}

	}
}

/*
func checkForUpdate() {

	// check for new version of eco-engine from github then verify the checksum
	// if checksum isn't the same as current version, download the new version and execve() it
	// if checksum is the same, continue with the current version

	// get checksum of current version
	var hasher = sha256.New()
	var f, err = os.Open("./eco-engine")
	if err != nil {
		log.Println("Unable to check for updates:", err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		log.Println("Unable to check for updates:", err)
		return
	}

	var currentChecksum = hex.EncodeToString(hasher.Sum(nil))


}
*/
