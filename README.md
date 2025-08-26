# fangcli
A CLI interface to fangbreaker.zone for Everquest. This application is not associated with the site nor with the video game Everquest.

This is a work in progress while I learn Go.

### Arguments:
 - bonus - Limit the output to just the bonus listed (experience, AA, coin, loot, rare, skill, respawn, faction)
 - minlevel - Show zones >= to this level
 - maxlevel - Show zones <= to this level
 - sortbylevel - Sort direction. ASC by default. (ASC or DESC)
 - expansion - Limit to the expansion name (classic, velious etc)
 - quiet - Hide extra output

### Examples

#### Basic running
fangcli

#### Only showing experience bonus zones
fangcli -bonus experience

#### Show zones at least level 10
fangcli -minlevel 10

#### Show rare bonus zones between the levels 10 and 30
fangcli -minlevel 10 -maxlevel 30 -bonus rare
