import {
  faCircleNotch,
  faWindowMaximize,
} from "@fortawesome/free-solid-svg-icons";

import { Card } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { ReactElement } from "react";
import { ServiceSummaryType } from "types/summary";
import { faGithub } from "@fortawesome/free-brands-svg-icons";
import { useDelayedRender } from "hooks/delayed-render";

interface ServiceImageData {
  service: ServiceSummaryType;
  visible: boolean;
}

export const ServiceImage = ({
  service,
  visible,
}: ServiceImageData): ReactElement => {
  const delayedRender = useDelayedRender(500);
  return service?.icon ? (
    <Card.Img
      variant="top"
      src={service.icon}
      alt={`${service.id} Image`}
      className="service-image"
      style={{ display: visible ? "inline" : "none" }}
    />
  ) : service?.loading === false ? (
    <FontAwesomeIcon
      icon={service.type === "github" ? faGithub : faWindowMaximize}
      style={{
        minWidth: "fit-content",
        display: visible ? "inline" : "none",
        height: "6rem",
      }}
      className={"service-image"}
    />
  ) : (
    <div
      className="service-image"
      style={{ display: visible ? "inline" : "none" }}
    >
      {delayedRender(() => (
        <FontAwesomeIcon
          icon={faCircleNotch}
          style={{
            minWidth: "fit-content",
            display: visible ? "inline" : "none",
            height: "6rem",
            padding: "0",
          }}
          className={"service-image fa-spin"}
        />
      ))}
    </div>
  );
};
