import PropTypes from "prop-types";
import { Button, Card } from "antd";
import { useState } from "react";
import "./BuildingPicker.css";

function BuildingPicker(props) {
  // const [style, setStyle] = useState(true);
  const [checkedList, setCheckedList] = useState([]);
  if (props.todayData.code != 0) {
    return null;
  }

  if (props.selectedCampus == "") {
    return null;
  }

  const building_name_id_map =
    props.todayData.data.campus_info_map[props.selectedCampus].building_id_map;

  const options = [];
  for (const [key, value] of Object.entries(building_name_id_map)) {
    options.push({
      label: key,
      value: value,
    });
  }
  options.sort((a, b) => {
    if (a.label.length - b.label.length != 0) {
      return a.label.length - b.label.length;
    } else {
      return a.label.localeCompare(b.label);
    }
  });

  return (
    <Card
      style={{
        maxWidth: 400,
        width: "90%",
        boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
      }}
      className="building-picker"
      bodyStyle={{
        maxWidth: "350px",
      }}
    >
      {options.map((item) => (
        <Button
          key={item.value}
          type={checkedList.includes(item.value) ? "primary" : "outline"}
          onClick={() => {
            if (checkedList.includes(item.value)) {
              setCheckedList(checkedList.filter((x) => x != item.value));
              props.setSelectedBuildings(
                props.selectedBuildings.filter((x) => x != item.value)
              );
            } else {
              setCheckedList([...checkedList, item.value]);
              props.setSelectedBuildings([
                ...props.selectedBuildings,
                item.value,
              ]);
            }
          }}
          style={{
            borderRadius: "0px",
            minWidth: "6em",
          }}
        >
          {item.label}
        </Button>
      ))}
    </Card>
  );
}

BuildingPicker.propTypes = {
  todayData: PropTypes.object,
  selectedBuildings: PropTypes.array,
  setSelectedBuildings: PropTypes.func,
  selectedCampus: PropTypes.string,
};

export default BuildingPicker;
