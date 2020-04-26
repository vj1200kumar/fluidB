<p align="center">
  <a href="https://fluentbase.org"><img src="/img/Logo.png" alt="fluiDB"></a>
</p>
<p align="center">
<a href="https://tile38.com/slack"><img src="https://img.shields.io/badge/slack-channel-orange.svg" alt="Slack Channel"></a>
<a href="https://travis-ci.org/tidwall/tile38"><img src="https://travis-ci.org/tidwall/tile38.svg?branch=master" alt="Build Status"></a>
<a href="https://hub.docker.com/r/tile38/tile38"><img src="https://img.shields.io/badge/docker-ready-blue.svg" alt="Docker Ready"></a>
</p>

Fluentbase is an open source (MIT licensed), in-memory geolocation data store, spatial index, and realtime geofence, distributed uder MIT LICENSE using code of [Tile-38](https://github.com/tidwall/tile38), [Algernon](https://github.com/xyproto/algernon), which also distributed uder MIT LICENSE.  It supports a variety of object types including lat/lon points, bounding boxes, XYZ tiles, Geohashes, and GeoJSON. 

<p align="center">
<i>This README is quick start document. You can find detailed documentation at <a href="https://fluentbase.org">https://fluentbase.org</a>.</i><br><br>
<a href="#searching"><img src="/img/search-nearby.png" alt="Nearby" border="0" width="120" height="120"></a>
<a href="#searching"><img src="/img/search-within.png" alt="Within" border="0" width="120" height="120"></a>
<a href="#searching"><img src="/img/search-intersects.png" alt="Intersects" border="0" width="120" height="120"></a>
<a href="https://fluentbase.org/topics/roaming-geofences"><img src="/img/roaming.gif" alt="Roaming Geofences" border="0" width="120" height="120"></a>
</p>

## Features

- Spatial index with [search](#searching) methods such as Nearby, Within, and Intersects.
- Realtime [geofencing](#geofencing) through [webhooks](https://fluentbase.org/commands/sethook) or [pub/sub channels](#pubsub-channels).
- Object types of [lat/lon](#latlon-point), [bbox](#bounding-box), [Geohash](#geohash), [GeoJSON](#geojson), [QuadKey](#quadkey), and [XYZ tile](#xyz-tile).
- Support for lots of [Clients Libraries](#client-libraries) written in many different languages.
- Variety of protocols, including [http](#http) (curl), [websockets](#websockets), [telnet](#telnet), and the [Redis RESP](http://redis.io/topics/protocol).
- Server responses are [RESP](http://redis.io/topics/protocol) or [JSON](http://www.json.org).
- Full [command line interface](#cli).
- Leader / follower [replication](#replication).
- In-memory database that persists on disk.

## Components
- `fluentbase-server    ` - The server
- `fluentbase-cli       ` - Command line interface tool

## Getting Started

### Getting fluentbase

Perhaps the easiest way to get the latest fluentbase is to use one of the pre-built release binaries which are available for OSX, Linux, FreeBSD, and Windows. Instructions for using these binaries are on the GitHub [releases page](https://github.com/tidwall/fluentbase/releases).


### Building fluentbase 

fluentbase can be compiled and used on Linux, OSX, Windows, FreeBSD, and probably others since the codebase is 100% Go. We support both 32 bit and 64 bit systems. [Go](https://golang.org/dl/) must be installed on the build machine.

To build everything simply:
```
$ make
```

To test:
```
$ make test
```

### Running 
For command line options invoke:
```
$ ./fluentbase-server -h
```

To run a single server:

```
$ ./fluentbase-server

# The fluentbase shell connects to localhost:9470
$ ./fluentbase-cli
> help
```

## Starts working

Basic operations:
```
$ ./fluentbase-cli

# add a couple of points named 'truck1' and 'truck2' to a collection named 'fleet'.
> set fleet truck1 point 33.5123 -112.2693   # on the Loop 101 in Phoenix
> set fleet truck2 point 33.4626 -112.1695   # on the I-10 in Phoenix

# search the 'fleet' collection.
> scan fleet                                 # returns both trucks in 'fleet'
> nearby fleet point 33.462 -112.268 6000    # search 6 kilometers around a point. returns one truck.

# key value operations
> get fleet truck1                           # returns 'truck1'
> del fleet truck2                           # deletes 'truck2'
> drop fleet                                 # removes all 
```

The full list of commands you can found on [full command list](https://fluentbase.org/%D0%BA%D0%BE%D0%BC%D0%B0%D0%BD%D0%B4%D1%8B/).

## Fields
Fields are extra data that belongs to an object. A field is always a double precision floating point. There is no limit to the number of fields that an object can have. 

To set a field when setting an object:
```
> set fleet truck1 field speed 90 point 33.5123 -112.2693             
> set fleet truck1 field speed 90 field age 21 point 33.5123 -112.2693
```

To set a field when an object already exists:
```
> fset fleet truck1 speed 90
```

## Searching

fluentbase has support to search for objects and points that are within or intersects other objects. All object types can be searched including Polygons, MultiPolygons, GeometryCollections, etc.

<img src="/img/search-within.png" width="200" height="200" border="0" alt="Search Within" align="left">

#### Within 
WITHIN searches a collection for objects that are fully contained inside a specified bounding area.
<BR CLEAR="ALL">

<img src="/img/search-intersects.png" width="200" height="200" border="0" alt="Search Intersects" align="left">

#### Intersects
INTERSECTS searches a collection for objects that intersect a specified bounding area.
<BR CLEAR="ALL">

<img src="/img/search-nearby.png" width="200" height="200" border="0" alt="Search Nearby" align="left">

#### Nearby
NEARBY searches a collection for objects that intersect a specified radius.
<BR CLEAR="ALL">

### Search options
**SPARSE** - This option will distribute the results of a search evenly across the requested area.  
This is very helpful for example; when you have many (perhaps millions) of objects and do not want them all clustered together on a map. Sparse will limit the number of objects returned and provide them evenly distributed so that your map looks clean.<br><br>
You can choose a value between 1 and 8. The value 1 will result in no more than 4 items. The value 8 will result in no more than 65536. *1=4, 2=16, 3=64, 4=256, 5=1024, 6=4098, 7=16384, 8=65536.*<br><br>
<table>
<td>No Sparsing<img src="/img/sparse-none.png" width="100" height="100" border="0" alt="Search Within"></td>
<td>Sparse 1<img src="/img/sparse-1.png" width="100" height="100" border="0" alt="Search Within"></td>
<td>Sparse 2<img src="/img/sparse-2.png" width="100" height="100" border="0" alt="Search Within"></td>
<td>Sparse 3<img src="/img/sparse-3.png" width="100" height="100" border="0" alt="Search Within"></td>
<td>Sparse 4<img src="/img/sparse-4.png" width="100" height="100" border="0" alt="Search Within"></td>
<td>Sparse 5<img src="/img/sparse-5.png" width="100" height="100" border="0" alt="Search Within"></td>
</table>

*Please note that the higher the sparse value, the slower the performance. Also, LIMIT and CURSOR are not available when using SPARSE.* 

**WHERE** - This option allows for filtering out results based on [field](#fields) values. For example<br>```nearby fleet where speed 70 +inf point 33.462 -112.268 6000``` will return only the objects in the 'fleet' collection that are within the 6 km radius **and** have a field named `speed` that is greater than `70`. <br><br>Multiple WHEREs are concatenated as **and** clauses. ```WHERE speed 70 +inf WHERE age -inf 24``` would be interpreted as *speed is over 70 <b>and</b> age is less than 24.*<br><br>The default value for a field is always `0`. Thus if you do a WHERE on the field `speed` and an object does not have that field set, the server will pretend that the object does and that the value is Zero.

**MATCH** - MATCH is similar to WHERE except that it works on the object id instead of fields.<br>```nearby fleet match truck* point 33.462 -112.268 6000``` will return only the objects in the 'fleet' collection that are within the 6 km radius **and** have an object id that starts with `truck`. There can be multiple MATCH options in a single search. The MATCH value is a simple [glob pattern](https://en.wikipedia.org/wiki/Glob_(programming)).

**CURSOR** - CURSOR is used to iterate though many objects from the search results. An iteration begins when the CURSOR is set to Zero or not included with the request, and completes when the cursor returned by the server is Zero.

**NOFIELDS** - NOFIELDS tells the server that you do not want field values returned with the search results.

**LIMIT** - LIMIT can be used to limit the number of objects returned for a single search request.


## Geofencing

<img src="/img/geofence.gif" width="200" height="200" border="0" alt="Geofence animation" align="left">
A <a href="https://en.wikipedia.org/wiki/Geo-fence">geofence</a> is a virtual boundary that can detect when an object enters or exits the area. This boundary can be a radius, bounding box, or a polygon. fluentbase can turn any standard search into a geofence monitor by adding the FENCE keyword to the search. 

*fluentbase also allows for [Webhooks](http://fluentbase.com/commands/sethook) to be assigned to Geofences.*

<br clear="all">

A simple example:
```
> nearby fleet fence point 33.462 -112.268 6000
```
This command opens a geofence that monitors the 'fleet' collection. The server will respond with:
```
{"ok":true,"live":true}
```
And the connection will be kept open. If any object enters or exits the 6 km radius around `33.462,-112.268` the server will respond in realtime with a message such as:

```
{"command":"set","detect":"enter","id":"truck02","object":{"type":"Point","coordinates":[-112.2695,33.4626]}}
```

The server will notify the client if the `command` is `del | set | drop`. 

- `del` notifies the client that an object has been deleted from the collection that is being fenced.
- `drop` notifies the client that the entire collection is dropped.
- `set` notifies the client that an object has been added or updated, and when it's position is detected by the fence.

The `detect` may be one of the following values.

- `inside` is when an object is inside the specified area.
- `outside` is when an object is outside the specified area.
- `enter` is when an object that **was not** previously in the fence has entered the area.
- `exit` is when an object that **was** previously in the fence has exited the area.
- `cross` is when an object that **was not** previously in the fence has entered **and** exited the area.

These can be used when establishing a geofence, to pre-filter responses. For instance, to limit responses to `enter` and `exit` detections:

```
> nearby fleet fence detect enter,exit point 33.462 -112.268 6000
```

### Pub/sub channels

fluentbase supports delivering geofence notications over pub/sub channels. 

To create a static geofence that sends notifications when a bus is within 200 meters of a point and sends to the `busstop` channel:

```
> setchan busstop nearby buses fence point 33.5123 -112.2693 200
```

Subscribe on the `busstop` channel:

```
> subscribe busstop
```

## Object types

All object types except for XYZ Tiles and QuadKeys can be stored in a collection. XYZ Tiles and QuadKeys are reserved for the SEARCH keyword only.

#### Lat/lon point
The most basic object type is a point that is composed of a latitude and a longitude. There is an optional `z` member that may be used for auxiliary data such as elevation or a timestamp.
```
set fleet truck1 point 33.5123 -112.2693     # plain lat/lon
set fleet truck1 point 33.5123 -112.2693 225 # lat/lon with z member
```

#### Bounding box
A bounding box consists of two points. The first being the southwestern most point and the second is the northeastern most point.
```
set fleet truck1 bounds 30 -110 40 -100
```
#### Geohash
A [geohash](https://en.wikipedia.org/wiki/Geohash) is a string respresentation of a point. With the length of the string indicating the precision of the point. 
```
set fleet truck1 hash 9tbnthxzr # this would be equivlent to 'point 33.5123 -112.2693'
```

#### GeoJSON
[GeoJSON](http://geojson.org/) is an industry standard format for representing a variety of object types including a point, multipoint, linestring, multilinestring, polygon, multipolygon, geometrycollection, feature, and featurecollection.

<i>* All ignored members will not persist.</i>

**Important to note that all coordinates are in Longitude, Latitude order.**

```
set city tempe object {"type":"Polygon","coordinates":[[[0,0],[10,10],[10,0],[0,0]]]}
```

#### XYZ Tile
An XYZ tile is rectangle bounding area on earth that is represented by an X, Y coordinate and a Z (zoom) level.
Check out [maptiler.org](http://www.maptiler.org/google-maps-coordinates-tile-bounds-projection/) for an interactive example.

#### QuadKey
A QuadKey used the same coordinate system as an XYZ tile except that the string representation is a string characters composed of 0, 1, 2, or 3. For a detailed explanation checkout [The Bing Maps Tile System](https://msdn.microsoft.com/en-us/library/bb259689.aspx).


## Network protocols

It's recommended to use a [client library](#client-libraries) or the [fluentbase CLI](#running), but there are times when only HTTP is available or when you need to test from a remote terminal. In those cases we provide an HTTP and telnet options.

#### HTTP
One of the simplest ways to call a fluentbase command is to use HTTP. From the command line you can use [curl](https://curl.haxx.se/). For example:

```
# call with request in the body
curl --data "set fleet truck3 point 33.4762 -112.10923" localhost:9470

# call with request in the url path
curl localhost:9470/set+fleet+truck3+point+33.4762+-112.10923
```

#### Websockets
Websockets can be used when you need to Geofence and keep the connection alive. It works just like the HTTP example above, with the exception that the connection stays alive and the data is sent from the server as text websocket messages.

#### Telnet
There is the option to use a plain telnet connection. The default output through telnet is [RESP](http://redis.io/topics/protocol).

```
telnet localhost 9470
set fleet truck3 point 33.4762 -112.10923
+OK

```

The server will respond in [JSON](http://json.org) or [RESP](http://redis.io/topics/protocol) depending on which protocol is used when initiating the first command.

- HTTP and Websockets use JSON. 
- Telnet and RESP clients use RESP.

## Client Libraries in most popular languages

fluentbase uses the [Redis RESP](http://redis.io/topics/protocol) protocol natively. Therefore most clients that support basic Redis commands will in turn support fluentbase. Below are a few of the popular clients. 

- C: [hiredis](https://github.com/redis/hiredis)
- C#: [StackExchange.Redis](https://github.com/StackExchange/StackExchange.Redis)
- C++: [redox](https://github.com/hmartiro/redox)
- Go: [go-redis](https://github.com/go-redis/redis) https://github.com/tidwall/fluentbase/wiki/Go-example-(go-redis)
- Go: [redigo](https://github.com/gomodule/redigo) (https://github.com/tidwall/fluentbase/wiki/Go-example-(redigo)
- Node.js: [node-fluentbase](https://github.com/phulst/node-fluentbase) (https://github.com/tidwall/fluentbase/wiki/Node.js-example-(node-fluentbase)
- Node.js: [node_redis](https://github.com/NodeRedis/node_redis) (https://github.com/tidwall/fluentbase/wiki/Node.js-example-(node-redis)
- PHP: [tinyredisclient](https://github.com/ptrofimov/tinyredisclient) (https://github.com/tidwall/fluentbase/wiki/PHP-example-(tinyredisclient)
- PHP: [phpredis](https://github.com/phpredis/phpredis)
- Python: [redis-py](https://github.com/andymccurdy/redis-py) (https://github.com/tidwall/fluentbase/wiki/Python-example))
- Scala: [scala-redis](https://github.com/debasishg/scala-redis)
- Swift: [Redbird](https://github.com/czechboy0/Redbird)

## Contact

Grigoriy Safronov <gvsafronov@gmail.com>


## Gratitude

Josh Baker, thank you for sharing your experience and knowledge with me.


## License

fluentbase source code is available under the MIT [License](/LICENSE).
