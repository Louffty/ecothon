"use client";
import styles from "../components/styles/LoginForm.module.css";
import OpacitedButton from "../components/ui/opacitedButton";
import PurpleButton from "../components/ui/purpleButton";
import stylesForInput from "../components/styles/DefaultInput.module.css";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Toastify from "toastify-js";
import "toastify-js/src/toastify.css";
import VkButton from "../components/ui/vkButton";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        "https://nothypeproduction.space/api/v1/user/auth",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ username, password }),
        }
      );
      const data = await response.json();
      if (response.ok) {
        localStorage.setItem("authToken", data.body);
        console.log(data.body);
        Toastify({
          text: "Успешная авторизация",
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
        router.push("/");
      } else {
        console.error(data.error);
        Toastify({
          text: "Проверьте корректность данных",
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
      console.error(error);
    }
  };
  const handleRedirect = () => {
    window.location.href =
      "https://nothypeproduction.space/api/v1/oauth/vk_login";
  };

  return (
    <main>
      <form onSubmit={handleSubmit}>
        <div className={styles.login_form}>
          <div className={styles.welcome_message_box}>
            <p className={styles.welcome_message}>
              Вход в{" "}
              <span className={styles.welcome_message_title}>ЭКОТИПЫ</span>
            </p>
          </div>

          <div className={stylesForInput.input_box}>
            <p className={stylesForInput.input_title}>Имя пользователя</p>
            <input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={stylesForInput.input}
              type="text"
            />
          </div>
          <div className={stylesForInput.input_box}>
            <p className={stylesForInput.input_title}>Пароль</p>
            <input
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={stylesForInput.input}
              type="password"
            />
          </div>

          <div className={styles.button_box}>
            <PurpleButton
              style={{width: "200%"}}
              title={"Войти"}
              type={"submit"}
            />
          </div>
          <div className={styles.button_box}>
            <OpacitedButton
              title={"Создать аккаунт"}
              onClick={() => {
                window.location.href = "/register";
              }}
            />
          </div>

          <div className={styles.button_box}>
            <VkButton
              style={{width: "200%"}}
              title={"Авторизация через VK"}
              onClick={handleRedirect}
            />
          </div>
        </div>
      </form>
    </main>
  );
}
