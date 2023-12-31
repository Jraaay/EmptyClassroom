import PropTypes from "prop-types";
import { DatePicker as AntdDatePicker } from "antd";
import "./DatePicker.css";

function DatePicker(props) {
  if (props.todayData.code != 0) {
    return null;
  }
  return (
    <div className="date-picker">
      <AntdDatePicker
        value={props.selectedDate}
        allowClear={false}
        inputReadOnly={true}
        onChange={(date) => {
          props.setSelectedDate(date);
        }}
        disabled={props.todayData.data?.class_table == null}
        disabledDate={(_date) => {
          const date = _date.toDate();
          date.setHours(0, 0, 0, 0);
          if (props.todayData.data?.class_table == null) {
            return true;
          }
          const startWeek = new Date(
            props.todayData.data?.class_table.start_week
          );
          startWeek.setHours(0, 0, 0, 0);
          const endWeek = new Date(props.todayData.data?.class_table.end_week);
          endWeek.setHours(0, 0, 0, 0);
          if (date < startWeek || date > endWeek) {
            return true;
          }
        }}
        popupStyle={{
          position: "absolute",
          left: "calc(50% - 144px)",
        }}
      />
    </div>
  );
}

DatePicker.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  setSelectedDate: PropTypes.func,
};

export default DatePicker;
