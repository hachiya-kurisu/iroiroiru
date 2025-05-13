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

* a few scripts for ingesting data into mongodb
* a go api microservice
* a very simple front-end allowing you to find occurrences close to where you're at

## todo

* improve this readme 🙃
* show occurrence images (?) or images based on species (?)
* improve icons
* find more suitable free map tiles?
* ingest data directly from a dwca
* drop data with uncertainty > 250 meters? 100m? 500 feels too... something 🤔
