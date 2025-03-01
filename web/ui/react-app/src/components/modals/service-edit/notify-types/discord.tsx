import {
  FormItem,
  FormItemWithPreview,
  FormLabel,
} from "components/generic/form";

import { NotifyDiscordType } from "types/config";
import NotifyOptions from "./generic";
import { memo } from "react";
import { useGlobalOrDefault } from "./util";

const DISCORD = ({
  name,

  global,
  defaults,
  hard_defaults,
}: {
  name: string;

  global?: NotifyDiscordType;
  defaults?: NotifyDiscordType;
  hard_defaults?: NotifyDiscordType;
}) => (
  <>
    <NotifyOptions
      name={name}
      global={global?.options}
      defaults={defaults?.options}
      hard_defaults={hard_defaults?.options}
    />
    <>
      <FormLabel text="URL Fields" heading />
      <FormItem
        name={`${name}.url_fields.webhookid`}
        required
        label="WebHook ID"
        tooltip={
          <>
            e.g. https://discord.com/api/webhooks/
            <span className="bold-underline">webhook_id</span>
            /token
          </>
        }
        placeholder={useGlobalOrDefault(
          global?.url_fields?.webhookid,
          defaults?.url_fields?.webhookid,
          hard_defaults?.url_fields?.webhookid
        )}
      />
      <FormItem
        name={`${name}.url_fields.token`}
        required
        label="Token"
        tooltip={
          <>
            e.g. https://discord.com/api/webhooks/webhook_id/
            <span className="bold-underline">token</span>
          </>
        }
        placeholder={useGlobalOrDefault(
          global?.url_fields?.token,
          defaults?.url_fields?.token,
          hard_defaults?.url_fields?.token
        )}
        onRight
      />
    </>
    <>
      <FormLabel text="Params" heading />
      <FormItemWithPreview
        name={`${name}.params.avatar`}
        label="Avatar"
        tooltip="Override WebHook avatar with this URL"
        placeholder={
          global?.params?.avatar ||
          defaults?.params?.avatar ||
          hard_defaults?.params?.avatar
        }
      />
      <FormItem
        name={`${name}.params.username`}
        label="Username"
        tooltip="Override the WebHook username"
        placeholder={useGlobalOrDefault(
          global?.params?.username,
          defaults?.params?.username,
          hard_defaults?.params?.username
        )}
      />
      <FormItem
        name={`${name}.params.title`}
        label="Title"
        placeholder={useGlobalOrDefault(
          global?.params?.title,
          defaults?.params?.title,
          hard_defaults?.params?.title
        )}
        onRight
      />
    </>
  </>
);

export default memo(DISCORD);
