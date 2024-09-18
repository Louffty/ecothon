"use client";

import React, { useEffect, useState } from "react";
import { fetchWithAuth } from "../utils/api";
import styles from "./styles/Profile.module.css";

const Profile = () => {
  const [user, setUser] = useState(null);

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
        
        {user?.events?.events.map((event, index) => (
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
            }}
          >
            <div style={{ position: "absolute", top: "10px", right: "10px" }}>
              type
            </div>

            <div>{event?.title}</div>
            <div>{event?.description}</div>
            <div>time</div>
            <div>{event?.address}</div>
          </div>
        ))}
      </div>
    </>
  );
};

export default Profile;
