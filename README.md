## Chartprobe
--

![build](https://github.com/ezbuy/chartprobe/actions/workflows/build-image.yaml/badge.svg)

Chartprobe is a command-line app for [chartmuseum](github.com/helm/chartmuseum),which provides some common request templates .

### Using Config File

```shell
~ > cat museum.yaml

# specify your host
CHARTPROBE_HOST: http://YOUR_CHARTMUSEUM_HOST/YOUR_REPO

```

### In the box

* Get
    * Get by chartname prefix: `chartprobe get --prefix your_chart_prefix`
    * Get all: `chartprobe get -a`
* Delete
    * Delete by chartname prefix: `chartprobe -c museum.yaml delele --prefix your_chart_prefix`
    * Delete all: `chartprobe -c museum.yaml delete -a`
    * Delete charts spawned during 24h:  `chartprobe -c museum.yaml delete --period 24h`

> Tap chartprobe -h to seek more help .

### Using as a Docker Container to swipe museum

```shell
~ > docker pull ghcr.io/ezbuy/chartprobe:latest
~ > docker run  -name chartprobe -e CHARTPROBE_HOST="your_museum_host" -e CHARTPROBE_PERIOD="-168h" ghcr.io/ezbuy/chartprobe:latest delete -a
```
