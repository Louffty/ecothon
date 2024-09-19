import React, { useEffect, useState } from "react";
import { fetchWithAuth } from "../utils/api";
import OpacitedButton from "./ui/opacitedButton";
import UserEventInfo from "./UserEventInfo";

const AllUserEvents = () => {
  const [events, setEvents] = useState(null);
  const [selectedEvent, setSelectedEvent] = useState(null);

  useEffect(() => {
    const fetchEvent = async () => {
      try {
        const response = await fetchWithAuth(
          `https://nothypeproduction.space/api/v1/event/allMy`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
          }
        );

        const data = await response.json();
        console.log("Data from API:", data);
        setEvents(data.body);
      } catch (error) {
        console.error("Error fetching events:", error);
      }
    };

    if (typeof window !== "undefined") {
      fetchEvent();
    }
  }, []);

  const answerClick = (event) => {
    setSelectedEvent(event);
  };

  if (!events) {
    return <p>Загрузка данных...</p>;
  }

  return (
    <>
      {selectedEvent ? (
        <UserEventInfo event={selectedEvent} />
      ) : (
        events.map((event, index) => (
          <div
            key={index}
            style={{
              backgroundColor: "#f7f5ee",
              width: "94%",
              margin: "auto",
              borderRadius: "15px",
              padding: "1vh",
              display: "flex",
              flexDirection: "column",
              position: "relative",
              marginTop:'2vh'
            }}
          >
            <div style={{ position: "absolute", top: "10px", right: "10px" }}>
              {event?.type || "Тип не указан"}
            </div>

            <div>{event?.title || "Заголовок отсутствует"}</div>
            <div>{event?.description || "Описание отсутствует"}</div>
            <OpacitedButton
              onClick={() => answerClick(event.uuid)}
              title={"Смотреь пользователей"}
            ></OpacitedButton>
          </div>
        ))
      )}
    </>
  );
};

export default AllUserEvents;
