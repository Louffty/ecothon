import React, { useState, useRef, useEffect } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  useMapEvents,
  useMap,
} from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import axios from "axios";
import styles from "./styles/Card.module.css";
import OpacitedButton from "./ui/opacitedButton";
import PurpleButton from "./ui/purpleButton";
import DefaultInput from "./ui/defaultInput";
import { fetchWithAuth } from "../utils/api";
import classesStyles from "./styles/Schedule.module.css";
import Toastify from "toastify-js";
import "toastify-js/src/toastify.css";

const coordinatesList = [
  { name: "Москва", coords: [55.7522, 37.6156] },
  { name: "Санкт-Петербург", coords: [59.9343, 30.3351] },
  { name: "Новосибирск", coords: [55.0084, 82.9357] },
];

const MapApp = () => {
  const [markerPosition, setMarkerPosition] = useState([]);
  const [markerData, setMarkerData] = useState([]);
  const [isAddingMarker, setIsAddingMarker] = useState(false);
  const [searchAddress, setSearchAddress] = useState("");
  const [newMarkerData, setNewMarkerData] = useState({
    title: "",
    description: "",
    time: null,
    address: "",
    lat: "",
    lng: "",
  });
  const [center, setCenter] = useState([55.7522, 37.6156]);

  const [selectedOptions, setSelectedOptions] = useState([]);

  const selectorChange = (event) => {
    const options = Array.from(event.target.selectedOptions).map(
      (option) => option.value
    );
    setSelectedOptions(options);
  };

  useEffect(() => {
    const addEvents = async () => {
      try {
        const url = `https://nothypeproduction.space/api/v1/event/all`;
        const response = await fetchWithAuth(url, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });
        const data = await response.json();
        if (Array.isArray(data.body)) {
          setMarkerData(data.body);
        }
      } catch (error) {
        console.error(error);
      }
    };
    addEvents();
  }, []);

  const mapRef = useRef(null);

  const handleSelectChange = (event) => {
    const selectedCoords = coordinatesList[event.target.value].coords;
    setCenter(selectedCoords);
    if (mapRef.current) {
      mapRef.current.setView(selectedCoords);
    }
  };

  const registerToEvent = async (uuid) => {
    console.log(uuid)
    try {
      const response = await fetchWithAuth(
        "https://nothypeproduction.space/api/v1/eventsusers/create",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({eventUUID: uuid })
        }
      );
      
      const responseData = await response.json();
      if (response.ok) {
        console.log("all okay");
      } else {
        console.error("Error");
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const AddMarkerButton = () => {
    const map = useMapEvents({
      click: async (e) => {
        const { lat, lng } = e.latlng;
        const addressData = await getAddressFromCoordinates(lat, lng);
        setMarkerPosition((prevMarkers) => [...prevMarkers, [lat, lng]]);
        setMarkerData((prevData) => [
          ...prevData,
          {
            position: [lat, lng],
            data: {
              title: "",
              description: "",
              time: "",
              address: addressData.display_name,
            },
          },
        ]);
        setNewMarkerData({
          title: "",
          description: "",
          time: "",
          address: addressData.display_name,
          lat: lat,
          lng: lng,
        });
        setIsAddingMarker(true);
      },
    });

    return (
      <div>
        <button onClick={() => map.locate({ setView: true })}>
          Add Marker
        </button>
      </div>
    );
  };

  const getAddressFromCoordinates = async (lat, lng) => {
    Toastify({
      text: "Получение адреса, пожалуйста, ожидайте",
      duration: 3000,
      gravity: "bottom",
      position: "right",
      style: {
        background: "#009605",
        width: "100%",
      },
    }).showToast();
    try {
      const response = await axios.get(
        `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&zoom=18&addressdetails=1`
      );
      return response.data;
    } catch (error) {
      Toastify({
        text: "Произошла ошибка",
        duration: 3000,
        gravity: "bottom",
        position: "right",
        style: {
          background: "#009605",
          width: "100%",
        },
      }).showToast();
      console.error("Error getting address:", error);
      return "Unknown address";
    }
  };

  const handleCloseModal = () => {
    setIsAddingMarker(false);
    setSearchAddress("");
  };

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    const newMarkerInfo = {
      position: [newMarkerData.lat, newMarkerData.lng],
      data: newMarkerData,
    };
    console.log(newMarkerData);
    setMarkerData((prevData) => [...prevData, newMarkerInfo]);
    try {
      const response = await fetchWithAuth(
        "https://nothypeproduction.space/api/v1/event/create",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(newMarkerData),
        }
      );

      const responseData = await response.json();

      Toastify({
        text: "Данные успешно добавлены",
        duration: 3000,
        gravity: "bottom",
        position: "right",
        style: {
          background: "#009605",
          width: "100%",
        },
      }).showToast();
    } catch (error) {
      console.error("Error:", error);
      Toastify({
        text: "Ошибка, проверьте корректность данных",
        duration: 3000,
        gravity: "bottom",
        position: "right",
        style: {
          background: "#009605",
          width: "100%",
        },
      }).showToast();
    }
    handleCloseModal();
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewMarkerData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSearchAddressChange = (e) => {
    setSearchAddress(e.target.value);
  };

  const getCoordinatesFromAddress = async () => {
    try {
      const response = await axios.get(
        `https://nominatim.openstreetmap.org/search?q=${searchAddress}&format=json&limit=1`
      );

      if (response.data.length > 0) {
        const { lat, lon } = response.data[0];
        setNewMarkerData({
          ...newMarkerData,
          lat: parseFloat(lat),
          lng: parseFloat(lon),
          address: response.data[0].display_name,
        });
        setIsAddingMarker(true);
        mapRef.current.flyTo([parseFloat(lat), parseFloat(lon)], 13);
        Toastify({
          text: "Данные отправлены, ожидайте",
          duration: 3000,
          newWindow: true,
          gravity: "bottom",
          position: "right",
          stopOnFocus: true,
          style: {
            background: "#009605",
            width: "100%",
          },
          onClick: function () {},
        }).showToast();
      } else {
        console.error("No results found for the address");
        Toastify({
          text: "Ошибка, попробуйте позже",
          duration: 3000,
          newWindow: true,
          gravity: "bottom",
          position: "right",
          stopOnFocus: true,
          style: {
            background: "#009605",
            width: "100%",
          },
          onClick: function () {},
        }).showToast();
      }
    } catch (error) {
      console.error("Error getting coordinates:", error);
    }
  };

  function convertTime(timeString) {
    const date = new Date(timeString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    const hour = String(date.getHours()).padStart(2, "0");
    const minute = String(date.getMinutes()).padStart(2, "0");
    const convertedTime = `${year}-${month}-${day} ${hour}:${minute}`;
    return convertedTime;
  }
  return (
    <>
      {isAddingMarker && (
        <div className={styles.modalOverlay} onClick={handleCloseModal}>
          <div
            className={styles.classesCreateClassBox}
            style={{ marginTop: "100px" }}
            onClick={(e) => e.stopPropagation()}
          >
            {/* <button
              type="button"
              onClick={handleCloseModal}
              className={styles.closeButton}
            >
              &times;
            </button> */}
            <form onSubmit={handleFormSubmit}>
              <DefaultInput
                type={"text"}
                title={"Название"}
                name={"title"}
                value={newMarkerData.title}
                onChange={handleInputChange}
              />
              <DefaultInput
                type={"text"}
                value={newMarkerData.description}
                title={"Описание"}
                name={"description"}
                onChange={handleInputChange}
              />
              <DefaultInput
                type={"text"}
                value={newMarkerData.address}
                title={"Адресс"}
                name={"address"}
                onChange={handleInputChange}
              />

              <input
                className={classesStyles.classes_create_date}
                style={{ marginLeft: "40px" }}
                name="time"
                type="datetime-local"
                value={newMarkerData.time}
                onChange={handleInputChange}
                required
              />

              <div
                style={{
                  float: "left",
                  marginLeft: "40px",
                  marginTop: "20px",
                }}
              >
                <PurpleButton type={"submit"} title={"Создать"} />
              </div>
            </form>
          </div>
        </div>
      )}
      <div className={styles.map}>
        <select onChange={handleSelectChange}>
          {coordinatesList.map((location, index) => (
            <option key={index} value={index}>
              {location.name}
            </option>
          ))}
        </select>
        <MapContainer
          ref={mapRef}
          center={[55.7522, 37.6156]}
          zoom={13}
          style={{ width: "100%", height: "500px", borderRadius: "16px" }}
          className={styles.map_container}
          attributionControl={false}
        >
          <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
          <AddMarkerButton />
          {markerData.map((marker, index) => (
            <Marker
              key={index}
              position={marker.position}
              draggable={false}
              icon={L.icon({
                iconUrl:
                  "https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png",
              })}
            >
              <Popup>
                <div>
                  <p>Название: {marker?.data?.title}</p>
                  <p>Описание: {marker?.data?.description}</p>
                  <p>Начало: {convertTime(marker?.data?.start_time)}</p>
                  <p>Адрес: {marker?.data?.address}</p>
                  <button onClick={() => registerToEvent(marker?.data?.uuid)}>Зарегистрироваться</button>
                </div>
              </Popup>
            </Marker>
          ))}
          {isAddingMarker && (
            <Marker
              position={[newMarkerData.lat, newMarkerData.lng]}
              draggable={false}
              icon={L.icon({
                iconUrl:
                  "https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png",
              })}
            >
              <Popup>
                <div>
                  <p>Название: {newMarkerData?.title}</p>
                  <p>Описание: {newMarkerData?.description}</p>
                  <p>Начало: {convertTime(newMarkerData?.start_time)}</p>
                  <p>Адрес: {newMarkerData?.address}</p>
                </div>
              </Popup>
            </Marker>
          )}
        </MapContainer>

        <div
          className={styles.map_find_box}
          style={{
            margin: "auto",
            marginTop: "20px",
            width: "100%",
            display: "flex",
            height: "50px",
          }}
        >
          <DefaultInput
            type={"text"}
            title={"Адресс"}
            value={searchAddress}
            onChange={handleSearchAddressChange}
          />
          <div style={{ marginTop: "45px", marginLeft: "15px" }}>
            <OpacitedButton
              title={"Добавить"}
              onClick={getCoordinatesFromAddress}
            />
          </div>
          <div className={styles.list}>
        {markerData.map((marker, index) => (
          <div key={index}>
            <p>Название: {marker?.data?.title}</p>
            <p>Описание: {marker?.data?.description}</p>
            <p>Начало: {convertTime(marker?.data?.start_time)}</p>
            <p>Адрес: {marker?.data?.address}</p>
            <button onClick={() => registerToEvent(marker?.data?.uuid)}>Зарегистрироваться</button>
          </div>
        ))}
      </div>
        </div>
      </div>

     
    </>
  );
};

export default MapApp;
