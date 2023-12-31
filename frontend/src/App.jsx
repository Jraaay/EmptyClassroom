import "./App.css";
import logo from "./assets/logo.png";
import { Typography, Spin, ConfigProvider, theme } from "antd";
import { useState, useEffect } from "react";
import dayjs from "dayjs";
import Notification from "./components/Notification";
import CachedFetch from "./utils/cache_get";
import CampusButtonGroup from "./components/CampusButtonGroup";
import DatePicker from "./components/DatePicker";
import BuildingPicker from "./components/BuildingPicker";
import ClassTimePicker from "./components/ClassTimePicker";
import EmptyClassroomTable from "./components/EmptyClassroomTable";
import GlobalEmpty from "./components/GlobalEmpty";
import Footer from "./components/Footer";

function App() {
  const [spining, setSpining] = useState(true);
  const [isError, setIsError] = useState(false);
  const [resp, setResp] = useState({ code: 1 });
  const [selectedCampus, setSelectedCampus] = useState("");
  const [selectedDate, setSelectedDate] = useState(dayjs());
  const [selectedBuildings, setSelectedBuildings] = useState([]);
  const [selectedClassTimes, setSelectedClassTimes] = useState([]);
  const [showClassTime, setShowClassTime] = useState(false);
  const [canSelectAllDay, setCanSelectAllDay] = useState(false);
  const [isDark, setIsDark] = useState(false);

  const { Title } = Typography;

  useEffect(() => {
    const mql = window.matchMedia("(prefers-color-scheme: dark)");

    function matchMode(e) {
      const body = document.body;
      if (e.matches) {
        body.classList.add("dark");
        setIsDark(true);
        localStorage.setItem("darkMode", "true");
      } else {
        body.classList.remove("dark");
        setIsDark(false);
        localStorage.setItem("darkMode", "false");
      }
    }

    mql.addEventListener("change", matchMode);

    matchMode(mql);
    CachedFetch("/api/get_data")
      .then((resp) => resp.json())
      .then((resp) => {
        setResp(resp);
        setIsError(false);
        setSpining(false);
      })
      .catch(() => {
        setIsError(true);
        setSpining(false);
      });
    setShowClassTime(localStorage.getItem("showClassTime") !== "false");
    setCanSelectAllDay(localStorage.getItem("canSelectAllDay") === "true");
  }, []);

  return (
    <ConfigProvider
      theme={{
        algorithm:
          localStorage.getItem("darkMode") === "true"
            ? theme.darkAlgorithm
            : theme.defaultAlgorithm,
      }}
    >
      <Spin spinning={spining}>
        <div className="App">
          <img src={logo} className="logo" />
          <Title
            level={3}
            style={{
              marginBottom: "15px",
            }}
          >
            BUPT 空教室查询
          </Title>
          <Notification todayData={resp} />
          <CampusButtonGroup
            todayData={resp}
            selectedCampus={selectedCampus}
            setSelectedCampus={setSelectedCampus}
            setSelectedBuildings={setSelectedBuildings}
            showClassTime={showClassTime}
            setShowClassTime={setShowClassTime}
            canSelectAllDay={canSelectAllDay}
            setCanSelectAllDay={setCanSelectAllDay}
          />
          <DatePicker
            todayData={resp}
            selectedDate={selectedDate}
            setSelectedDate={setSelectedDate}
          />
          <BuildingPicker
            todayData={resp}
            selectedBuildings={selectedBuildings}
            setSelectedBuildings={setSelectedBuildings}
            selectedCampus={selectedCampus}
          />
          <ClassTimePicker
            todayData={resp}
            selectedClassTimes={selectedClassTimes}
            setSelectedClassTimes={setSelectedClassTimes}
            selectedCampus={selectedCampus}
            selectedDate={selectedDate}
            showClassTime={showClassTime}
            canSelectAllDay={canSelectAllDay}
            isDark={isDark}
          />
          <EmptyClassroomTable
            todayData={resp}
            selectedDate={selectedDate}
            selectedCampus={selectedCampus}
            selectedBuildings={selectedBuildings}
            selectedClassTimes={selectedClassTimes}
            setIsError={setIsError}
          />
          <GlobalEmpty todayData={resp} isError={isError} />
          <Footer />
        </div>
      </Spin>
    </ConfigProvider>
  );
}

export default App;
