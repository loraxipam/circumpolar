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

# -kilo
echo kilo
go run circumpolar.go -kilo $TIKI \
	$FIJI	$MIA	$LAX	$NYC

# -mile
echo mile
go run circumpolar.go -mile $TIKI \
	$FIJI	$MIA	$LAX	$NYC

# Mars
echo mars
go run circumpolar.go -kilo -radius 3390 $TIKI \
	$FIJI	$MIA	$LAX	$NYC

# -json
echo json
go run circumpolar.go -json $TIKI \
	$FIJI	$MIA	$LAX	$NYC
echo ""

# Mars -json
echo mars json
go run circumpolar.go -json -kilo -radius 3390 $TIKI \
	$FIJI	$MIA	$LAX	$NYC
echo ""
