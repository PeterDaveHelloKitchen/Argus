import { Accordion, Row } from "react-bootstrap";
import { FC, memo } from "react";
import { FormItem, FormSelect } from "components/generic/form";

import { BooleanWithDefault } from "components/generic";
import EditServiceLatestVersionRequire from "./latest-version-require";
import FormURLCommands from "./latest-version-urlcommands";
import { LatestVersionLookupEditType } from "types/service-edit";
import { LatestVersionLookupType } from "types/config";
import VersionWithRefresh from "./version-with-refresh";
import { useWatch } from "react-hook-form";

interface Props {
  serviceName: string;
  original?: LatestVersionLookupEditType;
  defaults?: LatestVersionLookupType;
  hard_defaults?: LatestVersionLookupType;
}

const EditServiceLatestVersion: FC<Props> = ({
  serviceName,
  original,
  defaults,
  hard_defaults,
}) => {
  const latestVersionTypeOptions = [
    { label: "GitHub", value: "github" },
    { label: "URL", value: "url" },
  ];

  const latestVersionType = useWatch({ name: `latest_version.type` });

  return (
    <Accordion>
      <Accordion.Header>Latest Version:</Accordion.Header>
      <Accordion.Body>
        <Row>
          <FormSelect
            name="latest_version.type"
            col_xs={4}
            col_sm={4}
            label="Type"
            options={latestVersionTypeOptions}
          />
          <FormItem
            name="latest_version.url"
            required
            col_sm={8}
            col_xs={8}
            label={latestVersionType === "github" ? "Repo" : "URL"}
            onRight
          />
          {latestVersionType === "github" ? (
            <>
              <FormItem
                key="access_token"
                name="latest_version.access_token"
                col_sm={12}
                label="Access Token"
                placeholder={
                  defaults?.access_token || hard_defaults?.access_token
                }
                isURL={latestVersionType !== "github"}
              />
              <>
                <BooleanWithDefault
                  name="latest_version.use_prerelease"
                  label="Use pre-releases"
                  value={original?.use_prerelease}
                  defaultValue={
                    defaults?.use_prerelease || hard_defaults?.use_prerelease
                  }
                />
              </>
            </>
          ) : (
            <BooleanWithDefault
              name="latest_version.allow_invalid_certs"
              label="Allow Invalid Certs"
              value={original?.allow_invalid_certs}
              defaultValue={
                defaults?.allow_invalid_certs ||
                hard_defaults?.allow_invalid_certs
              }
            />
          )}
          <FormURLCommands />
          <EditServiceLatestVersionRequire />

          <VersionWithRefresh
            vType={0}
            serviceName={serviceName}
            original={original}
          />
        </Row>
      </Accordion.Body>
    </Accordion>
  );
};

export default memo(EditServiceLatestVersion);
