import styles from "../styles/VkButton.module.css"
export default function VkButton({title, type, onClick}) {
    return <button className={styles.vk_button} onClick={onClick} type={type}>{title}</button>
}
