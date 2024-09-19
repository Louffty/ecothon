
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import styles from "./styles/Sidebar.module.css";
import { CalendarFold, CircleHelp, MailQuestion, MapPinned, MessageCircleQuestion, SquareArrowLeft, UserRound } from "lucide-react";


const Sidebar = ({ setView }) => {
  const router = useRouter();
  const [authToken, setAuthToken] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("authToken");
    setAuthToken(token);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("authToken");

    router.push("/login");
  };

  return (
    <div className={styles.sidebar_box}>
      <nav className={styles.sidebar_nav}>
        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("profile")}
            className={styles.sidebar_button}
          >
            <UserRound size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Профиль
            </p>
          </button>
        </div>
      

        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("question")}
            className={styles.sidebar_button}
          >
            <MessageCircleQuestion size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Задать вопрос
            </p>
          </button>
        </div>


        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("answer")}
            className={styles.sidebar_button}
          >
            <CircleHelp size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Форум
            </p>
          </button>
        </div>


        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("userquestion")}
            className={styles.sidebar_button}
          >
            <MailQuestion size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Мои посты
            </p>
          </button>
        </div>

        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("shedule")}
            className={styles.sidebar_button}
          >
            <CalendarFold size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Мои мероприятия
            </p>
          </button>
        </div>


        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("card")}
            className={styles.sidebar_button}
          >
            <MapPinned size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Оффлайн
            </p>
          </button>
        </div>

        <div className={styles.sidebar_button_box}>
          <button
            onClick={() => setView("userevent")}
            className={styles.sidebar_button}
          >
           
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Онлайн
            </p>
          </button>
        </div>

        <div className={styles.sidebar_button_box}>
          <button onClick={handleLogout} className={styles.sidebar_button}>
             <SquareArrowLeft size={27} color="#009605" />
            <p
              className={styles.sidebar_button_p}
              style={{ marginLeft: "10px", marginTop: "5px" }}
            >
              Выйти
            </p>
          </button>
        </div>

      </nav>
    </div>
  );
};

export default Sidebar;
