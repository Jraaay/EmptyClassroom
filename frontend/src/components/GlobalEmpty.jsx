import PropTypes from "prop-types";
import { Card, Empty } from "antd";
import "./GlobalEmpty.css";

function GlobalEmpty(props) {
  if (props.todayData.code == 0) {
    return null;
  }

  return (
    <Card
      className="global-empty"
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
          props.isError
            ? "数据获取失败，请刷新重试，若一直失败，可以点击上方按钮反馈"
            : "加载中"
        }
      />
    </Card>
  );
}

GlobalEmpty.propTypes = {
  todayData: PropTypes.object.isRequired,
  isError: PropTypes.bool.isRequired,
};

export default GlobalEmpty;
