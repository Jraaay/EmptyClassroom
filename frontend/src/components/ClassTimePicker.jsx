import PropTypes from "prop-types";
import { Card, Button } from "antd";
import dayjs from "dayjs";
import "./ClassTimePicker.css";

function ClassTimePicker(props) {
  if (props.todayData.code != 0) {
    return null;
  }

  if (props.selectedCampus == "") {
    return null;
  }

  const classes = [
    "01",
    "02",
    "03",
    "04",
    "05",
    "06",
    "07",
    "08",
    "09",
    "10",
    "11",
    "12",
    "13",
    "14",
  ];

  const class_time = [
    "08:45",
    "09:35",
    "10:35",
    "11:25",
    "12:15",
    "13:45",
    "14:35",
    "15:30",
    "16:25",
    "17:20",
    "18:10",
    "19:15",
    "20:05",
    "20:55",
  ];

  const class_start_time = [
    "08:00",
    "08:50",
    "09:50",
    "10:40",
    "11:30",
    "13:00",
    "13:50",
    "14:45",
    "15:40",
    "16:35",
    "17:25",
    "18:30",
    "19:20",
    "20:10",
  ];

  const options = [];
  const now = new Date();
  const now_hour = now.getHours();
  const now_minute = now.getMinutes();

  for (let i = 0; i <= 13; i++) {
    options.push({
      label: classes[i],
      value: i,
      disabled:
        class_time[i].localeCompare(`${now_hour}:${now_minute}`) < 0 &&
        class_time[class_time.length - 1].localeCompare(
          `${now_hour}:${now_minute}`
        ) >= 0 &&
        !props.canSelectAllDay &&
        props.selectedDate.isSame(dayjs(), "day"),
    });
  }

  for (let i = 0; i <= 13; i++) {
    if (options[i].disabled && props.selectedClassTimes.includes(i)) {
      props.setSelectedClassTimes(
        props.selectedClassTimes.filter((x) => x != i)
      );
    }
  }

  function onCheckAllChange() {
    if (!isAllChecked()) {
      let newSelectedClassTimes = [];
      for (let i = 0; i <= 13; i++) {
        if (options[i].disabled) {
          continue;
        }
        newSelectedClassTimes.push(i);
      }
      props.setSelectedClassTimes(newSelectedClassTimes);
    } else {
      props.setSelectedClassTimes([]);
    }
  }

  function isAllChecked() {
    for (let i = 0; i <= 13; i++) {
      if (options[i].disabled) {
        continue;
      }
      if (!props.selectedClassTimes.includes(i)) {
        return false;
      }
    }
    return true;
  }

  return (
    <Card
      className="class-time-picker"
      style={{
        maxWidth: 400,
        width: "90%",
        boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
      }}
      bodyStyle={{
        maxWidth: "300px",
      }}
    >
      <div
        style={{
          display: "flex",
          flexWrap: "wrap",
          justifyContent: "center",
        }}
      >
        {options.map((x) => (
          <Button
            key={x.value}
            type={
              props.selectedClassTimes.includes(x.value) ? "primary" : "outline"
            }
            onClick={() => {
              if (props.selectedClassTimes.includes(x.value)) {
                props.setSelectedClassTimes(
                  props.selectedClassTimes.filter((y) => y != x.value)
                );
              } else {
                props.setSelectedClassTimes([
                  ...props.selectedClassTimes,
                  x.value,
                ]);
              }
            }}
            style={{
              borderRadius: "0px",
              width: "45px",
              margin: "2px",
              height: props.showClassTime ? "45px" : "30px",
              padding: "0px",
              color: x.disabled ? "#00000073" : null,
            }}
            disabled={x.disabled}
          >
            <div>
              {props.showClassTime ? (
                <div
                  style={{
                    fontSize: "0.7em",
                    marginBottom: "-0.5em",
                  }}
                >
                  {class_start_time[x.label - 1]}
                </div>
              ) : null}
              {x.label}
              {props.showClassTime ? (
                <div
                  style={{
                    fontSize: "0.7em",
                    marginTop: "-0.5em",
                  }}
                >
                  {class_time[x.label - 1]}
                </div>
              ) : null}
            </div>
          </Button>
        ))}
        <Button
          type={isAllChecked() ? "primary" : "outline"}
          onClick={onCheckAllChange}
          style={{
            borderRadius: "0px",
            width: "45px",
            margin: "2px",
            height: props.showClassTime ? "45px" : "30px",
            padding: "0px",
          }}
        >
          {isAllChecked() ? "全不选" : "全选"}
        </Button>
      </div>
    </Card>
  );
}

ClassTimePicker.propTypes = {
  todayData: PropTypes.object,
  selectedClassTimes: PropTypes.array,
  setSelectedClassTimes: PropTypes.func,
  selectedCampus: PropTypes.string,
  selectedDate: PropTypes.object,
  showClassTime: PropTypes.bool,
  canSelectAllDay: PropTypes.bool,
};

export default ClassTimePicker;
