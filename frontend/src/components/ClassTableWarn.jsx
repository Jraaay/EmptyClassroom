import PropTypes from "prop-types";
import dayjs from "dayjs";
import { Alert, Button, Popconfirm } from "antd";
import { CloseOutlined, QuestionCircleOutlined } from "@ant-design/icons";
import "./ClassTableWarn.css";

function ClassTableWarn(props) {
  let hasClassTableData = false;
  if (props.selectedCampus == "海南") {
    hasClassTableData = true;
  }
  if (props.useClassTable) {
    hasClassTableData = true;
  }
  if (!props.selectedDate.isSame(dayjs(), "day")) {
    hasClassTableData = true;
  }
  if (props.todayData.code == 0) {
    if (props.todayData.data.is_fallback != undefined) {
      if (props.todayData.data.is_fallback[props.selectedCampus]) {
        hasClassTableData = true;
      }
    }
  }
  if (!hasClassTableData) {
    return null;
  }
  return props.dontWarnClassTable ? null : (
    <Alert
      description={
        <>
          当前空教室数据包含来自课表的数据，相比教务数据，可能不准确：
          <Button
            size="small"
            onClick={() => {
              window.open(
                "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-UykqdD8nboWEi9xo6jTcbNAunL2"
              );
            }}
            icon={<QuestionCircleOutlined />}
          >
            为什么
          </Button>
        </>
      }
      type="warning"
      showIcon
      action={
        <Popconfirm
          title="永久关闭提醒？"
          description="请确认已经完全了解了【课表】和【教务】数据来源的区别，关闭提醒后将不再显示。"
          onConfirm={() => {
            localStorage.setItem("dontWarnClassTable", "true");
            props.setDontWarnClassTable(true);
          }}
          okText="确定"
          cancelText="取消"
          overlayStyle={{
            maxWidth: 300,
            width: "90%",
          }}
        >
          <Button size="small" icon={<CloseOutlined />} type="text"></Button>
        </Popconfirm>
      }
      className="class-table-warn"
      style={{
        maxWidth: 400,
        width: "90%",
        boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
        textAlign: "left",
      }}
    />
  );
}

ClassTableWarn.propTypes = {
  selectedDate: PropTypes.object,
  selectedCampus: PropTypes.string,
  useClassTable: PropTypes.bool,
  todayData: PropTypes.object,
  dontWarnClassTable: PropTypes.bool,
  setDontWarnClassTable: PropTypes.func,
};

export default ClassTableWarn;
