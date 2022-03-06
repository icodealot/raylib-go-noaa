## Raylib-NOAA Weather Client
> The source data for this application is maintained by the National Weather Service under the umbrella of the 
National Oceanic and Atmospheric Administration (NOAA).

This is a hobby project and a simple graphical client for the NOAA weather API hosted at weather.gov. This software 
uses the (excellent) Raylib library with Go bindings to display weather forecast data. The forecast data is rendered 
to an underlying OpenGL context setup by Raylib.

Go was chosen for its ease of use and robust HTTP support (but also because I already had a NOAA API wrapper handy.)

### Requirements
This project depends on the following external or open source components:
- api.weather.gov (for weather observations)
- [raylib](https://www.raylib.com/), a simple and easy-to-use library to enjoy videogames programming. (For Window 
creation and ease of rendering)
- Golang bindings for [raylib-go](https://github.com/gen2brain/raylib-go)
- C/C++ compiler (for `cgo`)
- ...TBD other for parsing YAML, etc.

Check out raylib and the raylib-go bindings projects for prerequisites and initial setup.

### Configuration
There are two different ways to configure the program to pull in a weather forecast and both options require a 
latitude and longitude value.

##### Config File
`config.yml`: The first option is to edit the `config.yml` file included with this project. The example shown here 
points the application at the Chicago, IL weather observations.

```yaml
noaa:
  latitude: 41.837
  longitude: -87.685
```
Then, once you have setup your coordinates you can pass a configuration option to the software such as:
```bash
noaawc.exe -config config.yml
```

##### Command Line Arguments
`-lat <latitude> -lon <longitude>`: The second option is to pass in the `lat` and `lon` arguments when starting the program.
```bash
noaawc.exe -lat "41.837" -lon "-87.685"
```
