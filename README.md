# iroiroiru

ever wonder who’s living quietly beside you?
tiny mushrooms tucked beneath leaves,
wildflowers whispering in the breeze,
beetles, birds, and beasts passing through unseen paths...

with free biodiversity data from gbif,
iroiroiru helps you uncover the secret lives all around you,
right here, right now.

## instances

[iroiroiru.jp - serves up biodiversity data from japan](https://iroiroiru.jp/)

## requirements

* go 1.24.2
* mongodb 8

## architecture

iroiroiru consists of the following:

* iroiroserve - a go api microservice
* iroiroload - command line tool to ingest occurrence and multimedia data
* iroiroview - minimalist front-end allowing you to find nearby occurrences

## todo

* improve this readme 🙃
* handle thumbnails from non-inaturalist images (first step: count images)
* show vernacular names when available
* allow linking directly to occurrences
* grab all inaturalist images (even nd is ok - inaturalist made the thumbnails)
* all medium thumbs, about ~20-25gb?

## might be nice to have

* improve icons
* iroiroload: bulk insert...?
* deal with overlapping markers...
