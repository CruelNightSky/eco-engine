// TODO
// 1 HQ Storage is always 5x of normal storage
// 2 Fix routh finding unction so it checks for existing node first before adding vertex

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
	"time"

	"github.com/RyanCarrier/dijkstra"
	"github.com/gookit/goutil/arrutil"
)

var (
	t                 map[string]*table.Territory // loaded
	loadedTerritories = make(map[string]*table.Territory)
	upgrades          *table.CostTable
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
				Upgrades: table.TerritoryPropertyUpgradeData{
					Damage:  0,
					Attack:  0,
					Health:  0,
					Defence: 0,
				},
				Bonuses: table.TerritoryPropertyBonusesData{
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
		// init territories specified in the request body
		var bytes, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		/*
					// accept websocket connection
					c, err := websocket.Accept(w, r, nil)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
						os.Exit(1)
						return
					}
					defer c.Close(websocket.StatusInternalError, "Internal Error")

					ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
					defer cancel()

					// read ws message
					var v interface{}
					err = wsjson.Read(ctx, c, &v)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
					}

					// unmarshal message
					var data struct {
						Type      string `json:"type"`
						Territory string
						Data      struct {
							Upgrades     table.TerritoryPropertyUpgradeData `json:"upgrades"`
							Bonuses      table.TerritoryPropertyBonusesData `json:"bonuses"`
							Tax          table.Tax                          `json:"tax"`
							Border       string                             `json:"border"`
							TradingStyle string                             `json:"tradingStyle"`
							Claim        bool                               `json:"claimed"` // if true, ally cannot be true
							Ally         bool                               `json:"ally"`
							HQ           bool                               `json:"hq"` // if true, we have to recalculate PathToHQ and Tax for all territories
						} `json:"data"`
					}

				json.Unmarshal(bytes, &data)


			// if no territory provided
			if data.Territory == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		*/
		var territories struct {
			Territories []string `json:"territories"`
			HQ          string   `json:"hq"`
		}
		json.Unmarshal(bytes, &territories)

		log.Printf("Received %v", territories)
		log.Println("initialising territories")

		var hq = territories.HQ

		// if no hq provided
		if hq == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var zeroData = table.TerritoryUpdateData{
			Property: table.TerritoryProperty{
				Upgrades: table.TerritoryPropertyUpgradeData{
					Damage:  0,
					Attack:  0,
					Health:  0,
					Defence: 0,
				},
				Bonuses: table.TerritoryPropertyBonusesData{
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

		for _, name := range territories.Territories {

			// set all the territory properties to 0 or default
			t[name].Set(zeroData).SetAllyTax(5).SetOthersTax(60).OpenBorder().Cheapest()
			loadedTerritories[name] = t[name]
			if name == hq {
				t[name].SetHQ()
			}
		}

		var terrList = make(map[string]*table.Territory)
		for _, name := range territories.Territories {
			terrList[name] = t[name]
		}

		getPathToHQCheapest(&t, hq)
		defer CalculateRouteToHQTax(&t, hq)

		log.Println("initialised territories")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"message":"initialised"}`))

		startTimer(&t)
	})

	http.ListenAndServe(":"+port, nil)

}

func startTimer(t *map[string]*table.Territory) {
	// run generateResource every 1s and resTick every 60s using goroutine
	go func(t *map[string]*table.Territory) {
		var counter = 0
		for {
			time.Sleep(time.Second * 1)
			generateResorce(t)
			counter++
			log.Println("tick")
			// every 60s
			if counter%60 == 0 {
				resourceTick(t)
				log.Println("resource tick")
				counter = 0
			}
		}
	}(t)
}

func generateResorce(territories *map[string]*table.Territory) {
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
		var efficientEmerald = (*territory).Property.Bonuses.EfficientEmerald
		var emeraldRate = (*territory).Property.Bonuses.EmeraldRate

		var emeraldMultiplier = float64(upgrades.Bonuses.EfficientEmeralds.Value[efficientEmerald]) * (4 / float64(upgrades.Bonuses.EmeraldsRate.Value[emeraldRate]))

		var emeraldGenerationPerSec = (float64(baseEmeraldGeneration) * emeraldMultiplier) / 3600
		var emeraldStorage = (*territory).Storage.Current.Emerald
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

		if territory.Name != "Maltic Plains" {

			// get struct field by string
			var baseResourceGeneration = reflect.ValueOf((*territory).BaseResourceProduction).FieldByName(territoryType).Int()
			var efficientResource = (*territory).Property.Bonuses.EfficientResource
			var resourceRate = (*territory).Property.Bonuses.ResourceRate

			var resourceMultiplier = float64(upgrades.Bonuses.EfficientResource.Value[efficientResource]) * (4 / float64(upgrades.Bonuses.ResourceRate.Value[resourceRate]))
			var resourceGenerationPerSec = (float64(baseResourceGeneration) * resourceMultiplier) / 3600

			// add the resource to storage using runtime reflection since we dont know what kind of resource it is
			var resourceStorage = reflect.ValueOf((*territory).Storage.Current).FieldByName(territoryType).Float()
			var resourceStorageCapacity = reflect.ValueOf((*territory).Storage.Capacity).FieldByName(territoryType).Float()

			// use reflection to set the value of the field
			var v = reflect.ValueOf(&territory.Storage.Current).Elem()
			var f = v.FieldByName(territoryType)
			if f.IsValid() && f.CanSet() {
				log.Println("setting field")
				if resourceStorage < resourceStorageCapacity {
					f.SetFloat(resourceStorage + float64(resourceGenerationPerSec))
				} else {
					f.SetFloat(resourceStorageCapacity)
				}
			} else {
				_ = fmt.Errorf("field is not valid or cannot be set")
			}
			log.Println("Emerald Production :", efficientEmerald, emeraldRate, emeraldMultiplier, emeraldGenerationPerSec, "storage :", emeraldStorage)
			log.Println("Resource Production :", efficientResource, resourceRate, resourceMultiplier, resourceGenerationPerSec, "storage :", resourceStorage, resourceStorageCapacity)

		} else {

			// for the maltic plains
			var baseResourceGenerationCrop = (*territory).BaseResourceProduction.Crop
			var baseResourceGenerationFish = (*territory).BaseResourceProduction.Fish

			var efficientResource = (*territory).Property.Bonuses.EfficientResource
			var resourceRate = (*territory).Property.Bonuses.ResourceRate

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
		damage = territory.Property.Upgrades.Damage
		attack = territory.Property.Upgrades.Attack
		hp = territory.Property.Upgrades.Health
		defence = territory.Property.Upgrades.Defence

		// calculate the upgrade usage cost of the territory
		usage.Ore += upgrades.UpgradesCost.Damage.Value[damage]
		usage.Crop += upgrades.UpgradesCost.Attack.Value[attack]
		usage.Wood += upgrades.UpgradesCost.Health.Value[hp]
		usage.Fish += upgrades.UpgradesCost.Defence.Value[defence]

		// bonuses
		var strongerMinions = territory.Property.Bonuses.StrongerMinions
		var towerMultiAttacks = territory.Property.Bonuses.TowerMultiAttack
		var aura = territory.Property.Bonuses.TowerAura
		var volley = territory.Property.Bonuses.TowerVolley

		var efficientResource = territory.Property.Bonuses.EfficientResource
		var resourceRate = territory.Property.Bonuses.ResourceRate
		var efficientEmerald = territory.Property.Bonuses.EfficientEmerald
		var emeraldRate = territory.Property.Bonuses.EmeraldRate

		var emStorage = territory.Property.Bonuses.LargerEmeraldStorage
		var resStorage = territory.Property.Bonuses.LargerResourceStorage

		usage.Ore += int(upgrades.Bonuses.StrongerMinions.Cost[volley] + upgrades.Bonuses.EfficientEmeralds.Cost[efficientEmerald])
		usage.Crop += int(upgrades.Bonuses.TowerAura.Cost[aura] + upgrades.Bonuses.EmeraldsRate.Cost[emeraldRate])
		usage.Wood += int(upgrades.Bonuses.StrongerMinions.Cost[towerMultiAttacks] + upgrades.Bonuses.LargerEmeraldsStorage.Cost[emStorage] + upgrades.Bonuses.StrongerMinions.Cost[strongerMinions])
		usage.Fish += int(upgrades.Bonuses.TowerMultiAttacks.Cost[efficientResource])
		usage.Emerald += int(upgrades.Bonuses.LargerResourceStorage.Cost[resStorage] + upgrades.Bonuses.ResourceRate.Cost[resourceRate] + upgrades.Bonuses.EfficientResource.Cost[efficientResource])

		// set the usage cost of the territory
		(*territory).TerritoryUsage = usage
	}
}

func resourceTick(territories *map[string]*table.Territory) {
	// move all resources and emeralds onto next terr's transversing to hq
	// if the terr is hq then ignore
	for _, territory := range *territories {
		if territory.Property.HQ {
			continue
		}

		// or if the terr is not claimed then do not move the storage into transversing res
		// but if the unclaimed terr has transversing res then move it to the next terr's transversing res anyways

		// we need to move res and emse from terrs closest to hq first using territory.RouteFromHQ

	}
}

func getPathToHQCheapest(territories *map[string]*table.Territory, HQ string) {
	// name is the HQ territory name
	// get path to hq using dijkstra, depending on the trading style
	// fastest  int
	// will find the shortest path while ccheapest will find the shortest path with the least GLOBAL tax
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
		territory.RouteToHQ = path
		var reversePath = make([]string, len(path))
		copy(reversePath, path)
		arrutil.Reverse(reversePath)
		territory.RouteFromHQ = reversePath
	}
}

func getPathToHQFastest(t *map[string]*table.Territory, HQ string) {

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
		routeToHQ := territory.RouteToHQ

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
