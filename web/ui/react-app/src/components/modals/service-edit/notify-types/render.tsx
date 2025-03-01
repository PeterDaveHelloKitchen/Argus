import {
  DISCORD,
  GOOGLE_CHAT,
  GOTIFY,
  IFTTT,
  JOIN,
  MATRIX,
  MATTERMOST,
  OPSGENIE,
  PUSHBULLET,
  PUSHOVER,
  ROCKET_CHAT,
  SLACK,
  SMTP,
  TEAMS,
  TELEGRAM,
  ZULIP,
} from "components/modals/service-edit/notify-types";
import { FC, memo } from "react";
import { NotifyType, NotifyTypes, ServiceDict } from "types/config";

interface RenderTypeProps {
  name: string;
  type: string;
  globalNotify?: NotifyType;
  defaults?: ServiceDict<NotifyType>;
  hard_defaults?: NotifyType;
}

const RENDER_TYPE_COMPONENTS: {
  [key in NotifyTypes]: FC<{
    name: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    global: any;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    defaults: any;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    hard_defaults: any;
  }>;
} = {
  discord: DISCORD,
  smtp: SMTP,
  googlechat: GOOGLE_CHAT,
  gotify: GOTIFY,
  ifttt: IFTTT,
  join: JOIN,
  mattermost: MATTERMOST,
  matrix: MATRIX,
  opsgenie: OPSGENIE,
  pushbullet: PUSHBULLET,
  pushover: PUSHOVER,
  rocketchat: ROCKET_CHAT,
  slack: SLACK,
  teams: TEAMS,
  telegram: TELEGRAM,
  zulip: ZULIP,
};

const RenderNotify: FC<RenderTypeProps> = ({
  name,
  type,
  globalNotify,
  defaults,
  hard_defaults,
}) => {
  const RenderTypeComponent =
    RENDER_TYPE_COMPONENTS[(type as NotifyTypes) || "discord"];

  return (
    <RenderTypeComponent
      name={name}
      global={globalNotify}
      defaults={defaults && defaults[type]}
      hard_defaults={hard_defaults && hard_defaults[type]}
    />
  );
};

export default memo(RenderNotify);
