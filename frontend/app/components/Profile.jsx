"use client";

import React, { useEffect, useState } from "react";
import { fetchWithAuth } from "../utils/api";
import styles from "./styles/Profile.module.css";
import OpacitedButton from "./ui/opacitedButton";

const Profile = () => {
  const [user, setUser] = useState(null);
  const [isModalOpen, setModalOpen] = useState(false);
  const [organisation, setOrganisation] = useState("");

  const handleVerifyClick = () => {
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  const addOrganization = async () => {
    console.log(organisation)
    try {
      const response = await fetchWithAuth(
        "https://nothypeproduction.space/api/v1/user/verified",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            organisation: organisation,
          }),
        }
      );

      if (response.ok) {
        setModalOpen(false);
      } else {
        console.error("Error submitting question");
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  useEffect(() => {
    if (typeof window !== "undefined") {
      const fetchUser = async () => {
        try {
          const response = await fetchWithAuth(
            `https://nothypeproduction.space/api/v1/user/me`,
            {
              method: "GET",
              headers: {
                "Content-Type": "application/json",
              },
            }
          );
          const data = await response.json();
          setUser(data.body);
          console.log(data);
        } catch (error) {
          console.error(error);
        }
      };
      fetchUser();
    }
  }, []);
  console.log(user);


  function convertTime(timeString) {
    const date = new Date(timeString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hour = String(date.getHours()).padStart(2, '0');
    const minute = String(date.getMinutes()).padStart(2, '0');
    const convertedTime = `${year}-${month}-${day} ${hour}:${minute}`;
    return convertedTime;
  }

  return (
    <>
      <div className={styles.profile_box}>
        <div className={styles.profile_user_data_box}>
          <div style={{ display: "flex" }}>
            <img
              src="profile_kitten.png"
              className={styles.profile_avatar}
              alt=""
            />
            <div className={styles.profile_nickname_and_rank}>
              <p className={styles.profile_nickname}>{user?.username}</p>
              <div style={{ display: "block" }}>
                <img
                  src="https://img.icons8.com/?size=100&id=99885&format=png&color=009605"
                  style={{ marginLeft: "5px", marginTop: "5px" }}
                  width={"15"}
                  height={"15"}
                  alt=""
                />

                <p className={styles.profile_rank}>{user?.rate}</p>
                <div style={{ marginRight: "2%" }}>
                  {user?.is_verified ? (
                    <p style={{color:'green'}}>Верифицировано</p>
                  ) : (
                    <>
                      <button onClick={handleVerifyClick}>Верификация</button>
                      {isModalOpen && (
                        <div className={styles.modal}>
                          <div className={styles.modalContent}>
                            <span className={styles.close} onClick={closeModal}>
                              &times;
                            </span>
                            <p>Пройти верификацию</p>
                            <input
                              value={organisation}
                              onChange={(e) => setOrganisation(e.target.value)}
                            />
                            <OpacitedButton
                              onClick={addOrganization}
                              title={"Отправить"}
                            ></OpacitedButton>
                          </div>
                        </div>
                      )}
                    </>
                  )}
                </div>
              </div>
            </div>
          </div>
          <div className={styles.profile_cookies_count}>
            <img
              src="biscuit.png"
              style={{
                width: "70px",
                height: "70px",
              }}
              alt=""
            />

            <p style={{ marginTop: "15px" }}>{user?.coins_amount}</p>
          </div>
        </div>
        <hr className={styles.profile_hr} />

        <div style={{ textAlign: "center" }}>
          <p
            style={{
              fontSize: "1.2em",
              color: "black",
              fontWeight: "bold",
              marginBottom: "2vh",
            }}
          >
            Запланированные мероприятия
          </p>
        </div>

        {user?.events?.events?.map((event, index) => (
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
              marginTop: '2vh'
            }}
          >
            <div style={{ position: "absolute", top: "10px", right: "10px" }}>
              Оффлайн
            </div>

            <div>{event?.title}</div>
            <div>{event?.description}</div>
            <div>{convertTime(event?.start_time)}</div>
            <div>{event?.address}</div>
          </div>
        ))}
      </div>
    </>
  );
};

export default Profile;
