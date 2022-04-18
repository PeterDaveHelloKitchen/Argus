import { ReactElement, useCallback, useState } from "react";

import { Card } from "react-bootstrap";
import { ServiceImage } from "./service-image";
import { ServiceInfo } from "./service-info";
import { ServiceSummaryType } from "types/summary";
import { UpdateInfo } from "./service-update-info";

interface ServiceData {
  service: ServiceSummaryType;
}

export const Service = ({ service }: ServiceData): ReactElement => {
  const [showUpdateInfo, setShowUpdateInfoMain] = useState(false);

  const setShowUpdateInfo = useCallback(() => {
    setShowUpdateInfoMain(!showUpdateInfo);
  }, [showUpdateInfo]);

  // Current and latest version are defined and differ
  const updateAvailable =
    service?.status?.current_version &&
    service?.status?.latest_version &&
    service?.status?.current_version !== service?.status?.latest_version
      ? true
      : false;

  // Latest version has been skipped
  const updateSkipped =
    updateAvailable &&
    service?.status?.approved_version &&
    service.status.approved_version.slice("SKIP_".length) ===
      service?.status?.latest_version
      ? true
      : false;

  return (
    <Card key={service.id} bg="secondary" className={"service shadow"}>
      <Card.Title className="service-title" key={service.id + "-title"}>
        <a className="same-color" href={service.url}>
          <strong>{service.id}</strong>
        </a>
      </Card.Title>

      <Card key={service.id} bg="secondary" className="service-inner">
        <UpdateInfo
          service={service}
          visible={updateAvailable && showUpdateInfo && !updateSkipped}
        />
        <ServiceImage
          service={service}
          visible={!(updateAvailable && showUpdateInfo && !updateSkipped)}
        />
        <ServiceInfo
          service={service}
          setShowUpdateInfo={setShowUpdateInfo}
          updateAvailable={updateAvailable}
          updateSkipped={updateSkipped}
        />
      </Card>
    </Card>
  );
};
