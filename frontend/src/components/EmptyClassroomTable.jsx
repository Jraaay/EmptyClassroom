import PropTypes from "prop-types";
import { Empty, Card, Table, Button, Tag, Modal, Descriptions } from "antd";
import { QuestionCircleOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import "./EmptyClassroomTable.css";
import CalculateEmptyClassroom from "../utils/calculte";
import dayjs from "dayjs";

function EmptyClassroomTable(props) {
  console.log(props.selectedBuildings);
  const [emptyClassroom, setEmptyClassroom] = useState([]);
  const [modalTitle, setModalTitle] = useState("");
  const [modalContent, setModalContent] = useState([]);
  const [openModal, setOpenModal] = useState(false);

  useEffect(() => {
    if (
      props.todayData.code == 0 &&
      props.selectedCampus != "" &&
      props.selectedBuildings.length != 0 &&
      props.selectedClassTimes.length != 0
    ) {
      let newEmptyClassroom = {};
      try {
        newEmptyClassroom = CalculateEmptyClassroom(
          props.todayData.data,
          props.selectedCampus,
          props.selectedDate.toDate(),
          props.selectedBuildings,
          props.selectedClassTimes
        );
      } catch (e) {
        props.setIsError(true);
        return;
      }
      setEmptyClassroom(newEmptyClassroom);
    }
  }, [
    props,
    props.selectedBuildings,
    props.selectedCampus,
    props.selectedClassTimes,
    props.selectedDate,
    props.todayData.code,
    props.todayData.data,
  ]);
  if (props.todayData.code != 0) {
    return null;
  }

  if (props.selectedCampus == "") {
    return null;
  }

  if (
    props.selectedBuildings.length == 0 ||
    props.selectedClassTimes.length == 0 ||
    emptyClassroom.length == 0
  ) {
    return (
      <Card
        className="empty-classroom-table"
        style={{
          maxWidth: 400,
          width: "90%",
          boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
        }}
        bodyStyle={{
          maxWidth: "300px",
        }}
      >
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            props.selectedBuildings.length == 0
              ? props.selectedClassTimes.length == 0
                ? "请选择教学楼和上课时间"
                : "请选择教学楼"
              : props.selectedClassTimes.length == 0
              ? "请选择上课时间"
              : "没有空教室了😭"
          }
        />
      </Card>
    );
  }

  function ShowClassroomEmptyInfo(name) {
    let classroomInfo = {};
    for (let i = 0; i < emptyClassroom.length; i++) {
      if (emptyClassroom[i].name == name) {
        classroomInfo = emptyClassroom[i];
        break;
      }
    }
    const emptyTimeList = classroomInfo.empty_class_time;
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
    let emptyTimeListStr = "";
    if (emptyTimeList[0] == 0) {
      emptyTimeListStr += "00:00";
    } else {
      emptyTimeListStr += "00:00-08:00, " + class_time[emptyTimeList[0] - 1];
    }
    for (let i = 1; i < emptyTimeList.length; i++) {
      if (emptyTimeList[i] - emptyTimeList[i - 1] == 1) {
        continue;
      } else {
        emptyTimeListStr +=
          "-" +
          class_start_time[emptyTimeList[i - 1] + 1] +
          ", " +
          class_time[emptyTimeList[i] - 1];
      }
    }
    if (emptyTimeList[emptyTimeList.length - 1] != 13) {
      emptyTimeListStr +=
        "-" + class_start_time[emptyTimeList[emptyTimeList.length - 1] + 1];
      emptyTimeListStr += ", " + class_time[class_time.length - 1] + "-24:00";
    } else {
      emptyTimeListStr += "-24:00";
    }
    const data = [
      {
        key: "座位数",
        value: "100",
      },
      {
        key: "类型",
        value: "多媒体教室",
      },
      {
        key: "空闲时间",
        value: emptyTimeListStr,
      },
      {
        key: "数据来源",
        value: classroomInfo.can_trust ? "教务（可信）" : "课表（参考）",
      },
    ];
    setModalTitle(name);
    setModalContent(data);
    setOpenModal(true);
  }

  const columns = [
    {
      title: "教室",
      key: "name",
      dataIndex: "name",
      align: "center",
      render: (text) => {
        return (
          <span style={{ display: "flex", justifyContent: "center" }}>
            <Button
              size="small"
              onClick={() => {
                ShowClassroomEmptyInfo(text);
              }}
            >
              {text}
            </Button>
          </span>
        );
      },
    },
    {
      title: "座位数",
      key: "size",
      dataIndex: "size",
      align: "center",
    },
    {
      title: "类型",
      key: "type",
      dataIndex: "type",
      align: "center",
    },
    {
      title: (
        <>
          来源
          <Button
            size="small"
            type="text"
            icon={<QuestionCircleOutlined />}
            onClick={() => {
              window.open(
                "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Bwf5dnkr0o3G6ExcedZcLzaSnBe"
              );
            }}
          />
        </>
      ),
      key: "can_trust",
      dataIndex: "can_trust",
      align: "center",
      render: (text) => {
        return text ? (
          <Tag color="green" bordered={false}>
            教务
          </Tag>
        ) : (
          <Tag color="red" bordered={false}>
            课表
          </Tag>
        );
      },
    },
  ];

  return (
    <div className="empty-classroom-table">
      <Card
        style={{
          maxWidth: 400,
          width: "90%",
          boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
        }}
        bodyStyle={{
          padding: "0px",
        }}
      >
        <Table
          dataSource={
            !props.useClassTable &&
            !(props.selectedCampus == "海南") &&
            props.selectedDate.isSame(dayjs(), "day") &&
            (props.todayData.data.is_fallback == undefined ||
              !props.todayData.data.is_fallback[props.selectedCampus])
              ? emptyClassroom.filter((item) => item.can_trust)
              : emptyClassroom
          }
          columns={columns}
          pagination={false}
          bordered={false}
          tableLayout="auto"
          size="small"
          rowKey={(record) => record.name}
          style={{
            width: "100%",
          }}
        />
      </Card>
      <Modal
        title={modalTitle}
        open={openModal}
        footer={null}
        onCancel={() => {
          setOpenModal(false);
        }}
      >
        <div>
          <Descriptions column={1} size="small" layout="vertical">
            {modalContent.map((item, index) => {
              return (
                <Descriptions.Item key={index} label={item.key}>
                  {item.value}
                </Descriptions.Item>
              );
            })}
          </Descriptions>
        </div>
      </Modal>
    </div>
  );
}

EmptyClassroomTable.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
  selectedClassTimes: PropTypes.array,
  setIsError: PropTypes.func,
  useClassTable: PropTypes.bool,
};

export default EmptyClassroomTable;
