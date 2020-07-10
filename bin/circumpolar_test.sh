#!/bin/bash --noprofile
# Test calls to circumpolar
#
# Cities used on Larry's Tiki sign
# Larry's - 29.1N, 80.9W
# Suva - 18.14S, 178.44E
# Miami - 25:46N, 80:12W
# LA  - 34:03N,  118:15W
# NYC - 40:46N, 73:59W

TIKI="29.10 -80.93"
FIJI="-18.14 178.44"
MIA="25.77 -80.2"
LAX="34.05 -118.25"
NYC="40.75 -73.9"
BAD="140.75 -873.9"

CMD='circumpolar'
[[ -f circumpolar.go ]] && CMD='go run circumpolar.go'

# none
$CMD

# default
echo default
$CMD $TIKI \
	$FIJI	$MIA	$LAX	$NYC	$BAD

# -mile
echo mile
$CMD -mile $TIKI \
	$FIJI	$MIA	$LAX	$NYC	$BAD

# -kilo
echo kilo
$CMD -kilo $TIKI \
	$FIJI	$MIA	$LAX	$NYC	$BAD

# Mars
echo mars
$CMD -kilo -radius 3390 $TIKI \
	$FIJI	$MIA	$LAX	$NYC	$BAD

# Mars -json
echo mars json
$CMD -json -kilo -radius 3390 $TIKI \
	$FIJI	$MIA	$LAX	$NYC	$BAD
echo ""

