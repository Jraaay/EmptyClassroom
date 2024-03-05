import PropTypes from "prop-types";
import {
  Radio,
  Button,
  Modal,
  Input,
  message,
  Switch,
  Typography,
  Divider,
} from "antd";
import { useEffect, useState } from "react";
import {
  MessageOutlined,
  SettingOutlined,
  GithubOutlined,
  HeartFilled,
} from "@ant-design/icons";
import "./CampusButtonGroup.css";

function CampusButtonGroup(props) {
  const [campusList, setCampusList] = useState([]);
  const [messageApi, contextHolder] = message.useMessage();
  useEffect(() => {
    if (props.todayData.code == 0) {
      let list = [];
      for (let campus in props.todayData.data.campus_info_map) {
        list.push(campus);
      }
      // 排序，西土城在第一，沙河在第二，其他按照字典序
      const order = ["西土城", "沙河"];
      list.sort((a, b) => {
        if (order.indexOf(a) == -1) {
          if (order.indexOf(b) == -1) {
            return a.localeCompare(b);
          } else {
            return 1;
          }
        } else {
          if (order.indexOf(b) == -1) {
            return -1;
          }
          return order.indexOf(a) - order.indexOf(b);
        }
      });
      setCampusList(list);
      if (props.selectedCampus == "" && list.length > 0) {
        props.setSelectedCampus(list[0]);
      }
    }
  }, [props, props.todayData.code, props.todayData.data?.campus_info_map]);

  const [openReportModal, setOpenReportModal] = useState(false);
  const [openSettingModal, setOpenSettingModal] = useState(false);
  const [textValue, setTextValue] = useState("");
  const [reportModalOkLoading, setReportModalOkLoading] = useState(false);

  function OpenReportModal() {
    setTextValue("");
    setOpenReportModal(true);
  }

  async function ReportModalOk() {
    if (textValue == "") {
      messageApi.error("请输入反馈内容");
      return;
    }
    setReportModalOkLoading(true);
    const resp = await fetch("/api/report", {
      method: "POST",
      body: JSON.stringify({
        text: textValue,
      }),
    });
    setReportModalOkLoading(false);
    setOpenReportModal(false);
    if (resp.status == 200) {
      messageApi.success("提交成功");
    } else {
      Modal.info({
        title: "提交失败",
        content: (
          // 发邮件
          <div>
            <div>提交失败，可以发送邮件至 jray@bupt.edu.cn</div>
          </div>
        ),
        okText: "复制",
        onOk() {
          navigator.clipboard.writeText("jray@bupt.edu.cn");
        },
        cancelText: "取消",
        icon: null,
      });
    }
  }

  function OpenSettingModal() {
    setOpenSettingModal(true);
  }

  return (
    <div className="campus-button-group">
      {contextHolder}
      <Button
        style={{
          marginRight: "10px",
        }}
        icon={<MessageOutlined />}
        onClick={OpenReportModal}
      />
      <Radio.Group
        value={props.selectedCampus}
        onChange={(e) => {
          props.setSelectedCampus(e.target.value);
          props.setSelectedBuildings([]);
        }}
        buttonStyle="solid"
        size="middle"
      >
        {campusList.map((campus) => {
          return (
            <Radio.Button value={campus} key={campus}>
              {campus}
            </Radio.Button>
          );
        })}
      </Radio.Group>
      <Button
        style={{
          marginLeft: "10px",
        }}
        icon={<SettingOutlined />}
        onClick={OpenSettingModal}
      />
      <Modal
        title="反馈提交"
        open={openReportModal}
        onOk={ReportModalOk}
        confirmLoading={reportModalOkLoading}
        onCancel={() => {
          setOpenReportModal(false);
        }}
        okText="提交"
        cancelText="取消"
      >
        <div>
          在反馈咨询前请先看看这个问答：
          <Button
            size="small"
            onClick={() => {
              window.open(
                "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj"
              );
            }}
          >
            空教室查询Q&A
          </Button>
          <Divider />
          如果没能解决你的问题，请在下方输入您的反馈，建议附上联系方式以便回复。
          <Input.TextArea
            rows={3}
            style={{
              marginTop: "10px",
            }}
            onChange={(v) => {
              setTextValue(v.target.value);
            }}
          />
        </div>
      </Modal>
      <Modal
        title="设置"
        open={openSettingModal}
        closable={true}
        footer={null}
        onCancel={() => {
          setOpenSettingModal(false);
        }}
      >
        <div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.showClassTime}
              onChange={(v) => {
                localStorage.setItem("showClassTime", v ? "true" : "false");
                props.setShowClassTime(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              显示课程时间
            </Typography.Title>
          </div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.canSelectAllDay}
              onChange={(v) => {
                localStorage.setItem("canSelectAllDay", v ? "true" : "false");
                props.setCanSelectAllDay(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              全选时选全天
            </Typography.Title>
          </div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.useClassTable}
              onChange={(v) => {
                localStorage.setItem("useClassTable", v ? "true" : "false");
                props.setUseClassTable(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              非必要情况下也使用课表数据
            </Typography.Title>
          </div>
          <Divider plain>
            <HeartFilled />
          </Divider>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            数据来源：
            <Button
              size="small"
              onClick={() => {
                window.open(
                  "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Zip8dx2rlobE5hxW00CcHwOOnre"
                );
              }}
            >
              了解更多
            </Button>
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            问答Q&A：
            <Button
              size="small"
              onClick={() => {
                window.open(
                  "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj"
                );
              }}
            >
              空教室查询Q&A
            </Button>
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            当前数据刷新时间：
            {new Date(props.todayData.data?.update_at).toLocaleString()}
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            项目已开源：
            <Button
              onClick={() =>
                window.open("https://github.com/Jraaay/EmptyClassroom")
              }
              icon={<GithubOutlined />}
              size="small"
            >
              Github
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

CampusButtonGroup.propTypes = {
  todayData: PropTypes.object.isRequired,
  selectedCampus: PropTypes.string,
  setSelectedCampus: PropTypes.func,
  setSelectedBuildings: PropTypes.func,
  showClassTime: PropTypes.bool,
  setShowClassTime: PropTypes.func,
  canSelectAllDay: PropTypes.bool,
  setCanSelectAllDay: PropTypes.func,
  useClassTable: PropTypes.bool,
  setUseClassTable: PropTypes.func,
};

export default CampusButtonGroup;
