import PropTypes from "prop-types";
import { notification } from "antd";
import { Button } from "antd";
import { InfoCircleOutlined } from "@ant-design/icons";
import React, { useEffect, useMemo, useCallback } from "react";
import "./Notification.css";

function Notification(props) {
  const [api, contextHolder] = notification.useNotification();

  const ShowNotification = useCallback(
    (notification) => {
      api.info({
        message: notification.title,
        description: (
          <Context.Consumer>
            {() => (
              <div
                dangerouslySetInnerHTML={{
                  __html: notification.content,
                }}
              ></div>
            )}
          </Context.Consumer>
        ),
        duration: notification.duration == 0 ? null : notification.duration,
      });
    },
    [api]
  );
  const Context = React.createContext({ name: "Default" });
  const contextValue = useMemo(
    () => ({ content: props.todayData.data?.content }),
    [props.todayData.data?.content]
  );
  useEffect(() => {
    if (
      props.todayData.code == 0 &&
      props.todayData.data?.notification != undefined
    ) {
      if (props.todayData.data.notification.showNotification) {
        setTimeout(() => {
          ShowNotification(props.todayData.data.notification);
        }, 100);
      }
    }
  }, [
    ShowNotification,
    props.todayData.code,
    props.todayData.data?.notification,
  ]);

  if (props.todayData.code != 0) {
    return null;
  }

  if (props.todayData.code == 0 && props.todayData.data.notification) {
    return (
      <Context.Provider value={contextValue}>
        {contextHolder}
        <Button
          icon={<InfoCircleOutlined />}
          onClick={() => {
            ShowNotification(props.todayData.data.notification);
          }}
          className="notification-button"
        >
          {props.todayData.data.notification.title}
        </Button>
      </Context.Provider>
    );
  }
  return null;
}

Notification.propTypes = {
  todayData: PropTypes.object.isRequired,
};

export default Notification;
