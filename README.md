# circumpolar

   The circumpolar tool is a command line tool to provide distance measurements between a point
   and a set of points, on a sphere of a given radius.

### Origin

   This was inspired by the "direction pole" at Jerry's Tiki Restaurant in Ponce Inlet. It is one of those poles with arrows and cities and miles on them, each pointing in the general direction that was supposedly correct for each arrow. Now, you too, can recreate your very own custom "direction pole" but yours will have **pinpoint accuracy**. Aloha, baby.

## Usage

   Canonical use is ```circumpolar {latA lonA} {latX lonX} [{latY lonY}...]``` where latA/lonA are the starting
   location point values and all subsequent pairs get reported relative to the first pair.
   
   Typical use is as follows: 

   ```circumpolar 40.75 -73.9 51.45 1.15``` would produce the distance from New York City to
   Oxford University in nautical miles here on Earth.

   ```circumpolar -- -1.28 36.82 51.45 -1.15``` would produce the distance from Nairobi to Oxford.
   You must escape an initial negative latitude value with ```--```.

   ```circumpolar -kilo -radius 3390 -- -1.28 36.82 51.45 -1.15``` would produce the same Nairobi/Oxford results
   in kilometers, though, if those cities were on Mars, whose radius is 3390 kilometers.

   ```circumpolar -json -- -1.28 36.82 51.45 -1.15``` would produce the same Nairobi/Oxford results
   but would return JSON instead of rows.

   NOAA provides magnetic declination queries for given points, so now results also show compass position with declination included.

## Options

   * ```-json``` - output results in JSON format
   * ```-kilo``` - output result distances in kilometers
   * ```-mile``` - output result distances in statute miles
   * ```-home``` - Stay home. Do not query NOAA for magnetic declination.
   * ```-radius N``` - use N as the sphere's radius rather than the Earth's

   If no distance flags are provided, the default distance unit is the nautical mile and calculations use the Earth's radius.

   Notice that the ```-radius``` flag doesn't really affect the kilo/mile flags. You can pass just a radius and the output will be _labeled_ as the default "NM" but that's really moot when the radius is overridden -- all output is relative to the value that was passed in. If you want distances reported in _The Register_ Vulture Central-approved standard ```linguine```s, then just pass the Earth's radius as 45,506,300 and only you will know that the "NM" label really means "lg" and not nautical miles.

   The ```-kilo``` and ```-mile``` flags _do_ affect the radius **when used alone**, though, because they will assign the radius to the Earth's value for the chosen unit. But understand that when you use a ```-radius``` value with these, you'll simply get the "km" and "mi" labels and calculations will use the passed-in radius.

## Next

   - ~~Output JSON~~
   - Parse d.m.s format and NSEW designators
   - Provide distance to lat/lon lines
   - ~~Magnetic declination~~
