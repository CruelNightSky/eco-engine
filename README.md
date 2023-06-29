# eco-engine

Eco engine is a dynamic simulator designed to simulate Wynncraft guild economy in real-time. Unlike other eco simulators out there, eco engine accurately replicates the guild's ecosystem, including the real-time travel of resources. Upgrading territories or boosting resource production will result in a delay before the resources arrive, just like in the game.

The project is written in Go (Golang), providing a powerful and efficient implementation for simulating the guild economy.

## Features

- Real-time simulation of guild economy
- Dynamic resource travel system
- Territory upgrades and bonuses
- Load and unload any territory on the fly
- Fast forward feature (Not yet implemented)
- Embeddable using FFI (Not yet implemented)
- Or communicate over WebSocket or HTTP

## Usage
Once eco engine is running, you can interact with it through WebSocket or HTTP POST/GET requests to simulate the guild economy. Refer to the documentation or examples provided in the project's repository for detailed instructions on how to communicate with the engine and utilize its features.

## Example Usage

To initialize territories:

**Endpoint:** `POST http://localhost:$PORT/init`

**Request Body:**
```json
{
    "territories" : [ "Mine Base Plains", "Abandoned Pass", "Plains Lake", "Mining Base Lower", "Mining Base Upper", "Ternaves Plains Lower", "Detlas Savannah Transition", "The Silent Road" ],
    "hq" : "Mine Base Plains"
}
```

To set territory upgrades and bonuses:

**Endpoint:** `POST http://localhost:$PORT/manage/Mine%20Base%20Plains/set`

**Request Body:**
```json
{
    "upgrades" : {
        "damage" : 6,
        "attack" : 6,
        "health" : 6,
        "defence" : 6
    },
    "bonuses" : {
        "minionsDamage" : 3,
        "multi" : 1,
        "aura" : 3,
        "volley" : 3,
        "largerEmeraldStorage" : 6,
        "LargerResourceStorage" : 6,
        "efficientResource" : 3,
        "resourceRate" : 3
    },
    "style" : "cheapest",
    "tax" : {
        "ally" : null,
        "global" : 60
    }
}
```

Please make sure to replace `$PORT` in the endpoint URLs with the actual port number you're using for the eco engine server.



## Installation

To install and run eco engine, follow these steps:

Make sure you have Golang >= 1.18 installed

1. Clone the repository:

```shell
git clone https://github.com/CruelNightSky/eco-engine.git
```

2. Change to the project directory:

```shell
cd eco-engine
```

3. Build the Go executable

```shell
go build
```

4. Run the executable

```shell
./main
```

By default, eco engine will start a web server and listen on a port specified in argv for WebSocket and HTTP communication.


## Documentations
Coming soon™️

## Contributing
Contributions to Eco engine are welcome! If you encounter any issues, have suggestions, or would like to contribute new features, please feel free to open a GitHub issue or submit a pull request.

## License
This project is licensed under the [Affero General Public License](https://en.wikipedia.org/wiki/Affero_General_Public_License).
