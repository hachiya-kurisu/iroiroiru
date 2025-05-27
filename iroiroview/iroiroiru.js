const key = '<protomaps api key>';
const map = new maplibregl.Map({
  container: 'map',
  style: 'https://api.protomaps.com/styles/v5/light/ja.json?key=' + key,
  center: [139.805, 35.749],
  zoom: 14,
});

const icons = {
  'Animalia': '/images/animalia.png',
  'Chromista': '/images/chromista.png',
  'Plantae': '/images/plantae.png',
  'Fungi': '/images/fungi.png',
  'Fallback': '/images/fallback.png',
}

map.on('load', () => {
  Promise.all(
    Object.entries(icons).map(([name, url]) =>
      new Promise((resolve, reject) => {
        const img = new Image();
        img.src = url;
        img.onload = function(e) {
          map.addImage(url, e.target);
          resolve();
        }
      })
    )
  ).then(() => {
    map.addSource('occurrences', {
      type: 'geojson',
      data: {
        type: 'FeatureCollection',
        features: [],
      }
    });
    map.addLayer({
      id: 'occurrenceLayer',
      type: 'symbol',
      source: 'occurrences',
      layout: {
        'icon-image': ['get', 'icon'],
        'icon-size': 0.5,
      }
    });
    geolocate.trigger();
  }).catch(err => console.error(err))
});

map.on('click', 'occurrenceLayer', function(e) {
  const feature = e.features[0];
  const o = JSON.parse(feature.properties.raw);
  const template = document.getElementById('occurrenceTemplate');
  const clone = template.content.cloneNode(true);
  const setText = (sel, val) => {
    const el = clone.querySelector(sel);
    val ? el.textContent = val : el.parentElement.remove();
  };

  setText(".coordinates", o.location.coordinates.reverse());
  setText(".occurrenceRemarks", o.occurrenceRemarks);
  setText(".locationRemarks", o.locationRemarks);
  setText(".scientificName", o.acceptedScientificName);
  setText(".locality", o.locality || o.verbatimLocality);
  setText(".recordedBy", o.recordedBy);
  setText(".identifiedBy", o.identifiedBy);
  setText(".eventDate", o.eventDate);
  setText(".license", o.license);
  setText(".lifestage", o.lifestage);
  setText(".datasetName", o.datasetName);
  setText(".occurrenceID", o.occurrenceID);
  setText(".uncertaintyInMeters", o.coordinateUncertaintyInMeters);
  setText(".distance", o.dist.calculated.toFixed(0) + "メートル先");

  const multimedia = clone.querySelector(".multimedia")
  if(o.multimedia) {
    o.multimedia.forEach(row => {
      const link = document.createElement("a");
      const img = document.createElement("img");

      // quick hack until i decide how to do proper thumbnails
      let src = row.identifier;
      if(src.startsWith("https://inaturalist-open-data.s3")) {
        src = src.replace("original.", "medium.");
      }
      img.src = src;
      link.href = row.identifier;
      link.appendChild(img);
      multimedia.appendChild(link);
    })
  }

  const uri = "https://www.gbif.org/occurrence/" + o.gbifID;
  clone.querySelector(".gbifLink").href = uri;
  const modal = document.getElementById('overlay');
  while (modal.firstChild) {
    modal.removeChild(modal.firstChild);
  }
  const closer = clone.querySelector('button');
  closer.addEventListener('click', () => {
    modal.close();
  });
  modal.appendChild(clone);
  modal.showModal();
});

const geolocate = new maplibregl.GeolocateControl({
  positionOptions: {
    enableHighAccuracy: true
  },
  trackUserLocation: true,
  showUserHeading: true
});
map.addControl(geolocate, 'top-left');
map.addControl(new maplibregl.GlobeControl(), 'top-left');

let api = "https://api.iroiroiru.jp/here";
if (window.location.hostname === 'localhost') {
  api = "http://localhost:8080/here";
}

geolocate.on('geolocate', (e) => {
  const params = { lat: e.coords.latitude, lon: e.coords.longitude };
  const query = new URLSearchParams(params).toString();

  document.getElementById("spinner").setAttribute("aria-busy", true);

  fetch(`${api}?${query}`).then(response => {
    document.getElementById("spinner").setAttribute("aria-busy", false);
    if (!response.ok) {
      throw new Error("oops: " + response.statusText);
    }
    return response.json();
  }).then(occurrences => {
    const features = occurrences.map((o, i) => {
      const marker = {
        type: 'Feature',
        geometry: o.location,
        properties: {
          icon: icons[o.kingdom] || icons.Fallback,
          raw: JSON.stringify(o),
        }
      }
      return marker;
    });

    map.getSource('occurrences').setData({
      type: 'FeatureCollection',
      features
    });
  });
});
