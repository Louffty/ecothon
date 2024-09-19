import React, { useState } from 'react'
import { fetchWithAuth } from '../utils/api';
import OpacitedButton from './ui/opacitedButton';
import Toastify from 'toastify-js'
import "toastify-js/src/toastify.css"
import styles from './styles/Conference.module.css'


const Conference = ({ conference }) => {
  const [conferenceData, setConferenceData] = useState(conference);


  const handleSubmit = async () => {
    Toastify({
      text: "Идет обработка, пожалуйста, ожидайте",
      duration: 3000,
      gravity: "bottom",
      position: "right",
      style: {
        background: "#009605",
        width: "100%",
      },
    }).showToast();
    try {
      const response = await fetchWithAuth(
        "https://nothypeproduction.space/summarize/generate",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            content: conferenceData.url,
          }),
        }
      );
      const responseData = await response.json();
     
      
      document.getElementById('summary-container').innerHTML = responseData.content.replaceAll("<a href", "<br> <a href");

      document.getElementById("summary-button-container").hidden = true;
      document.getElementById("summary-container").hidden = false;
    } catch (error) {
      console.error("Error:", error);
      
    }
  };


  return (
    <>
      <div style={{
        margin: "30px auto",
        marginTop: "3.5rem",
        display: "flex",
        flexDirection: "column",
        alignItems: "center"
      }}>
        <div style={{
          fontFamily: "'Inter', sans-serif",
          fontSize: "22px",
          fontWeight: "bold",
          color: "black"
        }}>{conferenceData.title}</div>
        <div style={{
          fontFamily: "'Inter', sans-serif",
          fontSize: "22px",
          fontWeight: "bold",
          color: "grey"
        }}>{conferenceData.description}</div>
        <div>
      <iframe
        width="720"
        height="405"
        src="https://rutube.ru/play/embed/1164b4268c912eabd400e11387fef321"
        frameBorder="0"
        allow="clipboard-write; autoplay"
        webkitAllowFullScreen
        mozallowfullscreen
        allowFullScreen
      ></iframe>
    
        </div>
        <div id="summary-button-container" style={{ marginTop: "20px" }}>
          <OpacitedButton title={"Краткое содержание"} onClick={handleSubmit}></OpacitedButton>
        </div>
        <div id='summary-container' className={styles.summary_decription} style={{ marginTop: "20px" }} hidden={true}></div>
      </div>
    </>
  );
}

export default Conference
