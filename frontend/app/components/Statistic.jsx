import React from "react";
import { Radar, Line, Doughnut } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  RadialLinearScale,
  Filler,
  Tooltip,
  Legend,
  Title,
  ArcElement,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  RadialLinearScale,
  Filler,
  Tooltip,
  Legend,
  Title,
  ArcElement
);

const Statistic = () => {
  const radarData = {
    labels: [">18", "18-25", "25-35", "35-45", "45<"],
    datasets: [
      {
        label: "Age",
        data: [20, 10, 30, 15, 25],
        backgroundColor: "rgba(255, 99, 132, 0.2)",
        borderColor: "rgba(255, 99, 132, 1)",
        borderWidth: 1,
      },
    ],
  };

  const lineData = {
    labels: [">18", "18-25", "25-35", "35-45", "45<"],
    datasets: [
      {
        label: "Age",
        data: [65, 59, 80, 81, 56],
        fill: false,
        backgroundColor: "rgba(54, 162, 235, 1)",
        borderColor: "rgba(54, 162, 235, 1)",
        borderWidth: 2,
      },
    ],
  };

  const doughnutData = {
    labels: ["Male", "Female"],
    datasets: [
      {
        label: "Gender Distribution",
        data: [55, 45],
        backgroundColor: ["rgba(255, 99, 132, 0.6)", "rgba(54, 162, 235, 0.6)"],
        borderColor: ["rgba(255, 99, 132, 1)", "rgba(54, 162, 235, 1)"],
        borderWidth: 1,
      },
    ],
  };
  const lineOptions = {
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  return (
    <div style={{ display: "block", padding: "20px" }}>
      <div style={{ width: "60%" }}>
        <h2>Radar Chart</h2>
        <Radar data={radarData} />
      </div>
      <div style={{ width: "60%" }}>
        <Line data={lineData} options={lineOptions} />
      </div>
      <div style={{ width: "60%" }}>
        <Doughnut data={doughnutData} />
      </div>
    </div>
  );
};

export default Statistic;
