import React, { useState, useEffect } from 'react'
import { fetchWithAuth } from '../utils/api';
import DefaultInput from './ui/defaultInput';
import PurpleButton from './ui/purpleButton';
import OpacitedButton from './ui/opacitedButton';
import Toastify from 'toastify-js'
import "toastify-js/src/toastify.css"
const Conference = ({ conference }) => {
  const [conferenceData, setConferenceData] = useState(conference);


  const handleSubmit = async () => {
    console.log(conferenceData.url)
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
     
      
      document.getElementById('summary-container').innerHTML = responseData.content;
    } catch (error) {
      console.error("Error:", error);
      
    }
  };


  return (
    <>
      <div style={{
        width: "300px",
        margin: "30px auto",
        display: "flex",
        flexDirection: "column",
        alignItems: "center"
      }}>
        <div style={{
          fontFamily: "'Inter', sans-serif",
          fontSize: "22px",
          fontWeight: "bold",
          color: "grey"
        }}>{conferenceData.title}</div>
        <div style={{
          fontFamily: "'Inter', sans-serif",
          fontSize: "22px",
          fontWeight: "bold",
          color: "grey"
        }}>{conferenceData.description}</div>
        <div>
          <iframe
            width="300"
            height="215"
            style={{
              borderStyle: "none",
              borderRadius: "16px"
            }}
            src={`https://www.youtube.com/embed/${conferenceData.url.split('/').pop()}`}
            allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
          />



          
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
        <div id='summary-container'></div>
        <div style={{ marginTop: "20px" }}>
          <OpacitedButton title={"Краткое содержание"} onClick={handleSubmit}></OpacitedButton>
        </div>
      </div>
    </>
  );
}

export default Conference