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
	"net/http"
	"os"
	"time"

	"github.com/RyanCarrier/dijkstra"
)

var (
	t                 map[string]*table.Territory // loaded
	ut                map[string]*table.Territory // unloaded
	st                int                         // second tick
	loadedTerritories = make(map[string]*table.Territory)
)

// setInterval function to use later for resgen
func setInterval(f func(map[string]*table.Territory), milliseconds int) chan bool {
	ticker := time.NewTicker(time.Duration(milliseconds) * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				// run func to work on the data
				f(t)
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
	return done
}

func clearInterval(done chan bool) {
	done <- true
}

func setTimeout(f func(args ...interface{}), ms int, args ...interface{}) chan bool {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	done := make(chan bool)
	go func() {
		select {
		case <-ticker.C:
			f(args...)
		case <-done:
			ticker.Stop()
			return
		}
	}()
	return done
}

func init() {
	// load all upgrades data
	var bytes, err = os.ReadFile("./upgrades.json")
	if err != nil {
		panic(err)
	}

	var upgrades table.CostTable

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
				Capacity: table.TerritoryResource{
					Emerald: 3000,
					Ore:     300,
					Wood:    300,
					Fish:    300,
					Crop:    300,
				},
				Current: table.TerritoryResource{
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

		log.Println("initialised territories")
		defer CalculateRouteToHQTax(&t, hq)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"message":"initialised"}`))
	})

	http.ListenAndServe(":"+port, nil)

	// run generateResource every 1s and resTick every 60s
	var done = make(chan struct{})
	var secTicker = time.NewTicker(time.Second)
	var minTicker = time.NewTicker(time.Minute)
	go func(t map[string]*table.Territory) {
		for {
			select {
			case <-secTicker.C:
				generateResorce(t)
				log.Println("tick")
			case <-done:
				secTicker.Stop()
				minTicker.Stop()
				return
			}
		}
	}(t)

	fmt.Printf("test")
}

func generateResorce(territories map[string]*table.Territory) {
	// rate means how many seconds it takes to generate n resource
	// n resource is calculated like this
	// nres = base res prod * efficient resource
	// so if rate is level 0 then it takes 1 second to generate 1/4 of the resource
	// and if resource stored in the storage excees the capacity then the excess resource will be lost
	// stoarge capacity is calculated like this
	// cap = base cap * larger storage

	// emerald generation
	for name, territory := range territories {

		// calculate the resource production
		var emeraldRate = float64(territory.Property.Bonuses.EmeraldRate)
		var emeraldProduction = float64(territory.ResourceProduction.Emerald) * (1 + emeraldRate/100)
		var emeraldStorage = float64(territory.Storage.Capacity.Emerald) * (1 + float64(territory.Property.Bonuses.LargerEmeraldStorage)/100)

		// if the storage is full then do nothing
		if float64(territory.Storage.Current.Emerald) >= emeraldStorage {
			// send websocket message "A territory %s is producing more emeralds than it can store!"
			continue
		}

		// if the storage is not full then generate the resource

		if float64(territory.Storage.Current.Emerald)+emeraldProduction <= emeraldStorage && !territory.Property.HQ {
			currEms := float64((territories)[name].Storage.Current.Emerald)
			currEms += emeraldProduction
		} else if !territory.Property.HQ {
			currEms := float64((territories)[name].Storage.Current.Emerald)
			currEms += emeraldStorage
		} else if float64(territory.Storage.Current.Emerald)+emeraldProduction <= emeraldStorage {

			// normal terrs
			currEms := float64((territories)[name].Storage.Current.Emerald)
			currEms += emeraldProduction
		} else {
			currEms := float64((territories)[name].Storage.Current.Emerald)
			currEms += emeraldStorage
		}
		fmt.Println(territories)
	}

	if st < 60 {
		st++
	} else {
		for range territories {
			resourceTick(territories)
		}
		st = 0
	}
}

func resourceTick(territories map[string]*table.Territory) {

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
		log.Println(territory.Name, " ", path, " ", territory.RouteToHQ)
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
		log.Println("Territory: ", territory)
	}
}

func CalculateRouteToHQTax(territories *map[string]*table.Territory, from string) float64 {
	// the formular to calculate tax are as follows
	// 1 - ((1 - terr1Tax) * (1 - terr2Tax) * (1 - terr3Tax) * ... * (1 - terrnTax))
	// for example if there are 4 territories and the tax are as follows : 60 60 5 5
	// the tax will be 1 - (0.40 * 0.40 * 0.95 * 0.95) = 0.8556 or 85.56%
	var startingTerritory = (*territories)[from]
	log.Println("Starting territory: ", startingTerritory.RouteToHQ)
	var routeToHQ = startingTerritory.RouteToHQ

	var taxList = []float64{}
	log.Println(routeToHQ)
	// iterate through the route to hq and get the tax of each territory
	for _, territory := range routeToHQ {
		log.Println("123called : ", territory)
		if (*territories)[territory].Property.HQ {
			continue
		} else {
			taxList = append(taxList, 1-(float64((*territories)[territory].Property.Tax.Others)/100))
			log.Println("test", 1-(float64((*territories)[territory].Property.Tax.Others)/100))
		}
	}
	log.Println(taxList)

	// calculate the tax
	for _, tax := range taxList {
		startingTerritory.RouteTax *= tax
	}

	startingTerritory.RouteTax *= 100
	log.Println("Route tax from ", from, " to HQ is ", startingTerritory.RouteTax)

	return startingTerritory.RouteTax
}
