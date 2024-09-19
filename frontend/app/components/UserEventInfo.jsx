import React, { useEffect, useState } from "react";
import { fetchWithAuth } from "../utils/api";
import styles from './styles/AddAnswer.module.css'
import answerQuestionStyles from "./styles/QuestionPreview.module.css";

const UserEventInfo = ({ event }) => {
  const [data, setData] = useState(null);

  const fetchQuestions = async () => {
    try {
      const response = await fetchWithAuth(
        `https://nothypeproduction.space/api/v1/eventsusers/getallbyevent/${event}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "ngrok-skip-browser-warning": "true",
          },
        }
      );
      const data = await response.json();
      setData(data);
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  };
  useEffect(() => {
    fetchQuestions();
  }, []);
  return (<div>
    <div style={{marginTop:'70px'}}></div>

          
<div className={styles.question_box}>
        <div className={answerQuestionStyles.question_preview_user_box}>
          
          <div
            className={answerQuestionStyles.question_preview_cost}
            style={{
              marginTop: "17px",
              margin: "auto",
              float: "right",
              marginRight: "00px",
            }}
          >
            {data?.event?.title}
            <img
              src="biscuit.png"
              className={answerQuestionStyles.question_preview_cookie}
              alt=""
            />
          </div>
        </div>
        <div className={styles.question_title} style={{ marginTop: "10px" }}>
          {data ? data?.event?.description : "Loading..."}
        </div>
        <div className={styles.question_description}>
          {data ? data?.event?.start_time : "Loading..."}
        </div>
        <div className={styles.question_bottombar}>
        </div>
      </div>

  </div>);
};

export default UserEventInfo;
