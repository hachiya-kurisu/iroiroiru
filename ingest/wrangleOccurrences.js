const db = connect("mongodb://localhost:27017/iroiro");

print("get rid of some broken data (gbif should already have done this)");

db.occurrences.deleteMany({ decimalLatitude: "" })
db.occurrences.deleteMany({ decimalLongitude: "" })

print("let's convert lat/lon to geojson...");

const cursor = db.occurrences.find({})
cursor.forEach(function(doc) {
  db.occurrences.updateOne(
    { _id: doc._id },
    {
      $set: {
        location: {
          type: "Point",
          coordinates: [doc.decimalLongitude, doc.decimalLatitude]
        }
      }
    }
  );
});

print("then we add a 2dsphere index");

db.occurrences.createIndex({ location: "2dsphere" });

print("we want to have a simple count by kingdom, so let's index that as well");

db.occurrences.createIndex({ kingdom: 1});

print("not sure if that index really helps, but whatever");
print("let's rock");

