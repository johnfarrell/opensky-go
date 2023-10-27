# Getting Started

## Installing

To install opensky-go, run

```Bash
go get -u github.com/johnfarrell/opensky-go
```

## Basic Usage

The following code snippet shows a simple usage to fetch the flights around the Seattle, WA area.

```Go
func main() {
   osClient, err := opensky.NewClient()
   if err != nil {
      panic(err)
   }
   
   states, err := osClient.GetStates(&opensky.BoundingBox{
       LongitudeMin: -122.4845,
       LongitudeMax: -122.1693,
       LatitudeMin:  47.5014,
       LatitudeMax:  47.6954,
   })
   if err != nil {
      panic(err)
   }
   
   fmt.Println(len(states))
}
```

## Passing Client Options

To configure the Client, pass the client options to `NewClient` in the following manner:

```Go
client, _ := opensky.NewClient(opensky.WithTimeout(60 * time.Second))
```

Multiple options can be passed at the same time:
```Go
client, _ := opensky.NewClient(
    opensky.WithAuthorization("username", "password"),
    opensky.WithTimeout(60 * time.Second),
)
```